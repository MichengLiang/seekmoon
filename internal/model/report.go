package model

// Report is the generated adoption evidence report.
type Report struct {
	Goal        string            `json:"goal"`
	Environment ReportEnvironment `json:"environment"`
	DataSources []string          `json:"data_sources"`
	Candidates  []CandidateRef    `json:"candidates"`
	Inspection  []EvidenceRef     `json:"inspection,omitempty"`
	Validation  []EvidenceRef     `json:"validation,omitempty"`
	Decision    AdoptionRecord    `json:"decision"`
}

// ReportEnvironment records the environment context used by a report.
type ReportEnvironment struct {
	Toolchain string          `json:"toolchain,omitempty"`
	Snapshot  SnapshotRef     `json:"snapshot"`
	Project   ProjectIdentity `json:"project"`
}
