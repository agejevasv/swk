package text

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
)

func RenderMarkdown(input []byte, toHTML bool) ([]byte, error) {
	if toHTML {
		var buf bytes.Buffer
		md := goldmark.New()
		if err := md.Convert(input, &buf); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
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
