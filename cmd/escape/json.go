package escape

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	textLib "github.com/agejevasv/swk/internal/text"
)

var jsonUnescape bool

var jsonCmd = &cobra.Command{
	Use:   "json [input]",
	Short: "JSON string escape or unescape",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		var result string
		if jsonUnescape {
			result, err = textLib.Unescape(input, "json")
		} else {
			result, err = textLib.Escape(input, "json")
		}
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	jsonCmd.Flags().BoolVarP(&jsonUnescape, "unescape", "u", false, "unescape instead of escape")
	Cmd.AddCommand(jsonCmd)
}
