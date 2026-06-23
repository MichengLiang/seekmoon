package source

import (
	"context"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/platform"
	"github.com/yumiaura/seekmoon/internal/store"
)

type MoonCLI struct {
	Runner platform.Runner
	Paths  store.Paths
}

func (m MoonCLI) Run(ctx context.Context, cwd, logName string, command ...string) model.CommandResult {
	runner := m.Runner
	if runner == nil {
		runner = platform.ExecRunner{FS: platform.OSFS{}}
	}
	probe := store.ProbeStore{Paths: m.Paths}
	logPath := probe.LogPath("moon-cli", logName)
	result, err := runner.Run(ctx, platform.RunRequest{Command: command, CWD: cwd, LogPath: logPath})
	return commandResultFromRun(result, err)
}

func (m MoonCLI) Version(ctx context.Context, cwd string) model.CommandResult {
	return m.Run(ctx, cwd, "moon-version", "moon", "--version")
}

func commandResultFromRun(result platform.RunResult, err error) model.CommandResult {
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
