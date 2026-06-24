package cli

import (
	"github.com/MichengLiang/seekmoon/internal/app"
	"github.com/MichengLiang/seekmoon/internal/service"
	"github.com/spf13/cobra"
)

func newProbeCommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	var target string
	cmd := &cobra.Command{
		Use:   "probe <candidate>",
		Short: "Run local validation for a candidate",
		Args:  argsUnlessContract(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resolveOutputMode(cmd, &flags); err != nil {
				return err
			}
			if done, err := renderContractIfRequested(cmd, rt, flags, schemaForCommand("probe")); done {
				return err
			}
			candidate, err := parseCandidate(args[0])
			if err != nil {
				return err
			}
			result, err := rt.Services.Registry.Probe.Probe(cmd.Context(), service.ProbeInput{Candidate: candidate, Target: target})
			return renderCommand(cmd, rt, flags, schemaForCommand("probe"), result, err)
		},
	}
	cmd.Flags().StringVar(&target, "target", "", "target backend")
	addOutputFlags(cmd, &flags)
	return cmd
}
