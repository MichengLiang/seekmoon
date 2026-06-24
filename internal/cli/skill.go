package cli

import (
	"github.com/MichengLiang/seekmoon/internal/app"
	"github.com/MichengLiang/seekmoon/internal/service"
	"github.com/spf13/cobra"
)

func newSkillCommand(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Search or view executable skill entries",
	}
	cmd.AddCommand(newSkillSearchCommand(rt), newSkillViewCommand(rt))
	return cmd
}

func newSkillSearchCommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search skill entries",
		Args:  argsUnlessContract(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resolveOutputMode(cmd, &flags); err != nil {
				return err
			}
			if done, err := renderContractIfRequested(cmd, rt, flags, schemaForCommand("skill")); done {
				return err
			}
			result, err := rt.Services.Registry.Skill.Search(cmd.Context(), service.SkillSearchInput{Query: args[0]})
			return renderCommand(cmd, rt, flags, schemaForCommand("skill"), result, err)
		},
	}
	addOutputFlags(cmd, &flags)
	return cmd
}

func newSkillViewCommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	cmd := &cobra.Command{
		Use:   "view <entry-or-number>",
		Short: "View a skill profile",
		Args:  argsUnlessContract(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resolveOutputMode(cmd, &flags); err != nil {
				return err
			}
			if done, err := renderContractIfRequested(cmd, rt, flags, schemaForCommand("skill")); done {
				return err
			}
			entry, err := parseCandidate(args[0])
			if err != nil {
				return err
			}
			result, err := rt.Services.Registry.Skill.View(cmd.Context(), service.SkillViewInput{Entry: entry})
			return renderCommand(cmd, rt, flags, schemaForCommand("skill"), result, err)
		},
	}
	addOutputFlags(cmd, &flags)
	return cmd
}
