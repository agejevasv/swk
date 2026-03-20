package convert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func resetAllFlags() {
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

func TestBase_DecToHex(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("base", "255")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "ff") {
		t.Errorf("expected output to contain 'ff' or '0xff', got %q", out)
	}
}

func TestBase_BinToDec(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("base", "--from", "bin", "--to", "dec", "11111111")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.TrimSpace(out), "255") {
		t.Errorf("expected '255', got %q", out)
	}
}

func TestBase_InvalidBase(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("base", "--from", "xyz", "42")
	if err == nil {
		t.Fatal("expected error for invalid base, got nil")
	}
}

func TestCase_Snake(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("case", "--to", "snake", "helloWorld")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.TrimSpace(out), "hello_world") {
		t.Errorf("expected 'hello_world', got %q", out)
	}
}

func TestCase_Camel(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("case", "--to", "camel", "hello_world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.TrimSpace(out), "helloWorld") {
		t.Errorf("expected 'helloWorld', got %q", out)
	}
}

func TestColor_HexToAll(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("color", "#ff0000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "rgb") && !strings.Contains(lower, "hsl") {
		t.Errorf("expected output to contain 'rgb' or 'hsl', got %q", out)
	}
}

func TestDate_UnixToISO(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("date", "--from", "unix", "--tz", "UTC", "0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "1970-01-01") {
		t.Errorf("expected '1970-01-01', got %q", out)
	}
}

func TestDate_Now(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("date", "--tz", "UTC", "now")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		t.Fatal("expected non-empty output")
	}
	if !strings.Contains(trimmed, "T") {
		t.Errorf("expected output to contain 'T', got %q", trimmed)
	}
}

func TestYAML2JSON(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("yaml2json", "a: 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"a"`) {
		t.Errorf("expected '\"a\"' in output, got %q", out)
	}
}

func TestJSON2YAML(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json2yaml", `{"a":1}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "a: 1") {
		t.Errorf("expected 'a: 1', got %q", out)
	}
}

func TestCSV2JSON(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("csv2json", "name,age\nalice,30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"alice"`) {
		t.Errorf("expected '\"alice\"' in output, got %q", out)
	}
}

func TestJSON2CSV(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json2csv", `[{"name":"alice","age":"30"}]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "alice") {
		t.Errorf("expected 'alice' in output, got %q", out)
	}
}

func TestMarkdown_HTML(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("markdown", "--html", "# Hello\n\nWorld")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "<h1>") {
		t.Errorf("expected '<h1>' in output, got %q", out)
	}
}

func TestBytes_ToHuman(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("bytes", "1073741824")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "GiB") {
		t.Errorf("expected 'GiB' in output, got %q", out)
	}
}

func TestBytes_ToBytes(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("bytes", "1GiB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "1073741824") {
		t.Errorf("expected '1073741824' in output, got %q", out)
	}
}

func TestBytes_Decimal(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("bytes", "-d", "1000000000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "GB") {
		t.Errorf("expected 'GB' in output, got %q", out)
	}
}

func TestChmod_Numeric(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("chmod", "--to", "symbolic", "755")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "rwxr-xr-x") {
		t.Errorf("expected 'rwxr-xr-x', got %q", out)
	}
}

func TestChmod_Symbolic(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("chmod", "--to", "numeric", "rwxr-xr-x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "755") {
		t.Errorf("expected '755', got %q", out)
	}
}

func TestChmod_Explain(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("chmod", "755")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Owner:") {
		t.Errorf("expected 'Owner:' in explain output, got %q", out)
	}
}

func TestChmod_Setuid(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("chmod", "4755")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "setuid") {
		t.Errorf("expected 'setuid' in output, got %q", out)
	}
}

func TestColor_RGBToHex(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("color", "--from", "rgb", "--to", "hex", "255,87,51")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.ToUpper(out), "FF5733") {
		t.Errorf("expected 'FF5733' in output, got %q", out)
	}
}

func TestDuration_SecondsToHuman(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("duration", "86400")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "1d") {
		t.Errorf("expected '1d', got %q", out)
	}
}

func TestDuration_HumanToSeconds(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("duration", "2d 5h 30m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "192600") {
		t.Errorf("expected '192600', got %q", out)
	}
}

func TestDuration_ExplicitTo(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("duration", "--to", "seconds", "2d 5h")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "190800") {
		t.Errorf("expected '190800', got %q", out)
	}
}

func TestDate_CustomLayout(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("date", "--from", "unix", "--to", "2006-01-02", "--tz", "UTC", "1700000000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "2023-11-14") {
		t.Errorf("expected '2023-11-14', got %q", out)
	}
}

func TestMarkdown_PlainText(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("markdown", "**bold** text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "**") {
		t.Errorf("expected markdown stripped, got %q", out)
	}
}

func TestCSV2JSON_Basic(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("csv2json", "name,age\nalice,30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"alice"`) {
		t.Errorf("expected '\"alice\"' in output, got %q", out)
	}
}

func TestJSON2CSV_Basic(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json2csv", `[{"name":"alice","age":"30"}]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "alice") {
		t.Errorf("expected 'alice' in output, got %q", out)
	}
}

func TestParseResize(t *testing.T) {
	tests := []struct {
		input   string
		wantW   int
		wantH   int
		wantErr bool
	}{
		// Valid cases
		{input: "800x600", wantW: 800, wantH: 600},
		{input: "1920x1080", wantW: 1920, wantH: 1080},
		{input: "1x1", wantW: 1, wantH: 1},
		// Invalid cases
		{input: "abc", wantErr: true},
		{input: "800", wantErr: true},
		{input: "0x600", wantErr: true},
		{input: "800x0", wantErr: true},
		{input: "-1x100", wantErr: true},
		{input: "axb", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			w, h, err := parseResize(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseResize(%q) expected error, got w=%d h=%d", tt.input, w, h)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseResize(%q) unexpected error: %v", tt.input, err)
			}
			if w != tt.wantW || h != tt.wantH {
				t.Errorf("parseResize(%q) = (%d, %d), want (%d, %d)", tt.input, w, h, tt.wantW, tt.wantH)
			}
		})
	}
}
