package query

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HTMLQuery queries an HTML document using a CSS selector.
// If attr is non-empty, it extracts that attribute from matched elements.
// Otherwise it extracts text content.
func HTMLQuery(input string, selector string, attr string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	var results []string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if attr != "" {
			if val, exists := s.Attr(attr); exists {
				results = append(results, val)
			}
		} else {
			results = append(results, strings.TrimSpace(s.Text()))
		}
	})

	return results, nil
}
