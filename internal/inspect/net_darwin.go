//go:build darwin

package inspect

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ListSockets lists network sockets by parsing lsof -F output on macOS.
func ListSockets(opts NetFilterOptions) ([]SocketEntry, error) {
	args := []string{"-n", "-P", "-FpcLnTt"}

	switch {
	case opts.TCP && !opts.UDP:
		args = append(args, "-iTCP")
	case opts.UDP && !opts.TCP:
		args = append(args, "-iUDP")
	default:
		args = append(args, "-i")
	}

	out, err := exec.Command("lsof", args...).Output()
	if err != nil {
		// lsof exits 1 when no results found — not an error for us.
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, fmt.Errorf("lsof: %w", err)
	}

	entries := parseLsofOutput(string(out))

	// Apply filters
	var filtered []SocketEntry
	for _, e := range entries {
		if !opts.All {
			isTCP := e.Proto == "tcp" || e.Proto == "tcp6"
			if isTCP && e.State != "LISTEN" {
				continue
			}
			// UDP has no state — always show bound UDP sockets in default mode.
		}
		if opts.Port > 0 && int(e.LocalPort) != opts.Port {
			continue
		}
		filtered = append(filtered, e)
	}

	// Enrich with service names and Docker info
	dockerPorts := queryDockerPorts()
	for i := range filtered {
		filtered[i].Service = resolveService(filtered[i].LocalPort)
		filtered[i].Container = dockerPorts[filtered[i].LocalPort]
	}

	SortSocketEntries(filtered)

	return filtered, nil
}

// lsofProcess accumulates per-process fields while parsing lsof -F output.
type lsofProcess struct {
	pid  int
	comm string
	user string
}

// lsofFile accumulates per-fd fields while parsing lsof -F output.
type lsofFile struct {
	fileType string // "IPv4" or "IPv6"
	addr     string // raw n-field value
	state    string // TCP state from TST= field, empty for UDP
}

// parseLsofOutput parses the machine-readable output of lsof -FpcLnTt.
func parseLsofOutput(output string) []SocketEntry {
	var entries []SocketEntry
	var proc lsofProcess
	var file lsofFile
	hasFile := false

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		field := line[0]
		value := line[1:]

		switch field {
		case 'p':
			// New process — emit previous file if any
			if hasFile {
				if e, ok := buildLsofEntry(proc, file); ok {
					entries = append(entries, e)
				}
				hasFile = false
			}
			pid, _ := strconv.Atoi(value)
			proc = lsofProcess{pid: pid}
			file = lsofFile{}

		case 'c':
			proc.comm = value

		case 'L':
			proc.user = value

		case 'f':
			// New file descriptor — emit previous file if any
			if hasFile {
				if e, ok := buildLsofEntry(proc, file); ok {
					entries = append(entries, e)
				}
			}
			file = lsofFile{}
			hasFile = true

		case 't':
			file.fileType = value

		case 'n':
			file.addr = value

		case 'T':
			// T field values: "ST=LISTEN", "QR=0", "QS=0", etc.
			if strings.HasPrefix(value, "ST=") {
				file.state = value[3:]
			}
		}
	}

	// Emit the last file
	if hasFile {
		if e, ok := buildLsofEntry(proc, file); ok {
			entries = append(entries, e)
		}
	}

	return entries
}

// buildLsofEntry converts accumulated lsof fields into a SocketEntry.
// Returns false if the file is not a network socket we care about.
func buildLsofEntry(proc lsofProcess, file lsofFile) (SocketEntry, bool) {
	if file.fileType != "IPv4" && file.fileType != "IPv6" {
		return SocketEntry{}, false
	}
	if file.addr == "" {
		return SocketEntry{}, false
	}

	isIPv6 := file.fileType == "IPv6"
	hasTCPState := file.state != ""

	var proto string
	switch {
	case !isIPv6 && hasTCPState:
		proto = "tcp"
	case isIPv6 && hasTCPState:
		proto = "tcp6"
	case !isIPv6 && !hasTCPState:
		proto = "udp"
	case isIPv6 && !hasTCPState:
		proto = "udp6"
	}

	localIP, localPort, remoteIP, remotePort := parseLsofAddress(file.addr, isIPv6)

	return SocketEntry{
		Proto:      proto,
		LocalIP:    localIP,
		LocalPort:  localPort,
		RemoteIP:   remoteIP,
		RemotePort: remotePort,
		State:      file.state,
		PID:        proc.pid,
		Process:    proc.comm,
		User:       proc.user,
	}, true
}

// parseLsofAddress parses an lsof n-field network address.
// Formats: "*:port", "addr:port", "addr:port->remoteaddr:remoteport", "[::1]:port"
func parseLsofAddress(s string, isIPv6 bool) (localIP string, localPort uint16, remoteIP string, remotePort uint16) {
	local := s
	remote := ""

	if idx := strings.Index(s, "->"); idx >= 0 {
		local = s[:idx]
		remote = s[idx+2:]
	}

	localIP, localPort = splitHostPort(local, isIPv6)
	if remote != "" {
		remoteIP, remotePort = splitHostPort(remote, isIPv6)
	}

	return
}

// splitHostPort splits "addr:port" handling IPv6 brackets and wildcards.
func splitHostPort(s string, isIPv6 bool) (string, uint16) {
	// Handle bracketed IPv6: [::1]:port
	if strings.HasPrefix(s, "[") {
		if idx := strings.LastIndex(s, "]:"); idx >= 0 {
			host := s[1:idx]
			port := parsePort(s[idx+2:])
			return host, port
		}
		return s, 0
	}

	// For IPv6 without brackets: find last colon
	// For IPv4: find last colon
	idx := strings.LastIndex(s, ":")
	if idx < 0 {
		return s, 0
	}

	host := s[:idx]
	port := parsePort(s[idx+1:])

	// Map wildcard to appropriate zero address
	if host == "*" {
		if isIPv6 {
			host = "::"
		} else {
			host = "0.0.0.0"
		}
	}

	return host, port
}

func parsePort(s string) uint16 {
	n, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0
	}
	return uint16(n)
}
