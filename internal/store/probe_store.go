package store

import (
	"path/filepath"
)

// ProbeStore resolves probe directories and log paths.
type ProbeStore struct {
	Paths Paths
}

// ProbeDir returns the directory for a probe id.
func (s ProbeStore) ProbeDir(id string) string {
	return filepath.Join(s.Paths.Probes, SafeName(id))
}

// LogPath returns a named log path within a probe directory.
func (s ProbeStore) LogPath(probeID, name string) string {
	return filepath.Join(s.ProbeDir(probeID), "logs", SafeName(name)+".log")
}
