package cli

import (
	"github.com/MichengLiang/seekmoon/internal/app"
	"github.com/MichengLiang/seekmoon/internal/service"
	"github.com/spf13/cobra"
)

func newSearchCommand(rt *app.Runtime) *cobra.Command {
	var flags outputFlags
	var target string
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search library module candidates",
		Args: argsUnlessContract(func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return parseError("search requires a query")
			}
			return nil
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resolveOutputMode(cmd, &flags); err != nil {
				return err
			}
			if done, err := renderContractIfRequested(cmd, rt, flags, schemaForCommand("search")); done {
				return err
			}
			input := service.SearchInput{Kind: "library", Target: target}
			if len(args) > 0 {
				input.Query = args[0]
			}
			result, err := rt.Services.Registry.Search.Search(cmd.Context(), input)
			return renderCommand(cmd, rt, flags, schemaForCommand("search"), result, err)
		},
	}
	cmd.Flags().StringVar(&target, "target", "", "target backend")
	addOutputFlags(cmd, &flags)
	return cmd
}
