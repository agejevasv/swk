package convert

import (
	"testing"
)

func TestConvertDateTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		fromFmt string
		toFmt   string
		tz      string
		want    string
		wantErr bool
	}{
		// Unix epoch (timestamp 0).
		{
			name:    "epoch_to_iso",
			input:   "0",
			fromFmt: "unix",
			toFmt:   "iso",
			tz:      "UTC",
			want:    "1970-01-01T00:00:00Z",
		},
		{
			name:    "epoch_to_human",
			input:   "0",
			fromFmt: "unix",
			toFmt:   "human",
			tz:      "UTC",
			want:    "Thu, 01 Jan 1970 00:00:00 UTC",
		},

		// Negative timestamp (pre-1970).
		{
			name:    "negative_unix_to_iso",
			input:   "-86400",
			fromFmt: "unix",
			toFmt:   "iso",
			tz:      "UTC",
			want:    "1969-12-31T00:00:00Z",
		},

		// Known timestamp.
		{
			name:    "unix_to_iso",
			input:   "1700000000",
			fromFmt: "unix",
			toFmt:   "iso",
			tz:      "UTC",
			want:    "2023-11-14T22:13:20Z",
		},
		{
			name:    "iso_to_unix",
			input:   "2023-11-14T22:13:20Z",
			fromFmt: "iso",
			toFmt:   "unix",
			tz:      "UTC",
			want:    "1700000000",
		},

		// Millisecond timestamps.
		{
			name:    "unixms_to_iso",
			input:   "1700000000000",
			fromFmt: "unixms",
			toFmt:   "iso",
			tz:      "UTC",
			want:    "2023-11-14T22:13:20Z",
		},
		{
			name:    "unix_to_unixms",
			input:   "1700000000",
			fromFmt: "unix",
			toFmt:   "unixms",
			tz:      "UTC",
			want:    "1700000000000",
		},

		// ISO 8601 with timezone offset.
		{
			name:    "iso_with_offset_to_unix",
			input:   "2023-11-14T17:13:20-05:00",
			fromFmt: "iso",
			toFmt:   "unix",
			tz:      "UTC",
			want:    "1700000000",
		},

		// RFC 2822.
		{
			name:    "rfc2822_to_unix",
			input:   "Tue, 14 Nov 2023 22:13:20 +0000",
			fromFmt: "rfc2822",
			toFmt:   "unix",
			tz:      "UTC",
			want:    "1700000000",
		},
		{
			name:    "unix_to_rfc2822",
			input:   "1700000000",
			fromFmt: "unix",
			toFmt:   "rfc2822",
			tz:      "UTC",
			want:    "Tue, 14 Nov 2023 22:13:20 +0000",
		},

		// Timezone conversion.
		{
			name:    "utc_iso_to_new_york",
			input:   "2023-11-14T22:13:20Z",
			fromFmt: "iso",
			toFmt:   "iso",
			tz:      "America/New_York",
			want:    "2023-11-14T17:13:20-05:00",
		},
		{
			name:    "utc_iso_to_tokyo",
			input:   "2023-11-14T22:13:20Z",
			fromFmt: "iso",
			toFmt:   "iso",
			tz:      "Asia/Tokyo",
			want:    "2023-11-15T07:13:20+09:00",
		},

		// Auto-detect formats.
		{
			name:    "auto_detect_unix",
			input:   "1700000000",
			fromFmt: "auto",
			toFmt:   "iso",
			tz:      "UTC",
			want:    "2023-11-14T22:13:20Z",
		},
		{
			name:    "auto_detect_unixms",
			input:   "1700000000000",
			fromFmt: "auto",
			toFmt:   "iso",
			tz:      "UTC",
			want:    "2023-11-14T22:13:20Z",
		},
		{
			name:    "auto_detect_iso",
			input:   "2023-11-14T22:13:20Z",
			fromFmt: "auto",
			toFmt:   "unix",
			tz:      "UTC",
			want:    "1700000000",
		},
		{
			name:    "auto_detect_rfc2822",
			input:   "Tue, 14 Nov 2023 22:13:20 +0000",
			fromFmt: "auto",
			toFmt:   "unix",
			tz:      "UTC",
			want:    "1700000000",
		},

		// Strftime --to format.
		{
			name:    "unix_to_strftime_date",
			input:   "1700000000",
			fromFmt: "unix",
			toFmt:   "%Y-%m-%d",
			tz:      "UTC",
			want:    "2023-11-14",
		},
		{
			name:    "unix_to_strftime_time",
			input:   "1700000000",
			fromFmt: "unix",
			toFmt:   "%H:%M:%S",
			tz:      "UTC",
			want:    "22:13:20",
		},
		{
			name:    "unix_to_strftime_full",
			input:   "1700000000",
			fromFmt: "unix",
			toFmt:   "%Y-%m-%d %H:%M:%S",
			tz:      "UTC",
			want:    "2023-11-14 22:13:20",
		},
		{
			name:    "unix_to_strftime_weekday_month",
			input:   "1700000000",
			fromFmt: "unix",
			toFmt:   "%A, %B %d",
			tz:      "UTC",
			want:    "Tuesday, November 14",
		},
		{
			name:    "unix_to_strftime_short_weekday_month",
			input:   "1700000000",
			fromFmt: "unix",
			toFmt:   "%a, %b %d",
			tz:      "UTC",
			want:    "Tue, Nov 14",
		},
		{
			name:    "unix_to_strftime_12h",
			input:   "1700000000",
			fromFmt: "unix",
			toFmt:   "%I:%M %p",
			tz:      "UTC",
			want:    "10:13 PM",
		},
		{
			name:    "unix_to_strftime_timezone",
			input:   "1700000000",
			fromFmt: "unix",
			toFmt:   "%Y-%m-%d %Z",
			tz:      "UTC",
			want:    "2023-11-14 UTC",
		},
		{
			name:    "unix_to_strftime_literal_percent",
			input:   "1700000000",
			fromFmt: "unix",
			toFmt:   "%Y%%%m",
			tz:      "UTC",
			want:    "2023%11",
		},

		// Strftime --from format.
		{
			name:    "strftime_date_to_iso",
			input:   "2023-11-14",
			fromFmt: "%Y-%m-%d",
			toFmt:   "iso",
			tz:      "UTC",
			want:    "2023-11-14T00:00:00Z",
		},
		{
			name:    "strftime_datetime_to_unix",
			input:   "2023-11-14 22:13:20",
			fromFmt: "%Y-%m-%d %H:%M:%S",
			toFmt:   "unix",
			tz:      "UTC",
			want:    "1700000000",
		},
		{
			name:    "strftime_from_mismatch",
			input:   "not-a-date",
			fromFmt: "%Y-%m-%d",
			toFmt:   "iso",
			tz:      "UTC",
			wantErr: true,
		},

		// Strftime roundtrip.
		{
			name:    "strftime_roundtrip",
			input:   "2023-11-14",
			fromFmt: "%Y-%m-%d",
			toFmt:   "%Y-%m-%d",
			tz:      "UTC",
			want:    "2023-11-14",
		},

		// Go layouts still work as fallback.
		{
			name:    "go_layout_still_works",
			input:   "2023-11-14",
			fromFmt: "2006-01-02",
			toFmt:   "iso",
			tz:      "UTC",
			want:    "2023-11-14T00:00:00Z",
		},

		// Local timezone (empty string).
		{
			name:    "empty_tz_uses_local",
			input:   "0",
			fromFmt: "unix",
			toFmt:   "unix",
			tz:      "",
			want:    "0",
		},

		// Error cases.
		{
			name:    "invalid_unix_timestamp",
			input:   "not-a-number",
			fromFmt: "unix",
			toFmt:   "iso",
			tz:      "UTC",
			wantErr: true,
		},
		{
			name:    "invalid_iso_format",
			input:   "2023-13-45",
			fromFmt: "iso",
			toFmt:   "unix",
			tz:      "UTC",
			wantErr: true,
		},
		{
			name:    "invalid_timezone",
			input:   "1700000000",
			fromFmt: "unix",
			toFmt:   "iso",
			tz:      "Invalid/Zone",
			wantErr: true,
		},
		// Go layout mismatch error.
		{
			name:    "go_layout_from_mismatch",
			input:   "not-a-date",
			fromFmt: "2006-01-02",
			toFmt:   "iso",
			tz:      "UTC",
			wantErr: true,
		},
		{
			name:    "auto_detect_fails_on_garbage",
			input:   "not-any-format",
			fromFmt: "auto",
			toFmt:   "iso",
			tz:      "UTC",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertDateTime(tt.input, tt.fromFmt, tt.toFmt, tt.tz)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ConvertDateTime() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ConvertDateTime() = %q, want %q", got, tt.want)
			}
		})
	}
}
