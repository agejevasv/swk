package encode

import (
	"strings"
	"testing"
)

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		urlSafe   bool
		noPadding bool
		want      string
	}{
		{
			name:  "simple_string",
			input: []byte("hello world"),
			want:  "aGVsbG8gd29ybGQ=",
		},
		{
			name:  "empty_input",
			input: []byte(""),
			want:  "",
		},
		{
			name:  "multi_line_input_preserves_newlines",
			input: []byte("line1\nline2\nline3\n"),
			want:  "bGluZTEKbGluZTIKbGluZTMK",
		},
		{
			name:    "url_safe_no_plus_slash",
			input:   []byte{0xfb, 0xff, 0xfe},
			urlSafe: true,
			want:    "-__-",
		},
		{
			name:  "standard_encoding_has_plus_slash",
			input: []byte{0xfb, 0xff, 0xfe},
			want:  "+//+",
		},
		{
			name:      "no_padding_removes_equals",
			input:     []byte("hello world"),
			noPadding: true,
			want:      "aGVsbG8gd29ybGQ",
		},
		{
			name:      "url_safe_and_no_padding",
			input:     []byte{0xfb, 0xff, 0xfe},
			urlSafe:   true,
			noPadding: true,
			want:      "-__-",
		},
		{
			name:  "single_byte",
			input: []byte("a"),
			want:  "YQ==",
		},
		{
			name:      "single_byte_no_padding",
			input:     []byte("a"),
			noPadding: true,
			want:      "YQ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Base64Encode(tt.input, tt.urlSafe, tt.noPadding)
			if got != tt.want {
				t.Errorf("Base64Encode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBase64Encode_AllByteValues(t *testing.T) {
	// Encode all 256 byte values and verify roundtrip.
	input := make([]byte, 256)
	for i := range input {
		input[i] = byte(i)
	}
	encoded := Base64Encode(input, false, false)
	decoded, err := Base64Decode(encoded, false)
	if err != nil {
		t.Fatalf("roundtrip decode error: %v", err)
	}
	if len(decoded) != 256 {
		t.Fatalf("roundtrip length = %d, want 256", len(decoded))
	}
	for i, b := range decoded {
		if b != byte(i) {
			t.Fatalf("roundtrip mismatch at byte %d: got %d, want %d", i, b, i)
		}
	}
}

func TestBase64Decode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		urlSafe bool
		want    string
		wantErr bool
	}{
		{
			name:  "simple_string",
			input: "aGVsbG8gd29ybGQ=",
			want:  "hello world",
		},
		{
			name:  "no_padding_input",
			input: "aGVsbG8gd29ybGQ",
			want:  "hello world",
		},
		{
			name:  "empty_input",
			input: "",
			want:  "",
		},
		{
			name:    "url_safe_decoding",
			input:   "-__-",
			urlSafe: true,
			want:    "\xfb\xff\xfe",
		},
		{
			name:  "whitespace_trimmed",
			input: "  aGVsbG8gd29ybGQ=  ",
			want:  "hello world",
		},
		{
			name:  "multi_line_roundtrip",
			input: "bGluZTEKbGluZTIKbGluZTMK",
			want:  "line1\nline2\nline3\n",
		},
		{
			name:    "invalid_base64",
			input:   "!!!invalid!!!",
			wantErr: true,
		},
		{
			name:    "invalid_base64_single_char",
			input:   "A",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64Decode(tt.input, tt.urlSafe)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Base64Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && string(got) != tt.want {
				t.Errorf("Base64Decode() = %q, want %q", string(got), tt.want)
			}
		})
	}
}

func TestBase64Roundtrip(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		urlSafe bool
	}{
		{"ascii", "Hello, World!", false},
		{"binary_like", "\x00\x01\x02\xff\xfe\xfd", false},
		{"unicode", "Hello, World!", false},
		{"url_safe", "test+data/here=now", true},
		{"empty", "", false},
		{"newlines", "line1\nline2\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Base64Encode([]byte(tt.input), tt.urlSafe, false)
			decoded, err := Base64Decode(encoded, tt.urlSafe)
			if err != nil {
				t.Fatalf("roundtrip decode error: %v", err)
			}
			if string(decoded) != tt.input {
				t.Errorf("roundtrip failed: got %q, want %q", string(decoded), tt.input)
			}
		})
	}
}

func TestBase64URLSafe_NoForbiddenChars(t *testing.T) {
	// URL-safe encoding should never contain +, /, or = (with noPadding).
	input := make([]byte, 256)
	for i := range input {
		input[i] = byte(i)
	}
	encoded := Base64Encode(input, true, true)
	if strings.ContainsAny(encoded, "+/=") {
		t.Errorf("URL-safe no-padding encoding contains forbidden chars: %s", encoded)
	}
}
