package diff

import (
	"bytes"
	"os"
	"path/filepath"
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

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return p
}

func TestText_Identical(t *testing.T) {
	t.Cleanup(resetAllFlags)
	dir := t.TempDir()
	a := writeTempFile(t, dir, "a.txt", "hello\n")
	b := writeTempFile(t, dir, "b.txt", "hello\n")
	out, err := executeCommand("text", a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty output for identical files, got %q", out)
	}
}

func TestText_Different(t *testing.T) {
	t.Cleanup(resetAllFlags)
	dir := t.TempDir()
	a := writeTempFile(t, dir, "a.txt", "hello\n")
	b := writeTempFile(t, dir, "b.txt", "world\n")
	out, err := executeCommand("text", a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "-hello") {
		t.Errorf("expected '-hello' in diff output, got %q", out)
	}
	if !strings.Contains(out, "+world") {
		t.Errorf("expected '+world' in diff output, got %q", out)
	}
}

func TestText_MissingArgs(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("text", "only-one-arg")
	if err == nil {
		t.Fatal("expected error with only one argument")
	}
}

func TestJSON_Identical(t *testing.T) {
	t.Cleanup(resetAllFlags)
	dir := t.TempDir()
	a := writeTempFile(t, dir, "a.json", `{"b":2,"a":1}`)
	b := writeTempFile(t, dir, "b.json", `{"a":1,"b":2}`)
	out, err := executeCommand("json", a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty output for semantically identical JSON, got %q", out)
	}
}

func TestJSON_Different(t *testing.T) {
	t.Cleanup(resetAllFlags)
	dir := t.TempDir()
	a := writeTempFile(t, dir, "a.json", `{"a":1}`)
	b := writeTempFile(t, dir, "b.json", `{"a":2}`)
	out, err := executeCommand("json", a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "-") && !strings.Contains(out, "+") {
		t.Errorf("expected diff markers in output, got %q", out)
	}
}

func TestJSON_InvalidInput(t *testing.T) {
	t.Cleanup(resetAllFlags)
	dir := t.TempDir()
	a := writeTempFile(t, dir, "a.json", `not json`)
	b := writeTempFile(t, dir, "b.json", `{"a":1}`)
	_, err := executeCommand("json", a, b)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
