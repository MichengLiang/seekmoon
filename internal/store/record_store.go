package store

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
)

type RecordStore struct {
	FS    platform.FS
	Paths Paths
}

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
