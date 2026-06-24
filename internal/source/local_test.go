package source

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/MichengLiang/seekmoon/internal/model"
	"github.com/MichengLiang/seekmoon/internal/platform"
	"github.com/MichengLiang/seekmoon/internal/store"
	"github.com/MichengLiang/seekmoon/internal/testutil"
)

func TestMoonCLICommandResultMapping(t *testing.T) {
	paths := store.ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	runner := &testutil.FakeRunner{Result: platform.RunResult{ExitCode: 0}}
	cli := MoonCLI{Runner: runner, Paths: paths}
	result := cli.Version(context.Background(), paths.ProjectRoot)
	if result.Status != model.StatePresent || result.LogPath == "" || len(runner.Calls) != 1 {
		t.Fatalf("moon version result = %#v calls=%#v", result, runner.Calls)
	}

	runner.Err = errors.New("exit status 1")
	runner.Result.ExitCode = 1
	failed := cli.Run(context.Background(), paths.ProjectRoot, "moon-build", "moon", "build", "--target", "js")
	if failed.Status != model.StateFailed || failed.ExitCode != 1 {
		t.Fatalf("moon build failure result = %#v", failed)
	}
}

func TestLocalIndexSparseAndMalformed(t *testing.T) {
	reader := LocalIndexReader{}
	summary := reader.Parse([]byte(`{"name":"owner/a","version":"1.0.0","description":"","keywords":[]}
not-json
{"version":"missing-name"}`))
	if len(summary.Records) != 1 || summary.Malformed != 2 {
		t.Fatalf("summary = %#v", summary)
	}
	if summary.Records[0].Description.Status != model.StateMissing {
		t.Fatalf("description state = %s", summary.Records[0].Description.Status)
	}
}

func TestLocalIndexSourceEnvelopeComplete(t *testing.T) {
	root := t.TempDir()
	path := testutil.WriteFile(t, root, "index.jsonl", `{"name":"owner/a","version":"1.0.0"}`)
	reader := LocalIndexReader{
		FS:    platform.OSFS{},
		Clock: platform.FixedClock{Time: time.Date(2026, 6, 24, 1, 2, 3, 0, time.UTC)},
	}
	result := reader.Read(context.Background(), path)
	if result.Source != string(model.SourceLocalIndex) || result.Path == "" || result.FetchedAt != "2026-06-24T01:02:03Z" || result.ParseState != model.StatePresent || result.RawRef == "" {
		t.Fatalf("local index envelope = %#v", result)
	}
}

func TestLocalCacheCoreSourceStates(t *testing.T) {
	root := t.TempDir()
	present := testutil.WriteFile(t, root, "cache/source.zip", "payload")
	reader := LocalCacheReader{FS: platform.OSFS{}}
	if got := reader.CheckPath(context.Background(), model.SourceLocalCache, present); got.Status != model.StatePresent {
		t.Fatalf("present cache = %#v", got)
	}
	if got := reader.CheckPath(context.Background(), model.SourceCoreLocalSource, filepath.Join(root, "missing")); got.Status != model.StateUnavailable {
		t.Fatalf("missing cache = %#v", got)
	}
}

func TestProjectContextJSONAndTOML(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, root, "moon.mod.json", `{"name":"demo/project","deps":{"owner/a":"1.0.0"}}`)
	testutil.WriteFile(t, root, "moon.pkg", `target = "js"`)
	reader := ProjectReader{
		FS:    platform.OSFS{},
		Clock: platform.FixedClock{Time: time.Date(2026, 6, 24, 1, 2, 3, 0, time.UTC)},
	}
	result := reader.Read(context.Background(), root)
	if result.Status != model.StatePresent || result.Value == nil {
		t.Fatalf("project result = %#v", result)
	}
	if result.FetchedAt != "2026-06-24T01:02:03Z" || result.ParseState != model.StatePresent || result.RawRef == "" {
		t.Fatalf("project envelope = %#v", result)
	}
	if result.Value.Identity.Module != "demo/project" || result.Value.DeclaredTarget.Status != model.StatePresent || result.Value.ExistingDependencies.Status != model.StatePresent {
		t.Fatalf("project context = %#v", result.Value)
	}
}

func TestProjectContextPartialParseFailureObservable(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, root, "moon.mod.json", `{"name":"demo/project","deps":{"owner/a":"1.0.0"}}`)
	testutil.WriteFile(t, root, "moon.pkg.json", `{"target":`)
	reader := ProjectReader{
		FS:    platform.OSFS{},
		Clock: platform.FixedClock{Time: time.Date(2026, 6, 24, 1, 2, 3, 0, time.UTC)},
	}
	result := reader.Read(context.Background(), root)
	if result.Status != model.StateFailed || result.ParseState != model.StateFailed || result.Error == "" || result.Value == nil {
		t.Fatalf("project partial failure result = %#v", result)
	}
	if result.Value.Identity.Module != "demo/project" || result.Value.ModuleConfig.Status != model.StatePresent || result.Value.PackageConfig.Status != model.StateFailed {
		t.Fatalf("partial project context = %#v", result.Value)
	}
	if result.FetchedAt != "2026-06-24T01:02:03Z" || result.RawRef == "" {
		t.Fatalf("project partial envelope = %#v", result)
	}
}

func TestProbePathProjectBounded(t *testing.T) {
	paths := store.ResolvePaths(platform.Env{WorkingDir: t.TempDir(), XDGCacheHome: t.TempDir()})
	probe := store.ProbeStore{Paths: paths}
	dir := probe.ProbeDir("owner/module@1.0.0 js")
	if filepath.Dir(dir) != paths.Probes {
		t.Fatalf("probe dir %q not under %q", dir, paths.Probes)
	}
	logPath := probe.LogPath("owner/module@1.0.0 js", "moon check --target js")
	if filepath.Dir(filepath.Dir(logPath)) != dir {
		t.Fatalf("log path %q not under probe dir %q", logPath, dir)
	}
}
