package escape

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	textLib "github.com/agejevasv/swk/internal/text"
)

var xmlCmd = &cobra.Command{
	Use:   "xml [input]",
	Short: "XML escape or unescape",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		xmlUnescape, _ := cmd.Flags().GetBool("unescape")

		var result string
		if xmlUnescape {
			result, err = textLib.Unescape(input, "xml")
		} else {
			result, err = textLib.Escape(input, "xml")
		}
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	xmlCmd.Flags().BoolP("unescape", "u", false, "unescape instead of escape")
	Cmd.AddCommand(xmlCmd)
}
