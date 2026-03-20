package escape

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	textLib "github.com/agejevasv/swk/internal/text"
)

var shellCmd = &cobra.Command{
	Use:   "shell [input]",
	Short: "Shell escape or unescape",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		shellUnescape := ioutil.MustGetBool(cmd, "unescape")

		var result string
		if shellUnescape {
			result, err = textLib.Unescape(input, "shell")
		} else {
			result, err = textLib.Escape(input, "shell")
		}
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	shellCmd.Flags().BoolP("unescape", "u", false, "unescape instead of escape")
	Cmd.AddCommand(shellCmd)
}
