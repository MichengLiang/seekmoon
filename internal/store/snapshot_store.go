package store

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
)

// SnapshotStore persists sync snapshots.
type SnapshotStore struct {
	FS    platform.FS
	Paths Paths
}

// Path returns the path for a snapshot id.
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

// Latest returns the newest snapshot by sorted snapshot id.
func (s SnapshotStore) Latest(ctx context.Context) (model.Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return model.Snapshot{}, err
	}
	entries, err := os.ReadDir(s.Paths.Snapshots)
	if err != nil {
		return model.Snapshot{}, err
	}
	var ids []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		ids = append(ids, entry.Name()[:len(entry.Name())-len(".json")])
	}
	if len(ids) == 0 {
		return model.Snapshot{}, os.ErrNotExist
	}
	sort.Strings(ids)
	return s.Read(ctx, ids[len(ids)-1])
}
