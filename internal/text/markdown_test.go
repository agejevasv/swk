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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderMarkdown([]byte(tt.input), tt.toHTML, "github")
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
