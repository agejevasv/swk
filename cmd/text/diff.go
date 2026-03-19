package text

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	textLib "github.com/agejevasv/swk/internal/text"
)

var (
	diffFile1   string
	diffFile2   string
	diffContext int
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show unified diff between two files",
	RunE: func(cmd *cobra.Command, args []string) error {
		aBytes, err := os.ReadFile(diffFile1)
		if err != nil {
			return fmt.Errorf("reading file1: %w", err)
		}
		bBytes, err := os.ReadFile(diffFile2)
		if err != nil {
			return fmt.Errorf("reading file2: %w", err)
		}

		result := textLib.Diff(string(aBytes), string(bBytes), diffContext)
		fmt.Fprint(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	diffCmd.Flags().StringVarP(&diffFile1, "file1", "a", "", "First file")
	diffCmd.Flags().StringVarP(&diffFile2, "file2", "b", "", "Second file")
	diffCmd.Flags().IntVarP(&diffContext, "context", "C", 3, "Number of context lines")
	diffCmd.MarkFlagRequired("file1")
	diffCmd.MarkFlagRequired("file2")
	Cmd.AddCommand(diffCmd)
}
