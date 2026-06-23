package output

import (
	"errors"
	"fmt"
	"io"

	"github.com/yumiaura/seekmoon/internal/model"
)

type SurfaceError struct {
	Value model.SurfaceError
}

func (e SurfaceError) Error() string {
	return e.Value.Meaning
}

func ErrorFrom(command string, err error) model.SurfaceError {
	var surface SurfaceError
	if errors.As(err, &surface) {
		return surface.Value
	}
	return model.SurfaceError{
		Command: command,
		Object:  "command",
		Source:  "service",
		State:   model.StateFailed,
		Meaning: err.Error(),
	}
}

func RenderError(writer io.Writer, surface model.SurfaceError) error {
	rows := [][2]string{
		{"command", surface.Command},
		{"object", surface.Object},
		{"source", surface.Source},
		{"state", string(surface.State)},
		{"meaning", surface.Meaning},
	}
	if surface.Recovery != "" {
		rows = append(rows, [2]string{"recover", surface.Recovery})
	}
	if surface.LogPath != "" {
		rows = append(rows, [2]string{"log", surface.LogPath})
	}
	for _, row := range rows {
		if _, err := fmt.Fprintf(writer, "%-8s %s\n", row[0], row[1]); err != nil {
			return err
		}
	}
	return nil
}
