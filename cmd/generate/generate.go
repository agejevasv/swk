package generate

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "generate",
	Args:    cobra.NoArgs,
	RunE:    showHelp,
	Aliases: []string{"g"},
	Short:   "Data generators",
}

// showHelp runs for the bare group command. Its presence makes cobra validate
// arguments, so an unknown subcommand is reported instead of printing help.
func showHelp(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
