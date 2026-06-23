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

type SurfaceFailure struct {
	Value SurfaceError
}

func (e SurfaceFailure) Error() string {
	return e.Value.Meaning
}

func (e SurfaceFailure) SurfaceError() SurfaceError {
	return e.Value
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

type EnvironmentStatus struct {
	Schema    string                       `json:"schema"`
	Toolchain map[string]CommandResult     `json:"toolchain"`
	Commands  map[string]CommandResult     `json:"commands"`
	Paths     map[string]EvidenceString    `json:"paths"`
	Network   map[string]EvidenceString    `json:"network"`
	Project   SourceResult[ProjectContext] `json:"project"`
}

type Comparison struct {
	Schema     string            `json:"schema"`
	Candidates []CandidateRef    `json:"candidates"`
	Fields     []ComparisonField `json:"fields"`
}

type ComparisonField struct {
	Name   string            `json:"name"`
	Values map[string]string `json:"values"`
}

type RawEnvelope struct {
	Schema    string `json:"schema"`
	Source    string `json:"source"`
	URL       string `json:"url,omitempty"`
	Path      string `json:"path,omitempty"`
	FetchedAt string `json:"fetched_at,omitempty"`
	Status    State  `json:"status"`
	RawRef    string `json:"raw_ref,omitempty"`
	Payload   any    `json:"payload,omitempty"`
	Error     string `json:"error,omitempty"`
}
