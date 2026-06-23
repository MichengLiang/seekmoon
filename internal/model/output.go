package model

// OutputMode selects how command output is rendered.
type OutputMode string

// Output modes supported by the renderer.
const (
	OutputPretty OutputMode = "pretty"
	OutputJSON   OutputMode = "json"
	OutputJQ     OutputMode = "jq"
	OutputShape  OutputMode = "shape"
	OutputSchema OutputMode = "schema"
)

// SurfaceError is the structured error shape shown on command failure.
type SurfaceError struct {
	Command  string `json:"command"`
	Object   string `json:"object"`
	Source   string `json:"source"`
	State    State  `json:"state"`
	Meaning  string `json:"meaning"`
	Recovery string `json:"recovery,omitempty"`
	LogPath  string `json:"log_path,omitempty"`
}

// SurfaceFailure carries a SurfaceError through the error interface.
type SurfaceFailure struct {
	Value SurfaceError
}

func (e SurfaceFailure) Error() string {
	return e.Value.Meaning
}

// SurfaceError returns the structured command failure payload.
func (e SurfaceFailure) SurfaceError() SurfaceError {
	return e.Value
}

// CommandInput records the normalized input values for an output envelope.
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

// CandidateRequest identifies a candidate by number or module coordinate.
type CandidateRequest struct {
	Raw     string `json:"raw"`
	Number  int    `json:"number,omitempty"`
	Module  string `json:"module,omitempty"`
	Version string `json:"version,omitempty"`
}

// CommandEnvelope wraps a command result with schema and input metadata.
type CommandEnvelope struct {
	Schema string       `json:"schema"`
	Input  CommandInput `json:"input"`
	Result any          `json:"result"`
}

// EnvironmentStatus captures local toolchain, path, network, and project evidence.
type EnvironmentStatus struct {
	Schema    string                       `json:"schema"`
	Toolchain map[string]CommandResult     `json:"toolchain"`
	Commands  map[string]CommandResult     `json:"commands"`
	Paths     map[string]EvidenceString    `json:"paths"`
	Network   map[string]EvidenceString    `json:"network"`
	Project   SourceResult[ProjectContext] `json:"project"`
}

// Comparison aligns evidence fields across candidate modules.
type Comparison struct {
	Schema     string            `json:"schema"`
	Candidates []CandidateRef    `json:"candidates"`
	Fields     []ComparisonField `json:"fields"`
}

// ComparisonField records one compared field across candidates.
type ComparisonField struct {
	Name   string            `json:"name"`
	Values map[string]string `json:"values"`
}

// RawEnvelope preserves an upstream payload with source status metadata.
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
