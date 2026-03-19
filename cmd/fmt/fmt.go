package fmt

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "fmt",
	Aliases: []string{"format", "f"},
	Short:   "Code and data formatters",
}
