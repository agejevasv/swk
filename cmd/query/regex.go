package query

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	queryLib "github.com/agejevasv/swk/internal/query"
)

var regexCmd = &cobra.Command{
	Use:     "regex PATTERN [input]",
	Aliases: []string{"re"},
	Short:   "Match, extract, or replace with regular expressions",
	Long: `By default, prints lines that match the pattern (like grep).
Use -o to print only the matched parts, --groups for structured JSON output,
or --replace for substitution.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		regexPattern := args[0]
		input, err := ioutil.ReadFileInputString(args[1:], cmd.InOrStdin())
		if err != nil {
			return err
		}

		regexGlobal := ioutil.MustGetBool(cmd, "global")
		regexGroups := ioutil.MustGetBool(cmd, "groups")
		regexReplace := ioutil.MustGetString(cmd, "replace")
		regexOnly := ioutil.MustGetBool(cmd, "only-matching")

		// --replace: substitute on full input.
		if regexReplace != "" {
			result, err := queryLib.RegexReplace(input, regexPattern, regexReplace)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), result)
			return nil
		}

		// --groups: structured JSON output on full input.
		if regexGroups {
			result, err := queryLib.RegexTest(input, regexPattern, regexGlobal)
			if err != nil {
				return err
			}
			if !result.Matched {
				return ioutil.NoMatchError{}
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(result.Matches)
		}

		// -o: print only matched parts from full input.
		if regexOnly {
			result, err := queryLib.RegexTest(input, regexPattern, regexGlobal)
			if err != nil {
				return err
			}
			if !result.Matched {
				return ioutil.NoMatchError{}
			}
			for _, m := range result.Matches {
				fmt.Fprintln(cmd.OutOrStdout(), m.Value)
			}
			return nil
		}

		// Default: print matching lines (grep-like).
		re, err := regexp.Compile(regexPattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
		matched := false
		for _, line := range strings.Split(input, "\n") {
			if re.MatchString(line) {
				fmt.Fprintln(cmd.OutOrStdout(), line)
				matched = true
			}
		}
		if !matched {
			return ioutil.NoMatchError{}
		}
		return nil
	},
}

func init() {
	regexCmd.Flags().BoolP("global", "g", false, "find all matches (for -o and --groups)")
	regexCmd.Flags().BoolP("only-matching", "o", false, "print only matched parts, not full lines")
	regexCmd.Flags().Bool("groups", false, "show capture groups as JSON")
	regexCmd.Flags().StringP("replace", "r", "", "replacement string")
	Cmd.AddCommand(regexCmd)
}
