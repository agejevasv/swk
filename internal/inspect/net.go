package inspect

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
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

// NetSocketsJSON returns JSON-encoded output for the socket list.
func NetSocketsJSON(entries []SocketEntry) ([]byte, error) {
	return json.MarshalIndent(entries, "", "  ")
}

// SortSocketEntries sorts entries by protocol then local port.
func SortSocketEntries(entries []SocketEntry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Proto != entries[j].Proto {
			return entries[i].Proto < entries[j].Proto
		}
		return entries[i].LocalPort < entries[j].LocalPort
	})
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
