package gen

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net"
	"testing"
	"time"
)

func parseCert(t *testing.T, result *CertResult) *x509.Certificate {
	t.Helper()
	block, _ := pem.Decode(result.CertPEM)
	if block == nil {
		t.Fatal("failed to decode cert PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}
	return cert
}

func parseKey(t *testing.T, result *CertResult) any {
	t.Helper()
	block, _ := pem.Decode(result.KeyPEM)
	if block == nil {
		t.Fatal("failed to decode key PEM")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("failed to parse private key: %v", err)
	}
	return key
}

func TestGenerateCert_Default(t *testing.T) {
	result, err := GenerateCert(CertOptions{})
	if err != nil {
		t.Fatal(err)
	}

	cert := parseCert(t, result)
	if cert.Subject.CommonName != "localhost" {
		t.Errorf("expected CN localhost, got %s", cert.Subject.CommonName)
	}
	if len(cert.DNSNames) == 0 || cert.DNSNames[0] != "localhost" {
		t.Errorf("expected DNS SAN localhost, got %v", cert.DNSNames)
	}
	if len(cert.IPAddresses) == 0 || !cert.IPAddresses[0].Equal(net.ParseIP("127.0.0.1")) {
		t.Errorf("expected IP SAN 127.0.0.1, got %v", cert.IPAddresses)
	}

	key := parseKey(t, result)
	if _, ok := key.(*ecdsa.PrivateKey); !ok {
		t.Errorf("expected ECDSA key, got %T", key)
	}
}

func TestGenerateCert_CustomCN(t *testing.T) {
	result, err := GenerateCert(CertOptions{CN: "myapp.local"})
	if err != nil {
		t.Fatal(err)
	}
	cert := parseCert(t, result)
	if cert.Subject.CommonName != "myapp.local" {
		t.Errorf("expected CN myapp.local, got %s", cert.Subject.CommonName)
	}
	if cert.DNSNames[0] != "myapp.local" {
		t.Errorf("expected DNS SAN myapp.local, got %v", cert.DNSNames)
	}
}

func TestGenerateCert_DNSSANs(t *testing.T) {
	result, err := GenerateCert(CertOptions{DNS: []string{"a.example.com", "b.example.com"}})
	if err != nil {
		t.Fatal(err)
	}
	cert := parseCert(t, result)
	if len(cert.DNSNames) != 2 {
		t.Fatalf("expected 2 DNS SANs, got %d", len(cert.DNSNames))
	}
	if cert.DNSNames[0] != "a.example.com" || cert.DNSNames[1] != "b.example.com" {
		t.Errorf("unexpected DNS SANs: %v", cert.DNSNames)
	}
}

func TestGenerateCert_IPSANs(t *testing.T) {
	result, err := GenerateCert(CertOptions{IPs: []string{"10.0.0.1", "::1"}})
	if err != nil {
		t.Fatal(err)
	}
	cert := parseCert(t, result)
	if len(cert.IPAddresses) != 2 {
		t.Fatalf("expected 2 IP SANs, got %d", len(cert.IPAddresses))
	}
}

func TestGenerateCert_Days(t *testing.T) {
	result, err := GenerateCert(CertOptions{Days: 30})
	if err != nil {
		t.Fatal(err)
	}
	cert := parseCert(t, result)
	expected := time.Now().AddDate(0, 0, 30)
	diff := cert.NotAfter.Sub(expected)
	if diff < -time.Minute || diff > time.Minute {
		t.Errorf("expected NotAfter ~%v, got %v", expected, cert.NotAfter)
	}
}

func TestGenerateCert_RSA(t *testing.T) {
	result, err := GenerateCert(CertOptions{KeyType: "rsa"})
	if err != nil {
		t.Fatal(err)
	}
	key := parseKey(t, result)
	if _, ok := key.(*rsa.PrivateKey); !ok {
		t.Errorf("expected RSA key, got %T", key)
	}
}

func TestGenerateCert_InvalidKeyType(t *testing.T) {
	_, err := GenerateCert(CertOptions{KeyType: "dsa"})
	if err == nil {
		t.Fatal("expected error for invalid key type")
	}
}

func TestGenerateCert_InvalidIP(t *testing.T) {
	_, err := GenerateCert(CertOptions{IPs: []string{"not-an-ip"}})
	if err == nil {
		t.Fatal("expected error for invalid IP")
	}
}

func TestGenerateCert_Wildcard(t *testing.T) {
	result, err := GenerateCert(CertOptions{CN: "*.example.com", DNS: []string{"*.example.com"}})
	if err != nil {
		t.Fatal(err)
	}
	cert := parseCert(t, result)
	if cert.DNSNames[0] != "*.example.com" {
		t.Errorf("expected wildcard SAN, got %v", cert.DNSNames)
	}
}

func TestGenerateCert_KeyUsage(t *testing.T) {
	result, err := GenerateCert(CertOptions{})
	if err != nil {
		t.Fatal(err)
	}
	cert := parseCert(t, result)
	if cert.KeyUsage&x509.KeyUsageDigitalSignature == 0 {
		t.Error("expected DigitalSignature key usage")
	}
	if len(cert.ExtKeyUsage) < 2 {
		t.Error("expected ServerAuth and ClientAuth extended key usage")
	}
}
