package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// WriteGolden writes a named golden fixture below root/golden.
func WriteGolden(t *testing.T, root, name, body string) string {
	t.Helper()
	return WriteFile(t, root, filepath.Join("golden", name), body)
}

// ReadGolden reads a golden fixture for test assertions.
func ReadGolden(t *testing.T, path string) string {
	t.Helper()
	// #nosec G304 -- golden test paths are created by fixtures in the test tree.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadGolden: %v", err)
	}
	return string(data)
}
