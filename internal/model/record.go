package model

import (
	"fmt"
	"time"
)

type AdoptionConclusion string

const (
	ConclusionAdopt                AdoptionConclusion = "adopt"
	ConclusionAdoptWithAdapter     AdoptionConclusion = "adopt-with-adapter"
	ConclusionContinueVerification AdoptionConclusion = "continue-verification"
	ConclusionContributeUpstream   AdoptionConclusion = "contribute-upstream"
	ConclusionFork                 AdoptionConclusion = "fork"
	ConclusionBuildOwn             AdoptionConclusion = "build-own"
	ConclusionRejectForNow         AdoptionConclusion = "reject-for-now"
)

func ParseAdoptionConclusion(value string) (AdoptionConclusion, error) {
	conclusion := AdoptionConclusion(value)
	if conclusion.IsValid() {
		return conclusion, nil
	}
	return "", fmt.Errorf("unknown adoption conclusion %q", value)
}

func (c AdoptionConclusion) IsValid() bool {
	switch c {
	case ConclusionAdopt, ConclusionAdoptWithAdapter, ConclusionContinueVerification, ConclusionContributeUpstream, ConclusionFork, ConclusionBuildOwn, ConclusionRejectForNow:
		return true
	default:
		return false
	}
}

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

type CandidateRef struct {
	Kind    string `json:"kind"`
	Module  string `json:"module"`
	Version string `json:"version"`
	Package string `json:"package,omitempty"`
	Name    string `json:"name,omitempty"`
}

type EvidenceRef struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
	Path string `json:"path,omitempty"`
}
