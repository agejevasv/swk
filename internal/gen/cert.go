package gen

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

// CertOptions configures certificate generation.
type CertOptions struct {
	CN      string
	DNS     []string
	IPs     []string
	Days    int
	KeyType string // "ec" or "rsa"
}

// CertResult holds the generated PEM-encoded certificate and key.
type CertResult struct {
	CertPEM []byte
	KeyPEM  []byte
}

// GenerateCert creates a self-signed TLS certificate and private key.
func GenerateCert(opts CertOptions) (*CertResult, error) {
	if opts.CN == "" {
		opts.CN = "localhost"
	}
	if len(opts.DNS) == 0 {
		opts.DNS = []string{opts.CN}
	}
	if len(opts.IPs) == 0 {
		opts.IPs = []string{"127.0.0.1"}
	}
	if opts.Days <= 0 {
		opts.Days = 365
	}
	if opts.KeyType == "" {
		opts.KeyType = "ec"
	}

	// Parse IPs
	var ipAddrs []net.IP
	for _, s := range opts.IPs {
		ip := net.ParseIP(s)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP address %q", s)
		}
		ipAddrs = append(ipAddrs, ip)
	}

	// Generate key
	var privKey any
	var pubKey any
	switch opts.KeyType {
	case "ec":
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate EC key: %w", err)
		}
		privKey = key
		pubKey = &key.PublicKey
	case "rsa":
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, fmt.Errorf("failed to generate RSA key: %w", err)
		}
		privKey = key
		pubKey = &key.PublicKey
	default:
		return nil, fmt.Errorf("invalid --key-type %q: use ec or rsa", opts.KeyType)
	}

	// Serial number
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: opts.CN},
		NotBefore:    now,
		NotAfter:     now.AddDate(0, 0, opts.Days),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		DNSNames:     opts.DNS,
		IPAddresses:  ipAddrs,

		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Self-sign
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, pubKey, privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	// Encode private key as PKCS8
	keyDER, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})

	return &CertResult{CertPEM: certPEM, KeyPEM: keyPEM}, nil
}
