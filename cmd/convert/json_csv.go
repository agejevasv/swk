package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var (
	jcReverse   bool
	jcDelimiter string
)

var jsonCSVCmd = &cobra.Command{
	Use:     "json-csv [input]",
	Aliases: []string{"jc"},
	Short:   "Convert between JSON and CSV",
	Long:    "Converts JSON array of objects to CSV by default. Use --reverse to convert CSV to JSON.",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		delimiter := ','
		if len(jcDelimiter) > 0 {
			delimiter = rune(jcDelimiter[0])
		}

		var output []byte
		if jcReverse {
			output, err = convLib.CSVToJSON([]byte(input), delimiter)
		} else {
			output, err = convLib.JSONToCSV([]byte(input), delimiter)
		}
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), string(output))
		return nil
	},
}

func init() {
	jsonCSVCmd.Flags().BoolVarP(&jcReverse, "reverse", "r", false, "convert CSV to JSON instead")
	jsonCSVCmd.Flags().StringVarP(&jcDelimiter, "delimiter", "d", ",", "CSV delimiter character")
	Cmd.AddCommand(jsonCSVCmd)
}
