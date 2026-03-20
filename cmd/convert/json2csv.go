package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var json2csvCmd = &cobra.Command{
	Use:   "json2csv [input]",
	Short: "Convert JSON to CSV",
	Example: `  echo '[{"name":"alice","age":30}]' | swk convert json2csv
  swk convert json2csv --delimiter ';' data.json`,
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

		output, err := convLib.JSONToCSV([]byte(input), delimiter)
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), string(output))
		return nil
	},
}

func init() {
	json2csvCmd.Flags().StringP("delimiter", "d", ",", "CSV delimiter character")
	Cmd.AddCommand(json2csvCmd)
}
