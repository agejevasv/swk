package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var csv2jsonCmd = &cobra.Command{
	Use:   "csv2json [input]",
	Short: "Convert CSV to JSON",
	Example: `  echo 'name,age\nalice,30' | swk convert csv2json
  swk convert csv2json --delimiter ';' data.csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		delimiter, err := ioutil.ParseDelimiter(ioutil.MustGetString(cmd, "delimiter"))
		if err != nil {
			return err
		}

		output, err := convLib.CSVToJSON([]byte(input), delimiter)
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), string(output))
		return nil
	},
}

func init() {
	csv2jsonCmd.Flags().StringP("delimiter", "d", ",", "CSV delimiter character")
	Cmd.AddCommand(csv2jsonCmd)
}
