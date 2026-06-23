package cli

import (
	"github.com/spf13/cobra"
	"github.com/yumiaura/seekmoon/internal/app"
	"github.com/yumiaura/seekmoon/internal/service"
)

func newAPICommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	var packagePath string
	cmd := &cobra.Command{
		Use:   "api <module-or-number> --package <path>",
		Short: "View package API profile",
		Args:  argsUnlessContract(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resolveOutputMode(cmd, &flags); err != nil {
				return err
			}
			if done, err := renderContractIfRequested(cmd, rt, flags, schemaForCommand("api")); done {
				return err
			}
			if packagePath == "" {
				return parseError("api requires --package")
			}
			candidate, err := parseCandidate(args[0])
			if err != nil {
				return err
			}
			result, err := rt.Services.Registry.API.API(cmd.Context(), service.APIInput{Candidate: candidate, Package: packagePath})
			return renderCommand(cmd, rt, flags, schemaForCommand("api"), result, err)
		},
	}
	cmd.Flags().StringVar(&packagePath, "package", "", "package path")
	addOutputFlags(cmd, &flags)
	return cmd
}
