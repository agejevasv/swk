package convert

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

var binaryUnits = []struct {
	suffix string
	size   float64
}{
	{"PB", math.Pow(1024, 5)},
	{"TB", math.Pow(1024, 4)},
	{"GB", math.Pow(1024, 3)},
	{"MB", math.Pow(1024, 2)},
	{"KB", 1024},
	{"B", 1},
}

var decimalUnits = []struct {
	suffix string
	size   float64
}{
	{"PB", 1e15},
	{"TB", 1e12},
	{"GB", 1e9},
	{"MB", 1e6},
	{"KB", 1e3},
	{"B", 1},
}

// unitMap maps unit strings (case-insensitive lookup done separately) to byte values.
var unitMap = map[string]float64{
	"B":   1,
	"K":   1e3,
	"KB":  1e3,
	"M":   1e6,
	"MB":  1e6,
	"G":   1e9,
	"GB":  1e9,
	"T":   1e12,
	"TB":  1e12,
	"P":   1e15,
	"PB":  1e15,
	"KIB": 1024,
	"MIB": math.Pow(1024, 2),
	"GIB": math.Pow(1024, 3),
	"TIB": math.Pow(1024, 4),
	"PIB": math.Pow(1024, 5),
}

// BytesToHuman converts a byte count to a human-readable string.
func BytesToHuman(bytes int64, decimal bool) string {
	if bytes == 0 {
		return "0 B"
	}

	b := float64(bytes)
	units := binaryUnits
	if decimal {
		units = decimalUnits
	}

	for _, u := range units {
		if math.Abs(b) >= u.size {
			val := b / u.size
			formatted := formatFloat(val)
			return formatted + " " + u.suffix
		}
	}
	return fmt.Sprintf("%d B", bytes)
}

// HumanToBytes converts a human-readable size string to bytes.
func HumanToBytes(input string) (int64, error) {
	input = strings.TrimSpace(input)

	// Find where the number ends and unit begins
	i := 0
	for i < len(input) && (input[i] == '.' || input[i] == '-' || input[i] == '+' || (input[i] >= '0' && input[i] <= '9')) {
		i++
	}

	if i == 0 {
		return 0, fmt.Errorf("invalid size %q: no numeric value found", input)
	}

	numStr := strings.TrimSpace(input[:i])
	unitStr := strings.TrimSpace(input[i:])

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q: %w", numStr, err)
	}

	if unitStr == "" {
		unitStr = "B"
	}

	multiplier, ok := unitMap[strings.ToUpper(unitStr)]
	if !ok {
		return 0, fmt.Errorf("unknown unit %q", unitStr)
	}

	return int64(math.Round(num * multiplier)), nil
}

// BytesConvert auto-detects direction. If input is a plain number, converts to human.
// If input has units, converts to bytes.
func BytesConvert(input string, decimal bool) (string, error) {
	input = strings.TrimSpace(input)

	if isPureNumber(input) {
		n, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return "", fmt.Errorf("invalid number %q: %w", input, err)
		}
		return BytesToHuman(n, decimal), nil
	}

	bytes, err := HumanToBytes(input)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", bytes), nil
}

func isPureNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !unicode.IsDigit(c) && c != '-' && c != '+' {
			return false
		}
	}
	return true
}

func formatFloat(f float64) string {
	if f == math.Trunc(f) {
		return fmt.Sprintf("%.0f", f)
	}
	s := fmt.Sprintf("%.2f", f)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}
