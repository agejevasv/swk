package query

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	queryLib "github.com/agejevasv/swk/internal/query"
)

var regexCmd = &cobra.Command{
	Use:     "regex PATTERN [input]",
	Aliases: []string{"re"},
	Short:   "Test regular expressions against input",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		regexPattern := args[0]
		input, err := ioutil.ReadFileInputString(args[1:], cmd.InOrStdin())
		if err != nil {
			return err
		}

		regexGlobal := ioutil.MustGetBool(cmd, "global")
		regexGroups := ioutil.MustGetBool(cmd, "groups")
		regexReplace := ioutil.MustGetString(cmd, "replace")

		if regexReplace != "" {
			result, err := queryLib.RegexReplace(input, regexPattern, regexReplace)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), result)
			return nil
		}

		result, err := queryLib.RegexTest(input, regexPattern, regexGlobal)
		if err != nil {
			return err
		}

		if !result.Matched {
			return ioutil.NoMatchError{}
		}

		if regexGroups {
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(result.Matches)
		}

		for _, m := range result.Matches {
			fmt.Fprintln(cmd.OutOrStdout(), m.Value)
		}

		return nil
	},
}

func init() {
	regexCmd.Flags().BoolP("global", "g", false, "Find all matches")
	regexCmd.Flags().Bool("groups", false, "Show capture groups as JSON")
	regexCmd.Flags().StringP("replace", "r", "", "Replacement string")
	Cmd.AddCommand(regexCmd)
}
