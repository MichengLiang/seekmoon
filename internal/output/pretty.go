package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/yumiaura/seekmoon/internal/model"
)

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
	case model.ProbeResult:
		_, err := fmt.Fprintf(writer, "probe   %s@%s\ntarget  %s\nstatus  %s\n", v.Module, v.Version, v.Target, v.Result)
		return err
	case *model.ProbeResult:
		if v == nil {
			return nil
		}
		_, err := fmt.Fprintf(writer, "probe   %s@%s\ntarget  %s\nstatus  %s\n", v.Module, v.Version, v.Target, v.Result)
		return err
	default:
		_, err := fmt.Fprintf(writer, "%v\n", value)
		return err
	}
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
		var parts []string
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
