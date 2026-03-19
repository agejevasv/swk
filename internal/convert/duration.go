package convert

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// DurationConvert converts between seconds and human-readable duration formats.
// The `to` parameter can be "human", "seconds", "minutes", "hours", or "" (auto).
func DurationConvert(input string, to string) (string, error) {
	input = strings.TrimSpace(input)
	to = strings.ToLower(strings.TrimSpace(to))

	if isPureDurationNumber(input) {
		// Input is seconds
		secs, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return "", fmt.Errorf("invalid number %q: %w", input, err)
		}
		totalSeconds := int64(secs)

		switch to {
		case "", "human":
			return secondsToHuman(totalSeconds), nil
		case "seconds":
			return fmt.Sprintf("%d", totalSeconds), nil
		case "minutes":
			return formatFloat(secs/60) + "m", nil
		case "hours":
			return formatFloat(secs/3600) + "h", nil
		default:
			return "", fmt.Errorf("unknown target format %q (use human, seconds, minutes, hours)", to)
		}
	}

	// Input is human readable
	totalSeconds, err := humanToSeconds(input)
	if err != nil {
		return "", err
	}

	switch to {
	case "", "seconds":
		return fmt.Sprintf("%d", totalSeconds), nil
	case "human":
		return secondsToHuman(totalSeconds), nil
	case "minutes":
		return formatFloat(float64(totalSeconds)/60) + "m", nil
	case "hours":
		return formatFloat(float64(totalSeconds)/3600) + "h", nil
	default:
		return "", fmt.Errorf("unknown target format %q (use human, seconds, minutes, hours)", to)
	}
}

func secondsToHuman(totalSeconds int64) string {
	if totalSeconds == 0 {
		return "0s"
	}

	negative := false
	if totalSeconds < 0 {
		negative = true
		totalSeconds = -totalSeconds
	}

	years := totalSeconds / (365 * 24 * 3600)
	totalSeconds %= 365 * 24 * 3600
	months := totalSeconds / (30 * 24 * 3600)
	totalSeconds %= 30 * 24 * 3600
	weeks := totalSeconds / (7 * 24 * 3600)
	totalSeconds %= 7 * 24 * 3600
	days := totalSeconds / (24 * 3600)
	totalSeconds %= 24 * 3600
	hours := totalSeconds / 3600
	totalSeconds %= 3600
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	var parts []string
	if years > 0 {
		parts = append(parts, fmt.Sprintf("%dy", years))
	}
	if months > 0 {
		parts = append(parts, fmt.Sprintf("%dmo", months))
	}
	if weeks > 0 {
		parts = append(parts, fmt.Sprintf("%dw", weeks))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	result := strings.Join(parts, " ")
	if negative {
		result = "-" + result
	}

	return result
}

var unitSeconds = map[string]float64{
	"y":  365 * 24 * 3600,
	"mo": 30 * 24 * 3600,
	"w":  7 * 24 * 3600,
	"d":  24 * 3600,
	"h":  3600,
	"m":  60,
	"s":  1,
}

func humanToSeconds(input string) (int64, error) {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", "")
	if input == "" {
		return 0, fmt.Errorf("empty duration string")
	}

	var total int64
	i := 0

	for i < len(input) {
		// Collect digits
		numStart := i
		for i < len(input) && (input[i] == '.' || (input[i] >= '0' && input[i] <= '9')) {
			i++
		}
		if numStart == i {
			return 0, fmt.Errorf("invalid duration %q: expected number at position %d", input, i)
		}
		val, err := strconv.ParseFloat(input[numStart:i], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number in duration: %q", input[numStart:i])
		}

		// Collect unit
		unitStart := i
		for i < len(input) && input[i] >= 'a' && input[i] <= 'z' {
			i++
		}
		unit := strings.ToLower(input[unitStart:i])
		if unit == "" {
			// Trailing number with no unit, treat as seconds
			total += int64(val)
			continue
		}

		multiplier, ok := unitSeconds[unit]
		if !ok {
			return 0, fmt.Errorf("unknown duration unit %q in %q", unit, input)
		}
		total += int64(val * multiplier)
	}

	return total, nil
}

func isPureDurationNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	hasDot := false
	for i, c := range s {
		if c == '-' && i == 0 {
			continue
		}
		if c == '.' {
			if hasDot {
				return false
			}
			hasDot = true
			continue
		}
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
