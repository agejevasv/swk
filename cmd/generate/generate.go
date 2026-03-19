package generate

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen", "g"},
	Short:   "Data generators",
}
