package cli

import (
	"github.com/MichengLiang/seekmoon/internal/app"
	"github.com/MichengLiang/seekmoon/internal/service"
	"github.com/spf13/cobra"
)

func newDoctorCommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check the local environment",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := resolveOutputMode(cmd, &flags); err != nil {
				return err
			}
			if done, err := renderContractIfRequested(cmd, rt, flags, schemaForCommand("doctor")); done {
				return err
			}
			result, err := rt.Services.Registry.Doctor.Doctor(cmd.Context(), service.DoctorInput{})
			return renderCommand(cmd, rt, flags, schemaForCommand("doctor"), result, err)
		},
	}
	addOutputFlags(cmd, &flags)
	return cmd
}
