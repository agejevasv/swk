package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var yaml2jsonCmd = &cobra.Command{
	Use:   "yaml2json [input]",
	Short: "Convert YAML to JSON",
	Example: `  echo 'a: 1' | swk convert yaml2json
  swk convert yaml2json config.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		output, err := convLib.YAMLToJSON([]byte(input), ioutil.MustGetInt(cmd, "indent"))
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), string(output))
		return nil
	},
}

func init() {
	yaml2jsonCmd.Flags().IntP("indent", "i", 2, "indentation spaces")
	Cmd.AddCommand(yaml2jsonCmd)
}
