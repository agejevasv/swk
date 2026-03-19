package escape

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "escape",
	Aliases: []string{"esc"},
	Short:   "Escape and unescape strings",
}
