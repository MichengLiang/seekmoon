package cli

import (
	"github.com/MichengLiang/seekmoon/internal/app"
	"github.com/MichengLiang/seekmoon/internal/model"
	"github.com/MichengLiang/seekmoon/internal/service"
	"github.com/spf13/cobra"
)

func newRecordCommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	var conclusion string
	var note string
	var kind string
	cmd := &cobra.Command{
		Use:   "record <candidate> --conclusion <value>",
		Short: "Save an adoption judgment",
		Args:  argsUnlessContract(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resolveOutputMode(cmd, &flags); err != nil {
				return err
			}
			if done, err := renderContractIfRequested(cmd, rt, flags, schemaForCommand("record")); done {
				return err
			}
			if conclusion == "" {
				return parseError("record requires --conclusion")
			}
			parsed, err := model.ParseAdoptionConclusion(conclusion)
			if err != nil {
				return parseError("%s", err.Error())
			}
			candidate, err := parseCandidate(args[0])
			if err != nil {
				return err
			}
			result, err := rt.Services.Registry.Record.Record(cmd.Context(), service.RecordInput{
				Candidate:  candidate,
				Kind:       kind,
				Conclusion: parsed,
				Note:       note,
			})
			return renderCommand(cmd, rt, flags, schemaForCommand("record"), result, err)
		},
	}
	cmd.Flags().StringVar(&conclusion, "conclusion", "", "adoption conclusion")
	cmd.Flags().StringVar(&note, "note", "", "record note")
	cmd.Flags().StringVar(&kind, "kind", "library", "candidate kind")
	addOutputFlags(cmd, &flags)
	return cmd
}
