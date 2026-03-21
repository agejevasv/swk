package diff

import (
	"fmt"

	"github.com/spf13/cobra"

	diffLib "github.com/agejevasv/swk/internal/diff"
	"github.com/agejevasv/swk/internal/ioutil"
)

var textCmd = &cobra.Command{
	Use:     "text <file1> <file2>",
	Aliases: []string{"txt"},
	Short:   "Unified text diff",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, b, err := diffLib.ReadTwoInputs(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		ctx := ioutil.MustGetInt(cmd, "context")
		result := diffLib.DiffText(a, b, ctx)
		if result != "" {
			if shouldColorize(cmd) {
				result = diffLib.Colorize(result)
			}
			fmt.Fprint(cmd.OutOrStdout(), result)
		}
		return nil
	},
}

func init() {
	textCmd.Flags().IntP("context", "C", 3, "context lines around changes")
	textCmd.Flags().String("color", "auto", "color output: auto, always, never")
	Cmd.AddCommand(textCmd)
}
