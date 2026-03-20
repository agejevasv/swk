package text

import (
	"fmt"
	"html"
	"strconv"
	"strings"
)

func Escape(input string, mode string) (string, error) {
	switch mode {
	case "json":
		q := strconv.Quote(input)
		return q[1 : len(q)-1], nil
	case "xml":
		return escapeXML(input), nil
	case "html":
		return html.EscapeString(input), nil
	case "shell":
		return escapeShell(input), nil
	default:
		return "", fmt.Errorf("unsupported escape mode: %s (supported: json, xml, html, shell)", mode)
	}
}

func Unescape(input string, mode string) (string, error) {
	switch mode {
	case "json":
		s, err := strconv.Unquote(`"` + input + `"`)
		if err != nil {
			return "", fmt.Errorf("invalid escaped string: %w", err)
		}
		return s, nil
	case "xml":
		return unescapeXML(input), nil
	case "html":
		return html.UnescapeString(input), nil
	case "shell":
		return unescapeShell(input), nil
	default:
		return "", fmt.Errorf("unsupported unescape mode: %s (supported: json, xml, html, shell)", mode)
	}
}

func escapeXML(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&apos;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

var xmlUnescaper = strings.NewReplacer(
	"&amp;", "&",
	"&lt;", "<",
	"&gt;", ">",
	"&quot;", `"`,
	"&apos;", "'",
)

func unescapeXML(s string) string {
	return xmlUnescaper.Replace(s)
}

func escapeShell(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func unescapeShell(s string) string {
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		s = s[1 : len(s)-1]
	}
	return strings.ReplaceAll(s, `'\''`, "'")
}
