package text

import (
	"testing"
)

func TestConvertCase(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		to      string
		want    string
		wantErr bool
	}{
		// upper
		{name: "upper basic", input: "hello world", to: "upper", want: "HELLO WORLD"},
		{name: "upper preserves newlines", input: "hello\nworld\n", to: "upper", want: "HELLO\nWORLD\n"},
		{name: "upper preserves tabs", input: "hello\tworld", to: "upper", want: "HELLO\tWORLD"},
		{name: "upper preserves indentation", input: "\t  hello", to: "upper", want: "\t  HELLO"},
		{name: "upper multi-line code", input: "func main() {\n\tfmt.Println()\n}", to: "upper", want: "FUNC MAIN() {\n\tFMT.PRINTLN()\n}"},
		{name: "upper single word", input: "hello", to: "upper", want: "HELLO"},
		{name: "upper empty", input: "", to: "upper", want: ""},

		// lower
		{name: "lower basic", input: "HELLO WORLD", to: "lower", want: "hello world"},
		{name: "lower preserves structure", input: "HELLO\n\tWORLD", to: "lower", want: "hello\n\tworld"},

		// snake_case
		{name: "snake from camel", input: "helloWorld", to: "snake", want: "hello_world"},
		{name: "snake from pascal", input: "HelloWorld", to: "snake", want: "hello_world"},
		{name: "snake from kebab", input: "hello-world", to: "snake", want: "hello_world"},
		{name: "snake from dot", input: "hello.world", to: "snake", want: "hello_world"},
		{name: "snake acronym XMLParser", input: "XMLParser", to: "snake", want: "xml_parser"},

		// camelCase
		{name: "camel from snake", input: "hello_world", to: "camel", want: "helloWorld"},
		{name: "camel from kebab", input: "hello-world", to: "camel", want: "helloWorld"},

		// PascalCase
		{name: "pascal from snake", input: "hello_world", to: "pascal", want: "HelloWorld"},

		// kebab-case
		{name: "kebab from camel", input: "helloWorld", to: "kebab", want: "hello-world"},

		// Title Case
		{name: "title basic", input: "hello world", to: "title", want: "Hello World"},

		// Sentence case
		{name: "sentence basic", input: "hello world foo", to: "sentence", want: "Hello world foo"},

		// dot.case
		{name: "dot from space-separated", input: "hello world", to: "dot", want: "hello.world"},

		// path/case
		{name: "path from space-separated", input: "hello world", to: "path", want: "hello/world"},

		// From various formats
		{name: "camel to pascal", input: "helloWorld", to: "pascal", want: "HelloWorld"},
		{name: "snake to kebab", input: "hello_world", to: "kebab", want: "hello-world"},
		{name: "kebab to pascal", input: "hello-world", to: "pascal", want: "HelloWorld"},

		// Empty
		{name: "empty snake", input: "", to: "snake", want: ""},
		{name: "empty camel", input: "", to: "camel", want: ""},
		{name: "empty upper", input: "", to: "upper", want: ""},

		// Unsupported
		{name: "unsupported case returns error", input: "hello", to: "nosuchcase", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertCase(tt.input, tt.to)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ConvertCase(%q, %q) = %q, want %q", tt.input, tt.to, got, tt.want)
			}
		})
	}
}
