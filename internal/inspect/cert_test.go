package inspect

import (
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
)

func generateTestCert(t *testing.T, notBefore, notAfter time.Time, dnsNames []string) []byte {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "test.example.com",
			Organization: []string{"Test Org"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
		DNSNames:  dnsNames,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
}

func TestCertDecode(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
		checkFn func(t *testing.T, info *CertInfo)
	}{
		{
			name:  "valid_self_signed_cert",
			input: generateTestCert(t, time.Now().Add(-time.Hour), time.Now().Add(24*time.Hour), nil),
			checkFn: func(t *testing.T, info *CertInfo) {
				if info.IsExpired {
					t.Error("expected certificate to not be expired")
				}
				if !strings.Contains(info.Subject, "test.example.com") {
					t.Errorf("subject missing CN, got: %s", info.Subject)
				}
				if !strings.Contains(info.Issuer, "test.example.com") {
					t.Errorf("issuer missing CN (self-signed), got: %s", info.Issuer)
				}
				if info.SerialNumber == "" {
					t.Error("expected serial number to be set")
				}
				if info.SignatureAlgorithm == "" {
					t.Error("expected signature algorithm to be set")
				}
			},
		},
		{
			name:  "cert_with_sans",
			input: generateTestCert(t, time.Now().Add(-time.Hour), time.Now().Add(24*time.Hour), []string{"test.example.com", "*.example.com", "alt.example.org"}),
			checkFn: func(t *testing.T, info *CertInfo) {
				if len(info.DNSNames) != 3 {
					t.Fatalf("expected 3 DNS names, got %d: %v", len(info.DNSNames), info.DNSNames)
				}
				if info.DNSNames[0] != "test.example.com" {
					t.Errorf("DNS name 0 = %q, want test.example.com", info.DNSNames[0])
				}
				if info.DNSNames[1] != "*.example.com" {
					t.Errorf("DNS name 1 = %q, want *.example.com", info.DNSNames[1])
				}
				if info.DNSNames[2] != "alt.example.org" {
					t.Errorf("DNS name 2 = %q, want alt.example.org", info.DNSNames[2])
				}
			},
		},
		{
			name:  "expired_certificate",
			input: generateTestCert(t, time.Now().Add(-48*time.Hour), time.Now().Add(-24*time.Hour), nil),
			checkFn: func(t *testing.T, info *CertInfo) {
				if !info.IsExpired {
					t.Error("expected certificate to be expired")
				}
			},
		},
		{
			name:  "cert_with_single_san",
			input: generateTestCert(t, time.Now().Add(-time.Hour), time.Now().Add(time.Hour), []string{"single.example.com"}),
			checkFn: func(t *testing.T, info *CertInfo) {
				if len(info.DNSNames) != 1 {
					t.Fatalf("expected 1 DNS name, got %d", len(info.DNSNames))
				}
				if info.DNSNames[0] != "single.example.com" {
					t.Errorf("DNS name = %q, want single.example.com", info.DNSNames[0])
				}
			},
		},
		{
			name:  "cert_no_sans_has_nil_dns_names",
			input: generateTestCert(t, time.Now().Add(-time.Hour), time.Now().Add(time.Hour), nil),
			checkFn: func(t *testing.T, info *CertInfo) {
				if len(info.DNSNames) != 0 {
					t.Errorf("expected no DNS names, got %v", info.DNSNames)
				}
			},
		},
		{
			name:    "not_pem_input",
			input:   []byte("this is not a certificate"),
			wantErr: true,
		},
		{
			name: "invalid_certificate_data_in_pem",
			input: pem.EncodeToMemory(&pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte("invalid certificate bytes"),
			}),
			wantErr: true,
		},
		{
			name:    "empty_input",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := CertDecode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CertDecode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFn != nil && info != nil {
				tt.checkFn(t, info)
			}
		})
	}
}

func TestCertDecode_NotBefore_NotAfter(t *testing.T) {
	notBefore := time.Now().Add(-time.Hour).Truncate(time.Second)
	notAfter := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	pemData := generateTestCert(t, notBefore, notAfter, nil)

	info, err := CertDecode(pemData)
	if err != nil {
		t.Fatalf("CertDecode: %v", err)
	}

	// Time comparison with second precision.
	if info.NotBefore.Unix() != notBefore.Unix() {
		t.Errorf("NotBefore = %v, want %v", info.NotBefore, notBefore)
	}
	if info.NotAfter.Unix() != notAfter.Unix() {
		t.Errorf("NotAfter = %v, want %v", info.NotAfter, notAfter)
	}
}

func TestCertInfoJSON(t *testing.T) {
	info := &CertInfo{
		Subject:            "CN=test",
		Issuer:             "CN=test",
		SerialNumber:       "1",
		SignatureAlgorithm: "ECDSA-SHA256",
		DNSNames:           []string{"test.example.com"},
		IsExpired:          false,
	}
	out, err := CertInfoJSON(info)
	if err != nil {
		t.Fatalf("CertInfoJSON: %v", err)
	}
	if len(out) == 0 {
		t.Error("expected non-empty JSON output")
	}
	if !strings.Contains(string(out), "test.example.com") {
		t.Error("JSON output should contain DNS name")
	}
}
