package app

import (
	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
)

// CommandResultFromRun converts a platform run result into model command evidence.
func CommandResultFromRun(result platform.RunResult, err error) model.CommandResult {
	status := model.StatePresent
	if err != nil {
		status = model.StateFailed
	}
	return model.CommandResult{
		Command:  result.Command,
		CWD:      result.CWD,
		ExitCode: result.ExitCode,
		Status:   status,
		LogPath:  result.LogPath,
	}
}
