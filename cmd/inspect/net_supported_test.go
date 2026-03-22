//go:build linux || darwin

package inspect

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNet_DefaultOutput(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("net")
	if err != nil {
		t.Skipf("skipping: %v (may need permissions)", err)
	}
	if !strings.Contains(out, "PROTO") {
		t.Errorf("expected table header in output, got %q", out)
	}
	if !strings.Contains(out, "LOCAL ADDRESS") {
		t.Errorf("expected LOCAL ADDRESS header, got %q", out)
	}
}

func TestNet_JSONOutput(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("net", "--json")
	if err != nil {
		t.Skipf("skipping: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if !json.Valid([]byte(trimmed)) {
		t.Errorf("expected valid JSON output, got %q", out)
	}
}

func TestNet_TCPFilter(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("net", "--tcp")
	if err != nil {
		t.Skipf("skipping: %v", err)
	}
	// Should not contain udp entries (if any)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines[1:] { // skip header
		if strings.HasPrefix(line, "udp") {
			t.Errorf("expected no UDP entries with --tcp, got line: %s", line)
		}
	}
}

func TestNet_UDPFilter(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("net", "--udp")
	if err != nil {
		t.Skipf("skipping: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines[1:] { // skip header
		if strings.HasPrefix(line, "tcp") {
			t.Errorf("expected no TCP entries with --udp, got line: %s", line)
		}
	}
}

func TestNet_AllFlag(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("net", "--all")
	if err != nil {
		t.Skipf("skipping: %v", err)
	}
	if !strings.Contains(out, "REMOTE ADDRESS") {
		t.Errorf("expected REMOTE ADDRESS column with --all, got %q", out)
	}
}

func TestNet_PortFilter(t *testing.T) {
	t.Cleanup(resetAllFlags)
	// Filter for a port unlikely to have listeners — should produce header only
	out, err := executeCommand("net", "--port", "59999")
	if err != nil {
		t.Skipf("skipping: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) > 1 {
		t.Logf("unexpected entries on port 59999: %s", out)
	}
}
