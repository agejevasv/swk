package diff

import (
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var Cmd = &cobra.Command{
	Use:     "diff",
	Aliases: []string{"d"},
	Short:   "Compare files",
}

func shouldColorize(cmd *cobra.Command) bool {
	flag, _ := cmd.Flags().GetString("color")
	switch flag {
	case "always":
		return true
	case "never":
		return false
	default: // "auto"
		return term.IsTerminal(int(os.Stdout.Fd()))
	}
}
