package model

type OutputMode string

const (
	OutputPretty OutputMode = "pretty"
	OutputJSON   OutputMode = "json"
	OutputJQ     OutputMode = "jq"
	OutputShape  OutputMode = "shape"
	OutputSchema OutputMode = "schema"
)

type SurfaceError struct {
	Command  string `json:"command"`
	Object   string `json:"object"`
	Source   string `json:"source"`
	State    State  `json:"state"`
	Meaning  string `json:"meaning"`
	Recovery string `json:"recovery,omitempty"`
	LogPath  string `json:"log_path,omitempty"`
}

type CommandInput struct {
	Command    string             `json:"command"`
	Query      string             `json:"query,omitempty"`
	Kind       string             `json:"kind,omitempty"`
	Target     string             `json:"target,omitempty"`
	Module     string             `json:"module,omitempty"`
	Package    string             `json:"package,omitempty"`
	Candidates []CandidateRequest `json:"candidates,omitempty"`
	Conclusion AdoptionConclusion `json:"conclusion,omitempty"`
	Note       string             `json:"note,omitempty"`
	Format     string             `json:"format,omitempty"`
	Source     string             `json:"source,omitempty"`
}

type CandidateRequest struct {
	Raw     string `json:"raw"`
	Number  int    `json:"number,omitempty"`
	Module  string `json:"module,omitempty"`
	Version string `json:"version,omitempty"`
}

type CommandEnvelope struct {
	Schema string       `json:"schema"`
	Input  CommandInput `json:"input"`
	Result any          `json:"result"`
}
