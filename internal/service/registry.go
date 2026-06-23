package service

import (
	"context"
	"fmt"

	"github.com/yumiaura/seekmoon/internal/model"
)

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

type DoctorService interface {
	Doctor(context.Context, DoctorInput) (any, error)
}

type SyncRunner interface {
	Sync(context.Context) (model.Snapshot, error)
}

type SearchService interface {
	Search(context.Context, SearchInput) (model.SearchOutput, error)
}

type ViewService interface {
	View(context.Context, ViewInput) (model.ManifestProfile, error)
}

type APIService interface {
	API(context.Context, APIInput) (model.PackageData, error)
}

type SourceService interface {
	Source(context.Context, SourceInput) (model.SourceResolution, error)
}

type SkillService interface {
	Search(context.Context, SkillSearchInput) ([]model.SkillEntry, error)
	View(context.Context, SkillViewInput) (model.SkillProfile, error)
}

type CompareService interface {
	Compare(context.Context, CompareInput) (any, error)
}

type ProbeService interface {
	Probe(context.Context, ProbeInput) (model.ProbeResult, error)
}

type RecordService interface {
	Record(context.Context, RecordInput) (model.AdoptionRecord, error)
}

type ReportService interface {
	Report(context.Context, ReportInput) (model.Report, error)
}

type RawService interface {
	Raw(context.Context, RawInput) (any, error)
}

type DoctorInput struct{}

type SearchInput struct {
	Query  string
	Target string
	Kind   string
}

type ViewInput struct {
	Candidate model.CandidateRequest
}

type APIInput struct {
	Candidate model.CandidateRequest
	Package   string
}

type SourceInput struct {
	Candidate model.CandidateRequest
}

type SkillSearchInput struct {
	Query string
}

type SkillViewInput struct {
	Entry model.CandidateRequest
}

type CompareInput struct {
	Candidates []model.CandidateRequest
}

type ProbeInput struct {
	Candidate model.CandidateRequest
	Target    string
}

type RecordInput struct {
	Candidate  model.CandidateRequest
	Kind       string
	Conclusion model.AdoptionConclusion
	Note       string
}

type ReportInput struct {
	Format string
}

type RawInput struct {
	Source string
	Args   []string
}

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

type PendingService struct{}

func (PendingService) Doctor(context.Context, DoctorInput) (any, error) {
	return nil, pending("doctor")
}

func (PendingService) Search(context.Context, SearchInput) (model.SearchOutput, error) {
	return model.SearchOutput{}, pending("search")
}

func (PendingService) View(context.Context, ViewInput) (model.ManifestProfile, error) {
	return model.ManifestProfile{}, pending("view")
}

func (PendingService) API(context.Context, APIInput) (model.PackageData, error) {
	return model.PackageData{}, pending("api")
}

func (PendingService) Source(context.Context, SourceInput) (model.SourceResolution, error) {
	return model.SourceResolution{}, pending("source")
}

func (PendingService) Compare(context.Context, CompareInput) (any, error) {
	return nil, pending("compare")
}

func (PendingService) Probe(context.Context, ProbeInput) (model.ProbeResult, error) {
	return model.ProbeResult{}, pending("probe")
}

func (PendingService) Record(context.Context, RecordInput) (model.AdoptionRecord, error) {
	return model.AdoptionRecord{}, pending("record")
}

func (PendingService) Report(context.Context, ReportInput) (model.Report, error) {
	return model.Report{}, pending("report")
}

func (PendingService) Raw(context.Context, RawInput) (any, error) {
	return nil, pending("raw")
}

type PendingSkillService struct{}

func (PendingSkillService) Search(context.Context, SkillSearchInput) ([]model.SkillEntry, error) {
	return nil, pending("skill search")
}

func (PendingSkillService) View(context.Context, SkillViewInput) (model.SkillProfile, error) {
	return model.SkillProfile{}, pending("skill view")
}

func pending(command string) error {
	return fmt.Errorf("%s service behavior is outside Batch C", command)
}
