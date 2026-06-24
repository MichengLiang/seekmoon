package service

import (
	"context"
	"fmt"

	"github.com/MichengLiang/seekmoon/internal/model"
)

// Registry groups all command service interfaces.
type Registry struct {
	Doctor  DoctorService
	Sync    SyncRunner
	Search  SearchService
	View    ViewService
	API     APIService
	Source  SourceService
	Skill   SkillService
	Compare CompareService
	Probe   ProbeService
	Record  RecordService
	Report  ReportService
	Raw     RawService
}

// DoctorService handles the doctor command.
type DoctorService interface {
	Doctor(context.Context, DoctorInput) (any, error)
}

// SyncRunner handles the sync command.
type SyncRunner interface {
	Sync(context.Context) (model.Snapshot, error)
}

// SearchService handles library search.
type SearchService interface {
	Search(context.Context, SearchInput) (model.SearchOutput, error)
}

// ViewService handles manifest inspection.
type ViewService interface {
	View(context.Context, ViewInput) (model.ManifestProfile, error)
}

// APIService handles package API inspection.
type APIService interface {
	API(context.Context, APIInput) (model.PackageData, error)
}

// SourceService handles source resolution.
type SourceService interface {
	Source(context.Context, SourceInput) (model.SourceResolution, error)
}

// SkillService handles skill search and inspection.
type SkillService interface {
	Search(context.Context, SkillSearchInput) ([]model.SkillEntry, error)
	View(context.Context, SkillViewInput) (model.SkillProfile, error)
}

// CompareService handles candidate comparison.
type CompareService interface {
	Compare(context.Context, CompareInput) (any, error)
}

// ProbeService handles isolated probe execution.
type ProbeService interface {
	Probe(context.Context, ProbeInput) (model.ProbeResult, error)
}

// RecordService handles adoption record creation.
type RecordService interface {
	Record(context.Context, RecordInput) (model.AdoptionRecord, error)
}

// ReportService handles report generation.
type ReportService interface {
	Report(context.Context, ReportInput) (model.Report, error)
}

// RawService handles raw source projection.
type RawService interface {
	Raw(context.Context, RawInput) (any, error)
}

// DoctorInput is the doctor command input.
type DoctorInput struct{}

// SearchInput is the library search command input.
type SearchInput struct {
	Query  string
	Target string
	Kind   string
}

// ViewInput is the manifest view command input.
type ViewInput struct {
	Candidate model.CandidateRequest
}

// APIInput is the package API command input.
type APIInput struct {
	Candidate model.CandidateRequest
	Package   string
}

// SourceInput is the source command input.
type SourceInput struct {
	Candidate model.CandidateRequest
}

// SkillSearchInput is the skill search command input.
type SkillSearchInput struct {
	Query string
}

// SkillViewInput is the skill view command input.
type SkillViewInput struct {
	Entry model.CandidateRequest
}

// CompareInput is the compare command input.
type CompareInput struct {
	Candidates []model.CandidateRequest
}

// ProbeInput is the probe command input.
type ProbeInput struct {
	Candidate model.CandidateRequest
	Target    string
}

// RecordInput is the record command input.
type RecordInput struct {
	Candidate  model.CandidateRequest
	Kind       string
	Conclusion model.AdoptionConclusion
	Note       string
}

// ReportInput is the report command input.
type ReportInput struct {
	Format string
}

// RawInput is the raw command input.
type RawInput struct {
	Source string
	Args   []string
}

// NewPendingRegistry creates placeholder services around an available sync runner.
func NewPendingRegistry(sync SyncRunner) Registry {
	pending := PendingService{}
	return Registry{
		Doctor:  pending,
		Sync:    sync,
		Search:  pending,
		View:    pending,
		API:     pending,
		Source:  pending,
		Skill:   PendingSkillService{},
		Compare: pending,
		Probe:   pending,
		Record:  pending,
		Report:  pending,
		Raw:     pending,
	}
}

// PendingService returns explicit pending errors for unimplemented commands.
type PendingService struct{}

// Doctor returns a pending-service error.
func (PendingService) Doctor(context.Context, DoctorInput) (any, error) {
	return nil, pending("doctor")
}

// Search returns a pending-service error.
func (PendingService) Search(context.Context, SearchInput) (model.SearchOutput, error) {
	return model.SearchOutput{}, pending("search")
}

// View returns a pending-service error.
func (PendingService) View(context.Context, ViewInput) (model.ManifestProfile, error) {
	return model.ManifestProfile{}, pending("view")
}

// API returns a pending-service error.
func (PendingService) API(context.Context, APIInput) (model.PackageData, error) {
	return model.PackageData{}, pending("api")
}

// Source returns a pending-service error.
func (PendingService) Source(context.Context, SourceInput) (model.SourceResolution, error) {
	return model.SourceResolution{}, pending("source")
}

// Compare returns a pending-service error.
func (PendingService) Compare(context.Context, CompareInput) (any, error) {
	return nil, pending("compare")
}

// Probe returns a pending-service error.
func (PendingService) Probe(context.Context, ProbeInput) (model.ProbeResult, error) {
	return model.ProbeResult{}, pending("probe")
}

// Record returns a pending-service error.
func (PendingService) Record(context.Context, RecordInput) (model.AdoptionRecord, error) {
	return model.AdoptionRecord{}, pending("record")
}

// Report returns a pending-service error.
func (PendingService) Report(context.Context, ReportInput) (model.Report, error) {
	return model.Report{}, pending("report")
}

// Raw returns a pending-service error.
func (PendingService) Raw(context.Context, RawInput) (any, error) {
	return nil, pending("raw")
}

// PendingSkillService returns explicit pending errors for skill commands.
type PendingSkillService struct{}

// Search returns a pending-service error.
func (PendingSkillService) Search(context.Context, SkillSearchInput) ([]model.SkillEntry, error) {
	return nil, pending("skill search")
}

// View returns a pending-service error.
func (PendingSkillService) View(context.Context, SkillViewInput) (model.SkillProfile, error) {
	return model.SkillProfile{}, pending("skill view")
}

func pending(command string) error {
	return fmt.Errorf("%s service behavior is outside Batch C", command)
}
