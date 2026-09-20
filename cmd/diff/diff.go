package diff

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/agejevasv/swk/internal/ioutil"
)

var Cmd = &cobra.Command{
	Use:     "diff",
	Args:    cobra.NoArgs,
	RunE:    showHelp,
	Aliases: []string{"d"},
	Short:   "Compare files",
}

func validateColor(cmd *cobra.Command) error {
	switch ioutil.MustGetString(cmd, "color") {
	case "auto", "always", "never":
		return nil
	default:
		return fmt.Errorf("--color must be auto, always or never")
	}
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

// showHelp runs for the bare group command. Its presence makes cobra validate
// arguments, so an unknown subcommand is reported instead of printing help.
func showHelp(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
