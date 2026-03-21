package gen

import "testing"

func TestGenerateCron(t *testing.T) {
	tests := []struct {
		name    string
		opts    CronOptions
		want    string
		wantErr bool
	}{
		// --every
		{"every 1m", CronOptions{Every: "1m"}, "* * * * *", false},
		{"every 5m", CronOptions{Every: "5m"}, "*/5 * * * *", false},
		{"every 15m", CronOptions{Every: "15m"}, "*/15 * * * *", false},
		{"every 30m", CronOptions{Every: "30m"}, "*/30 * * * *", false},
		{"every 1h", CronOptions{Every: "1h"}, "0 * * * *", false},
		{"every 2h", CronOptions{Every: "2h"}, "0 */2 * * *", false},
		{"every 6h", CronOptions{Every: "6h"}, "0 */6 * * *", false},
		{"every 12h", CronOptions{Every: "12h"}, "0 */12 * * *", false},

		// --daily
		{"daily default", CronOptions{Daily: true}, "0 0 * * *", false},
		{"daily at 9:00", CronOptions{Daily: true, At: "9:00"}, "0 9 * * *", false},
		{"daily at 14:30", CronOptions{Daily: true, At: "14:30"}, "30 14 * * *", false},

		// --weekdays
		{"weekdays default", CronOptions{Weekdays: true}, "0 0 * * 1-5", false},
		{"weekdays at 9:00", CronOptions{Weekdays: true, At: "9:00"}, "0 9 * * 1-5", false},

		// --weekly
		{"weekly default", CronOptions{Weekly: true}, "0 0 * * 1", false},
		{"weekly MON at 9:00", CronOptions{Weekly: true, Day: "MON", At: "9:00"}, "0 9 * * 1", false},
		{"weekly FRI at 17:00", CronOptions{Weekly: true, Day: "FRI", At: "17:00"}, "0 17 * * 5", false},
		{"weekly SUN", CronOptions{Weekly: true, Day: "SUN"}, "0 0 * * 0", false},
		{"weekly numeric 3", CronOptions{Weekly: true, Day: "3"}, "0 0 * * 3", false},
		{"weekly lowercase", CronOptions{Weekly: true, Day: "mon"}, "0 0 * * 1", false},

		// --monthly
		{"monthly default", CronOptions{Monthly: true}, "0 0 1 * *", false},
		{"monthly 15th", CronOptions{Monthly: true, Day: "15"}, "0 0 15 * *", false},
		{"monthly at 6:00", CronOptions{Monthly: true, Day: "1", At: "6:00"}, "0 6 1 * *", false},

		// --yearly
		{"yearly default", CronOptions{Yearly: true}, "0 0 1 1 *", false},
		{"yearly Jun 1", CronOptions{Yearly: true, Month: "6", Day: "1"}, "0 0 1 6 *", false},
		{"yearly JAN", CronOptions{Yearly: true, Month: "JAN"}, "0 0 1 1 *", false},
		{"yearly DEC 25", CronOptions{Yearly: true, Month: "DEC", Day: "25"}, "0 0 25 12 *", false},
		{"yearly lowercase", CronOptions{Yearly: true, Month: "jun"}, "0 0 1 6 *", false},

		// Errors: no schedule / multiple schedules
		{"no schedule", CronOptions{}, "", true},
		{"multiple schedules", CronOptions{Daily: true, Weekly: true}, "", true},
		{"every and daily", CronOptions{Every: "5m", Daily: true}, "", true},

		// Errors: invalid flag combinations
		{"every with at", CronOptions{Every: "5m", At: "9:00"}, "", true},
		{"every with day", CronOptions{Every: "5m", Day: "MON"}, "", true},
		{"every with month", CronOptions{Every: "5m", Month: "JAN"}, "", true},
		{"daily with day", CronOptions{Daily: true, Day: "MON"}, "", true},
		{"daily with month", CronOptions{Daily: true, Month: "JAN"}, "", true},
		{"weekdays with day", CronOptions{Weekdays: true, Day: "MON"}, "", true},
		{"weekly with month", CronOptions{Weekly: true, Month: "JAN"}, "", true},

		// Errors: invalid --every
		{"every invalid unit", CronOptions{Every: "5d"}, "", true},
		{"every invalid number", CronOptions{Every: "xm"}, "", true},
		{"every too short", CronOptions{Every: "m"}, "", true},
		{"every 7m", CronOptions{Every: "7m"}, "", true},
		{"every 5h", CronOptions{Every: "5h"}, "", true},
		{"every 0m", CronOptions{Every: "0m"}, "", true},

		// Errors: invalid --at
		{"invalid at format", CronOptions{Daily: true, At: "9"}, "", true},
		{"invalid at hour", CronOptions{Daily: true, At: "25:00"}, "", true},
		{"invalid at minute", CronOptions{Daily: true, At: "9:60"}, "", true},

		// Errors: invalid --day / --month values
		{"invalid weekly day", CronOptions{Weekly: true, Day: "FOO"}, "", true},
		{"invalid weekly day num", CronOptions{Weekly: true, Day: "7"}, "", true},
		{"invalid monthly day 0", CronOptions{Monthly: true, Day: "0"}, "", true},
		{"invalid monthly day 32", CronOptions{Monthly: true, Day: "32"}, "", true},
		{"invalid monthly day name", CronOptions{Monthly: true, Day: "MON"}, "", true},
		{"invalid month", CronOptions{Yearly: true, Month: "13"}, "", true},
		{"invalid month name", CronOptions{Yearly: true, Month: "FOO"}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateCron(tt.opts)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
