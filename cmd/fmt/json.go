package fmt

import (
	goFmt "fmt"

	"github.com/spf13/cobra"

	fmtLib "github.com/agejevasv/swk/internal/fmt"
	"github.com/agejevasv/swk/internal/ioutil"
)

var jsonMinify bool
var jsonIndent int

var jsonCmd = &cobra.Command{
	Use:     "json [input]",
	Aliases: []string{"j"},
	Short:   "Format JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		opts := fmtLib.JSONOptions{
			Indent: jsonIndent,
			Minify: jsonMinify,
		}

		result, err := fmtLib.FormatJSON([]byte(input), opts)
		if err != nil {
			return err
		}

		goFmt.Fprintln(cmd.OutOrStdout(), string(result))
		return nil
	},
}

func init() {
	jsonCmd.Flags().BoolVarP(&jsonMinify, "minify", "m", false, "minify JSON")
	jsonCmd.Flags().IntVarP(&jsonIndent, "indent", "i", 2, "indentation spaces")

	Cmd.AddCommand(jsonCmd)
}
