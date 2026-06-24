package helpdoc

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Apply injects authored English help text into a Cobra command tree.
func Apply(root *cobra.Command) error {
	docs := Docs()
	commonFlags := CommonFlagDocs()
	var missing []string
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		if cmd.Hidden || isCobraInfrastructure(cmd) {
			return
		}
		key := cmd.CommandPath()
		doc, ok := docs[key]
		if ok {
			cmd.Short = doc.ShortEN
			cmd.Long = doc.LongEN
			cmd.Example = doc.ExampleEN
			applyFlagDocs(cmd, commonFlags)
			applyFlagDocs(cmd, doc.Flags)
		} else if isSeekMoonAuthored(cmd) {
			missing = append(missing, key)
		}
		for _, child := range cmd.Commands() {
			walk(child)
		}
	}
	walk(root)
	if len(missing) > 0 {
		return fmt.Errorf("missing help docs for %v", missing)
	}
	return nil
}

func applyFlagDocs(cmd *cobra.Command, flags map[string]FlagDoc) {
	for name, doc := range flags {
		if doc.UsageEN == "" {
			continue
		}
		if flag := cmd.Flags().Lookup(name); flag != nil {
			flag.Usage = doc.UsageEN
		}
		if flag := cmd.PersistentFlags().Lookup(name); flag != nil {
			flag.Usage = doc.UsageEN
		}
	}
}

func isSeekMoonAuthored(cmd *cobra.Command) bool {
	if isCobraInfrastructure(cmd) {
		return false
	}
	return true
}

func isCobraInfrastructure(cmd *cobra.Command) bool {
	switch cmd.Name() {
	case "completion", "help":
		return true
	default:
		return false
	}
}
