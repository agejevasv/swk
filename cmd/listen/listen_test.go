package listen

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/spf13/pflag"
)

func resetAllFlags() {
	Cmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Value.Set(f.DefValue)
		f.Changed = false
	})
}

func TestListen_RequestLogging(t *testing.T) {
	t.Cleanup(resetAllFlags)

	// Capture stderr to get the port and logged requests
	pr, pw, _ := os.Pipe()
	Cmd.SetOut(new(devNull))
	Cmd.SetErr(pw)
	Cmd.SetArgs([]string{"--port", "0", "--status", "201", "--body", `{"ok":true}`})

	errCh := make(chan error, 1)
	go func() {
		errCh <- Cmd.Execute()
	}()

	// Read startup message to get port
	buf := make([]byte, 512)
	n, _ := pr.Read(buf)
	msg := string(buf[:n])

	var port int
	for i := len(msg) - 1; i >= 0; i-- {
		if msg[i] == ':' {
			fmt.Sscanf(msg[i+1:], "%d", &port)
			break
		}
	}
	if port == 0 {
		t.Fatalf("could not parse port from startup message: %q", msg)
	}

	// Make a request
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Post(
		fmt.Sprintf("http://localhost:%d/test", port),
		"application/json",
		strings.NewReader(`{"hello":"world"}`),
	)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"ok":true}` {
		t.Errorf("expected response body '{\"ok\":true}', got %q", string(body))
	}

	// Shut down
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server did not shut down within timeout")
	}

	pw.Close()
	logged, _ := io.ReadAll(pr)
	pr.Close()

	logStr := string(logged)
	if !strings.Contains(logStr, "POST") {
		t.Errorf("expected 'POST' in logged output, got %q", logStr)
	}
	if !strings.Contains(logStr, "/test") {
		t.Errorf("expected '/test' in logged output, got %q", logStr)
	}
}

type devNull struct{}

func (devNull) Write(p []byte) (int, error) { return len(p), nil }

// An out-of-range status used to reach net/http and panic on every request.
func TestListen_RejectsInvalidStatus(t *testing.T) {
	for _, status := range []string{"-1", "0", "99", "700", "1000"} {
		cmd := *Cmd
		cmd.SetOut(new(bytes.Buffer))
		cmd.SetErr(new(bytes.Buffer))
		cmd.SetArgs([]string{"--status", status, "--port", "0"})
		if err := cmd.Execute(); err == nil {
			t.Errorf("--status %s should have been rejected", status)
		}
		cmd.Flags().Set("status", "200")
	}
}

func TestListen_RejectsInvalidPort(t *testing.T) {
	for _, port := range []string{"-1", "65536", "99999"} {
		cmd := *Cmd
		cmd.SetOut(new(bytes.Buffer))
		cmd.SetErr(new(bytes.Buffer))
		cmd.SetArgs([]string{"--port", port})
		if err := cmd.Execute(); err == nil {
			t.Errorf("--port %s should have been rejected", port)
		}
		cmd.Flags().Set("port", "8080")
	}
}
