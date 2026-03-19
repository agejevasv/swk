package encode

import (
	"fmt"

	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

var base64Decode bool
var base64URLSafe bool
var base64NoPadding bool

var base64Cmd = &cobra.Command{
	Use:     "base64 [input]",
	Aliases: []string{"b64"},
	Short:   "Base64 encode or decode",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		if base64Decode {
			result, err := encLib.Base64Decode(input, base64URLSafe)
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(result)
			return err
		}
		result := encLib.Base64Encode([]byte(input), base64URLSafe, base64NoPadding)
		fmt.Fprintln(cmd.OutOrStdout(), result)

		return nil
	},
}

func init() {
	base64Cmd.Flags().BoolVarP(&base64Decode, "decode", "d", false, "decode base64 input")
	base64Cmd.Flags().BoolVarP(&base64URLSafe, "url-safe", "u", false, "use URL-safe encoding")
	base64Cmd.Flags().BoolVar(&base64NoPadding, "no-padding", false, "omit padding characters")

	Cmd.AddCommand(base64Cmd)
}
