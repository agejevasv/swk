package format

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

func TestJSON_Invalid(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("json", `{invalid`)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
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

func TestJSON_CustomIndent(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json", "--indent", "4", `{"a":1}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "    ") {
		t.Errorf("expected 4-space indentation, got %q", out)
	}
}

func TestJSON2Table_Box(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json2table", `[{"name":"alice","age":30},{"name":"bob","age":25}]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "alice") {
		t.Errorf("expected 'alice' in output, got %q", out)
	}
	if !strings.Contains(out, "bob") {
		t.Errorf("expected 'bob' in output, got %q", out)
	}
}

func TestJSON2Table_Simple(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json2table", "--style", "simple", `[{"x":1}]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "x") {
		t.Errorf("expected 'x' in output, got %q", out)
	}
}

func TestJSON2Table_Plain(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json2table", "--style", "plain", `[{"x":1}]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "x") {
		t.Errorf("expected 'x' in output, got %q", out)
	}
}

func TestJSON2Table_Invalid(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("json2table", `not json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestCSV2Table(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("csv2table", "name,age\nalice,30\nbob,25")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "alice") {
		t.Errorf("expected 'alice' in output, got %q", out)
	}
}

func TestCSV2Table_CustomDelimiter(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("csv2table", "--delimiter", ";", "name;age\nalice;30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "alice") {
		t.Errorf("expected 'alice' in output, got %q", out)
	}
}

func TestXML_Invalid(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("xml", `not xml at all <<<>>>`)
	if err == nil {
		t.Fatal("expected error for invalid XML, got nil")
	}
}
