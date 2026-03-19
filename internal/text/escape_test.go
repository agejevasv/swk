package text

import (
	"testing"
)

func TestEscape(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		mode    string
		want    string
		wantErr bool
	}{
		// JSON escaping
		{name: "json quotes", input: `say "hello"`, mode: "json", want: `say \"hello\"`},
		{name: "json backslash", input: `path\to\file`, mode: "json", want: `path\\to\\file`},
		{name: "json newline", input: "line1\nline2", mode: "json", want: `line1\nline2`},
		{name: "json tab", input: "col1\tcol2", mode: "json", want: `col1\tcol2`},
		{name: "json empty", input: "", mode: "json", want: ""},
		{name: "json multi-line", input: "line1\nline2\nline3", mode: "json", want: `line1\nline2\nline3`},

		// XML escaping
		{name: "xml amp", input: "&", mode: "xml", want: "&amp;"},
		{name: "xml lt", input: "<", mode: "xml", want: "&lt;"},
		{name: "xml gt", input: ">", mode: "xml", want: "&gt;"},
		{name: "xml quot", input: `"`, mode: "xml", want: "&quot;"},
		{name: "xml apos", input: "'", mode: "xml", want: "&apos;"},
		{name: "xml all special", input: `& < > " '`, mode: "xml", want: `&amp; &lt; &gt; &quot; &apos;`},
		{name: "xml empty", input: "", mode: "xml", want: ""},

		// HTML escaping
		{name: "html amp", input: "&", mode: "html", want: "&amp;"},
		{name: "html lt", input: "<", mode: "html", want: "&lt;"},
		{name: "html gt", input: ">", mode: "html", want: "&gt;"},
		{name: "html empty", input: "", mode: "html", want: ""},

		// C escaping
		{name: "c newline", input: "hello\nworld", mode: "c", want: `hello\nworld`},
		{name: "c tab", input: "hello\tworld", mode: "c", want: `hello\tworld`},
		{name: "c empty", input: "", mode: "c", want: ""},

		// Shell escaping
		{name: "shell simple space", input: "hello world", mode: "shell", want: "'hello world'"},
		{name: "shell single quote", input: "it's", mode: "shell", want: `'it'\''s'`},
		{name: "shell special chars", input: "echo $HOME & rm -rf | yes!", mode: "shell", want: "'echo $HOME & rm -rf | yes!'"},
		{name: "shell empty", input: "", mode: "shell", want: "''"},

		// Invalid mode
		{name: "invalid mode", input: "test", mode: "invalid", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Escape(tt.input, tt.mode)
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
				t.Errorf("Escape(%q, %q) = %q, want %q", tt.input, tt.mode, got, tt.want)
			}
		})
	}
}

func TestUnescape(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		mode    string
		want    string
		wantErr bool
	}{
		{name: "json unescape quotes", input: `hello \"world\"`, mode: "json", want: `hello "world"`},
		{name: "json unescape newline", input: `line1\nline2`, mode: "json", want: "line1\nline2"},
		{name: "json unescape backslash", input: `path\\to\\file`, mode: "json", want: `path\to\file`},
		{name: "xml unescape all", input: "&amp; &lt; &gt; &quot; &apos;", mode: "xml", want: `& < > " '`},
		{name: "html unescape", input: "&amp; &lt; &gt;", mode: "html", want: "& < >"},
		{name: "shell unescape", input: "'hello world'", mode: "shell", want: "hello world"},
		{name: "c unescape backslash", input: `hello\\world`, mode: "c", want: `hello\world`},
		{name: "c unescape tab", input: `hello\tworld`, mode: "c", want: "hello\tworld"},
		{name: "invalid mode", input: "test", mode: "bogus", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unescape(tt.input, tt.mode)
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
				t.Errorf("Unescape(%q, %q) = %q, want %q", tt.input, tt.mode, got, tt.want)
			}
		})
	}
}

func TestEscapeUnescapeRoundtrip(t *testing.T) {
	modes := []string{"json", "xml", "html", "c", "shell"}
	inputs := []string{
		"hello world",
		"line1\nline2",
		`quotes "and" backslash \`,
		"special & < > chars",
		"it's a test",
		"",
		"tab\there",
		"Hello <world> & \"friends\"!\n\tGoodbye.",
	}

	for _, mode := range modes {
		for _, input := range inputs {
			t.Run(mode+"/"+input, func(t *testing.T) {
				escaped, err := Escape(input, mode)
				if err != nil {
					t.Fatalf("Escape error: %v", err)
				}
				unescaped, err := Unescape(escaped, mode)
				if err != nil {
					t.Fatalf("Unescape error: %v", err)
				}
				if unescaped != input {
					t.Errorf("roundtrip failed: %q -> %q -> %q", input, escaped, unescaped)
				}
			})
		}
	}
}
