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

func ParseState(value string) (State, error) {
	state := State(value)
	if state.IsValid() {
		return state, nil
	}
	return "", fmt.Errorf("unknown evidence state %q", value)
}

func (s State) IsValid() bool {
	_, ok := allStates[s]
	return ok
}

func (s State) Meaning() string {
	return allStates[s]
}

func (s State) IsAbsence() bool {
	return s == StateMissing || s == StateUnknown || s == StateUnavailable
}

func (s State) IsFailure() bool {
	return s == StateFailed
}

func (s State) MarshalJSON() ([]byte, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid evidence state %q", string(s))
	}
	return json.Marshal(string(s))
}

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

type Evidence[T any] struct {
	Status State   `json:"status"`
	Value  *T      `json:"value"`
	Source *string `json:"source"`
	Error  string  `json:"error,omitempty"`
}

func Present[T any](value T, source string) Evidence[T] {
	return Evidence[T]{Status: StatePresent, Value: &value, Source: sourcePtr(source)}
}

func Missing[T any](source string) Evidence[T] {
	return Evidence[T]{Status: StateMissing, Source: sourcePtr(source)}
}

func Unknown[T any]() Evidence[T] {
	return Evidence[T]{Status: StateUnknown}
}

func Failed[T any](source, message string) Evidence[T] {
	return Evidence[T]{Status: StateFailed, Source: sourcePtr(source), Error: message}
}

func Unavailable[T any](source string) Evidence[T] {
	return Evidence[T]{Status: StateUnavailable, Source: sourcePtr(source)}
}

func Derived[T any](value T, source string) Evidence[T] {
	return Evidence[T]{Status: StateDerived, Value: &value, Source: sourcePtr(source)}
}

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
	EvidenceString      = Evidence[string]
	EvidenceStringArray = Evidence[[]string]
	EvidenceInt         = Evidence[int]
	EvidenceObject      = Evidence[map[string]any]
)

type SourceLabel string

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
