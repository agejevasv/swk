package format

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "format",
	Aliases: []string{"fmt", "f"},
	Short:   "Prettify, minify, and render data for display",
}
