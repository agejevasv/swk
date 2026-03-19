package encode

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/pflag"
)

func resetAllFlags() {
	// Reset package-level vars bound via XxxVarP
	base64Decode = false
	base64URLSafe = false
	base64NoPadding = false
	urlDecode = false
	urlComponent = false
	htmlDecode = false
	gzipDecode = false
	gzipLevel = 6
	jwtDecode = false
	jwtSecret = ""
	jwtAlgo = "HS256"
	certCheckExpiry = false
	qrOutput = "terminal"
	qrSize = 256
	qrLevel = "M"

	// Also reset cobra flag Changed state on all subcommands
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

func executeCommandBytes(args ...string) ([]byte, error) {
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetErr(buf)
	Cmd.SetArgs(args)
	err := Cmd.Execute()
	return buf.Bytes(), err
}

// ── base64 ─────────────────────────────────────────────────────────────────

func TestBase64_Encode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("base64", "hello world")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if trimmed != "aGVsbG8gd29ybGQ=" {
		t.Errorf("expected 'aGVsbG8gd29ybGQ=', got %q", trimmed)
	}
}

func TestBase64_Decode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("base64", "-d", "aGVsbG8gd29ybGQ=")
	if err != nil {
		t.Fatal(err)
	}
	if out != "hello world" {
		t.Errorf("expected 'hello world', got %q", out)
	}
}

func TestBase64_URLSafe(t *testing.T) {
	t.Cleanup(resetAllFlags)
	// Encode bytes that produce +/ in standard base64
	// The bytes 0xfb, 0xef, 0xbe produce "++++=" in standard base64
	// Use a string that produces + or / in standard encoding
	input := "subjects?_d" // standard: c3ViamVjdHM/X2Q= contains ? which encodes with /
	out, err := executeCommand("base64", "--url-safe", input)
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	// URL-safe encoding should not contain + or /
	if strings.ContainsAny(trimmed, "+/") {
		t.Errorf("expected URL-safe encoding without + or /, got %q", trimmed)
	}
}

func TestBase64_NoPadding(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("base64", "--no-padding", "hello world")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if strings.Contains(trimmed, "=") {
		t.Errorf("expected no padding characters, got %q", trimmed)
	}
	if trimmed != "aGVsbG8gd29ybGQ" {
		t.Errorf("expected 'aGVsbG8gd29ybGQ', got %q", trimmed)
	}
}

func TestBase64_Roundtrip(t *testing.T) {
	t.Cleanup(resetAllFlags)
	original := "The quick brown fox jumps over the lazy dog"
	encoded, err := executeCommand("base64", original)
	if err != nil {
		t.Fatal(err)
	}
	resetAllFlags()
	decoded, err := executeCommand("base64", "-d", strings.TrimSpace(encoded))
	if err != nil {
		t.Fatal(err)
	}
	if decoded != original {
		t.Errorf("roundtrip failed: expected %q, got %q", original, decoded)
	}
}

// ── url ────────────────────────────────────────────────────────────────────

func TestURL_Encode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("url", "hello world&foo=bar")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if !strings.Contains(trimmed, "%") {
		t.Errorf("expected URL-encoded output with %% chars, got %q", trimmed)
	}
	// Space should be encoded
	if strings.Contains(trimmed, " ") {
		t.Errorf("expected spaces to be encoded, got %q", trimmed)
	}
}

func TestURL_Decode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("url", "-d", "hello%20world%26foo%3Dbar")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if trimmed != "hello world&foo=bar" {
		t.Errorf("expected 'hello world&foo=bar', got %q", trimmed)
	}
}

func TestURL_ComponentFlag(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("url", "--component", "hello world")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	// Component encoding uses + for spaces (QueryEscape)
	if !strings.Contains(trimmed, "+") {
		t.Errorf("expected component encoding with + for spaces, got %q", trimmed)
	}
}

func TestURL_ComponentDecode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("url", "-d", "--component", "hello+world")
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if trimmed != "hello world" {
		t.Errorf("expected 'hello world', got %q", trimmed)
	}
}

// ── html ───────────────────────────────────────────────────────────────────

func TestHTML_Encode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("html", `<div class="test">&</div>`)
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if !strings.Contains(trimmed, "&lt;") {
		t.Errorf("expected &lt; in output, got %q", trimmed)
	}
	if !strings.Contains(trimmed, "&amp;") {
		t.Errorf("expected &amp; in output, got %q", trimmed)
	}
	if !strings.Contains(trimmed, "&gt;") {
		t.Errorf("expected &gt; in output, got %q", trimmed)
	}
}

func TestHTML_Decode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("html", "-d", `&lt;div class=&quot;test&quot;&gt;&amp;&lt;/div&gt;`)
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSpace(out)
	if !strings.Contains(trimmed, "<div") {
		t.Errorf("expected decoded HTML with '<div', got %q", trimmed)
	}
	if !strings.Contains(trimmed, "&") {
		t.Errorf("expected decoded '&', got %q", trimmed)
	}
}

// ── gzip ───────────────────────────────────────────────────────────────────

func TestGzip_Roundtrip(t *testing.T) {
	t.Cleanup(resetAllFlags)
	original := "Hello, gzip compression test!"

	compressedBytes, err := executeCommandBytes("gzip", original)
	if err != nil {
		t.Fatal(err)
	}
	if len(compressedBytes) == 0 {
		t.Fatal("expected non-empty compressed output")
	}

	resetAllFlags()

	// Decompress using raw bytes via stdin
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetErr(buf)
	Cmd.SetIn(bytes.NewReader(compressedBytes))
	Cmd.SetArgs([]string{"gzip", "-d"})
	err = Cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != original {
		t.Errorf("gzip roundtrip failed: expected %q, got %q", original, buf.String())
	}
}

// ── jwt ────────────────────────────────────────────────────────────────────

func makeTestJWT(secret string, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func TestJWT_Encode(t *testing.T) {
	t.Cleanup(resetAllFlags)

	out, err := executeCommand("jwt", "--secret", "mykey", `{"sub":"user1","role":"admin"}`)
	if err != nil {
		t.Fatal(err)
	}
	out = strings.TrimSpace(out)
	parts := strings.Split(out, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3-part JWT, got %d parts: %s", len(parts), out)
	}
}

func TestJWT_EncodeRequiresSecret(t *testing.T) {
	t.Cleanup(resetAllFlags)

	_, err := executeCommand("jwt", `{"sub":"test"}`)
	if err == nil {
		t.Error("expected error when encoding without --secret")
	}
}

func TestJWT_Decode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	secret := "test-secret"
	tokenStr := makeTestJWT(secret, jwt.MapClaims{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  1516239022,
	})

	out, err := executeCommand("jwt", "-d", tokenStr)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"sub"`) {
		t.Errorf("expected JWT payload with 'sub' claim, got:\n%s", out)
	}
	if !strings.Contains(out, `"John Doe"`) {
		t.Errorf("expected JWT payload with 'John Doe', got:\n%s", out)
	}
	if !strings.Contains(out, `"header"`) {
		t.Errorf("expected 'header' section in output, got:\n%s", out)
	}
}

func TestJWT_VerifyWithCorrectSecret(t *testing.T) {
	t.Cleanup(resetAllFlags)
	secret := "my-secret-key"
	tokenStr := makeTestJWT(secret, jwt.MapClaims{
		"sub": "user1",
	})

	out, err := executeCommand("jwt", "-d", "--secret", secret, tokenStr)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"valid": true`) {
		t.Errorf("expected valid=true with correct secret, got:\n%s", out)
	}
}

func TestJWT_VerifyWithWrongSecret(t *testing.T) {
	t.Cleanup(resetAllFlags)
	tokenStr := makeTestJWT("correct-secret", jwt.MapClaims{
		"sub": "user1",
	})

	_, err := executeCommand("jwt", "-d", "--secret", "wrong-secret", tokenStr)
	if err == nil {
		t.Error("expected error with wrong secret")
	}
}

func TestJWT_Roundtrip(t *testing.T) {
	t.Cleanup(resetAllFlags)

	// Encode
	out, err := executeCommand("jwt", "--secret", "roundtrip-key", `{"sub":"test","role":"admin"}`)
	if err != nil {
		t.Fatal(err)
	}
	token := strings.TrimSpace(out)

	resetAllFlags()

	// Decode and verify
	out, err = executeCommand("jwt", "-d", "--secret", "roundtrip-key", token)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"valid": true`) {
		t.Errorf("expected valid=true, got:\n%s", out)
	}
	if !strings.Contains(out, `"sub"`) {
		t.Errorf("expected sub claim, got:\n%s", out)
	}
}

// ── cert ───────────────────────────────────────────────────────────────────

func generateTestCert(t *testing.T, notBefore, notAfter time.Time) string {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "test.example.com",
			Organization: []string{"Test Org"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
		DNSNames:  []string{"test.example.com", "*.example.com"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	return buf.String()
}

func TestCert_Decode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	certPEM := generateTestCert(t, time.Now().Add(-24*time.Hour), time.Now().Add(365*24*time.Hour))

	// Pass cert via stdin since it's multiline
	Cmd.SetIn(strings.NewReader(certPEM))
	out, err := executeCommand("cert")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "test.example.com") {
		t.Errorf("expected subject with 'test.example.com', got:\n%s", out)
	}
	if !strings.Contains(out, `"is_expired": false`) {
		t.Errorf("expected is_expired=false, got:\n%s", out)
	}
	if !strings.Contains(out, "*.example.com") {
		t.Errorf("expected DNS name '*.example.com', got:\n%s", out)
	}
}

func TestCert_CheckExpiry_Valid(t *testing.T) {
	t.Cleanup(resetAllFlags)
	certPEM := generateTestCert(t, time.Now().Add(-24*time.Hour), time.Now().Add(365*24*time.Hour))

	Cmd.SetIn(strings.NewReader(certPEM))
	out, err := executeCommand("cert", "--check-expiry")
	if err != nil {
		t.Errorf("expected no error for valid cert with --check-expiry, got: %v", err)
	}
	if !strings.Contains(out, `"is_expired": false`) {
		t.Errorf("expected is_expired=false, got:\n%s", out)
	}
}

func TestCert_CheckExpiry_Expired(t *testing.T) {
	t.Cleanup(resetAllFlags)
	certPEM := generateTestCert(t, time.Now().Add(-48*time.Hour), time.Now().Add(-24*time.Hour))

	Cmd.SetIn(strings.NewReader(certPEM))
	_, err := executeCommand("cert", "--check-expiry")
	if err == nil {
		t.Error("expected error for expired cert with --check-expiry")
	}
	if !strings.Contains(err.Error(), "expired") {
		t.Errorf("expected error message about expiry, got: %v", err)
	}
}

// ── qr ─────────────────────────────────────────────────────────────────────

func TestQR_TerminalOutput(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("qr", "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Error("expected non-empty QR terminal output")
	}
	// Terminal QR output should contain block characters (unicode block elements)
	hasBlocks := strings.ContainsAny(out, "\u2580\u2584\u2588\u2591\u2592\u2593\u2596\u2597\u2598\u2599\u259A\u259B\u259C\u259D\u259E\u259F\u2503\u2501\u250F\u2513\u2517\u251B")
	if !hasBlocks {
		// Some QR renderers use different characters; just check it has
		// non-ASCII content or at least substantial output
		if len(out) < 50 {
			t.Errorf("expected substantial QR output with block characters, got:\n%s", out)
		}
	}
}

func TestQR_InvalidLevel(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("qr", "--level", "Z", "hello")
	if err == nil {
		t.Error("expected error for invalid QR level")
	}
	if !strings.Contains(err.Error(), "invalid QR error correction level") {
		t.Errorf("expected specific error message, got: %v", err)
	}
}
