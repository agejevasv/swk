package encode

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTInfo holds decoded JWT information.
type JWTInfo struct {
	Header    map[string]interface{} `json:"header"`
	Payload   map[string]interface{} `json:"payload"`
	Signature string                 `json:"signature"`
	Valid     bool                   `json:"valid"`
	ExpiredAt *time.Time             `json:"expired_at,omitempty"`
}

func JWTEncode(payloadJSON string, secret string, algo string) (string, error) {
	var claims jwt.MapClaims
	if err := json.Unmarshal([]byte(payloadJSON), &claims); err != nil {
		return "", fmt.Errorf("invalid JSON payload: %w", err)
	}

	var method jwt.SigningMethod
	switch algo {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	case "HS512":
		method = jwt.SigningMethodHS512
	default:
		return "", fmt.Errorf("unsupported algorithm: %s (supported: HS256, HS384, HS512)", algo)
	}

	token := jwt.NewWithClaims(method, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenStr, nil
}

func JWTDecode(tokenStr string) (*JWTInfo, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
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
		Payload:   map[string]interface{}(claims),
		Signature: sig,
	}

	if exp, ok := claims["exp"]; ok {
		if expFloat, ok := exp.(float64); ok {
			expTime := time.Unix(int64(expFloat), 0)
			info.ExpiredAt = &expTime
		}
	}

	return info, nil
}

func JWTVerify(tokenStr string, secret string) (*JWTInfo, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	// Extract signature from raw token string
	sig := ""
	if parts := strings.SplitN(tokenStr, ".", 3); len(parts) == 3 {
		sig = parts[2]
	}

	info := &JWTInfo{
		Valid:     false,
		Signature: sig,
	}

	// If token is nil, the JWT was malformed — that's a hard error.
	if token == nil {
		return nil, fmt.Errorf("invalid JWT: %w", err)
	}

	info.Header = token.Header
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		info.Payload = map[string]interface{}(claims)

		if exp, ok := claims["exp"]; ok {
			if expFloat, ok := exp.(float64); ok {
				expTime := time.Unix(int64(expFloat), 0)
				info.ExpiredAt = &expTime
			}
		}
	}
	info.Valid = token.Valid

	// Signature mismatch or expiry is not an error — it's represented by Valid=false.
	return info, nil
}

func JWTInfoJSON(info *JWTInfo) ([]byte, error) {
	return json.MarshalIndent(info, "", "  ")
}
