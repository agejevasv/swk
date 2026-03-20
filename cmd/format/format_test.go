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
