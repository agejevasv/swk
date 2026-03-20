package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var json2yamlCmd = &cobra.Command{
	Use:   "json2yaml [input]",
	Short: "Convert JSON to YAML",
	Example: `  echo '{"a":1}' | swk convert json2yaml
  swk convert json2yaml data.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		output, err := convLib.JSONToYAML([]byte(input))
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), string(output))
		return nil
	},
}

func init() {
	Cmd.AddCommand(json2yamlCmd)
}
