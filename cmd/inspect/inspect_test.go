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

	"github.com/golang-jwt/jwt/v5"
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
	out, err := executeCommandWithStdin(certPEM, "cert", "--check-expiry")
	if err == nil {
		t.Fatal("expected error for expired cert, got nil")
	}
	// JSON output still goes to stdout; the error is silent (exit code only).
	if !strings.Contains(out, "test.example.com") {
		t.Errorf("expected JSON output on stdout, got %q", out)
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

func TestURL_Parse(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("url", "https://example.com:8080/api?page=1#top")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "example.com") {
		t.Errorf("expected 'example.com' in output, got %q", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected '8080' in output, got %q", out)
	}
}

func TestURL_Invalid(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("url", "://not-a-url")
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}

func TestCert_NotExpired(t *testing.T) {
	t.Cleanup(resetAllFlags)
	now := time.Now()
	certPEM := generateTestCert(now.Add(-1*time.Hour), now.Add(24*time.Hour))
	_, err := executeCommandWithStdin(certPEM, "cert", "--check-expiry")
	if err != nil {
		t.Fatalf("expected no error for valid cert, got %v", err)
	}
}

func TestCron_Next(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("cron", "--next", "3", "0 9 * * MON")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	count := 0
	for _, l := range lines {
		if strings.Contains(l, "202") {
			count++
		}
	}
	if count < 3 {
		t.Errorf("expected at least 3 next runs, got %d in output %q", count, out)
	}
}

func TestText_InspectJSON(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("text", "--json", "Hello World")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]any
	if jsonErr := json.Unmarshal([]byte(out), &result); jsonErr != nil {
		t.Errorf("expected valid JSON output, got parse error: %v\noutput: %q", jsonErr, out)
	}
}

func createTestJWT(t *testing.T, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("failed to create test JWT: %v", err)
	}
	return tokenStr
}

func TestJWT_BasicDecode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := createTestJWT(t, jwt.MapClaims{"sub": "user1", "role": "admin"})
	out, err := executeCommand("jwt", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "HS256") {
		t.Errorf("expected HS256 in output, got %q", out)
	}
	if !strings.Contains(out, "user1") {
		t.Errorf("expected user1 in output, got %q", out)
	}
	if !strings.Contains(out, "admin") {
		t.Errorf("expected admin in output, got %q", out)
	}
}

func TestJWT_ExpiredToken(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := createTestJWT(t, jwt.MapClaims{
		"sub": "user1",
		"exp": float64(time.Now().Add(-24 * time.Hour).Unix()),
	})
	out, err := executeCommand("jwt", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "(expired)") {
		t.Errorf("expected (expired) in output, got %q", out)
	}
}

func TestJWT_ValidToken(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := createTestJWT(t, jwt.MapClaims{
		"sub": "user1",
		"exp": float64(time.Now().Add(24 * time.Hour).Unix()),
	})
	out, err := executeCommand("jwt", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "(valid)") {
		t.Errorf("expected (valid) in output, got %q", out)
	}
}

func TestJWT_CheckExpiry_Expired(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := createTestJWT(t, jwt.MapClaims{
		"sub": "user1",
		"exp": float64(time.Now().Add(-24 * time.Hour).Unix()),
	})
	_, err := executeCommand("jwt", "--check-expiry", token)
	if err == nil {
		t.Fatal("expected error for expired token with --check-expiry")
	}
}

func TestJWT_CheckExpiry_Valid(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := createTestJWT(t, jwt.MapClaims{
		"sub": "user1",
		"exp": float64(time.Now().Add(24 * time.Hour).Unix()),
	})
	_, err := executeCommand("jwt", "--check-expiry", token)
	if err != nil {
		t.Fatalf("expected no error for valid token, got %v", err)
	}
}

func TestJWT_JSONOutput(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := createTestJWT(t, jwt.MapClaims{"sub": "user1"})
	out, err := executeCommand("jwt", "--json", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]any
	if jsonErr := json.Unmarshal([]byte(out), &result); jsonErr != nil {
		t.Errorf("expected valid JSON, got parse error: %v\noutput: %q", jsonErr, out)
	}
	if _, ok := result["header"]; !ok {
		t.Errorf("expected 'header' in JSON output")
	}
	if _, ok := result["payload"]; !ok {
		t.Errorf("expected 'payload' in JSON output")
	}
}

func TestJWT_ComplexClaims(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := createTestJWT(t, jwt.MapClaims{
		"sub":         "user1",
		"permissions": []string{"read", "write"},
		"metadata":    map[string]any{"tenant": "acme"},
	})
	out, err := executeCommand("jwt", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `["read","write"]`) {
		t.Errorf("expected JSON array in output, got %q", out)
	}
	if !strings.Contains(out, `{"tenant":"acme"}`) {
		t.Errorf("expected JSON object in output, got %q", out)
	}
}

func TestJWT_Stdin(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := createTestJWT(t, jwt.MapClaims{"sub": "piped"})
	out, err := executeCommandWithStdin(token, "jwt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "piped") {
		t.Errorf("expected 'piped' in output, got %q", out)
	}
}

func TestJWT_Invalid(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("jwt", "not-a-jwt")
	if err == nil {
		t.Fatal("expected error for invalid JWT")
	}
}

func TestJWT_TimestampClaims(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := createTestJWT(t, jwt.MapClaims{
		"sub": "user1",
		"iat": float64(1700000000),
		"nbf": float64(1700000000),
	})
	out, err := executeCommand("jwt", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "2023-11-14") {
		t.Errorf("expected formatted timestamp in output, got %q", out)
	}
}

func TestJWT_CheckExpiry_NoExp(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := createTestJWT(t, jwt.MapClaims{"sub": "user1"})
	_, err := executeCommand("jwt", "--check-expiry", token)
	if err != nil {
		t.Fatalf("expected no error when no exp claim, got %v", err)
	}
}

func TestJWT_ValueTypes(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := createTestJWT(t, jwt.MapClaims{
		"sub":     "user1",
		"score":   3.14,
		"active":  true,
		"nothing": nil,
		"count":   float64(42),
	})
	out, err := executeCommand("jwt", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "3.14") {
		t.Errorf("expected float 3.14 in output, got %q", out)
	}
	if !strings.Contains(out, "true") {
		t.Errorf("expected bool true in output, got %q", out)
	}
	if !strings.Contains(out, "null") {
		t.Errorf("expected null in output, got %q", out)
	}
	if !strings.Contains(out, "42") {
		t.Errorf("expected integer 42 in output, got %q", out)
	}
}

func TestJWT_AudArray(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := createTestJWT(t, jwt.MapClaims{
		"sub": "user1",
		"aud": []string{"api.example.com", "web.example.com"},
	})
	out, err := executeCommand("jwt", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "api.example.com") && !strings.Contains(out, "web.example.com") {
		t.Errorf("expected audience values in output, got %q", out)
	}
}

func TestJWT_RegisteredClaimsOrder(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := createTestJWT(t, jwt.MapClaims{
		"role": "admin",
		"sub":  "user1",
		"iss":  "test",
	})
	out, err := executeCommand("jwt", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	subIdx := strings.Index(out, "sub")
	issIdx := strings.Index(out, "iss")
	roleIdx := strings.Index(out, "role")
	if subIdx > issIdx {
		t.Errorf("expected sub before iss")
	}
	if issIdx > roleIdx {
		t.Errorf("expected registered claims before custom claims")
	}
}
