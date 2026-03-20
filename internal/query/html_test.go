package query

import (
	"testing"
)

func TestHTMLQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		selector string
		attr     string
		want     []string
		wantErr  bool
	}{
		// Text extraction: simple HTML
		{
			name:     "simple paragraph text",
			input:    `<p>Hello, world!</p>`,
			selector: "p",
			want:     []string{"Hello, world!"},
		},
		{
			name:     "multiple paragraphs",
			input:    `<p>First</p><p>Second</p><p>Third</p>`,
			selector: "p",
			want:     []string{"First", "Second", "Third"},
		},
		{
			name:     "text with inline tags stripped",
			input:    `<p>Hello <strong>bold</strong> world</p>`,
			selector: "p",
			want:     []string{"Hello bold world"},
		},
		// Text extraction: nested elements
		{
			name:     "nested div text",
			input:    `<div><span>inner</span></div>`,
			selector: "span",
			want:     []string{"inner"},
		},
		{
			name:     "deeply nested text",
			input:    `<div><ul><li>one</li><li>two</li></ul></div>`,
			selector: "li",
			want:     []string{"one", "two"},
		},
		// Text extraction: multiple matches
		{
			name:     "multiple divs",
			input:    `<div>A</div><div>B</div>`,
			selector: "div",
			want:     []string{"A", "B"},
		},

		// Attribute extraction
		{
			name:     "href attribute",
			input:    `<a href="https://example.com">Link</a>`,
			selector: "a",
			attr:     "href",
			want:     []string{"https://example.com"},
		},
		{
			name:     "multiple href attributes",
			input:    `<a href="/a">A</a><a href="/b">B</a>`,
			selector: "a",
			attr:     "href",
			want:     []string{"/a", "/b"},
		},
		{
			name:     "class attribute",
			input:    `<div class="main">content</div>`,
			selector: "div",
			attr:     "class",
			want:     []string{"main"},
		},
		{
			name:     "id attribute",
			input:    `<div id="header">Header</div>`,
			selector: "div",
			attr:     "id",
			want:     []string{"header"},
		},
		{
			name:     "src attribute on img",
			input:    `<img src="photo.jpg"><img src="logo.png">`,
			selector: "img",
			attr:     "src",
			want:     []string{"photo.jpg", "logo.png"},
		},
		{
			name:     "attribute not present on element",
			input:    `<div>no id here</div>`,
			selector: "div",
			attr:     "id",
			want:     nil,
		},

		// CSS selectors: tag
		{
			name:     "tag selector div",
			input:    `<div>D</div><p>P</p><div>D2</div>`,
			selector: "div",
			want:     []string{"D", "D2"},
		},
		{
			name:     "tag selector p",
			input:    `<p>one</p><span>skip</span><p>two</p>`,
			selector: "p",
			want:     []string{"one", "two"},
		},
		// CSS selectors: class
		{
			name:     "class selector",
			input:    `<p class="highlight">yes</p><p>no</p><p class="highlight">also</p>`,
			selector: ".highlight",
			want:     []string{"yes", "also"},
		},
		// CSS selectors: ID
		{
			name:     "id selector",
			input:    `<div id="main">Main</div><div id="sidebar">Side</div>`,
			selector: "#main",
			want:     []string{"Main"},
		},
		// CSS selectors: attribute selector
		{
			name:     "attribute selector [data-value]",
			input:    `<span data-value="x">X</span><span>Y</span>`,
			selector: "[data-value]",
			want:     []string{"X"},
		},
		{
			name:     "attribute value selector",
			input:    `<input type="text"><input type="password">`,
			selector: `[type="text"]`,
			attr:     "type",
			want:     []string{"text"},
		},
		// CSS selectors: descendant
		{
			name:     "descendant selector div p",
			input:    `<div><p>Inside</p></div><p>Outside</p>`,
			selector: "div p",
			want:     []string{"Inside"},
		},
		{
			name:     "descendant selector ul li",
			input:    `<ul><li>A</li><li>B</li></ul><ol><li>C</li></ol>`,
			selector: "ul li",
			want:     []string{"A", "B"},
		},

		// Edge cases
		{
			name:     "no matches returns empty slice",
			input:    `<div>Hello</div>`,
			selector: "span",
			want:     nil,
		},
		{
			name:     "empty input",
			input:    "",
			selector: "p",
			want:     nil,
		},
		{
			name:     "whitespace trimming",
			input:    `<p>  trimmed  </p>`,
			selector: "p",
			want:     []string{"trimmed"},
		},
		{
			name:     "newlines and tabs trimmed",
			input:    "<p>\n\tspaced\n</p>",
			selector: "p",
			want:     []string{"spaced"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HTMLQuery(tt.input, tt.selector, tt.attr)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d results %v, want %d results %v", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("result[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
