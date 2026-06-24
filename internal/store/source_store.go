package store

import (
	"path/filepath"

	"github.com/MichengLiang/seekmoon/internal/platform"
)

// SourceStore resolves local source inspection paths.
type SourceStore struct {
	Paths Paths
}

// SourceDir returns the local source directory for an id.
func (s SourceStore) SourceDir(id string) string {
	return filepath.Join(s.Paths.Sources, SafeName(id))
}

// Registry groups all stores used by the runtime.
type Registry struct {
	Snapshots SnapshotStore
	Sessions  SessionStore
	Records   RecordStore
	Reports   ReportStore
	Cache     CacheStore
	Probes    ProbeStore
	Sources   SourceStore
}

// NewRegistry creates stores for the supplied filesystem and paths.
func NewRegistry(fs platform.FS, paths Paths) Registry {
	return Registry{
		Snapshots: SnapshotStore{FS: fs, Paths: paths},
		Sessions:  SessionStore{FS: fs, Paths: paths},
		Records:   RecordStore{FS: fs, Paths: paths},
		Reports:   ReportStore{FS: fs, Paths: paths},
		Cache:     CacheStore{FS: fs, Paths: paths},
		Probes:    ProbeStore{Paths: paths},
		Sources:   SourceStore{Paths: paths},
	}
}
