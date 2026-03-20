package encode

import (
	"testing"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		algo    string
		want    string
		wantErr bool
	}{
		// Known test vectors: empty string.
		{
			name:  "md5_empty",
			input: []byte(""),
			algo:  "md5",
			want:  "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:  "sha1_empty",
			input: []byte(""),
			algo:  "sha1",
			want:  "da39a3ee5e6b4b0d3255bfef95601890afd80709",
		},
		{
			name:  "sha256_empty",
			input: []byte(""),
			algo:  "sha256",
			want:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},

		// Known test vectors: "hello".
		{
			name:  "md5_hello",
			input: []byte("hello"),
			algo:  "md5",
			want:  "5d41402abc4b2a76b9719d911017c592",
		},
		{
			name:  "sha1_hello",
			input: []byte("hello"),
			algo:  "sha1",
			want:  "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
		},
		{
			name:  "sha256_hello",
			input: []byte("hello"),
			algo:  "sha256",
			want:  "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:  "sha384_hello",
			input: []byte("hello"),
			algo:  "sha384",
			want:  "59e1748777448c69de6b800d7a33bbfb9ff1b463e44354c3553bcdb9c666fa90125a3c79f90397bdf5f6a13de828684f",
		},
		{
			name:  "sha512_hello",
			input: []byte("hello"),
			algo:  "sha512",
			want:  "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043",
		},

		// Known test vector: "abc" (commonly used in NIST examples).
		{
			name:  "sha256_abc",
			input: []byte("abc"),
			algo:  "sha256",
			want:  "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
		},

		// Case insensitive algo.
		{
			name:  "algo_uppercase",
			input: []byte("hello"),
			algo:  "SHA256",
			want:  "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:  "algo_mixed_case",
			input: []byte("hello"),
			algo:  "Sha256",
			want:  "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},

		// Error cases.
		{
			name:    "unsupported_algorithm",
			input:   []byte("hello"),
			algo:    "invalid",
			wantErr: true,
		},
		{
			name:    "empty_algo",
			input:   []byte("hello"),
			algo:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hash(tt.input, tt.algo)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Hash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Hash() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHashVerify(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		algo     string
		expected string
		want     bool
		wantErr  bool
	}{
		{
			name:     "correct_hash_returns_true",
			input:    []byte("hello"),
			algo:     "sha256",
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			want:     true,
		},
		{
			name:     "wrong_hash_returns_false",
			input:    []byte("hello"),
			algo:     "sha256",
			expected: "0000000000000000000000000000000000000000000000000000000000000000",
			want:     false,
		},
		{
			name:     "case_insensitive_comparison",
			input:    []byte("hello"),
			algo:     "sha256",
			expected: "2CF24DBA5FB0A30E26E83B2AC5B9E29E1B161E5C1FA7425E73043362938B9824",
			want:     true,
		},
		{
			name:     "md5_correct",
			input:    []byte("hello"),
			algo:     "md5",
			expected: "5d41402abc4b2a76b9719d911017c592",
			want:     true,
		},
		{
			name:     "empty_input_hash",
			input:    []byte(""),
			algo:     "md5",
			expected: "d41d8cd98f00b204e9800998ecf8427e",
			want:     true,
		},
		{
			name:     "unsupported_algo_returns_error",
			input:    []byte("hello"),
			algo:     "invalid",
			expected: "abc",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashVerify(tt.input, tt.algo, tt.expected)
			if (err != nil) != tt.wantErr {
				t.Fatalf("HashVerify() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("HashVerify() = %v, want %v", got, tt.want)
			}
		})
	}
}
