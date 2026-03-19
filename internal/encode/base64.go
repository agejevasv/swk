package encode

import (
	"encoding/base64"
	"strings"
)

func Base64Encode(input []byte, urlSafe, noPadding bool) string {
	enc := base64.StdEncoding
	if urlSafe {
		enc = base64.URLEncoding
	}
	if noPadding {
		enc = enc.WithPadding(base64.NoPadding)
	}
	return enc.EncodeToString(input)
}

func Base64Decode(input string, urlSafe bool) ([]byte, error) {
	input = strings.TrimSpace(input)
	enc := base64.StdEncoding
	if urlSafe {
		enc = base64.URLEncoding
	}
	// Try with padding first, then without
	result, err := enc.DecodeString(input)
	if err != nil {
		enc = enc.WithPadding(base64.NoPadding)
		result, err = enc.DecodeString(input)
	}
	return result, err
}
