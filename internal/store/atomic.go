package store

import (
	"context"
	"path/filepath"

	"github.com/yumiaura/seekmoon/internal/platform"
)

func AtomicWriteFile(ctx context.Context, fs platform.FS, path string, data []byte) error {
	if fs == nil {
		fs = platform.OSFS{}
	}
	if err := fs.MkdirAll(ctx, filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := fs.WriteFile(ctx, tmp, data, 0o644); err != nil {
		return err
	}
	if err := fs.Rename(ctx, tmp, path); err != nil {
		_ = fs.Remove(ctx, tmp)
		return err
	}
	return nil
}
