package cli

import (
	"github.com/MichengLiang/seekmoon/internal/app"
	"github.com/MichengLiang/seekmoon/internal/service"
	"github.com/spf13/cobra"
)

func newRawCommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	cmd := &cobra.Command{
		Use:   "raw <source> ...",
		Short: "Read a raw source payload",
		Args:  argsUnlessContract(cobra.MinimumNArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resolveOutputMode(cmd, &flags); err != nil {
				return err
			}
			if done, err := renderContractIfRequested(cmd, rt, flags, schemaForCommand("raw")); done {
				return err
			}
			result, err := rt.Services.Registry.Raw.Raw(cmd.Context(), service.RawInput{Source: args[0], Args: args[1:]})
			return renderCommand(cmd, rt, flags, schemaForCommand("raw"), result, err)
		},
	}
	addOutputFlags(cmd, &flags)
	return cmd
}
