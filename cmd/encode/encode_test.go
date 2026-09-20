package encode

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/pflag"

	"github.com/agejevasv/swk/internal/ioutil"
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

func makeTestJWT(secret string, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func TestBase64_Encode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("base64", "hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.TrimSpace(out), "aGVsbG8gd29ybGQ=") {
		t.Errorf("expected 'aGVsbG8gd29ybGQ=', got %q", out)
	}
}

func TestBase64_Decode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("base64", "-d", "aGVsbG8gd29ybGQ=")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.TrimSpace(out), "hello world") {
		t.Errorf("expected 'hello world', got %q", out)
	}
}

func TestBase64_Roundtrip(t *testing.T) {
	t.Cleanup(resetAllFlags)
	// Encode
	encoded, err := executeCommand("base64", "roundtrip test")
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	encoded = strings.TrimSpace(encoded)

	resetAllFlags()

	// Decode
	decoded, err := executeCommand("base64", "-d", encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if strings.TrimSpace(decoded) != "roundtrip test" {
		t.Errorf("roundtrip failed: got %q", decoded)
	}
}

func TestHash_SHA256(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("hash", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) != 64 {
		t.Errorf("expected 64-char SHA-256 hash, got %d chars: %q", len(trimmed), trimmed)
	}
}

func TestHash_MD5(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("hash", "--algo", "md5", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) != 32 {
		t.Errorf("expected 32-char MD5 hash, got %d chars: %q", len(trimmed), trimmed)
	}
}

func TestHash_Verify(t *testing.T) {
	t.Cleanup(resetAllFlags)
	// First compute hash
	out, err := executeCommand("hash", "hello")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	hash := strings.TrimSpace(out)

	resetAllFlags()

	// Verify it
	out, err = executeCommand("hash", "--verify", hash, "hello")
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "ok") && !strings.Contains(lower, "match") && !strings.Contains(lower, "✓") && !strings.Contains(lower, "valid") {
		t.Errorf("expected verification success in output, got %q", out)
	}
}

func TestJWT_Encode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("jwt", "--secret", "mykey", `{"sub":"user1","role":"admin"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parts := strings.Split(strings.TrimSpace(out), ".")
	if len(parts) != 3 {
		t.Errorf("expected 3 dot-separated parts, got %d: %q", len(parts), out)
	}
}

func TestJWT_Decode(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := makeTestJWT("testsecret", jwt.MapClaims{"sub": "user1"})
	out, err := executeCommand("jwt", "-d", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "sub") {
		t.Errorf("expected 'sub' in decoded output, got %q", out)
	}
}

func TestJWT_Roundtrip(t *testing.T) {
	t.Cleanup(resetAllFlags)
	// Encode
	out, err := executeCommand("jwt", "--secret", "roundkey", `{"sub":"roundtrip"}`)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	token := strings.TrimSpace(out)

	resetAllFlags()

	// Decode and verify
	out, err = executeCommand("jwt", "-d", "--secret", "roundkey", token)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !strings.Contains(out, "roundtrip") {
		t.Errorf("expected 'roundtrip' in decoded output, got %q", out)
	}
}

func TestBase64_URLSafe(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("base64", "--url-safe", "hello?world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if strings.ContainsAny(trimmed, "+/") {
		t.Errorf("expected URL-safe encoding without +/, got %q", trimmed)
	}
}

func TestBase64_NoPadding(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("base64", "--no-padding", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if strings.Contains(trimmed, "=") {
		t.Errorf("expected no padding, got %q", trimmed)
	}
}

func TestHash_VerifyFail(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("hash", "--verify", "0000000000000000000000000000000000000000000000000000000000000000", "hello")
	if err == nil {
		t.Fatal("expected error for failed verification, got nil")
	}
	if !strings.Contains(out, "FAILED") {
		t.Errorf("expected 'FAILED' in output, got %q", out)
	}
}

func TestHash_SHA512(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("hash", "--algo", "sha512", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	if len(trimmed) != 128 {
		t.Errorf("expected 128-char SHA-512 hash, got %d chars: %q", len(trimmed), trimmed)
	}
}

func TestJWT_EncodeNoSecret(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("jwt", `{"sub":"user1"}`)
	if err == nil {
		t.Fatal("expected error when no secret provided, got nil")
	}
}

func TestJWT_DecodeInvalid(t *testing.T) {
	t.Cleanup(resetAllFlags)
	_, err := executeCommand("jwt", "-d", "not.a.jwt")
	if err == nil {
		t.Fatal("expected error for invalid JWT, got nil")
	}
}

func TestJWT_VerifyWrongSecret(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := makeTestJWT("correct", jwt.MapClaims{"sub": "u"})
	out, err := executeCommand("jwt", "-d", "--secret", "wrong", token)
	// Should still decode but Valid should be false
	if err != nil {
		// Some implementations error on bad sig
		return
	}
	if strings.Contains(out, `"valid": true`) {
		t.Errorf("expected valid=false, got %q", out)
	}
}

func TestQR_Terminal(t *testing.T) {
	t.Cleanup(resetAllFlags)
	out, err := executeCommand("qr", "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) == "" {
		t.Error("expected non-empty QR output")
	}
}

// A failed signature check must exit 1, not 0.
func TestJWT_VerifyFailureExitsOne(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := makeTestJWT("right-secret", jwt.MapClaims{"sub": "u1"})

	out, err := executeCommand("jwt", "-d", "--secret", "wrong-secret", token)
	if err == nil {
		t.Fatal("expected an error for a bad signature")
	}

	var ec ioutil.ExitCoder
	if !errors.As(err, &ec) {
		t.Fatalf("expected an ExitCoder, got %v", err)
	}
	if ec.ExitCode() != 1 {
		t.Errorf("exit code = %d, want 1", ec.ExitCode())
	}
	if !strings.Contains(out, `"valid": false`) {
		t.Errorf("expected valid:false in output, got %q", out)
	}
}

func TestJWT_VerifySuccessExitsZero(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := makeTestJWT("right-secret", jwt.MapClaims{"sub": "u1"})

	out, err := executeCommand("jwt", "-d", "--secret", "right-secret", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"valid": true`) {
		t.Errorf("expected valid:true in output, got %q", out)
	}
}

// A plain decode verifies nothing, so it must not report a validity verdict.
func TestJWT_DecodeOmitsValidField(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := makeTestJWT("secret", jwt.MapClaims{"sub": "u1"})

	out, err := executeCommand("jwt", "-d", token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, `"valid"`) {
		t.Errorf("decode output should not claim a validity verdict, got %q", out)
	}
}

// Accepting both would let the token's own alg header choose the key.
func TestJWT_SecretAndKeyAreMutuallyExclusive(t *testing.T) {
	t.Cleanup(resetAllFlags)
	token := makeTestJWT("secret", jwt.MapClaims{"sub": "u1"})

	if _, err := executeCommand("jwt", "-d", "--secret", "s", "--key", "k.pem", token); err == nil {
		t.Fatal("expected an error when both --secret and --key are given")
	}
}
