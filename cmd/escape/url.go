package escape

import (
	"fmt"

	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

// NewURLCmd creates a URL percent-encode/decode command.
func NewURLCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "url [input]",
		Short: "URL percent-encode or decode",
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
			if err != nil {
				return err
			}

			urlDecode := ioutil.MustGetBool(cmd, "unescape")
			urlComponent := ioutil.MustGetBool(cmd, "component")

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
	cmd.Flags().BoolP("unescape", "u", false, "decode URL-encoded input")
	cmd.Flags().BoolP("component", "c", false, "use component encoding (QueryEscape)")
	return cmd
}

func init() {
	Cmd.AddCommand(NewURLCmd())
}
