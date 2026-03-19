package text

import (
	"strings"
	"testing"
)

func TestInspect(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantChars  int
		wantWords  int
		wantLines  int
		wantSent   int
		wantBytes  int
		wantASCII  bool
		wantUni    bool
		checkHuman string
	}{
		{
			name:       "empty string",
			input:      "",
			wantChars:  0,
			wantWords:  0,
			wantLines:  0,
			wantSent:   0,
			wantBytes:  0,
			wantASCII:  true,
			wantUni:    false,
			checkHuman: "0 B",
		},
		{
			name:      "single word",
			input:     "hello",
			wantChars: 5,
			wantWords: 1,
			wantLines: 1,
			wantBytes: 5,
			wantASCII: true,
			wantUni:   false,
		},
		{
			name:      "multi-word",
			input:     "hello world",
			wantChars: 11,
			wantWords: 2,
			wantLines: 1,
			wantBytes: 11,
			wantASCII: true,
			wantUni:   false,
		},
		{
			name:      "multi-line",
			input:     "line1\nline2\nline3",
			wantChars: 17,
			wantWords: 3,
			wantLines: 3,
			wantBytes: 17,
			wantASCII: true,
			wantUni:   false,
		},
		{
			name:      "unicode accented chars",
			input:     "h\u00e9llo",
			wantChars: 5,
			wantWords: 1,
			wantLines: 1,
			wantBytes: 6,
			wantASCII: false,
			wantUni:   true,
		},
		{
			name:      "emoji single char",
			input:     "\U0001F44B",
			wantChars: 1,
			wantWords: 1,
			wantLines: 1,
			wantBytes: 4,
			wantASCII: false,
			wantUni:   true,
		},
		{
			name:      "pure ASCII printable",
			input:     "Hello, World! 123",
			wantChars: 17,
			wantWords: 3,
			wantLines: 1,
			wantASCII: true,
			wantUni:   false,
		},
		{
			name:     "sentence counting",
			input:    "Hello. World! How?",
			wantSent: 3,
		},
		{
			name:       "BytesHuman for 1KB",
			input:      strings.Repeat("a", 1024),
			checkHuman: "1.0 KB",
		},
		{
			name:      "only whitespace has 0 words",
			input:     "   \t\n",
			wantWords: 0,
		},
		{
			name:      "tabs and spaces between words",
			input:     "word1\tword2  word3",
			wantWords: 3,
		},
		{
			name:      "trailing newline counts extra line",
			input:     "line1\n",
			wantLines: 2,
			wantChars: 6,
			wantWords: 1,
		},
		{
			name:      "Chinese characters",
			input:     "\u4f60\u597d\u4e16\u754c",
			wantChars: 4,
			wantBytes: 12,
			wantASCII: false,
			wantUni:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := Inspect(tt.input)

			// Only check fields that are explicitly set (non-zero or meaningful)
			if tt.name == "empty string" || tt.wantChars > 0 || tt.name == "single word" || tt.name == "multi-word" || tt.name == "multi-line" || tt.name == "unicode accented chars" || tt.name == "emoji single char" || tt.name == "pure ASCII printable" || tt.name == "trailing newline counts extra line" || tt.name == "Chinese characters" {
				if info.Characters != tt.wantChars {
					t.Errorf("Characters = %d, want %d", info.Characters, tt.wantChars)
				}
			}

			if tt.name != "sentence counting" && tt.name != "BytesHuman for 1KB" && tt.name != "Chinese characters" {
				if info.Words != tt.wantWords {
					t.Errorf("Words = %d, want %d", info.Words, tt.wantWords)
				}
			}

			if tt.wantLines > 0 {
				if info.Lines != tt.wantLines {
					t.Errorf("Lines = %d, want %d", info.Lines, tt.wantLines)
				}
			}

			if tt.wantSent > 0 {
				if info.Sentences != tt.wantSent {
					t.Errorf("Sentences = %d, want %d", info.Sentences, tt.wantSent)
				}
			}

			if tt.wantBytes > 0 {
				if info.Bytes != tt.wantBytes {
					t.Errorf("Bytes = %d, want %d", info.Bytes, tt.wantBytes)
				}
			}

			// Check ASCII/Unicode for tests that set those fields
			switch tt.name {
			case "empty string", "single word", "multi-word", "multi-line",
				"unicode accented chars", "emoji single char", "pure ASCII printable",
				"Chinese characters":
				if info.IsASCII != tt.wantASCII {
					t.Errorf("IsASCII = %v, want %v", info.IsASCII, tt.wantASCII)
				}
				if info.HasUnicode != tt.wantUni {
					t.Errorf("HasUnicode = %v, want %v", info.HasUnicode, tt.wantUni)
				}
			}

			if tt.checkHuman != "" {
				if info.BytesHuman != tt.checkHuman {
					t.Errorf("BytesHuman = %q, want %q", info.BytesHuman, tt.checkHuman)
				}
			}
		})
	}
}

func TestHumanBytes(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1572864, "1.5 MB"},
		{1073741824, "1.0 GB"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := humanBytes(tt.input)
			if got != tt.want {
				t.Errorf("humanBytes(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestInspectUnicodeBytesDiffer(t *testing.T) {
	info := Inspect("h\u00e9llo")
	if info.Characters != 5 {
		t.Errorf("Characters = %d, want 5", info.Characters)
	}
	if info.Bytes != 6 {
		t.Errorf("Bytes = %d, want 6", info.Bytes)
	}
	if info.Bytes <= info.Characters {
		t.Error("expected Bytes > Characters for unicode text")
	}
}
