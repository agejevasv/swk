package query

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	queryLib "github.com/agejevasv/swk/internal/query"
)

var regexCmd = &cobra.Command{
	Use:     "regex [input]",
	Aliases: []string{"re"},
	Short:   "Test regular expressions against input",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		regexPattern := ioutil.MustGetString(cmd, "pattern")
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
	regexCmd.Flags().StringP("pattern", "p", "", "Regex pattern")
	regexCmd.Flags().BoolP("global", "g", false, "Find all matches")
	regexCmd.Flags().Bool("groups", false, "Show capture groups as JSON")
	regexCmd.Flags().StringP("replace", "r", "", "Replacement string")
	regexCmd.MarkFlagRequired("pattern")
	Cmd.AddCommand(regexCmd)
}
