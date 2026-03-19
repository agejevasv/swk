package gen

import (
	"fmt"

	"github.com/spf13/cobra"

	genLib "github.com/agejevasv/swk/internal/gen"
)

var (
	loremWords      int
	loremSentences  int
	loremParagraphs int
)

var loremCmd = &cobra.Command{
	Use:     "lorem",
	Aliases: []string{"li"},
	Short:   "Generate lorem ipsum text",
	Long:    "Generate lorem ipsum placeholder text as words, sentences, or paragraphs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if loremWords == 0 && loremSentences == 0 && loremParagraphs == 0 {
			loremParagraphs = 1
		}

		if loremWords > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), genLib.GenerateWords(loremWords))
		}
		if loremSentences > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), genLib.GenerateSentences(loremSentences))
		}
		if loremParagraphs > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), genLib.GenerateParagraphs(loremParagraphs))
		}

		return nil
	},
}

func init() {
	loremCmd.Flags().IntVarP(&loremWords, "words", "w", 0, "number of words to generate")
	loremCmd.Flags().IntVarP(&loremSentences, "sentences", "s", 0, "number of sentences to generate")
	loremCmd.Flags().IntVarP(&loremParagraphs, "paragraphs", "p", 0, "number of paragraphs to generate")
	Cmd.AddCommand(loremCmd)
}
