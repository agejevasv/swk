package inspect

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"strings"
	"testing"
	"time"

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

func executeCommandWithStdin(stdin string, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetErr(buf)
	Cmd.SetIn(strings.NewReader(stdin))
	Cmd.SetArgs(args)
	err := Cmd.Execute()
	return buf.String(), err
}

func generateTestCert(notBefore, notAfter time.Time) string {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "test.example.com",
		},
		DNSNames:  []string{"test.example.com"},
		NotBefore: notBefore,
		NotAfter:  notAfter,
		KeyUsage:  x509.KeyUsageDigitalSignature,
	}
	certDER, _ := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})
	return string(certPEM)
}

func TestCert_Decode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	now := time.Now()
	certPEM := generateTestCert(now.Add(-1*time.Hour), now.Add(24*time.Hour))
	out, err := executeCommandWithStdin(certPEM, "cert")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test.example.com") {
		t.Errorf("expected 'test.example.com' in output, got %q", out)
	}
}

func TestCert_CheckExpiry_Expired(t *testing.T) {
	t.Cleanup(resetAllFlags)
	now := time.Now()
	certPEM := generateTestCert(now.Add(-48*time.Hour), now.Add(-24*time.Hour))
	_, err := executeCommandWithStdin(certPEM, "cert", "--check-expiry")
	if err == nil {
		t.Fatal("expected error for expired cert, got nil")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "expired") {
		t.Errorf("expected 'expired' in error message, got %q", err.Error())
	}
}

func TestCron_Default(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("cron", "*/5 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Next") {
		t.Errorf("expected 'Next' in output, got %q", out)
	}
}

func TestCron_ExplainOnly(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("cron", "--explain", "0 9 * * 1-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "Next") {
		t.Errorf("expected no 'Next' in explain-only output, got %q", out)
	}
}

func TestCron_Invalid(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("cron", "not-a-cron")
	if err == nil {
		t.Fatal("expected error for invalid cron expression, got nil")
	}
}

func TestText_Inspect(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("text", "Hello World")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Characters:") {
		t.Errorf("expected 'Characters:' in output, got %q", out)
	}
	if !strings.Contains(out, "Words:") {
		t.Errorf("expected 'Words:' in output, got %q", out)
	}
}

func TestText_InspectJSON(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("text", "--json", "Hello World")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]interface{}
	if jsonErr := json.Unmarshal([]byte(out), &result); jsonErr != nil {
		t.Errorf("expected valid JSON output, got parse error: %v\noutput: %q", jsonErr, out)
	}
}
