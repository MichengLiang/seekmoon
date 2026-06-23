package model

// SnapshotRef is the compact reference to a stored snapshot.
type SnapshotRef struct {
	ID      string   `json:"id"`
	Sources []string `json:"sources"`
}

// Snapshot captures synced source evidence at one point in time.
type Snapshot struct {
	ID         string              `json:"id"`
	CreatedAt  string              `json:"created_at"`
	Sources    []SourceResult[any] `json:"sources"`
	Statistics SnapshotStatistics  `json:"statistics"`
	Raw        map[string]any      `json:"raw,omitempty"`
}

// SnapshotStatistics records registry-wide counts from the snapshot.
type SnapshotStatistics struct {
	TotalModules   int `json:"total_modules"`
	TotalPackages  int `json:"total_packages"`
	TotalLines     int `json:"total_lines"`
	TotalDownloads int `json:"total_downloads"`
}
