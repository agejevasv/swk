package test

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	testLib "github.com/agejevasv/swk/internal/test"
)

var jsonpathQuery string

var jsonpathCmd = &cobra.Command{
	Use:     "jsonpath [json]",
	Aliases: []string{"jp"},
	Short:   "Query JSON with JSONPath expressions",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		result, err := testLib.JSONPathQuery([]byte(input), jsonpathQuery)
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), string(result))
		return nil
	},
}

func init() {
	jsonpathCmd.Flags().StringVarP(&jsonpathQuery, "query", "q", "", "JSONPath query expression")
	jsonpathCmd.MarkFlagRequired("query")
	Cmd.AddCommand(jsonpathCmd)
}
