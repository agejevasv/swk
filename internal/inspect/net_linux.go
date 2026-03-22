//go:build linux

package inspect

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SocketEntry represents a single network socket with process info.
type SocketEntry struct {
	Proto      string `json:"proto"`
	LocalIP    string `json:"local_ip"`
	LocalPort  uint16 `json:"local_port"`
	RemoteIP   string `json:"remote_ip"`
	RemotePort uint16 `json:"remote_port"`
	State      string `json:"state"`
	PID        int    `json:"pid"`
	Process    string `json:"process"`
	User       string `json:"user"`
	Service    string `json:"service,omitempty"`
	Container  string `json:"container,omitempty"`
}

// NetFilterOptions controls which sockets are returned.
type NetFilterOptions struct {
	All  bool
	TCP  bool
	UDP  bool
	Port int
}

var tcpStateNames = map[uint8]string{
	0x01: "ESTABLISHED",
	0x02: "SYN_SENT",
	0x03: "SYN_RECV",
	0x04: "FIN_WAIT1",
	0x05: "FIN_WAIT2",
	0x06: "TIME_WAIT",
	0x07: "CLOSE",
	0x08: "CLOSE_WAIT",
	0x09: "LAST_ACK",
	0x0A: "LISTEN",
	0x0B: "CLOSING",
}

// socketRaw holds parsed data from a /proc/net/* line before PID resolution.
type socketRaw struct {
	proto      string
	localIP    net.IP
	localPort  uint16
	remoteIP   net.IP
	remotePort uint16
	state      uint8
	uid        uint32
	inode      uint64
}

// ListSockets scans /proc for socket information and returns filtered entries.
// procRoot allows tests to supply a fake /proc tree (default: "/proc").
func ListSockets(procRoot string, opts NetFilterOptions) ([]SocketEntry, error) {
	showTCP := opts.TCP || (!opts.TCP && !opts.UDP)
	showUDP := opts.UDP || (!opts.TCP && !opts.UDP)

	var raws []socketRaw

	if showTCP {
		for _, f := range []struct {
			name  string
			proto string
			ipv6  bool
		}{
			{"net/tcp", "tcp", false},
			{"net/tcp6", "tcp6", true},
		} {
			path := filepath.Join(procRoot, f.name)
			entries, err := parseProcNet(path, f.proto, f.ipv6)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return nil, err
			}
			raws = append(raws, entries...)
		}
	}

	if showUDP {
		for _, f := range []struct {
			name  string
			proto string
			ipv6  bool
		}{
			{"net/udp", "udp", false},
			{"net/udp6", "udp6", true},
		} {
			path := filepath.Join(procRoot, f.name)
			entries, err := parseProcNet(path, f.proto, f.ipv6)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return nil, err
			}
			raws = append(raws, entries...)
		}
	}

	// Filter by state
	var filtered []socketRaw
	for _, r := range raws {
		if !opts.All {
			isTCP := r.proto == "tcp" || r.proto == "tcp6"
			isUDP := r.proto == "udp" || r.proto == "udp6"
			if isTCP && r.state != 0x0A {
				continue
			}
			if isUDP && r.state != 0x07 {
				continue
			}
		}
		if opts.Port > 0 && int(r.localPort) != opts.Port {
			continue
		}
		filtered = append(filtered, r)
	}

	// Build inode→PID map
	inodeMap := buildInodeMap(procRoot)

	// Query Docker for container port mappings
	dockerPorts := queryDockerPorts()

	// Build entries
	entries := make([]SocketEntry, 0, len(filtered))
	for _, r := range filtered {
		pid := inodeMap[r.inode]
		proc, user := getProcessInfo(procRoot, pid, r.uid)

		entries = append(entries, SocketEntry{
			Proto:      r.proto,
			LocalIP:    r.localIP.String(),
			LocalPort:  r.localPort,
			RemoteIP:   r.remoteIP.String(),
			RemotePort: r.remotePort,
			State:      stateName(r.state),
			PID:        pid,
			Process:    proc,
			User:       user,
			Service:    resolveService(r.localPort),
			Container:  dockerPorts[r.localPort],
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Proto != entries[j].Proto {
			return entries[i].Proto < entries[j].Proto
		}
		return entries[i].LocalPort < entries[j].LocalPort
	})

	return entries, nil
}

func stateName(s uint8) string {
	if name, ok := tcpStateNames[s]; ok {
		return name
	}
	return fmt.Sprintf("0x%02X", s)
}

func parseProcNet(path string, proto string, isIPv6 bool) ([]socketRaw, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) < 2 {
		return nil, nil
	}

	var entries []socketRaw
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		localIP, localPort, err := parseHexIPPort(fields[1], isIPv6)
		if err != nil {
			continue
		}

		remoteIP, remotePort, err := parseHexIPPort(fields[2], isIPv6)
		if err != nil {
			continue
		}

		state, err := strconv.ParseUint(fields[3], 16, 8)
		if err != nil {
			continue
		}

		uid, err := strconv.ParseUint(fields[7], 10, 32)
		if err != nil {
			continue
		}

		inode, err := strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			continue
		}

		entries = append(entries, socketRaw{
			proto:      proto,
			localIP:    localIP,
			localPort:  localPort,
			remoteIP:   remoteIP,
			remotePort: remotePort,
			state:      uint8(state),
			uid:        uint32(uid),
			inode:      inode,
		})
	}

	return entries, nil
}

func parseHexIPPort(s string, isIPv6 bool) (net.IP, uint16, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return nil, 0, fmt.Errorf("invalid address:port %q", s)
	}

	port, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid port %q: %w", parts[1], err)
	}

	ipHex := parts[0]
	ipBytes, err := hex.DecodeString(ipHex)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid IP hex %q: %w", ipHex, err)
	}

	if isIPv6 {
		if len(ipBytes) != 16 {
			return nil, 0, fmt.Errorf("invalid IPv6 length: %d", len(ipBytes))
		}
		// Four 4-byte groups, each in little-endian
		for i := 0; i < 16; i += 4 {
			binary.BigEndian.PutUint32(ipBytes[i:i+4], binary.LittleEndian.Uint32(ipBytes[i:i+4]))
		}
		return net.IP(ipBytes), uint16(port), nil
	}

	if len(ipBytes) != 4 {
		return nil, 0, fmt.Errorf("invalid IPv4 length: %d", len(ipBytes))
	}
	// Little-endian → network order
	reverseBytes(ipBytes)
	return net.IP(ipBytes), uint16(port), nil
}

func reverseBytes(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func buildInodeMap(procRoot string) map[uint64]int {
	inodeMap := make(map[uint64]int, 256)

	entries, err := os.ReadDir(procRoot)
	if err != nil {
		return inodeMap
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		fdDir := filepath.Join(procRoot, entry.Name(), "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			if !strings.HasPrefix(link, "socket:[") || !strings.HasSuffix(link, "]") {
				continue
			}
			inodeStr := link[8 : len(link)-1]
			inode, err := strconv.ParseUint(inodeStr, 10, 64)
			if err != nil {
				continue
			}
			inodeMap[inode] = pid
		}
	}

	return inodeMap
}

func getProcessInfo(procRoot string, pid int, uid uint32) (name, user string) {
	if pid == 0 {
		return "", resolveUser(uid)
	}

	pidDir := filepath.Join(procRoot, strconv.Itoa(pid))

	commBytes, err := os.ReadFile(filepath.Join(pidDir, "comm"))
	if err == nil {
		name = strings.TrimSpace(string(commBytes))
	}

	user = resolveUser(uid)
	return name, user
}

// dockerContainer represents relevant fields from the Docker API response.
type dockerContainer struct {
	Names []string `json:"Names"`
	Image string   `json:"Image"`
	Ports []struct {
		PublicPort uint16 `json:"PublicPort"`
	} `json:"Ports"`
}

// queryDockerPorts queries the Docker socket for container port mappings.
// Returns a map of host port → "container_name (image)" or empty map on failure.
func queryDockerPorts() map[uint16]string {
	return queryDockerPortsFromSocket("/var/run/docker.sock")
}

func queryDockerPortsFromSocket(socketPath string) map[uint16]string {
	result := make(map[uint16]string)

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.DialTimeout("unix", socketPath, 300*time.Millisecond)
			},
		},
		Timeout: 300 * time.Millisecond,
	}

	resp, err := client.Get("http://localhost/containers/json")
	if err != nil {
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result
	}

	var containers []dockerContainer
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return result
	}

	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		label := name
		if c.Image != "" {
			if label != "" {
				label += " (" + c.Image + ")"
			} else {
				label = c.Image
			}
		}

		for _, p := range c.Ports {
			if p.PublicPort > 0 {
				result[p.PublicPort] = label
			}
		}
	}

	return result
}

var (
	passwdOnce  sync.Once
	passwdCache map[string]string // uid → username
)

func loadPasswd() {
	passwdCache = make(map[string]string)
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.SplitN(line, ":", 4)
		if len(fields) >= 3 {
			passwdCache[fields[2]] = fields[0]
		}
	}
}

func resolveUser(uid uint32) string {
	passwdOnce.Do(loadPasswd)
	uidStr := strconv.FormatUint(uint64(uid), 10)
	if name, ok := passwdCache[uidStr]; ok {
		return name
	}
	return uidStr
}

var wellKnownPorts = map[uint16]string{
	21:    "ftp",
	22:    "ssh",
	23:    "telnet",
	25:    "smtp",
	53:    "dns",
	80:    "http",
	110:   "pop3",
	143:   "imap",
	443:   "https",
	465:   "smtps",
	587:   "submission",
	993:   "imaps",
	995:   "pop3s",
	1080:  "socks",
	1433:  "mssql",
	1521:  "oracle",
	2181:  "zookeeper",
	2379:  "etcd",
	3000:  "grafana",
	3306:  "mysql",
	4222:  "nats",
	5432:  "postgres",
	5672:  "amqp",
	5900:  "vnc",
	6379:  "redis",
	6443:  "kube-api",
	8080:  "http-alt",
	8443:  "https-alt",
	8500:  "consul",
	8888:  "http-alt",
	9090:  "prometheus",
	9092:  "kafka",
	9200:  "elasticsearch",
	9300:  "elasticsearch",
	11211: "memcached",
	15672: "rabbitmq-mgmt",
	27017: "mongodb",
	27018: "mongodb",
	28015: "rethinkdb",
	50051: "grpc",
}

func resolveService(port uint16) string {
	if s, ok := wellKnownPorts[port]; ok {
		return s
	}
	return ""
}

// NetSocketsJSON returns JSON-encoded output for the socket list.
func NetSocketsJSON(entries []SocketEntry) ([]byte, error) {
	return json.MarshalIndent(entries, "", "  ")
}
