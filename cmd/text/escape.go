package text

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	textLib "github.com/agejevasv/swk/internal/text"
)

var (
	escapeMode     string
	escapeUnescape bool
)

var escapeCmd = &cobra.Command{
	Use:     "escape [text]",
	Aliases: []string{"esc"},
	Short:   "Escape or unescape strings",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		var result string
		if escapeUnescape {
			result, err = textLib.Unescape(input, escapeMode)
		} else {
			result, err = textLib.Escape(input, escapeMode)
		}
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	escapeCmd.Flags().StringVarP(&escapeMode, "mode", "m", "json", "Escape mode (json, xml, html, c, shell)")
	escapeCmd.Flags().BoolVarP(&escapeUnescape, "unescape", "u", false, "Unescape instead of escape")
	Cmd.AddCommand(escapeCmd)
}
