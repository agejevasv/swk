package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	textLib "github.com/agejevasv/swk/internal/text"
)

var (
	mdHTML  bool
	mdTheme string
)

var markdownCmd = &cobra.Command{
	Use:     "markdown [text]",
	Aliases: []string{"md"},
	Short:   "Render markdown to HTML or plain text",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		result, err := textLib.RenderMarkdown([]byte(input), mdHTML, mdTheme)
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), string(result))
		return nil
	},
}

func init() {
	markdownCmd.Flags().BoolVar(&mdHTML, "html", false, "Output HTML (default is plain text)")
	markdownCmd.Flags().StringVar(&mdTheme, "theme", "github", "highlight.js theme (github, monokai, dracula, nord, tokyo-night-dark, etc.)")
	Cmd.AddCommand(markdownCmd)
}
