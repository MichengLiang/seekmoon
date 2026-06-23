package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/yumiaura/seekmoon/internal/contract"
	"github.com/yumiaura/seekmoon/internal/model"
)

// RenderShape writes the compact field shape for a schema identifier.
func RenderShape(writer io.Writer, schema string) error {
	shape, err := ShapeFor(schema)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer, shape.Schema); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer); err != nil {
		return err
	}
	for _, field := range shape.Fields {
		if err := writeField(writer, field, 0); err != nil {
			return err
		}
	}
	return nil
}

// ShapeFor resolves a schema identifier to its field shape.
func ShapeFor(schema string) (contract.Shape, error) {
	switch schema {
	case "", model.SchemaSearchResultsV1:
		return contract.SearchResultsShape(), nil
	case model.SchemaEnvironmentStatusV1:
		return contract.EnvironmentStatusShape(), nil
	case model.SchemaSnapshotV1:
		return contract.SnapshotShape(), nil
	case model.SchemaAdoptionRecordV1:
		return contract.AdoptionRecordShape(), nil
	case model.SchemaManifestProfileV1:
		return contract.ManifestProfileShape(), nil
	case model.SchemaPackageDataV1:
		return contract.PackageDataShape(), nil
	case model.SchemaSkillEntryV1:
		return contract.SkillEntryShape(), nil
	case model.SchemaSourceResolutionV1:
		return contract.SourceResolutionShape(), nil
	case model.SchemaComparisonV1:
		return contract.ComparisonShape(), nil
	case model.SchemaProbeResultV1:
		return contract.ProbeResultShape(), nil
	case model.SchemaReportV1:
		return contract.ReportShape(), nil
	case model.SchemaRawPayloadV1:
		return contract.RawPayloadShape(), nil
	default:
		return contract.Shape{}, fmt.Errorf("unknown shape schema %q", schema)
	}
}

func writeField(writer io.Writer, field contract.Field, depth int) error {
	prefix := strings.Repeat("  ", depth)
	name := field.Name
	if field.Type == "array" {
		name += "[]"
	}
	if _, err := fmt.Fprintf(writer, "%s%s: %s\n", prefix, name, field.Type); err != nil {
		return err
	}
	for _, child := range field.Fields {
		if err := writeField(writer, child, depth+1); err != nil {
			return err
		}
	}
	return nil
}
