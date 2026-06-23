// Package model defines SeekMoon canonical objects and evidence state vocabulary.
package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

// State is the closed evidence vocabulary shared by JSON output, schemas,
// stores, records, reports, and error surfaces.
type State string

// Evidence states used by source readers and output contracts.
const (
	StatePresent     State = "present"
	StateMissing     State = "missing"
	StateUnknown     State = "unknown"
	StateFailed      State = "failed"
	StateUnavailable State = "unavailable"
	StateDerived     State = "derived"
)

var allStates = map[State]string{
	StatePresent:     "source supplied a valid value",
	StateMissing:     "source position exists but has no value",
	StateUnknown:     "loaded evidence cannot answer the current question",
	StateFailed:      "request, command, or parse action executed and failed",
	StateUnavailable: "optional source is absent for the current object",
	StateDerived:     "SeekMoon calculated the value from current facts",
}

// ParseState validates and returns an evidence state.
func ParseState(value string) (State, error) {
	state := State(value)
	if state.IsValid() {
		return state, nil
	}
	return "", fmt.Errorf("unknown evidence state %q", value)
}

// IsValid reports whether the state belongs to the closed evidence vocabulary.
func (s State) IsValid() bool {
	_, ok := allStates[s]
	return ok
}

// Meaning returns the human-readable meaning of a state.
func (s State) Meaning() string {
	return allStates[s]
}

// IsAbsence reports whether the state means no usable value is available.
func (s State) IsAbsence() bool {
	return s == StateMissing || s == StateUnknown || s == StateUnavailable
}

// IsFailure reports whether the state represents a failed source action.
func (s State) IsFailure() bool {
	return s == StateFailed
}

// MarshalJSON rejects unknown states before writing JSON.
func (s State) MarshalJSON() ([]byte, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid evidence state %q", string(s))
	}
	return json.Marshal(string(s))
}

// UnmarshalJSON validates incoming state strings.
func (s *State) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	state, err := ParseState(value)
	if err != nil {
		return err
	}
	*s = state
	return nil
}

// Evidence wraps a value with source, status, and failure metadata.
type Evidence[T any] struct {
	Status State   `json:"status"`
	Value  *T      `json:"value"`
	Source *string `json:"source"`
	Error  string  `json:"error,omitempty"`
}

// Present returns evidence with an observed source value.
func Present[T any](value T, source string) Evidence[T] {
	return Evidence[T]{Status: StatePresent, Value: &value, Source: sourcePtr(source)}
}

// Missing returns evidence for a known field with no source value.
func Missing[T any](source string) Evidence[T] {
	return Evidence[T]{Status: StateMissing, Source: sourcePtr(source)}
}

// Unknown returns evidence for an unanswered question.
func Unknown[T any]() Evidence[T] {
	return Evidence[T]{Status: StateUnknown}
}

// Failed returns evidence for an attempted source action that failed.
func Failed[T any](source, message string) Evidence[T] {
	return Evidence[T]{Status: StateFailed, Source: sourcePtr(source), Error: message}
}

// Unavailable returns evidence for an optional absent source.
func Unavailable[T any](source string) Evidence[T] {
	return Evidence[T]{Status: StateUnavailable, Source: sourcePtr(source)}
}

// Derived returns evidence calculated from current facts.
func Derived[T any](value T, source string) Evidence[T] {
	return Evidence[T]{Status: StateDerived, Value: &value, Source: sourcePtr(source)}
}

// Validate checks evidence status/value consistency.
func (e Evidence[T]) Validate() error {
	if !e.Status.IsValid() {
		return fmt.Errorf("invalid evidence status %q", e.Status)
	}
	if e.Status == StateFailed && e.Error == "" {
		return errors.New("failed evidence requires error")
	}
	if (e.Status == StatePresent || e.Status == StateDerived) && e.Value == nil {
		return fmt.Errorf("%s evidence requires value", e.Status)
	}
	return nil
}

func sourcePtr(source string) *string {
	if source == "" {
		return nil
	}
	return &source
}

type (
	// EvidenceString is string evidence.
	EvidenceString = Evidence[string]
	// EvidenceStringArray is string-slice evidence.
	EvidenceStringArray = Evidence[[]string]
	// EvidenceInt is integer evidence.
	EvidenceInt = Evidence[int]
	// EvidenceObject is object-map evidence.
	EvidenceObject = Evidence[map[string]any]
)

// SourceLabel identifies a source family in evidence records.
type SourceLabel string

// Source labels used by source readers and stores.
const (
	SourceModulesAPI      SourceLabel = "modules_api"
	SourceStatisticsAPI   SourceLabel = "statistics_api"
	SourceManifestAPI     SourceLabel = "manifest_api"
	SourceSkillsAPI       SourceLabel = "skills_api"
	SourceModuleIndex     SourceLabel = "module_index_asset"
	SourcePackageData     SourceLabel = "package_data_asset"
	SourceResourceAsset   SourceLabel = "resource_asset"
	SourceSourceZip       SourceLabel = "source_zip"
	SourceMoonCLI         SourceLabel = "moon_cli"
	SourceLocalIndex      SourceLabel = "local_index"
	SourceLocalCache      SourceLabel = "local_cache"
	SourceCoreLocalSource SourceLabel = "core_local_source"
	SourceRepositoryAPI   SourceLabel = "repository_api"
	SourceProjectContext  SourceLabel = "project_context"
	SourceDerived         SourceLabel = "derived"
)

// SourceResult records value, status, and raw-reference data for one source.
type SourceResult[T any] struct {
	Source     string `json:"source"`
	URL        string `json:"url,omitempty"`
	Path       string `json:"path,omitempty"`
	FetchedAt  string `json:"fetched_at,omitempty"`
	Status     State  `json:"status"`
	ParseState State  `json:"parse_state,omitempty"`
	RawRef     string `json:"raw_ref,omitempty"`
	Error      string `json:"error,omitempty"`
	Value      *T     `json:"value,omitempty"`
}
