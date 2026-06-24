package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/MichengLiang/seekmoon/internal/model"
	"github.com/MichengLiang/seekmoon/internal/store"
)

const defaultSessionID = "default"

func surfaceFailure(object, source string, state model.State, meaning, recovery string) error {
	return model.SurfaceFailure{Value: model.SurfaceError{
		Object:   object,
		Source:   source,
		State:    state,
		Meaning:  meaning,
		Recovery: recovery,
	}}
}

func sourceFailure(object string, result model.SourceResult[any], recovery string) error {
	return surfaceFailure(object, result.Source, result.Status, firstNonEmpty(result.Error, "source did not return a usable value"), recovery)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func evidenceStringValue(value model.EvidenceString) string {
	if value.Value != nil {
		return *value.Value
	}
	if value.Status != "" {
		return string(value.Status)
	}
	return ""
}

func candidateFromRequest(ctx context.Context, sessions store.SessionStore, request model.CandidateRequest) (model.CandidateRef, error) {
	if request.Number > 0 {
		candidates, err := sessions.ReadCandidates(ctx, defaultSessionID)
		if err != nil {
			return model.CandidateRef{}, surfaceFailure(
				request.Raw,
				"session_store",
				model.StateFailed,
				fmt.Sprintf("candidate number %d is not available in the current session", request.Number),
				"run search first or pass a module coordinate",
			)
		}
		candidate, ok := candidates.Candidates[request.Number]
		if !ok {
			return model.CandidateRef{}, surfaceFailure(
				request.Raw,
				"session_store",
				model.StateFailed,
				fmt.Sprintf("candidate number %d is not present in the current session", request.Number),
				"run search again or pass a module coordinate",
			)
		}
		return candidate, nil
	}
	if request.Module == "" {
		return model.CandidateRef{}, surfaceFailure(request.Raw, "command input", model.StateFailed, "module coordinate is required", "pass owner/module or a search candidate number")
	}
	return model.CandidateRef{Kind: "library", Module: request.Module, Version: request.Version}, nil
}

func snapshotModules(ctx context.Context, snapshots store.SnapshotStore) (model.Snapshot, []model.ModuleSummary, error) {
	snapshot, err := snapshots.Latest(ctx)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return model.Snapshot{}, nil, os.ErrNotExist
		}
		return model.Snapshot{}, nil, err
	}
	return snapshot, modulesFromSnapshot(snapshot), nil
}

func modulesFromSnapshot(snapshot model.Snapshot) []model.ModuleSummary {
	if snapshot.Raw == nil {
		return nil
	}
	raw, ok := snapshot.Raw["modules"]
	if !ok {
		return nil
	}
	data, ok := raw.([]any)
	if !ok {
		return nil
	}
	modules := make([]model.ModuleSummary, 0, len(data))
	for _, item := range data {
		object, ok := item.(map[string]any)
		if !ok {
			continue
		}
		module, _ := object["module"].(string)
		version, _ := object["version"].(string)
		if module == "" {
			continue
		}
		modules = append(modules, model.ModuleSummary{
			Module:      module,
			Version:     version,
			Description: evidenceFromAnyString(object["description"]),
			Keywords:    evidenceFromAnyStrings(object["keywords"]),
			Repository:  evidenceFromAnyString(object["repository"]),
			License:     evidenceFromAnyString(object["license"]),
		})
	}
	return modules
}

func evidenceFromAnyString(raw any) model.EvidenceString {
	object, ok := raw.(map[string]any)
	if !ok {
		return model.Unknown[string]()
	}
	status, _ := object["status"].(string)
	source, _ := object["source"].(string)
	value, _ := object["value"].(string)
	switch model.State(status) {
	case model.StatePresent, model.StateDerived:
		return model.Present(value, source)
	case model.StateMissing:
		return model.Missing[string](source)
	case model.StateFailed:
		err, _ := object["error"].(string)
		return model.Failed[string](source, err)
	case model.StateUnavailable:
		return model.Unavailable[string](source)
	default:
		return model.Unknown[string]()
	}
}

func evidenceFromAnyStrings(raw any) model.EvidenceStringArray {
	object, ok := raw.(map[string]any)
	if !ok {
		return model.Unknown[[]string]()
	}
	status, _ := object["status"].(string)
	source, _ := object["source"].(string)
	var values []string
	if rawValues, ok := object["value"].([]any); ok {
		for _, rawValue := range rawValues {
			if value, ok := rawValue.(string); ok {
				values = append(values, value)
			}
		}
	}
	switch model.State(status) {
	case model.StatePresent, model.StateDerived:
		return model.Present(values, source)
	case model.StateMissing:
		return model.Missing[[]string](source)
	case model.StateFailed:
		err, _ := object["error"].(string)
		return model.Failed[[]string](source, err)
	case model.StateUnavailable:
		return model.Unavailable[[]string](source)
	default:
		return model.Unknown[[]string]()
	}
}

func packagePaths(tree model.ModuleIndexTree) []string {
	var paths []string
	var walk func(model.ModuleIndexTree)
	walk = func(node model.ModuleIndexTree) {
		if node.Package != nil && node.Package.Path != "" {
			paths = append(paths, node.Package.Path)
		}
		for _, child := range node.Childs {
			walk(child)
		}
	}
	walk(tree)
	sort.Strings(paths)
	return paths
}

func hasPackagePath(tree model.ModuleIndexTree, path string) bool {
	for _, known := range packagePaths(tree) {
		if known == path {
			return true
		}
	}
	return false
}
