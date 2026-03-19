package convert

import (
	"testing"
)

func TestConvertBase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		fromBase int
		toBase   int
		want     string
		wantErr  bool
	}{
		// Basic conversions between all base combinations.
		{
			name:     "dec_to_hex",
			input:    "255",
			fromBase: 10,
			toBase:   16,
			want:     "0xff",
		},
		{
			name:     "hex_to_dec",
			input:    "ff",
			fromBase: 16,
			toBase:   10,
			want:     "255",
		},
		{
			name:     "dec_to_bin",
			input:    "10",
			fromBase: 10,
			toBase:   2,
			want:     "0b1010",
		},
		{
			name:     "bin_to_dec",
			input:    "1010",
			fromBase: 2,
			toBase:   10,
			want:     "10",
		},
		{
			name:     "dec_to_oct",
			input:    "8",
			fromBase: 10,
			toBase:   8,
			want:     "0o10",
		},
		{
			name:     "oct_to_dec",
			input:    "10",
			fromBase: 8,
			toBase:   10,
			want:     "8",
		},
		{
			name:     "bin_to_hex",
			input:    "11111111",
			fromBase: 2,
			toBase:   16,
			want:     "0xff",
		},
		{
			name:     "hex_to_bin",
			input:    "a",
			fromBase: 16,
			toBase:   2,
			want:     "0b1010",
		},
		{
			name:     "oct_to_hex",
			input:    "77",
			fromBase: 8,
			toBase:   16,
			want:     "0x3f",
		},
		{
			name:     "hex_to_oct",
			input:    "3f",
			fromBase: 16,
			toBase:   8,
			want:     "0o77",
		},
		{
			name:     "bin_to_oct",
			input:    "111",
			fromBase: 2,
			toBase:   8,
			want:     "0o7",
		},
		{
			name:     "oct_to_bin",
			input:    "7",
			fromBase: 8,
			toBase:   2,
			want:     "0b111",
		},

		// Prefix stripping.
		{
			name:     "0x_prefix_hex_to_dec",
			input:    "0xFF",
			fromBase: 16,
			toBase:   10,
			want:     "255",
		},
		{
			name:     "0b_prefix_bin_to_dec",
			input:    "0b1010",
			fromBase: 2,
			toBase:   10,
			want:     "10",
		},
		{
			name:     "0o_prefix_oct_to_dec",
			input:    "0o10",
			fromBase: 8,
			toBase:   10,
			want:     "8",
		},

		// Prefix stripping (input has 0x but fromBase is already specified).
		{
			name:     "hex_with_prefix",
			input:    "0xff",
			fromBase: 16,
			toBase:   10,
			want:     "255",
		},
		{
			name:     "bin_with_prefix",
			input:    "0b1010",
			fromBase: 2,
			toBase:   10,
			want:     "10",
		},
		{
			name:     "oct_with_prefix",
			input:    "0o10",
			fromBase: 8,
			toBase:   10,
			want:     "8",
		},

		// Zero in all bases.
		{
			name:     "zero_dec_to_hex",
			input:    "0",
			fromBase: 10,
			toBase:   16,
			want:     "0x0",
		},
		{
			name:     "zero_dec_to_bin",
			input:    "0",
			fromBase: 10,
			toBase:   2,
			want:     "0b0",
		},
		{
			name:     "zero_dec_to_oct",
			input:    "0",
			fromBase: 10,
			toBase:   8,
			want:     "0o0",
		},
		{
			name:     "zero_hex_to_dec",
			input:    "0",
			fromBase: 16,
			toBase:   10,
			want:     "0",
		},

		// Large number (max int64 area).
		{
			name:     "large_number_dec_to_hex",
			input:    "9223372036854775807",
			fromBase: 10,
			toBase:   16,
			want:     "0x7fffffffffffffff",
		},
		{
			name:     "large_number_hex_to_dec",
			input:    "7fffffffffffffff",
			fromBase: 16,
			toBase:   10,
			want:     "9223372036854775807",
		},

		// Negative numbers.
		{
			name:     "negative_dec_to_hex",
			input:    "-1",
			fromBase: 10,
			toBase:   16,
			want:     "0x-1",
		},
		{
			name:     "negative_dec_to_bin",
			input:    "-10",
			fromBase: 10,
			toBase:   2,
			want:     "0b-1010",
		},

		// Error cases.
		{
			name:     "invalid_characters_for_base",
			input:    "xyz",
			fromBase: 10,
			toBase:   16,
			wantErr:  true,
		},
		{
			name:     "invalid_binary_digit",
			input:    "2",
			fromBase: 2,
			toBase:   10,
			wantErr:  true,
		},
		{
			name:     "unsupported_from_base",
			input:    "10",
			fromBase: 99,
			toBase:   10,
			wantErr:  true,
		},
		{
			name:     "unsupported_to_base",
			input:    "10",
			fromBase: 10,
			toBase:   99,
			wantErr:  true,
		},
		{
			name:     "from_base_1_unsupported",
			input:    "1",
			fromBase: 1,
			toBase:   10,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertBase(tt.input, tt.fromBase, tt.toBase)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ConvertBase() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ConvertBase() = %q, want %q", got, tt.want)
			}
		})
	}
}
