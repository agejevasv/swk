package escape

import (
	"fmt"

	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

var urlDecode bool
var urlComponent bool

var urlCmd = &cobra.Command{
	Use:   "url [input]",
	Short: "URL percent-encode or decode",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		if urlDecode {
			result, err := encLib.URLDecode(input, urlComponent)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), result)
		} else {
			result := encLib.URLEncode(input, urlComponent)
			fmt.Fprintln(cmd.OutOrStdout(), result)
		}

		return nil
	},
}

func init() {
	urlCmd.Flags().BoolVarP(&urlDecode, "unescape", "u", false, "decode URL-encoded input")
	urlCmd.Flags().BoolVarP(&urlComponent, "component", "c", false, "use component encoding (QueryEscape)")
	Cmd.AddCommand(urlCmd)
}
