package service

import (
	"context"

	"github.com/MichengLiang/seekmoon/internal/model"
	"github.com/MichengLiang/seekmoon/internal/source"
)

// RawFlow projects upstream source payloads without normalization.
type RawFlow struct {
	Mooncakes source.MooncakesClient
	Assets    source.AssetClient
	Skills    source.SkillsClient
}

// Raw returns the selected upstream payload with source status metadata.
func (s RawFlow) Raw(ctx context.Context, input RawInput) (any, error) {
	switch input.Source {
	case "modules":
		result := s.Mooncakes.FetchRawModules(ctx)
		return rawEnvelope(result.Source, result.URL, result.Path, result.FetchedAt, result.Status, result.RawRef, rawPayload(result.Value), result.Error), nil
	case "manifest":
		module := firstArg(input.Args)
		result := s.Mooncakes.FetchRawManifest(ctx, module)
		return rawEnvelope(result.Source, result.URL, result.Path, result.FetchedAt, result.Status, result.RawRef, rawPayload(result.Value), result.Error), nil
	case "module-index":
		module, version := splitModuleVersion(firstArg(input.Args))
		result := s.Assets.FetchRawModuleIndex(ctx, module, version)
		return rawEnvelope(result.Source, result.URL, result.Path, result.FetchedAt, result.Status, result.RawRef, rawPayload(result.Value), result.Error), nil
	case "package-data":
		module, version := splitModuleVersion(firstArg(input.Args))
		pkg := firstArg(input.Args[1:])
		result := s.Assets.FetchRawPackageData(ctx, module, version, pkg)
		return rawEnvelope(result.Source, result.URL, result.Path, result.FetchedAt, result.Status, result.RawRef, rawPayload(result.Value), result.Error), nil
	case "skills":
		result := s.Skills.FetchRawSkills(ctx)
		return rawEnvelope(result.Source, result.URL, result.Path, result.FetchedAt, result.Status, result.RawRef, rawPayload(result.Value), result.Error), nil
	default:
		return nil, surfaceFailure(input.Source, "command input", model.StateFailed, "unknown raw source selector", "choose modules, manifest, module-index, package-data, or skills")
	}
}

func rawPayload(value *any) any {
	if value == nil {
		return nil
	}
	return *value
}

func rawEnvelope(sourceName, url, path, fetchedAt string, status model.State, rawRef string, payload any, err string) model.RawEnvelope {
	return model.RawEnvelope{
		Schema:    model.SchemaRawPayloadV1,
		Source:    sourceName,
		URL:       url,
		Path:      path,
		FetchedAt: fetchedAt,
		Status:    status,
		RawRef:    rawRef,
		Payload:   payload,
		Error:     err,
	}
}

func firstArg(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

func splitModuleVersion(value string) (string, string) {
	module := value
	version := "latest"
	for i := len(value) - 1; i >= 0; i-- {
		if value[i] == '@' {
			module = value[:i]
			version = value[i+1:]
			break
		}
	}
	return module, version
}
