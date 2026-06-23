// Package acceptance provides black-box CLI fixtures for WP13 acceptance tests.
package acceptance

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/yumiaura/seekmoon/internal/app"
	"github.com/yumiaura/seekmoon/internal/cli"
	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/output"
	"github.com/yumiaura/seekmoon/internal/platform"
	"github.com/yumiaura/seekmoon/internal/service"
)

// Harness owns a real CLI runtime wired to fake services.
type Harness struct {
	Runtime *app.Runtime
	Fakes   *Services
}

// CommandResult records process-style command output from the CLI harness.
type CommandResult struct {
	Code   int
	Output string
}

// NewHarness creates an offline SeekMoon runtime for acceptance tests.
func NewHarness(t *testing.T) *Harness {
	t.Helper()
	rt, err := app.NewRuntime(app.WithEnv(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir(), Home: t.TempDir()}))
	if err != nil {
		t.Fatalf("NewRuntime: %v", err)
	}
	fakes := &Services{}
	rt.Renderer = output.DefaultRenderer{}
	rt.Services.Registry = service.Registry{
		Doctor:  fakes,
		Sync:    fakes,
		Search:  fakes,
		View:    fakes,
		API:     fakes,
		Source:  fakes,
		Skill:   SkillService{Fakes: fakes},
		Compare: fakes,
		Probe:   fakes,
		Record:  fakes,
		Report:  fakes,
		Raw:     fakes,
	}
	return &Harness{Runtime: rt, Fakes: fakes}
}

// Run executes the real Cobra command tree with fake services.
func (h *Harness) Run(args ...string) CommandResult {
	var out bytes.Buffer
	code := cli.ExecuteWithCode(context.Background(), h.Runtime, cli.Options{Out: &out, Err: &out}, args...)
	return CommandResult{Code: code, Output: out.String()}
}

// Services implements all non-skill service interfaces with deterministic data.
type Services struct {
	SearchInputs      []service.SearchInput
	SkillSearchInputs []service.SkillSearchInput
	ViewInputs        []service.ViewInput
	APIInputs         []service.APIInput
	SourceInputs      []service.SourceInput
	ProbeInputs       []service.ProbeInput
	RecordInputs      []service.RecordInput
	ReportInputs      []service.ReportInput
}

// Doctor returns deterministic environment status.
func (s *Services) Doctor(context.Context, service.DoctorInput) (any, error) {
	return model.EnvironmentStatus{Schema: model.SchemaEnvironmentStatusV1}, nil
}

// Sync returns a deterministic snapshot.
func (s *Services) Sync(context.Context) (model.Snapshot, error) {
	return model.Snapshot{ID: "snapshot-acceptance", CreatedAt: "2026-06-24T00:00:00Z"}, nil
}

// Search records search input and returns one deterministic library result.
func (s *Services) Search(_ context.Context, input service.SearchInput) (model.SearchOutput, error) {
	s.SearchInputs = append(s.SearchInputs, input)
	return model.SearchOutput{
		Schema:   model.SchemaSearchResultsV1,
		Snapshot: model.SnapshotRef{ID: "snapshot-acceptance", Sources: []string{string(model.SourceModulesAPI), string(model.SourceManifestAPI)}},
		Query:    model.SearchQuery{Text: input.Query, Kind: "library", Target: input.Target},
		Results: []model.SearchResult{{
			Rank:        1,
			Module:      "mizchi/markdown",
			Version:     "0.6.2",
			Description: model.Present("Incremental Markdown parser", string(model.SourceModulesAPI)),
			Keywords:    model.Present([]string{"markdown"}, string(model.SourceModulesAPI)),
			License:     model.Present("MIT", string(model.SourceModulesAPI)),
			Repository:  model.Present("https://github.com/mizchi/markdown.mbt", string(model.SourceModulesAPI)),
			Downloads:   model.Present(7567, string(model.SourceManifestAPI)),
			BuildStatus: model.Present("success", string(model.SourceManifestAPI)),
			Target:      model.Unknown[map[string]any](),
			Match:       model.MatchEvidence{Fields: []string{"module", "description"}, Token: input.Query},
			SnapshotID:  "snapshot-acceptance",
		}},
	}, nil
}

// View records view input and returns a deterministic manifest profile.
func (s *Services) View(_ context.Context, input service.ViewInput) (model.ManifestProfile, error) {
	s.ViewInputs = append(s.ViewInputs, input)
	return model.ManifestProfile{
		Module:        candidateModule(input.Candidate),
		Version:       "0.6.2",
		LatestVersion: "0.6.2",
		Downloads:     7567,
		HasPackage:    true,
		BuildStatus:   model.Present("success", string(model.SourceManifestAPI)),
		Metadata: model.ManifestMetadata{
			Description:      model.Present("Incremental Markdown parser", string(model.SourceManifestAPI)),
			Keywords:         model.Present([]string{"markdown"}, string(model.SourceManifestAPI)),
			Repository:       model.Present("https://github.com/mizchi/markdown.mbt", string(model.SourceManifestAPI)),
			License:          model.Present("MIT", string(model.SourceManifestAPI)),
			Checksum:         model.Missing[string](string(model.SourceManifestAPI)),
			PreferredTarget:  model.Missing[string](string(model.SourceManifestAPI)),
			SupportedTargets: model.Missing[[]string](string(model.SourceManifestAPI)),
			Raw:              map[string]any{"package_count": 1, "index_state": string(model.StatePresent)},
		},
		Versions:      []model.ManifestVersion{{Version: "0.6.2"}},
		VersionsCount: model.Derived(1, string(model.SourceDerived)),
		DocsURL:       model.Derived("https://mooncakes.io/docs/mizchi/markdown", string(model.SourceDerived)),
	}, nil
}

// API records API input and returns deterministic package data.
func (s *Services) API(_ context.Context, input service.APIInput) (model.PackageData, error) {
	s.APIInputs = append(s.APIInputs, input)
	if input.Package != "mizchi/markdown/src/api" {
		return model.PackageData{}, model.SurfaceFailure{Value: model.SurfaceError{Command: "seekmoon api", Object: input.Package, Source: "module_index.json", State: model.StateFailed, Meaning: "package path is not present in module index", Recovery: "choose one of mizchi/markdown/src/api"}}
	}
	return model.PackageData{
		Name:   "src/api",
		Types:  []model.APIEntry{{Name: "Document", Docstring: model.Present("Markdown document", string(model.SourcePackageData)), Signature: "<a>Document</a>", PlainSignature: model.Derived("Document", string(model.SourceDerived)), Loc: model.Present(map[string]any{"path": "src/api/doc.mbt", "line": 1}, string(model.SourcePackageData))}},
		Values: []model.APIEntry{{Name: "parse", Docstring: model.Present("Parse Markdown", string(model.SourcePackageData)), Signature: "parse(input : String) -> Document", PlainSignature: model.Derived("parse(input : String) -> Document", string(model.SourceDerived)), Loc: model.Present(map[string]any{"path": "src/api/parse.mbt", "line": 7}, string(model.SourcePackageData))}},
	}, nil
}

// Source records source input and returns deterministic source attempts.
func (s *Services) Source(_ context.Context, input service.SourceInput) (model.SourceResolution, error) {
	s.SourceInputs = append(s.SourceInputs, input)
	return model.SourceResolution{
		Module:          candidateModule(input.Candidate),
		Version:         candidateVersion(input.Candidate),
		MoonFetch:       model.SourceAttempt{Status: model.StateFailed, Error: "fake runner disabled"},
		SourceZip:       model.SourceAttempt{Status: model.StatePresent, URL: "https://download.mooncakes.io/user/mizchi/markdown/0.6.2.zip", Path: "source.zip"},
		LocalCache:      model.SourceAttempt{Status: model.StateUnavailable, Path: "/tmp/empty-cache"},
		CoreLocalSource: model.SourceAttempt{Status: model.StateUnavailable, Path: "/tmp/moon/lib/core"},
		RepositorySource: model.SourceAttempt{
			Status: model.StateUnknown,
			URL:    "https://github.com/mizchi/markdown.mbt",
		},
		SelectedSource: model.SelectedSource{Method: "source_zip", URL: "https://download.mooncakes.io/user/mizchi/markdown/0.6.2.zip"},
		FilesSummary:   model.FilesSummary{MoonMod: true, Readme: true, License: true, Sources: 4, Tests: 1},
	}, nil
}

// Compare returns deterministic comparison output.
func (s *Services) Compare(context.Context, service.CompareInput) (any, error) {
	return model.Comparison{Schema: model.SchemaComparisonV1}, nil
}

// Probe records probe input and returns deterministic probe evidence.
func (s *Services) Probe(_ context.Context, input service.ProbeInput) (model.ProbeResult, error) {
	s.ProbeInputs = append(s.ProbeInputs, input)
	return model.ProbeResult{
		Module:          candidateModule(input.Candidate),
		Version:         candidateVersion(input.Candidate),
		Target:          input.Target,
		ProbePath:       "/tmp/seekmoon-probe",
		MoonNew:         command("moon", "new"),
		MoonAdd:         command("moon", "add", candidateModule(input.Candidate)),
		MoonCheck:       command("moon", "check"),
		MoonTest:        command("moon", "test"),
		MoonCheckTarget: command("moon", "check", "--target", input.Target),
		MoonBuildTarget: command("moon", "build", "--target", input.Target),
		Logs:            map[string]string{"moon_build_target": "/tmp/seekmoon-probe/build.log"},
		Result:          "verifiable",
	}, nil
}

// Record records record input and returns deterministic adoption evidence.
func (s *Services) Record(_ context.Context, input service.RecordInput) (model.AdoptionRecord, error) {
	s.RecordInputs = append(s.RecordInputs, input)
	return model.AdoptionRecord{RecordID: "record-acceptance", CreatedAt: time.Date(2026, 6, 24, 0, 0, 0, 0, time.UTC), Project: model.ProjectIdentity{Root: "/tmp/project", Module: "demo/app"}, SnapshotID: "snapshot-acceptance", Candidate: model.CandidateRef{Kind: "library", Module: candidateModule(input.Candidate), Version: candidateVersion(input.Candidate)}, Version: candidateVersion(input.Candidate), Conclusion: input.Conclusion, EvidenceRefs: []model.EvidenceRef{{Kind: "manifest", ID: candidateModule(input.Candidate)}, {Kind: "source", ID: candidateModule(input.Candidate)}}, Note: input.Note}, nil
}

// Report records report input and returns deterministic report evidence.
func (s *Services) Report(_ context.Context, input service.ReportInput) (model.Report, error) {
	s.ReportInputs = append(s.ReportInputs, input)
	record := model.AdoptionRecord{RecordID: "record-acceptance", CreatedAt: time.Date(2026, 6, 24, 0, 0, 0, 0, time.UTC), Project: model.ProjectIdentity{Root: "/tmp/project", Module: "demo/app"}, SnapshotID: "snapshot-acceptance", Candidate: model.CandidateRef{Kind: "library", Module: "mizchi/markdown", Version: "0.6.2"}, Version: "0.6.2", Conclusion: model.ConclusionContinueVerification, EvidenceRefs: []model.EvidenceRef{{Kind: "manifest", ID: "mizchi/markdown"}}}
	return model.Report{Goal: "Assess MoonBit dependency", Environment: model.ReportEnvironment{Toolchain: "moon fake", Snapshot: model.SnapshotRef{ID: "snapshot-acceptance", Sources: []string{string(model.SourceManifestAPI)}}, Project: record.Project}, DataSources: []string{string(model.SourceManifestAPI)}, Candidates: []model.CandidateRef{record.Candidate}, Inspection: []model.EvidenceRef{{Kind: "manifest", ID: "mizchi/markdown"}}, Decision: record}, nil
}

// Raw returns deterministic raw payload output.
func (s *Services) Raw(context.Context, service.RawInput) (any, error) {
	return model.RawEnvelope{Schema: model.SchemaRawPayloadV1, Source: string(model.SourceModulesAPI), Status: model.StatePresent, Payload: map[string]any{"name": "mizchi/markdown"}}, nil
}

// SkillService implements the skill-specific service interface.
type SkillService struct {
	Fakes *Services
}

// Search records skill search input and returns one deterministic skill.
func (s SkillService) Search(_ context.Context, input service.SkillSearchInput) ([]model.SkillEntry, error) {
	s.Fakes.SkillSearchInputs = append(s.Fakes.SkillSearchInputs, input)
	return []model.SkillEntry{skillEntry()}, nil
}

// View returns deterministic skill profile evidence.
func (s SkillService) View(context.Context, service.SkillViewInput) (model.SkillProfile, error) {
	return model.SkillProfile{Entry: skillEntry(), SkillMD: model.Present("# cowsay", string(model.SourceSkillsAPI)), WasmAsset: model.Present(map[string]any{"url": "https://mooncakes.io/cowsay.wasm"}, string(model.SourceSkillsAPI)), ChecksumAsset: model.Present(map[string]any{"url": "https://mooncakes.io/cowsay.sha256"}, string(model.SourceSkillsAPI)), RunwasmCoordinate: model.Derived("Yoorkin/cowsay/cowsay@0.1.0", string(model.SourceDerived))}, nil
}

func candidateModule(candidate model.CandidateRequest) string {
	if candidate.Module != "" {
		return candidate.Module
	}
	return "mizchi/markdown"
}

func candidateVersion(candidate model.CandidateRequest) string {
	if candidate.Version != "" {
		return candidate.Version
	}
	return "0.6.2"
}

func skillEntry() model.SkillEntry {
	return model.SkillEntry{Module: "Yoorkin/cowsay", Author: "Yoorkin", AuthorAvatar: model.Missing[string](string(model.SourceSkillsAPI)), Version: "0.1.0", Package: "cowsay", Name: "cowsay", DetailURL: "skills/Yoorkin/cowsay", WasmURL: "https://mooncakes.io/cowsay.wasm", ChecksumURL: "https://mooncakes.io/cowsay.sha256", Metadata: map[string]any{"description": "cowsay"}, Repository: model.Present("https://github.com/Yoorkin/cowsay", string(model.SourceSkillsAPI))}
}

func command(parts ...string) model.CommandResult {
	return model.CommandResult{Command: parts, CWD: "/tmp/seekmoon-probe", ExitCode: 0, Status: model.StatePresent, LogPath: "/tmp/seekmoon-probe/command.log"}
}
