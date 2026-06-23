package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/yumiaura/seekmoon/internal/contract"
	"github.com/yumiaura/seekmoon/internal/model"
)

// RenderSchema writes the JSON Schema contract for a schema identifier.
func RenderSchema(writer io.Writer, schema string) error {
	object, err := SchemaFor(schema)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(object)
}

// SchemaFor resolves a schema identifier to its JSON Schema contract.
func SchemaFor(schema string) (contract.Schema, error) {
	switch schema {
	case "", model.SchemaSearchResultsV1:
		return contract.SearchResultsSchema(), nil
	case model.SchemaEnvironmentStatusV1:
		return contract.EnvironmentStatusSchema(), nil
	case model.SchemaSnapshotV1:
		return contract.SnapshotSchema(), nil
	case model.SchemaAdoptionRecordV1:
		return contract.AdoptionRecordSchema(), nil
	case model.SchemaManifestProfileV1:
		return contract.ManifestProfileSchema(), nil
	case model.SchemaPackageDataV1:
		return contract.PackageDataSchema(), nil
	case model.SchemaSkillEntryV1:
		return contract.SkillEntrySchema(), nil
	case model.SchemaSourceResolutionV1:
		return contract.SourceResolutionSchema(), nil
	case model.SchemaComparisonV1:
		return contract.ComparisonSchema(), nil
	case model.SchemaProbeResultV1:
		return contract.ProbeResultSchema(), nil
	case model.SchemaReportV1:
		return contract.ReportSchema(), nil
	case model.SchemaRawPayloadV1:
		return contract.RawPayloadSchema(), nil
	default:
		return nil, fmt.Errorf("unknown JSON schema %q", schema)
	}
}
