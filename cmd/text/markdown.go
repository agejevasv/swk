package text

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	textLib "github.com/agejevasv/swk/internal/text"
)

var mdHTML bool

var markdownCmd = &cobra.Command{
	Use:     "markdown [text]",
	Aliases: []string{"md"},
	Short:   "Render markdown to HTML or plain text",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		result, err := textLib.RenderMarkdown([]byte(input), mdHTML)
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), string(result))
		return nil
	},
}

func init() {
	markdownCmd.Flags().BoolVar(&mdHTML, "html", false, "Output HTML (default is plain text)")
	Cmd.AddCommand(markdownCmd)
}
