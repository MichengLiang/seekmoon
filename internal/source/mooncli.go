package source

import (
	"context"

	"github.com/MichengLiang/seekmoon/internal/model"
	"github.com/MichengLiang/seekmoon/internal/platform"
	"github.com/MichengLiang/seekmoon/internal/store"
)

// MoonCLI runs Moon CLI commands through a platform runner.
type MoonCLI struct {
	Runner platform.Runner
	Paths  store.Paths
}

// Run executes a Moon-related command and records command evidence.
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

// Version runs moon --version.
func (m MoonCLI) Version(ctx context.Context, cwd string) model.CommandResult {
	return m.Run(ctx, cwd, "moon-version", "moon", "--version")
}

// Fetch runs moon fetch for a module.
func (m MoonCLI) Fetch(ctx context.Context, cwd, module string) model.CommandResult {
	return m.Run(ctx, cwd, "moon-fetch", "moon", "fetch", module)
}

// Update runs moon update.
func (m MoonCLI) Update(ctx context.Context, cwd string) model.CommandResult {
	return m.Run(ctx, cwd, "moon-update", "moon", "update")
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
