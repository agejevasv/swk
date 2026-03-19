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
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/{{theme}}.min.css">
<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
<script>hljs.highlightAll();</script>
</head>
<body>
%s
</body>
</html>`

func RenderMarkdown(input []byte, toHTML bool, theme string) ([]byte, error) {
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
		if theme == "" {
			theme = "github"
		}
		tmpl := strings.ReplaceAll(htmlTemplate, "{{theme}}", theme)
		page := fmt.Sprintf(tmpl, buf.String())
		return []byte(page), nil
	}
	return []byte(stripMarkdown(string(input))), nil
}

func stripMarkdown(s string) string {
	re := regexp.MustCompile(`(?m)^#{1,6}\s+`)
	s = re.ReplaceAllString(s, "")

	reBold := regexp.MustCompile(`\*\*(.+?)\*\*`)
	s = reBold.ReplaceAllString(s, "$1")
	reBold2 := regexp.MustCompile(`__(.+?)__`)
	s = reBold2.ReplaceAllString(s, "$1")

	reItalic := regexp.MustCompile(`\*(.+?)\*`)
	s = reItalic.ReplaceAllString(s, "$1")
	reItalic2 := regexp.MustCompile(`_(.+?)_`)
	s = reItalic2.ReplaceAllString(s, "$1")

	reStrike := regexp.MustCompile(`~~(.+?)~~`)
	s = reStrike.ReplaceAllString(s, "$1")

	reCode := regexp.MustCompile("`([^`]+)`")
	s = reCode.ReplaceAllString(s, "$1")

	reCodeBlock := regexp.MustCompile("(?s)```[a-z]*\n?(.*?)```")
	s = reCodeBlock.ReplaceAllString(s, "$1")

	reLink := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	s = reLink.ReplaceAllString(s, "$1")

	reImg := regexp.MustCompile(`!\[([^\]]*)\]\([^)]+\)`)
	s = reImg.ReplaceAllString(s, "$1")

	reHR := regexp.MustCompile(`(?m)^[-*_]{3,}\s*$`)
	s = reHR.ReplaceAllString(s, "")

	reBlockquote := regexp.MustCompile(`(?m)^>\s?`)
	s = reBlockquote.ReplaceAllString(s, "")

	reBlankLines := regexp.MustCompile(`\n{3,}`)
	s = reBlankLines.ReplaceAllString(s, "\n\n")

	return strings.TrimSpace(s) + "\n"
}
