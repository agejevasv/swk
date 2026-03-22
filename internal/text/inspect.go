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
		KiB = 1024
		MiB = 1024 * KiB
		GiB = 1024 * MiB
	)
	switch {
	case b >= GiB:
		return fmt.Sprintf("%.1f GiB", float64(b)/float64(GiB))
	case b >= MiB:
		return fmt.Sprintf("%.1f MiB", float64(b)/float64(MiB))
	case b >= KiB:
		return fmt.Sprintf("%.1f KiB", float64(b)/float64(KiB))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
