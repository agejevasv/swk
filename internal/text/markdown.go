package text

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
body { max-width: 800px; margin: 40px auto; padding: 0 20px; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; line-height: 1.6; color: #222; }
h1, h2, h3, h4, h5, h6 { margin-top: 1.5em; margin-bottom: 0.5em; line-height: 1.3; }
h1 { font-size: 2em; border-bottom: 1px solid #ddd; padding-bottom: 0.3em; }
h2 { font-size: 1.5em; border-bottom: 1px solid #eee; padding-bottom: 0.3em; }
code { background: #f4f4f4; padding: 2px 6px; border-radius: 3px; font-size: 0.9em; }
pre { background: #f4f4f4; padding: 16px; border-radius: 6px; overflow-x: auto; }
pre code { background: none; padding: 0; }
blockquote { border-left: 4px solid #ddd; margin: 0; padding: 0 16px; color: #555; }
table { border-collapse: collapse; width: 100%%; }
th, td { border: 1px solid #ddd; padding: 8px 12px; text-align: left; }
th { background: #f4f4f4; }
a { color: #0366d6; text-decoration: none; }
a:hover { text-decoration: underline; }
img { max-width: 100%%; }
hr { border: none; border-top: 1px solid #ddd; margin: 2em 0; }
</style>
%s</head>
<body>
%s
</body>
</html>`

const highlightSnippet = `<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/{{theme}}.min.css">
<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
<script>hljs.highlightAll();</script>
`

// Precompiled regexes for markdown stripping.
var (
	reHeading = regexp.MustCompile(`(?m)^#{1,6}\s+`)
	reBold    = regexp.MustCompile(`\*\*(.+?)\*\*`)
	reBold2   = regexp.MustCompile(`__(.+?)__`)
	reItalic  = regexp.MustCompile(`\*(.+?)\*`)
	// Underscore emphasis only applies at word boundaries; CommonMark treats
	// intra-word underscores as literal, so snake_case_names survive.
	reItalic2    = regexp.MustCompile(`(^|[^\pL\pN_])_([^_\n]+)_($|[^\pL\pN_])`)
	reStrike     = regexp.MustCompile(`~~(.+?)~~`)
	reCode       = regexp.MustCompile("`([^`]+)`")
	reCodeBlock  = regexp.MustCompile("(?s)```[a-z]*\n?(.*?)```")
	reLink       = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	reImg        = regexp.MustCompile(`!\[([^\]]*)\]\([^)]+\)`)
	reHR         = regexp.MustCompile(`(?m)^[-*_]{3,}\s*$`)
	reBlockquote = regexp.MustCompile(`(?m)^>\s?`)
	reBlankLines = regexp.MustCompile(`\n{3,}`)
)

// themePattern limits --theme to characters that are safe in a URL path
// segment, so the value cannot break out of the stylesheet link.
var themePattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

func RenderMarkdown(input []byte, toHTML bool, syntaxHighlight bool, theme string) ([]byte, error) {
	if toHTML {
		var buf bytes.Buffer
		md := goldmark.New(goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			extension.Typographer,
		))
		if err := md.Convert(input, &buf); err != nil {
			return nil, err
		}

		highlight := ""
		if syntaxHighlight {
			if theme == "" {
				theme = "github"
			}
			if !themePattern.MatchString(theme) {
				return nil, fmt.Errorf("invalid theme %q: use letters, digits, dot, dash or underscore", theme)
			}
			highlight = strings.ReplaceAll(highlightSnippet, "{{theme}}", theme)
		}
		page := fmt.Sprintf(htmlTemplate, highlight, buf.String())
		return []byte(page + "\n"), nil
	}
	return []byte(stripMarkdown(string(input))), nil
}

// stripUnderscoreEmphasis removes _emphasis_ at word boundaries. The pattern
// consumes the boundary character on each side, so adjacent spans such as
// "_a_ _b_" need more than one pass; it runs until the text stops changing.
func stripUnderscoreEmphasis(s string) string {
	for i := 0; i < 8; i++ {
		next := reItalic2.ReplaceAllString(s, "$1$2$3")
		if next == s {
			break
		}
		s = next
	}
	return s
}

func stripMarkdown(s string) string {
	s = reHeading.ReplaceAllString(s, "")
	s = reBold.ReplaceAllString(s, "$1")
	s = reBold2.ReplaceAllString(s, "$1")
	s = reItalic.ReplaceAllString(s, "$1")
	s = stripUnderscoreEmphasis(s)
	s = reStrike.ReplaceAllString(s, "$1")
	s = reCode.ReplaceAllString(s, "$1")
	s = reCodeBlock.ReplaceAllString(s, "$1")
	// Images first: the link pattern also matches the "[alt](src)" tail of an
	// image and would leave a stray "!" behind.
	s = reImg.ReplaceAllString(s, "$1")
	s = reLink.ReplaceAllString(s, "$1")
	s = reHR.ReplaceAllString(s, "")
	s = reBlockquote.ReplaceAllString(s, "")
	s = reBlankLines.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s) + "\n"
}
