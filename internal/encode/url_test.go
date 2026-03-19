package encode

import (
	"testing"
)

func TestURLEncode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		component bool
		want      string
	}{
		{
			name:      "component_space_as_plus",
			input:     "hello world",
			component: true,
			want:      "hello+world",
		},
		{
			name:      "path_space_as_percent20",
			input:     "hello world",
			component: false,
			want:      "hello%20world",
		},
		{
			name:      "special_chars_ampersand",
			input:     "a=1&b=2",
			component: true,
			want:      "a%3D1%26b%3D2",
		},
		{
			name:      "question_mark_and_hash",
			input:     "page?q=test#section",
			component: true,
			want:      "page%3Fq%3Dtest%23section",
		},
		{
			name:      "percent_sign",
			input:     "100%",
			component: true,
			want:      "100%25",
		},
		{
			name:      "plus_sign",
			input:     "a+b",
			component: true,
			want:      "a%2Bb",
		},
		{
			name:      "forward_slash_component",
			input:     "path/to/file",
			component: true,
			want:      "path%2Fto%2Ffile",
		},
		{
			name:      "forward_slash_path_mode",
			input:     "path/to/file",
			component: false,
			want:      "path%2Fto%2Ffile",
		},
		{
			name:      "unicode_emoji",
			input:     "hello \xf0\x9f\x91\x8b",
			component: true,
			want:      "hello+%F0%9F%91%8B",
		},
		{
			name:      "unicode_cjk",
			input:     "\xe4\xb8\x96\xe7\x95\x8c",
			component: true,
			want:      "%E4%B8%96%E7%95%8C",
		},
		{
			name:      "empty_string",
			input:     "",
			component: true,
			want:      "",
		},
		{
			name:      "empty_string_path",
			input:     "",
			component: false,
			want:      "",
		},
		{
			name:      "already_safe_chars",
			input:     "abc123",
			component: true,
			want:      "abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := URLEncode(tt.input, tt.component)
			if tt.want != "" && got != tt.want {
				t.Errorf("URLEncode() = %q, want %q", got, tt.want)
			}
			// For unicode tests, just verify it encodes without error and contains percent encoding.
			if tt.want == "" && tt.input != "" && got == tt.input {
				t.Errorf("URLEncode() did not encode input %q", tt.input)
			}
		})
	}
}

func TestURLDecode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		component bool
		want      string
		wantErr   bool
	}{
		{
			name:      "component_plus_as_space",
			input:     "hello+world",
			component: true,
			want:      "hello world",
		},
		{
			name:      "path_plus_literal",
			input:     "hello+world",
			component: false,
			want:      "hello+world",
		},
		{
			name:      "percent_20_space",
			input:     "hello%20world",
			component: false,
			want:      "hello world",
		},
		{
			name:      "percent_20_component",
			input:     "hello%20world",
			component: true,
			want:      "hello world",
		},
		{
			name:      "special_chars",
			input:     "a%3D1%26b%3D2",
			component: true,
			want:      "a=1&b=2",
		},
		{
			name:      "empty_string",
			input:     "",
			component: true,
			want:      "",
		},
		{
			name:      "no_encoding_needed",
			input:     "abc123",
			component: true,
			want:      "abc123",
		},
		{
			name:      "invalid_percent_encoding",
			input:     "%ZZ",
			component: true,
			wantErr:   true,
		},
		{
			name:      "incomplete_percent_encoding",
			input:     "%2",
			component: true,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URLDecode(tt.input, tt.component)
			if (err != nil) != tt.wantErr {
				t.Fatalf("URLDecode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("URLDecode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestURLRoundtrip(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		component bool
	}{
		{"special_chars_component", "key=value&foo=bar#baz", true},
		{"spaces_component", "hello world", true},
		{"spaces_path", "hello world", false},
		{"unicode", "cafe\u0301", true},
		{"empty", "", true},
		{"slashes_component", "/path/to/resource", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := URLEncode(tt.input, tt.component)
			decoded, err := URLDecode(encoded, tt.component)
			if err != nil {
				t.Fatalf("roundtrip decode error: %v", err)
			}
			if decoded != tt.input {
				t.Errorf("roundtrip failed: got %q, want %q", decoded, tt.input)
			}
		})
	}
}
