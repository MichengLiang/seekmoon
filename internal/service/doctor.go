package service

import (
	"context"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

// DoctorFlow reports local environment and project context evidence.
type DoctorFlow struct {
	Project source.ProjectReader
	MoonCLI source.MoonCLI
	Paths   store.Paths
}

// Doctor returns environment status without mutating project records.
func (s DoctorFlow) Doctor(ctx context.Context, input DoctorInput) (any, error) {
	_ = input
	root := s.Paths.ProjectRoot
	toolchain := map[string]model.CommandResult{}
	commands := map[string]model.CommandResult{}
	if s.MoonCLI.Runner != nil {
		toolchain["moon"] = s.MoonCLI.Version(ctx, root)
		for _, command := range []string{"moon", "moonc", "moonrun", "mooncake"} {
			commands[command] = s.MoonCLI.Run(ctx, root, "doctor-"+command, command, "--version")
		}
	} else {
		toolchain["moon"] = unknownCommand("moon", "--version")
		for _, command := range []string{"moon", "moonc", "moonrun", "mooncake"} {
			commands[command] = unknownCommand(command, "--version")
		}
	}
	return model.EnvironmentStatus{
		Schema:    model.SchemaEnvironmentStatusV1,
		Toolchain: toolchain,
		Commands:  commands,
		Paths: map[string]model.EvidenceString{
			"registry_index": model.Present(s.Paths.MoonIndex, string(model.SourceLocalIndex)),
			"registry_cache": model.Present(s.Paths.Sources, string(model.SourceLocalCache)),
			"assets_cache":   model.Present(s.Paths.Assets, string(model.SourceResourceAsset)),
			"core_source":    model.Unknown[string](),
		},
		Network: map[string]model.EvidenceString{
			"mooncakes_api": model.Unknown[string](),
			"assets":        model.Unknown[string](),
		},
		Project: s.Project.Read(ctx, root),
	}, nil
}
