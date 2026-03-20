package query

import (
	"testing"
)

func TestValidateXML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "valid simple XML",
			input: "<root><item>test</item></root>",
		},
		{
			name:  "valid XML with attributes",
			input: `<root attr="val"/>`,
		},
		{
			name:  "valid self-closing tag",
			input: "<br/>",
		},
		{
			name:    "missing closing tag returns error",
			input:   "<root><item>test</item>",
			wantErr: true,
		},
		{
			name:    "mismatched tags returns error",
			input:   "<a></b>",
			wantErr: true,
		},
		{
			name:    "plain text is not valid XML",
			input:   "just plain text",
			wantErr: true,
		},
		{
			name:    "empty input is not valid XML",
			input:   "",
			wantErr: true,
		},
		{
			name:  "XML with declaration",
			input: `<?xml version="1.0"?><root/>`,
		},
		{
			name:    "malformed XML unclosed angle bracket",
			input:   "<root><unclosed",
			wantErr: true,
		},
		{
			name:  "nested elements valid",
			input: `<root><a><b><c>deep</c></b></a></root>`,
		},
		{
			name:    "invalid entity reference",
			input:   `<root>&invalid;</root>`,
			wantErr: true,
		},
		{
			name:  "XML with CDATA section",
			input: `<root><![CDATA[some <data>]]></root>`,
		},
		{
			name:  "XML with namespace",
			input: `<root xmlns:ns="http://example.com"><ns:item>test</ns:item></root>`,
		},
		{
			name:    "unclosed root tag",
			input:   "<root>",
			wantErr: true,
		},
		{
			name: "multiple root-level elements",
			// Go xml decoder processes tokens sequentially; multiple roots are fine token-by-token.
			input: "<a/><b/>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateXML([]byte(tt.input))
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
