package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func WriteGolden(t *testing.T, root, name, body string) string {
	t.Helper()
	return WriteFile(t, root, filepath.Join("golden", name), body)
}

func ReadGolden(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadGolden: %v", err)
	}
	return string(data)
}
