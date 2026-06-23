package model

type SourceResolution struct {
	Module           string         `json:"module"`
	Version          string         `json:"version"`
	MoonFetch        SourceAttempt  `json:"moon_fetch"`
	SourceZip        SourceAttempt  `json:"source_zip"`
	LocalCache       SourceAttempt  `json:"local_cache"`
	CoreLocalSource  SourceAttempt  `json:"core_local_source"`
	RepositorySource SourceAttempt  `json:"repository_source"`
	SelectedSource   SelectedSource `json:"selected_source"`
	FilesSummary     FilesSummary   `json:"files_summary"`
}

type SourceAttempt struct {
	Status State  `json:"status"`
	Path   string `json:"path,omitempty"`
	URL    string `json:"url,omitempty"`
	Error  string `json:"error,omitempty"`
}

type SelectedSource struct {
	Method string `json:"method"`
	Path   string `json:"path,omitempty"`
	URL    string `json:"url,omitempty"`
}

type FilesSummary struct {
	MoonMod  bool `json:"moon_mod"`
	Readme   bool `json:"readme"`
	License  bool `json:"license"`
	Sources  int  `json:"sources"`
	Tests    int  `json:"tests"`
	Examples int  `json:"examples"`
	Benches  int  `json:"benches"`
}
