package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
	"github.com/yumiaura/seekmoon/internal/testutil"
)

func TestDoctorReportsProjectContextWithoutWritingRecords(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	cache := t.TempDir()
	paths := store.ResolvePaths(platform.Env{WorkingDir: root, XDGCacheHome: cache})
	fs := platform.OSFS{}
	if err := fs.WriteFile(ctx, filepath.Join(root, "moon.mod.json"), []byte(`{"name":"demo/app","deps":{"mizchi/markdown":"0.6.2"}}`), 0o644); err != nil {
		t.Fatalf("write project config: %v", err)
	}
	svc := DoctorFlow{Project: source.ProjectReader{FS: fs}, Paths: paths}
	got, err := svc.Doctor(ctx, DoctorInput{})
	if err != nil {
		t.Fatalf("Doctor: %v", err)
	}
	status := got.(model.EnvironmentStatus)
	if status.Project.Value == nil || status.Project.Value.Identity.Module != "demo/app" {
		t.Fatalf("project status = %#v", status.Project)
	}
	if _, err := fs.ReadFile(ctx, filepath.Join(paths.Records, "anything.json")); err == nil {
		t.Fatal("doctor should not create records")
	}
}

func TestSyncRecordsPartialSourceFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()
		switch path {
		case "/api/v0/modules":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[{"name":"mizchi/markdown","version":"0.6.2"}]`))
		case "/api/v0/modules/statistics":
			http.Error(w, "down", http.StatusBadGateway)
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	paths := store.ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir(), Home: t.TempDir()})
	service := SyncService{
		Mooncakes:  source.MooncakesClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}},
		LocalIndex: source.LocalIndexReader{FS: platform.OSFS{}},
		Snapshots:  store.SnapshotStore{FS: platform.OSFS{}, Paths: paths},
		Paths:      paths,
		Now:        func() time.Time { return time.Date(2026, 6, 24, 1, 2, 3, 0, time.UTC) },
	}
	snapshot, err := service.Sync(context.Background())
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if len(snapshot.Sources) != 3 {
		t.Fatalf("sources = %#v", snapshot.Sources)
	}
	if snapshot.Sources[0].Status != model.StatePresent || snapshot.Sources[1].Status != model.StateFailed || snapshot.Sources[2].Status != model.StateUnavailable {
		t.Fatalf("partial source states = %#v", snapshot.Sources)
	}
	localIndex, ok := snapshot.Raw["local_index"].(map[string]any)
	if !ok || localIndex["status"] != model.StateUnavailable || localIndex["error"] == "" {
		t.Fatalf("local index failure summary = %#v", snapshot.Raw["local_index"])
	}
	read, err := service.Snapshots.Read(context.Background(), snapshot.ID)
	if err != nil {
		t.Fatalf("read snapshot: %v", err)
	}
	if len(read.Sources) != 3 || read.Raw["moon_version"] == nil {
		t.Fatalf("stored snapshot = %#v", read)
	}
}

func TestSyncRecordsLocalIndexSummaryWhenPresent(t *testing.T) {
	server := mooncakesServer(t)
	defer server.Close()

	paths := store.ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir(), Home: t.TempDir()})
	testutil.WriteFile(t, paths.MoonIndex, "mizchi/markdown.index", `{"name":"mizchi/markdown","version":"0.6.1"}
{"name":"mizchi/markdown","version":"0.6.2"}`)
	testutil.WriteFile(t, paths.MoonIndex, "moonbit-community/cmark.index", `{"name":"moonbit-community/cmark","version":"0.4.4"}`)
	service := SyncService{
		Mooncakes:  source.MooncakesClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}},
		LocalIndex: source.LocalIndexReader{FS: platform.OSFS{}},
		Snapshots:  store.SnapshotStore{FS: platform.OSFS{}, Paths: paths},
		Paths:      paths,
		Now:        func() time.Time { return time.Date(2026, 6, 24, 1, 2, 3, 0, time.UTC) },
	}
	snapshot, err := service.Sync(context.Background())
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if snapshot.Sources[2].Source != string(model.SourceLocalIndex) || snapshot.Sources[2].Status != model.StatePresent {
		t.Fatalf("local index source = %#v", snapshot.Sources[2])
	}
	localIndex, ok := snapshot.Raw["local_index"].(map[string]any)
	if !ok {
		t.Fatalf("local index raw summary = %#v", snapshot.Raw["local_index"])
	}
	if localIndex["file_count"] != 2 || localIndex["record_count"] != 3 || localIndex["malformed"] != 0 {
		t.Fatalf("local index counts = %#v", localIndex)
	}
}

func TestSearchUsesTransientModulesFetchAndWritesSessionMap(t *testing.T) {
	ctx := context.Background()
	server := mooncakesServer(t)
	defer server.Close()
	paths := store.ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	svc := SearchFlow{
		Mooncakes: source.MooncakesClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}},
		Snapshots: store.SnapshotStore{FS: platform.OSFS{}, Paths: paths},
		Sessions:  store.SessionStore{FS: platform.OSFS{}, Paths: paths},
	}
	output, err := svc.Search(ctx, SearchInput{Query: "markdown", Target: "js"})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if output.Snapshot.ID != "transient" || len(output.Results) != 1 {
		t.Fatalf("search output = %#v", output)
	}
	if output.Results[0].Target.Status != model.StateUnknown {
		t.Fatalf("target evidence = %#v", output.Results[0].Target)
	}
	session, err := svc.Sessions.ReadCandidates(ctx, defaultSessionID)
	if err != nil {
		t.Fatalf("session: %v", err)
	}
	if session.Candidates[1].Module != "mizchi/markdown" {
		t.Fatalf("session candidates = %#v", session.Candidates)
	}
}

func TestViewReadsManifestAndModuleIndex(t *testing.T) {
	ctx := context.Background()
	server := mooncakesServer(t)
	defer server.Close()
	paths := store.ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	svc := ViewFlow{
		Mooncakes: source.MooncakesClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}},
		Assets:    source.AssetClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}},
		Sessions:  store.SessionStore{FS: platform.OSFS{}, Paths: paths},
	}
	profile, err := svc.View(ctx, ViewInput{Candidate: model.CandidateRequest{Module: "mizchi/markdown"}})
	if err != nil {
		t.Fatalf("View: %v", err)
	}
	if profile.Module != "mizchi/markdown" || profile.Metadata.Raw["package_count"].(int) != 1 {
		t.Fatalf("profile = %#v", profile)
	}
}

func TestCompareAlignsEvidenceWithoutScore(t *testing.T) {
	ctx := context.Background()
	server := mooncakesServer(t)
	defer server.Close()
	svc := CompareFlow{Mooncakes: source.MooncakesClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}}}
	got, err := svc.Compare(ctx, CompareInput{Candidates: []model.CandidateRequest{{Module: "mizchi/markdown"}, {Module: "moonbit-community/cmark"}}})
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	comparison := got.(model.Comparison)
	if len(comparison.Fields) == 0 || strings.Contains(strings.ToLower(mustJSON(comparison)), "score") || strings.Contains(strings.ToLower(mustJSON(comparison)), "recommendation") {
		t.Fatalf("comparison = %#v", comparison)
	}
}

func TestAPIReadsModuleIndexBeforePackageDataAndReportsKnownPackages(t *testing.T) {
	ctx := context.Background()
	server := mooncakesServer(t)
	defer server.Close()
	svc := APIFlow{
		Mooncakes: source.MooncakesClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}},
		Assets:    source.AssetClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}},
	}
	data, err := svc.API(ctx, APIInput{Candidate: model.CandidateRequest{Module: "mizchi/markdown", Version: "0.6.2"}, Package: "mizchi/markdown/src/api"})
	if err != nil {
		t.Fatalf("API: %v", err)
	}
	if len(data.Types) != 1 || data.Types[0].PlainSignature.Value == nil {
		t.Fatalf("package data = %#v", data)
	}
	_, err = svc.API(ctx, APIInput{Candidate: model.CandidateRequest{Module: "mizchi/markdown", Version: "0.6.2"}, Package: "mizchi/markdown/missing"})
	if err == nil || !strings.Contains(err.Error(), "known packages") {
		t.Fatalf("known package error = %v", err)
	}
}

func TestSourceRecordsAttemptsAndSelectsDeterministically(t *testing.T) {
	ctx := context.Background()
	server := mooncakesServer(t)
	defer server.Close()
	paths := store.ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	svc := SourceFlow{
		Mooncakes:  source.MooncakesClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}},
		Assets:     source.AssetClient{BaseURL: server.URL, DownloadURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}},
		LocalCache: source.LocalCacheReader{FS: platform.OSFS{}},
		Paths:      paths,
	}
	got, err := svc.Source(ctx, SourceInput{Candidate: model.CandidateRequest{Module: "mizchi/markdown", Version: "0.6.2"}})
	if err != nil {
		t.Fatalf("Source: %v", err)
	}
	if got.SourceZip.Status != model.StatePresent || got.LocalCache.Status == "" || got.RepositorySource.URL == "" {
		t.Fatalf("source attempts = %#v", got)
	}
	if got.SelectedSource.Method != "source_zip" {
		t.Fatalf("selected source = %#v", got.SelectedSource)
	}
}

func TestSkillSearchUsesSkillsAPIAndViewDerivesRunwasmCoordinate(t *testing.T) {
	ctx := context.Background()
	server := mooncakesServer(t)
	defer server.Close()
	paths := store.ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	svc := SkillFlow{Skills: source.SkillsClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}}, Sessions: store.SessionStore{FS: platform.OSFS{}, Paths: paths}}
	results, err := svc.Search(ctx, SkillSearchInput{Query: "cowsay"})
	if err != nil {
		t.Fatalf("Skill search: %v", err)
	}
	if len(results) != 1 || results[0].Module != "Yoorkin/cowsay" {
		t.Fatalf("skill results = %#v", results)
	}
	profile, err := svc.View(ctx, SkillViewInput{Entry: model.CandidateRequest{Module: "Yoorkin/cowsay"}})
	if err != nil {
		t.Fatalf("Skill view: %v", err)
	}
	if profile.RunwasmCoordinate.Value == nil || *profile.RunwasmCoordinate.Value != "Yoorkin/cowsay/cowsay@0.1.0" {
		t.Fatalf("profile = %#v", profile)
	}
}

func TestProbeUsesIsolatedPathAndRecordsFailedTargetLogs(t *testing.T) {
	ctx := context.Background()
	paths := store.ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	sessions := store.SessionStore{FS: platform.OSFS{}, Paths: paths}
	if err := sessions.WriteCandidates(ctx, store.SessionCandidates{SessionID: defaultSessionID, Candidates: map[int]model.CandidateRef{1: {Kind: "library", Module: "mizchi/markdown", Version: "0.6.2"}}}); err != nil {
		t.Fatalf("write session: %v", err)
	}
	svc := ProbeFlow{
		MoonCLI:  source.MoonCLI{Runner: fakeRunner{failContains: "build"}, Paths: paths},
		Sessions: sessions,
		Probes:   store.ProbeStore{Paths: paths},
		Paths:    paths,
	}
	result, err := svc.Probe(ctx, ProbeInput{Candidate: model.CandidateRequest{Number: 1, Raw: "1"}, Target: "js"})
	if err != nil {
		t.Fatalf("Probe: %v", err)
	}
	if !strings.Contains(result.ProbePath, ".seekmoon") || result.MoonBuildTarget.Status != model.StateFailed || result.Logs["moon_build_target"] == "" || result.Result != "failed" {
		t.Fatalf("probe result = %#v", result)
	}
}

func TestRecordValidatesConclusionAndWritesEvidenceRefs(t *testing.T) {
	ctx := context.Background()
	paths := store.ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	sessions := store.SessionStore{FS: platform.OSFS{}, Paths: paths}
	if err := sessions.WriteCandidates(ctx, store.SessionCandidates{SessionID: defaultSessionID, Candidates: map[int]model.CandidateRef{1: {Kind: "library", Module: "mizchi/markdown", Version: "0.6.2"}}}); err != nil {
		t.Fatalf("write session: %v", err)
	}
	svc := RecordFlow{Sessions: sessions, Records: store.RecordStore{FS: platform.OSFS{}, Paths: paths}, Project: source.ProjectReader{FS: platform.OSFS{}}, Paths: paths, Now: func() time.Time { return time.Date(2026, 6, 24, 0, 0, 0, 0, time.UTC) }}
	record, err := svc.Record(ctx, RecordInput{Candidate: model.CandidateRequest{Number: 1, Raw: "1"}, Conclusion: model.ConclusionContinueVerification, Note: "needs probe"})
	if err != nil {
		t.Fatalf("Record: %v", err)
	}
	if len(record.EvidenceRefs) == 0 || record.Note != "needs probe" {
		t.Fatalf("record = %#v", record)
	}
}

func TestReportOmitsUnobservedSources(t *testing.T) {
	ctx := context.Background()
	paths := store.ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	records := store.RecordStore{FS: platform.OSFS{}, Paths: paths}
	record := model.AdoptionRecord{RecordID: "one", Project: model.ProjectIdentity{Root: paths.ProjectRoot}, Candidate: model.CandidateRef{Kind: "library", Module: "mizchi/markdown"}, Conclusion: model.ConclusionContinueVerification, EvidenceRefs: []model.EvidenceRef{{Kind: "manifest", ID: "mizchi/markdown"}}}
	if err := records.Write(ctx, record); err != nil {
		t.Fatalf("write record: %v", err)
	}
	svc := ReportFlow{Records: records, Reports: store.ReportStore{FS: platform.OSFS{}, Paths: paths}, Paths: paths}
	report, err := svc.Report(ctx, ReportInput{Format: "json"})
	if err != nil {
		t.Fatalf("Report: %v", err)
	}
	if strings.Contains(strings.Join(report.DataSources, ","), "repository") || len(report.Validation) != 0 {
		t.Fatalf("report = %#v", report)
	}
}

func TestRawReturnsSourceStatusAndUpstreamPayloadShape(t *testing.T) {
	ctx := context.Background()
	server := mooncakesServer(t)
	defer server.Close()
	svc := RawFlow{
		Mooncakes: source.MooncakesClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}},
		Assets:    source.AssetClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}},
		Skills:    source.SkillsClient{BaseURL: server.URL, Fetcher: source.Fetcher{Client: server.Client()}},
	}
	got, err := svc.Raw(ctx, RawInput{Source: "modules"})
	if err != nil {
		t.Fatalf("Raw: %v", err)
	}
	envelope := got.(model.RawEnvelope)
	if envelope.Status != model.StatePresent || envelope.Source != string(model.SourceModulesAPI) || envelope.Payload == nil {
		t.Fatalf("raw envelope = %#v", envelope)
	}
	modules, ok := envelope.Payload.([]any)
	if !ok || len(modules) == 0 {
		t.Fatalf("raw modules payload = %#v", envelope.Payload)
	}
	firstModule := modules[0].(map[string]any)
	if firstModule["name"] != "mizchi/markdown" || firstModule["module"] != nil {
		t.Fatalf("modules raw fields normalized = %#v", firstModule)
	}

	got, err = svc.Raw(ctx, RawInput{Source: "manifest", Args: []string{"mizchi/markdown"}})
	if err != nil {
		t.Fatalf("Raw manifest: %v", err)
	}
	manifest := got.(model.RawEnvelope).Payload.(map[string]any)
	if manifest["latest_version"] != "0.6.2" || manifest["docs_url"] != nil || manifest["versions_count"] != nil {
		t.Fatalf("manifest raw fields normalized = %#v", manifest)
	}

	got, err = svc.Raw(ctx, RawInput{Source: "module-index", Args: []string{"mizchi/markdown@0.6.2"}})
	if err != nil {
		t.Fatalf("Raw module-index: %v", err)
	}
	index := got.(model.RawEnvelope).Payload.(map[string]any)
	if index["childs"] != nil || index["children"] == nil {
		t.Fatalf("module index raw shape hidden = %#v", index)
	}

	got, err = svc.Raw(ctx, RawInput{Source: "package-data", Args: []string{"mizchi/markdown@0.6.2", "mizchi/markdown/src/api"}})
	if err != nil {
		t.Fatalf("Raw package-data: %v", err)
	}
	packageData := got.(model.RawEnvelope).Payload.(map[string]any)
	types := packageData["types"].([]any)
	firstType := types[0].(map[string]any)
	if firstType["signature"] != "<a>Document</a>" || firstType["plain_signature"] != nil {
		t.Fatalf("package data raw fields normalized = %#v", firstType)
	}
}

func mooncakesServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v0/modules":
			_, _ = w.Write([]byte(`[{"name":"mizchi/markdown","version":"0.6.2","description":"Incremental Markdown parser","keywords":["markdown"],"repository":"https://github.com/mizchi/markdown.mbt","license":"MIT"},{"name":"moonbit-community/cmark","version":"0.4.4","description":"CommonMark parser","keywords":["cmark"],"repository":"https://github.com/moonbit-community/cmark","license":"Apache-2"}]`))
		case "/api/v0/modules/statistics":
			_, _ = w.Write([]byte(`{"total_modules":2,"total_packages":3,"total_lines":4,"total_downloads":5}`))
		case "/api/v0/manifest/mizchi/markdown":
			_, _ = w.Write([]byte(`{"module":"mizchi/markdown","version":"0.6.2","latest_version":"0.6.2","downloads":7567,"has_package":true,"build_status":"success","metadata":{"description":"Incremental Markdown parser","license":"MIT","repository":"https://github.com/mizchi/markdown.mbt"},"versions":[{"version":"0.6.2"}]}`))
		case "/api/v0/manifest/moonbit-community/cmark":
			_, _ = w.Write([]byte(`{"module":"moonbit-community/cmark","version":"0.4.4","latest_version":"0.4.4","downloads":3120,"has_package":true,"build_status":"success","metadata":{"description":"CommonMark parser","license":"Apache-2","repository":"https://github.com/moonbit-community/cmark"},"versions":[{"version":"0.4.4"}]}`))
		case "/assets/mizchi/markdown@0.6.2/module_index.json":
			_, _ = w.Write([]byte(`{"name":"markdown","package":{"path":"mizchi/markdown/src/api","types":[{"name":"Document"}],"values":[{"name":"parse"}]},"children":[]}`))
		case "/assets/mizchi/markdown@0.6.2/src/api/package_data.json":
			_, _ = w.Write([]byte(`{"name":"src/api","types":[{"name":"Document","signature":"<a>Document</a>"}],"values":[{"name":"parse","signature":"parse(input : String) -> Document"}]}`))
		case "/user/mizchi/markdown/0.6.2.zip":
			w.Header().Set("Content-Type", "application/zip")
			_, _ = w.Write(testZip(t))
		case "/api/v0/skills":
			_, _ = w.Write([]byte(`[{"module":"Yoorkin/cowsay","version":"0.1.0","package":"cowsay","name":"cowsay","detail_url":"skills/Yoorkin/cowsay","wasm_url":"` + "http://" + r.Host + `/skill.wasm","checksum_url":"` + "http://" + r.Host + `/skill.sha256","metadata":{"description":"cowsay"}}]`))
		case "/api/v0/skills/Yoorkin%2Fcowsay", "/api/v0/skills/Yoorkin/cowsay":
			_, _ = w.Write([]byte(`{"module":"Yoorkin/cowsay","version":"0.1.0","package":"cowsay","name":"cowsay","detail_url":"skills/Yoorkin/cowsay","wasm_url":"` + "http://" + r.Host + `/skill.wasm","checksum_url":"` + "http://" + r.Host + `/skill.sha256","metadata":{"description":"cowsay"}}`))
		case "/skill.wasm", "/skill.sha256":
			_, _ = w.Write([]byte("asset"))
		default:
			http.NotFound(w, r)
		}
	}))
}

func testZip(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, name := range []string{"moon.mod.json", "README.md", "LICENSE", "src/lib.mbt", "test/lib_test.mbt"} {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("zip create: %v", err)
		}
		_, _ = w.Write([]byte("x"))
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zip close: %v", err)
	}
	return buf.Bytes()
}

func mustJSON(value any) string {
	data, _ := json.Marshal(value)
	return string(data)
}

type fakeRunner struct {
	failContains string
}

func (r fakeRunner) Run(ctx context.Context, request platform.RunRequest) (platform.RunResult, error) {
	result := platform.RunResult{Command: request.Command, CWD: request.CWD, LogPath: request.LogPath}
	if strings.Contains(strings.Join(request.Command, " "), r.failContains) {
		result.ExitCode = 1
		return result, errFakeRunner
	}
	return result, nil
}

var errFakeRunner = fakeRunError{}

type fakeRunError struct{}

func (fakeRunError) Error() string { return "fake command failed" }
