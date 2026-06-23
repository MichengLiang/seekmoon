package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/yumiaura/seekmoon/internal/app"
)

func TestRootCommandThinStartupSurface(t *testing.T) {
	rt, err := app.NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime: %v", err)
	}
	var out bytes.Buffer
	if err := Execute(context.Background(), rt, Options{Out: &out}); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(out.String(), "pending Batch C") {
		t.Fatalf("root output = %q", out.String())
	}
}

func TestPlaceholderCommandDoesNotImplementBatchBehavior(t *testing.T) {
	rt, err := app.NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime: %v", err)
	}
	cmd := NewRoot(rt, Options{})
	cmd.SetArgs([]string{"search"})
	err = cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "outside Batch A") {
		t.Fatalf("placeholder command error = %v", err)
	}
}
