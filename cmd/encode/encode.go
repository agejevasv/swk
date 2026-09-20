package encode

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "encode",
	Args:    cobra.NoArgs,
	RunE:    showHelp,
	Aliases: []string{"enc"},
	Short:   "Encoders and decoders",
}

// showHelp runs for the bare group command. Its presence makes cobra validate
// arguments, so an unknown subcommand is reported instead of printing help.
func showHelp(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
