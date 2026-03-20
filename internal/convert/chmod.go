package convert

import (
	"fmt"
	"strconv"
	"strings"
)

// ChmodToSymbolic converts numeric permission like "755" or "4755" to symbolic like "rwxr-xr-x".
func ChmodToSymbolic(numeric string) (string, error) {
	numeric = strings.TrimSpace(numeric)

	if len(numeric) < 3 || len(numeric) > 4 {
		return "", fmt.Errorf("invalid numeric permission %q: expected 3 or 4 digits", numeric)
	}

	for _, c := range numeric {
		if c < '0' || c > '7' {
			return "", fmt.Errorf("invalid numeric permission %q: digits must be 0-7", numeric)
		}
	}

	// Pad to 4 digits
	padded := numeric
	if len(padded) == 3 {
		padded = "0" + padded
	}

	special, _ := strconv.Atoi(string(padded[0]))
	owner, _ := strconv.Atoi(string(padded[1]))
	group, _ := strconv.Atoi(string(padded[2]))
	other, _ := strconv.Atoi(string(padded[3]))

	var sb strings.Builder
	// Owner
	writePermTriple(&sb, owner, special&4 != 0, false, 's', 'S')
	// Group
	writePermTriple(&sb, group, false, special&2 != 0, 's', 'S')
	// Other
	writePermTriple(&sb, other, false, special&1 != 0, 't', 'T')

	return sb.String(), nil
}

func writePermTriple(sb *strings.Builder, perm int, setuid, sticky bool, setChar, setCharUpper byte) {
	if perm&4 != 0 {
		sb.WriteByte('r')
	} else {
		sb.WriteByte('-')
	}
	if perm&2 != 0 {
		sb.WriteByte('w')
	} else {
		sb.WriteByte('-')
	}

	hasExec := perm&1 != 0
	if setuid || sticky {
		if hasExec {
			sb.WriteByte(setChar)
		} else {
			sb.WriteByte(setCharUpper)
		}
	} else {
		if hasExec {
			sb.WriteByte('x')
		} else {
			sb.WriteByte('-')
		}
	}
}

// ChmodToNumeric converts symbolic permission like "rwxr-xr-x" to numeric like "755".
func ChmodToNumeric(symbolic string) (string, error) {
	symbolic = strings.TrimSpace(symbolic)

	if len(symbolic) != 9 {
		return "", fmt.Errorf("invalid symbolic permission %q: expected 9 characters", symbolic)
	}

	special := 0
	owner := parseTriple(symbolic[0:3], &special, 4)
	group := parseTriple(symbolic[3:6], &special, 2)
	other := parseTriple(symbolic[6:9], &special, 1)

	if special != 0 {
		return fmt.Sprintf("%d%d%d%d", special, owner, group, other), nil
	}
	return fmt.Sprintf("%d%d%d", owner, group, other), nil
}

func parseTriple(s string, special *int, specialBit int) int {
	val := 0
	if s[0] == 'r' {
		val |= 4
	}
	if s[1] == 'w' {
		val |= 2
	}
	switch s[2] {
	case 'x':
		val |= 1
	case 's':
		val |= 1
		*special |= specialBit
	case 'S':
		*special |= specialBit
	case 't':
		val |= 1
		*special |= specialBit
	case 'T':
		*special |= specialBit
	}
	return val
}

// ChmodExplain auto-detects the input format and returns a detailed explanation.
func ChmodExplain(input string) (string, error) {
	input = strings.TrimSpace(input)

	var numeric, symbolic string
	var err error

	if isNumericChmod(input) {
		numeric = input
		symbolic, err = ChmodToSymbolic(input)
		if err != nil {
			return "", err
		}
	} else {
		symbolic = input
		numeric, err = ChmodToNumeric(input)
		if err != nil {
			return "", err
		}
	}

	// Pad numeric for parsing special bits
	padded := numeric
	if len(padded) == 3 {
		padded = "0" + padded
	}

	owner := describePerms(symbolic[0:3])
	group := describePerms(symbolic[3:6])
	other := describePerms(symbolic[6:9])

	special, _ := strconv.Atoi(string(padded[0]))
	var lines []string
	lines = append(lines, fmt.Sprintf("Numeric:   %s", numeric))
	lines = append(lines, fmt.Sprintf("Symbolic:  %s", symbolic))
	lines = append(lines, fmt.Sprintf("Owner:     %s", owner))
	lines = append(lines, fmt.Sprintf("Group:     %s", group))
	lines = append(lines, fmt.Sprintf("Other:     %s", other))

	if special&4 != 0 {
		lines = append(lines, "Special:   setuid")
	}
	if special&2 != 0 {
		lines = append(lines, "Special:   setgid")
	}
	if special&1 != 0 {
		lines = append(lines, "Special:   sticky bit")
	}

	return strings.Join(lines, "\n"), nil
}

func describePerms(triple string) string {
	var parts []string
	if triple[0] == 'r' {
		parts = append(parts, "read")
	}
	if triple[1] == 'w' {
		parts = append(parts, "write")
	}
	switch triple[2] {
	case 'x', 's', 't':
		parts = append(parts, "execute")
	}
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, ", ")
}

func isNumericChmod(s string) bool {
	if len(s) < 3 || len(s) > 4 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '7' {
			return false
		}
	}
	return true
}
