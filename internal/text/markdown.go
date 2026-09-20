package text

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	gtext "github.com/yuin/goldmark/text"
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
	return []byte(stripMarkdown(input)), nil
}

// stripMarkdown renders markdown as plain text by walking the parsed document.
// Pattern matching cannot tell markup from content: it rewrote text inside code
// spans, ate literal asterisks in prose, and left fence characters behind.
func stripMarkdown(source []byte) string {
	md := goldmark.New(goldmark.WithExtensions(
		extension.GFM,
		extension.DefinitionList,
		extension.Footnote,
	))
	doc := md.Parser().Parse(gtext.NewReader(source))

	blocks := renderBlocks(doc, source, "")
	return strings.TrimSpace(strings.Join(blocks, "\n\n")) + "\n"
}

// renderBlocks turns each block-level child into one plain-text chunk.
func renderBlocks(parent ast.Node, src []byte, indent string) []string {
	var out []string

	for n := parent.FirstChild(); n != nil; n = n.NextSibling() {
		switch node := n.(type) {
		case *ast.FencedCodeBlock:
			out = appendBlock(out, indentLines(codeText(node.Lines(), src), indent))
		case *ast.CodeBlock:
			out = appendBlock(out, indentLines(codeText(node.Lines(), src), indent))
		case *ast.Blockquote:
			out = append(out, renderBlocks(node, src, indent)...)
		case *ast.List:
			out = appendBlock(out, renderList(node, src, indent))
		case *east.Table:
			out = appendBlock(out, renderTable(node, src))
		case *ast.ThematicBreak, *ast.HTMLBlock:
			// Markup with no text of its own.
		default:
			out = appendBlock(out, indentLines(inlineText(n, src), indent))
		}
	}

	return out
}

func appendBlock(out []string, block string) []string {
	if strings.TrimSpace(block) == "" {
		return out
	}
	return append(out, block)
}

// renderList keeps the markers, which carry the structure of the text.
func renderList(list *ast.List, src []byte, indent string) string {
	var lines []string
	number := list.Start
	if number == 0 {
		number = 1
	}

	for item := list.FirstChild(); item != nil; item = item.NextSibling() {
		marker := "- "
		if list.IsOrdered() {
			marker = fmt.Sprintf("%d. ", number)
			number++
		}

		body := strings.Join(renderBlocks(item, src, indent+"  "), "\n")
		body = strings.TrimSpace(body)
		if body == "" {
			continue
		}

		first, rest, _ := strings.Cut(body, "\n")
		lines = append(lines, indent+marker+first)
		if rest != "" {
			lines = append(lines, rest)
		}
	}

	return strings.Join(lines, "\n")
}

func renderTable(table *east.Table, src []byte) string {
	var rows []string

	for row := table.FirstChild(); row != nil; row = row.NextSibling() {
		var cells []string
		for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
			cells = append(cells, strings.TrimSpace(inlineText(cell, src)))
		}
		if len(cells) > 0 {
			rows = append(rows, strings.Join(cells, " | "))
		}
	}

	return strings.Join(rows, "\n")
}

func codeText(lines *gtext.Segments, src []byte) string {
	var b strings.Builder
	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)
		b.Write(seg.Value(src))
	}
	return strings.TrimRight(b.String(), "\n")
}

func indentLines(s, indent string) string {
	if indent == "" || s == "" {
		return s
	}
	parts := strings.Split(s, "\n")
	for i, p := range parts {
		parts[i] = indent + p
	}
	return strings.Join(parts, "\n")
}

// inlineText collects the text of an inline subtree, dropping the markup.
func inlineText(n ast.Node, src []byte) string {
	var b strings.Builder
	writeInline(&b, n, src)
	return strings.TrimRight(b.String(), " \t")
}

func writeInline(b *strings.Builder, n ast.Node, src []byte) {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		switch node := c.(type) {
		case *ast.Text:
			seg := node.Segment
			b.Write(seg.Value(src))
			if node.SoftLineBreak() || node.HardLineBreak() {
				b.WriteByte('\n')
			}
		case *ast.String:
			b.Write(node.Value)
		case *ast.AutoLink:
			b.Write(node.URL(src))
		case *ast.RawHTML:
			// Markup with no text of its own.
		default:
			// Code spans, emphasis, links and images all reduce to their contents.
			writeInline(b, c, src)
		}
	}
}
