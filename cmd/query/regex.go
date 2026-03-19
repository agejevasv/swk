package query

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	testLib "github.com/agejevasv/swk/internal/test"
)

var (
	regexPattern string
	regexGlobal  bool
	regexGroups  bool
	regexReplace string
)

var regexCmd = &cobra.Command{
	Use:     "regex [input]",
	Aliases: []string{"re"},
	Short:   "Test regular expressions against input",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		if regexReplace != "" {
			result, err := testLib.RegexReplace(input, regexPattern, regexReplace)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), result)
			return nil
		}

		result, err := testLib.RegexTest(input, regexPattern, regexGlobal)
		if err != nil {
			return err
		}

		if !result.Matched {
			fmt.Fprintln(cmd.OutOrStdout(), "No matches found.")
			return nil
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
	regexCmd.Flags().StringVarP(&regexPattern, "pattern", "p", "", "Regex pattern")
	regexCmd.Flags().BoolVarP(&regexGlobal, "global", "g", false, "Find all matches")
	regexCmd.Flags().BoolVar(&regexGroups, "groups", false, "Show capture groups as JSON")
	regexCmd.Flags().StringVarP(&regexReplace, "replace", "r", "", "Replacement string")
	regexCmd.MarkFlagRequired("pattern")
	Cmd.AddCommand(regexCmd)
}
