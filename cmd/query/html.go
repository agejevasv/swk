package query

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	queryLib "github.com/agejevasv/swk/internal/query"
)

var htmlCmd = &cobra.Command{
	Use:   "html SELECTOR [input]",
	Short: "Query HTML with CSS selectors",
	Example: `  # Extract all links
  curl -s https://example.com | swk query html 'a' --attr href

  # Get text content of all paragraphs
  cat page.html | swk query html 'p'

  # Extract specific element
  echo '<div class="x"><span>hi</span></div>' | swk query html 'div.x span'`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]
		input, err := ioutil.ReadFileInputString(args[1:], cmd.InOrStdin())
		if err != nil {
			return err
		}

		htmlQueryAttr := ioutil.MustGetString(cmd, "attr")

		results, err := queryLib.HTMLQuery(input, selector, htmlQueryAttr)
		if err != nil {
			return err
		}

		if len(results) == 0 {
			return ioutil.NoMatchError{}
		}

		for _, r := range results {
			fmt.Fprintln(cmd.OutOrStdout(), r)
		}

		return nil
	},
}

func init() {
	htmlCmd.Flags().String("attr", "", "extract attribute value instead of text")
	Cmd.AddCommand(htmlCmd)
}
