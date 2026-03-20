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

func TestJSON_Prettify(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json", `{"a":1,"b":2}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\n") {
		t.Errorf("expected newlines in prettified output, got %q", out)
	}
	if !strings.Contains(out, "  ") {
		t.Errorf("expected indentation in prettified output, got %q", out)
	}
}

func TestJSON_Minify(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json", "--minify", `{"a": 1}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if strings.Contains(trimmed, " ") {
		t.Errorf("expected no whitespace in minified output, got %q", trimmed)
	}
}

func TestJSON_ToYAML(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json", "--to", "yaml", `{"a":1}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "a: 1") {
		t.Errorf("expected 'a: 1', got %q", out)
	}
}

func TestJSON_FromYAML(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json", "--from", "yaml", "a: 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"a"`) {
		t.Errorf("expected '\"a\"' in output, got %q", out)
	}
}

func TestJSON_ToCSV(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json", "--to", "csv", `[{"name":"alice","age":"30"}]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "alice") {
		t.Errorf("expected 'alice' in output, got %q", out)
	}
}

func TestJSON_FromCSV(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json", "--from", "csv", "name,age\nalice,30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"alice"`) {
		t.Errorf("expected '\"alice\"' in output, got %q", out)
	}
}

func TestJSON_Invalid(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("json", `{invalid`)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
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

func TestXML_Prettify(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("xml", `<root><a>1</a></root>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\n") {
		t.Errorf("expected newlines in prettified output, got %q", out)
	}
}

func TestXML_Minify(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("xml", "--minify", "<root>\n  <a>1</a>\n</root>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if strings.Contains(trimmed, "\n") {
		t.Errorf("expected no newlines in minified output, got %q", trimmed)
	}
}
