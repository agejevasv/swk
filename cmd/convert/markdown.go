package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	textLib "github.com/agejevasv/swk/internal/text"
)

var markdownCmd = &cobra.Command{
	Use:     "markdown [text]",
	Aliases: []string{"md"},
	Short:   "Convert markdown to HTML or plain text",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		mdHTML := ioutil.MustGetBool(cmd, "html")
		syntaxHL := ioutil.MustGetBool(cmd, "syntax-highlight")
		theme := ioutil.MustGetString(cmd, "theme")

		result, err := textLib.RenderMarkdown([]byte(input), mdHTML, syntaxHL, theme)
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), string(result))
		return nil
	},
}

func init() {
	markdownCmd.Flags().Bool("html", false, "Output HTML (default is plain text)")
	markdownCmd.Flags().Bool("syntax-highlight", false, "Include highlight.js for syntax highlighting (requires --html)")
	markdownCmd.Flags().String("theme", "github", "highlight.js theme (github, monokai, dracula, nord, tokyo-night-dark, etc.)")
	Cmd.AddCommand(markdownCmd)
}
