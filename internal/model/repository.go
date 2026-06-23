package model

type RepositorySignal struct {
	URL           string         `json:"url"`
	Status        State          `json:"status"`
	IsArchived    Evidence[bool] `json:"is_archived"`
	PushedAt      EvidenceString `json:"pushed_at"`
	License       EvidenceString `json:"license"`
	HasReleases   Evidence[bool] `json:"has_releases"`
	OpenIssues    EvidenceInt    `json:"open_issues"`
	OpenPulls     EvidenceInt    `json:"open_pulls"`
	HasWorkflows  Evidence[bool] `json:"has_workflows"`
	HasReadme     Evidence[bool] `json:"has_readme"`
	HasTests      Evidence[bool] `json:"has_tests"`
	HasExamples   Evidence[bool] `json:"has_examples"`
	DefaultBranch EvidenceString `json:"default_branch"`
	Error         string         `json:"error,omitempty"`
}
