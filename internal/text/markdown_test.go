package text

import (
	"strings"
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		toHTML       bool
		wantContains []string
		wantMissing  []string
	}{
		// HTML mode
		{
			name:         "html heading h1",
			input:        "# Hello",
			toHTML:       true,
			wantContains: []string{"<h1>Hello</h1>"},
		},
		{
			name:         "html bold",
			input:        "**bold**",
			toHTML:       true,
			wantContains: []string{"<strong>bold</strong>"},
		},
		{
			name:         "html link",
			input:        "[link](http://example.com)",
			toHTML:       true,
			wantContains: []string{`<a href="http://example.com">link</a>`},
		},
		{
			name:         "html code block",
			input:        "```\ncode here\n```",
			toHTML:       true,
			wantContains: []string{"<code>"},
		},
		{
			name:         "html unordered list",
			input:        "- item1\n- item2",
			toHTML:       true,
			wantContains: []string{"<li>"},
		},
		{
			name:         "html italic",
			input:        "*italic*",
			toHTML:       true,
			wantContains: []string{"<em>italic</em>"},
		},
		{
			name:         "html heading h2",
			input:        "## Subtitle",
			toHTML:       true,
			wantContains: []string{"<h2>Subtitle</h2>"},
		},

		// Plain text mode
		{
			name:         "plain text heading stripped",
			input:        "# Hello",
			toHTML:       false,
			wantContains: []string{"Hello"},
			wantMissing:  []string{"#"},
		},
		{
			name:         "plain text bold stripped",
			input:        "**bold**",
			toHTML:       false,
			wantContains: []string{"bold"},
			wantMissing:  []string{"**"},
		},
		{
			name:         "plain text link shows text only",
			input:        "[link](http://example.com)",
			toHTML:       false,
			wantContains: []string{"link"},
			wantMissing:  []string{"]("},
		},

		// Empty input
		{
			name:   "empty input html",
			input:  "",
			toHTML: true,
		},
		{
			name:   "empty input plain",
			input:  "",
			toHTML: false,
		},

		// Code block with language in HTML
		{
			name:         "html code block with language",
			input:        "```go\nfmt.Println(\"hi\")\n```",
			toHTML:       true,
			wantContains: []string{"<pre>", "fmt.Println"},
		},
	}

	// Syntax highlight tests
	highlightTests := []struct {
		name  string
		input string
		theme string
		want  string
	}{
		{
			name:  "syntax highlight default theme",
			input: "# Hello\n```go\nfmt.Println()\n```",
			theme: "",
			want:  "highlight.js",
		},
		{
			name:  "syntax highlight custom theme",
			input: "# Hello",
			theme: "monokai",
			want:  "monokai",
		},
	}

	for _, tt := range highlightTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderMarkdown([]byte(tt.input), true, true, tt.theme)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(string(got), tt.want) {
				t.Errorf("expected output to contain %q", tt.want)
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderMarkdown([]byte(tt.input), tt.toHTML, false, "github")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			output := string(got)
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q.\nGot:\n%s", want, output)
				}
			}
			for _, notWant := range tt.wantMissing {
				if strings.Contains(output, notWant) {
					t.Errorf("output should not contain %q.\nGot:\n%s", notWant, output)
				}
			}
		})
	}
}

// Images are stripped before links, so no stray "!" is left behind.
func TestStripMarkdown_Images(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"image alt text", "![alt text](img.png)", "alt text"},
		{"image with empty alt", "![](img.png)", ""},
		{"link is still stripped", "[link](http://x)", "link"},
		{"image inside a sentence", "see ![a](b.png) here", "see a here"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderMarkdown([]byte(tt.input), false, false, "")
			if err != nil {
				t.Fatalf("RenderMarkdown: %v", err)
			}
			if strings.TrimRight(string(got), "\n") != tt.want {
				t.Errorf("got %q, want %q", string(got), tt.want)
			}
		})
	}
}

// Underscore emphasis applies at word boundaries only, so identifiers survive.
func TestStripMarkdown_Underscores(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"snake case preserved", "snake_case_name", "snake_case_name"},
		{"identifier preserved", "foo_bar_baz", "foo_bar_baz"},
		// CommonMark: the outer underscores flank the word so they are emphasis,
		// while the inner one sits between letters and stays literal.
		{"surrounded identifier is emphasis", "_private_field_", "private_field"},
		{"emphasis stripped", "_emphasis_", "emphasis"},
		{"emphasis in sentence", "a _b_ c", "a b c"},
		{"adjacent emphasis", "_a_ _b_", "a b"},
		{"emphasis next to punctuation", "(_x_)", "(x)"},
		{"mixed identifier and emphasis", "use _foo_ with bar_baz", "use foo with bar_baz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderMarkdown([]byte(tt.input), false, false, "")
			if err != nil {
				t.Fatalf("RenderMarkdown: %v", err)
			}
			if strings.TrimRight(string(got), "\n") != tt.want {
				t.Errorf("got %q, want %q", string(got), tt.want)
			}
		})
	}
}

func TestRenderMarkdown_ThemeValidation(t *testing.T) {
	if _, err := RenderMarkdown([]byte("# x"), true, true, "tokyo-night-dark"); err != nil {
		t.Errorf("valid theme rejected: %v", err)
	}

	bad := []string{`x"><script>alert(1)</script>`, "../../etc/passwd", "a b"}
	for _, theme := range bad {
		if _, err := RenderMarkdown([]byte("# x"), true, true, theme); err == nil {
			t.Errorf("theme %q should have been rejected", theme)
		}
	}
}

func TestRenderMarkdown_HTMLEndsWithNewline(t *testing.T) {
	got, err := RenderMarkdown([]byte("# x"), true, false, "")
	if err != nil {
		t.Fatalf("RenderMarkdown: %v", err)
	}
	if !strings.HasSuffix(string(got), "</html>\n") {
		t.Errorf("HTML output should end with a newline, got %q", string(got)[len(got)-20:])
	}
}
