package cli

import (
	"github.com/MichengLiang/seekmoon/internal/app"
	"github.com/spf13/cobra"
)

func newSyncCommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Create a data snapshot",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := resolveOutputMode(cmd, &flags); err != nil {
				return err
			}
			if done, err := renderContractIfRequested(cmd, rt, flags, schemaForCommand("sync")); done {
				return err
			}
			result, err := rt.Services.Registry.Sync.Sync(cmd.Context())
			return renderCommand(cmd, rt, flags, schemaForCommand("sync"), result, err)
		},
	}
	addOutputFlags(cmd, &flags)
	return cmd
}
