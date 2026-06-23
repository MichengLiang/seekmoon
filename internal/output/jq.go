package output

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/itchyny/gojq"
	"github.com/yumiaura/seekmoon/internal/model"
)

func RenderJQ(ctx context.Context, writer io.Writer, schema string, value any, expression, command string) error {
	if expression == "" {
		return renderJQError(writer, jqSurfaceError(command, "jq expression is required"))
	}
	projected, err := JSONProjection(schema, value)
	if err != nil {
		return renderJQError(writer, jqSurfaceError(command, err.Error()))
	}
	query, err := gojq.Parse(expression)
	if err != nil {
		return renderJQError(writer, jqSurfaceError(command, err.Error()))
	}
	iter := query.RunWithContext(ctx, projected)
	for {
		item, ok := iter.Next()
		if !ok {
			return nil
		}
		if err, ok := item.(error); ok {
			return renderJQError(writer, jqSurfaceError(command, err.Error()))
		}
		switch v := item.(type) {
		case string:
			if _, err := fmt.Fprintln(writer, v); err != nil {
				return err
			}
		default:
			data, err := json.Marshal(v)
			if err != nil {
				return renderJQError(writer, jqSurfaceError(command, err.Error()))
			}
			if _, err := fmt.Fprintln(writer, string(data)); err != nil {
				return err
			}
		}
	}
}

func jqSurfaceError(command, meaning string) model.SurfaceError {
	return model.SurfaceError{
		Command:  command,
		Object:   "jq_expression",
		Source:   "embedded gojq interpreter",
		State:    model.StateFailed,
		Meaning:  meaning,
		Recovery: "inspect the command shape, correct the expression, or use --json",
	}
}

func renderJQError(writer io.Writer, surface model.SurfaceError) error {
	if err := RenderError(writer, surface); err != nil {
		return err
	}
	return SurfaceError{Value: surface}
}
