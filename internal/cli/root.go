// Package cli owns the Cobra command surface and keeps handlers thin.
package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/MichengLiang/seekmoon/internal/app"
	"github.com/MichengLiang/seekmoon/internal/output"
	"github.com/spf13/cobra"
)

// Options configures command output streams.
type Options struct {
	Out io.Writer
	Err io.Writer
}

// NewRoot builds the root Cobra command and subcommands.
func NewRoot(rt *app.Runtime, options Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "seekmoon",
		Short:         "MoonBit package discovery workbench",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	if options.Out != nil {
		cmd.SetOut(options.Out)
	}
	if options.Err != nil {
		cmd.SetErr(options.Err)
	}
	cmd.SetContext(context.Background())
	cmd.AddCommand(
		newDoctorCommand(rt),
		newSyncCommand(rt),
		newSearchCommand(rt),
		newViewCommand(rt),
		newAPICommand(rt),
		newSourceCommand(rt),
		newSkillCommand(rt),
		newCompareCommand(rt),
		newProbeCommand(rt),
		newRecordCommand(rt),
		newReportCommand(rt),
		newRawCommand(rt),
	)
	return cmd
}

// Execute runs the root command with the supplied context.
func Execute(ctx context.Context, rt *app.Runtime, options Options) error {
	root := NewRoot(rt, options)
	root.SetContext(ctx)
	return root.Execute()
}

// ExecuteWithCode runs the root command and returns a process-style exit code.
func ExecuteWithCode(ctx context.Context, rt *app.Runtime, options Options, args ...string) int {
	root := NewRoot(rt, options)
	root.SetContext(ctx)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		if options.Err != nil {
			_, _ = fmt.Fprintln(options.Err, err)
		}
		if isUsageFailure(err) {
			return exitCodeUsage
		}
		return exitCodeError
	}
	return exitCodeOK
}

func isUsageFailure(err error) bool {
	if isParseFailure(err) {
		return true
	}
	message := err.Error()
	return strings.Contains(message, "unknown flag") ||
		strings.Contains(message, "flag needs an argument") ||
		strings.Contains(message, "accepts ") ||
		strings.Contains(message, "requires ")
}

func renderCommand(cmd *cobra.Command, rt *app.Runtime, flags outputFlags, schema string, value any, err error) error {
	if rt.Renderer == nil {
		return fmt.Errorf("runtime renderer is not configured")
	}
	return rt.Renderer.Render(cmd.Context(), outputRequest(cmd.CommandPath(), flags, schema, cmd.OutOrStdout(), value, err))
}

func outputRequest(command string, flags outputFlags, schema string, writer io.Writer, value any, err error) output.Request {
	return output.Request{
		Command:      command,
		Mode:         flags.mode,
		JQExpression: flags.jq,
		Schema:       schema,
		Writer:       writer,
		Value:        value,
		Err:          err,
	}
}
