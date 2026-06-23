package blackbox_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/tests/acceptance"
)

func TestA1SearchGeneratesCandidatesWithoutHandWrittenMooncakesURLs(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("search", "markdown", "--json")
	requireOK(t, result)
	requireContains(t, result.Output, "mizchi/markdown")
	if len(h.Fakes.SearchInputs) != 1 || h.Fakes.SearchInputs[0].Query != "markdown" {
		t.Fatalf("search input = %#v", h.Fakes.SearchInputs)
	}
	if strings.Contains(result.Output, "/api/v0/modules") {
		t.Fatalf("black-box command exposed hand-written API URL: %s", result.Output)
	}
}

func TestA2SearchResultCanBeReferencedBySessionLocalNumber(t *testing.T) {
	h := acceptance.NewHarness(t)
	requireOK(t, h.Run("search", "markdown", "--json"))
	result := h.Run("view", "1", "--json")
	requireOK(t, result)
	if len(h.Fakes.ViewInputs) != 1 || h.Fakes.ViewInputs[0].Candidate.Number != 1 {
		t.Fatalf("view input = %#v", h.Fakes.ViewInputs)
	}
	requireContains(t, result.Output, "mizchi/markdown")
}

func TestA3LibraryModuleAndSkillEntryUseDifferentCommandSurfaces(t *testing.T) {
	h := acceptance.NewHarness(t)
	requireOK(t, h.Run("search", "markdown", "--json"))
	requireOK(t, h.Run("skill", "search", "cowsay", "--json"))
	if len(h.Fakes.SearchInputs) != 1 || h.Fakes.SearchInputs[0].Query != "markdown" {
		t.Fatalf("library search input = %#v", h.Fakes.SearchInputs)
	}
	if len(h.Fakes.SkillSearchInputs) != 1 || h.Fakes.SkillSearchInputs[0].Query != "cowsay" {
		t.Fatalf("skill search input = %#v", h.Fakes.SkillSearchInputs)
	}
}

func TestA4ModuleProfileContainsManifestEvidenceAndPackageIndexState(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("view", "mizchi/markdown", "--json")
	requireOK(t, result)
	var object map[string]any
	decodeJSON(t, result.Output, &object)
	if object["schema"] != model.SchemaManifestProfileV1 {
		t.Fatalf("schema = %#v", object["schema"])
	}
	metadata := object["metadata"].(map[string]any)
	license := metadata["license"].(map[string]any)
	if license["status"] != string(model.StatePresent) {
		t.Fatalf("manifest license evidence = %#v", license)
	}
	raw := metadata["raw"].(map[string]any)
	if raw["index_state"] != string(model.StatePresent) || raw["package_count"].(float64) != 1 {
		t.Fatalf("package index state = %#v", raw)
	}
}

func TestA5PackageAPIProfileComesFromModuleIndexAndPackageData(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("api", "mizchi/markdown", "--package", "mizchi/markdown/src/api", "--json")
	requireOK(t, result)
	var object map[string]any
	decodeJSON(t, result.Output, &object)
	types := object["types"].([]any)
	firstType := types[0].(map[string]any)
	if firstType["signature"] != "<a>Document</a>" {
		t.Fatalf("raw signature was not preserved: %#v", firstType)
	}
	if firstType["plain_signature"].(map[string]any)["status"] != string(model.StateDerived) {
		t.Fatalf("plain signature should be derived: %#v", firstType["plain_signature"])
	}
	if len(h.Fakes.APIInputs) != 1 || h.Fakes.APIInputs[0].Package != "mizchi/markdown/src/api" {
		t.Fatalf("api input = %#v", h.Fakes.APIInputs)
	}
}

func TestA6PublishedSourceCanBeFetchedOrLocatedThroughSourceResolution(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("source", "mizchi/markdown@0.6.2", "--json")
	requireOK(t, result)
	var object map[string]any
	decodeJSON(t, result.Output, &object)
	if object["selected_source"].(map[string]any)["method"] != "source_zip" {
		t.Fatalf("source resolution = %#v", object)
	}
	if object["source_zip"].(map[string]any)["status"] != string(model.StatePresent) {
		t.Fatalf("source zip attempt = %#v", object["source_zip"])
	}
}

func TestA7TargetSupportRemainsUnknownBeforeEvidence(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("search", "markdown", "--target", "js", "--json")
	requireOK(t, result)
	var object map[string]any
	decodeJSON(t, result.Output, &object)
	first := object["results"].([]any)[0].(map[string]any)
	target := first["target"].(map[string]any)
	if target["status"] != string(model.StateUnknown) {
		t.Fatalf("target state = %#v", target)
	}
}

func TestA8ProbeProducesLocalDerivedEvidenceAndDoesNotMutateUpstreamFacts(t *testing.T) {
	h := acceptance.NewHarness(t)
	before := h.Run("view", "mizchi/markdown", "--json")
	requireOK(t, before)
	probe := h.Run("probe", "1", "--target", "js", "--json")
	requireOK(t, probe)
	after := h.Run("view", "mizchi/markdown", "--json")
	requireOK(t, after)
	requireContains(t, probe.Output, "moon_build_target")
	if before.Output != after.Output {
		t.Fatalf("probe mutated upstream view output\nbefore=%s\nafter=%s", before.Output, after.Output)
	}
}

func TestA9AdoptionDecisionPersistsAsRecordWithEvidenceRefs(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("record", "1", "--conclusion", "continue-verification", "--json")
	requireOK(t, result)
	var object map[string]any
	decodeJSON(t, result.Output, &object)
	if object["conclusion"] != string(model.ConclusionContinueVerification) {
		t.Fatalf("record conclusion = %#v", object["conclusion"])
	}
	if len(object["evidence_refs"].([]any)) == 0 {
		t.Fatalf("record missing evidence refs: %#v", object)
	}
}

func TestA10ReportListsOnlyActuallyUsedSources(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("report", "--format", "markdown", "--json")
	requireOK(t, result)
	requireContains(t, result.Output, string(model.SourceManifestAPI))
	if strings.Contains(result.Output, string(model.SourceRepositoryAPI)) {
		t.Fatalf("report listed unobserved repository source: %s", result.Output)
	}
}

func TestA11PrettyTextIsLowNoiseAndNotParsingInterface(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("search", "markdown")
	requireOK(t, result)
	requireContains(t, result.Output, "Search:")
	requireContains(t, result.Output, "mizchi/markdown")
	if strings.Contains(result.Output, "{") || strings.Contains(result.Output, "schema") {
		t.Fatalf("pretty output looks like machine contract output: %s", result.Output)
	}
}

func TestA12JSONOutputContainsSchemaIDAndEvidenceStates(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("search", "markdown", "--json")
	requireOK(t, result)
	var object map[string]any
	decodeJSON(t, result.Output, &object)
	if object["schema"] != model.SchemaSearchResultsV1 {
		t.Fatalf("schema id = %#v", object["schema"])
	}
	first := object["results"].([]any)[0].(map[string]any)
	if first["description"].(map[string]any)["status"] != string(model.StatePresent) {
		t.Fatalf("description evidence = %#v", first["description"])
	}
}

func TestA13BuiltInJQEvaluatesCommandJSONOutput(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("search", "markdown", "--jq", ".results[].module")
	requireOK(t, result)
	if strings.TrimSpace(result.Output) != "mizchi/markdown" {
		t.Fatalf("jq output = %q", result.Output)
	}
}

func TestA14ShapeExplainsJSONFieldsWithoutRealQueryExecution(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("search", "--shape")
	requireOK(t, result)
	requireContains(t, result.Output, model.SchemaSearchResultsV1)
	requireContains(t, result.Output, "results[]")
	if len(h.Fakes.SearchInputs) != 0 {
		t.Fatalf("shape mode executed search: %#v", h.Fakes.SearchInputs)
	}
}

func TestA15SchemaProvidesJSONSchemaForStrictConsumers(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("search", "--schema")
	requireOK(t, result)
	var object map[string]any
	decodeJSON(t, result.Output, &object)
	if object["$id"] != model.SchemaSearchResultsV1 || object["type"] != "object" {
		t.Fatalf("schema output = %#v", object)
	}
}

func TestA16CommandFailureUsesErrorSurfaceWithSourceStateMeaningAndRecovery(t *testing.T) {
	h := acceptance.NewHarness(t)
	result := h.Run("api", "mizchi/markdown", "--package", "mizchi/markdown/missing")
	if result.Code != 1 {
		t.Fatalf("exit code = %d output=%s", result.Code, result.Output)
	}
	for _, text := range []string{"source", "module_index.json", "state", "failed", "meaning", "recover"} {
		requireContains(t, result.Output, text)
	}
}

func requireOK(t *testing.T, result acceptance.CommandResult) {
	t.Helper()
	if result.Code != 0 {
		t.Fatalf("exit code = %d want 0 output=%s", result.Code, result.Output)
	}
}

func requireContains(t *testing.T, value, want string) {
	t.Helper()
	if !strings.Contains(value, want) {
		t.Fatalf("output missing %q: %s", want, value)
	}
}

func decodeJSON(t *testing.T, data string, target any) {
	t.Helper()
	if err := json.Unmarshal([]byte(data), target); err != nil {
		t.Fatalf("decode JSON: %v\n%s", err, data)
	}
}
