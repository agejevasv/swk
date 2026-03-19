package fmt

import (
	"github.com/spf13/cobra"

	fmtLib "github.com/agejevasv/swk/internal/fmt"
	"github.com/agejevasv/swk/internal/ioutil"
)

var xmlMinify bool
var xmlIndent int

var xmlCmd = &cobra.Command{
	Use:     "xml [input]",
	Aliases: []string{"x"},
	Short:   "Format XML",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		opts := fmtLib.XMLOptions{
			Indent: xmlIndent,
			Minify: xmlMinify,
		}

		result, err := fmtLib.FormatXML([]byte(input), opts)
		if err != nil {
			return err
		}

		if len(result) > 0 && result[len(result)-1] != '\n' {
			result = append(result, '\n')
		}
		_, err = cmd.OutOrStdout().Write(result)
		return err
	},
}

func init() {
	xmlCmd.Flags().BoolVarP(&xmlMinify, "minify", "m", false, "minify XML")
	xmlCmd.Flags().IntVarP(&xmlIndent, "indent", "i", 2, "indentation spaces")

	Cmd.AddCommand(xmlCmd)
}
