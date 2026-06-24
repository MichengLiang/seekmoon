package store

import (
	"context"
	"path/filepath"

	"github.com/MichengLiang/seekmoon/internal/platform"
)

// CacheStore resolves cache paths and writes cached payloads.
type CacheStore struct {
	FS    platform.FS
	Paths Paths
}

// MooncakesPath returns the cache path for a Mooncakes API key.
func (s CacheStore) MooncakesPath(key string) string {
	return filepath.Join(s.Paths.Mooncakes, SafeName(key)+".json")
}

// AssetPath returns the cache path for a downloaded asset key.
func (s CacheStore) AssetPath(key string) string {
	return filepath.Join(s.Paths.Assets, SafeName(key))
}

// GitHubPath returns the cache path for a GitHub API key.
func (s CacheStore) GitHubPath(key string) string {
	return filepath.Join(s.Paths.GitHub, SafeName(key)+".json")
}

func (s CacheStore) Write(ctx context.Context, path string, data []byte) error {
	return AtomicWriteFile(ctx, s.FS, path, data)
}
