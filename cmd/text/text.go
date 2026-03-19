package text

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "text",
	Aliases: []string{"txt"},
	Short:   "Text utilities",
}
