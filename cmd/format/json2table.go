package format

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var json2tableCmd = &cobra.Command{
	Use:   "json2table [input]",
	Short: "Format JSON array as a table",
	Example: `  echo '[{"name":"alice","age":30}]' | swk format json2table
  swk format json2table --style simple data.json
  swk format json2table --style plain data.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		result, err := convLib.ToTable([]byte(input), ioutil.MustGetString(cmd, "style"), "json", ',')
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	json2tableCmd.Flags().String("style", "box", "table style (box, simple, plain)")
	Cmd.AddCommand(json2tableCmd)
}
