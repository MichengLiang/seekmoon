// Package testutil contains fixtures shared by SeekMoon package tests.
package testutil

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/yumiaura/seekmoon/internal/platform"
)

// WriteFile writes a test fixture below root using private file permissions.
func WriteFile(t *testing.T, root, rel, body string) string {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}

// FakeRunner records Run requests and returns configured command results.
type FakeRunner struct {
	Result platform.RunResult
	Err    error
	Calls  []platform.RunRequest
}

// Run records the request and returns the configured fake result.
func (r *FakeRunner) Run(_ context.Context, request platform.RunRequest) (platform.RunResult, error) {
	r.Calls = append(r.Calls, request)
	result := r.Result
	if result.Command == nil {
		result.Command = request.Command
	}
	if result.CWD == "" {
		result.CWD = request.CWD
	}
	if result.LogPath == "" {
		result.LogPath = request.LogPath
	}
	return result, r.Err
}
