package model

type SnapshotRef struct {
	ID      string   `json:"id"`
	Sources []string `json:"sources"`
}

type Snapshot struct {
	ID         string              `json:"id"`
	CreatedAt  string              `json:"created_at"`
	Sources    []SourceResult[any] `json:"sources"`
	Statistics SnapshotStatistics  `json:"statistics"`
	Raw        map[string]any      `json:"raw,omitempty"`
}

type SnapshotStatistics struct {
	TotalModules   int `json:"total_modules"`
	TotalPackages  int `json:"total_packages"`
	TotalLines     int `json:"total_lines"`
	TotalDownloads int `json:"total_downloads"`
}
