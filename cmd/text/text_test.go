package text

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	pflag "github.com/spf13/pflag"
)

func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetErr(buf)
	Cmd.SetArgs(args)
	err := Cmd.Execute()
	return buf.String(), err
}

func resetFlagChanged(flags *pflag.FlagSet) {
	flags.VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

// --- Inspect tests ---

func TestInspectDefaultTable(t *testing.T) {
	t.Cleanup(func() {
		inspectJSON = false
		resetFlagChanged(inspectCmd.Flags())
	})

	out, err := executeCommand("inspect", "Hello World")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Characters:") {
		t.Fatalf("expected 'Characters:' in table output, got: %s", out)
	}
	if !strings.Contains(out, "Words:") {
		t.Fatalf("expected 'Words:' in table output, got: %s", out)
	}
}

func TestInspectJSON(t *testing.T) {
	t.Cleanup(func() {
		inspectJSON = false
		resetFlagChanged(inspectCmd.Flags())
	})

	out, err := executeCommand("inspect", "--json", "Hello World")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !json.Valid([]byte(out)) {
		t.Fatalf("expected valid JSON output, got: %s", out)
	}
}

// --- Escape tests ---

func TestEscapeDefaultJSON(t *testing.T) {
	t.Cleanup(func() {
		escapeMode = "json"
		escapeUnescape = false
		resetFlagChanged(escapeCmd.Flags())
	})

	out, err := executeCommand("escape", `Hello "World"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `\"`) {
		t.Fatalf("expected escaped double quotes in JSON mode, got: %s", out)
	}
}

func TestEscapeModeXML(t *testing.T) {
	t.Cleanup(func() {
		escapeMode = "json"
		escapeUnescape = false
		resetFlagChanged(escapeCmd.Flags())
	})

	out, err := executeCommand("escape", "--mode", "xml", `<tag>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "&lt;") {
		t.Fatalf("expected XML-escaped '<', got: %s", out)
	}
}

func TestEscapeModeHTML(t *testing.T) {
	t.Cleanup(func() {
		escapeMode = "json"
		escapeUnescape = false
		resetFlagChanged(escapeCmd.Flags())
	})

	out, err := executeCommand("escape", "--mode", "html", `<div>&</div>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "&amp;") || !strings.Contains(out, "&lt;") {
		t.Fatalf("expected HTML-escaped output, got: %s", out)
	}
}

func TestEscapeUnescape(t *testing.T) {
	t.Cleanup(func() {
		escapeMode = "json"
		escapeUnescape = false
		resetFlagChanged(escapeCmd.Flags())
	})

	// First escape
	escaped, err := executeCommand("escape", "--mode", "xml", `<tag>`)
	if err != nil {
		t.Fatalf("unexpected error escaping: %v", err)
	}

	// Reset flags before next call
	escapeMode = "json"
	escapeUnescape = false

	// Then unescape
	unescaped, err := executeCommand("escape", "--mode", "xml", "--unescape", strings.TrimSpace(escaped))
	if err != nil {
		t.Fatalf("unexpected error unescaping: %v", err)
	}
	if !strings.Contains(unescaped, "<tag>") {
		t.Fatalf("expected unescaped '<tag>', got: %s", unescaped)
	}
}

// --- Case tests ---

func TestCaseSnake(t *testing.T) {
	t.Cleanup(func() {
		caseTo = ""
		resetFlagChanged(caseCmd.Flags())
	})

	out, err := executeCommand("case", "--to", "snake", "helloWorld")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "hello_world") {
		t.Fatalf("expected 'hello_world', got: %s", out)
	}
}

func TestCaseCamel(t *testing.T) {
	t.Cleanup(func() {
		caseTo = ""
		resetFlagChanged(caseCmd.Flags())
	})

	out, err := executeCommand("case", "--to", "camel", "hello_world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "helloWorld") {
		t.Fatalf("expected 'helloWorld', got: %s", out)
	}
}

func TestCaseKebab(t *testing.T) {
	t.Cleanup(func() {
		caseTo = ""
		resetFlagChanged(caseCmd.Flags())
	})

	out, err := executeCommand("case", "--to", "kebab", "helloWorld")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "hello-world") {
		t.Fatalf("expected 'hello-world', got: %s", out)
	}
}

func TestCaseToRequired(t *testing.T) {
	t.Cleanup(func() {
		caseTo = ""
		resetFlagChanged(caseCmd.Flags())
	})

	_, err := executeCommand("case", "helloWorld")
	if err == nil {
		t.Fatal("expected error when --to is not provided")
	}
}

// --- Diff tests ---

func TestDiffOutput(t *testing.T) {
	t.Cleanup(func() {
		diffFile1 = ""
		diffFile2 = ""
		diffContext = 3
		resetFlagChanged(diffCmd.Flags())
	})

	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.txt")
	f2 := filepath.Join(dir, "b.txt")
	os.WriteFile(f1, []byte("line1\nline2\nline3\n"), 0644)
	os.WriteFile(f2, []byte("line1\nmodified\nline3\n"), 0644)

	out, err := executeCommand("diff", "--file1", f1, "--file2", f2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "line2") || !strings.Contains(out, "modified") {
		t.Fatalf("expected diff showing changes, got: %s", out)
	}
}

func TestDiffContextFlag(t *testing.T) {
	t.Cleanup(func() {
		diffFile1 = ""
		diffFile2 = ""
		diffContext = 3
		resetFlagChanged(diffCmd.Flags())
	})

	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.txt")
	f2 := filepath.Join(dir, "b.txt")
	os.WriteFile(f1, []byte("1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n"), 0644)
	os.WriteFile(f2, []byte("1\n2\n3\n4\nFIVE\n6\n7\n8\n9\n10\n"), 0644)

	out, err := executeCommand("diff", "--file1", f1, "--file2", f2, "--context", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// With context=1, we should see limited surrounding lines
	if !strings.Contains(out, "FIVE") {
		t.Fatalf("expected diff to contain 'FIVE', got: %s", out)
	}
}

// --- Markdown tests ---

func TestMarkdownHTML(t *testing.T) {
	t.Cleanup(func() {
		mdHTML = false
		resetFlagChanged(markdownCmd.Flags())
	})

	out, err := executeCommand("markdown", "--html", "# Hello\n\nWorld")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "<h1>") || !strings.Contains(out, "Hello") {
		t.Fatalf("expected HTML h1 tag, got: %s", out)
	}
}

func TestMarkdownDefaultPlainText(t *testing.T) {
	t.Cleanup(func() {
		mdHTML = false
		resetFlagChanged(markdownCmd.Flags())
	})

	out, err := executeCommand("markdown", "**bold** text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Default strips markdown - should not contain markdown syntax
	if strings.Contains(out, "**") {
		t.Fatalf("expected markdown stripped, but found ** in output: %s", out)
	}
}
