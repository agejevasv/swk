package query

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

func TestJSON_BasicQuery(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json", "$.name", `{"name":"alice"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "alice") {
		t.Errorf("expected 'alice' in output, got %q", out)
	}
}

func TestJSON_NestedQuery(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("json", "$.a.b", `{"a":{"b":"deep"}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "deep") {
		t.Errorf("expected 'deep' in output, got %q", out)
	}
}

func TestRegex_DefaultPrintsMatchingLines(t *testing.T) {
	t.Cleanup(resetAllFlags)
	// Default mode is line-oriented: prints full lines that match.
	out, err := executeCommand("regex", `\d+`, "abc123def")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "abc123def" {
		t.Errorf("expected full line 'abc123def', got %q", out)
	}
}

func TestRegex_OnlyMatchingExtractsValues(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("regex", "-o", "--global", `\d+`, "abc123def456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "123") {
		t.Errorf("expected '123' in output, got %q", out)
	}
	if !strings.Contains(out, "456") {
		t.Errorf("expected '456' in output, got %q", out)
	}
}

func TestRegex_Replace(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("regex", "--replace", "NUM", `\d+`, "abc123def")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "abcNUMdef") {
		t.Errorf("expected 'abcNUMdef' in output, got %q", out)
	}
}

func TestHTML_Query(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("html", "span", "<div><span>hello</span></div>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("expected 'hello' in output, got %q", out)
	}
}

func TestJSON_NoMatch(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("json", "$.nonexistent", `{"name":"alice"}`)
	if err == nil {
		t.Fatal("expected error for no match, got nil")
	}
}

func TestRegex_NoMatch(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("regex", `zzz`, "abc123")
	if err == nil {
		t.Fatal("expected error for no match, got nil")
	}
}

func TestRegex_Groups(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("regex", "--groups", `(\w+):(\d+)`, "John:30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "John") {
		t.Errorf("expected 'John' in groups output, got %q", out)
	}
	if !strings.Contains(out, "30") {
		t.Errorf("expected '30' in groups output, got %q", out)
	}
}

func TestRegex_GroupsNoMatch(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("regex", "--groups", `zzz`, "abc")
	if err == nil {
		t.Fatal("expected error for no match with --groups, got nil")
	}
}

func TestRegex_OnlyMatchingNoMatch(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("regex", "-o", `zzz`, "abc")
	if err == nil {
		t.Fatal("expected error for no match with -o, got nil")
	}
}

func TestRegex_InvalidPattern(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("regex", `[invalid`, "abc")
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestHTML_NoMatch(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("html", "span.nonexistent", "<div>hello</div>")
	if err == nil {
		t.Fatal("expected error for no match, got nil")
	}
}

func TestHTML_QueryAttr(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("html", "a", "--attr", "href", `<a href="http://example.com">link</a>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "http://example.com") {
		t.Errorf("expected 'http://example.com' in output, got %q", out)
	}
}
