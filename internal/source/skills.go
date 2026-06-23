package source

import (
	"context"
	"net/url"
	"strings"

	"github.com/yumiaura/seekmoon/internal/model"
)

type SkillsClient struct {
	BaseURL string
	Fetcher Fetcher
}

func (c SkillsClient) FetchSkills(ctx context.Context) model.SourceResult[[]model.SkillEntry] {
	var raw []skillPayload
	fetch := c.Fetcher.FetchJSON(ctx, c.apiURL("/api/v0/skills"), &raw)
	if fetch.Status != model.StatePresent {
		return SourceResult[[]model.SkillEntry](string(model.SourceSkillsAPI), fetch, nil)
	}
	entries := make([]model.SkillEntry, 0, len(raw))
	for _, item := range raw {
		entries = append(entries, normalizeSkill(item))
	}
	return SourceResult(string(model.SourceSkillsAPI), fetch, &entries)
}

func (c SkillsClient) FetchSkill(ctx context.Context, entry string) model.SourceResult[model.SkillEntry] {
	var raw skillPayload
	fetch := c.Fetcher.FetchJSON(ctx, c.apiURL("/api/v0/skills/"+url.PathEscape(entry)), &raw)
	if fetch.Status != model.StatePresent {
		return SourceResult[model.SkillEntry](string(model.SourceSkillsAPI), fetch, nil)
	}
	value := normalizeSkill(raw)
	return SourceResult(string(model.SourceSkillsAPI), fetch, &value)
}

func (c SkillsClient) FetchSkillAsset(ctx context.Context, assetURL string) model.EvidenceObject {
	fetch := c.Fetcher.Fetch(ctx, assetURL)
	if fetch.Status == model.StatePresent {
		return model.Present(map[string]any{"url": fetch.URL, "bytes": len(fetch.Body)}, string(model.SourceSkillsAPI))
	}
	if fetch.Status == model.StateUnavailable {
		return model.Unavailable[map[string]any](string(model.SourceSkillsAPI))
	}
	return model.Failed[map[string]any](string(model.SourceSkillsAPI), fetch.Error)
}

func (c SkillsClient) apiURL(path string) string {
	base := strings.TrimRight(c.BaseURL, "/")
	if base == "" {
		base = DefaultMooncakesBaseURL
	}
	return base + path
}

type skillPayload struct {
	Module       string         `json:"module"`
	Author       string         `json:"author"`
	AuthorAvatar string         `json:"author_avatar"`
	Version      string         `json:"version"`
	Package      string         `json:"package"`
	Name         string         `json:"name"`
	DetailURL    string         `json:"detail_url"`
	WasmURL      string         `json:"wasm_url"`
	ChecksumURL  string         `json:"checksum_url"`
	Metadata     map[string]any `json:"metadata"`
	Repository   string         `json:"repository"`
	CreatedAt    string         `json:"created_at"`
}

func normalizeSkill(raw skillPayload) model.SkillEntry {
	return model.SkillEntry{
		Module:       raw.Module,
		Author:       raw.Author,
		AuthorAvatar: evidenceString(raw.AuthorAvatar, string(model.SourceSkillsAPI)),
		Version:      raw.Version,
		Package:      raw.Package,
		Name:         raw.Name,
		DetailURL:    raw.DetailURL,
		WasmURL:      raw.WasmURL,
		ChecksumURL:  raw.ChecksumURL,
		Metadata:     raw.Metadata,
		Repository:   evidenceString(raw.Repository, string(model.SourceSkillsAPI)),
		CreatedAt:    raw.CreatedAt,
	}
}
