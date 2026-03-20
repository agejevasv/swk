package generate

import (
	"fmt"

	"github.com/spf13/cobra"

	genLib "github.com/agejevasv/swk/internal/gen"
	"github.com/agejevasv/swk/internal/ioutil"
)

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Generate lorem ipsum text",
	Long:  "Generate lorem ipsum placeholder text as words, sentences, or paragraphs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		textWords := ioutil.MustGetInt(cmd, "words")
		textSentences := ioutil.MustGetInt(cmd, "sentences")
		textParagraphs := ioutil.MustGetInt(cmd, "paragraphs")

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
	textCmd.Flags().IntP("words", "w", 0, "number of words to generate")
	textCmd.Flags().IntP("sentences", "s", 0, "number of sentences to generate")
	textCmd.Flags().IntP("paragraphs", "p", 0, "number of paragraphs to generate")
	Cmd.AddCommand(textCmd)
}
