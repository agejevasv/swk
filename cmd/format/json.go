package format

import (
	"fmt"

	"github.com/spf13/cobra"

	fmtLib "github.com/agejevasv/swk/internal/fmt"
	"github.com/agejevasv/swk/internal/ioutil"
)

var jsonCmd = &cobra.Command{
	Use:   "json [input]",
	Short: "Prettify or minify JSON",
	Example: `  # Prettify JSON
  echo '{"a":1}' | swk format json

  # Minify JSON
  echo '{"a": 1}' | swk format json --minify

  # Custom indent
  echo '{"a":1}' | swk format json --indent 4`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		opts := fmtLib.JSONOptions{
			Indent: ioutil.MustGetInt(cmd, "indent"),
			Minify: ioutil.MustGetBool(cmd, "minify"),
		}
		result, err := fmtLib.FormatJSON([]byte(input), opts)
		if err != nil {
			return err
		}
		fmt.Fprint(cmd.OutOrStdout(), string(result))
		return nil
	},
}

func init() {
	jsonCmd.Flags().BoolP("minify", "m", false, "minify JSON output")
	jsonCmd.Flags().IntP("indent", "i", 2, "indentation spaces")
	Cmd.AddCommand(jsonCmd)
}
