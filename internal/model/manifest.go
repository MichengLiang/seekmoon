package model

type ManifestProfile struct {
	Module        string            `json:"module"`
	Version       string            `json:"version"`
	LatestVersion string            `json:"latest_version"`
	Downloads     int               `json:"downloads"`
	HasPackage    bool              `json:"has_package"`
	BuildStatus   EvidenceString    `json:"build_status"`
	Metadata      ManifestMetadata  `json:"metadata"`
	Versions      []ManifestVersion `json:"versions"`
	VersionsCount EvidenceInt       `json:"versions_count"`
	DocsURL       EvidenceString    `json:"docs_url"`
}

type ManifestMetadata struct {
	Description      EvidenceString      `json:"description,omitempty"`
	Keywords         EvidenceStringArray `json:"keywords,omitempty"`
	Repository       EvidenceString      `json:"repository,omitempty"`
	License          EvidenceString      `json:"license,omitempty"`
	Checksum         EvidenceString      `json:"checksum,omitempty"`
	Deps             map[string]string   `json:"deps,omitempty"`
	PreferredTarget  EvidenceString      `json:"preferred_target,omitempty"`
	SupportedTargets EvidenceStringArray `json:"supported_targets,omitempty"`
	Raw              map[string]any      `json:"raw,omitempty"`
}

type ManifestVersion struct {
	Version string `json:"version"`
	Yanked  bool   `json:"yanked,omitempty"`
}
