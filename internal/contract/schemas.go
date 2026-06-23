package contract

import "github.com/yumiaura/seekmoon/internal/model"

// Schema is a JSON Schema object rendered by the schema output mode.
type Schema map[string]any

// SearchResultsSchema returns the contract for search result JSON output.
func SearchResultsSchema() Schema {
	return Schema{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"$id":     model.SchemaSearchResultsV1,
		"type":    "object",
		"required": []string{
			"schema",
			"snapshot",
			"query",
			"results",
		},
		"properties": map[string]any{
			"schema": map[string]any{"const": model.SchemaSearchResultsV1},
			"snapshot": map[string]any{
				"type":     "object",
				"required": []string{"id", "sources"},
				"properties": map[string]any{
					"id":      map[string]any{"type": "string"},
					"sources": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
				},
			},
			"query": map[string]any{
				"type":     "object",
				"required": []string{"text", "kind"},
				"properties": map[string]any{
					"text":   map[string]any{"type": "string"},
					"kind":   map[string]any{"enum": []string{"library", "skill"}},
					"target": map[string]any{"type": []string{"string", "null"}},
				},
			},
			"results": map[string]any{
				"type":  "array",
				"items": map[string]any{"$ref": "#/$defs/searchResult"},
			},
		},
		"$defs": map[string]any{
			"state": map[string]any{
				"enum": []string{
					string(model.StatePresent),
					string(model.StateMissing),
					string(model.StateUnknown),
					string(model.StateFailed),
					string(model.StateUnavailable),
					string(model.StateDerived),
				},
			},
			"evidence": map[string]any{
				"type":     "object",
				"required": []string{"status", "value", "source"},
				"properties": map[string]any{
					"status": map[string]any{"$ref": "#/$defs/state"},
					"value":  map[string]any{},
					"source": map[string]any{"type": []string{"string", "null"}},
					"error":  map[string]any{"type": "string"},
				},
			},
			"searchResult": map[string]any{
				"type":     "object",
				"required": []string{"rank", "module", "version", "description", "license", "repository", "match", "snapshot_id"},
				"properties": map[string]any{
					"rank":         map[string]any{"type": "integer"},
					"module":       map[string]any{"type": "string"},
					"version":      map[string]any{"type": "string"},
					"description":  map[string]any{"$ref": "#/$defs/evidence"},
					"keywords":     map[string]any{"$ref": "#/$defs/evidence"},
					"license":      map[string]any{"$ref": "#/$defs/evidence"},
					"repository":   map[string]any{"$ref": "#/$defs/evidence"},
					"downloads":    map[string]any{"$ref": "#/$defs/evidence"},
					"build_status": map[string]any{"$ref": "#/$defs/evidence"},
					"target":       map[string]any{"$ref": "#/$defs/evidence"},
					"match":        map[string]any{"type": "object"},
					"snapshot_id":  map[string]any{"type": "string"},
				},
			},
		},
	}
}

// AdoptionRecordSchema returns the contract for adoption record JSON output.
func AdoptionRecordSchema() Schema {
	return Schema{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"$id":     model.SchemaAdoptionRecordV1,
		"type":    "object",
		"required": []string{
			"record_id",
			"created_at",
			"project",
			"snapshot_id",
			"candidate",
			"version",
			"conclusion",
			"evidence_refs",
		},
		"properties": map[string]any{
			"record_id":   map[string]any{"type": "string"},
			"created_at":  map[string]any{"type": "string", "format": "date-time"},
			"project":     map[string]any{"type": "object"},
			"snapshot_id": map[string]any{"type": "string"},
			"candidate":   map[string]any{"type": "object"},
			"version":     map[string]any{"type": "string"},
			"conclusion": map[string]any{"enum": []string{
				string(model.ConclusionAdopt),
				string(model.ConclusionAdoptWithAdapter),
				string(model.ConclusionContinueVerification),
				string(model.ConclusionContributeUpstream),
				string(model.ConclusionFork),
				string(model.ConclusionBuildOwn),
				string(model.ConclusionRejectForNow),
			}},
			"evidence_refs": map[string]any{"type": "array"},
			"not_confirmed": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"note":          map[string]any{"type": "string"},
		},
	}
}

// EnvironmentStatusSchema returns the contract for doctor environment output.
func EnvironmentStatusSchema() Schema {
	return genericObjectSchema(model.SchemaEnvironmentStatusV1)
}

// SnapshotSchema returns the contract for sync snapshot output.
func SnapshotSchema() Schema {
	return genericObjectSchema(model.SchemaSnapshotV1)
}

// ManifestProfileSchema returns the contract for manifest profile output.
func ManifestProfileSchema() Schema {
	return genericObjectSchema(model.SchemaManifestProfileV1)
}

// PackageDataSchema returns the contract for package API data output.
func PackageDataSchema() Schema {
	return genericObjectSchema(model.SchemaPackageDataV1)
}

// SkillEntrySchema returns the contract for skill entry output.
func SkillEntrySchema() Schema {
	return genericObjectSchema(model.SchemaSkillEntryV1)
}

// SourceResolutionSchema returns the contract for source resolution output.
func SourceResolutionSchema() Schema {
	return genericObjectSchema(model.SchemaSourceResolutionV1)
}

// ComparisonSchema returns the contract for comparison output.
func ComparisonSchema() Schema {
	return genericObjectSchema(model.SchemaComparisonV1)
}

// ProbeResultSchema returns the contract for probe result output.
func ProbeResultSchema() Schema {
	return genericObjectSchema(model.SchemaProbeResultV1)
}

// ReportSchema returns the contract for report output.
func ReportSchema() Schema {
	return genericObjectSchema(model.SchemaReportV1)
}

// RawPayloadSchema returns the contract for raw source payload output.
func RawPayloadSchema() Schema {
	return genericObjectSchema(model.SchemaRawPayloadV1)
}

func genericObjectSchema(id string) Schema {
	return Schema{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"$id":     id,
		"type":    "object",
		"required": []string{
			"schema",
		},
		"properties": map[string]any{
			"schema": map[string]any{"const": id},
		},
	}
}
