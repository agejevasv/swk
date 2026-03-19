package generate

import (
	"fmt"

	"github.com/spf13/cobra"

	genLib "github.com/agejevasv/swk/internal/gen"
)

var (
	textWords      int
	textSentences  int
	textParagraphs int
)

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Generate lorem ipsum text",
	Long:  "Generate lorem ipsum placeholder text as words, sentences, or paragraphs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if textWords == 0 && textSentences == 0 && textParagraphs == 0 {
			textParagraphs = 1
		}

		if textWords > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), genLib.GenerateWords(textWords))
		}
		if textSentences > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), genLib.GenerateSentences(textSentences))
		}
		if textParagraphs > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), genLib.GenerateParagraphs(textParagraphs))
		}

		return nil
	},
}

func init() {
	textCmd.Flags().IntVarP(&textWords, "words", "w", 0, "number of words to generate")
	textCmd.Flags().IntVarP(&textSentences, "sentences", "s", 0, "number of sentences to generate")
	textCmd.Flags().IntVarP(&textParagraphs, "paragraphs", "p", 0, "number of paragraphs to generate")
	Cmd.AddCommand(textCmd)
}
