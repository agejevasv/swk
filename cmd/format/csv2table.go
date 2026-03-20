package format

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var csv2tableCmd = &cobra.Command{
	Use:   "csv2table [input]",
	Short: "Format CSV as a table",
	Example: `  echo 'name,age\nalice,30' | swk format csv2table
  swk format csv2table --delimiter ';' data.csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		csvDelimiter := ioutil.MustGetString(cmd, "delimiter")
		delimiter := ','
		if len(csvDelimiter) > 0 {
			delimiter = rune(csvDelimiter[0])
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
