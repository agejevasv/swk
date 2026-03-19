package encode

import (
	"strings"
	"testing"
)

func TestQRTerminal(t *testing.T) {
	tests := []struct {
		name  string
		input string
		level string
	}{
		{
			name:  "short_text",
			input: "hello",
			level: "M",
		},
		{
			name:  "url",
			input: "https://example.com/path?q=test&lang=en",
			level: "M",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := QRTerminal(tt.input, tt.level)
			if err != nil {
				t.Fatalf("QRTerminal() error = %v", err)
			}

			// Should contain block characters.
			if !strings.ContainsAny(result, "\u2588\u2580\u2584") {
				t.Error("QRTerminal() output does not contain expected block characters")
			}

			// Should have multiple lines.
			lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
			if len(lines) < 5 {
				t.Errorf("QRTerminal() output too few lines: %d", len(lines))
			}
		})
	}
}

func TestQRTerminal_AllLevels(t *testing.T) {
	levels := []string{"L", "M", "Q", "H"}
	for _, level := range levels {
		t.Run("level_"+level, func(t *testing.T) {
			result, err := QRTerminal("test", level)
			if err != nil {
				t.Fatalf("QRTerminal(level=%s) error = %v", level, err)
			}
			if len(result) == 0 {
				t.Errorf("QRTerminal(level=%s) returned empty output", level)
			}
		})
	}
}

func TestQRTerminal_CaseInsensitiveLevel(t *testing.T) {
	// parseQRLevel uses strings.ToUpper, so lowercase should work.
	_, err := QRTerminal("test", "l")
	if err != nil {
		t.Fatalf("QRTerminal(level=l) error = %v", err)
	}
}

func TestQRGenerate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		size  int
		level string
	}{
		{
			name:  "short_text",
			input: "hello",
			size:  256,
			level: "M",
		},
		{
			name:  "url",
			input: "https://example.com",
			size:  512,
			level: "H",
		},
	}

	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := QRGenerate(tt.input, tt.size, tt.level)
			if err != nil {
				t.Fatalf("QRGenerate() error = %v", err)
			}

			if len(data) < 8 {
				t.Fatal("QRGenerate() output too short to be a PNG")
			}

			// Verify PNG header bytes.
			for i, b := range pngHeader {
				if data[i] != b {
					t.Fatalf("QRGenerate() byte %d: got 0x%02x, want 0x%02x", i, data[i], b)
				}
			}
		})
	}
}

func TestQRGenerate_AllLevels(t *testing.T) {
	levels := []string{"L", "M", "Q", "H"}
	for _, level := range levels {
		t.Run("level_"+level, func(t *testing.T) {
			data, err := QRGenerate("test", 128, level)
			if err != nil {
				t.Fatalf("QRGenerate(level=%s) error = %v", level, err)
			}
			// Higher error correction produces more data.
			if len(data) == 0 {
				t.Errorf("QRGenerate(level=%s) returned empty output", level)
			}
		})
	}
}

func TestQRInvalidLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"invalid_X", "X"},
		{"invalid_Z", "Z"},
		{"invalid_empty", ""},
		{"invalid_number", "1"},
	}

	for _, tt := range tests {
		t.Run("generate_"+tt.name, func(t *testing.T) {
			_, err := QRGenerate("hello", 256, tt.level)
			if err == nil {
				t.Error("QRGenerate() expected error for invalid level")
			}
		})

		t.Run("terminal_"+tt.name, func(t *testing.T) {
			_, err := QRTerminal("hello", tt.level)
			if err == nil {
				t.Error("QRTerminal() expected error for invalid level")
			}
		})
	}
}
