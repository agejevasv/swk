package inspect

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

// --- NetSocketsJSON tests ---

func TestNetSocketsJSON(t *testing.T) {
	entries := []SocketEntry{
		{
			Proto:     "tcp",
			LocalIP:   "127.0.0.1",
			LocalPort: 80,
			RemoteIP:  "0.0.0.0",
			State:     "LISTEN",
			PID:       1234,
			Process:   "nginx",
			User:      "www-data",
		},
	}

	out, err := NetSocketsJSON(entries)
	if err != nil {
		t.Fatal(err)
	}

	s := string(out)
	for _, want := range []string{`"proto": "tcp"`, `"local_port": 80`, `"process": "nginx"`} {
		if !strings.Contains(s, want) {
			t.Errorf("expected JSON to contain %q, got %s", want, s)
		}
	}
}

func TestNetSocketsJSON_Empty(t *testing.T) {
	out, err := NetSocketsJSON([]SocketEntry{})
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "[]" {
		t.Errorf("expected [], got %s", string(out))
	}
}

// --- resolveService tests ---

func TestResolveService_Known(t *testing.T) {
	tests := []struct {
		port uint16
		want string
	}{
		{22, "ssh"},
		{80, "http"},
		{443, "https"},
		{3306, "mysql"},
		{5432, "postgres"},
		{6379, "redis"},
		{9200, "elasticsearch"},
		{27017, "mongodb"},
	}
	for _, tt := range tests {
		got := resolveService(tt.port)
		if got != tt.want {
			t.Errorf("resolveService(%d) = %q, want %q", tt.port, got, tt.want)
		}
	}
}

func TestResolveService_Unknown(t *testing.T) {
	got := resolveService(59999)
	if got != "" {
		t.Errorf("expected empty string for unknown port, got %q", got)
	}
}

// --- SortSocketEntries tests ---

func TestSortSocketEntries(t *testing.T) {
	entries := []SocketEntry{
		{Proto: "udp", LocalPort: 53},
		{Proto: "tcp", LocalPort: 443},
		{Proto: "tcp", LocalPort: 80},
	}
	SortSocketEntries(entries)
	if entries[0].Proto != "tcp" || entries[0].LocalPort != 80 {
		t.Errorf("expected tcp:80 first, got %s:%d", entries[0].Proto, entries[0].LocalPort)
	}
	if entries[1].Proto != "tcp" || entries[1].LocalPort != 443 {
		t.Errorf("expected tcp:443 second, got %s:%d", entries[1].Proto, entries[1].LocalPort)
	}
	if entries[2].Proto != "udp" || entries[2].LocalPort != 53 {
		t.Errorf("expected udp:53 third, got %s:%d", entries[2].Proto, entries[2].LocalPort)
	}
}

// --- queryDockerPortsFromSocket tests ---

func TestQueryDockerPorts_FakeServer(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "docker.sock")

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/containers/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"Names": ["/my-es"],
				"Image": "wgroup/es:7.17.9",
				"Ports": [{"PublicPort": 9200}, {"PublicPort": 9300}]
			},
			{
				"Names": ["/my-redis"],
				"Image": "redis:7",
				"Ports": [{"PublicPort": 6379}]
			}
		]`)
	})
	go http.Serve(ln, mux)

	result := queryDockerPortsFromSocket(socketPath)

	if result[9200] != "my-es (wgroup/es:7.17.9)" {
		t.Errorf("expected 'my-es (wgroup/es:7.17.9)' for port 9200, got %q", result[9200])
	}
	if result[9300] != "my-es (wgroup/es:7.17.9)" {
		t.Errorf("expected 'my-es (wgroup/es:7.17.9)' for port 9300, got %q", result[9300])
	}
	if result[6379] != "my-redis (redis:7)" {
		t.Errorf("expected 'my-redis (redis:7)' for port 6379, got %q", result[6379])
	}
}

func TestQueryDockerPorts_NoSocket(t *testing.T) {
	result := queryDockerPortsFromSocket("/nonexistent/docker.sock")
	if len(result) != 0 {
		t.Errorf("expected empty map for missing socket, got %d entries", len(result))
	}
}

func TestQueryDockerPorts_BadResponse(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "docker.sock")

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/containers/json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	go http.Serve(ln, mux)

	result := queryDockerPortsFromSocket(socketPath)
	if len(result) != 0 {
		t.Errorf("expected empty map for 500 response, got %d entries", len(result))
	}
}

func TestQueryDockerPorts_NoPublicPorts(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "docker.sock")

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/containers/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"Names": ["/internal"], "Image": "app:latest", "Ports": [{"PublicPort": 0}]}]`)
	})
	go http.Serve(ln, mux)

	result := queryDockerPortsFromSocket(socketPath)
	if len(result) != 0 {
		t.Errorf("expected empty map for no public ports, got %d entries", len(result))
	}
}
