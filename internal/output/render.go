// Package output projects canonical objects into terminal, JSON, jq, shape, schema, and error surfaces.
package output

import (
	"context"
	"fmt"
	"io"

	"github.com/MichengLiang/seekmoon/internal/model"
)

// Renderer renders command output requests.
type Renderer interface {
	Render(ctx context.Context, request Request) error
}

// Request carries all data needed to render a command response.
type Request struct {
	Command      string
	Mode         model.OutputMode
	JQExpression string
	Schema       string
	Writer       io.Writer
	Value        any
	Err          error
}

// RendererFunc adapts a function to the Renderer interface.
type RendererFunc func(context.Context, Request) error

// Render calls the wrapped renderer function.
func (f RendererFunc) Render(ctx context.Context, request Request) error {
	return f(ctx, request)
}

// DefaultRenderer dispatches to the built-in output modes.
type DefaultRenderer struct {
	Writer io.Writer
}

// Render writes either a command result or a structured command error.
func (r DefaultRenderer) Render(ctx context.Context, request Request) error {
	if request.Writer == nil {
		request.Writer = r.Writer
	}
	if request.Writer == nil {
		request.Writer = io.Discard
	}
	if request.Mode == "" {
		request.Mode = model.OutputPretty
	}
	if request.Err != nil {
		surface := ErrorFrom(request.Command, request.Err)
		if err := RenderError(request.Writer, surface); err != nil {
			return err
		}
		return SurfaceError{Value: surface}
	}
	switch request.Mode {
	case model.OutputPretty:
		return RenderPretty(request.Writer, request.Value)
	case model.OutputJSON:
		return RenderJSON(request.Writer, request.Schema, request.Value)
	case model.OutputJQ:
		return RenderJQ(ctx, request.Writer, request.Schema, request.Value, request.JQExpression, request.Command)
	case model.OutputShape:
		return RenderShape(request.Writer, request.Schema)
	case model.OutputSchema:
		return RenderSchema(request.Writer, request.Schema)
	default:
		surface := model.SurfaceError{
			Command:  request.Command,
			Object:   "output_mode",
			Source:   "output renderer",
			State:    model.StateFailed,
			Meaning:  fmt.Sprintf("unsupported output mode %q", request.Mode),
			Recovery: "choose one of pretty, json, jq, shape, or schema",
		}
		if err := RenderError(request.Writer, surface); err != nil {
			return err
		}
		return SurfaceError{Value: surface}
	}
}
