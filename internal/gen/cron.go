package gen

import (
	"fmt"
	"strconv"
	"strings"
)

// CronOptions configures cron expression generation.
type CronOptions struct {
	Every    string
	Daily    bool
	Weekly   bool
	Monthly  bool
	Yearly   bool
	Weekdays bool
	At       string
	Day      string
	Month    string
}

var dayNames = map[string]int{
	"SUN": 0, "MON": 1, "TUE": 2, "WED": 3,
	"THU": 4, "FRI": 5, "SAT": 6,
}

var monthNames = map[string]int{
	"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4,
	"MAY": 5, "JUN": 6, "JUL": 7, "AUG": 8,
	"SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12,
}

// GenerateCron builds a cron expression from the given options.
func GenerateCron(opts CronOptions) (string, error) {
	schedCount := countTrue(opts.Daily, opts.Weekly, opts.Monthly, opts.Yearly, opts.Weekdays) + boolFromStr(opts.Every)
	if schedCount == 0 {
		return "", fmt.Errorf("specify a schedule: --every, --daily, --weekly, --monthly, --yearly, or --weekdays")
	}
	if schedCount > 1 {
		return "", fmt.Errorf("specify only one schedule type")
	}

	if err := validateFlagCombinations(opts); err != nil {
		return "", err
	}

	if opts.Every != "" {
		return generateEvery(opts.Every)
	}

	hour, minute, err := parseAt(opts.At)
	if err != nil {
		return "", err
	}

	if opts.Daily {
		return fmt.Sprintf("%d %d * * *", minute, hour), nil
	}

	if opts.Weekdays {
		return fmt.Sprintf("%d %d * * 1-5", minute, hour), nil
	}

	if opts.Weekly {
		dow, err := parseDayOfWeek(opts.Day, 1)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d %d * * %d", minute, hour, dow), nil
	}

	if opts.Monthly {
		dom, err := parseDayOfMonth(opts.Day, 1)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d %d %d * *", minute, hour, dom), nil
	}

	// Yearly
	mon, err := parseMonth(opts.Month, 1)
	if err != nil {
		return "", err
	}
	dom, err := parseDayOfMonth(opts.Day, 1)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d %d %d %d *", minute, hour, dom, mon), nil
}

func validateFlagCombinations(opts CronOptions) error {
	if opts.Every != "" {
		if opts.At != "" {
			return fmt.Errorf("--at cannot be used with --every")
		}
		if opts.Day != "" {
			return fmt.Errorf("--day cannot be used with --every")
		}
		if opts.Month != "" {
			return fmt.Errorf("--month cannot be used with --every")
		}
	}
	if opts.Daily || opts.Weekdays {
		if opts.Day != "" {
			return fmt.Errorf("--day cannot be used with --%s", schedName(opts))
		}
		if opts.Month != "" {
			return fmt.Errorf("--month cannot be used with --%s", schedName(opts))
		}
	}
	if opts.Weekly {
		if opts.Month != "" {
			return fmt.Errorf("--month cannot be used with --weekly")
		}
	}
	return nil
}

func schedName(opts CronOptions) string {
	if opts.Daily {
		return "daily"
	}
	return "weekdays"
}

func generateEvery(s string) (string, error) {
	if len(s) < 2 {
		return "", fmt.Errorf("invalid --every %q: use format like 5m or 2h", s)
	}

	unit := s[len(s)-1]
	n, err := strconv.Atoi(s[:len(s)-1])
	if err != nil || n <= 0 {
		return "", fmt.Errorf("invalid --every %q: use format like 5m or 2h", s)
	}

	switch unit {
	case 'm':
		if 60%n != 0 {
			return "", fmt.Errorf("--every %s: minutes must divide 60 evenly (1,2,3,4,5,6,10,12,15,20,30)", s)
		}
		if n == 1 {
			return "* * * * *", nil
		}
		return fmt.Sprintf("*/%d * * * *", n), nil
	case 'h':
		if 24%n != 0 {
			return "", fmt.Errorf("--every %s: hours must divide 24 evenly (1,2,3,4,6,8,12)", s)
		}
		if n == 1 {
			return "0 * * * *", nil
		}
		return fmt.Sprintf("0 */%d * * *", n), nil
	default:
		return "", fmt.Errorf("invalid --every unit %q: use m (minutes) or h (hours)", string(unit))
	}
}

func parseAt(s string) (hour, minute int, err error) {
	if s == "" {
		return 0, 0, nil
	}
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid --at %q: use HH:MM format", s)
	}
	hour, err = strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return 0, 0, fmt.Errorf("invalid --at hour %q: must be 0-23", parts[0])
	}
	minute, err = strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute > 59 {
		return 0, 0, fmt.Errorf("invalid --at minute %q: must be 0-59", parts[1])
	}
	return hour, minute, nil
}

func parseDayOfWeek(s string, defaultVal int) (int, error) {
	if s == "" {
		return defaultVal, nil
	}
	if v, ok := dayNames[strings.ToUpper(s)]; ok {
		return v, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 || n > 6 {
		return 0, fmt.Errorf("invalid --day %q: use a day name (MON-SUN) or number (0-6)", s)
	}
	return n, nil
}

func parseDayOfMonth(s string, defaultVal int) (int, error) {
	if s == "" {
		return defaultVal, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 31 {
		return 0, fmt.Errorf("invalid --day %q: use a number 1-31", s)
	}
	return n, nil
}

func parseMonth(s string, defaultVal int) (int, error) {
	if s == "" {
		return defaultVal, nil
	}
	if v, ok := monthNames[strings.ToUpper(s)]; ok {
		return v, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 12 {
		return 0, fmt.Errorf("invalid --month %q: use a month name (JAN-DEC) or number (1-12)", s)
	}
	return n, nil
}

func countTrue(vals ...bool) int {
	n := 0
	for _, v := range vals {
		if v {
			n++
		}
	}
	return n
}

func boolFromStr(s string) int {
	if s != "" {
		return 1
	}
	return 0
}
