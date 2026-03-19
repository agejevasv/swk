package text

import (
	"fmt"
	"strings"
	"unicode"
)

func ConvertCase(input string, to string) (string, error) {
	words := splitWords(input)

	switch to {
	case "upper":
		return strings.ToUpper(input), nil
	case "lower":
		return strings.ToLower(input), nil
	case "camel":
		return toCamel(words), nil
	case "pascal":
		return toPascal(words), nil
	case "snake":
		return joinLower(words, "_"), nil
	case "kebab":
		return joinLower(words, "-"), nil
	case "title":
		return toTitle(words), nil
	case "sentence":
		return toSentence(words), nil
	case "dot":
		return joinLower(words, "."), nil
	case "path":
		return joinLower(words, "/"), nil
	default:
		return "", fmt.Errorf("unsupported case: %s (supported: camel, pascal, snake, kebab, upper, lower, title, sentence, dot, path)", to)
	}
}

func splitWords(s string) []string {
	var words []string
	var current strings.Builder

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == ' ' || r == '_' || r == '-' || r == '.' || r == '/' || r == '\t' || r == '\n' {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
			continue
		}

		if i > 0 && current.Len() > 0 {
			prev := runes[i-1]
			if unicode.IsLower(prev) && unicode.IsUpper(r) {
				words = append(words, current.String())
				current.Reset()
			}
			if unicode.IsUpper(prev) && unicode.IsUpper(r) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				words = append(words, current.String())
				current.Reset()
			}
		}

		current.WriteRune(r)
	}

	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}

func toCamel(words []string) string {
	if len(words) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(strings.ToLower(words[0]))
	for _, w := range words[1:] {
		b.WriteString(capitalize(w))
	}
	return b.String()
}

func toPascal(words []string) string {
	var b strings.Builder
	for _, w := range words {
		b.WriteString(capitalize(w))
	}
	return b.String()
}

func toTitle(words []string) string {
	parts := make([]string, len(words))
	for i, w := range words {
		parts[i] = capitalize(w)
	}
	return strings.Join(parts, " ")
}

func toSentence(words []string) string {
	if len(words) == 0 {
		return ""
	}
	parts := make([]string, len(words))
	parts[0] = capitalize(words[0])
	for i, w := range words[1:] {
		parts[i+1] = strings.ToLower(w)
	}
	return strings.Join(parts, " ")
}

func joinLower(words []string, sep string) string {
	parts := make([]string, len(words))
	for i, w := range words {
		parts[i] = strings.ToLower(w)
	}
	return strings.Join(parts, sep)
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
