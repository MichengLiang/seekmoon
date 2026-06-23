package source

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/yumiaura/seekmoon/internal/model"
)

type AssetClient struct {
	BaseURL     string
	DownloadURL string
	Fetcher     Fetcher
}

func (c AssetClient) FetchModuleIndex(ctx context.Context, module, version string) model.SourceResult[model.ModuleIndexTree] {
	var raw moduleIndexNode
	fetch := c.Fetcher.FetchJSON(ctx, c.assetURL(module, version, "module_index.json"), &raw)
	if fetch.Status != model.StatePresent {
		return SourceResult[model.ModuleIndexTree](string(model.SourceModuleIndex), fetch, nil)
	}
	tree, err := normalizeModuleIndex(module, raw)
	if err != nil {
		fetch.Status = model.StateFailed
		fetch.ParseState = model.StateFailed
		fetch.Error = err.Error()
		return SourceResult[model.ModuleIndexTree](string(model.SourceModuleIndex), fetch, nil)
	}
	return SourceResult(string(model.SourceModuleIndex), fetch, &tree)
}

func (c AssetClient) FetchPackageData(ctx context.Context, module, version, packagePath string) model.SourceResult[model.PackageData] {
	relpath, err := model.PackageRelPath(module, packagePath)
	if err != nil {
		fetch := FetchResult{URL: c.assetURL(module, version, packagePath), FetchedAt: sourceNow(c.Fetcher.Clock), Status: model.StateFailed, ParseState: model.StateFailed, RawRef: "inline:" + c.assetURL(module, version, packagePath), Error: err.Error()}
		return SourceResult[model.PackageData](string(model.SourcePackageData), fetch, nil)
	}
	var raw packageDataPayload
	fetch := c.Fetcher.FetchJSON(ctx, c.assetURL(module, version, relpath, "package_data.json"), &raw)
	if fetch.Status != model.StatePresent {
		return SourceResult[model.PackageData](string(model.SourcePackageData), fetch, nil)
	}
	data := normalizePackageData(raw)
	return SourceResult(string(model.SourcePackageData), fetch, &data)
}

func (c AssetClient) FetchResource(ctx context.Context, module, version, packagePath string) model.SourceResult[map[string]any] {
	relpath, err := model.PackageRelPath(module, packagePath)
	if err != nil {
		fetch := FetchResult{URL: c.assetURL(module, version, packagePath), FetchedAt: sourceNow(c.Fetcher.Clock), Status: model.StateFailed, ParseState: model.StateFailed, RawRef: "inline:" + c.assetURL(module, version, packagePath), Error: err.Error()}
		return SourceResult[map[string]any](string(model.SourceResourceAsset), fetch, nil)
	}
	var raw map[string]any
	fetch := c.Fetcher.FetchJSON(ctx, c.assetURL(module, version, relpath, "resource.json"), &raw)
	if fetch.Status == model.StateUnavailable {
		return SourceResult[map[string]any](string(model.SourceResourceAsset), fetch, nil)
	}
	if fetch.Status != model.StatePresent {
		return SourceResult[map[string]any](string(model.SourceResourceAsset), fetch, nil)
	}
	return SourceResult(string(model.SourceResourceAsset), fetch, &raw)
}

func (c AssetClient) FetchSourceZipAttempt(ctx context.Context, module, version string) model.SourceAttempt {
	fetch := c.Fetcher.Fetch(ctx, c.sourceZipURL(module, version))
	attempt := model.SourceAttempt{Status: fetch.Status, URL: fetch.URL, Error: fetch.Error}
	if fetch.Status != model.StatePresent {
		return attempt
	}
	summary, err := SummarizeZip(fetch.Body)
	if err != nil {
		attempt.Status = model.StateFailed
		attempt.Error = err.Error()
		return attempt
	}
	attempt.Path = filesSummaryString(summary)
	return attempt
}

func (c AssetClient) assetURL(module, version string, parts ...string) string {
	base := strings.TrimRight(c.BaseURL, "/")
	if base == "" {
		base = DefaultMooncakesBaseURL
	}
	coord, _ := model.ParseModuleCoordinate(module)
	pathParts := []string{"assets", coord.Owner, coord.Name + "@" + version}
	pathParts = append(pathParts, parts...)
	escaped := make([]string, 0, len(pathParts))
	for _, part := range pathParts {
		if part == "" {
			continue
		}
		for _, segment := range strings.Split(part, "/") {
			if segment == "" {
				continue
			}
			escaped = append(escaped, url.PathEscape(segment))
		}
	}
	return base + "/" + strings.Join(escaped, "/")
}

func (c AssetClient) sourceZipURL(module, version string) string {
	base := strings.TrimRight(c.DownloadURL, "/")
	if base == "" {
		base = "https://download.mooncakes.io"
	}
	coord, _ := model.ParseModuleCoordinate(module)
	return fmt.Sprintf("%s/user/%s/%s/%s.zip", base, url.PathEscape(coord.Owner), url.PathEscape(coord.Name), url.PathEscape(version))
}

type moduleIndexNode struct {
	Name     string             `json:"name"`
	Package  *packageSummaryRaw `json:"package"`
	Childs   []moduleIndexNode  `json:"childs"`
	Children []moduleIndexNode  `json:"children"`
}

type packageSummaryRaw struct {
	Path      string              `json:"path"`
	Traits    []model.APISummary  `json:"traits"`
	Errors    []model.APISummary  `json:"errors"`
	Types     []model.TypeSummary `json:"types"`
	TypeAlias []model.APISummary  `json:"typealias"`
	Values    []model.APISummary  `json:"values"`
	Misc      []model.APISummary  `json:"misc"`
}

func normalizeModuleIndex(module string, raw moduleIndexNode) (model.ModuleIndexTree, error) {
	var pkg *model.PackageSummary
	if raw.Package != nil {
		relpath, err := model.PackageRelPath(module, raw.Package.Path)
		if err != nil {
			return model.ModuleIndexTree{}, err
		}
		pkg = &model.PackageSummary{
			Path:      raw.Package.Path,
			RelPath:   model.Derived(relpath, string(model.SourceDerived)),
			Traits:    raw.Package.Traits,
			Errors:    raw.Package.Errors,
			Types:     raw.Package.Types,
			TypeAlias: raw.Package.TypeAlias,
			Values:    raw.Package.Values,
			Misc:      raw.Package.Misc,
		}
	}
	children := raw.Childs
	if len(children) == 0 {
		children = raw.Children
	}
	out := model.ModuleIndexTree{Name: raw.Name, Package: pkg, Childs: make([]model.ModuleIndexTree, 0, len(children))}
	for _, child := range children {
		normalized, err := normalizeModuleIndex(module, child)
		if err != nil {
			return model.ModuleIndexTree{}, err
		}
		out.Childs = append(out.Childs, normalized)
	}
	return out, nil
}

type packageDataPayload struct {
	Name      string        `json:"name"`
	Traits    []apiEntryRaw `json:"traits"`
	Errors    []apiEntryRaw `json:"errors"`
	Types     []apiEntryRaw `json:"types"`
	TypeAlias []apiEntryRaw `json:"typealias"`
	Values    []apiEntryRaw `json:"values"`
	Misc      []apiEntryRaw `json:"misc"`
}

type apiEntryRaw struct {
	Name      string         `json:"name"`
	Docstring string         `json:"docstring"`
	Signature string         `json:"signature"`
	Loc       map[string]any `json:"loc"`
	Methods   []apiEntryRaw  `json:"methods"`
	Impls     []apiEntryRaw  `json:"impls"`
}

func normalizePackageData(raw packageDataPayload) model.PackageData {
	return model.PackageData{
		Name:      raw.Name,
		Traits:    normalizeAPIEntries(raw.Traits),
		Errors:    normalizeAPIEntries(raw.Errors),
		Types:     normalizeAPIEntries(raw.Types),
		TypeAlias: normalizeAPIEntries(raw.TypeAlias),
		Values:    normalizeAPIEntries(raw.Values),
		Misc:      normalizeAPIEntries(raw.Misc),
	}
}

func normalizeAPIEntries(raw []apiEntryRaw) []model.APIEntry {
	out := make([]model.APIEntry, 0, len(raw))
	for _, item := range raw {
		loc := model.Missing[map[string]any](string(model.SourcePackageData))
		if item.Loc != nil {
			loc = model.Present(item.Loc, string(model.SourcePackageData))
		}
		out = append(out, model.APIEntry{
			Name:           item.Name,
			Docstring:      evidenceString(item.Docstring, string(model.SourcePackageData)),
			Signature:      item.Signature,
			PlainSignature: model.Derived(PlainSignature(item.Signature), string(model.SourceDerived)),
			Loc:            loc,
			Methods:        normalizeAPIEntries(item.Methods),
			Impls:          normalizeAPIEntries(item.Impls),
		})
	}
	return out
}

var htmlTagRE = regexp.MustCompile(`<[^>]*>`)

func PlainSignature(signature string) string {
	plain := htmlTagRE.ReplaceAllString(signature, "")
	plain = strings.ReplaceAll(plain, "&lt;", "<")
	plain = strings.ReplaceAll(plain, "&gt;", ">")
	plain = strings.ReplaceAll(plain, "&amp;", "&")
	return strings.TrimSpace(plain)
}

func SummarizeZip(data []byte) (model.FilesSummary, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return model.FilesSummary{}, err
	}
	var summary model.FilesSummary
	for _, file := range reader.File {
		name := file.Name
		switch {
		case strings.HasSuffix(name, "moon.mod") || strings.HasSuffix(name, "moon.mod.json"):
			summary.MoonMod = true
		case strings.Contains(strings.ToLower(name), "readme"):
			summary.Readme = true
		case strings.Contains(strings.ToLower(name), "license"):
			summary.License = true
		case strings.Contains(name, "test"):
			summary.Tests++
		case strings.Contains(name, "example"):
			summary.Examples++
		case strings.Contains(name, "bench"):
			summary.Benches++
		case strings.HasSuffix(name, ".mbt"):
			summary.Sources++
		}
	}
	return summary, nil
}

func filesSummaryString(s model.FilesSummary) string {
	data, _ := json.Marshal(s)
	return string(data)
}

func StatusFromHTTP(statusCode int) model.State {
	if statusCode == http.StatusNotFound {
		return model.StateUnavailable
	}
	if statusCode >= 200 && statusCode < 300 {
		return model.StatePresent
	}
	return model.StateFailed
}
