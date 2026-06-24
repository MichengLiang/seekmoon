package journey_test

import (
	"strings"
	"testing"

	"github.com/MichengLiang/seekmoon/tests/acceptance"
)

func TestLibraryDiscoveryJourneySearchViewAPIProbeRecordReport(t *testing.T) {
	h := acceptance.NewHarness(t)
	for _, step := range []struct {
		name string
		args []string
		want string
	}{
		{name: "search", args: []string{"search", "markdown"}, want: "mizchi/markdown"},
		{name: "view", args: []string{"view", "1"}, want: "docs"},
		{name: "api", args: []string{"api", "1", "--package", "mizchi/markdown/src/api"}, want: "Document"},
		{name: "source", args: []string{"source", "1"}, want: "source_zip"},
		{name: "probe", args: []string{"probe", "1", "--target", "js"}, want: "verifiable"},
		{name: "record", args: []string{"record", "1", "--conclusion", "continue-verification"}, want: "continue-verification"},
		{name: "report", args: []string{"report", "--format", "markdown"}, want: "Assess MoonBit dependency"},
	} {
		t.Run(step.name, func(t *testing.T) {
			result := h.Run(step.args...)
			if result.Code != 0 {
				t.Fatalf("exit code = %d output=%s", result.Code, result.Output)
			}
			if !strings.Contains(result.Output, step.want) {
				t.Fatalf("output missing %q: %s", step.want, result.Output)
			}
		})
	}
}

func TestSkillDiscoveryJourneyUsesSkillSurfaceAndRunwasmProfile(t *testing.T) {
	h := acceptance.NewHarness(t)
	search := h.Run("skill", "search", "cowsay", "--json")
	if search.Code != 0 || !strings.Contains(search.Output, "Yoorkin/cowsay") {
		t.Fatalf("skill search result = %#v", search)
	}
	view := h.Run("skill", "view", "Yoorkin/cowsay", "--json")
	if view.Code != 0 || !strings.Contains(view.Output, "runwasm_coordinate") || !strings.Contains(view.Output, "Yoorkin/cowsay/cowsay@0.1.0") {
		t.Fatalf("skill view result = %#v", view)
	}
	if len(h.Fakes.SearchInputs) != 0 {
		t.Fatalf("skill journey should not use library search: %#v", h.Fakes.SearchInputs)
	}
}

func TestPipelineJourneyCombinesJSONShapeSchemaAndJQ(t *testing.T) {
	h := acceptance.NewHarness(t)
	for _, args := range [][]string{
		{"search", "markdown", "--json"},
		{"search", "--shape"},
		{"search", "--schema"},
		{"search", "markdown", "--jq", ".results[].license.value"},
	} {
		result := h.Run(args...)
		if result.Code != 0 {
			t.Fatalf("%v exit code = %d output=%s", args, result.Code, result.Output)
		}
		if strings.TrimSpace(result.Output) == "" {
			t.Fatalf("%v produced empty output", args)
		}
	}
}

func TestFailureRecoveryJourneySurfacesActionableAPIPathFailure(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("api", "mizchi/markdown", "--package", "mizchi/markdown/missing")
	if result.Code != 1 {
		t.Fatalf("exit code = %d output=%s", result.Code, result.Output)
	}
	for _, want := range []string{"module_index.json", "failed", "package path is not present", "choose one"} {
		if !strings.Contains(result.Output, want) {
			t.Fatalf("failure output missing %q: %s", want, result.Output)
		}
	}
}
