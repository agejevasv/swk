package convert

import (
	"fmt"
	"strconv"
	"strings"
)

func ConvertBase(input string, fromBase, toBase int) (string, error) {
	if fromBase < 2 || fromBase > 16 {
		return "", fmt.Errorf("unsupported from-base: %d", fromBase)
	}
	if toBase < 2 || toBase > 16 {
		return "", fmt.Errorf("unsupported to-base: %d", toBase)
	}

	cleaned := strings.TrimSpace(input)

	// Strip only the prefix that matches the source base. "0b1" is a valid
	// hexadecimal literal (177), so stripping "0b" from it would be wrong.
	lower := strings.ToLower(cleaned)
	switch {
	case fromBase == 16 && strings.HasPrefix(lower, "0x"),
		fromBase == 2 && strings.HasPrefix(lower, "0b"),
		fromBase == 8 && strings.HasPrefix(lower, "0o"):
		cleaned = cleaned[2:]
	}

	n, err := strconv.ParseInt(cleaned, fromBase, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse %q as base %d: %w", input, fromBase, err)
	}

	result := strconv.FormatInt(n, toBase)

	switch toBase {
	case 2:
		result = "0b" + result
	case 8:
		result = "0o" + result
	case 16:
		result = "0x" + result
	}

	return result, nil
}
