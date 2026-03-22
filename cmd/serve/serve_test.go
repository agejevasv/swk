package serve

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	genLib "github.com/agejevasv/swk/internal/gen"
	"github.com/spf13/pflag"
)

func resetAllFlags() {
	Cmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Value.Set(f.DefValue)
		f.Changed = false
	})
}

func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetErr(buf)
	Cmd.SetArgs(args)
	err := Cmd.Execute()
	return buf.String(), err
}

func TestServe_InvalidDir(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("/nonexistent/path")
	if err == nil {
		t.Fatal("expected error for nonexistent directory")
	}
}

func TestServe_FileNotDir(t *testing.T) {
	t.Cleanup(resetAllFlags)
	f := filepath.Join(t.TempDir(), "file.txt")
	os.WriteFile(f, []byte("hello"), 0o644)

	_, err := executeCommand(f)
	if err == nil {
		t.Fatal("expected error when path is a file, not a directory")
	}
}

func TestServe_TLS_MissingCert(t *testing.T) {
	t.Cleanup(resetAllFlags)
	dir := t.TempDir()
	_, err := executeCommand(dir, "--tls", "--cert", "/nonexistent/cert.pem", "--key", "/nonexistent/key.pem", "--port", "0")
	if err == nil {
		t.Fatal("expected error for missing TLS cert")
	}
}

func TestServe_TLS_Integration(t *testing.T) {
	t.Cleanup(resetAllFlags)

	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "hello.txt"), []byte("hello tls"), 0o644)

	// Generate cert
	result, err := genLib.GenerateCert(genLib.CertOptions{CN: "localhost"})
	if err != nil {
		t.Fatal(err)
	}
	certPath := filepath.Join(dir, "cert.pem")
	keyPath := filepath.Join(dir, "key.pem")
	os.WriteFile(certPath, result.CertPEM, 0o644)
	os.WriteFile(keyPath, result.KeyPEM, 0o600)

	// Start server in background
	Cmd.SetArgs([]string{dir, "--tls", "--cert", certPath, "--key", keyPath, "--port", "0", "--no-log"})
	Cmd.SetOut(new(devNull))

	// Capture stderr to get the port
	pr, pw, _ := os.Pipe()
	Cmd.SetErr(pw)

	errCh := make(chan error, 1)
	go func() {
		errCh <- Cmd.Execute()
	}()

	// Read startup message to get port
	buf := make([]byte, 512)
	n, _ := pr.Read(buf)
	msg := string(buf[:n])
	pw.Close()
	pr.Close()

	// Extract port from "Serving ... on https://[::]:PORT"
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

	// Make HTTPS request (skip cert verification since self-signed)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("https://localhost:%d/hello.txt", port))
	if err != nil {
		t.Fatalf("HTTPS request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "hello tls" {
		t.Errorf("expected 'hello tls', got %q", string(body))
	}
	if resp.TLS == nil {
		t.Error("expected TLS connection")
	}

	// Shut down the server cleanly via SIGINT
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("server returned unexpected error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server did not shut down within timeout")
	}
}

type devNull struct{}

func (devNull) Write(p []byte) (int, error) { return len(p), nil }
