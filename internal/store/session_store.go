package store

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
)

type SessionStore struct {
	FS    platform.FS
	Paths Paths
}

type SessionCandidates struct {
	SessionID  string                     `json:"session_id"`
	SnapshotID string                     `json:"snapshot_id"`
	Candidates map[int]model.CandidateRef `json:"candidates"`
}

func (s SessionStore) Path(sessionID string) string {
	return filepath.Join(s.Paths.Sessions, SafeName(sessionID)+".json")
}

func (s SessionStore) WriteCandidates(ctx context.Context, candidates SessionCandidates) error {
	data, err := json.MarshalIndent(candidates, "", "  ")
	if err != nil {
		return err
	}
	return AtomicWriteFile(ctx, s.FS, s.Path(candidates.SessionID), data)
}

func (s SessionStore) ReadCandidates(ctx context.Context, sessionID string) (SessionCandidates, error) {
	fs := s.FS
	if fs == nil {
		fs = platform.OSFS{}
	}
	data, err := fs.ReadFile(ctx, s.Path(sessionID))
	if err != nil {
		return SessionCandidates{}, err
	}
	var candidates SessionCandidates
	if err := json.Unmarshal(data, &candidates); err != nil {
		return SessionCandidates{}, err
	}
	return candidates, nil
}
