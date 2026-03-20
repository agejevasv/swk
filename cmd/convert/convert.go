package convert

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"c"},
	Short:   "Data format converters",
}
