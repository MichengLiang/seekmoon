package cli

import (
	"github.com/spf13/cobra"
	"github.com/yumiaura/seekmoon/internal/app"
	"github.com/yumiaura/seekmoon/internal/service"
)

func newDoctorCommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check the local environment",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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
