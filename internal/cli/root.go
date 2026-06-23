// Package cli owns the Cobra command surface and keeps handlers thin.
package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/yumiaura/seekmoon/internal/app"
)

type Options struct {
	Out io.Writer
	Err io.Writer
}

func NewRoot(rt *app.Runtime, options Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "seekmoon",
		Short:         "MoonBit package discovery workbench",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "seekmoon: command surface pending Batch C")
			return nil
		},
	}
	if options.Out != nil {
		cmd.SetOut(options.Out)
	}
	if options.Err != nil {
		cmd.SetErr(options.Err)
	}
	cmd.SetContext(context.Background())
	cmd.AddCommand(
		placeholder("doctor"),
		placeholder("sync"),
		placeholder("search"),
		placeholder("view"),
		placeholder("api"),
		placeholder("source"),
		placeholder("skill"),
		placeholder("compare"),
		placeholder("probe"),
		placeholder("record"),
		placeholder("report"),
		placeholder("raw"),
	)
	_ = rt
	return cmd
}

func Execute(ctx context.Context, rt *app.Runtime, options Options) error {
	root := NewRoot(rt, options)
	root.SetContext(ctx)
	return root.Execute()
}

func placeholder(use string) *cobra.Command {
	return &cobra.Command{
		Use:          use,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("%s command behavior is outside Batch A", cmd.CommandPath())
		},
	}
}
