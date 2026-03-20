package inspect

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "inspect",
	Aliases: []string{"i"},
	Short:   "Inspect and analyze data",
}
