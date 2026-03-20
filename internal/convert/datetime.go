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

func strftimeToGo(format string) string {
	result := format
	for _, s := range strftimeMap {
		result = strings.ReplaceAll(result, s.directive, s.goLayout)
	}
	return result
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
		return time.Unix(0, n*int64(time.Millisecond)), nil
	case "iso":
		return time.Parse(time.RFC3339, input)
	case "rfc2822":
		return time.Parse(rfc2822Format, input)
	case "human":
		return time.Parse(humanFormat, input)
	default:
		layout := format
		if strings.Contains(format, "%") {
			layout = strftimeToGo(format)
		}
		return time.Parse(layout, input)
	}
}

func autoDetect(input string) (time.Time, error) {
	if n, err := strconv.ParseInt(input, 10, 64); err == nil {
		if n > 1e12 {
			return time.Unix(0, n*int64(time.Millisecond)), nil
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
		layout := format
		if strings.Contains(format, "%") {
			layout = strftimeToGo(format)
		}
		return t.Format(layout), nil
	}
}
