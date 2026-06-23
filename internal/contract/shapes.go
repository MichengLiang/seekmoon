// Package contract owns output contract definitions independent of renderers.
package contract

import "github.com/yumiaura/seekmoon/internal/model"

type Field struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Required bool    `json:"required"`
	Fields   []Field `json:"fields,omitempty"`
}

type Shape struct {
	Schema string  `json:"schema"`
	Fields []Field `json:"fields"`
}

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
