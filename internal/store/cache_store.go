package store

import (
	"context"
	"path/filepath"

	"github.com/yumiaura/seekmoon/internal/platform"
)

type CacheStore struct {
	FS    platform.FS
	Paths Paths
}

func (s CacheStore) MooncakesPath(key string) string {
	return filepath.Join(s.Paths.Mooncakes, SafeName(key)+".json")
}

func (s CacheStore) AssetPath(key string) string {
	return filepath.Join(s.Paths.Assets, SafeName(key))
}

func (s CacheStore) GitHubPath(key string) string {
	return filepath.Join(s.Paths.GitHub, SafeName(key)+".json")
}

func (s CacheStore) Write(ctx context.Context, path string, data []byte) error {
	return AtomicWriteFile(ctx, s.FS, path, data)
}
