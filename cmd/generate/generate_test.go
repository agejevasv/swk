package generate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func resetAllFlags() {
	genWidth = 800
	genHeight = 600
	genStyle = "mixed"
	genOutput = ""
	pwLength = 16
	pwCount = 1
	pwNoUpper = false
	pwNoLower = false
	pwNoDigits = false
	pwNoSymbols = false
	pwExclude = ""
	textWords = 0
	textSentences = 0
	textParagraphs = 0
	uuidVersion = 4
	uuidCount = 1
	uuidNamespace = ""
	uuidName = ""

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

func TestPassword_Default(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) != 16 {
		t.Errorf("expected 16-char password, got %d chars: %q", len(trimmed), trimmed)
	}
}

func TestPassword_Length32(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("password", "--length", "32")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) != 32 {
		t.Errorf("expected 32-char password, got %d chars: %q", len(trimmed), trimmed)
	}
}

func TestPassword_Count5(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("password", "--count", "5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d", len(lines))
	}
}

func TestText_Default(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) == "" {
		t.Error("expected non-empty text output")
	}
}

func TestText_Words10(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("text", "--words", "10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	words := strings.Fields(strings.TrimSpace(out))
	if len(words) != 10 {
		t.Errorf("expected 10 words, got %d: %q", len(words), out)
	}
}

func TestUUID_Default(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) != 36 {
		t.Errorf("expected 36-char UUID, got %d chars: %q", len(trimmed), trimmed)
	}
	parts := strings.Split(trimmed, "-")
	if len(parts) != 5 {
		t.Errorf("expected 5 hyphen-separated parts, got %d", len(parts))
	}
}

func TestUUID_Count3(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("uuid", "--count", "3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestUUID_InvalidVersion(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("uuid", "--version", "99")
	if err == nil {
		t.Fatal("expected error for invalid UUID version, got nil")
	}
}
