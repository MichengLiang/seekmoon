package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseStateClosedVocabulary(t *testing.T) {
	valid := []State{StatePresent, StateMissing, StateUnknown, StateFailed, StateUnavailable, StateDerived}
	for _, want := range valid {
		got, err := ParseState(string(want))
		if err != nil {
			t.Fatalf("ParseState(%q): %v", want, err)
		}
		if got != want {
			t.Fatalf("ParseState(%q) = %q", want, got)
		}
	}

	if _, err := ParseState("not-fetched"); err == nil {
		t.Fatal("ParseState accepted a state outside the WBS vocabulary")
	}
}

func TestStateSemanticsDoNotCollapseAbsenceOrFailure(t *testing.T) {
	if StateMissing.Meaning() == StateUnknown.Meaning() {
		t.Fatal("missing and unknown must retain distinct meanings")
	}
	if StateUnavailable.IsFailure() {
		t.Fatal("unavailable is optional-source absence, not source failure")
	}
	if !StateFailed.IsFailure() {
		t.Fatal("failed must remain the explicit action failure state")
	}
	if !StateMissing.IsAbsence() || !StateUnknown.IsAbsence() || !StateUnavailable.IsAbsence() {
		t.Fatal("absence helpers must preserve missing/unknown/unavailable")
	}
}

func TestEvidenceWrapperJSONRoundTrip(t *testing.T) {
	original := Present("MIT", string(SourceModulesAPI))
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got EvidenceString
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got.Status != StatePresent || got.Value == nil || *got.Value != "MIT" || got.Source == nil || *got.Source != string(SourceModulesAPI) {
		t.Fatalf("round trip mismatch: %#v", got)
	}
}

func TestEvidenceWrapperPreservesNullSource(t *testing.T) {
	data, err := json.Marshal(Unknown[string]())
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !strings.Contains(string(data), `"source":null`) {
		t.Fatalf("unknown evidence must emit nullable source field, got %s", data)
	}
	var got EvidenceString
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got.Source != nil {
		t.Fatalf("unknown evidence source = %q, want nil", *got.Source)
	}
}

func TestFailedEvidenceRequiresError(t *testing.T) {
	source := string(SourceManifestAPI)
	if err := (EvidenceString{Status: StateFailed, Source: &source}).Validate(); err == nil {
		t.Fatal("failed evidence without error should not validate")
	}
}
