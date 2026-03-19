package gen

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func resetAllFlags() {
	// uuid.go package-level vars
	uuidVersion = 4
	uuidCount = 1
	uuidNamespace = ""
	uuidName = ""

	// hash.go package-level vars
	hashAlgo = "sha256"
	hashVerify = ""

	// password.go package-level vars
	pwLength = 16
	pwCount = 1
	pwNoUpper = false
	pwNoLower = false
	pwNoDigits = false
	pwNoSymbols = false
	pwExclude = ""

	// lorem.go package-level vars
	loremWords = 0
	loremSentences = 0
	loremParagraphs = 0

	// Reset all cobra subcommand flags to defaults and clear Changed state
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

// --- UUID tests ---

func TestUUIDDefault(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	// UUID v4 format: 8-4-4-4-12 hex chars
	if len(trimmed) != 36 {
		t.Fatalf("expected 36-char UUID, got %d chars: %s", len(trimmed), trimmed)
	}
	parts := strings.Split(trimmed, "-")
	if len(parts) != 5 {
		t.Fatalf("expected 5 dash-separated parts, got %d", len(parts))
	}
}

func TestUUIDCount3(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("uuid", "--count", "3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 UUIDs, got %d", len(lines))
	}
}

func TestUUIDVersion1(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("uuid", "--version", "1", "--count", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) != 36 {
		t.Fatalf("expected valid UUID, got: %s", trimmed)
	}
}

func TestUUIDVersion4(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("uuid", "--version", "4", "--count", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) != 36 {
		t.Fatalf("expected valid UUID, got: %s", trimmed)
	}
}

func TestUUIDVersion7(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("uuid", "--version", "7", "--count", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) != 36 {
		t.Fatalf("expected valid UUID, got: %s", trimmed)
	}
}

func TestUUIDInvalidVersion(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("uuid", "--version", "99")
	if err == nil {
		t.Fatal("expected error for invalid UUID version")
	}
}

// --- Hash tests ---

func TestHashDefaultSHA256(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("hash", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	// SHA256 produces 64 hex chars
	if len(trimmed) != 64 {
		t.Fatalf("expected 64-char sha256 hash, got %d chars: %s", len(trimmed), trimmed)
	}
}

func TestHashMD5(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("hash", "--algo", "md5", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	// MD5 produces 32 hex chars
	if len(trimmed) != 32 {
		t.Fatalf("expected 32-char md5 hash, got %d chars: %s", len(trimmed), trimmed)
	}
}

func TestHashVerifyCorrect(t *testing.T) {
	t.Cleanup(resetAllFlags)
	// First get the hash
	hashOut, err := executeCommand("hash", "--algo", "sha256", "hello")
	if err != nil {
		t.Fatalf("unexpected error getting hash: %v", err)
	}
	hash := strings.TrimSpace(hashOut)

	// Reset flags before next call
	resetAllFlags()

	// Now verify
	out, err := executeCommand("hash", "--algo", "sha256", "--verify", hash, "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "OK") {
		t.Fatalf("expected 'OK' for correct hash verification, got: %s", out)
	}
}

func TestHashVerifyWrong(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("hash", "--algo", "sha256", "--verify", "0000000000000000000000000000000000000000000000000000000000000000", "hello")
	if err == nil {
		t.Fatal("expected error for wrong hash verification")
	}
	if !strings.Contains(err.Error(), "mismatch") {
		t.Fatalf("expected mismatch error, got: %v", err)
	}
}

// --- Password tests ---

func TestPasswordDefaultLength(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) != 16 {
		t.Fatalf("expected 16-char password by default, got %d chars: %s", len(trimmed), trimmed)
	}
}

func TestPasswordLength32(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("password", "--length", "32")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) != 32 {
		t.Fatalf("expected 32-char password, got %d chars: %s", len(trimmed), trimmed)
	}
}

func TestPasswordNoSymbols(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("password", "--no-symbols", "--length", "50")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	symbols := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
	for _, ch := range trimmed {
		if strings.ContainsRune(symbols, ch) {
			t.Fatalf("found symbol %c in password with --no-symbols: %s", ch, trimmed)
		}
	}
}

func TestPasswordCount5(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("password", "--count", "5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 passwords, got %d", len(lines))
	}
}

// --- Lorem tests ---

func TestLoremDefaultParagraph(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("lorem")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) == 0 {
		t.Fatal("expected non-empty lorem ipsum output")
	}
}

func TestLoremWords10(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("lorem", "--words", "10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	words := strings.Fields(trimmed)
	if len(words) != 10 {
		t.Fatalf("expected 10 words, got %d: %s", len(words), trimmed)
	}
}

func TestLoremSentences3(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("lorem", "--sentences", "3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) == 0 {
		t.Fatal("expected non-empty output for 3 sentences")
	}
	// Sentences end with periods
	dotCount := strings.Count(trimmed, ".")
	if dotCount < 3 {
		t.Fatalf("expected at least 3 periods for 3 sentences, got %d", dotCount)
	}
}

func TestLoremParagraphs2(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("lorem", "--paragraphs", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) == 0 {
		t.Fatal("expected non-empty output for 2 paragraphs")
	}
	// Paragraphs are separated by blank lines
	if !strings.Contains(trimmed, "\n\n") {
		t.Fatal("expected paragraphs separated by blank lines")
	}
}
