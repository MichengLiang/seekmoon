package service

import (
	"context"
	"path/filepath"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

// ProbeFlow runs isolated Moon CLI verification.
type ProbeFlow struct {
	MoonCLI  source.MoonCLI
	Sessions store.SessionStore
	Probes   store.ProbeStore
	Paths    store.Paths
}

// Probe executes or records probe commands for a candidate.
func (s ProbeFlow) Probe(ctx context.Context, input ProbeInput) (model.ProbeResult, error) {
	candidate, err := candidateFromRequest(ctx, s.Sessions, input.Candidate)
	if err != nil {
		return model.ProbeResult{}, err
	}
	version := firstNonEmpty(candidate.Version, "latest")
	targetSuffix := input.Target
	if targetSuffix == "" {
		targetSuffix = "default"
	}
	probeID := candidate.Module + "@" + version + "-" + targetSuffix
	probePath := s.Probes.ProbeDir(probeID)
	results := model.ProbeResult{
		Module:    candidate.Module,
		Version:   version,
		Target:    input.Target,
		ProbePath: probePath,
		Logs:      map[string]string{},
	}
	if s.MoonCLI.Runner == nil {
		results.MoonNew = unknownCommand("moon", "new")
		results.MoonAdd = unknownCommand("moon", "add", candidate.Module)
		results.MoonCheck = unknownCommand("moon", "check")
		results.MoonTest = unknownCommand("moon", "test")
		results.MoonCheckTarget = unknownCommand("moon", "check", "--target", input.Target)
		results.MoonBuildTarget = unknownCommand("moon", "build", "--target", input.Target)
		results.Result = "incomplete"
		return results, nil
	}
	results.MoonNew = s.MoonCLI.Run(ctx, probePath, "moon-new", "moon", "new", "seekmoon-probe")
	results.MoonAdd = s.MoonCLI.Run(ctx, probePath, "moon-add", "moon", "add", candidate.Module)
	results.MoonCheck = s.MoonCLI.Run(ctx, probePath, "moon-check", "moon", "check")
	results.MoonTest = s.MoonCLI.Run(ctx, probePath, "moon-test", "moon", "test")
	if input.Target != "" {
		results.MoonCheckTarget = s.MoonCLI.Run(ctx, probePath, "moon-check-target", "moon", "check", "--target", input.Target)
		results.MoonBuildTarget = s.MoonCLI.Run(ctx, probePath, "moon-build-target", "moon", "build", "--target", input.Target)
	} else {
		results.MoonCheckTarget = unknownCommand("moon", "check", "--target")
		results.MoonBuildTarget = unknownCommand("moon", "build", "--target")
	}
	results.Logs = commandLogs(results)
	results.Result = probeResultStatus(results)
	return results, nil
}

func commandLogs(result model.ProbeResult) map[string]string {
	logs := map[string]string{}
	commands := map[string]model.CommandResult{
		"moon_new":          result.MoonNew,
		"moon_add":          result.MoonAdd,
		"moon_check":        result.MoonCheck,
		"moon_test":         result.MoonTest,
		"moon_check_target": result.MoonCheckTarget,
		"moon_build_target": result.MoonBuildTarget,
	}
	for key, command := range commands {
		if command.LogPath != "" {
			logs[key] = filepath.Clean(command.LogPath)
		}
	}
	return logs
}

func probeResultStatus(result model.ProbeResult) string {
	for _, command := range []model.CommandResult{result.MoonNew, result.MoonAdd, result.MoonCheck, result.MoonTest, result.MoonCheckTarget, result.MoonBuildTarget} {
		if command.Status == model.StateFailed {
			return "failed"
		}
		if command.Status == model.StateUnknown {
			return "incomplete"
		}
	}
	return "verifiable"
}
