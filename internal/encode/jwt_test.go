package encode

import (
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

func TestJWTVerify(t *testing.T) {
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
			name:      "wrong_secret_returns_error",
			token:     validFutureToken,
			secret:    "wrong-secret",
			wantValid: false,
			wantErr:   true,
		},
		{
			name:      "expired_token_returns_error",
			token:     expiredToken,
			secret:    secret,
			wantValid: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := JWTVerify(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Fatalf("JWTVerify() error = %v, wantErr %v", err, tt.wantErr)
			}
			if info != nil && info.Valid != tt.wantValid {
				t.Errorf("JWTVerify() valid = %v, want %v", info.Valid, tt.wantValid)
			}
		})
	}
}

func TestJWTVerify_WrongSecret_StillReturnsInfo(t *testing.T) {
	secret := "correct-secret"
	token := createTestJWT(t, jwt.MapClaims{
		"sub": "user1",
		"exp": float64(time.Now().Add(time.Hour).Unix()),
	}, secret)

	info, err := JWTVerify(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
	// Info should still be populated even on verification failure.
	if info == nil {
		t.Fatal("expected info to be non-nil even on verification failure")
	}
	if info.Valid {
		t.Error("expected Valid=false for wrong secret")
	}
	if info.Header == nil {
		t.Error("expected Header to be populated")
	}
}

func TestJWTEncode(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		secret  string
		algo    string
		wantErr bool
		checkFn func(t *testing.T, token string)
	}{
		{
			name:    "basic HS256",
			payload: `{"sub":"user1","role":"admin"}`,
			secret:  "test-secret",
			algo:    "HS256",
			checkFn: func(t *testing.T, token string) {
				info, err := JWTDecode(token)
				if err != nil {
					t.Fatalf("failed to decode created token: %v", err)
				}
				if info.Payload["sub"] != "user1" {
					t.Errorf("expected sub=user1, got %v", info.Payload["sub"])
				}
				if info.Payload["role"] != "admin" {
					t.Errorf("expected role=admin, got %v", info.Payload["role"])
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
				info, err := JWTDecode(token)
				if err != nil {
					t.Fatalf("failed to decode: %v", err)
				}
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
				info, err := JWTDecode(token)
				if err != nil {
					t.Fatalf("failed to decode: %v", err)
				}
				if info.Header["alg"] != "HS512" {
					t.Errorf("expected alg=HS512, got %v", info.Header["alg"])
				}
			},
		},
		{
			name:    "roundtrip: encode then verify",
			payload: `{"sub":"roundtrip","exp":9999999999}`,
			secret:  "my-secret",
			algo:    "HS256",
			checkFn: func(t *testing.T, token string) {
				info, err := JWTVerify(token, "my-secret")
				if err != nil {
					t.Fatalf("verify failed: %v", err)
				}
				if !info.Valid {
					t.Error("expected Valid=true")
				}
				if info.Payload["sub"] != "roundtrip" {
					t.Errorf("expected sub=roundtrip, got %v", info.Payload["sub"])
				}
			},
		},
		{
			name:    "invalid JSON payload",
			payload: `not json`,
			secret:  "secret",
			algo:    "HS256",
			wantErr: true,
		},
		{
			name:    "unsupported algorithm",
			payload: `{"sub":"test"}`,
			secret:  "secret",
			algo:    "RS256",
			wantErr: true,
		},
		{
			name:    "empty secret still signs",
			payload: `{"sub":"test"}`,
			secret:  "",
			algo:    "HS256",
			checkFn: func(t *testing.T, token string) {
				if token == "" {
					t.Error("expected non-empty token")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := JWTEncode(tt.payload, tt.secret, tt.algo)
			if (err != nil) != tt.wantErr {
				t.Fatalf("JWTEncode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFn != nil && err == nil {
				tt.checkFn(t, token)
			}
		})
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
