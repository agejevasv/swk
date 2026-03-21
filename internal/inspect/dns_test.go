package inspect

import (
	"net"
	"testing"
)

func skipIfNoDNS(t *testing.T) {
	t.Helper()
	if _, err := net.LookupHost("localhost"); err != nil {
		t.Skip("DNS not available")
	}
}

func TestLookupDNS_Localhost(t *testing.T) {
	skipIfNoDNS(t)
	result, err := LookupDNS("localhost", "")
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "localhost" {
		t.Errorf("expected name localhost, got %s", result.Name)
	}
	if len(result.Records) == 0 {
		t.Error("expected at least one record for localhost")
	}
}

func TestLookupDNS_TypeA(t *testing.T) {
	skipIfNoDNS(t)
	result, err := LookupDNS("localhost", "A")
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range result.Records {
		if r.Type != "A" {
			t.Errorf("expected only A records, got %s", r.Type)
		}
	}
}

func TestLookupDNS_ReverseLookup(t *testing.T) {
	skipIfNoDNS(t)
	result, err := LookupDNS("127.0.0.1", "")
	if err != nil {
		t.Skipf("reverse lookup not configured: %v", err)
	}
	if result.Name != "127.0.0.1" {
		t.Errorf("expected name 127.0.0.1, got %s", result.Name)
	}
	for _, r := range result.Records {
		if r.Type != "PTR" {
			t.Errorf("expected PTR records, got %s", r.Type)
		}
	}
}

func TestLookupDNS_InvalidType(t *testing.T) {
	_, err := LookupDNS("localhost", "INVALID")
	if err == nil {
		t.Fatal("expected error for invalid record type")
	}
}

func TestLookupDNS_PTRWithHostname(t *testing.T) {
	_, err := LookupDNS("example.com", "PTR")
	if err == nil {
		t.Fatal("expected error for PTR with hostname")
	}
}

func TestLookupDNS_IPDetection(t *testing.T) {
	// IPv4
	if ip := net.ParseIP("192.168.1.1"); ip == nil {
		t.Error("expected valid IPv4")
	}
	// IPv6
	if ip := net.ParseIP("::1"); ip == nil {
		t.Error("expected valid IPv6")
	}
	// Not an IP
	if ip := net.ParseIP("example.com"); ip != nil {
		t.Error("expected nil for hostname")
	}
}

func TestLookupDNS_EmptyInput(t *testing.T) {
	_, err := LookupDNS("", "")
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestDNSResultJSON(t *testing.T) {
	result := &DNSResult{
		Name: "example.com",
		Records: []DNSRecord{
			{Type: "A", Value: "93.184.216.34"},
		},
	}
	out, err := DNSResultJSON(result)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !containsStr(s, `"type": "A"`) || !containsStr(s, `"value": "93.184.216.34"`) {
		t.Errorf("unexpected JSON: %s", s)
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
