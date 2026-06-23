package store

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
)

func TestResolvePathsUsesProjectAndXDGCacheBoundaries(t *testing.T) {
	root := t.TempDir()
	cache := t.TempDir()
	paths := ResolvePaths(platform.Env{WorkingDir: root, XDGCacheHome: cache})

	if paths.ProjectDir != filepath.Join(root, ".seekmoon") {
		t.Fatalf("project dir = %q", paths.ProjectDir)
	}
	if paths.Snapshots != filepath.Join(root, ".seekmoon", "snapshots") {
		t.Fatalf("snapshots = %q", paths.Snapshots)
	}
	if paths.CacheRoot != filepath.Join(cache, "seekmoon") {
		t.Fatalf("cache root = %q", paths.CacheRoot)
	}
	if paths.Assets != filepath.Join(cache, "seekmoon", "assets") {
		t.Fatalf("assets = %q", paths.Assets)
	}
}

func TestSafeNameSanitizesRecordAndReportNames(t *testing.T) {
	got := SafeName("mizchi/markdown@0.6.2 js")
	if got != "mizchi-markdown-0.6.2-js" {
		t.Fatalf("SafeName = %q", got)
	}
	if strings.ContainsAny(got, `/\ `) {
		t.Fatalf("SafeName kept path separators or spaces: %q", got)
	}
}

func TestAtomicWriteAndSnapshotRoundTrip(t *testing.T) {
	ctx := context.Background()
	paths := ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	store := SnapshotStore{FS: platform.OSFS{}, Paths: paths}
	snapshot := model.Snapshot{
		ID:        "2026-06-24T12:00:00+08:00",
		CreatedAt: "2026-06-24T12:00:00+08:00",
	}
	if err := store.Write(ctx, snapshot); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, err := store.Read(ctx, snapshot.ID)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got.ID != snapshot.ID {
		t.Fatalf("snapshot ID = %q", got.ID)
	}
	if _, err := os.Stat(store.Path(snapshot.ID) + ".tmp"); !os.IsNotExist(err) {
		t.Fatalf("temporary file should not remain after atomic write")
	}
}

func TestSessionCandidateMappingRoundTrip(t *testing.T) {
	ctx := context.Background()
	paths := ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	store := SessionStore{FS: platform.OSFS{}, Paths: paths}
	candidates := SessionCandidates{
		SessionID:  "default",
		SnapshotID: "snapshot-1",
		Candidates: map[int]model.CandidateRef{
			1: {Kind: "library", Module: "mizchi/markdown", Version: "0.6.2"},
		},
	}
	if err := store.WriteCandidates(ctx, candidates); err != nil {
		t.Fatalf("WriteCandidates: %v", err)
	}
	got, err := store.ReadCandidates(ctx, "default")
	if err != nil {
		t.Fatalf("ReadCandidates: %v", err)
	}
	if got.Candidates[1].Module != "mizchi/markdown" {
		t.Fatalf("candidate mapping mismatch: %#v", got.Candidates)
	}
}

func TestRecordReportAndLogPathsAreSanitized(t *testing.T) {
	paths := ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	record := RecordStore{Paths: paths}
	report := ReportStore{Paths: paths}
	probe := ProbeStore{Paths: paths}

	if strings.Contains(record.Path("mizchi/markdown@0.6.2"), "mizchi/markdown") {
		t.Fatalf("record path was not sanitized: %q", record.Path("mizchi/markdown@0.6.2"))
	}
	if filepath.Ext(report.Path("markdown report", ".md")) != ".md" {
		t.Fatalf("report path extension mismatch")
	}
	logPath := probe.LogPath("mizchi/markdown@0.6.2-js", "moon check --target js")
	if !strings.HasSuffix(logPath, filepath.Join("logs", "moon-check-target-js.log")) {
		t.Fatalf("log path = %q", logPath)
	}
}

func TestRecordStoreRoundTrip(t *testing.T) {
	ctx := context.Background()
	paths := ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	store := RecordStore{FS: platform.OSFS{}, Paths: paths}
	record := model.AdoptionRecord{
		RecordID:     "mizchi/markdown@0.6.2",
		CreatedAt:    time.Date(2026, 6, 24, 12, 0, 0, 0, time.UTC),
		Project:      model.ProjectIdentity{Root: paths.ProjectRoot, Module: "demo/project"},
		SnapshotID:   "snapshot-1",
		Candidate:    model.CandidateRef{Kind: "library", Module: "mizchi/markdown", Version: "0.6.2"},
		Version:      "0.6.2",
		Conclusion:   model.ConclusionContinueVerification,
		EvidenceRefs: []model.EvidenceRef{{Kind: "manifest", ID: "manifest-1"}},
	}
	if err := store.Write(ctx, record); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, err := store.Read(ctx, record.RecordID)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got.Conclusion != model.ConclusionContinueVerification {
		t.Fatalf("conclusion = %q", got.Conclusion)
	}
}
