package query

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	queryLib "github.com/agejevasv/swk/internal/query"
)

var htmlCmd = &cobra.Command{
	Use:   "html [input]",
	Short: "Query HTML with CSS selectors",
	Example: `  # Extract all links
  curl -s https://example.com | swk query html -q 'a' --attr href

  # Get text content of all paragraphs
  cat page.html | swk query html -q 'p'

  # Extract specific element
  echo '<div class="x"><span>hi</span></div>' | swk query html -q 'div.x span'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		htmlQuerySelector := ioutil.MustGetString(cmd, "query")
		htmlQueryAttr := ioutil.MustGetString(cmd, "attr")

		results, err := queryLib.HTMLQuery(input, htmlQuerySelector, htmlQueryAttr)
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
	htmlCmd.Flags().StringP("query", "q", "", "CSS selector")
	htmlCmd.Flags().String("attr", "", "extract attribute value instead of text")
	htmlCmd.MarkFlagRequired("query")
	Cmd.AddCommand(htmlCmd)
}
