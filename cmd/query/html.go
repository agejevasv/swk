package query

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
)

var (
	htmlQuerySelector string
	htmlQueryAttr     string
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
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(input))
		if err != nil {
			return fmt.Errorf("parsing HTML: %w", err)
		}

		doc.Find(htmlQuerySelector).Each(func(i int, s *goquery.Selection) {
			if htmlQueryAttr != "" {
				if val, exists := s.Attr(htmlQueryAttr); exists {
					fmt.Fprintln(cmd.OutOrStdout(), val)
				}
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), strings.TrimSpace(s.Text()))
			}
		})

		return nil
	},
}

func init() {
	htmlCmd.Flags().StringVarP(&htmlQuerySelector, "query", "q", "", "CSS selector")
	htmlCmd.Flags().StringVar(&htmlQueryAttr, "attr", "", "extract attribute value instead of text")
	htmlCmd.MarkFlagRequired("query")
	Cmd.AddCommand(htmlCmd)
}
