package app

import (
	"errors"
	"testing"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
)

func TestCommandResultFromRunMapsPlatformResultToModelState(t *testing.T) {
	result := CommandResultFromRun(platform.RunResult{
		Command:  []string{"moon", "check"},
		CWD:      "/tmp/project",
		ExitCode: 0,
		LogPath:  ".seekmoon/logs/moon-check.log",
	}, nil)
	if result.Status != model.StatePresent {
		t.Fatalf("status = %q", result.Status)
	}

	failed := CommandResultFromRun(platform.RunResult{
		Command:  []string{"moon", "check"},
		CWD:      "/tmp/project",
		ExitCode: 1,
		LogPath:  ".seekmoon/logs/moon-check.log",
	}, errors.New("exit status 1"))
	if failed.Status != model.StateFailed {
		t.Fatalf("failed status = %q", failed.Status)
	}
}
