package encode

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"strings"
)

func Hash(input []byte, algo string) (string, error) {
	h, err := newHash(algo)
	if err != nil {
		return "", err
	}
	h.Write(input)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func HashVerify(input []byte, algo, expected string) (bool, error) {
	computed, err := Hash(input, algo)
	if err != nil {
		return false, err
	}
	return strings.EqualFold(computed, expected), nil
}

func newHash(algo string) (hash.Hash, error) {
	switch strings.ToLower(algo) {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha384":
		return sha512.New384(), nil
	case "sha512":
		return sha512.New(), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %q (supported: md5, sha1, sha256, sha384, sha512)", algo)
	}
}
