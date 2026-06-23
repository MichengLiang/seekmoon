package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

// ReportFlow builds reports from stored adoption records.
type ReportFlow struct {
	Records store.RecordStore
	Reports store.ReportStore
	Project source.ProjectReader
	Paths   store.Paths
	Now     func() time.Time
}

// Report returns and optionally writes a report.
func (s ReportFlow) Report(ctx context.Context, input ReportInput) (model.Report, error) {
	records, err := s.Records.List(ctx)
	if err != nil {
		return model.Report{}, err
	}
	if len(records) == 0 {
		return model.Report{}, surfaceFailure("report", "record_store", model.StateUnavailable, "no adoption records are available for report", "create a record first")
	}
	decision := records[len(records)-1]
	dataSources := usedSources(records)
	report := model.Report{
		Goal:        "SeekMoon investigation report",
		Environment: model.ReportEnvironment{Snapshot: model.SnapshotRef{ID: decision.SnapshotID}, Project: decision.Project},
		DataSources: dataSources,
		Candidates:  candidatesFromRecords(records),
		Inspection:  inspectionRefs(records),
		Validation:  validationRefs(records),
		Decision:    decision,
	}
	if input.Format != "" {
		name := "report-" + time.Now().Format("2006-01-02")
		if s.Now != nil {
			name = "report-" + s.Now().Format("2006-01-02")
		}
		ext := ".json"
		if input.Format == "markdown" || input.Format == "md" {
			ext = ".md"
		}
		_, _ = s.Reports.Write(ctx, name, ext, reportBytes(report, ext))
	}
	return report, nil
}

func usedSources(records []model.AdoptionRecord) []string {
	seen := map[string]bool{}
	var out []string
	for _, record := range records {
		for _, ref := range record.EvidenceRefs {
			if ref.Kind == "" || seen[ref.Kind] {
				continue
			}
			seen[ref.Kind] = true
			out = append(out, ref.Kind)
		}
	}
	return out
}

func candidatesFromRecords(records []model.AdoptionRecord) []model.CandidateRef {
	out := make([]model.CandidateRef, 0, len(records))
	for _, record := range records {
		out = append(out, record.Candidate)
	}
	return out
}

func inspectionRefs(records []model.AdoptionRecord) []model.EvidenceRef {
	var out []model.EvidenceRef
	for _, record := range records {
		for _, ref := range record.EvidenceRefs {
			switch ref.Kind {
			case "manifest", "api", "source", "skill", "candidate":
				out = append(out, ref)
			}
		}
	}
	return out
}

func validationRefs(records []model.AdoptionRecord) []model.EvidenceRef {
	var out []model.EvidenceRef
	for _, record := range records {
		for _, ref := range record.EvidenceRefs {
			if strings.HasPrefix(ref.Kind, "probe") {
				out = append(out, ref)
			}
		}
	}
	return out
}

func reportBytes(report model.Report, ext string) []byte {
	if ext == ".md" {
		return []byte("# SeekMoon Report\n\nGoal: " + report.Goal + "\n")
	}
	data, _ := json.MarshalIndent(report, "", "  ")
	return data
}
