package store

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
)

// RecordStore persists adoption records.
type RecordStore struct {
	FS    platform.FS
	Paths Paths
}

// Path returns the path for a record id.
func (s RecordStore) Path(recordID string) string {
	return filepath.Join(s.Paths.Records, SafeName(recordID)+".json")
}

func (s RecordStore) Write(ctx context.Context, record model.AdoptionRecord) error {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	return AtomicWriteFile(ctx, s.FS, s.Path(record.RecordID), data)
}

func (s RecordStore) Read(ctx context.Context, recordID string) (model.AdoptionRecord, error) {
	fs := s.FS
	if fs == nil {
		fs = platform.OSFS{}
	}
	data, err := fs.ReadFile(ctx, s.Path(recordID))
	if err != nil {
		return model.AdoptionRecord{}, err
	}
	var record model.AdoptionRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return model.AdoptionRecord{}, err
	}
	return record, nil
}

// List returns all stored adoption records.
func (s RecordStore) List(ctx context.Context) ([]model.AdoptionRecord, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(s.Paths.Records)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	records := make([]model.AdoptionRecord, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := platformFS(s.FS).ReadFile(ctx, filepath.Join(s.Paths.Records, entry.Name()))
		if err != nil {
			return nil, err
		}
		var record model.AdoptionRecord
		if err := json.Unmarshal(data, &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func platformFS(fs platform.FS) platform.FS {
	if fs == nil {
		return platform.OSFS{}
	}
	return fs
}
