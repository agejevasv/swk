package convert

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	humanFormat   = "Mon, 02 Jan 2006 15:04:05 MST"
	rfc2822Format = "Mon, 02 Jan 2006 15:04:05 -0700"
)

// probeTime differs from Go's reference time (Mon Jan 2 15:04:05 MST 2006) in
// every component, so formatting a run of literal text with it reveals whether
// Go read any of that text as a layout token.
var probeTime = time.Date(2021, 7, 9, 8, 23, 47, 0, time.UTC)

var strftimeMap = []struct{ directive, goLayout string }{
	{"%Y", "2006"},
	{"%m", "01"},
	{"%d", "02"},
	{"%H", "15"},
	{"%I", "03"},
	{"%M", "04"},
	{"%S", "05"},
	{"%p", "PM"},
	{"%Z", "MST"},
	{"%z", "-0700"},
	{"%A", "Monday"},
	{"%a", "Mon"},
	{"%B", "January"},
	{"%b", "Jan"},
	{"%n", "\n"},
	{"%t", "\t"},
	{"%%", "%"},
}

func strftimeLayout(directive string) (string, bool) {
	for _, s := range strftimeMap {
		if s.directive == directive {
			return s.goLayout, true
		}
	}
	return "", false
}

// formatStrftime renders each directive on its own and copies everything else
// through untouched. Translating the whole string into one Go layout would let
// literal text such as "Jan", "12" or "MST" be read as layout tokens.
func formatStrftime(t time.Time, format string) (string, error) {
	var b strings.Builder

	for i := 0; i < len(format); i++ {
		if format[i] != '%' || i+1 >= len(format) {
			b.WriteByte(format[i])
			continue
		}

		directive := format[i : i+2]
		layout, ok := strftimeLayout(directive)
		if !ok {
			return "", fmt.Errorf("unsupported strftime directive %q", directive)
		}
		b.WriteString(t.Format(layout))
		i++
	}

	return b.String(), nil
}

// strftimeToGo builds a Go layout for parsing. Parsing has to hand the whole
// layout to time.Parse, so literal text that Go would read as a token is
// rejected rather than silently misparsed.
func strftimeToGo(format string) (string, error) {
	var result, literal strings.Builder

	flush := func() error {
		lit := literal.String()
		literal.Reset()
		if lit == "" {
			return nil
		}
		if probeTime.Format(lit) != lit {
			return fmt.Errorf("literal text %q in a strftime format would be read as a time layout; "+
				"use only directives and punctuation when parsing", lit)
		}
		result.WriteString(lit)
		return nil
	}

	for i := 0; i < len(format); i++ {
		if format[i] == '%' && i+1 < len(format) {
			directive := format[i : i+2]
			layout, ok := strftimeLayout(directive)
			if !ok {
				return "", fmt.Errorf("unsupported strftime directive %q", directive)
			}
			if err := flush(); err != nil {
				return "", err
			}
			result.WriteString(layout)
			i++
			continue
		}
		literal.WriteByte(format[i])
	}

	if err := flush(); err != nil {
		return "", err
	}
	return result.String(), nil
}

func ConvertDateTime(input string, fromFmt, toFmt, tz string) (string, error) {
	var loc *time.Location
	if tz == "" || strings.EqualFold(tz, "Local") {
		loc = time.Local
	} else {
		var err error
		loc, err = time.LoadLocation(tz)
		if err != nil {
			return "", fmt.Errorf("invalid timezone %q: %w", tz, err)
		}
	}

	var t time.Time
	var err error

	if strings.EqualFold(fromFmt, "auto") {
		t, err = autoDetect(input)
	} else {
		t, err = parseFormat(input, fromFmt)
	}
	if err != nil {
		return "", err
	}

	t = t.In(loc)

	return formatTime(t, toFmt)
}

func parseFormat(input, format string) (time.Time, error) {
	switch strings.ToLower(format) {
	case "unix":
		n, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid unix timestamp: %w", err)
		}
		return time.Unix(n, 0), nil
	case "unixms":
		n, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid unix millisecond timestamp: %w", err)
		}
		return time.UnixMilli(n), nil
	case "iso":
		return time.Parse(time.RFC3339, input)
	case "rfc2822":
		return time.Parse(rfc2822Format, input)
	case "human":
		return time.Parse(humanFormat, input)
	default:
		layout := format
		if strings.Contains(format, "%") {
			var err error
			layout, err = strftimeToGo(format)
			if err != nil {
				return time.Time{}, err
			}
		} else if !isTimeLayout(format) {
			return time.Time{}, fmt.Errorf("unknown format %q: use unix, unixms, iso, rfc2822, human, auto, "+
				"a strftime format such as %%Y-%%m-%%d, or a Go layout such as 2006-01-02", format)
		}
		return time.Parse(layout, input)
	}
}

func autoDetect(input string) (time.Time, error) {
	if n, err := strconv.ParseInt(input, 10, 64); err == nil {
		if n > 1e12 {
			return time.UnixMilli(n), nil
		}
		return time.Unix(n, 0), nil
	}

	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t, nil
	}

	if t, err := time.Parse(rfc2822Format, input); err == nil {
		return t, nil
	}

	if t, err := time.Parse(humanFormat, input); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("could not auto-detect format for %q", input)
}

func formatTime(t time.Time, format string) (string, error) {
	switch strings.ToLower(format) {
	case "unix":
		return strconv.FormatInt(t.Unix(), 10), nil
	case "unixms":
		return strconv.FormatInt(t.UnixMilli(), 10), nil
	case "iso":
		return t.Format(time.RFC3339), nil
	case "rfc2822":
		return t.Format(rfc2822Format), nil
	case "human":
		return t.Format(humanFormat), nil
	default:
		if strings.Contains(format, "%") {
			return formatStrftime(t, format)
		}
		if !isTimeLayout(format) {
			return "", fmt.Errorf("unknown format %q: use unix, unixms, iso, rfc2822, human, "+
				"a strftime format such as %%Y-%%m-%%d, or a Go layout such as 2006-01-02", format)
		}
		return t.Format(format), nil
	}
}

// isTimeLayout reports whether a format carries any Go layout token. Text with
// none of them, such as "epoch", is a typo rather than a layout: Go would
// return it unchanged and the caller would never learn it was wrong.
func isTimeLayout(format string) bool {
	return probeTime.Format(format) != format
}
