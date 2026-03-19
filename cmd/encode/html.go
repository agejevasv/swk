package encode

import (
	"fmt"

	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

var htmlDecode bool

var htmlCmd = &cobra.Command{
	Use:   "html [input]",
	Short: "HTML encode or decode",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		if htmlDecode {
			result := encLib.HTMLDecode(input)
			fmt.Fprintln(cmd.OutOrStdout(), result)
		} else {
			result := encLib.HTMLEncode(input)
			fmt.Fprintln(cmd.OutOrStdout(), result)
		}

		return nil
	},
}

func init() {
	htmlCmd.Flags().BoolVarP(&htmlDecode, "decode", "d", false, "decode HTML entities")

	Cmd.AddCommand(htmlCmd)
}
