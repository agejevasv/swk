package generate

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

func TestUUID_V7(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("uuid", "--version", "7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) != 36 {
		t.Errorf("expected 36-char UUID, got %d chars: %q", len(trimmed), trimmed)
	}
}

func TestUUID_V1(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("uuid", "--version", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) != 36 {
		t.Errorf("expected 36-char UUID, got %d chars: %q", len(trimmed), trimmed)
	}
}

func TestPassword_NoSymbols(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("password", "--no-symbols", "--length", "50")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	for _, r := range trimmed {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			t.Errorf("expected no symbols, got char %q in %q", string(r), trimmed)
			break
		}
	}
}

func TestText_Paragraphs(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("text", "--paragraphs", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) == "" {
		t.Error("expected non-empty paragraph output")
	}
}

func TestText_Sentences(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("text", "--sentences", "3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) == "" {
		t.Error("expected non-empty sentence output")
	}
}

func TestUUID_InvalidVersion(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("uuid", "--version", "99")
	if err == nil {
		t.Fatal("expected error for invalid UUID version, got nil")
	}
}

func TestCron_Every5m(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("cron", "--every", "5m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "*/5 * * * *" {
		t.Errorf("expected '*/5 * * * *', got %q", strings.TrimSpace(out))
	}
}

func TestCron_DailyAt(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("cron", "--daily", "--at", "9:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "0 9 * * *" {
		t.Errorf("expected '0 9 * * *', got %q", strings.TrimSpace(out))
	}
}

func TestCron_Weekdays(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("cron", "--weekdays", "--at", "9:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "0 9 * * 1-5" {
		t.Errorf("expected '0 9 * * 1-5', got %q", strings.TrimSpace(out))
	}
}

func TestCron_Weekly(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("cron", "--weekly", "--day", "FRI", "--at", "17:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "0 17 * * 5" {
		t.Errorf("expected '0 17 * * 5', got %q", strings.TrimSpace(out))
	}
}

func TestCron_NoSchedule(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("cron")
	if err == nil {
		t.Fatal("expected error with no schedule flag")
	}
}

func TestCron_DailyWithDay(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("cron", "--daily", "--day", "MON")
	if err == nil {
		t.Fatal("expected error for --daily with --day")
	}
}
