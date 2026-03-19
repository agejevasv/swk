package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var (
	jyReverse bool
	jyIndent  int
)

var jsonYAMLCmd = &cobra.Command{
	Use:     "json-yaml [input]",
	Aliases: []string{"jy"},
	Short:   "Convert between JSON and YAML",
	Long:    "Converts JSON to YAML by default. Use --reverse to convert YAML to JSON.",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		var output []byte
		if jyReverse {
			output, err = convLib.YAMLToJSON([]byte(input), jyIndent)
		} else {
			output, err = convLib.JSONToYAML([]byte(input))
		}
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), string(output))
		return nil
	},
}

func init() {
	jsonYAMLCmd.Flags().BoolVarP(&jyReverse, "reverse", "r", false, "convert YAML to JSON instead")
	jsonYAMLCmd.Flags().IntVarP(&jyIndent, "indent", "i", 2, "indentation level for JSON output")
	Cmd.AddCommand(jsonYAMLCmd)
}
