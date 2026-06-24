package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MichengLiang/seekmoon/internal/app"
	"github.com/MichengLiang/seekmoon/internal/model"
	"github.com/spf13/cobra"
)

const (
	exitCodeOK    = 0
	exitCodeError = 1
	exitCodeUsage = 2
)

type outputFlags struct {
	json       bool
	shape      bool
	schema     bool
	jq         string
	jqProvided bool
	mode       model.OutputMode
}

func addOutputFlags(cmd *cobra.Command, flags *outputFlags) {
	cmd.Flags().BoolVar(&flags.json, "json", false, "render JSON output")
	cmd.Flags().StringVar(&flags.jq, "jq", "", "evaluate a jq expression against JSON output")
	cmd.Flags().BoolVar(&flags.shape, "shape", false, "render the command output shape")
	cmd.Flags().BoolVar(&flags.schema, "schema", false, "render the command JSON Schema")
}

func resolveOutputMode(cmd *cobra.Command, flags *outputFlags) error {
	flags.jqProvided = cmd.Flags().Changed("jq")
	selected := 0
	if flags.json {
		selected++
		flags.mode = model.OutputJSON
	}
	if flags.jqProvided {
		selected++
		flags.mode = model.OutputJQ
		if strings.TrimSpace(flags.jq) == "" {
			return parseError("--jq requires an expression")
		}
	}
	if flags.shape {
		selected++
		flags.mode = model.OutputShape
	}
	if flags.schema {
		selected++
		flags.mode = model.OutputSchema
	}
	if selected > 1 {
		return parseError("choose only one output mode")
	}
	if selected == 0 {
		flags.mode = model.OutputPretty
	}
	return nil
}

func contractProjectionRequested(cmd *cobra.Command) bool {
	return cmd.Flags().Changed("shape") || cmd.Flags().Changed("schema")
}

func argsUnlessContract(inner cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if contractProjectionRequested(cmd) {
			return nil
		}
		return inner(cmd, args)
	}
}

func renderContractIfRequested(cmd *cobra.Command, rt *app.Runtime, flags outputFlags, schema string) (bool, error) {
	if flags.mode != model.OutputShape && flags.mode != model.OutputSchema {
		return false, nil
	}
	return true, renderCommand(cmd, rt, flags, schema, nil, nil)
}

type parseFailure struct {
	message string
}

func (e parseFailure) Error() string {
	return e.message
}

func parseError(format string, args ...any) error {
	return parseFailure{message: fmt.Sprintf(format, args...)}
}

func isParseFailure(err error) bool {
	_, ok := err.(parseFailure)
	return ok
}

func parseCandidate(raw string) (model.CandidateRequest, error) {
	if raw == "" {
		return model.CandidateRequest{}, parseError("candidate is required")
	}
	if n, err := strconv.Atoi(raw); err == nil {
		if n <= 0 {
			return model.CandidateRequest{}, parseError("candidate number must be positive")
		}
		return model.CandidateRequest{Raw: raw, Number: n}, nil
	}
	module := raw
	version := ""
	if before, after, ok := strings.Cut(raw, "@"); ok {
		module = before
		version = after
	}
	if _, err := model.ParseModuleCoordinate(module); err != nil {
		return model.CandidateRequest{}, parseError("%s", err.Error())
	}
	return model.CandidateRequest{Raw: raw, Module: module, Version: version}, nil
}

func schemaForCommand(name string) string {
	switch name {
	case "doctor":
		return model.SchemaEnvironmentStatusV1
	case "sync":
		return model.SchemaSnapshotV1
	case "search":
		return model.SchemaSearchResultsV1
	case "view":
		return model.SchemaManifestProfileV1
	case "api":
		return model.SchemaPackageDataV1
	case "source":
		return model.SchemaSourceResolutionV1
	case "skill":
		return model.SchemaSkillEntryV1
	case "compare":
		return model.SchemaComparisonV1
	case "probe":
		return model.SchemaProbeResultV1
	case "record":
		return model.SchemaAdoptionRecordV1
	case "report":
		return model.SchemaReportV1
	case "raw":
		return model.SchemaRawPayloadV1
	default:
		return model.SchemaSearchResultsV1
	}
}
