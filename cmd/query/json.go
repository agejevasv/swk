package query

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	queryLib "github.com/agejevasv/swk/internal/query"
)

var jsonCmd = &cobra.Command{
	Use:     "json [input]",
	Aliases: []string{"jp"},
	Short:   "Query JSON with JSONPath expressions",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		jsonpathQuery := ioutil.MustGetString(cmd, "query")

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
	jsonCmd.Flags().StringP("query", "q", "", "JSONPath query expression")
	jsonCmd.MarkFlagRequired("query")
	Cmd.AddCommand(jsonCmd)
}
