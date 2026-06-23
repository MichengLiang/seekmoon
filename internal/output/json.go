package output

import (
	"encoding/json"
	"fmt"
	"io"
)

func JSONProjection(schema string, value any) (any, error) {
	if schema == "" {
		return nil, fmt.Errorf("schema id is required")
	}
	if object, ok := value.(interface{ SchemaID() string }); ok && object.SchemaID() != "" {
		schema = object.SchemaID()
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var decoded any
	if err := json.Unmarshal(data, &decoded); err != nil {
		return nil, err
	}
	if object, ok := decoded.(map[string]any); ok {
		if _, exists := object["schema"]; !exists {
			object["schema"] = schema
		}
		return object, nil
	}
	return map[string]any{"schema": schema, "result": decoded}, nil
}

func RenderJSON(writer io.Writer, schema string, value any) error {
	projected, err := JSONProjection(schema, value)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(projected)
}
