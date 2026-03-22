//go:build linux

package inspect

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
)

type fakePID struct {
	comm   string
	uid    uint32
	inodes []uint64
}

func buildFakeProc(t *testing.T, nets map[string]string, pids map[int]fakePID) string {
	t.Helper()
	root := t.TempDir()

	// Create net/ files
	netDir := filepath.Join(root, "net")
	if err := os.MkdirAll(netDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range nets {
		if err := os.WriteFile(filepath.Join(netDir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// Create PID directories
	for pid, info := range pids {
		pidDir := filepath.Join(root, fmt.Sprintf("%d", pid))
		fdDir := filepath.Join(pidDir, "fd")
		if err := os.MkdirAll(fdDir, 0o755); err != nil {
			t.Fatal(err)
		}

		// comm
		if err := os.WriteFile(filepath.Join(pidDir, "comm"), []byte(info.comm+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		// status with Uid line
		statusContent := fmt.Sprintf("Name:\t%s\nUid:\t%d\t%d\t%d\t%d\n", info.comm, info.uid, info.uid, info.uid, info.uid)
		if err := os.WriteFile(filepath.Join(pidDir, "status"), []byte(statusContent), 0o644); err != nil {
			t.Fatal(err)
		}

		// fd symlinks
		for i, inode := range info.inodes {
			link := filepath.Join(fdDir, fmt.Sprintf("%d", i+3))
			if err := os.Symlink(fmt.Sprintf("socket:[%d]", inode), link); err != nil {
				t.Fatal(err)
			}
		}
	}

	return root
}

// --- parseHexIPPort tests ---

func TestParseHexIPPort_IPv4Loopback(t *testing.T) {
	ip, port, err := parseHexIPPort("0100007F:0050", false)
	if err != nil {
		t.Fatal(err)
	}
	if !ip.Equal(net.ParseIP("127.0.0.1")) {
		t.Errorf("expected 127.0.0.1, got %s", ip)
	}
	if port != 80 {
		t.Errorf("expected port 80, got %d", port)
	}
}

func TestParseHexIPPort_IPv4Any(t *testing.T) {
	ip, port, err := parseHexIPPort("00000000:1F90", false)
	if err != nil {
		t.Fatal(err)
	}
	if !ip.Equal(net.IPv4zero) {
		t.Errorf("expected 0.0.0.0, got %s", ip)
	}
	if port != 8080 {
		t.Errorf("expected port 8080, got %d", port)
	}
}

func TestParseHexIPPort_IPv4Regular(t *testing.T) {
	// 192.168.1.100 → in little-endian hex: 6401A8C0
	ip, port, err := parseHexIPPort("6401A8C0:0016", false)
	if err != nil {
		t.Fatal(err)
	}
	if !ip.Equal(net.ParseIP("192.168.1.100")) {
		t.Errorf("expected 192.168.1.100, got %s", ip)
	}
	if port != 22 {
		t.Errorf("expected port 22, got %d", port)
	}
}

func TestParseHexIPPort_IPv6Loopback(t *testing.T) {
	// ::1 in /proc format: 00000000000000000000000001000000
	ip, port, err := parseHexIPPort("00000000000000000000000001000000:0050", true)
	if err != nil {
		t.Fatal(err)
	}
	if !ip.Equal(net.ParseIP("::1")) {
		t.Errorf("expected ::1, got %s", ip)
	}
	if port != 80 {
		t.Errorf("expected port 80, got %d", port)
	}
}

func TestParseHexIPPort_IPv6Any(t *testing.T) {
	ip, port, err := parseHexIPPort("00000000000000000000000000000000:1F90", true)
	if err != nil {
		t.Fatal(err)
	}
	if !ip.Equal(net.IPv6zero) {
		t.Errorf("expected ::, got %s", ip)
	}
	if port != 8080 {
		t.Errorf("expected port 8080, got %d", port)
	}
}

func TestParseHexIPPort_InvalidHex(t *testing.T) {
	_, _, err := parseHexIPPort("ZZZZZZZZ:0050", false)
	if err == nil {
		t.Fatal("expected error for invalid hex")
	}
}

func TestParseHexIPPort_MissingColon(t *testing.T) {
	_, _, err := parseHexIPPort("0100007F0050", false)
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseHexIPPort_InvalidPort(t *testing.T) {
	_, _, err := parseHexIPPort("0100007F:ZZZZ", false)
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestParseHexIPPort_WrongIPv4Length(t *testing.T) {
	_, _, err := parseHexIPPort("0100007F00:0050", false)
	if err == nil {
		t.Fatal("expected error for wrong IPv4 length")
	}
}

func TestParseHexIPPort_WrongIPv6Length(t *testing.T) {
	_, _, err := parseHexIPPort("0100007F:0050", true)
	if err == nil {
		t.Fatal("expected error for wrong IPv6 length")
	}
}

// --- parseProcNet tests ---

func TestParseProcNet_TCPListen(t *testing.T) {
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 12345 1 0000000000000000 100 0 0 10 0`

	root := t.TempDir()
	netDir := filepath.Join(root, "net")
	os.MkdirAll(netDir, 0o755)
	os.WriteFile(filepath.Join(netDir, "tcp"), []byte(content), 0o644)

	entries, err := parseProcNet(filepath.Join(netDir, "tcp"), "tcp", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	e := entries[0]
	if e.proto != "tcp" {
		t.Errorf("expected proto tcp, got %s", e.proto)
	}
	if !e.localIP.Equal(net.ParseIP("127.0.0.1")) {
		t.Errorf("expected 127.0.0.1, got %s", e.localIP)
	}
	if e.localPort != 80 {
		t.Errorf("expected port 80, got %d", e.localPort)
	}
	if e.state != 0x0A {
		t.Errorf("expected state 0x0A, got 0x%02X", e.state)
	}
	if e.uid != 1000 {
		t.Errorf("expected uid 1000, got %d", e.uid)
	}
	if e.inode != 12345 {
		t.Errorf("expected inode 12345, got %d", e.inode)
	}
}

func TestParseProcNet_Established(t *testing.T) {
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0050 0100007F:C000 01 00000000:00000000 00:00000000 00000000  1000        0 12345 1 0000000000000000 100 0 0 10 0`

	root := t.TempDir()
	netDir := filepath.Join(root, "net")
	os.MkdirAll(netDir, 0o755)
	os.WriteFile(filepath.Join(netDir, "tcp"), []byte(content), 0o644)

	entries, err := parseProcNet(filepath.Join(netDir, "tcp"), "tcp", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].state != 0x01 {
		t.Errorf("expected state 0x01 (ESTABLISHED), got 0x%02X", entries[0].state)
	}
	if entries[0].remotePort != 0xC000 {
		t.Errorf("expected remote port %d, got %d", 0xC000, entries[0].remotePort)
	}
}

func TestParseProcNet_HeaderOnly(t *testing.T) {
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode`

	root := t.TempDir()
	netDir := filepath.Join(root, "net")
	os.MkdirAll(netDir, 0o755)
	os.WriteFile(filepath.Join(netDir, "tcp"), []byte(content), 0o644)

	entries, err := parseProcNet(filepath.Join(netDir, "tcp"), "tcp", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for header-only, got %d", len(entries))
	}
}

func TestParseProcNet_MalformedLine(t *testing.T) {
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 12345 1 0000000000000000 100 0 0 10 0
   bad line here
   2: 0100007F:0051 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 12346 1 0000000000000000 100 0 0 10 0`

	root := t.TempDir()
	netDir := filepath.Join(root, "net")
	os.MkdirAll(netDir, 0o755)
	os.WriteFile(filepath.Join(netDir, "tcp"), []byte(content), 0o644)

	entries, err := parseProcNet(filepath.Join(netDir, "tcp"), "tcp", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries (skipping malformed), got %d", len(entries))
	}
}

func TestParseProcNet_FileNotExist(t *testing.T) {
	_, err := parseProcNet("/nonexistent/path", "tcp", false)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

// --- buildInodeMap tests ---

func TestBuildInodeMap_Basic(t *testing.T) {
	root := buildFakeProc(t, nil, map[int]fakePID{
		1234: {comm: "nginx", uid: 33, inodes: []uint64{5555, 6666}},
	})

	m := buildInodeMap(root)
	if m[5555] != 1234 {
		t.Errorf("expected inode 5555 → PID 1234, got %d", m[5555])
	}
	if m[6666] != 1234 {
		t.Errorf("expected inode 6666 → PID 1234, got %d", m[6666])
	}
}

func TestBuildInodeMap_MultiplePIDs(t *testing.T) {
	root := buildFakeProc(t, nil, map[int]fakePID{
		100: {comm: "a", uid: 0, inodes: []uint64{111}},
		200: {comm: "b", uid: 0, inodes: []uint64{222}},
	})

	m := buildInodeMap(root)
	if m[111] != 100 {
		t.Errorf("expected inode 111 → PID 100, got %d", m[111])
	}
	if m[222] != 200 {
		t.Errorf("expected inode 222 → PID 200, got %d", m[222])
	}
}

func TestBuildInodeMap_EmptyProc(t *testing.T) {
	root := t.TempDir()
	m := buildInodeMap(root)
	if len(m) != 0 {
		t.Errorf("expected empty map, got %d entries", len(m))
	}
}

// --- ListSockets integration tests ---

func TestListSockets_ListenOnly(t *testing.T) {
	tcp := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 11111 1 0000000000000000 100 0 0 10 0
   1: 0100007F:0051 0100007F:C000 01 00000000:00000000 00:00000000 00000000  1000        0 22222 1 0000000000000000 100 0 0 10 0`

	root := buildFakeProc(t, map[string]string{"tcp": tcp}, map[int]fakePID{
		42: {comm: "myapp", uid: 1000, inodes: []uint64{11111}},
	})

	entries, err := listSocketsFrom(root, NetFilterOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 LISTEN entry, got %d", len(entries))
	}
	if entries[0].State != "LISTEN" {
		t.Errorf("expected LISTEN state, got %s", entries[0].State)
	}
	if entries[0].PID != 42 {
		t.Errorf("expected PID 42, got %d", entries[0].PID)
	}
	if entries[0].Process != "myapp" {
		t.Errorf("expected process myapp, got %s", entries[0].Process)
	}
}

func TestListSockets_AllStates(t *testing.T) {
	tcp := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 11111 1 0000000000000000 100 0 0 10 0
   1: 0100007F:0051 0100007F:C000 01 00000000:00000000 00:00000000 00000000  1000        0 22222 1 0000000000000000 100 0 0 10 0`

	root := buildFakeProc(t, map[string]string{"tcp": tcp}, nil)

	entries, err := listSocketsFrom(root, NetFilterOptions{All: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries with --all, got %d", len(entries))
	}
}

func TestListSockets_TCPFilter(t *testing.T) {
	tcp := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 11111 1 0000000000000000 100 0 0 10 0`
	udp := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode ref pointer drops
   0: 0100007F:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000  1000        0 33333 1 0000000000000000 100 0 0 10 0`

	root := buildFakeProc(t, map[string]string{"tcp": tcp, "udp": udp}, nil)

	entries, err := listSocketsFrom(root, NetFilterOptions{TCP: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 TCP entry, got %d", len(entries))
	}
	if entries[0].Proto != "tcp" {
		t.Errorf("expected proto tcp, got %s", entries[0].Proto)
	}
}

func TestListSockets_UDPFilter(t *testing.T) {
	tcp := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 11111 1 0000000000000000 100 0 0 10 0`
	udp := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode ref pointer drops
   0: 0100007F:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000  1000        0 33333 1 0000000000000000 100 0 0 10 0`

	root := buildFakeProc(t, map[string]string{"tcp": tcp, "udp": udp}, nil)

	entries, err := listSocketsFrom(root, NetFilterOptions{UDP: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 UDP entry, got %d", len(entries))
	}
	if entries[0].Proto != "udp" {
		t.Errorf("expected proto udp, got %s", entries[0].Proto)
	}
}

func TestListSockets_PortFilter(t *testing.T) {
	tcp := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 11111 1 0000000000000000 100 0 0 10 0
   1: 0100007F:01BB 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 22222 1 0000000000000000 100 0 0 10 0`

	root := buildFakeProc(t, map[string]string{"tcp": tcp}, nil)

	entries, err := listSocketsFrom(root, NetFilterOptions{Port: 80})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry for port 80, got %d", len(entries))
	}
	if entries[0].LocalPort != 80 {
		t.Errorf("expected port 80, got %d", entries[0].LocalPort)
	}
}

func TestListSockets_NoMatchingPID(t *testing.T) {
	tcp := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 99999 1 0000000000000000 100 0 0 10 0`

	root := buildFakeProc(t, map[string]string{"tcp": tcp}, nil)

	entries, err := listSocketsFrom(root, NetFilterOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].PID != 0 {
		t.Errorf("expected PID 0 for unmapped inode, got %d", entries[0].PID)
	}
	if entries[0].Process != "" {
		t.Errorf("expected empty process, got %q", entries[0].Process)
	}
}

func TestListSockets_SortedOutput(t *testing.T) {
	tcp := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:01BB 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 11111 1 0000000000000000 100 0 0 10 0
   1: 0100007F:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 22222 1 0000000000000000 100 0 0 10 0`

	root := buildFakeProc(t, map[string]string{"tcp": tcp}, nil)

	entries, err := listSocketsFrom(root, NetFilterOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].LocalPort >= entries[1].LocalPort {
		t.Errorf("expected sorted by port: %d < %d", entries[0].LocalPort, entries[1].LocalPort)
	}
}

func TestListSockets_UDPDefaultState(t *testing.T) {
	udp := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode ref pointer drops
   0: 0100007F:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000  1000        0 11111 1 0000000000000000 100 0 0 10 0
   1: 0100007F:0036 0100007F:C000 01 00000000:00000000 00:00000000 00000000  1000        0 22222 1 0000000000000000 100 0 0 10 0`

	root := buildFakeProc(t, map[string]string{"udp": udp}, nil)

	entries, err := listSocketsFrom(root, NetFilterOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 UDP entry (state 07 only), got %d", len(entries))
	}
	if entries[0].State != "CLOSE" {
		t.Errorf("expected CLOSE state, got %s", entries[0].State)
	}
}

func TestListSockets_MissingNetFiles(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "net"), 0o755)

	entries, err := listSocketsFrom(root, NetFilterOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

// --- stateName tests ---

func TestStateName_Known(t *testing.T) {
	tests := []struct {
		state uint8
		want  string
	}{
		{0x01, "ESTABLISHED"},
		{0x0A, "LISTEN"},
		{0x06, "TIME_WAIT"},
		{0x07, "CLOSE"},
	}
	for _, tt := range tests {
		got := stateName(tt.state)
		if got != tt.want {
			t.Errorf("stateName(0x%02X) = %q, want %q", tt.state, got, tt.want)
		}
	}
}

func TestStateName_Unknown(t *testing.T) {
	got := stateName(0xFF)
	if got != "0xFF" {
		t.Errorf("expected 0xFF, got %s", got)
	}
}

// --- service field via listSocketsFrom ---

func TestListSockets_ServiceField(t *testing.T) {
	tcp := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 11111 1 0000000000000000 100 0 0 10 0
   1: 0100007F:1538 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 22222 1 0000000000000000 100 0 0 10 0`

	root := buildFakeProc(t, map[string]string{"tcp": tcp}, nil)

	entries, err := listSocketsFrom(root, NetFilterOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	for _, e := range entries {
		if e.LocalPort == 80 && e.Service != "http" {
			t.Errorf("expected service http for port 80, got %q", e.Service)
		}
		if e.LocalPort == 5432 && e.Service != "postgres" {
			t.Errorf("expected service postgres for port 5432, got %q", e.Service)
		}
	}
}
