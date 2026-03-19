package text

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// TextInfo holds analysis results for a text string.
type TextInfo struct {
	Characters int    `json:"characters"`
	Words      int    `json:"words"`
	Lines      int    `json:"lines"`
	Sentences  int    `json:"sentences"`
	Bytes      int    `json:"bytes"`
	BytesHuman string `json:"bytes_human"`
	IsASCII    bool   `json:"is_ascii"`
	HasUnicode bool   `json:"has_unicode"`
}

func Inspect(input string) *TextInfo {
	info := &TextInfo{}

	info.Bytes = len(input)
	info.BytesHuman = humanBytes(info.Bytes)
	info.Characters = utf8.RuneCountInString(input)

	words := strings.Fields(input)
	info.Words = len(words)

	if input == "" {
		info.Lines = 0
	} else {
		info.Lines = strings.Count(input, "\n") + 1
	}

	for _, r := range input {
		if r == '.' || r == '!' || r == '?' {
			info.Sentences++
		}
	}

	info.IsASCII = true
	for _, r := range input {
		if r > 127 {
			info.IsASCII = false
			info.HasUnicode = true
			break
		}
	}

	return info
}

func humanBytes(b int) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case b >= GB:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
