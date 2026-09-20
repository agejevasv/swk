package format

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var csv2tableCmd = &cobra.Command{
	Use:   "csv2table [input]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Format CSV as a table",
	Example: `  printf 'name,age\nalice,30\n' | swk format csv2table
  swk format csv2table --delimiter ';' data.csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		delimiter, err := ioutil.ParseDelimiter(ioutil.MustGetString(cmd, "delimiter"))
		if err != nil {
			return err
		}

		result, err := convLib.ToTable([]byte(input), ioutil.MustGetString(cmd, "style"), "csv", delimiter)
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	csv2tableCmd.Flags().String("style", "box", "table style (box, simple, plain)")
	csv2tableCmd.Flags().StringP("delimiter", "d", ",", "CSV delimiter character")
	Cmd.AddCommand(csv2tableCmd)
}
