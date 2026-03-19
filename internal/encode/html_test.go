package encode

import (
	"testing"
)

func TestHTMLEncode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "ampersand",
			input: "A & B",
			want:  "A &amp; B",
		},
		{
			name:  "less_than",
			input: "a < b",
			want:  "a &lt; b",
		},
		{
			name:  "greater_than",
			input: "a > b",
			want:  "a &gt; b",
		},
		{
			name:  "double_quote",
			input: `say "hello"`,
			want:  "say &#34;hello&#34;",
		},
		{
			name:  "single_quote",
			input: "it's",
			want:  "it&#39;s",
		},
		{
			name:  "all_five_entities",
			input: `& < > " '`,
			want:  "&amp; &lt; &gt; &#34; &#39;",
		},
		{
			name:  "nested_html_tags",
			input: `<div class="test"><p>Hello & World</p></div>`,
			want:  "&lt;div class=&#34;test&#34;&gt;&lt;p&gt;Hello &amp; World&lt;/p&gt;&lt;/div&gt;",
		},
		{
			name:  "unicode_passthrough",
			input: "Hello World",
			want:  "Hello World",
		},
		{
			name:  "multi_line_html",
			input: "<h1>Title</h1>\n<p>Body</p>\n",
			want:  "&lt;h1&gt;Title&lt;/h1&gt;\n&lt;p&gt;Body&lt;/p&gt;\n",
		},
		{
			name:  "empty_string",
			input: "",
			want:  "",
		},
		{
			name:  "no_special_chars",
			input: "hello world 123",
			want:  "hello world 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HTMLEncode(tt.input)
			if got != tt.want {
				t.Errorf("HTMLEncode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTMLDecode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "named_entities",
			input: "&amp; &lt; &gt;",
			want:  "& < >",
		},
		{
			name:  "numeric_entities",
			input: "&#34;hello&#34;",
			want:  `"hello"`,
		},
		{
			name:  "single_quote_entity",
			input: "it&#39;s",
			want:  "it's",
		},
		{
			name:  "mixed_entities_and_text",
			input: "&lt;div&gt;Hello &amp; World&lt;/div&gt;",
			want:  "<div>Hello & World</div>",
		},
		{
			name:  "no_entities",
			input: "hello world",
			want:  "hello world",
		},
		{
			name:  "empty_string",
			input: "",
			want:  "",
		},
		{
			name:  "multi_line",
			input: "&lt;h1&gt;Title&lt;/h1&gt;\n&lt;p&gt;Body&lt;/p&gt;",
			want:  "<h1>Title</h1>\n<p>Body</p>",
		},
		{
			name:  "unicode_passthrough",
			input: "Hello World",
			want:  "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HTMLDecode(tt.input)
			if got != tt.want {
				t.Errorf("HTMLDecode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTMLRoundtrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"all_entities", `& < > " '`},
		{"nested_tags", `<div class="x"><p>A & B</p></div>`},
		{"multi_line", "<h1>Title</h1>\n<p>Body & more</p>"},
		{"empty", ""},
		{"plain_text", "no special chars here"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := HTMLEncode(tt.input)
			decoded := HTMLDecode(encoded)
			if decoded != tt.input {
				t.Errorf("roundtrip failed: got %q, want %q", decoded, tt.input)
			}
		})
	}
}
