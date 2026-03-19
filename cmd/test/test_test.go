package test

import (
	"bytes"
	"encoding/json"
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

func resetRegexFlags() {
	regexPattern = ""
	regexGlobal = false
	regexGroups = false
	regexReplace = ""
	regexCmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

func resetJSONPathFlags() {
	jsonpathQuery = ""
	jsonpathCmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

func resetXMLValFlags() {
	xmlvalCmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

// --- Regex tests ---

func TestRegexSimpleMatch(t *testing.T) {
	t.Cleanup(resetRegexFlags)

	out, err := executeCommand("regex", "--pattern", `\d+`, "abc123def")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "123") {
		t.Fatalf("expected match '123', got: %s", out)
	}
}

func TestRegexGlobalFindsAll(t *testing.T) {
	t.Cleanup(resetRegexFlags)

	out, err := executeCommand("regex", "--pattern", `\d+`, "--global", "abc123def456ghi789")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "123") || !strings.Contains(out, "456") || !strings.Contains(out, "789") {
		t.Fatalf("expected all matches, got: %s", out)
	}
}

func TestRegexGroupsShowsJSON(t *testing.T) {
	t.Cleanup(resetRegexFlags)

	out, err := executeCommand("regex", "--pattern", `(\w+)@(\w+)`, "--groups", "user@host")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !json.Valid([]byte(out)) {
		t.Fatalf("expected valid JSON output for --groups, got: %s", out)
	}
}

func TestRegexReplace(t *testing.T) {
	t.Cleanup(resetRegexFlags)

	out, err := executeCommand("regex", "--pattern", `\d+`, "--replace", "NUM", "abc123def")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "abcNUMdef") {
		t.Fatalf("expected replaced output 'abcNUMdef', got: %s", out)
	}
}

func TestRegexPatternRequired(t *testing.T) {
	t.Cleanup(resetRegexFlags)

	out, err := executeCommand("regex", "some input")
	if err == nil && !strings.Contains(out, "required") && !strings.Contains(out, "pattern") {
		t.Fatal("expected error or message about required --pattern flag")
	}
}

func TestRegexInvalidPattern(t *testing.T) {
	t.Cleanup(resetRegexFlags)

	_, err := executeCommand("regex", "--pattern", `[invalid`, "test")
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

// --- JSONPath tests ---

func TestJSONPathBasicQuery(t *testing.T) {
	t.Cleanup(resetJSONPathFlags)

	out, err := executeCommand("jsonpath", "--query", "$.name", `{"name":"alice","age":30}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "alice") {
		t.Fatalf("expected 'alice' in output, got: %s", out)
	}
}

func TestJSONPathNestedQuery(t *testing.T) {
	t.Cleanup(resetJSONPathFlags)

	out, err := executeCommand("jsonpath", "--query", "$.a.b", `{"a":{"b":"deep"}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "deep") {
		t.Fatalf("expected 'deep' in output, got: %s", out)
	}
}

func TestJSONPathArrayAccess(t *testing.T) {
	t.Cleanup(resetJSONPathFlags)

	out, err := executeCommand("jsonpath", "--query", "$.items[0]", `{"items":["first","second"]}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "first") {
		t.Fatalf("expected 'first' in output, got: %s", out)
	}
}

func TestJSONPathQueryRequired(t *testing.T) {
	t.Cleanup(resetJSONPathFlags)

	out, err := executeCommand("jsonpath", `{"a":1}`)
	if err == nil && !strings.Contains(out, "required") && !strings.Contains(out, "query") {
		t.Fatal("expected error or message about required --query flag")
	}
}

// --- XMLVal tests ---

func TestXMLValValid(t *testing.T) {
	t.Cleanup(resetXMLValFlags)

	out, err := executeCommand("xmlval", `<root><child>text</child></root>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Valid") {
		t.Fatalf("expected 'Valid' for well-formed XML, got: %s", out)
	}
}

func TestXMLValInvalid(t *testing.T) {
	t.Cleanup(resetXMLValFlags)

	_, err := executeCommand("xmlval", `<root><unclosed>`)
	if err == nil {
		t.Fatal("expected error for invalid XML")
	}
}
