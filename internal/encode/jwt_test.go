package encode

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createTestJWT(t *testing.T, claims jwt.MapClaims, secret string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to create test JWT: %v", err)
	}
	return tokenStr
}

func encodePEMPrivateKey(t *testing.T, key interface{}) []byte {
	t.Helper()
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("failed to marshal private key: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
}

func encodePEMPublicKey(t *testing.T, key interface{}) []byte {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
}

func TestJWTDecode(t *testing.T) {
	secret := "test-secret"

	validToken := createTestJWT(t, jwt.MapClaims{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  float64(1516239022),
	}, secret)

	expiredToken := createTestJWT(t, jwt.MapClaims{
		"sub": "123",
		"exp": float64(time.Now().Add(-time.Hour).Unix()),
	}, secret)

	allClaimsToken := createTestJWT(t, jwt.MapClaims{
		"sub":  "user123",
		"iss":  "test-issuer",
		"aud":  "test-audience",
		"exp":  float64(time.Now().Add(time.Hour).Unix()),
		"nbf":  float64(time.Now().Add(-time.Minute).Unix()),
		"iat":  float64(time.Now().Unix()),
		"jti":  "unique-id-123",
		"role": "admin",
	}, secret)

	tests := []struct {
		name    string
		token   string
		wantErr bool
		checkFn func(t *testing.T, info *JWTInfo)
	}{
		{
			name:  "valid_hs256_token",
			token: validToken,
			checkFn: func(t *testing.T, info *JWTInfo) {
				if info.Payload["sub"] != "1234567890" {
					t.Errorf("expected sub=1234567890, got %v", info.Payload["sub"])
				}
				if info.Payload["name"] != "John Doe" {
					t.Errorf("expected name=John Doe, got %v", info.Payload["name"])
				}
				if info.Header["alg"] != "HS256" {
					t.Errorf("expected alg=HS256, got %v", info.Header["alg"])
				}
			},
		},
		{
			name:  "expired_token_has_exp_field",
			token: expiredToken,
			checkFn: func(t *testing.T, info *JWTInfo) {
				if info.ExpiredAt == nil {
					t.Error("expected ExpiredAt to be set")
				}
				if info.ExpiredAt != nil && info.ExpiredAt.After(time.Now()) {
					t.Error("expected ExpiredAt to be in the past")
				}
			},
		},
		{
			name:  "token_with_all_standard_claims",
			token: allClaimsToken,
			checkFn: func(t *testing.T, info *JWTInfo) {
				if info.Payload["sub"] != "user123" {
					t.Errorf("expected sub=user123, got %v", info.Payload["sub"])
				}
				if info.Payload["iss"] != "test-issuer" {
					t.Errorf("expected iss=test-issuer, got %v", info.Payload["iss"])
				}
				if info.Payload["role"] != "admin" {
					t.Errorf("expected role=admin, got %v", info.Payload["role"])
				}
				if info.ExpiredAt == nil {
					t.Error("expected ExpiredAt to be set")
				}
			},
		},
		{
			name:    "malformed_not_three_parts",
			token:   "only.two",
			wantErr: true,
		},
		{
			name:    "completely_invalid",
			token:   "garbage",
			wantErr: true,
		},
		{
			name:    "empty_string",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := JWTDecode(tt.token)
			if (err != nil) != tt.wantErr {
				t.Fatalf("JWTDecode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFn != nil && info != nil {
				tt.checkFn(t, info)
			}
		})
	}
}

func TestJWTVerify_HMAC(t *testing.T) {
	secret := "my-secret-key"

	validFutureToken := createTestJWT(t, jwt.MapClaims{
		"sub": "123",
		"exp": float64(time.Now().Add(time.Hour).Unix()),
	}, secret)

	expiredToken := createTestJWT(t, jwt.MapClaims{
		"sub": "123",
		"exp": float64(time.Now().Add(-time.Hour).Unix()),
	}, secret)

	tests := []struct {
		name      string
		token     string
		secret    string
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "correct_secret_valid_token",
			token:     validFutureToken,
			secret:    secret,
			wantValid: true,
		},
		{
			name:      "wrong_secret_returns_invalid",
			token:     validFutureToken,
			secret:    "wrong-secret",
			wantValid: false,
		},
		{
			name:      "expired_token_returns_invalid",
			token:     expiredToken,
			secret:    secret,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := JWTVerify(tt.token, tt.secret, nil)
			if (err != nil) != tt.wantErr {
				t.Fatalf("JWTVerify() error = %v, wantErr %v", err, tt.wantErr)
			}
			if info != nil && info.Valid != tt.wantValid {
				t.Errorf("JWTVerify() valid = %v, want %v", info.Valid, tt.wantValid)
			}
		})
	}
}

func TestJWTEncode_HMAC(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		secret  string
		algo    string
		wantErr bool
		checkFn func(t *testing.T, token string)
	}{
		{
			name:    "basic_HS256",
			payload: `{"sub":"user1","role":"admin"}`,
			secret:  "test-secret",
			algo:    "HS256",
			checkFn: func(t *testing.T, token string) {
				info, err := JWTDecode(token)
				if err != nil {
					t.Fatalf("failed to decode: %v", err)
				}
				if info.Payload["sub"] != "user1" {
					t.Errorf("expected sub=user1, got %v", info.Payload["sub"])
				}
				if info.Header["alg"] != "HS256" {
					t.Errorf("expected alg=HS256, got %v", info.Header["alg"])
				}
			},
		},
		{
			name:    "HS384",
			payload: `{"sub":"test"}`,
			secret:  "secret",
			algo:    "HS384",
			checkFn: func(t *testing.T, token string) {
				info, _ := JWTDecode(token)
				if info.Header["alg"] != "HS384" {
					t.Errorf("expected alg=HS384, got %v", info.Header["alg"])
				}
			},
		},
		{
			name:    "HS512",
			payload: `{"sub":"test"}`,
			secret:  "secret",
			algo:    "HS512",
			checkFn: func(t *testing.T, token string) {
				info, _ := JWTDecode(token)
				if info.Header["alg"] != "HS512" {
					t.Errorf("expected alg=HS512, got %v", info.Header["alg"])
				}
			},
		},
		{
			name:    "roundtrip_encode_then_verify",
			payload: `{"sub":"roundtrip","exp":9999999999}`,
			secret:  "my-secret",
			algo:    "HS256",
			checkFn: func(t *testing.T, token string) {
				info, err := JWTVerify(token, "my-secret", nil)
				if err != nil {
					t.Fatalf("verify failed: %v", err)
				}
				if !info.Valid {
					t.Error("expected Valid=true")
				}
			},
		},
		{
			name:    "invalid_JSON_payload",
			payload: `not json`,
			secret:  "secret",
			algo:    "HS256",
			wantErr: true,
		},
		{
			name:    "unsupported_algorithm",
			payload: `{"sub":"test"}`,
			secret:  "secret",
			algo:    "NONE",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := JWTEncode(tt.payload, tt.secret, nil, tt.algo)
			if (err != nil) != tt.wantErr {
				t.Fatalf("JWTEncode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFn != nil && err == nil {
				tt.checkFn(t, token)
			}
		})
	}
}

func TestJWT_RSA_Roundtrip(t *testing.T) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	privPEM := encodePEMPrivateKey(t, privKey)
	pubPEM := encodePEMPublicKey(t, &privKey.PublicKey)

	// Sign with private key.
	token, err := JWTEncode(`{"sub":"rsa-test","exp":9999999999}`, "", privPEM, "RS256")
	if err != nil {
		t.Fatalf("JWTEncode RS256 failed: %v", err)
	}

	// Decode without verification.
	info, err := JWTDecode(token)
	if err != nil {
		t.Fatalf("JWTDecode failed: %v", err)
	}
	if info.Header["alg"] != "RS256" {
		t.Errorf("expected alg=RS256, got %v", info.Header["alg"])
	}
	if info.Payload["sub"] != "rsa-test" {
		t.Errorf("expected sub=rsa-test, got %v", info.Payload["sub"])
	}

	// Verify with public key.
	info, err = JWTVerify(token, "", pubPEM)
	if err != nil {
		t.Fatalf("JWTVerify RS256 failed: %v", err)
	}
	if !info.Valid {
		t.Error("expected Valid=true for correct public key")
	}

	// Verify with wrong key should return Valid=false.
	wrongKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	wrongPEM := encodePEMPublicKey(t, &wrongKey.PublicKey)
	info, err = JWTVerify(token, "", wrongPEM)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Valid {
		t.Error("expected Valid=false for wrong public key")
	}
}

func TestJWT_ECDSA_Roundtrip(t *testing.T) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate EC key: %v", err)
	}

	privPEM := encodePEMPrivateKey(t, privKey)
	pubPEM := encodePEMPublicKey(t, &privKey.PublicKey)

	token, err := JWTEncode(`{"sub":"ec-test","exp":9999999999}`, "", privPEM, "ES256")
	if err != nil {
		t.Fatalf("JWTEncode ES256 failed: %v", err)
	}

	info, err := JWTDecode(token)
	if err != nil {
		t.Fatalf("JWTDecode failed: %v", err)
	}
	if info.Header["alg"] != "ES256" {
		t.Errorf("expected alg=ES256, got %v", info.Header["alg"])
	}

	info, err = JWTVerify(token, "", pubPEM)
	if err != nil {
		t.Fatalf("JWTVerify ES256 failed: %v", err)
	}
	if !info.Valid {
		t.Error("expected Valid=true")
	}
}

func TestJWT_Ed25519_Roundtrip(t *testing.T) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate Ed25519 key: %v", err)
	}

	privPEM := encodePEMPrivateKey(t, privKey)
	pubPEM := encodePEMPublicKey(t, pubKey)

	token, err := JWTEncode(`{"sub":"ed-test","exp":9999999999}`, "", privPEM, "EdDSA")
	if err != nil {
		t.Fatalf("JWTEncode EdDSA failed: %v", err)
	}

	info, err := JWTDecode(token)
	if err != nil {
		t.Fatalf("JWTDecode failed: %v", err)
	}
	if info.Header["alg"] != "EdDSA" {
		t.Errorf("expected alg=EdDSA, got %v", info.Header["alg"])
	}

	info, err = JWTVerify(token, "", pubPEM)
	if err != nil {
		t.Fatalf("JWTVerify EdDSA failed: %v", err)
	}
	if !info.Valid {
		t.Error("expected Valid=true")
	}
}

func TestJWT_KeyErrors(t *testing.T) {
	// RS256 without key file should error.
	_, err := JWTEncode(`{"sub":"test"}`, "", nil, "RS256")
	if err == nil {
		t.Error("expected error for RS256 without key")
	}

	// HMAC without secret returns Valid=false (not an error — token is still parseable).
	hmacToken := createTestJWT(t, jwt.MapClaims{"sub": "test"}, "secret")
	info, err := JWTVerify(hmacToken, "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Valid {
		t.Error("expected Valid=false for HMAC verify without secret")
	}

	// Bad PEM data should error.
	_, err = JWTEncode(`{"sub":"test"}`, "", []byte("not a pem"), "RS256")
	if err == nil {
		t.Error("expected error for invalid PEM")
	}
}

func TestJWTInfoJSON(t *testing.T) {
	info := &JWTInfo{
		Header:  map[string]interface{}{"alg": "HS256", "typ": "JWT"},
		Payload: map[string]interface{}{"sub": "123"},
		Valid:   true,
	}
	out, err := JWTInfoJSON(info)
	if err != nil {
		t.Fatalf("JWTInfoJSON: %v", err)
	}
	if len(out) == 0 {
		t.Error("expected non-empty JSON output")
	}
}
