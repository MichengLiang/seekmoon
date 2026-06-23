package model

import (
	"fmt"
	"time"
)

// AdoptionConclusion is the closed recommendation vocabulary for records.
type AdoptionConclusion string

// Adoption conclusions supported by records and schemas.
const (
	ConclusionAdopt                AdoptionConclusion = "adopt"
	ConclusionAdoptWithAdapter     AdoptionConclusion = "adopt-with-adapter"
	ConclusionContinueVerification AdoptionConclusion = "continue-verification"
	ConclusionContributeUpstream   AdoptionConclusion = "contribute-upstream"
	ConclusionFork                 AdoptionConclusion = "fork"
	ConclusionBuildOwn             AdoptionConclusion = "build-own"
	ConclusionRejectForNow         AdoptionConclusion = "reject-for-now"
)

// ParseAdoptionConclusion validates an adoption conclusion string.
func ParseAdoptionConclusion(value string) (AdoptionConclusion, error) {
	conclusion := AdoptionConclusion(value)
	if conclusion.IsValid() {
		return conclusion, nil
	}
	return "", fmt.Errorf("unknown adoption conclusion %q", value)
}

// IsValid reports whether the conclusion belongs to the closed vocabulary.
func (c AdoptionConclusion) IsValid() bool {
	switch c {
	case ConclusionAdopt, ConclusionAdoptWithAdapter, ConclusionContinueVerification, ConclusionContributeUpstream, ConclusionFork, ConclusionBuildOwn, ConclusionRejectForNow:
		return true
	default:
		return false
	}
}

// AdoptionRecord records the dependency adoption decision and evidence refs.
type AdoptionRecord struct {
	RecordID     string             `json:"record_id"`
	CreatedAt    time.Time          `json:"created_at"`
	Project      ProjectIdentity    `json:"project"`
	SnapshotID   string             `json:"snapshot_id"`
	Candidate    CandidateRef       `json:"candidate"`
	Version      string             `json:"version"`
	Conclusion   AdoptionConclusion `json:"conclusion"`
	EvidenceRefs []EvidenceRef      `json:"evidence_refs"`
	NotConfirmed []string           `json:"not_confirmed,omitempty"`
	Note         string             `json:"note,omitempty"`
}

// CandidateRef identifies a library or skill candidate.
type CandidateRef struct {
	Kind    string `json:"kind"`
	Module  string `json:"module"`
	Version string `json:"version"`
	Package string `json:"package,omitempty"`
	Name    string `json:"name,omitempty"`
}

// EvidenceRef points to source evidence used by a record or report.
type EvidenceRef struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
	Path string `json:"path,omitempty"`
}
