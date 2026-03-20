package ioutil

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadInput(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		stdin   io.Reader
		want    string
		wantErr bool
	}{
		{
			name:  "args returns exact bytes",
			args:  []string{"hello world"},
			stdin: strings.NewReader("ignored"),
			want:  "hello world",
		},
		{
			name:  "stdin returns exact bytes including trailing newline",
			args:  nil,
			stdin: strings.NewReader("hello\n"),
			want:  "hello\n",
		},
		{
			name:  "args takes priority over stdin",
			args:  []string{"from-args"},
			stdin: strings.NewReader("from-stdin"),
			want:  "from-args",
		},
		{
			name:    "nil stdin returns error",
			args:    nil,
			stdin:   nil,
			wantErr: true,
		},
		{
			name:  "empty args falls through to stdin",
			args:  []string{},
			stdin: strings.NewReader("from-stdin"),
			want:  "from-stdin",
		},
		{
			name:  "binary-like bytes preserved",
			args:  []string{"\x00\x01\x02"},
			stdin: nil,
			want:  "\x00\x01\x02",
		},
		{
			name:  "empty stdin returns empty",
			args:  nil,
			stdin: strings.NewReader(""),
			want:  "",
		},
		{
			name:  "multi-line stdin preserved",
			args:  nil,
			stdin: strings.NewReader("line1\nline2\nline3\n"),
			want:  "line1\nline2\nline3\n",
		},
		{
			name:  "dash reads stdin",
			args:  []string{"-"},
			stdin: strings.NewReader("from-stdin"),
			want:  "from-stdin",
		},
		{
			name:    "dash with nil stdin returns error",
			args:    []string{"-"},
			stdin:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadInput(tt.args, tt.stdin)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("got %q, want %q", string(got), tt.want)
			}
		})
	}
}

func TestReadInputString(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		stdin   io.Reader
		want    string
		wantErr bool
	}{
		{
			name:  "from args returns string without trailing newline",
			args:  []string{"hello"},
			stdin: nil,
			want:  "hello",
		},
		{
			name:  "trims trailing LF from stdin",
			args:  nil,
			stdin: strings.NewReader("hello\n"),
			want:  "hello",
		},
		{
			name:  "trims trailing CRLF (Windows)",
			args:  nil,
			stdin: strings.NewReader("hello\r\n"),
			want:  "hello",
		},
		{
			name:  "preserves internal newlines",
			args:  nil,
			stdin: strings.NewReader("line1\nline2\n"),
			want:  "line1\nline2",
		},
		{
			name:  "no trailing newline returns as-is",
			args:  nil,
			stdin: strings.NewReader("hello"),
			want:  "hello",
		},
		{
			name:  "multiple trailing newlines all trimmed",
			args:  nil,
			stdin: strings.NewReader("hello\n\n\n"),
			want:  "hello",
		},
		{
			name:  "only newlines returns empty",
			args:  nil,
			stdin: strings.NewReader("\n\n"),
			want:  "",
		},
		{
			name:  "preserves internal CRLF but trims trailing",
			args:  nil,
			stdin: strings.NewReader("line1\r\nline2\r\n"),
			want:  "line1\r\nline2",
		},
		{
			name:    "nil stdin returns error",
			args:    nil,
			stdin:   nil,
			wantErr: true,
		},
		{
			name:  "empty input returns empty",
			args:  nil,
			stdin: strings.NewReader(""),
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadInputString(tt.args, tt.stdin)
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

func TestReadFileInput(t *testing.T) {
	// Create a temp file for file-based tests.
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(filePath, []byte("file content\n"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	tests := []struct {
		name    string
		args    []string
		stdin   io.Reader
		want    string
		wantErr bool
	}{
		{
			name:  "file path reads file content",
			args:  []string{filePath},
			stdin: nil,
			want:  "file content\n",
		},
		{
			name:  "dash reads stdin",
			args:  []string{"-"},
			stdin: strings.NewReader("from-stdin"),
			want:  "from-stdin",
		},
		{
			name:    "dash with nil stdin returns error",
			args:    []string{"-"},
			stdin:   nil,
			wantErr: true,
		},
		{
			name:  "non-existent path treated as literal",
			args:  []string{"/no/such/path/file.txt"},
			stdin: nil,
			want:  "/no/such/path/file.txt",
		},
		{
			name:  "directory path treated as literal",
			args:  []string{dir},
			stdin: nil,
			want:  dir,
		},
		{
			name:  "no args reads stdin",
			args:  nil,
			stdin: strings.NewReader("from-stdin"),
			want:  "from-stdin",
		},
		{
			name:    "no args nil stdin returns error",
			args:    nil,
			stdin:   nil,
			wantErr: true,
		},
		{
			name:  "literal content still works",
			args:  []string{`{"key":"value"}`},
			stdin: nil,
			want:  `{"key":"value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFileInput(tt.args, tt.stdin)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("got %q, want %q", string(got), tt.want)
			}
		})
	}
}

func TestReadFileInputString(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(filePath, []byte("file content\n"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	tests := []struct {
		name    string
		args    []string
		stdin   io.Reader
		want    string
		wantErr bool
	}{
		{
			name:  "file path reads and trims content",
			args:  []string{filePath},
			stdin: nil,
			want:  "file content",
		},
		{
			name:  "dash reads stdin and trims",
			args:  []string{"-"},
			stdin: strings.NewReader("from-stdin\n"),
			want:  "from-stdin",
		},
		{
			name:  "non-existent path treated as literal",
			args:  []string{"not-a-file"},
			stdin: nil,
			want:  "not-a-file",
		},
		{
			name:  "literal JSON content still works",
			args:  []string{`{"a":1}`},
			stdin: nil,
			want:  `{"a":1}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFileInputString(tt.args, tt.stdin)
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

func TestReadStdin(t *testing.T) {
	t.Run("nil returns error", func(t *testing.T) {
		_, err := ReadStdin(nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("valid reader returns content", func(t *testing.T) {
		got, err := ReadStdin(strings.NewReader("test data"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "test data" {
			t.Errorf("got %q, want %q", string(got), "test data")
		}
	})
}
