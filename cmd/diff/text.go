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
	Long:    "Print a unified diff of two files. Exits 1 when they differ, like diff(1).",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, b, err := diffLib.ReadTwoInputs(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		ctx := ioutil.MustGetInt(cmd, "context")
		result := diffLib.DiffText(a, b, ctx)
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
	textCmd.Flags().IntP("context", "C", 3, "context lines around changes")
	textCmd.Flags().String("color", "auto", "color output: auto, always, never")
	Cmd.AddCommand(textCmd)
}
