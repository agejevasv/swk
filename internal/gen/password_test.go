package gen

import (
	"strings"
	"testing"
	"unicode"
)

func TestGeneratePassword(t *testing.T) {
	tests := []struct {
		name    string
		opts    PasswordOpts
		wantErr bool
		check   func(string) bool
	}{
		{
			name: "default_length_16",
			opts: PasswordOpts{Length: 16, Upper: true, Lower: true, Digits: true, Symbols: true},
			check: func(s string) bool {
				return len(s) == 16
			},
		},
		{
			name: "length_32",
			opts: PasswordOpts{Length: 32, Upper: true, Lower: true, Digits: true, Symbols: true},
			check: func(s string) bool {
				return len(s) == 32
			},
		},
		{
			name: "no_upper_excludes_uppercase",
			opts: PasswordOpts{Length: 100, Lower: true, Digits: true, Symbols: true},
			check: func(s string) bool {
				for _, c := range s {
					if unicode.IsUpper(c) {
						return false
					}
				}
				return true
			},
		},
		{
			name: "no_lower_excludes_lowercase",
			opts: PasswordOpts{Length: 100, Upper: true, Digits: true, Symbols: true},
			check: func(s string) bool {
				for _, c := range s {
					if unicode.IsLower(c) {
						return false
					}
				}
				return true
			},
		},
		{
			name: "no_digits_excludes_digits",
			opts: PasswordOpts{Length: 100, Upper: true, Lower: true, Symbols: true},
			check: func(s string) bool {
				for _, c := range s {
					if unicode.IsDigit(c) {
						return false
					}
				}
				return true
			},
		},
		{
			name: "no_symbols_excludes_symbols",
			opts: PasswordOpts{Length: 100, Upper: true, Lower: true, Digits: true},
			check: func(s string) bool {
				for _, c := range s {
					if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
						return false
					}
				}
				return true
			},
		},
		{
			name: "only_lowercase",
			opts: PasswordOpts{Length: 50, Lower: true},
			check: func(s string) bool {
				for _, c := range s {
					if !unicode.IsLower(c) {
						return false
					}
				}
				return len(s) == 50
			},
		},
		{
			name: "only_digits",
			opts: PasswordOpts{Length: 20, Digits: true},
			check: func(s string) bool {
				for _, c := range s {
					if !unicode.IsDigit(c) {
						return false
					}
				}
				return len(s) == 20
			},
		},
		{
			name: "exclude_specific_chars",
			opts: PasswordOpts{Length: 200, Upper: true, Lower: true, Digits: true, Exclude: "aeiouAEIOU01"},
			check: func(s string) bool {
				return !strings.ContainsAny(s, "aeiouAEIOU01")
			},
		},
		{
			name: "exclude_all_but_one",
			opts: PasswordOpts{Length: 10, Lower: true, Exclude: "abcdefghijklmnopqrstuvwxy"},
			check: func(s string) bool {
				// Only 'z' should remain.
				for _, c := range s {
					if c != 'z' {
						return false
					}
				}
				return len(s) == 10
			},
		},
		{
			name: "length_1",
			opts: PasswordOpts{Length: 1, Upper: true},
			check: func(s string) bool {
				return len(s) == 1 && unicode.IsUpper(rune(s[0]))
			},
		},

		// Error cases.
		{
			name:    "zero_length_returns_error",
			opts:    PasswordOpts{Length: 0, Upper: true},
			wantErr: true,
		},
		{
			name:    "negative_length_returns_error",
			opts:    PasswordOpts{Length: -5, Upper: true},
			wantErr: true,
		},
		{
			name:    "no_character_types_returns_error",
			opts:    PasswordOpts{Length: 10},
			wantErr: true,
		},
		{
			name:    "all_chars_excluded_returns_error",
			opts:    PasswordOpts{Length: 10, Digits: true, Exclude: "0123456789"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePassword(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GeneratePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.check != nil && !tt.check(got) {
				t.Errorf("GeneratePassword() = %q, failed check", got)
			}
		})
	}
}

func TestGeneratePassword_50AllUnique(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 50; i++ {
		got, err := GeneratePassword(PasswordOpts{
			Length:  32,
			Upper:   true,
			Lower:   true,
			Digits:  true,
			Symbols: true,
		})
		if err != nil {
			t.Fatalf("GeneratePassword: %v", err)
		}
		if seen[got] {
			t.Fatalf("duplicate password found: %s", got)
		}
		seen[got] = true
	}
}
