package store

import (
	"context"
	"path/filepath"

	"github.com/MichengLiang/seekmoon/internal/platform"
)

// ReportStore writes generated reports.
type ReportStore struct {
	FS    platform.FS
	Paths Paths
}

// Path returns the report path for a name and extension.
func (s ReportStore) Path(name string, ext string) string {
	if ext == "" {
		ext = ".md"
	}
	return filepath.Join(s.Paths.Reports, SafeName(name)+ext)
}

func (s ReportStore) Write(ctx context.Context, name, ext string, data []byte) (string, error) {
	path := s.Path(name, ext)
	return path, AtomicWriteFile(ctx, s.FS, path, data)
}
