package test

import (
	"testing"
)

func TestRegexTest(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		pattern    string
		global     bool
		wantMatch  bool
		wantValues []string
		wantGroups [][]string
		wantErr    bool
	}{
		{
			name:       "simple match digits",
			input:      "abc 123 def",
			pattern:    `\d+`,
			global:     false,
			wantMatch:  true,
			wantValues: []string{"123"},
		},
		{
			name:       "global match digits",
			input:      "abc 123 def 456",
			pattern:    `\d+`,
			global:     true,
			wantMatch:  true,
			wantValues: []string{"123", "456"},
		},
		{
			name:       "capture groups email-like",
			input:      "user@host",
			pattern:    `(\w+)@(\w+)`,
			global:     false,
			wantMatch:  true,
			wantValues: []string{"user@host"},
			wantGroups: [][]string{{"user", "host"}},
		},
		{
			name:       "date pattern with groups",
			input:      "Today is 2024-01-15 ok",
			pattern:    `(\d{4})-(\d{2})-(\d{2})`,
			global:     false,
			wantMatch:  true,
			wantValues: []string{"2024-01-15"},
			wantGroups: [][]string{{"2024", "01", "15"}},
		},
		{
			name:      "no match returns Matched=false",
			input:     "hello world",
			pattern:   `\d+`,
			global:    false,
			wantMatch: false,
		},
		{
			name:    "invalid pattern returns error",
			input:   "test",
			pattern: `[invalid`,
			global:  false,
			wantErr: true,
		},
		{
			name:       "multi-line input with (?m)^",
			input:      "hello world\nfoo bar\nbaz qux",
			pattern:    `(?m)^(\w+)`,
			global:     true,
			wantMatch:  true,
			wantValues: []string{"hello", "foo", "baz"},
			wantGroups: [][]string{{"hello"}, {"foo"}, {"baz"}},
		},
		{
			name:       "special regex chars in input (not pattern)",
			input:      "price is $100.00 (USD)",
			pattern:    `\$[\d.]+`,
			global:     false,
			wantMatch:  true,
			wantValues: []string{"$100.00"},
		},
		{
			name:      "empty input no match",
			input:     "",
			pattern:   `\w+`,
			global:    false,
			wantMatch: false,
		},
		{
			name:      "empty input global no match",
			input:     "",
			pattern:   `\w+`,
			global:    true,
			wantMatch: false,
		},
		{
			name:       "global with groups finds multiple",
			input:      "cat:3 dog:5 bird:12",
			pattern:    `(\w+):(\d+)`,
			global:     true,
			wantMatch:  true,
			wantValues: []string{"cat:3", "dog:5", "bird:12"},
			wantGroups: [][]string{{"cat", "3"}, {"dog", "5"}, {"bird", "12"}},
		},
		{
			name:       "match at start and end",
			input:      "123abc456",
			pattern:    `\d+`,
			global:     true,
			wantMatch:  true,
			wantValues: []string{"123", "456"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RegexTest(tt.input, tt.pattern, tt.global)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Matched != tt.wantMatch {
				t.Fatalf("Matched = %v, want %v", result.Matched, tt.wantMatch)
			}
			if !tt.wantMatch {
				return
			}
			if len(result.Matches) != len(tt.wantValues) {
				t.Fatalf("got %d matches, want %d", len(result.Matches), len(tt.wantValues))
			}
			for i, wv := range tt.wantValues {
				if result.Matches[i].Value != wv {
					t.Errorf("match[%d].Value = %q, want %q", i, result.Matches[i].Value, wv)
				}
			}
			if tt.wantGroups != nil {
				for i, wg := range tt.wantGroups {
					if len(result.Matches[i].Groups) != len(wg) {
						t.Errorf("match[%d] got %d groups, want %d", i, len(result.Matches[i].Groups), len(wg))
						continue
					}
					for j, g := range wg {
						if result.Matches[i].Groups[j] != g {
							t.Errorf("match[%d].Groups[%d] = %q, want %q", i, j, result.Matches[i].Groups[j], g)
						}
					}
				}
			}
		})
	}
}

func TestRegexReplace(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		pattern     string
		replacement string
		want        string
		wantErr     bool
	}{
		{
			name:        "simple replace all occurrences",
			input:       "foo baz foo",
			pattern:     `foo`,
			replacement: "bar",
			want:        "bar baz bar",
		},
		{
			name:        "replace with backreference swap",
			input:       "hello world",
			pattern:     `(\w+) (\w+)`,
			replacement: "$2 $1",
			want:        "world hello",
		},
		{
			name:    "invalid pattern returns error",
			input:   "test",
			pattern: `[bad`,
			wantErr: true,
		},
		{
			name:        "no match returns original",
			input:       "hello",
			pattern:     `xyz`,
			replacement: "abc",
			want:        "hello",
		},
		{
			name:        "replace digits with placeholder",
			input:       "order 123 and 456",
			pattern:     `\d+`,
			replacement: "NUM",
			want:        "order NUM and NUM",
		},
		{
			name:        "replace with group reformat date",
			input:       "2024-01-15",
			pattern:     `(\d{4})-(\d{2})-(\d{2})`,
			replacement: "$2/$3/$1",
			want:        "01/15/2024",
		},
		{
			name:        "empty replacement deletes matches",
			input:       "hello 123 world",
			pattern:     `\d+`,
			replacement: "",
			want:        "hello  world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegexReplace(tt.input, tt.pattern, tt.replacement)
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
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
