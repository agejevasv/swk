package convert

import (
	"math"
	"testing"
)

func TestBytesToHuman(t *testing.T) {
	tests := []struct {
		name   string
		bytes  int64
		binary bool
		want   string
	}{
		// Zero.
		{
			name:   "zero_binary_false",
			bytes:  0,
			binary: false,
			want:   "0 B",
		},
		{
			name:   "zero_binary_true",
			bytes:  0,
			binary: true,
			want:   "0 B",
		},

		// Small values (stays as bytes).
		{
			name:   "one_byte",
			bytes:  1,
			binary: false,
			want:   "1 B",
		},
		{
			name:   "small_bytes",
			bytes:  512,
			binary: false,
			want:   "512 B",
		},

		// 1024-based units (binary=false uses binaryUnits).
		{
			name:   "exactly_1KB_1024",
			bytes:  1024,
			binary: false,
			want:   "1 KB",
		},
		{
			name:   "exactly_1MB_1024",
			bytes:  1048576,
			binary: false,
			want:   "1 MB",
		},
		{
			name:   "exactly_1GB_1024",
			bytes:  1073741824,
			binary: false,
			want:   "1 GB",
		},
		{
			name:   "exactly_1TB_1024",
			bytes:  int64(math.Pow(1024, 4)),
			binary: false,
			want:   "1 TB",
		},
		{
			name:   "exactly_1PB_1024",
			bytes:  int64(math.Pow(1024, 5)),
			binary: false,
			want:   "1 PB",
		},
		{
			name:   "fractional_KB_1024",
			bytes:  1536,
			binary: false,
			want:   "1.5 KB",
		},

		// 1000-based units (binary=true uses decimalUnits).
		{
			name:   "exactly_1KB_1000",
			bytes:  1000,
			binary: true,
			want:   "1 KB",
		},
		{
			name:   "exactly_1MB_1000",
			bytes:  1000000,
			binary: true,
			want:   "1 MB",
		},
		{
			name:   "exactly_1GB_1000",
			bytes:  1000000000,
			binary: true,
			want:   "1 GB",
		},
		{
			name:   "fractional_KB_1000",
			bytes:  1500,
			binary: true,
			want:   "1.5 KB",
		},

		// Large value.
		{
			name:   "large_value_1024",
			bytes:  5368709120,
			binary: false,
			want:   "5 GB",
		},
		{
			name:   "large_value_1000",
			bytes:  5000000000,
			binary: true,
			want:   "5 GB",
		},

		// Non-round values.
		{
			name:   "non_round_MB",
			bytes:  1572864,
			binary: false,
			want:   "1.5 MB",
		},
		{
			name:   "non_round_decimal",
			bytes:  1234567890,
			binary: true,
			want:   "1.23 GB",
		},

		// Negative values.
		{
			name:   "negative_bytes",
			bytes:  -1024,
			binary: false,
			want:   "-1 KB",
		},
		{
			name:   "negative_large",
			bytes:  -1073741824,
			binary: false,
			want:   "-1 GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BytesToHuman(tt.bytes, tt.binary)
			if got != tt.want {
				t.Errorf("BytesToHuman(%d, %v) = %q, want %q", tt.bytes, tt.binary, got, tt.want)
			}
		})
	}
}

func TestHumanToBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		// Basic units (decimal/SI).
		{
			name:  "bytes_explicit",
			input: "100B",
			want:  100,
		},
		{
			name:  "kilobytes_KB",
			input: "1KB",
			want:  1000,
		},
		{
			name:  "megabytes_MB",
			input: "1MB",
			want:  1000000,
		},
		{
			name:  "gigabytes_GB",
			input: "1GB",
			want:  1000000000,
		},
		{
			name:  "terabytes_TB",
			input: "1TB",
			want:  1000000000000,
		},
		{
			name:  "petabytes_PB",
			input: "1PB",
			want:  1000000000000000,
		},

		// Short unit names.
		{
			name:  "short_K",
			input: "1K",
			want:  1000,
		},
		{
			name:  "short_M",
			input: "1M",
			want:  1000000,
		},
		{
			name:  "short_G",
			input: "1G",
			want:  1000000000,
		},
		{
			name:  "short_T",
			input: "1T",
			want:  1000000000000,
		},
		{
			name:  "short_P",
			input: "1P",
			want:  1000000000000000,
		},

		// Binary (IEC) units.
		{
			name:  "kibibytes_KiB",
			input: "1KiB",
			want:  1024,
		},
		{
			name:  "mebibytes_MiB",
			input: "1MiB",
			want:  1048576,
		},
		{
			name:  "gibibytes_GiB",
			input: "1GiB",
			want:  1073741824,
		},
		{
			name:  "tebibytes_TiB",
			input: "1TiB",
			want:  int64(math.Pow(1024, 4)),
		},
		{
			name:  "pebibytes_PiB",
			input: "1PiB",
			want:  int64(math.Pow(1024, 5)),
		},

		// Fractional values.
		{
			name:  "fractional_MB",
			input: "1.5MB",
			want:  1500000,
		},
		{
			name:  "fractional_KB",
			input: "2.5KB",
			want:  2500,
		},
		{
			name:  "fractional_GiB",
			input: "1.5GiB",
			want:  int64(math.Round(1.5 * math.Pow(1024, 3))),
		},

		// Whitespace handling.
		{
			name:  "space_between_number_and_unit",
			input: "100 MB",
			want:  100000000,
		},
		{
			name:  "leading_trailing_whitespace",
			input: "  1KB  ",
			want:  1000,
		},

		// No unit defaults to bytes.
		{
			name:  "plain_number_defaults_to_bytes",
			input: "1024",
			want:  1024,
		},
		{
			name:  "zero_no_unit",
			input: "0",
			want:  0,
		},

		// Case insensitivity.
		{
			name:  "lowercase_kb",
			input: "1kb",
			want:  1000,
		},
		{
			name:  "mixed_case_Mb",
			input: "1Mb",
			want:  1000000,
		},

		// Error cases.
		{
			name:    "empty_string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "no_numeric_value",
			input:   "KB",
			wantErr: true,
		},
		{
			name:    "unknown_unit",
			input:   "1XB",
			wantErr: true,
		},
		{
			name:    "garbage_input",
			input:   "hello",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HumanToBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("HumanToBytes(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("HumanToBytes(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestBytesConvert(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		binary  bool
		want    string
		wantErr bool
	}{
		// Pure number -> human-readable (auto-detects direction).
		{
			name:   "number_to_human_1024",
			input:  "1024",
			binary: false,
			want:   "1 KB",
		},
		{
			name:   "number_to_human_1000",
			input:  "1000",
			binary: true,
			want:   "1 KB",
		},
		{
			name:   "number_to_human_zero",
			input:  "0",
			binary: false,
			want:   "0 B",
		},
		{
			name:   "number_to_human_large",
			input:  "1073741824",
			binary: false,
			want:   "1 GB",
		},
		{
			name:   "number_to_human_small",
			input:  "42",
			binary: false,
			want:   "42 B",
		},

		// Human-readable -> bytes (auto-detects direction).
		{
			name:   "human_to_bytes_KB",
			input:  "1KB",
			binary: false,
			want:   "1000",
		},
		{
			name:   "human_to_bytes_MB",
			input:  "1MB",
			binary: false,
			want:   "1000000",
		},
		{
			name:   "human_to_bytes_GiB",
			input:  "1GiB",
			binary: false,
			want:   "1073741824",
		},
		{
			name:   "human_to_bytes_fractional",
			input:  "1.5MB",
			binary: false,
			want:   "1500000",
		},
		{
			name:   "human_to_bytes_with_space",
			input:  "100 GB",
			binary: false,
			want:   "100000000000",
		},

		// isPureNumber tested indirectly: signs make it a pure number.
		{
			name:   "negative_number_to_human",
			input:  "-1024",
			binary: false,
			want:   "-1 KB",
		},
		{
			name:   "positive_sign_number",
			input:  "+1048576",
			binary: false,
			want:   "1 MB",
		},

		// formatFloat tested indirectly: fractional formatting.
		{
			name:   "format_float_whole_number",
			input:  "1048576",
			binary: false,
			want:   "1 MB",
		},
		{
			name:   "format_float_one_decimal",
			input:  "1536",
			binary: false,
			want:   "1.5 KB",
		},
		{
			name:   "format_float_two_decimals",
			input:  "1234567890",
			binary: true,
			want:   "1.23 GB",
		},

		// Whitespace handling.
		{
			name:   "whitespace_around_number",
			input:  "  1024  ",
			binary: false,
			want:   "1 KB",
		},
		{
			name:   "whitespace_around_human",
			input:  "  1KB  ",
			binary: false,
			want:   "1000",
		},

		// Error cases.
		{
			name:    "empty_input",
			input:   "",
			binary:  false,
			wantErr: true,
		},
		{
			name:    "invalid_unit",
			input:   "5XB",
			binary:  false,
			wantErr: true,
		},
		{
			name:    "non_numeric_non_unit",
			input:   "hello",
			binary:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BytesConvert(tt.input, tt.binary)
			if (err != nil) != tt.wantErr {
				t.Fatalf("BytesConvert(%q, %v) error = %v, wantErr %v", tt.input, tt.binary, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("BytesConvert(%q, %v) = %q, want %q", tt.input, tt.binary, got, tt.want)
			}
		})
	}
}
