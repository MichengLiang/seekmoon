package cli

import (
	"github.com/spf13/cobra"
	"github.com/yumiaura/seekmoon/internal/app"
	"github.com/yumiaura/seekmoon/internal/service"
)

func newReportCommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	var format string
	cmd := &cobra.Command{
		Use:   "report --format <format>",
		Short: "Render an investigation report",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resolveOutputMode(cmd, &flags); err != nil {
				return err
			}
			if done, err := renderContractIfRequested(cmd, rt, flags, schemaForCommand("report")); done {
				return err
			}
			if format == "" {
				return parseError("report requires --format")
			}
			result, err := rt.Services.Registry.Report.Report(cmd.Context(), service.ReportInput{Format: format})
			return renderCommand(cmd, rt, flags, schemaForCommand("report"), result, err)
		},
	}
	cmd.Flags().StringVar(&format, "format", "", "report format")
	addOutputFlags(cmd, &flags)
	return cmd
}
