package diff

import (
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/agejevasv/swk/internal/ioutil"
)

var Cmd = &cobra.Command{
	Use:     "diff",
	Aliases: []string{"d"},
	Short:   "Compare files",
}

func shouldColorize(cmd *cobra.Command) bool {
	switch ioutil.MustGetString(cmd, "color") {
	case "always":
		return true
	case "never":
		return false
	default: // "auto"
		if f, ok := cmd.OutOrStdout().(*os.File); ok {
			return term.IsTerminal(int(f.Fd()))
		}
		return false
	}
}
