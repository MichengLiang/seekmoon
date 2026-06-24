package contract

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/MichengLiang/seekmoon/internal/model"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

func TestSearchShapeUsesStableSchemaID(t *testing.T) {
	shape := SearchResultsShape()
	if shape.Schema != model.SchemaSearchResultsV1 {
		t.Fatalf("shape schema = %q", shape.Schema)
	}
	if len(shape.Fields) == 0 || shape.Fields[0].Name != "schema" {
		t.Fatalf("shape fields not explicit: %#v", shape.Fields)
	}
}

func TestSearchSchemaCompiles(t *testing.T) {
	data, err := json.Marshal(SearchResultsSchema())
	if err != nil {
		t.Fatalf("Marshal schema: %v", err)
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("schema was not valid JSON: %v", err)
	}
	compiler := jsonschema.NewCompiler()
	compiler.DefaultDraft(jsonschema.Draft2020)
	if err := compiler.AddResource(model.SchemaSearchResultsV1, doc); err != nil {
		t.Fatalf("AddResource: %v", err)
	}
	if _, err := compiler.Compile(model.SchemaSearchResultsV1); err != nil {
		t.Fatalf("schema did not compile as JSON Schema: %v", err)
	}
}

func TestAdoptionSchemaCarriesConclusionEnum(t *testing.T) {
	schema := AdoptionRecordSchema()
	props := schema["properties"].(map[string]any)
	conclusion := props["conclusion"].(map[string]any)
	values := conclusion["enum"].([]string)
	if len(values) != 7 {
		t.Fatalf("conclusion enum length = %d", len(values))
	}
}
