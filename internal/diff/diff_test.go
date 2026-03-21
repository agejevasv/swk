package diff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- ReadTwoInputs ---

func TestReadTwoInputs_TwoFiles(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.txt")
	f2 := filepath.Join(dir, "b.txt")
	os.WriteFile(f1, []byte("hello"), 0o644)
	os.WriteFile(f2, []byte("world"), 0o644)

	a, b, err := ReadTwoInputs([]string{f1, f2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != "hello" || string(b) != "world" {
		t.Errorf("got %q and %q", a, b)
	}
}

func TestReadTwoInputs_Stdin(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "b.txt")
	os.WriteFile(f, []byte("file"), 0o644)

	stdin := strings.NewReader("stdin")
	a, b, err := ReadTwoInputs([]string{"-", f}, stdin)
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != "stdin" || string(b) != "file" {
		t.Errorf("got %q and %q", a, b)
	}
}

func TestReadTwoInputs_StdinSecond(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "a.txt")
	os.WriteFile(f, []byte("file"), 0o644)

	stdin := strings.NewReader("stdin")
	a, b, err := ReadTwoInputs([]string{f, "-"}, stdin)
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != "file" || string(b) != "stdin" {
		t.Errorf("got %q and %q", a, b)
	}
}

func TestReadTwoInputs_BothStdin(t *testing.T) {
	_, _, err := ReadTwoInputs([]string{"-", "-"}, nil)
	if err == nil {
		t.Fatal("expected error for both stdin")
	}
}

func TestReadTwoInputs_TooFewArgs(t *testing.T) {
	_, _, err := ReadTwoInputs([]string{"one"}, nil)
	if err == nil {
		t.Fatal("expected error for too few args")
	}
}

// --- DiffText ---

func TestDiffText_Identical(t *testing.T) {
	result := DiffText([]byte("hello\n"), []byte("hello\n"), 3)
	if result != "" {
		t.Errorf("expected empty diff for identical input, got %q", result)
	}
}

func TestDiffText_Changed(t *testing.T) {
	a := []byte("line1\nline2\nline3\n")
	b := []byte("line1\nchanged\nline3\n")
	result := DiffText(a, b, 3)
	if !strings.Contains(result, "-line2") {
		t.Errorf("expected removed line, got %q", result)
	}
	if !strings.Contains(result, "+changed") {
		t.Errorf("expected added line, got %q", result)
	}
}

func TestDiffText_Added(t *testing.T) {
	a := []byte("line1\n")
	b := []byte("line1\nline2\n")
	result := DiffText(a, b, 3)
	if !strings.Contains(result, "+line2") {
		t.Errorf("expected added line, got %q", result)
	}
}

func TestDiffText_Removed(t *testing.T) {
	a := []byte("line1\nline2\n")
	b := []byte("line1\n")
	result := DiffText(a, b, 3)
	if !strings.Contains(result, "-line2") {
		t.Errorf("expected removed line, got %q", result)
	}
}

func TestDiffText_Context(t *testing.T) {
	var lines []string
	for i := 0; i < 20; i++ {
		lines = append(lines, "line")
	}
	a := strings.Join(lines, "\n") + "\n"
	lines[10] = "changed"
	b := strings.Join(lines, "\n") + "\n"

	result := DiffText([]byte(a), []byte(b), 1)
	// With context=1, should not show all 20 lines
	resultLines := strings.Split(strings.TrimSpace(result), "\n")
	if len(resultLines) > 10 {
		t.Errorf("expected limited context, got %d lines", len(resultLines))
	}
}

// --- DiffJSON ---

func TestDiffJSON_Identical(t *testing.T) {
	j := []byte(`{"a":1,"b":2}`)
	result, err := DiffJSON(j, j, 3)
	if err != nil {
		t.Fatal(err)
	}
	if result != "" {
		t.Errorf("expected empty diff, got %q", result)
	}
}

func TestDiffJSON_DifferentKeyOrder(t *testing.T) {
	a := []byte(`{"b":2,"a":1}`)
	b := []byte(`{"a":1,"b":2}`)
	result, err := DiffJSON(a, b, 3)
	if err != nil {
		t.Fatal(err)
	}
	if result != "" {
		t.Errorf("expected empty diff for same content different order, got %q", result)
	}
}

func TestDiffJSON_ValueChanged(t *testing.T) {
	a := []byte(`{"name":"alice","age":30}`)
	b := []byte(`{"name":"alice","age":31}`)
	result, err := DiffJSON(a, b, 3)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "-") || !strings.Contains(result, "+") {
		t.Errorf("expected diff with changes, got %q", result)
	}
	if !strings.Contains(result, "30") && !strings.Contains(result, "31") {
		t.Errorf("expected age values in diff, got %q", result)
	}
}

func TestDiffJSON_KeyAdded(t *testing.T) {
	a := []byte(`{"a":1}`)
	b := []byte(`{"a":1,"b":2}`)
	result, err := DiffJSON(a, b, 3)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "+") {
		t.Errorf("expected addition in diff, got %q", result)
	}
}

func TestDiffJSON_KeyRemoved(t *testing.T) {
	a := []byte(`{"a":1,"b":2}`)
	b := []byte(`{"a":1}`)
	result, err := DiffJSON(a, b, 3)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "-") {
		t.Errorf("expected removal in diff, got %q", result)
	}
}

func TestDiffJSON_NestedChange(t *testing.T) {
	a := []byte(`{"user":{"name":"alice","role":"user"}}`)
	b := []byte(`{"user":{"name":"alice","role":"admin"}}`)
	result, err := DiffJSON(a, b, 3)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "user") && !strings.Contains(result, "admin") {
		t.Errorf("expected nested change in diff, got %q", result)
	}
}

func TestDiffJSON_ArrayChange(t *testing.T) {
	a := []byte(`{"items":[1,2,3]}`)
	b := []byte(`{"items":[1,2,4]}`)
	result, err := DiffJSON(a, b, 3)
	if err != nil {
		t.Fatal(err)
	}
	if result == "" {
		t.Error("expected diff for array change")
	}
}

func TestDiffJSON_InvalidFirst(t *testing.T) {
	_, err := DiffJSON([]byte("not json"), []byte(`{"a":1}`), 3)
	if err == nil {
		t.Fatal("expected error for invalid first input")
	}
}

func TestDiffJSON_InvalidSecond(t *testing.T) {
	_, err := DiffJSON([]byte(`{"a":1}`), []byte("not json"), 3)
	if err == nil {
		t.Fatal("expected error for invalid second input")
	}
}

// --- Colorize ---

func TestColorize_Empty(t *testing.T) {
	if Colorize("") != "" {
		t.Error("expected empty output for empty input")
	}
}

func TestColorize_Lines(t *testing.T) {
	input := "--- a\n+++ b\n@@ -1 +1 @@\n-old\n+new\n context\n"
	result := Colorize(input)
	if !strings.Contains(result, "\033[31m-old") {
		t.Error("expected red for removed line")
	}
	if !strings.Contains(result, "\033[32m+new") {
		t.Error("expected green for added line")
	}
	if !strings.Contains(result, "\033[36m@@") {
		t.Error("expected cyan for hunk header")
	}
	if !strings.Contains(result, "\033[1m--- a") {
		t.Error("expected bold for file header")
	}
}
