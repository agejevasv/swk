package encode

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/agejevasv/swk/internal/jsonx"
)

// JWTInfo holds decoded JWT information.
type JWTInfo struct {
	Header    map[string]any `json:"header"`
	Payload   map[string]any `json:"payload"`
	Signature string         `json:"signature"`
	// Valid is nil when the token was decoded without signature verification,
	// so "not verified" is never reported as "invalid".
	Valid *bool `json:"valid,omitempty"`
	// Error holds why verification failed, when it failed.
	Error     string     `json:"error,omitempty"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
}

// IsVerified reports whether signature verification ran and succeeded.
func (i *JWTInfo) IsVerified() bool {
	return i != nil && i.Valid != nil && *i.Valid
}

func JWTEncode(payloadJSON string, secret string, keyPEM []byte, algo string) (string, error) {
	// Decoding through encoding/json would turn every claim into a float64 and
	// bake the rounding into a signed token.
	var claims jwt.MapClaims
	if err := jsonx.Decode([]byte(payloadJSON), &claims); err != nil {
		return "", fmt.Errorf("invalid JSON payload: %w", err)
	}

	method := jwt.GetSigningMethod(algo)
	if method == nil {
		return "", fmt.Errorf("unsupported algorithm: %s", algo)
	}

	token := jwt.NewWithClaims(method, claims)

	signingKey, err := resolveSigningKey(method, secret, keyPEM)
	if err != nil {
		return "", err
	}

	tokenStr, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenStr, nil
}

func JWTDecode(tokenStr string) (*JWTInfo, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation(), jwt.WithJSONNumber())
	token, parts, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("invalid JWT: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	sig := ""
	if len(parts) == 3 {
		sig = parts[2]
	}

	info := &JWTInfo{
		Header:    token.Header,
		Payload:   map[string]any(claims),
		Signature: sig,
	}

	extractExpiry(claims, info)

	return info, nil
}

func JWTVerify(tokenStr string, secret string, keyPEM []byte) (*JWTInfo, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return resolveVerifyKey(token.Method, secret, keyPEM)
	}, jwt.WithJSONNumber())

	sig := ""
	if parts := strings.SplitN(tokenStr, ".", 3); len(parts) == 3 {
		sig = parts[2]
	}

	info := &JWTInfo{Signature: sig}

	if token == nil {
		return nil, fmt.Errorf("invalid JWT: %w", err)
	}

	info.Header = token.Header
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		info.Payload = map[string]any(claims)
		extractExpiry(claims, info)
	}

	valid := token.Valid
	info.Valid = &valid
	if !valid {
		info.Error = "signature verification failed"
		if err != nil {
			info.Error = err.Error()
		}
	}

	return info, nil
}

func JWTInfoJSON(info *JWTInfo) ([]byte, error) {
	return json.MarshalIndent(info, "", "  ")
}

// resolveSigningKey returns the appropriate signing key for the given method.
func resolveSigningKey(method jwt.SigningMethod, secret string, keyPEM []byte) (any, error) {
	switch method.(type) {
	case *jwt.SigningMethodHMAC:
		if secret == "" {
			return nil, fmt.Errorf("--secret is required for %s", method.Alg())
		}
		return []byte(secret), nil
	case *jwt.SigningMethodRSA:
		return parsePrivateKey(keyPEM, "RSA")
	case *jwt.SigningMethodECDSA:
		return parsePrivateKey(keyPEM, "EC")
	case *jwt.SigningMethodEd25519:
		return parsePrivateKey(keyPEM, "Ed25519")
	default:
		return nil, fmt.Errorf("unsupported signing method: %s", method.Alg())
	}
}

// resolveVerifyKey returns the appropriate verification key for the given method.
func resolveVerifyKey(method jwt.SigningMethod, secret string, keyPEM []byte) (any, error) {
	switch method.(type) {
	case *jwt.SigningMethodHMAC:
		if secret == "" {
			return nil, fmt.Errorf("--secret is required for %s", method.Alg())
		}
		return []byte(secret), nil
	case *jwt.SigningMethodRSA:
		return parsePublicKey(keyPEM, "RSA")
	case *jwt.SigningMethodECDSA:
		return parsePublicKey(keyPEM, "EC")
	case *jwt.SigningMethodEd25519:
		return parsePublicKey(keyPEM, "Ed25519")
	default:
		return nil, fmt.Errorf("unsupported signing method: %s", method.Alg())
	}
}

func parsePrivateKey(keyPEM []byte, expect string) (any, error) {
	if len(keyPEM) == 0 {
		return nil, fmt.Errorf("--key is required for %s algorithms", expect)
	}

	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Try PKCS8 first (works for all key types).
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// Fall back to type-specific parsers.
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("failed to parse private key")
}

func parsePublicKey(keyPEM []byte, expect string) (any, error) {
	if len(keyPEM) == 0 {
		return nil, fmt.Errorf("--key is required for %s algorithms", expect)
	}

	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Try PKIX (standard public key format).
	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		return key, nil
	}

	// Try PKCS1 RSA public key.
	if key, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		return key, nil
	}

	// Try parsing as certificate and extracting public key.
	if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
		return cert.PublicKey, nil
	}

	return nil, fmt.Errorf("failed to parse public key")
}

func extractExpiry(claims jwt.MapClaims, info *JWTInfo) {
	exp, ok := claims["exp"]
	if !ok {
		return
	}
	if secs, ok := ClaimSeconds(exp); ok {
		expTime := time.Unix(secs, 0)
		info.ExpiredAt = &expTime
	}
}

// ClaimSeconds reads a numeric time claim, which is a json.Number when the
// parser preserves literals and a float64 when it does not.
func ClaimSeconds(v any) (int64, bool) {
	switch n := v.(type) {
	case json.Number:
		if i, err := n.Int64(); err == nil {
			return i, true
		}
		if f, err := n.Float64(); err == nil {
			return int64(f), true
		}
	case float64:
		return int64(n), true
	case int64:
		return n, true
	}
	return 0, false
}
