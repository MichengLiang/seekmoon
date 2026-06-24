package source

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MichengLiang/seekmoon/internal/model"
	"github.com/MichengLiang/seekmoon/internal/platform"
	"github.com/pelletier/go-toml/v2"
)

// ProjectReader reads local Moon project configuration.
type ProjectReader struct {
	FS    platform.FS
	Clock platform.Clock
}

func (r ProjectReader) Read(ctx context.Context, root string) model.SourceResult[model.ProjectContext] {
	fs := r.FS
	if fs == nil {
		fs = platform.OSFS{}
	}
	project := model.ProjectContext{
		Identity:             model.ProjectIdentity{Root: root},
		ModuleConfig:         model.Unavailable[map[string]any](string(model.SourceProjectContext)),
		PackageConfig:        model.Unavailable[map[string]any](string(model.SourceProjectContext)),
		DeclaredTarget:       model.Unknown[string](),
		ExistingDependencies: model.Unknown[map[string]any](),
	}
	moduleConfig := readFirstConfig(ctx, fs, root, "moon.mod.json", "moon.mod")
	switch moduleConfig.Status {
	case model.StatePresent:
		project.ModuleConfig = moduleConfig.Evidence
		if module, ok := moduleConfig.Value["name"].(string); ok {
			project.Identity.Module = module
		}
		if deps, ok := moduleConfig.Value["deps"].(map[string]any); ok {
			project.ExistingDependencies = model.Present(deps, string(model.SourceProjectContext))
		}
	case model.StateFailed:
		project.ModuleConfig = moduleConfig.Evidence
	}
	packageConfig := readFirstConfig(ctx, fs, root, "moon.pkg.json", "moon.pkg")
	switch packageConfig.Status {
	case model.StatePresent:
		project.PackageConfig = packageConfig.Evidence
		if target, ok := packageConfig.Value["target"].(string); ok {
			project.DeclaredTarget = model.Present(target, string(model.SourceProjectContext))
		}
	case model.StateFailed:
		project.PackageConfig = packageConfig.Evidence
	}
	status := model.StatePresent
	parseState := model.StatePresent
	var readErr string
	if moduleConfig.Status == model.StateFailed || packageConfig.Status == model.StateFailed {
		status = model.StateFailed
		parseState = model.StateFailed
		readErr = strings.Join(failedConfigErrors(moduleConfig, packageConfig), "; ")
	}
	cleanRoot := filepath.Clean(root)
	return model.SourceResult[model.ProjectContext]{
		Source:     string(model.SourceProjectContext),
		Path:       cleanRoot,
		FetchedAt:  sourceNow(r.Clock),
		Status:     status,
		ParseState: parseState,
		RawRef:     fmt.Sprintf("project:%s", cleanRoot),
		Error:      readErr,
		Value:      &project,
	}
}

type configRead struct {
	Status   model.State
	Path     string
	Value    map[string]any
	Evidence model.EvidenceObject
	Error    string
}

func readFirstConfig(ctx context.Context, fs platform.FS, root string, names ...string) configRead {
	for _, name := range names {
		path := filepath.Join(root, name)
		data, err := fs.ReadFile(ctx, path)
		if err != nil {
			continue
		}
		value := map[string]any{}
		if strings.HasSuffix(name, ".json") {
			if err := json.Unmarshal(data, &value); err != nil {
				return configRead{Status: model.StateFailed, Path: path, Evidence: model.Failed[map[string]any](string(model.SourceProjectContext), err.Error()), Error: fmt.Sprintf("%s: %s", path, err.Error())}
			}
		} else if err := toml.Unmarshal(data, &value); err != nil {
			return configRead{Status: model.StateFailed, Path: path, Evidence: model.Failed[map[string]any](string(model.SourceProjectContext), err.Error()), Error: fmt.Sprintf("%s: %s", path, err.Error())}
		}
		return configRead{Status: model.StatePresent, Path: path, Value: value, Evidence: model.Present(value, string(model.SourceProjectContext))}
	}
	return configRead{Status: model.StateUnavailable, Evidence: model.Unavailable[map[string]any](string(model.SourceProjectContext))}
}

func failedConfigErrors(configs ...configRead) []string {
	var out []string
	for _, config := range configs {
		if config.Status == model.StateFailed && config.Error != "" {
			out = append(out, config.Error)
		}
	}
	return out
}
