package diff

import (
	"fmt"

	"github.com/spf13/cobra"

	diffLib "github.com/agejevasv/swk/internal/diff"
	"github.com/agejevasv/swk/internal/ioutil"
)

var jsonCmd = &cobra.Command{
	Use:   "json <file1> <file2>",
	Short: "Semantic JSON diff (normalizes key order)",
	Long:  "Compare two JSON documents ignoring key order. Exits 1 when they differ, like diff(1).",
	Example: `  swk diff json old.json new.json
  curl -s api/v1 | swk diff json - saved.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		a, b, err := diffLib.ReadTwoInputs(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		ctx := ioutil.MustGetInt(cmd, "context")
		result, err := diffLib.DiffJSON(a, b, ctx)
		if err != nil {
			return err
		}
		if result == "" {
			return nil
		}

		if shouldColorize(cmd) {
			result = diffLib.Colorize(result)
		}
		fmt.Fprint(cmd.OutOrStdout(), result)

		// Exit 1 when the inputs differ, the way diff(1) does.
		return ioutil.DifferError{}
	},
}

func init() {
	jsonCmd.Flags().IntP("context", "C", 3, "context lines around changes")
	jsonCmd.Flags().String("color", "auto", "color output: auto, always, never")
	Cmd.AddCommand(jsonCmd)
}
