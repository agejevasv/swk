package convert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func resetAllFlags() {
	// datetime.go package-level vars
	dtFrom = "auto"
	dtTo = "iso"
	dtTz = "Local"

	// cron.go package-level vars
	cronNext = 5
	cronExplain = false

	// json_yaml.go package-level vars
	jyReverse = false
	jyIndent = 2

	// json_csv.go package-level vars
	jcReverse = false
	jcDelimiter = ","

	// numbase.go package-level vars
	nbFrom = "dec"
	nbTo = "hex"

	// Reset all cobra subcommand flags to defaults and clear Changed state
	for _, sub := range Cmd.Commands() {
		sub.Flags().VisitAll(func(f *pflag.Flag) {
			f.Value.Set(f.DefValue)
			f.Changed = false
		})
	}
}

func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetErr(buf)
	Cmd.SetArgs(args)
	err := Cmd.Execute()
	return buf.String(), err
}

// ── json-yaml ──────────────────────────────────────────────────────────────

func TestJsonYaml_JSONToYAML(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json-yaml", `{"name":"alice","age":30}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "name: alice") {
		t.Errorf("expected YAML key 'name: alice', got:\n%s", out)
	}
	if !strings.Contains(out, "age: 30") {
		t.Errorf("expected YAML key 'age: 30', got:\n%s", out)
	}
}

func TestJsonYaml_YAMLToJSON_Reverse(t *testing.T) {
	t.Cleanup(resetAllFlags)
	yamlInput := "name: alice\nage: 30"
	out, err := executeCommand("json-yaml", "--reverse", yamlInput)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"name"`) {
		t.Errorf("expected JSON key '\"name\"', got:\n%s", out)
	}
	if !strings.Contains(out, `"age"`) {
		t.Errorf("expected JSON key '\"age\"', got:\n%s", out)
	}
}

func TestJsonYaml_IndentFlag(t *testing.T) {
	t.Cleanup(resetAllFlags)
	yamlInput := "name: alice"
	out, err := executeCommand("json-yaml", "--reverse", "--indent", "4", yamlInput)
	if err != nil {
		t.Fatal(err)
	}
	// With indent=4, JSON should have 4-space indentation
	if !strings.Contains(out, "    ") {
		t.Errorf("expected 4-space indentation, got:\n%s", out)
	}
}

func TestJsonYaml_InvalidJSON(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("json-yaml", `{invalid json!!!`)
	if err == nil {
		t.Error("expected error for invalid JSON input")
	}
}

// ── json-csv ───────────────────────────────────────────────────────────────

func TestJsonCSV_JSONToCSV(t *testing.T) {
	t.Cleanup(resetAllFlags)
	input := `[{"name":"alice","age":30},{"name":"bob","age":25}]`
	out, err := executeCommand("json-csv", input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "alice") || !strings.Contains(out, "bob") {
		t.Errorf("expected CSV output with alice and bob, got:\n%s", out)
	}
	// Should have header row and data rows
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines (header + 2 data), got %d", len(lines))
	}
}

func TestJsonCSV_CSVToJSON_Reverse(t *testing.T) {
	t.Cleanup(resetAllFlags)
	csvInput := "name,age\nalice,30\nbob,25"
	out, err := executeCommand("json-csv", "--reverse", csvInput)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"name"`) {
		t.Errorf("expected JSON with 'name' key, got:\n%s", out)
	}
	if !strings.Contains(out, `"alice"`) {
		t.Errorf("expected JSON with 'alice' value, got:\n%s", out)
	}
}

func TestJsonCSV_DelimiterFlag(t *testing.T) {
	t.Cleanup(resetAllFlags)
	// First convert JSON to semicolon-delimited CSV
	input := `[{"name":"alice","age":"30"}]`
	out, err := executeCommand("json-csv", "-d", ";", input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, ";") {
		t.Errorf("expected semicolon-delimited output, got:\n%s", out)
	}
}

func TestJsonCSV_DelimiterFlag_Reverse(t *testing.T) {
	t.Cleanup(resetAllFlags)
	csvInput := "name;age\nalice;30"
	out, err := executeCommand("json-csv", "--reverse", "--delimiter", ";", csvInput)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"alice"`) {
		t.Errorf("expected JSON with alice, got:\n%s", out)
	}
}

// ── numbase ────────────────────────────────────────────────────────────────

func TestNumbase_DecToHex_Default(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("numbase", "255")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if !strings.EqualFold(trimmed, "ff") && !strings.EqualFold(trimmed, "0xff") {
		t.Errorf("expected 'ff' or '0xff', got %q", trimmed)
	}
}

func TestNumbase_BinToDec(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("numbase", "--from", "bin", "--to", "dec", "11111111")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if trimmed != "255" {
		t.Errorf("expected '255', got %q", trimmed)
	}
}

func TestNumbase_HexToOct(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("numbase", "--from", "hex", "--to", "oct", "ff")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if trimmed != "377" && trimmed != "0o377" {
		t.Errorf("expected '377' or '0o377', got %q", trimmed)
	}
}

func TestNumbase_DecToBin(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("numbase", "--from", "dec", "--to", "bin", "10")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if trimmed != "1010" && trimmed != "0b1010" {
		t.Errorf("expected '1010' or '0b1010', got %q", trimmed)
	}
}

func TestNumbase_OctToHex(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("numbase", "--from", "oct", "--to", "hex", "377")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if !strings.EqualFold(trimmed, "ff") && !strings.EqualFold(trimmed, "0xff") {
		t.Errorf("expected 'ff' or '0xff', got %q", trimmed)
	}
}

func TestNumbase_InvalidBase(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("numbase", "--from", "xyz", "42")
	if err == nil {
		t.Error("expected error for invalid base name")
	}
}

// ── datetime ───────────────────────────────────────────────────────────────

func TestDatetime_UnixToISO(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("datetime", "--from", "unix", "--tz", "UTC", "0")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if !strings.Contains(trimmed, "1970-01-01") {
		t.Errorf("expected ISO date containing 1970-01-01, got %q", trimmed)
	}
}

func TestDatetime_ISOToUnix(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("datetime", "--from", "iso", "--to", "unix", "1970-01-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if trimmed != "0" {
		t.Errorf("expected '0', got %q", trimmed)
	}
}

func TestDatetime_Now(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("datetime", "--tz", "UTC", "now")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		t.Error("expected non-empty output for 'now'")
	}
	// The output should be an ISO date string
	if !strings.Contains(trimmed, "T") {
		t.Errorf("expected ISO format with 'T' separator, got %q", trimmed)
	}
}

func TestDatetime_TzFlag(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("datetime", "--from", "unix", "--tz", "America/New_York", "0")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	// Unix epoch in New York is 1969-12-31T19:00:00-05:00
	if !strings.Contains(trimmed, "1969-12-31") {
		t.Errorf("expected date 1969-12-31 in New York timezone, got %q", trimmed)
	}
}

func TestDatetime_AutoDetect(t *testing.T) {
	t.Cleanup(resetAllFlags)
	// Default --from is "auto"; passing a unix timestamp should auto-detect
	out, err := executeCommand("datetime", "--to", "iso", "--tz", "UTC", "0")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if !strings.Contains(trimmed, "1970") {
		t.Errorf("expected year 1970 in output, got %q", trimmed)
	}
}

// ── cron ───────────────────────────────────────────────────────────────────

func TestCron_DefaultShowsBoth(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("cron", "*/5 * * * *")
	if err != nil {
		t.Fatal(err)
	}
	// Default shows explanation and next 5 runs
	if !strings.Contains(out, "Next 5 runs") {
		t.Errorf("expected 'Next 5 runs' in output, got:\n%s", out)
	}
	// Should contain some explanation text
	if len(out) < 20 {
		t.Errorf("expected substantial output with both explain + next, got:\n%s", out)
	}
}

func TestCron_ExplainOnly(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("cron", "--explain", "0 9 * * 1-5")
	if err != nil {
		t.Fatal(err)
	}
	// With --explain only (and no --next), should show explanation but not "Next N runs"
	if strings.Contains(out, "Next") {
		t.Errorf("expected explain-only output without 'Next', got:\n%s", out)
	}
	if out == "" {
		t.Error("expected non-empty explanation")
	}
}

func TestCron_NextOnly(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("cron", "--next", "3", "*/10 * * * *")
	if err != nil {
		t.Fatal(err)
	}
	// With --next 3 only (no --explain), should show timestamps
	// The output may include a "Next N runs:" header line
	if !strings.Contains(out, "T") {
		t.Errorf("expected RFC3339 timestamps with 'T', got:\n%s", out)
	}
	// Count lines containing timestamps (with 'T' separator from RFC3339)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, "T") {
			count++
		}
	}
	if count != 3 {
		t.Errorf("expected 3 timestamp lines, got %d:\n%s", count, out)
	}
}

func TestCron_NextWithExplain(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("cron", "--next", "2", "--explain", "0 0 * * *")
	if err != nil {
		t.Fatal(err)
	}
	// When both --explain and --next are given, shows both
	if !strings.Contains(out, "T") {
		t.Errorf("expected timestamp output with 'T' in RFC3339, got:\n%s", out)
	}
}

func TestCron_InvalidExpression(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("cron", "not-a-cron")
	if err == nil {
		t.Error("expected error for invalid cron expression")
	}
}
