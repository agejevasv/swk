package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var (
	tableStyle     string
	tableFrom      string
	tableDelimiter string
)

var tableCmd = &cobra.Command{
	Use:   "table [input]",
	Short: "Render JSON or CSV as a formatted table",
	Example: `  # JSON array to box table
  echo '[{"name":"alice","age":30},{"name":"bob","age":25}]' | swk convert table

  # Simple ASCII style
  echo '[{"name":"alice"}]' | swk convert table --style simple

  # Plain (no borders)
  echo '[{"name":"alice"}]' | swk convert table --style plain

  # CSV input
  echo 'name,age\nalice,30' | swk convert table --from csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		delimiter := ','
		if len(tableDelimiter) > 0 {
			delimiter = rune(tableDelimiter[0])
		}

		result, err := convLib.ToTable([]byte(input), tableStyle, tableFrom, delimiter)
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	tableCmd.Flags().StringVar(&tableStyle, "style", "box", "table style (box, simple, plain)")
	tableCmd.Flags().StringVar(&tableFrom, "from", "json", "input format (json, csv)")
	tableCmd.Flags().StringVarP(&tableDelimiter, "delimiter", "d", ",", "CSV delimiter character")
	Cmd.AddCommand(tableCmd)
}
