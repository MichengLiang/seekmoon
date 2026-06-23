package store

import (
	"path/filepath"
)

type ProbeStore struct {
	Paths Paths
}

func (s ProbeStore) ProbeDir(id string) string {
	return filepath.Join(s.Paths.Probes, SafeName(id))
}

func (s ProbeStore) LogPath(probeID, name string) string {
	return filepath.Join(s.ProbeDir(probeID), "logs", SafeName(name)+".log")
}
