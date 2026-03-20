package convert

import (
	"testing"
)

func TestDurationConvert(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		to      string
		want    string
		wantErr bool
	}{
		// Pure number inputs (seconds → human), auto mode.
		{
			name:  "60_seconds_to_human_auto",
			input: "60",
			to:    "",
			want:  "1m",
		},
		{
			name:  "3600_seconds_to_human_auto",
			input: "3600",
			to:    "",
			want:  "1h",
		},
		{
			name:  "90061_seconds_to_human_auto",
			input: "90061",
			to:    "",
			want:  "1d 1h 1m 1s",
		},
		{
			name:  "zero_seconds_to_human",
			input: "0",
			to:    "",
			want:  "0s",
		},
		{
			name:  "negative_seconds_to_human",
			input: "-3661",
			to:    "",
			want:  "-1h 1m 1s",
		},
		{
			name:  "1_second",
			input: "1",
			to:    "",
			want:  "1s",
		},
		{
			name:  "seconds_exactly_one_day",
			input: "86400",
			to:    "",
			want:  "1d",
		},
		{
			name:  "seconds_exactly_one_week",
			input: "604800",
			to:    "",
			want:  "1w",
		},
		{
			name:  "seconds_exactly_one_year",
			input: "31536000",
			to:    "",
			want:  "1y",
		},
		{
			name:  "seconds_one_month",
			input: "2592000",
			to:    "",
			want:  "1mo",
		},

		// Pure number inputs with explicit `to` = "human".
		{
			name:  "60_seconds_to_human_explicit",
			input: "60",
			to:    "human",
			want:  "1m",
		},
		{
			name:  "3600_seconds_to_human_explicit",
			input: "3600",
			to:    "human",
			want:  "1h",
		},

		// Pure number inputs with `to` = "seconds" (identity).
		{
			name:  "seconds_to_seconds",
			input: "3600",
			to:    "seconds",
			want:  "3600",
		},
		{
			name:  "float_seconds_to_seconds_truncates",
			input: "3600.7",
			to:    "seconds",
			want:  "3600",
		},

		// Pure number inputs with `to` = "minutes".
		{
			name:  "seconds_to_minutes",
			input: "120",
			to:    "minutes",
			want:  "2m",
		},
		{
			name:  "seconds_to_minutes_fractional",
			input: "90",
			to:    "minutes",
			want:  "1.5m",
		},
		{
			name:  "seconds_to_minutes_zero",
			input: "0",
			to:    "minutes",
			want:  "0m",
		},

		// Pure number inputs with `to` = "hours".
		{
			name:  "seconds_to_hours",
			input: "7200",
			to:    "hours",
			want:  "2h",
		},
		{
			name:  "seconds_to_hours_fractional",
			input: "5400",
			to:    "hours",
			want:  "1.5h",
		},

		// Human-readable inputs → seconds (auto mode).
		{
			name:  "1h30m_to_seconds",
			input: "1h30m",
			to:    "",
			want:  "5400",
		},
		{
			name:  "1d_to_seconds",
			input: "1d",
			to:    "",
			want:  "86400",
		},
		{
			name:  "1w_to_seconds",
			input: "1w",
			to:    "",
			want:  "604800",
		},
		{
			name:  "1y_to_seconds",
			input: "1y",
			to:    "",
			want:  "31536000",
		},
		{
			name:  "1mo_to_seconds",
			input: "1mo",
			to:    "",
			want:  "2592000",
		},
		{
			name:  "combined_duration_to_seconds",
			input: "1h30m45s",
			to:    "",
			want:  "5445",
		},
		{
			name:  "complex_combined_duration",
			input: "1d2h3m4s",
			to:    "",
			want:  "93784",
		},
		{
			name:  "human_with_spaces",
			input: "1h 30m",
			to:    "",
			want:  "5400",
		},

		// Human-readable inputs with explicit `to` = "seconds".
		{
			name:  "human_to_seconds_explicit",
			input: "2h",
			to:    "seconds",
			want:  "7200",
		},

		// Human-readable inputs with `to` = "human" (round-trip).
		{
			name:  "human_to_human",
			input: "1h30m",
			to:    "human",
			want:  "1h 30m",
		},

		// Human-readable inputs to minutes.
		{
			name:  "human_to_minutes",
			input: "2h",
			to:    "minutes",
			want:  "120m",
		},
		{
			name:  "human_to_minutes_fractional",
			input: "1h30s",
			to:    "minutes",
			want:  "60.5m",
		},

		// Human-readable inputs to hours.
		{
			name:  "human_to_hours",
			input: "2d",
			to:    "hours",
			want:  "48h",
		},
		{
			name:  "human_to_hours_fractional",
			input: "1h30m",
			to:    "hours",
			want:  "1.5h",
		},

		// Edge cases.
		{
			name:    "empty_input",
			input:   "",
			to:      "",
			wantErr: true,
		},
		{
			name:    "invalid_format",
			input:   "abc",
			to:      "",
			wantErr: true,
		},
		{
			name:    "unknown_unit",
			input:   "5x",
			to:      "",
			wantErr: true,
		},
		{
			name:    "unknown_to_format_number",
			input:   "100",
			to:      "banana",
			wantErr: true,
		},
		{
			name:    "unknown_to_format_human",
			input:   "1h",
			to:      "banana",
			wantErr: true,
		},
		{
			name:  "very_large_number",
			input: "100000000",
			to:    "human",
			want:  "3y 2mo 2d 9h 46m 40s",
		},
		{
			name:  "whitespace_padding",
			input: "  3600  ",
			to:    "  human  ",
			want:  "1h",
		},
		{
			name:  "just_seconds_unit",
			input: "30s",
			to:    "",
			want:  "30",
		},
		{
			name:  "just_minutes_unit",
			input: "5m",
			to:    "",
			want:  "300",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DurationConvert(tt.input, tt.to)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DurationConvert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("DurationConvert() = %q, want %q", got, tt.want)
			}
		})
	}
}
