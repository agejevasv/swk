package query

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	queryLib "github.com/agejevasv/swk/internal/query"
)

var jsonCmd = &cobra.Command{
	Use:   "json EXPRESSION [input]",
	Short: "Query JSON with JSONPath expressions",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonpathQuery := args[0]
		input, err := ioutil.ReadFileInputString(args[1:], cmd.InOrStdin())
		if err != nil {
			return err
		}

		result, err := queryLib.JSONPathQuery([]byte(input), jsonpathQuery)
		if err != nil {
			return err
		}

		if result == nil {
			return ioutil.NoMatchError{}
		}

		fmt.Fprint(cmd.OutOrStdout(), string(result))
		return nil
	},
}

func init() {
	Cmd.AddCommand(jsonCmd)
}
