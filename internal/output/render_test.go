package output

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/MichengLiang/seekmoon/internal/model"
)

func sampleSearchOutput() model.SearchOutput {
	return model.SearchOutput{
		Schema:   model.SchemaSearchResultsV1,
		Snapshot: model.SnapshotRef{ID: "2026-06-21T22:35:10+08:00", Sources: []string{"modules_api"}},
		Query:    model.SearchQuery{Text: "markdown", Kind: "library", Target: "js"},
		Results: []model.SearchResult{{
			Rank:        1,
			Module:      "mizchi/markdown",
			Version:     "0.6.2",
			Description: model.Present("Incremental Markdown parser and compiler", string(model.SourceModulesAPI)),
			Keywords:    model.Present([]string{"markdown"}, string(model.SourceModulesAPI)),
			License:     model.Present("MIT", string(model.SourceModulesAPI)),
			Repository:  model.Present("https://github.com/mizchi/markdown.mbt", string(model.SourceModulesAPI)),
			Downloads:   model.Present(7567, string(model.SourceManifestAPI)),
			BuildStatus: model.Present("success", string(model.SourceManifestAPI)),
			Target:      model.Unknown[map[string]any](),
			Match:       model.MatchEvidence{Fields: []string{"module", "description"}},
			SnapshotID:  "2026-06-21T22:35:10+08:00",
		}},
	}
}

func TestJSONProjectionCarriesSchemaAndEvidence(t *testing.T) {
	var out bytes.Buffer
	err := DefaultRenderer{}.Render(context.Background(), Request{
		Mode:   model.OutputJSON,
		Schema: model.SchemaSearchResultsV1,
		Writer: &out,
		Value:  sampleSearchOutput(),
	})
	if err != nil {
		t.Fatalf("Render JSON: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
		t.Fatalf("JSON output invalid: %v", err)
	}
	if decoded["schema"] != model.SchemaSearchResultsV1 {
		t.Fatalf("schema = %#v", decoded["schema"])
	}
	results := decoded["results"].([]any)
	first := results[0].(map[string]any)
	license := first["license"].(map[string]any)
	if license["status"] != string(model.StatePresent) || license["source"] != string(model.SourceModulesAPI) {
		t.Fatalf("license evidence = %#v", license)
	}
}

func TestJQProjectionConsumesJSONProjection(t *testing.T) {
	var out bytes.Buffer
	err := DefaultRenderer{}.Render(context.Background(), Request{
		Command:      "seekmoon search",
		Mode:         model.OutputJQ,
		Schema:       model.SchemaSearchResultsV1,
		JQExpression: ".results[].module",
		Writer:       &out,
		Value:        sampleSearchOutput(),
	})
	if err != nil {
		t.Fatalf("Render jq: %v", err)
	}
	if strings.TrimSpace(out.String()) != "mizchi/markdown" {
		t.Fatalf("jq output = %q", out.String())
	}
}

func TestJQFailureMapsToErrorSurface(t *testing.T) {
	var out bytes.Buffer
	err := DefaultRenderer{}.Render(context.Background(), Request{
		Command:      "seekmoon search --jq",
		Mode:         model.OutputJQ,
		Schema:       model.SchemaSearchResultsV1,
		JQExpression: ".results[",
		Writer:       &out,
		Value:        sampleSearchOutput(),
	})
	if err == nil {
		t.Fatalf("expected jq error")
	}
	if !strings.Contains(out.String(), "jq_expression") || !strings.Contains(out.String(), "embedded gojq interpreter") {
		t.Fatalf("error surface = %q", out.String())
	}
}

func TestShapeAndSchemaUseContractDefinitions(t *testing.T) {
	var shape bytes.Buffer
	if err := RenderShape(&shape, model.SchemaSearchResultsV1); err != nil {
		t.Fatalf("RenderShape: %v", err)
	}
	if !strings.Contains(shape.String(), "seekmoon.search-results.v1") || !strings.Contains(shape.String(), "results[]") {
		t.Fatalf("shape = %q", shape.String())
	}
	var schema bytes.Buffer
	if err := RenderSchema(&schema, model.SchemaSearchResultsV1); err != nil {
		t.Fatalf("RenderSchema: %v", err)
	}
	if !strings.Contains(schema.String(), "\"$schema\"") || !strings.Contains(schema.String(), model.SchemaSearchResultsV1) {
		t.Fatalf("schema = %q", schema.String())
	}
}

func TestPrettyTextOmitsRecoveryTutorial(t *testing.T) {
	var out bytes.Buffer
	if err := RenderPretty(&out, sampleSearchOutput()); err != nil {
		t.Fatalf("RenderPretty: %v", err)
	}
	if strings.Contains(out.String(), "recover") || strings.Contains(out.String(), "shape") {
		t.Fatalf("pretty text includes recovery/tutorial text: %q", out.String())
	}
}
