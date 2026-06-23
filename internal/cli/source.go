package cli

import (
	"github.com/spf13/cobra"
	"github.com/yumiaura/seekmoon/internal/app"
	"github.com/yumiaura/seekmoon/internal/service"
)

func newSourceCommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	cmd := &cobra.Command{
		Use:   "source <module-or-number[@version]>",
		Short: "Locate published source",
		Args:  argsUnlessContract(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resolveOutputMode(cmd, &flags); err != nil {
				return err
			}
			if done, err := renderContractIfRequested(cmd, rt, flags, schemaForCommand("source")); done {
				return err
			}
			candidate, err := parseCandidate(args[0])
			if err != nil {
				return err
			}
			result, err := rt.Services.Registry.Source.Source(cmd.Context(), service.SourceInput{Candidate: candidate})
			return renderCommand(cmd, rt, flags, schemaForCommand("source"), result, err)
		},
	}
	addOutputFlags(cmd, &flags)
	return cmd
}
