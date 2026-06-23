package store

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
)

type SnapshotStore struct {
	FS    platform.FS
	Paths Paths
}

func (s SnapshotStore) Path(id string) string {
	return filepath.Join(s.Paths.Snapshots, SafeName(id)+".json")
}

func (s SnapshotStore) Write(ctx context.Context, snapshot model.Snapshot) error {
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}
	return AtomicWriteFile(ctx, s.FS, s.Path(snapshot.ID), data)
}

func (s SnapshotStore) Read(ctx context.Context, id string) (model.Snapshot, error) {
	fs := s.FS
	if fs == nil {
		fs = platform.OSFS{}
	}
	data, err := fs.ReadFile(ctx, s.Path(id))
	if err != nil {
		return model.Snapshot{}, err
	}
	var snapshot model.Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return model.Snapshot{}, err
	}
	return snapshot, nil
}
