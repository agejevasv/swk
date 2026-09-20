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

	// Keep the sign aside so "-0xff" parses the same as "0xff".
	sign := ""
	if strings.HasPrefix(cleaned, "-") || strings.HasPrefix(cleaned, "+") {
		sign, cleaned = cleaned[:1], cleaned[1:]
	}

	// Strip only the prefix that matches the source base. "0b1" is a valid
	// hexadecimal literal (177), so stripping "0b" from it would be wrong.
	lower := strings.ToLower(cleaned)
	switch {
	case fromBase == 16 && strings.HasPrefix(lower, "0x"),
		fromBase == 2 && strings.HasPrefix(lower, "0b"),
		fromBase == 8 && strings.HasPrefix(lower, "0o"):
		cleaned = cleaned[2:]
	}

	n, err := strconv.ParseInt(sign+cleaned, fromBase, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse %q as base %d: %w", input, fromBase, err)
	}

	digits := strconv.FormatInt(n, toBase)

	// FormatInt already emitted the sign, so the prefix belongs after it.
	outSign := ""
	if rest, ok := strings.CutPrefix(digits, "-"); ok {
		outSign, digits = "-", rest
	}

	switch toBase {
	case 2:
		digits = "0b" + digits
	case 8:
		digits = "0o" + digits
	case 16:
		digits = "0x" + digits
	}

	return outSign + digits, nil
}
