// Package contract owns output contract definitions independent of renderers.
package contract

import "github.com/yumiaura/seekmoon/internal/model"

// Field describes one projected output field for shape rendering.
type Field struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Required bool    `json:"required"`
	Fields   []Field `json:"fields,omitempty"`
}

// Shape describes the stable field outline for an output schema.
type Shape struct {
	Schema string  `json:"schema"`
	Fields []Field `json:"fields"`
}

// SearchResultsShape returns the field outline for search result output.
func SearchResultsShape() Shape {
	return Shape{
		Schema: model.SchemaSearchResultsV1,
		Fields: []Field{
			{Name: "schema", Type: "string", Required: true},
			{Name: "snapshot", Type: "object", Required: true, Fields: []Field{
				{Name: "id", Type: "string", Required: true},
				{Name: "sources", Type: "string[]", Required: true},
			}},
			{Name: "query", Type: "object", Required: true, Fields: []Field{
				{Name: "text", Type: "string", Required: true},
				{Name: "kind", Type: "library|skill", Required: true},
				{Name: "target", Type: "string|null", Required: false},
			}},
			{Name: "results", Type: "array", Required: true, Fields: []Field{
				{Name: "rank", Type: "int", Required: true},
				{Name: "module", Type: "string", Required: true},
				{Name: "version", Type: "string", Required: true},
				{Name: "description", Type: "evidence<string>", Required: true},
				{Name: "license", Type: "evidence<string>", Required: true},
				{Name: "repository", Type: "evidence<string>", Required: true},
				{Name: "target", Type: "evidence<object>", Required: false},
				{Name: "match", Type: "object", Required: true},
				{Name: "snapshot_id", Type: "string", Required: true},
			}},
		},
	}
}

// AdoptionRecordShape returns the field outline for adoption record output.
func AdoptionRecordShape() Shape {
	return Shape{
		Schema: model.SchemaAdoptionRecordV1,
		Fields: []Field{
			{Name: "record_id", Type: "string", Required: true},
			{Name: "created_at", Type: "datetime", Required: true},
			{Name: "project", Type: "object", Required: true},
			{Name: "snapshot_id", Type: "string", Required: true},
			{Name: "candidate", Type: "object", Required: true},
			{Name: "version", Type: "string", Required: true},
			{Name: "conclusion", Type: "adoption-conclusion", Required: true},
			{Name: "evidence_refs", Type: "array", Required: true},
			{Name: "not_confirmed", Type: "array", Required: false},
			{Name: "note", Type: "string", Required: false},
		},
	}
}

// EnvironmentStatusShape returns the field outline for doctor environment output.
func EnvironmentStatusShape() Shape {
	return genericShape(model.SchemaEnvironmentStatusV1, "toolchain", "paths", "network", "project")
}

// SnapshotShape returns the field outline for sync snapshot output.
func SnapshotShape() Shape {
	return genericShape(model.SchemaSnapshotV1, "id", "created_at", "sources", "statistics", "raw")
}

// ManifestProfileShape returns the field outline for manifest profile output.
func ManifestProfileShape() Shape {
	return genericShape(model.SchemaManifestProfileV1, "module", "version", "metadata", "downloads", "build_status", "docs_url")
}

// PackageDataShape returns the field outline for package API data output.
func PackageDataShape() Shape {
	return genericShape(model.SchemaPackageDataV1, "name", "traits", "errors", "types", "typealias", "values", "misc")
}

// SkillEntryShape returns the field outline for skill entry output.
func SkillEntryShape() Shape {
	return genericShape(model.SchemaSkillEntryV1, "module", "version", "package", "name", "wasm_url", "checksum_url", "repository")
}

// SourceResolutionShape returns the field outline for source resolution output.
func SourceResolutionShape() Shape {
	return genericShape(model.SchemaSourceResolutionV1, "module", "version", "moon_fetch", "source_zip", "local_cache", "selected_source")
}

// ComparisonShape returns the field outline for comparison output.
func ComparisonShape() Shape {
	return genericShape(model.SchemaComparisonV1, "candidates", "evidence", "differences", "recommendation")
}

// ProbeResultShape returns the field outline for probe result output.
func ProbeResultShape() Shape {
	return genericShape(model.SchemaProbeResultV1, "module", "version", "target", "probe_path", "moon_add", "moon_check", "result")
}

// ReportShape returns the field outline for report output.
func ReportShape() Shape {
	return genericShape(model.SchemaReportV1, "goal", "environment", "data_sources", "candidates", "decision")
}

// RawPayloadShape returns the field outline for raw source payload output.
func RawPayloadShape() Shape {
	return genericShape(model.SchemaRawPayloadV1, "source", "status", "payload", "metadata")
}

func genericShape(schema string, names ...string) Shape {
	fields := make([]Field, 0, 1+len(names))
	fields = append(fields, Field{Name: "schema", Type: "string", Required: true})
	for _, name := range names {
		fields = append(fields, Field{Name: name, Type: "object|string", Required: false})
	}
	return Shape{Schema: schema, Fields: fields}
}
