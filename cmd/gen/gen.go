package gen

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "gen",
	Aliases: []string{"generate", "g"},
	Short:   "Data generators",
}
