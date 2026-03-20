package escape

import (
	"fmt"

	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

// NewHTMLCmd creates an HTML entity escape/unescape command.
func NewHTMLCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "html [input]",
		Short: "HTML entity escape or unescape",
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
			if err != nil {
				return err
			}

			htmlDecode := ioutil.MustGetBool(cmd, "unescape")

			if htmlDecode {
				fmt.Fprintln(cmd.OutOrStdout(), encLib.HTMLDecode(input))
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), encLib.HTMLEncode(input))
			}

			return nil
		},
	}
	cmd.Flags().BoolP("unescape", "u", false, "unescape HTML entities")
	return cmd
}

func init() {
	Cmd.AddCommand(NewHTMLCmd())
}
