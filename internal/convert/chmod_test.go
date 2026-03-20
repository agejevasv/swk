package convert

import (
	"strings"
	"testing"
)

func TestChmodToSymbolic(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// Standard 3-digit permissions.
		{
			name:  "755_standard",
			input: "755",
			want:  "rwxr-xr-x",
		},
		{
			name:  "644_standard",
			input: "644",
			want:  "rw-r--r--",
		},
		{
			name:  "777_all_permissions",
			input: "777",
			want:  "rwxrwxrwx",
		},
		{
			name:  "000_no_permissions",
			input: "000",
			want:  "---------",
		},
		{
			name:  "700_owner_only",
			input: "700",
			want:  "rwx------",
		},
		{
			name:  "600_owner_rw",
			input: "600",
			want:  "rw-------",
		},
		{
			name:  "444_read_only",
			input: "444",
			want:  "r--r--r--",
		},
		{
			name:  "111_execute_only",
			input: "111",
			want:  "--x--x--x",
		},

		// 4-digit with special bits.
		{
			name:  "4755_setuid",
			input: "4755",
			want:  "rwsr-xr-x",
		},
		{
			name:  "2755_setgid",
			input: "2755",
			want:  "rwxr-sr-x",
		},
		{
			name:  "1755_sticky",
			input: "1755",
			want:  "rwxr-xr-t",
		},
		{
			name:  "4644_setuid_no_exec",
			input: "4644",
			want:  "rwSr--r--",
		},
		{
			name:  "2644_setgid_no_exec",
			input: "2644",
			want:  "rw-r-Sr--",
		},
		{
			name:  "1644_sticky_no_exec",
			input: "1644",
			want:  "rw-r--r-T",
		},
		{
			name:  "7777_all_special_all_perms",
			input: "7777",
			want:  "rwsrwsrwt",
		},
		{
			name:  "7000_all_special_no_perms",
			input: "7000",
			want:  "--S--S--T",
		},
		{
			name:  "0755_explicit_zero_special",
			input: "0755",
			want:  "rwxr-xr-x",
		},

		// Whitespace trimming.
		{
			name:  "whitespace_trimmed",
			input: "  755  ",
			want:  "rwxr-xr-x",
		},

		// Error cases.
		{
			name:    "too_short",
			input:   "75",
			wantErr: true,
		},
		{
			name:    "too_long",
			input:   "07550",
			wantErr: true,
		},
		{
			name:    "non_digit_chars",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "out_of_range_888",
			input:   "888",
			wantErr: true,
		},
		{
			name:    "out_of_range_digit_9",
			input:   "759",
			wantErr: true,
		},
		{
			name:    "empty_string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChmodToSymbolic(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ChmodToSymbolic() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ChmodToSymbolic() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestChmodToNumeric(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// Standard permissions.
		{
			name:  "rwxr-xr-x_to_755",
			input: "rwxr-xr-x",
			want:  "755",
		},
		{
			name:  "rw-r--r--_to_644",
			input: "rw-r--r--",
			want:  "644",
		},
		{
			name:  "rwxrwxrwx_to_777",
			input: "rwxrwxrwx",
			want:  "777",
		},
		{
			name:  "no_perms_to_000",
			input: "---------",
			want:  "000",
		},
		{
			name:  "rwx------_to_700",
			input: "rwx------",
			want:  "700",
		},
		{
			name:  "r--r--r--_to_444",
			input: "r--r--r--",
			want:  "444",
		},

		// Special bits with execute (lowercase).
		{
			name:  "setuid_with_exec",
			input: "rwsr-xr-x",
			want:  "4755",
		},
		{
			name:  "setgid_with_exec",
			input: "rwxr-sr-x",
			want:  "2755",
		},
		{
			name:  "sticky_with_exec",
			input: "rwxr-xr-t",
			want:  "1755",
		},
		{
			name:  "all_special_all_perms",
			input: "rwsrwsrwt",
			want:  "7777",
		},

		// Special bits without execute (uppercase).
		{
			name:  "setuid_no_exec",
			input: "rwSr--r--",
			want:  "4644",
		},
		{
			name:  "setgid_no_exec",
			input: "rw-r-Sr--",
			want:  "2644",
		},
		{
			name:  "sticky_no_exec",
			input: "rw-r--r-T",
			want:  "1644",
		},
		{
			name:  "all_special_no_perms",
			input: "--S--S--T",
			want:  "7000",
		},

		// Whitespace trimming.
		{
			name:  "whitespace_trimmed",
			input: "  rwxr-xr-x  ",
			want:  "755",
		},

		// Error cases.
		{
			name:    "too_short",
			input:   "rwxr-x",
			wantErr: true,
		},
		{
			name:    "too_long",
			input:   "rwxr-xr-xx",
			wantErr: true,
		},
		{
			name:    "empty_string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChmodToNumeric(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ChmodToNumeric() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ChmodToNumeric() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestChmodExplain(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantContains []string
		wantErr      bool
	}{
		// Numeric input auto-detection.
		{
			name:  "numeric_755",
			input: "755",
			wantContains: []string{
				"Numeric:   755",
				"Symbolic:  rwxr-xr-x",
				"Owner:     read, write, execute",
				"Group:     read, execute",
				"Other:     read, execute",
			},
		},
		{
			name:  "numeric_644",
			input: "644",
			wantContains: []string{
				"Numeric:   644",
				"Symbolic:  rw-r--r--",
				"Owner:     read, write",
				"Group:     read",
				"Other:     read",
			},
		},
		{
			name:  "numeric_000",
			input: "000",
			wantContains: []string{
				"Numeric:   000",
				"Symbolic:  ---------",
				"Owner:     none",
				"Group:     none",
				"Other:     none",
			},
		},

		// Symbolic input auto-detection.
		{
			name:  "symbolic_rwxr-xr-x",
			input: "rwxr-xr-x",
			wantContains: []string{
				"Numeric:   755",
				"Symbolic:  rwxr-xr-x",
				"Owner:     read, write, execute",
			},
		},
		{
			name:  "symbolic_rw-r--r--",
			input: "rw-r--r--",
			wantContains: []string{
				"Numeric:   644",
				"Symbolic:  rw-r--r--",
			},
		},

		// Special bits in explanation.
		{
			name:  "setuid_4755",
			input: "4755",
			wantContains: []string{
				"Numeric:   4755",
				"Special:   setuid",
			},
		},
		{
			name:  "setgid_2755",
			input: "2755",
			wantContains: []string{
				"Numeric:   2755",
				"Special:   setgid",
			},
		},
		{
			name:  "sticky_1755",
			input: "1755",
			wantContains: []string{
				"Numeric:   1755",
				"Special:   sticky bit",
			},
		},
		{
			name:  "all_special_7777",
			input: "7777",
			wantContains: []string{
				"Special:   setuid",
				"Special:   setgid",
				"Special:   sticky bit",
			},
		},

		// Symbolic with special bits.
		{
			name:  "symbolic_setuid",
			input: "rwsr-xr-x",
			wantContains: []string{
				"Numeric:   4755",
				"Special:   setuid",
			},
		},

		// Error cases.
		{
			name:    "invalid_numeric",
			input:   "999",
			wantErr: true,
		},
		{
			name:    "invalid_symbolic_length",
			input:   "rwx",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChmodExplain(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ChmodExplain() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				for _, want := range tt.wantContains {
					if !strings.Contains(got, want) {
						t.Errorf("ChmodExplain() output missing %q\ngot:\n%s", want, got)
					}
				}
			}
		})
	}
}

func TestChmodRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		numeric string
	}{
		{"755", "755"},
		{"644", "644"},
		{"777", "777"},
		{"000", "000"},
		{"4755", "4755"},
		{"2755", "2755"},
		{"1755", "1755"},
		{"7777", "7777"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			symbolic, err := ChmodToSymbolic(tt.numeric)
			if err != nil {
				t.Fatalf("ChmodToSymbolic(%q) unexpected error: %v", tt.numeric, err)
			}

			gotNumeric, err := ChmodToNumeric(symbolic)
			if err != nil {
				t.Fatalf("ChmodToNumeric(%q) unexpected error: %v", symbolic, err)
			}

			if gotNumeric != tt.numeric {
				t.Errorf("round trip %q -> %q -> %q, want %q", tt.numeric, symbolic, gotNumeric, tt.numeric)
			}
		})
	}
}
