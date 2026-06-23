package store

import (
	"path/filepath"

	"github.com/yumiaura/seekmoon/internal/platform"
)

type SourceStore struct {
	Paths Paths
}

func (s SourceStore) SourceDir(id string) string {
	return filepath.Join(s.Paths.Sources, SafeName(id))
}

type Registry struct {
	Snapshots SnapshotStore
	Sessions  SessionStore
	Records   RecordStore
	Reports   ReportStore
	Cache     CacheStore
	Probes    ProbeStore
	Sources   SourceStore
}

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
