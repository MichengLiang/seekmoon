package platform

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
)

// RunRequest describes one command execution request.
type RunRequest struct {
	Command []string
	CWD     string
	Env     []string
	LogPath string
}

// Runner executes host commands for Moon CLI workflows.
type Runner interface {
	Run(ctx context.Context, request RunRequest) (RunResult, error)
}

// RunResult records command execution metadata.
type RunResult struct {
	Command  []string
	CWD      string
	ExitCode int
	LogPath  string
}

// ExecRunner executes commands through os/exec.
type ExecRunner struct {
	FS FS
}

// Run executes the request command and optionally writes combined logs.
func (r ExecRunner) Run(ctx context.Context, request RunRequest) (RunResult, error) {
	if len(request.Command) == 0 {
		return RunResult{Command: request.Command, CWD: request.CWD, LogPath: request.LogPath}, exec.ErrNotFound
	}
	// #nosec G204 -- command argv comes from SeekMoon service command construction,
	// not a shell string; exec.CommandContext preserves argument boundaries.
	cmd := exec.CommandContext(ctx, request.Command[0], request.Command[1:]...)
	cmd.Dir = request.CWD
	cmd.Env = request.Env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := RunResult{
		Command: request.Command,
		CWD:     request.CWD,
		LogPath: request.LogPath,
	}
	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}
	if request.LogPath != "" {
		fs := r.FS
		if fs == nil {
			fs = OSFS{}
		}
		if mkErr := fs.MkdirAll(ctx, filepath.Dir(request.LogPath), 0o755); mkErr != nil && err == nil {
			err = mkErr
		}
		payload := append(stdout.Bytes(), stderr.Bytes()...)
		if writeErr := fs.WriteFile(ctx, request.LogPath, payload, 0o644); writeErr != nil && err == nil {
			err = writeErr
		}
	}
	return result, err
}
