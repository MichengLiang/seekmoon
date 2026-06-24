package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MichengLiang/seekmoon/internal/app"
	"github.com/MichengLiang/seekmoon/internal/model"
	"github.com/MichengLiang/seekmoon/internal/output"
	"github.com/MichengLiang/seekmoon/internal/service"
)

func TestRootCommandRegistersBatchCCommands(t *testing.T) {
	rt, err := app.NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime: %v", err)
	}
	var out bytes.Buffer
	if code := ExecuteWithCode(context.Background(), rt, Options{Out: &out, Err: &out}); code != exitCodeOK {
		t.Fatalf("root exit code = %d output=%q", code, out.String())
	}
	for _, name := range []string{"doctor", "sync", "search", "view", "api", "source", "skill", "compare", "probe", "record", "report", "raw"} {
		if !strings.Contains(out.String(), name) {
			t.Fatalf("help missing %q: %s", name, out.String())
		}
	}
}

func TestRequiredArgumentFailureMapsToExitCode2(t *testing.T) {
	rt, err := app.NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime: %v", err)
	}
	var out bytes.Buffer
	if code := ExecuteWithCode(context.Background(), rt, Options{Out: &out, Err: &out}, "search"); code != exitCodeUsage {
		t.Fatalf("exit code = %d output=%q", code, out.String())
	}
}

func TestJQWithoutExpressionMapsToExitCode2(t *testing.T) {
	rt, err := app.NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime: %v", err)
	}
	var out bytes.Buffer
	if code := ExecuteWithCode(context.Background(), rt, Options{Out: &out, Err: &out}, "search", "markdown", "--jq"); code != exitCodeUsage {
		t.Fatalf("exit code = %d output=%q", code, out.String())
	}
}

func TestServiceErrorRendersAndMapsToExitCode1(t *testing.T) {
	rt, err := app.NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime: %v", err)
	}
	rt.Services.Registry.Search = &fakeSearchService{err: errors.New("search unavailable")}
	var out bytes.Buffer
	if code := ExecuteWithCode(context.Background(), rt, Options{Out: &out, Err: &out}, "search", "markdown"); code != exitCodeError {
		t.Fatalf("exit code = %d output=%q", code, out.String())
	}
	if !strings.Contains(out.String(), "search unavailable") || !strings.Contains(out.String(), "state") {
		t.Fatalf("service error surface = %q", out.String())
	}
}

func TestContractProjectionModesBypassOperandsAndPendingServices(t *testing.T) {
	cases := []struct {
		name   string
		args   []string
		schema string
	}{
		{name: "search shape", args: []string{"search", "--shape"}, schema: model.SchemaSearchResultsV1},
		{name: "view schema", args: []string{"view", "--schema"}, schema: model.SchemaManifestProfileV1},
		{name: "api shape", args: []string{"api", "--shape"}, schema: model.SchemaPackageDataV1},
		{name: "source shape", args: []string{"source", "--shape"}, schema: model.SchemaSourceResolutionV1},
		{name: "skill search shape", args: []string{"skill", "search", "--shape"}, schema: model.SchemaSkillEntryV1},
		{name: "skill view schema", args: []string{"skill", "view", "--schema"}, schema: model.SchemaSkillEntryV1},
		{name: "compare shape", args: []string{"compare", "--shape"}, schema: model.SchemaComparisonV1},
		{name: "probe shape", args: []string{"probe", "--shape"}, schema: model.SchemaProbeResultV1},
		{name: "record shape", args: []string{"record", "--shape"}, schema: model.SchemaAdoptionRecordV1},
		{name: "report shape", args: []string{"report", "--shape"}, schema: model.SchemaReportV1},
		{name: "raw shape", args: []string{"raw", "--shape"}, schema: model.SchemaRawPayloadV1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rt, err := app.NewRuntime()
			if err != nil {
				t.Fatalf("NewRuntime: %v", err)
			}
			var out bytes.Buffer
			if code := ExecuteWithCode(context.Background(), rt, Options{Out: &out, Err: &out}, tc.args...); code != exitCodeOK {
				t.Fatalf("exit code = %d output=%q", code, out.String())
			}
			if !strings.Contains(out.String(), tc.schema) {
				t.Fatalf("contract output missing schema %q: %s", tc.schema, out.String())
			}
			if strings.Contains(out.String(), "outside Batch C") || strings.Contains(out.String(), "requires ") || strings.Contains(out.String(), "accepts ") {
				t.Fatalf("contract mode ran ordinary command path: %s", out.String())
			}
		})
	}
}

func TestSearchThinHandlerCallsServiceOnce(t *testing.T) {
	rt, err := app.NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime: %v", err)
	}
	search := &fakeSearchService{result: model.SearchOutput{
		Schema:   model.SchemaSearchResultsV1,
		Snapshot: model.SnapshotRef{ID: "snapshot", Sources: []string{"modules_api"}},
		Query:    model.SearchQuery{Text: "markdown", Kind: "library"},
	}}
	rt.Services.Registry.Search = search
	var rendered output.Request
	rt.Renderer = output.RendererFunc(func(_ context.Context, request output.Request) error {
		rendered = request
		return nil
	})
	if code := ExecuteWithCode(context.Background(), rt, Options{}, "search", "markdown", "--target", "js", "--json"); code != exitCodeOK {
		t.Fatalf("exit code = %d", code)
	}
	if search.calls != 1 || search.input.Query != "markdown" || search.input.Target != "js" {
		t.Fatalf("search calls=%d input=%#v", search.calls, search.input)
	}
	if rendered.Mode != model.OutputJSON || rendered.Schema != model.SchemaSearchResultsV1 {
		t.Fatalf("rendered request = %#v", rendered)
	}
}

func TestCandidateNumberAndModuleInputsAreDistinct(t *testing.T) {
	number, err := parseCandidate("1")
	if err != nil {
		t.Fatalf("parse number: %v", err)
	}
	module, err := parseCandidate("moonbitlang/core@0.1.0")
	if err != nil {
		t.Fatalf("parse module: %v", err)
	}
	if number.Number != 1 || number.Module != "" {
		t.Fatalf("number candidate = %#v", number)
	}
	if module.Module != "moonbitlang/core" || module.Version != "0.1.0" || module.Number != 0 {
		t.Fatalf("module candidate = %#v", module)
	}
}

type fakeSearchService struct {
	calls  int
	input  service.SearchInput
	result model.SearchOutput
	err    error
}

func (s *fakeSearchService) Search(_ context.Context, input service.SearchInput) (model.SearchOutput, error) {
	s.calls++
	s.input = input
	return s.result, s.err
}
