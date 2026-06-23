package source

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
)

type LocalCacheReader struct {
	FS platform.FS
}

func (r LocalCacheReader) CheckPath(ctx context.Context, label model.SourceLabel, path string) model.SourceAttempt {
	fs := r.FS
	if fs == nil {
		fs = platform.OSFS{}
	}
	_, err := fs.ReadFile(ctx, path)
	attempt := model.SourceAttempt{Path: filepath.Clean(path)}
	if err == nil {
		attempt.Status = model.StatePresent
		return attempt
	}
	if errors.Is(err, context.Canceled) {
		attempt.Status = model.StateFailed
		attempt.Error = err.Error()
		return attempt
	}
	attempt.Status = model.StateUnavailable
	attempt.Error = err.Error()
	_ = label
	return attempt
}
