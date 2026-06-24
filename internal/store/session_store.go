package store

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/MichengLiang/seekmoon/internal/model"
	"github.com/MichengLiang/seekmoon/internal/platform"
)

// SessionStore persists command session candidate maps.
type SessionStore struct {
	FS    platform.FS
	Paths Paths
}

// SessionCandidates records numbered candidates for a command session.
type SessionCandidates struct {
	SessionID  string                     `json:"session_id"`
	SnapshotID string                     `json:"snapshot_id"`
	Candidates map[int]model.CandidateRef `json:"candidates"`
}

// Path returns the path for a session id.
func (s SessionStore) Path(sessionID string) string {
	return filepath.Join(s.Paths.Sessions, SafeName(sessionID)+".json")
}

// WriteCandidates writes numbered candidates for a session.
func (s SessionStore) WriteCandidates(ctx context.Context, candidates SessionCandidates) error {
	data, err := json.MarshalIndent(candidates, "", "  ")
	if err != nil {
		return err
	}
	return AtomicWriteFile(ctx, s.FS, s.Path(candidates.SessionID), data)
}

// ReadCandidates reads numbered candidates for a session.
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
