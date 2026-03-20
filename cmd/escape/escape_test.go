package escape

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

func TestHTML_Escape(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("html", "<div>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "&lt;") {
		t.Errorf("expected '&lt;' in output, got %q", out)
	}
}

func TestHTML_Unescape(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("html", "-u", "&lt;div&gt;")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "<div>") {
		t.Errorf("expected '<div>' in output, got %q", out)
	}
}

func TestJSON_Escape(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json", `Hello "World"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `\"`) {
		t.Errorf("expected escaped quotes in output, got %q", out)
	}
}

func TestJSON_Unescape(t *testing.T) {
	t.Cleanup(resetAllFlags)
	// First escape
	escaped, err := executeCommand("json", `Hello "World"`)
	if err != nil {
		t.Fatalf("escape error: %v", err)
	}
	escaped = strings.TrimSpace(escaped)

	resetAllFlags()

	// Then unescape
	out, err := executeCommand("json", "-u", escaped)
	if err != nil {
		t.Fatalf("unescape error: %v", err)
	}
	if !strings.Contains(out, `Hello "World"`) {
		t.Errorf("expected roundtrip to original, got %q", out)
	}
}

func TestShell_Escape(t *testing.T) {
	t.Cleanup(resetAllFlags)
	input := "it's a test"
	out, err := executeCommand("shell", input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		t.Error("expected non-empty output")
	}
	if trimmed == input {
		t.Error("expected escaped output to differ from input")
	}
}

func TestURL_Encode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("url", "hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "%") {
		t.Errorf("expected '%%' in URL-encoded output, got %q", out)
	}
}

func TestURL_Decode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("url", "-u", "hello%20world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "hello world") {
		t.Errorf("expected 'hello world', got %q", out)
	}
}

func TestURL_Component(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("url", "--component", "hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "+") {
		t.Errorf("expected '+' in component-encoded output, got %q", out)
	}
}

func TestXML_Escape(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("xml", "<tag>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "&lt;") {
		t.Errorf("expected '&lt;' in output, got %q", out)
	}
}

func TestXML_Unescape(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("xml", "-u", "&lt;tag&gt;")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "<tag>") {
		t.Errorf("expected '<tag>' in output, got %q", out)
	}
}

func TestShell_Unescape(t *testing.T) {
	t.Cleanup(resetAllFlags)
	// Escape first
	escaped, err := executeCommand("shell", "it's a test")
	if err != nil {
		t.Fatalf("escape error: %v", err)
	}
	escaped = strings.TrimSpace(escaped)

	resetAllFlags()

	// Then unescape
	out, err := executeCommand("shell", "-u", escaped)
	if err != nil {
		t.Fatalf("unescape error: %v", err)
	}
	if !strings.Contains(out, "it's a test") {
		t.Errorf("expected roundtrip to original, got %q", out)
	}
}

func TestURL_ComponentDecode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("url", "-u", "--component", "hello+world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "hello world") {
		t.Errorf("expected 'hello world', got %q", out)
	}
}
