package convert

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func CronExplain(expr string) (string, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(expr)
	if err != nil {
		return "", fmt.Errorf("invalid cron expression: %w", err)
	}

	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return "", fmt.Errorf("expected 5 fields, got %d", len(fields))
	}

	minute := explainField(fields[0], "minute", 0, 59)
	hour := explainField(fields[1], "hour", 0, 23)
	dom := explainField(fields[2], "day of month", 1, 31)
	month := explainField(fields[3], "month", 1, 12)
	dow := explainField(fields[4], "day of week", 0, 6)

	parts := []string{minute, hour, dom, month, dow}
	return strings.Join(parts, "\n"), nil
}

func explainField(field, name string, min, max int) string {
	if field == "*" {
		return fmt.Sprintf("%-15s every %s (%d-%d)", name+":", name, min, max)
	}
	return fmt.Sprintf("%-15s %s", name+":", field)
}

func CronNext(expr string, n int, from time.Time) ([]time.Time, error) {
	if n <= 0 {
		return nil, fmt.Errorf("number of runs must be positive, got %d", n)
	}

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	sched, err := parser.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("invalid cron expression: %w", err)
	}

	results := make([]time.Time, 0, n)
	t := from
	for i := 0; i < n; i++ {
		t = sched.Next(t)
		// A schedule such as "0 0 30 2 *" never fires; Next returns the zero time.
		if t.IsZero() {
			if i == 0 {
				return nil, fmt.Errorf("schedule %q never fires", expr)
			}
			break
		}
		results = append(results, t)
	}
	return results, nil
}
