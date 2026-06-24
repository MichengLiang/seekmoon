package cli

import (
	"github.com/MichengLiang/seekmoon/internal/app"
	"github.com/MichengLiang/seekmoon/internal/service"
	"github.com/spf13/cobra"
)

func newViewCommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	cmd := &cobra.Command{
		Use:   "view <module-or-number>",
		Short: "View a library module profile",
		Args:  argsUnlessContract(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resolveOutputMode(cmd, &flags); err != nil {
				return err
			}
			if done, err := renderContractIfRequested(cmd, rt, flags, schemaForCommand("view")); done {
				return err
			}
			candidate, err := parseCandidate(args[0])
			if err != nil {
				return err
			}
			result, err := rt.Services.Registry.View.View(cmd.Context(), service.ViewInput{Candidate: candidate})
			return renderCommand(cmd, rt, flags, schemaForCommand("view"), result, err)
		},
	}
	addOutputFlags(cmd, &flags)
	return cmd
}
