//go:build darwin

package inspect

import (
	"os/exec"
	"testing"
)

func TestParseLsofOutput_TCPListen(t *testing.T) {
	output := "p1234\ncnginx\nLroot\nf4\ntIPv4\nn*:80\nTST=LISTEN\nTQR=0\nTQS=0\n"

	entries := parseLsofOutput(output)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	e := entries[0]
	if e.Proto != "tcp" {
		t.Errorf("expected proto tcp, got %s", e.Proto)
	}
	if e.LocalIP != "0.0.0.0" {
		t.Errorf("expected local IP 0.0.0.0, got %s", e.LocalIP)
	}
	if e.LocalPort != 80 {
		t.Errorf("expected local port 80, got %d", e.LocalPort)
	}
	if e.State != "LISTEN" {
		t.Errorf("expected state LISTEN, got %s", e.State)
	}
	if e.PID != 1234 {
		t.Errorf("expected PID 1234, got %d", e.PID)
	}
	if e.Process != "nginx" {
		t.Errorf("expected process nginx, got %s", e.Process)
	}
	if e.User != "root" {
		t.Errorf("expected user root, got %s", e.User)
	}
}

func TestParseLsofOutput_TCPEstablished(t *testing.T) {
	output := "p5678\ncnode\nLuser1\nf12\ntIPv4\nn127.0.0.1:3000->10.0.0.5:54321\nTST=ESTABLISHED\nTQR=0\nTQS=0\n"

	entries := parseLsofOutput(output)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	e := entries[0]
	if e.Proto != "tcp" {
		t.Errorf("expected proto tcp, got %s", e.Proto)
	}
	if e.LocalIP != "127.0.0.1" {
		t.Errorf("expected local IP 127.0.0.1, got %s", e.LocalIP)
	}
	if e.LocalPort != 3000 {
		t.Errorf("expected local port 3000, got %d", e.LocalPort)
	}
	if e.RemoteIP != "10.0.0.5" {
		t.Errorf("expected remote IP 10.0.0.5, got %s", e.RemoteIP)
	}
	if e.RemotePort != 54321 {
		t.Errorf("expected remote port 54321, got %d", e.RemotePort)
	}
	if e.State != "ESTABLISHED" {
		t.Errorf("expected state ESTABLISHED, got %s", e.State)
	}
}

func TestParseLsofOutput_UDP(t *testing.T) {
	output := "p9999\ncdnsmasq\nLnobody\nf5\ntIPv4\nn*:53\n"

	entries := parseLsofOutput(output)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	e := entries[0]
	if e.Proto != "udp" {
		t.Errorf("expected proto udp, got %s", e.Proto)
	}
	if e.LocalPort != 53 {
		t.Errorf("expected local port 53, got %d", e.LocalPort)
	}
	if e.State != "" {
		t.Errorf("expected empty state for UDP, got %s", e.State)
	}
}

func TestParseLsofOutput_IPv6Listen(t *testing.T) {
	output := "p1000\ncapache\nLwww\nf6\ntIPv6\nn[::1]:443\nTST=LISTEN\nTQR=0\nTQS=0\n"

	entries := parseLsofOutput(output)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	e := entries[0]
	if e.Proto != "tcp6" {
		t.Errorf("expected proto tcp6, got %s", e.Proto)
	}
	if e.LocalIP != "::1" {
		t.Errorf("expected local IP ::1, got %s", e.LocalIP)
	}
	if e.LocalPort != 443 {
		t.Errorf("expected local port 443, got %d", e.LocalPort)
	}
}

func TestParseLsofOutput_IPv6Wildcard(t *testing.T) {
	output := "p100\ncmyapp\nLroot\nf3\ntIPv6\nn*:8080\nTST=LISTEN\n"

	entries := parseLsofOutput(output)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].LocalIP != "::" {
		t.Errorf("expected :: for IPv6 wildcard, got %s", entries[0].LocalIP)
	}
}

func TestParseLsofOutput_IPv6UDP(t *testing.T) {
	output := "p200\ncresolved\nLroot\nf7\ntIPv6\nn[::1]:5353\n"

	entries := parseLsofOutput(output)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	e := entries[0]
	if e.Proto != "udp6" {
		t.Errorf("expected proto udp6, got %s", e.Proto)
	}
}

func TestParseLsofOutput_MultipleProcesses(t *testing.T) {
	output := "p100\ncnginx\nLroot\nf4\ntIPv4\nn*:80\nTST=LISTEN\n" +
		"f5\ntIPv4\nn*:443\nTST=LISTEN\n" +
		"p200\ncnode\nLuser\nf3\ntIPv4\nn127.0.0.1:3000\nTST=LISTEN\n"

	entries := parseLsofOutput(output)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	// First two belong to PID 100
	if entries[0].PID != 100 || entries[1].PID != 100 {
		t.Errorf("expected first two entries from PID 100")
	}
	if entries[0].LocalPort != 80 || entries[1].LocalPort != 443 {
		t.Errorf("expected ports 80 and 443 for PID 100")
	}

	// Third belongs to PID 200
	if entries[2].PID != 200 {
		t.Errorf("expected PID 200, got %d", entries[2].PID)
	}
}

func TestParseLsofOutput_Empty(t *testing.T) {
	entries := parseLsofOutput("")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for empty input, got %d", len(entries))
	}
}

func TestParseLsofOutput_IgnoresNonNetworkFiles(t *testing.T) {
	// A regular file entry (tREG) should be skipped
	output := "p100\nccat\nLuser\nf3\ntREG\nn/etc/passwd\n"

	entries := parseLsofOutput(output)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for non-network file, got %d", len(entries))
	}
}

func TestParseLsofOutput_IgnoresExtraTFields(t *testing.T) {
	// Only TST= should be extracted; TQR=, TQS= should be ignored
	output := "p100\ncapp\nLroot\nf3\ntIPv4\nn*:8080\nTST=LISTEN\nTQR=5\nTQS=10\n"

	entries := parseLsofOutput(output)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].State != "LISTEN" {
		t.Errorf("expected state LISTEN, got %s", entries[0].State)
	}
}

func TestParseLsofOutput_BoundAddress(t *testing.T) {
	output := "p100\ncapp\nLroot\nf3\ntIPv4\nn192.168.1.1:9090\nTST=LISTEN\n"

	entries := parseLsofOutput(output)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].LocalIP != "192.168.1.1" {
		t.Errorf("expected local IP 192.168.1.1, got %s", entries[0].LocalIP)
	}
}

func TestListSockets_Integration(t *testing.T) {
	// Skip if lsof is not available
	if _, err := exec.LookPath("lsof"); err != nil {
		t.Skip("lsof not available")
	}

	entries, err := ListSockets(NetFilterOptions{})
	if err != nil {
		t.Skipf("ListSockets failed (may need permissions): %v", err)
	}

	// Just verify it returns without error and entries have valid shape
	for _, e := range entries {
		if e.Proto == "" {
			t.Errorf("entry has empty proto")
		}
		if e.PID == 0 {
			t.Errorf("entry has zero PID")
		}
	}
}
