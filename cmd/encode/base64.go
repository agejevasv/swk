package encode

import (
	"fmt"

	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

var base64Cmd = &cobra.Command{
	Use:     "base64 [input]",
	Aliases: []string{"b64"},
	Short:   "Base64 encode or decode",
	RunE: func(cmd *cobra.Command, args []string) error {
		decode := ioutil.MustGetBool(cmd, "decode")
		urlSafe := ioutil.MustGetBool(cmd, "url-safe")
		noPadding := ioutil.MustGetBool(cmd, "no-padding")

		if decode {
			// Decode expects a base64 string — trim whitespace but not binary data.
			input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
			if err != nil {
				return err
			}
			result, err := encLib.Base64Decode(input, urlSafe)
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(result)
			return err
		}

		// Encode: use raw bytes to preserve all input including trailing newlines.
		input, err := ioutil.ReadFileInput(args, cmd.InOrStdin())
		if err != nil {
			return err
		}
		result := encLib.Base64Encode(input, urlSafe, noPadding)
		fmt.Fprintln(cmd.OutOrStdout(), result)

		return nil
	},
}

func init() {
	base64Cmd.Flags().BoolP("decode", "d", false, "decode base64 input")
	base64Cmd.Flags().BoolP("url-safe", "u", false, "use URL-safe encoding")
	base64Cmd.Flags().Bool("no-padding", false, "omit padding characters")
	Cmd.AddCommand(base64Cmd)
}
