package gen

import (
	"strings"
	"testing"
	"unicode"
)

func TestGenerateWords(t *testing.T) {
	tests := []struct {
		name  string
		count int
		want  int
	}{
		{"one_word", 1, 1},
		{"five_words", 5, 5},
		{"ten_words", 10, 10},
		{"hundred_words", 100, 100},
		{"zero_words", 0, 0},
		{"negative_words", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateWords(tt.count)
			if tt.want == 0 {
				if got != "" {
					t.Errorf("GenerateWords(%d) = %q, want empty", tt.count, got)
				}
				return
			}
			words := strings.Fields(got)
			if len(words) != tt.want {
				t.Errorf("GenerateWords(%d) returned %d words, want %d", tt.count, len(words), tt.want)
			}
		})
	}
}

func TestGenerateWords_FirstWordCapitalized(t *testing.T) {
	for i := 0; i < 10; i++ {
		got := GenerateWords(5)
		if got == "" {
			t.Fatal("GenerateWords(5) returned empty")
		}
		if !unicode.IsUpper(rune(got[0])) {
			t.Errorf("first character not uppercase: %q", got)
		}
	}
}

func TestGenerateWords_ExactCount(t *testing.T) {
	for _, n := range []int{1, 3, 7, 15, 50} {
		got := GenerateWords(n)
		words := strings.Fields(got)
		if len(words) != n {
			t.Errorf("GenerateWords(%d) returned %d words", n, len(words))
		}
	}
}

func TestGenerateSentences(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{"one_sentence", 1},
		{"three_sentences", 3},
		{"five_sentences", 5},
		{"zero_sentences", 0},
		{"negative_sentences", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSentences(tt.count)
			if tt.count <= 0 {
				if got != "" {
					t.Errorf("GenerateSentences(%d) = %q, want empty", tt.count, got)
				}
				return
			}
			// Count periods as sentence endings.
			periods := strings.Count(got, ".")
			if periods != tt.count {
				t.Errorf("GenerateSentences(%d) has %d periods, want %d\ngot: %s", tt.count, periods, tt.count, got)
			}
		})
	}
}

func TestGenerateSentences_EndWithPeriod(t *testing.T) {
	for i := 0; i < 10; i++ {
		got := GenerateSentences(1)
		if !strings.HasSuffix(got, ".") {
			t.Errorf("sentence does not end with period: %q", got)
		}
	}
}

func TestGenerateSentences_StartWithCapital(t *testing.T) {
	for i := 0; i < 10; i++ {
		got := GenerateSentences(1)
		if got == "" {
			t.Fatal("empty sentence")
		}
		if !unicode.IsUpper(rune(got[0])) {
			t.Errorf("sentence does not start with capital: %q", got)
		}
	}
}

func TestGenerateSentences_MultipleStartWithCapital(t *testing.T) {
	got := GenerateSentences(3)
	// Split on ". " to find sentence boundaries.
	// Each sentence after the first starts after ". "
	parts := strings.Split(got, ". ")
	for i, part := range parts {
		if part == "" {
			continue
		}
		// Last part may end with just ".", trim it.
		part = strings.TrimSuffix(part, ".")
		if part == "" {
			continue
		}
		if !unicode.IsUpper(rune(part[0])) {
			t.Errorf("sentence %d does not start with capital: %q", i, part)
		}
	}
}

func TestGenerateParagraphs(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{"one_paragraph", 1},
		{"three_paragraphs", 3},
		{"five_paragraphs", 5},
		{"zero_paragraphs", 0},
		{"negative_paragraphs", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateParagraphs(tt.count)
			if tt.count <= 0 {
				if got != "" {
					t.Errorf("GenerateParagraphs(%d) = %q, want empty", tt.count, got)
				}
				return
			}
			paragraphs := strings.Split(got, "\n\n")
			if len(paragraphs) != tt.count {
				t.Errorf("GenerateParagraphs(%d) returned %d paragraphs, want %d", tt.count, len(paragraphs), tt.count)
			}
		})
	}
}

func TestGenerateParagraphs_SeparatedByDoubleNewline(t *testing.T) {
	got := GenerateParagraphs(3)
	if !strings.Contains(got, "\n\n") {
		t.Error("paragraphs not separated by double newline")
	}
	parts := strings.Split(got, "\n\n")
	if len(parts) != 3 {
		t.Errorf("expected 3 paragraphs separated by \\n\\n, got %d", len(parts))
	}
	for i, p := range parts {
		if strings.TrimSpace(p) == "" {
			t.Errorf("paragraph %d is empty", i)
		}
	}
}

func TestGenerateWords_LargeGeneration(t *testing.T) {
	got := GenerateWords(100)
	words := strings.Fields(got)
	if len(words) != 100 {
		t.Errorf("GenerateWords(100) returned %d words, want 100", len(words))
	}
}
