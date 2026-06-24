package cli

import (
	"github.com/MichengLiang/seekmoon/internal/app"
	"github.com/MichengLiang/seekmoon/internal/model"
	"github.com/MichengLiang/seekmoon/internal/service"
	"github.com/spf13/cobra"
)

func newCompareCommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	cmd := &cobra.Command{
		Use:   "compare <candidate>...",
		Short: "Compare candidate evidence",
		Args:  argsUnlessContract(cobra.MinimumNArgs(2)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resolveOutputMode(cmd, &flags); err != nil {
				return err
			}
			if done, err := renderContractIfRequested(cmd, rt, flags, schemaForCommand("compare")); done {
				return err
			}
			candidates := make([]model.CandidateRequest, 0, len(args))
			for _, arg := range args {
				candidate, err := parseCandidate(arg)
				if err != nil {
					return err
				}
				candidates = append(candidates, candidate)
			}
			result, err := rt.Services.Registry.Compare.Compare(cmd.Context(), service.CompareInput{Candidates: candidates})
			return renderCommand(cmd, rt, flags, schemaForCommand("compare"), result, err)
		},
	}
	addOutputFlags(cmd, &flags)
	return cmd
}
