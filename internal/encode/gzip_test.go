package encode

import (
	"bytes"
	"compress/gzip"
	"strings"
	"testing"
)

func TestGzipRoundtrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
		level int
	}{
		{
			name:  "simple_string_default_level",
			input: "hello world",
			level: gzip.DefaultCompression,
		},
		{
			name:  "empty_input",
			input: "",
			level: gzip.DefaultCompression,
		},
		{
			name:  "best_speed",
			input: "compress me fast",
			level: gzip.BestSpeed,
		},
		{
			name:  "best_compression",
			input: "compress me well " + strings.Repeat("repetition ", 20),
			level: gzip.BestCompression,
		},
		{
			name:  "multi_line",
			input: "line1\nline2\nline3\n",
			level: gzip.DefaultCompression,
		},
		{
			name:  "level_1",
			input: "test level 1",
			level: 1,
		},
		{
			name:  "level_9",
			input: "test level 9",
			level: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := GzipCompress([]byte(tt.input), tt.level)
			if err != nil {
				t.Fatalf("GzipCompress() error = %v", err)
			}

			decompressed, err := GzipDecompress(compressed)
			if err != nil {
				t.Fatalf("GzipDecompress() error = %v", err)
			}

			if string(decompressed) != tt.input {
				t.Errorf("roundtrip: got %q, want %q", string(decompressed), tt.input)
			}
		})
	}
}

func TestGzipCompress_LargeInput(t *testing.T) {
	// Generate 1MB of repeated text.
	chunk := "The quick brown fox jumps over the lazy dog. "
	var sb strings.Builder
	for sb.Len() < 1024*1024 {
		sb.WriteString(chunk)
	}
	input := sb.String()

	compressed, err := GzipCompress([]byte(input), gzip.BestCompression)
	if err != nil {
		t.Fatalf("GzipCompress() error = %v", err)
	}

	// Compressed size should be significantly smaller than input.
	if len(compressed) >= len(input)/2 {
		t.Errorf("compressed size (%d) not significantly smaller than input (%d)", len(compressed), len(input))
	}

	decompressed, err := GzipDecompress(compressed)
	if err != nil {
		t.Fatalf("GzipDecompress() error = %v", err)
	}
	if string(decompressed) != input {
		t.Error("large input roundtrip failed: content mismatch")
	}
}

func TestGzipCompress_BinaryRoundtrip(t *testing.T) {
	// All 256 byte values.
	input := make([]byte, 256)
	for i := range input {
		input[i] = byte(i)
	}

	compressed, err := GzipCompress(input, gzip.DefaultCompression)
	if err != nil {
		t.Fatalf("GzipCompress() error = %v", err)
	}

	decompressed, err := GzipDecompress(compressed)
	if err != nil {
		t.Fatalf("GzipDecompress() error = %v", err)
	}
	if !bytes.Equal(decompressed, input) {
		t.Error("binary roundtrip failed")
	}
}

func TestGzipCompress_DifferentLevelsProduceDifferentSizes(t *testing.T) {
	input := []byte(strings.Repeat("abcdefghijklmnop", 1000))

	compressed1, err := GzipCompress(input, gzip.BestSpeed)
	if err != nil {
		t.Fatalf("BestSpeed: %v", err)
	}

	compressed9, err := GzipCompress(input, gzip.BestCompression)
	if err != nil {
		t.Fatalf("BestCompression: %v", err)
	}

	// BestCompression should produce smaller or equal output.
	if len(compressed9) > len(compressed1) {
		t.Errorf("BestCompression (%d bytes) > BestSpeed (%d bytes)", len(compressed9), len(compressed1))
	}
}

func TestGzipDecompress_InvalidData(t *testing.T) {
	_, err := GzipDecompress([]byte("not gzip data"))
	if err == nil {
		t.Error("GzipDecompress() expected error for invalid data")
	}
}

func TestGzipDecompress_EmptyBytes(t *testing.T) {
	_, err := GzipDecompress([]byte{})
	if err == nil {
		t.Error("GzipDecompress() expected error for empty input")
	}
}

func TestGzipCompress_InvalidLevel(t *testing.T) {
	_, err := GzipCompress([]byte("test"), 99)
	if err == nil {
		t.Error("GzipCompress() expected error for invalid level")
	}
}
