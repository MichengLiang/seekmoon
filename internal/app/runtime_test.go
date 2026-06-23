package app

import (
	"testing"

	"github.com/yumiaura/seekmoon/internal/platform"
)

func TestNewRuntimeCompositionOrderArtifacts(t *testing.T) {
	root := t.TempDir()
	cache := t.TempDir()
	rt, err := NewRuntime(WithEnv(platform.Env{WorkingDir: root, XDGCacheHome: cache}))
	if err != nil {
		t.Fatalf("NewRuntime: %v", err)
	}
	if rt.Env.WorkingDir != root {
		t.Fatalf("env working dir = %q", rt.Env.WorkingDir)
	}
	if rt.Paths.ProjectRoot != root {
		t.Fatalf("paths project root = %q", rt.Paths.ProjectRoot)
	}
	if rt.HTTP == nil || rt.Clock == nil || rt.FS == nil || rt.Runner == nil {
		t.Fatalf("runtime host capabilities not initialized: %#v", rt)
	}
	if rt.Stores.Records.Paths.ProjectRoot != root {
		t.Fatalf("stores not composed from resolved paths")
	}
	if rt.Sources.Mooncakes.Fetcher.Client == nil || rt.Services.Sync.Snapshots.Paths.ProjectRoot != root {
		t.Fatalf("Batch B source/service registries not initialized")
	}
}
