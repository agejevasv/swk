package gen

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars  = "abcdefghijklmnopqrstuvwxyz"
	digitChars  = "0123456789"
	symbolChars = "!@#$%^&*()-_=+[]{}|;:',.<>?/`~"
)

// PasswordOpts configures password generation.
type PasswordOpts struct {
	Length  int
	Upper   bool
	Lower   bool
	Digits  bool
	Symbols bool
	Exclude string
}

func GeneratePassword(opts PasswordOpts) (string, error) {
	if opts.Length <= 0 {
		return "", fmt.Errorf("password length must be positive")
	}

	var charset strings.Builder
	if opts.Upper {
		charset.WriteString(upperChars)
	}
	if opts.Lower {
		charset.WriteString(lowerChars)
	}
	if opts.Digits {
		charset.WriteString(digitChars)
	}
	if opts.Symbols {
		charset.WriteString(symbolChars)
	}

	chars := charset.String()
	if opts.Exclude != "" {
		var filtered strings.Builder
		for _, c := range chars {
			if !strings.ContainsRune(opts.Exclude, c) {
				filtered.WriteRune(c)
			}
		}
		chars = filtered.String()
	}

	if len(chars) == 0 {
		return "", fmt.Errorf("no characters available for password generation")
	}

	password := make([]byte, opts.Length)
	for i := range password {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		password[i] = chars[idx.Int64()]
	}

	return string(password), nil
}
