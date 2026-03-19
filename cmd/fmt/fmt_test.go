package fmt

import (
	"bytes"
	"strings"
	"testing"

	pflag "github.com/spf13/pflag"
)

func resetFlagChanged(flags *pflag.FlagSet) {
	flags.VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetErr(buf)
	Cmd.SetArgs(args)
	err := Cmd.Execute()
	return buf.String(), err
}

// --- JSON tests ---

func TestJSONPrettyPrint(t *testing.T) {
	t.Cleanup(func() {
		jsonMinify = false
		jsonIndent = 2
		resetFlagChanged(jsonCmd.Flags())
	})

	out, err := executeCommand("json", `{"a":1,"b":2}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\n") {
		t.Fatal("expected indented output with newlines")
	}
	if !strings.Contains(out, "  ") {
		t.Fatal("expected 2-space indentation by default")
	}
}

func TestJSONMinify(t *testing.T) {
	t.Cleanup(func() {
		jsonMinify = false
		jsonIndent = 2
		resetFlagChanged(jsonCmd.Flags())
	})

	out, err := executeCommand("json", "--minify", `{"a": 1, "b": 2}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if strings.Contains(trimmed, " ") || strings.Contains(trimmed, "\n  ") {
		t.Fatalf("expected minified output with no extra whitespace, got: %s", trimmed)
	}
}

func TestJSONIndent4(t *testing.T) {
	t.Cleanup(func() {
		jsonMinify = false
		jsonIndent = 2
		resetFlagChanged(jsonCmd.Flags())
	})

	out, err := executeCommand("json", "--indent", "4", `{"a":1}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "    ") {
		t.Fatal("expected 4-space indentation")
	}
}


func TestJSONInvalidInput(t *testing.T) {
	t.Cleanup(func() {
		jsonMinify = false
		jsonIndent = 2
		resetFlagChanged(jsonCmd.Flags())
	})

	_, err := executeCommand("json", `{not valid json}`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// --- XML tests ---

func TestXMLPrettyPrint(t *testing.T) {
	t.Cleanup(func() {
		xmlMinify = false
		xmlIndent = 2
		resetFlagChanged(xmlCmd.Flags())
	})

	out, err := executeCommand("xml", `<root><a>1</a></root>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\n") {
		t.Fatal("expected indented output with newlines")
	}
}

func TestXMLMinify(t *testing.T) {
	t.Cleanup(func() {
		xmlMinify = false
		xmlIndent = 2
		resetFlagChanged(xmlCmd.Flags())
	})

	input := "<root>\n  <a>1</a>\n</root>"
	out, err := executeCommand("xml", "--minify", input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if strings.Contains(trimmed, "\n") {
		t.Fatalf("expected minified output on one line, got: %s", trimmed)
	}
}

func TestXMLIndent4(t *testing.T) {
	t.Cleanup(func() {
		xmlMinify = false
		xmlIndent = 2
		resetFlagChanged(xmlCmd.Flags())
	})

	out, err := executeCommand("xml", "--indent", "4", `<root><a>1</a></root>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "    ") {
		t.Fatal("expected 4-space indentation")
	}
}

func TestXMLInvalidInput(t *testing.T) {
	t.Cleanup(func() {
		xmlMinify = false
		xmlIndent = 2
		resetFlagChanged(xmlCmd.Flags())
	})

	_, err := executeCommand("xml", `<not valid xml`)
	if err == nil {
		t.Fatal("expected error for invalid XML")
	}
}

// --- SQL tests ---

func TestSQLBasicFormatting(t *testing.T) {
	t.Cleanup(func() {
		sqlUppercase = false
		resetFlagChanged(sqlCmd.Flags())
	})

	out, err := executeCommand("sql", "select a, b from t where a = 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) == "" {
		t.Fatal("expected non-empty formatted SQL output")
	}
}

func TestSQLUppercase(t *testing.T) {
	t.Cleanup(func() {
		sqlUppercase = false
		resetFlagChanged(sqlCmd.Flags())
	})

	out, err := executeCommand("sql", "--uppercase", "select a from t where a = 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "SELECT") {
		t.Fatalf("expected uppercase SELECT keyword, got: %s", out)
	}
	if !strings.Contains(out, "FROM") {
		t.Fatalf("expected uppercase FROM keyword, got: %s", out)
	}
	if !strings.Contains(out, "WHERE") {
		t.Fatalf("expected uppercase WHERE keyword, got: %s", out)
	}
}
