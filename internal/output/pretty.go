package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/yumiaura/seekmoon/internal/model"
)

// RenderPretty writes the default human-oriented command output.
func RenderPretty(writer io.Writer, value any) error {
	switch v := value.(type) {
	case model.SearchOutput:
		return renderSearch(writer, v)
	case *model.SearchOutput:
		if v == nil {
			return nil
		}
		return renderSearch(writer, *v)
	case model.Snapshot:
		_, err := fmt.Fprintf(writer, "snapshot  %s\nsources   %d\n", v.ID, len(v.Sources))
		return err
	case *model.Snapshot:
		if v == nil {
			return nil
		}
		_, err := fmt.Fprintf(writer, "snapshot  %s\nsources   %d\n", v.ID, len(v.Sources))
		return err
	case model.EnvironmentStatus:
		return renderEnvironment(writer, v)
	case *model.EnvironmentStatus:
		if v == nil {
			return nil
		}
		return renderEnvironment(writer, *v)
	case model.ManifestProfile:
		return renderManifest(writer, v)
	case *model.ManifestProfile:
		if v == nil {
			return nil
		}
		return renderManifest(writer, *v)
	case model.PackageData:
		return renderPackageData(writer, v)
	case *model.PackageData:
		if v == nil {
			return nil
		}
		return renderPackageData(writer, *v)
	case model.SourceResolution:
		return renderSourceResolution(writer, v)
	case *model.SourceResolution:
		if v == nil {
			return nil
		}
		return renderSourceResolution(writer, *v)
	case []model.SkillEntry:
		return renderSkillEntries(writer, v)
	case model.SkillProfile:
		return renderSkillProfile(writer, v)
	case *model.SkillProfile:
		if v == nil {
			return nil
		}
		return renderSkillProfile(writer, *v)
	case model.Comparison:
		return renderComparison(writer, v)
	case *model.Comparison:
		if v == nil {
			return nil
		}
		return renderComparison(writer, *v)
	case model.ProbeResult:
		return renderProbe(writer, v)
	case *model.ProbeResult:
		if v == nil {
			return nil
		}
		return renderProbe(writer, *v)
	case model.AdoptionRecord:
		return renderRecord(writer, v)
	case *model.AdoptionRecord:
		if v == nil {
			return nil
		}
		return renderRecord(writer, *v)
	case model.Report:
		_, err := fmt.Fprintf(writer, "report  %s\nsources %s\n", v.Goal, strings.Join(v.DataSources, ","))
		return err
	case *model.Report:
		if v == nil {
			return nil
		}
		_, err := fmt.Fprintf(writer, "report  %s\nsources %s\n", v.Goal, strings.Join(v.DataSources, ","))
		return err
	case model.RawEnvelope:
		_, err := fmt.Fprintf(writer, "raw     %s\nstate   %s\nsource  %s\n", v.Source, v.Status, firstPretty(v.URL, v.Path, v.RawRef))
		return err
	case *model.RawEnvelope:
		if v == nil {
			return nil
		}
		_, err := fmt.Fprintf(writer, "raw     %s\nstate   %s\nsource  %s\n", v.Source, v.Status, firstPretty(v.URL, v.Path, v.RawRef))
		return err
	default:
		_, err := fmt.Fprintf(writer, "%v\n", value)
		return err
	}
}

func renderEnvironment(writer io.Writer, status model.EnvironmentStatus) error {
	if _, err := fmt.Fprintf(writer, "moon        %s\n", commandStatus(status.Toolchain["moon"])); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "registry    %s\nnetwork     %s\nproject     %s\n", evidenceText(status.Paths["registry_index"]), evidenceText(status.Network["mooncakes_api"]), status.Project.Status); err != nil {
		return err
	}
	return nil
}

func renderSearch(writer io.Writer, output model.SearchOutput) error {
	if _, err := fmt.Fprintf(writer, "Search: %s    target: %s    kind: %s    snapshot: %s\n\n", output.Query.Text, output.Query.Target, output.Query.Kind, output.Snapshot.ID); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer, "#  module                         version   license   build    downloads   target"); err != nil {
		return err
	}
	for _, result := range output.Results {
		license := evidenceText(result.License)
		build := evidenceText(result.BuildStatus)
		downloads := evidenceIntText(result.Downloads)
		target := evidenceObjectText(result.Target)
		if _, err := fmt.Fprintf(writer, "%-2d %-30s %-9s %-9s %-8s %-11s %s\n", result.Rank, result.Module, result.Version, license, build, downloads, target); err != nil {
			return err
		}
	}
	return nil
}

func renderManifest(writer io.Writer, profile model.ManifestProfile) error {
	_, err := fmt.Fprintf(
		writer, "%s %s\n\ndescription  %s\nlicense      %s\nrepository   %s\ndownloads    %d\nbuild        %s\ndocs         %s\n",
		profile.Module,
		profile.Version,
		evidenceText(profile.Metadata.Description),
		evidenceText(profile.Metadata.License),
		evidenceText(profile.Metadata.Repository),
		profile.Downloads,
		evidenceText(profile.BuildStatus),
		evidenceText(profile.DocsURL),
	)
	return err
}

func renderPackageData(writer io.Writer, data model.PackageData) error {
	if _, err := fmt.Fprintf(writer, "API: %s\n\n", data.Name); err != nil {
		return err
	}
	if err := renderAPISection(writer, "types", data.Types); err != nil {
		return err
	}
	return renderAPISection(writer, "functions", data.Values)
}

func renderSourceResolution(writer io.Writer, resolution model.SourceResolution) error {
	_, err := fmt.Fprintf(
		writer,
		"source  %s@%s\nstatus  %s\nmethod  %s\npath    %s\n",
		resolution.Module,
		resolution.Version,
		resolution.SelectedSource.Method,
		resolution.SelectedSource.Method,
		firstPretty(resolution.SelectedSource.Path, resolution.SelectedSource.URL),
	)
	return err
}

func renderSkillEntries(writer io.Writer, entries []model.SkillEntry) error {
	if _, err := fmt.Fprintln(writer, "#  skill                         version   package   wasm       checksum"); err != nil {
		return err
	}
	for index, entry := range entries {
		if _, err := fmt.Fprintf(writer, "%-2d %-29s %-9s %-9s %-10s %s\n", index+1, entry.Module, entry.Version, entry.Package, stateByPresence(entry.WasmURL), stateByPresence(entry.ChecksumURL)); err != nil {
			return err
		}
	}
	return nil
}

func renderSkillProfile(writer io.Writer, profile model.SkillProfile) error {
	_, err := fmt.Fprintf(
		writer,
		"skill   %s\nversion %s\npackage %s\nwasm    %s\nsha256  %s\nrun     %s\n",
		profile.Entry.Module,
		profile.Entry.Version,
		profile.Entry.Package,
		profile.WasmAsset.Status,
		profile.ChecksumAsset.Status,
		evidenceText(profile.RunwasmCoordinate),
	)
	return err
}

func renderComparison(writer io.Writer, comparison model.Comparison) error {
	for _, field := range comparison.Fields {
		if _, err := fmt.Fprintf(writer, "%s", field.Name); err != nil {
			return err
		}
		for _, candidate := range comparison.Candidates {
			if _, err := fmt.Fprintf(writer, "  %s=%s", candidate.Module, field.Values[candidate.Module]); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
	}
	return nil
}

func renderProbe(writer io.Writer, result model.ProbeResult) error {
	_, err := fmt.Fprintf(writer, "probe   %s@%s\ntarget  %s\npath    %s\nstatus  %s\n", result.Module, result.Version, result.Target, result.ProbePath, result.Result)
	return err
}

func renderRecord(writer io.Writer, record model.AdoptionRecord) error {
	_, err := fmt.Fprintf(writer, "recorded  %s@%s\nstatus    %s\n", record.Candidate.Module, record.Version, record.Conclusion)
	return err
}

func renderAPISection(writer io.Writer, label string, entries []model.APIEntry) error {
	if len(entries) == 0 {
		return nil
	}
	if _, err := fmt.Fprintln(writer, label); err != nil {
		return err
	}
	for _, entry := range entries {
		text := entry.Name
		if entry.PlainSignature.Value != nil && *entry.PlainSignature.Value != "" {
			text = *entry.PlainSignature.Value
		}
		if _, err := fmt.Fprintf(writer, "  %s\n", text); err != nil {
			return err
		}
	}
	return nil
}

func evidenceText(e model.EvidenceString) string {
	if e.Value != nil && *e.Value != "" {
		return *e.Value
	}
	if e.Status != "" {
		return string(e.Status)
	}
	return "unknown"
}

func evidenceIntText(e model.EvidenceInt) string {
	if e.Value != nil {
		return fmt.Sprintf("%d", *e.Value)
	}
	if e.Status != "" {
		return string(e.Status)
	}
	return "unknown"
}

func evidenceObjectText(e model.EvidenceObject) string {
	if e.Value != nil && len(*e.Value) > 0 {
		parts := make([]string, 0, len(*e.Value))
		for key := range *e.Value {
			parts = append(parts, key)
		}
		return strings.Join(parts, ",")
	}
	if e.Status != "" {
		return string(e.Status)
	}
	return "unknown"
}

func commandStatus(result model.CommandResult) string {
	if result.Status != "" {
		return string(result.Status)
	}
	return "unknown"
}

func stateByPresence(value string) string {
	if value == "" {
		return string(model.StateMissing)
	}
	return string(model.StatePresent)
}

func firstPretty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return "unknown"
}
