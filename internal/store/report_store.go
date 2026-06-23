package store

import (
	"context"
	"path/filepath"

	"github.com/yumiaura/seekmoon/internal/platform"
)

type ReportStore struct {
	FS    platform.FS
	Paths Paths
}

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
