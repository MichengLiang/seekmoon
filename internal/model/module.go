package model

import (
	"fmt"
	"strings"
)

type ModuleCoordinate struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

func ParseModuleCoordinate(value string) (ModuleCoordinate, error) {
	owner, name, ok := strings.Cut(value, "/")
	if !ok || owner == "" || name == "" || strings.Contains(name, "/") {
		return ModuleCoordinate{}, fmt.Errorf("module coordinate must be owner/module: %q", value)
	}
	return ModuleCoordinate{Owner: owner, Name: name}, nil
}

func (c ModuleCoordinate) String() string {
	if c.Owner == "" || c.Name == "" {
		return ""
	}
	return c.Owner + "/" + c.Name
}

type ModuleSummary struct {
	Module      string              `json:"module"`
	Version     string              `json:"version"`
	Description EvidenceString      `json:"description"`
	Keywords    EvidenceStringArray `json:"keywords"`
	Repository  EvidenceString      `json:"repository"`
	License     EvidenceString      `json:"license"`
	IsNew       bool                `json:"is_new"`
	CreatedAt   string              `json:"created_at"`
	Raw         map[string]any      `json:"raw,omitempty"`
}

type SearchResult struct {
	Rank        int                 `json:"rank"`
	Module      string              `json:"module"`
	Version     string              `json:"version"`
	Description EvidenceString      `json:"description"`
	Keywords    EvidenceStringArray `json:"keywords"`
	License     EvidenceString      `json:"license"`
	Repository  EvidenceString      `json:"repository"`
	Downloads   EvidenceInt         `json:"downloads,omitempty"`
	BuildStatus EvidenceString      `json:"build_status,omitempty"`
	Target      EvidenceObject      `json:"target,omitempty"`
	Match       MatchEvidence       `json:"match"`
	SnapshotID  string              `json:"snapshot_id"`
}

type MatchEvidence struct {
	Fields []string           `json:"fields,omitempty"`
	Token  string             `json:"token,omitempty"`
	Score  *Evidence[float64] `json:"score,omitempty"`
}

type SearchQuery struct {
	Text   string `json:"text"`
	Kind   string `json:"kind"`
	Target string `json:"target,omitempty"`
}

type SearchOutput struct {
	Schema   string         `json:"schema"`
	Snapshot SnapshotRef    `json:"snapshot"`
	Query    SearchQuery    `json:"query"`
	Results  []SearchResult `json:"results"`
}
