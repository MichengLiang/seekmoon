package source

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/MichengLiang/seekmoon/internal/model"
)

// DefaultMooncakesBaseURL is the default public Mooncakes API origin.
const DefaultMooncakesBaseURL = "https://mooncakes.io"

// MooncakesClient reads public Mooncakes API endpoints.
type MooncakesClient struct {
	BaseURL string
	Fetcher Fetcher
}

// FetchModules reads the public module list.
func (c MooncakesClient) FetchModules(ctx context.Context) model.SourceResult[[]model.ModuleSummary] {
	var raw []moduleItem
	fetch := c.Fetcher.FetchJSON(ctx, c.apiURL("/api/v0/modules"), &raw)
	if fetch.Status != model.StatePresent {
		return sourceResult[[]model.ModuleSummary](string(model.SourceModulesAPI), fetch, nil)
	}
	modules := make([]model.ModuleSummary, 0, len(raw))
	for _, item := range raw {
		if item.Name == "" {
			continue
		}
		modules = append(modules, normalizeModule(item))
	}
	return sourceResult(string(model.SourceModulesAPI), fetch, &modules)
}

// FetchStatistics reads module registry aggregate statistics.
func (c MooncakesClient) FetchStatistics(ctx context.Context) model.SourceResult[model.SnapshotStatistics] {
	var stats model.SnapshotStatistics
	fetch := c.Fetcher.FetchJSON(ctx, c.apiURL("/api/v0/modules/statistics"), &stats)
	if fetch.Status != model.StatePresent {
		return sourceResult[model.SnapshotStatistics](string(model.SourceStatisticsAPI), fetch, nil)
	}
	return sourceResult(string(model.SourceStatisticsAPI), fetch, &stats)
}

// FetchManifest reads and normalizes a module manifest.
func (c MooncakesClient) FetchManifest(ctx context.Context, module string) model.SourceResult[model.ManifestProfile] {
	coord, err := model.ParseModuleCoordinate(module)
	if err != nil {
		fetch := FetchResult{URL: c.apiURL("/api/v0/manifest/" + module), FetchedAt: sourceNow(c.Fetcher.Clock), Status: model.StateFailed, ParseState: model.StateFailed, RawRef: "inline:" + c.apiURL("/api/v0/manifest/"+module), Error: err.Error()}
		return sourceResult[model.ManifestProfile](string(model.SourceManifestAPI), fetch, nil)
	}
	var raw manifestPayload
	fetch := c.Fetcher.FetchJSON(ctx, c.apiURL("/api/v0/manifest/"+url.PathEscape(coord.Owner)+"/"+url.PathEscape(coord.Name)), &raw)
	if fetch.Status != model.StatePresent {
		return sourceResult[model.ManifestProfile](string(model.SourceManifestAPI), fetch, nil)
	}
	profile, err := normalizeManifest(raw)
	if err != nil {
		fetch.Status = model.StateFailed
		fetch.ParseState = model.StateFailed
		fetch.Error = err.Error()
		return sourceResult[model.ManifestProfile](string(model.SourceManifestAPI), fetch, nil)
	}
	return sourceResult(string(model.SourceManifestAPI), fetch, &profile)
}

// FetchRawModules returns the raw modules API payload.
func (c MooncakesClient) FetchRawModules(ctx context.Context) model.SourceResult[any] {
	fetch := c.Fetcher.Fetch(ctx, c.apiURL("/api/v0/modules"))
	return RawJSONSourceResult(string(model.SourceModulesAPI), fetch)
}

// FetchRawManifest returns the raw manifest API payload.
func (c MooncakesClient) FetchRawManifest(ctx context.Context, module string) model.SourceResult[any] {
	coord, err := model.ParseModuleCoordinate(module)
	rawURL := c.apiURL("/api/v0/manifest/" + module)
	if err != nil {
		fetch := FetchResult{URL: rawURL, FetchedAt: sourceNow(c.Fetcher.Clock), Status: model.StateFailed, ParseState: model.StateFailed, RawRef: "inline:" + rawURL, Error: err.Error()}
		return sourceResult[any](string(model.SourceManifestAPI), fetch, nil)
	}
	fetch := c.Fetcher.Fetch(ctx, c.apiURL("/api/v0/manifest/"+url.PathEscape(coord.Owner)+"/"+url.PathEscape(coord.Name)))
	return RawJSONSourceResult(string(model.SourceManifestAPI), fetch)
}

func (c MooncakesClient) apiURL(path string) string {
	base := strings.TrimRight(c.BaseURL, "/")
	if base == "" {
		base = DefaultMooncakesBaseURL
	}
	return base + path
}

type moduleItem struct {
	Name        string         `json:"name"`
	Version     string         `json:"version"`
	Description string         `json:"description"`
	Keywords    []string       `json:"keywords"`
	Repository  string         `json:"repository"`
	License     string         `json:"license"`
	IsNew       bool           `json:"is_new"`
	CreatedAt   string         `json:"created_at"`
	Raw         map[string]any `json:"-"`
}

func (m *moduleItem) UnmarshalJSON(data []byte) error {
	type alias moduleItem
	var raw map[string]any
	if err := jsonUnmarshal(data, &raw); err != nil {
		return err
	}
	var item alias
	if err := jsonUnmarshal(data, &item); err != nil {
		return err
	}
	*(*alias)(m) = item
	m.Raw = raw
	return nil
}

func normalizeModule(item moduleItem) model.ModuleSummary {
	return model.ModuleSummary{
		Module:      item.Name,
		Version:     item.Version,
		Description: evidenceString(item.Description, string(model.SourceModulesAPI)),
		Keywords:    evidenceStrings(item.Keywords, string(model.SourceModulesAPI)),
		Repository:  evidenceString(item.Repository, string(model.SourceModulesAPI)),
		License:     evidenceString(item.License, string(model.SourceModulesAPI)),
		IsNew:       item.IsNew,
		CreatedAt:   item.CreatedAt,
		Raw:         item.Raw,
	}
}

type manifestPayload struct {
	Module        string                  `json:"module"`
	Name          string                  `json:"name"`
	Version       string                  `json:"version"`
	LatestVersion string                  `json:"latest_version"`
	Downloads     int                     `json:"downloads"`
	HasPackage    bool                    `json:"has_package"`
	BuildStatus   *string                 `json:"build_status"`
	Metadata      map[string]any          `json:"metadata"`
	Versions      []model.ManifestVersion `json:"versions"`
}

func normalizeManifest(raw manifestPayload) (model.ManifestProfile, error) {
	module := raw.Module
	if module == "" {
		module = raw.Name
	}
	if raw.Module != "" && raw.Name != "" && raw.Module != raw.Name {
		return model.ManifestProfile{}, fmt.Errorf("manifest module %q does not match name %q", raw.Module, raw.Name)
	}
	versionCount := len(raw.Versions)
	return model.ManifestProfile{
		Module:        module,
		Version:       raw.Version,
		LatestVersion: raw.LatestVersion,
		Downloads:     raw.Downloads,
		HasPackage:    raw.HasPackage,
		BuildStatus:   evidenceOptionalString(raw.BuildStatus, string(model.SourceManifestAPI)),
		Metadata:      normalizeMetadata(raw.Metadata),
		Versions:      raw.Versions,
		VersionsCount: model.Derived(versionCount, string(model.SourceDerived)),
		DocsURL:       model.Derived("https://mooncakes.io/docs/"+module, string(model.SourceDerived)),
	}, nil
}

func normalizeMetadata(raw map[string]any) model.ManifestMetadata {
	if raw == nil {
		raw = map[string]any{}
	}
	return model.ManifestMetadata{
		Description:      evidenceString(stringFromRaw(raw, "description"), string(model.SourceManifestAPI)),
		Keywords:         evidenceStrings(stringsFromRaw(raw, "keywords"), string(model.SourceManifestAPI)),
		Repository:       evidenceString(stringFromRaw(raw, "repository"), string(model.SourceManifestAPI)),
		License:          evidenceString(stringFromRaw(raw, "license"), string(model.SourceManifestAPI)),
		Checksum:         evidenceString(stringFromRaw(raw, "checksum"), string(model.SourceManifestAPI)),
		Deps:             stringMapFromRaw(raw, "deps"),
		PreferredTarget:  evidenceString(firstStringFromRaw(raw, "preferred-target", "preferred_target", "preferred-backend"), string(model.SourceManifestAPI)),
		SupportedTargets: evidenceStrings(firstStringsFromRaw(raw, "supported-targets", "supported_targets", "targets"), string(model.SourceManifestAPI)),
		Raw:              raw,
	}
}
