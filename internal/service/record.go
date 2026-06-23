package service

import (
	"context"
	"time"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

type RecordFlow struct {
	Sessions store.SessionStore
	Records  store.RecordStore
	Project  source.ProjectReader
	Paths    store.Paths
	Now      func() time.Time
}

func (s RecordFlow) Record(ctx context.Context, input RecordInput) (model.AdoptionRecord, error) {
	if !input.Conclusion.IsValid() {
		return model.AdoptionRecord{}, surfaceFailure("record conclusion", "command input", model.StateFailed, "unknown adoption conclusion", "choose a fixed adoption conclusion enum value")
	}
	candidate, err := candidateFromRequest(ctx, s.Sessions, input.Candidate)
	if err != nil {
		return model.AdoptionRecord{}, err
	}
	if input.Kind != "" {
		candidate.Kind = input.Kind
	}
	now := time.Now()
	if s.Now != nil {
		now = s.Now()
	}
	project := s.Project.Read(ctx, s.Paths.ProjectRoot)
	record := model.AdoptionRecord{
		RecordID:     store.SafeName(candidate.Kind + "-" + candidate.Module + "@" + candidate.Version + "-" + candidate.Package),
		CreatedAt:    now,
		Project:      project.Value.Identity,
		SnapshotID:   "default",
		Candidate:    candidate,
		Version:      candidate.Version,
		Conclusion:   input.Conclusion,
		EvidenceRefs: []model.EvidenceRef{{Kind: "candidate", ID: candidate.Module, Path: candidate.Package}},
		NotConfirmed: []string{"repository activity", "target support"},
		Note:         input.Note,
	}
	if record.Project.Root == "" {
		record.Project.Root = s.Paths.ProjectRoot
	}
	if err := s.Records.Write(ctx, record); err != nil {
		return model.AdoptionRecord{}, err
	}
	return record, nil
}
