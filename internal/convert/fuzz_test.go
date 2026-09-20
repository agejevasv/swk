package convert

import (
	"testing"
	"time"
)

// These targets assert only that a parser never panics on arbitrary input.
// Every one of them is reachable from a command line argument.

func FuzzBytesConvert(f *testing.F) {
	for _, s := range []string{"1024", "1.5GiB", "1.5GB", "0", "-1", "1e5", "", " ", "B", "1KIB", "+5", "1.2.3MB"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = BytesConvert(s, false)
		_, _ = BytesConvert(s, true)
		_, _ = HumanToBytes(s)
	})
}

func FuzzChmodExplain(f *testing.F) {
	for _, s := range []string{"755", "4755", "rwxr-xr-x", "rwsr-sr-t", "---------", "0", "77777", "", "rw"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = ChmodExplain(s)
		_, _ = ChmodToSymbolic(s)
		_, _ = ChmodToNumeric(s)
	})
}

func FuzzDurationConvert(f *testing.F) {
	for _, s := range []string{"86400", "2d 5h 30m", "1y 6mo", "-1", "0", "", "1.5h", "99999999999999999999d", "m", "5x"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		for _, to := range []string{"", "human", "seconds", "minutes", "hours", "bogus"} {
			_, _ = DurationConvert(s, to)
		}
	})
}

func FuzzConvertBase(f *testing.F) {
	f.Add("255", 10, 16)
	f.Add("0xff", 16, 10)
	f.Add("0b1", 16, 10)
	f.Add("", 2, 2)
	f.Add("-1", 10, 2)
	f.Fuzz(func(t *testing.T, s string, from, to int) {
		_, _ = ConvertBase(s, from, to)
	})
}

func FuzzConvertDateTime(f *testing.F) {
	f.Add("1700000000", "unix", "iso", "UTC")
	f.Add("2023-11-14", "%Y-%m-%d", "unix", "")
	f.Add("", "auto", "iso", "Local")
	f.Add("99999999999999999999", "auto", "human", "UTC")
	f.Fuzz(func(t *testing.T, in, from, to, tz string) {
		_, _ = ConvertDateTime(in, from, to, tz)
	})
}

func FuzzCron(f *testing.F) {
	for _, s := range []string{"* * * * *", "*/5 * * * *", "0 9 * * MON", "", "* * * *", "@daily", "60 * * * *"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = CronExplain(s)
		_, _ = CronNext(s, 3, time.Unix(0, 0))
	})
}

func FuzzJSONToCSV(f *testing.F) {
	for _, s := range []string{`[{"a":1}]`, `[]`, `{}`, ``, `[{"a":{"b":1}}]`, `[1,2]`, `[{"a":1},{"b":2}]`} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = JSONToCSV([]byte(s), ',')
		_, _ = CSVToJSON([]byte(s), ',')
	})
}

func FuzzJSONToYAML(f *testing.F) {
	for _, s := range []string{`{"a":1}`, `[]`, ``, `{"a":{"b":[1,2]}}`, `1e400`, `{"a":12345678901234567890}`} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = JSONToYAML([]byte(s))
		_, _ = YAMLToJSON([]byte(s), 2)
	})
}

func FuzzToTable(f *testing.F) {
	for _, s := range []string{`[{"a":1}]`, `[[{"a":1}]]`, `[]`, ``, `[{"a":"日本語"}]`, "a,b\n1,2"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		for _, style := range []string{"box", "simple", "plain"} {
			_, _ = ToTable([]byte(s), style, "json", ',')
			_, _ = ToTable([]byte(s), style, "csv", ',')
		}
	})
}
