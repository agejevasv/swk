package fmt

import (
	"strings"
	"testing"
)

func TestFormatXML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		opts    XMLOptions
		wantErr bool
		checkFn func(t *testing.T, got string)
	}{
		{
			name:  "pretty_print_indent_2",
			input: `<root><child>text</child></root>`,
			opts:  XMLOptions{Indent: 2},
			checkFn: func(t *testing.T, got string) {
				if !strings.Contains(got, "\n") {
					t.Error("expected indented output with newlines")
				}
				if !strings.Contains(got, "  <child>") {
					t.Errorf("expected 2-space indented child element, got:\n%s", got)
				}
				if !strings.Contains(got, "text") {
					t.Error("content missing")
				}
			},
		},
		{
			name:  "minify_removes_whitespace",
			input: "<root>\n  <child>text</child>\n</root>",
			opts:  XMLOptions{Minify: true},
			checkFn: func(t *testing.T, got string) {
				got = strings.TrimSpace(got)
				if strings.Contains(got, "\n") {
					t.Errorf("expected minified output without newlines, got:\n%s", got)
				}
				if !strings.Contains(got, "<root>") || !strings.Contains(got, "<child>") {
					t.Error("tags missing from minified output")
				}
			},
		},
		{
			name:  "self_closing_tags",
			input: `<root><empty/></root>`,
			opts:  XMLOptions{Indent: 2},
			checkFn: func(t *testing.T, got string) {
				if !strings.Contains(got, "<root>") {
					t.Error("root tag missing")
				}
				// Self-closing tags may be rendered as <empty></empty> by Go's xml encoder.
			},
		},
		{
			name:  "attributes_preserved",
			input: `<root attr="val"><child id="1">text</child></root>`,
			opts:  XMLOptions{Indent: 2},
			checkFn: func(t *testing.T, got string) {
				if !strings.Contains(got, `attr="val"`) {
					t.Errorf("attribute not preserved, got:\n%s", got)
				}
				if !strings.Contains(got, `id="1"`) {
					t.Errorf("child attribute not preserved, got:\n%s", got)
				}
			},
		},
		{
			name:  "xml_declaration_preserved",
			input: `<?xml version="1.0" encoding="UTF-8"?><root><child>text</child></root>`,
			opts:  XMLOptions{Indent: 2},
			checkFn: func(t *testing.T, got string) {
				if !strings.Contains(got, "<?xml") {
					t.Errorf("XML declaration not preserved, got:\n%s", got)
				}
			},
		},
		{
			name:  "deeply_nested_5_levels",
			input: `<a><b><c><d><e>deep</e></d></c></b></a>`,
			opts:  XMLOptions{Indent: 2},
			checkFn: func(t *testing.T, got string) {
				if !strings.Contains(got, "        <e>") {
					t.Errorf("expected 8-space indent for level 4, got:\n%s", got)
				}
			},
		},
		{
			name:  "pretty_print_indent_4",
			input: `<root><child><nested>val</nested></child></root>`,
			opts:  XMLOptions{Indent: 4},
			checkFn: func(t *testing.T, got string) {
				if !strings.Contains(got, "    <child>") {
					t.Errorf("expected 4-space indent, got:\n%s", got)
				}
				if !strings.Contains(got, "        <nested>") {
					t.Errorf("expected 8-space indent for nested, got:\n%s", got)
				}
			},
		},
		{
			name:  "pretty_print_has_trailing_newline",
			input: `<root><child>text</child></root>`,
			opts:  XMLOptions{Indent: 2},
			checkFn: func(t *testing.T, got string) {
				if !strings.HasSuffix(got, "\n") {
					t.Error("expected trailing newline for pretty-printed output")
				}
			},
		},
		{
			name:  "minify_no_trailing_newline",
			input: `<root><child>text</child></root>`,
			opts:  XMLOptions{Minify: true},
			checkFn: func(t *testing.T, got string) {
				if strings.HasSuffix(got, "\n") {
					t.Error("minified output should not have trailing newline")
				}
			},
		},
		{
			name:  "multiple_children",
			input: `<root><a>1</a><b>2</b><c>3</c></root>`,
			opts:  XMLOptions{Indent: 2},
			checkFn: func(t *testing.T, got string) {
				if !strings.Contains(got, "  <a>") || !strings.Contains(got, "  <b>") || !strings.Contains(got, "  <c>") {
					t.Errorf("expected all children indented, got:\n%s", got)
				}
			},
		},

		// Error cases.
		{
			name:    "invalid_xml_unclosed_tag",
			input:   `<root><unclosed>`,
			opts:    XMLOptions{Indent: 2},
			wantErr: true,
		},
		{
			name:    "empty_input",
			input:   "",
			opts:    XMLOptions{Indent: 2},
			wantErr: true,
		},
		{
			name:  "whitespace_only",
			input: "   \n\t  ",
			opts:  XMLOptions{Indent: 2},
			checkFn: func(t *testing.T, got string) {
				// Whitespace-only input produces empty output (whitespace CharData is skipped)
				if strings.TrimSpace(got) != "" {
					t.Errorf("expected empty output for whitespace-only input, got %q", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatXML([]byte(tt.input), tt.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FormatXML() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFn != nil && !tt.wantErr {
				tt.checkFn(t, string(got))
			}
		})
	}
}
