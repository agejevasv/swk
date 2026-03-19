package convert

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"conv", "c"},
	Short:   "Data format converters",
}
