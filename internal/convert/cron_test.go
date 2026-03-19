package convert

import (
	"strings"
	"testing"
	"time"
)

func TestCronExplain(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		contains []string
		wantErr  bool
	}{
		{
			name:     "every_minute",
			expr:     "* * * * *",
			contains: []string{"minute:", "every minute", "hour:", "every hour", "day of month:", "month:", "day of week:"},
		},
		{
			name:     "specific_time",
			expr:     "30 9 * * 1",
			contains: []string{"minute:", "30", "hour:", "9", "day of week:", "1"},
		},
		{
			name:     "step_values",
			expr:     "*/5 * * * *",
			contains: []string{"minute:", "*/5"},
		},
		{
			name:     "ranges",
			expr:     "0 9-17 * * *",
			contains: []string{"hour:", "9-17"},
		},
		{
			name:     "lists",
			expr:     "0 0 1,15 * *",
			contains: []string{"day of month:", "1,15"},
		},
		{
			name:     "day_names_MON_to_FRI",
			expr:     "0 9 * * MON-FRI",
			contains: []string{"day of week:", "MON-FRI"},
		},
		{
			name:     "complex_expression",
			expr:     "*/15 9-17 * 1-6 MON-FRI",
			contains: []string{"minute:", "*/15", "hour:", "9-17", "month:", "1-6", "day of week:", "MON-FRI"},
		},
		{
			name:    "invalid_expression",
			expr:    "invalid",
			wantErr: true,
		},
		{
			name:    "too_few_fields",
			expr:    "* * *",
			wantErr: true,
		},
		{
			name:    "too_many_fields",
			expr:    "* * * * * *",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CronExplain(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CronExplain() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				for _, want := range tt.contains {
					if !strings.Contains(got, want) {
						t.Errorf("CronExplain() output missing %q\ngot:\n%s", want, got)
					}
				}
			}
		})
	}
}

func TestCronExplain_OutputHasFiveLines(t *testing.T) {
	got, err := CronExplain("* * * * *")
	if err != nil {
		t.Fatalf("CronExplain: %v", err)
	}
	lines := strings.Split(got, "\n")
	if len(lines) != 5 {
		t.Errorf("CronExplain() output has %d lines, want 5\ngot:\n%s", len(lines), got)
	}
}

func TestCronNext(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		expr    string
		n       int
		wantErr bool
	}{
		{
			name: "every_minute_3_results",
			expr: "* * * * *",
			n:    3,
		},
		{
			name: "every_hour_5_results",
			expr: "0 * * * *",
			n:    5,
		},
		{
			name: "daily_at_noon",
			expr: "0 12 * * *",
			n:    3,
		},
		{
			name: "step_every_5_minutes",
			expr: "*/5 * * * *",
			n:    5,
		},
		{
			name: "single_result",
			expr: "0 0 * * *",
			n:    1,
		},
		{
			name:    "invalid_expression",
			expr:    "bad cron",
			n:       1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CronNext(tt.expr, tt.n, from)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CronNext() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(got) != tt.n {
				t.Fatalf("CronNext() returned %d results, want %d", len(got), tt.n)
			}
			// All results must be after 'from' and in ascending order.
			prev := from
			for i, tm := range got {
				if !tm.After(prev) {
					t.Errorf("CronNext() result[%d] = %v, not after %v", i, tm, prev)
				}
				prev = tm
			}
		})
	}
}

func TestCronNext_EveryMinuteCorrectTimes(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	got, err := CronNext("* * * * *", 3, from)
	if err != nil {
		t.Fatalf("CronNext: %v", err)
	}
	expected := []time.Time{
		time.Date(2024, 1, 1, 0, 1, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 0, 2, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 0, 3, 0, 0, time.UTC),
	}
	for i, want := range expected {
		if !got[i].Equal(want) {
			t.Errorf("CronNext() result[%d] = %v, want %v", i, got[i], want)
		}
	}
}

func TestCronNext_HourlyCorrectTimes(t *testing.T) {
	from := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	got, err := CronNext("0 * * * *", 2, from)
	if err != nil {
		t.Fatalf("CronNext: %v", err)
	}
	if !got[0].Equal(time.Date(2024, 6, 15, 11, 0, 0, 0, time.UTC)) {
		t.Errorf("CronNext() first = %v, want 2024-06-15 11:00 UTC", got[0])
	}
	if !got[1].Equal(time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)) {
		t.Errorf("CronNext() second = %v, want 2024-06-15 12:00 UTC", got[1])
	}
}
