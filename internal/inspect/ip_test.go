package inspect

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLookupPublicIP_Valid(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("203.0.113.42\n"))
	}))
	defer srv.Close()

	ip, err := lookupPublicIP(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if ip != "203.0.113.42" {
		t.Errorf("expected 203.0.113.42, got %q", ip)
	}
}

func TestLookupPublicIP_IPv6(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("2001:db8::1\n"))
	}))
	defer srv.Close()

	ip, err := lookupPublicIP(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if ip != "2001:db8::1" {
		t.Errorf("expected 2001:db8::1, got %q", ip)
	}
}

func TestLookupPublicIP_Trimmed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("  10.0.0.1  \n"))
	}))
	defer srv.Close()

	ip, err := lookupPublicIP(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if ip != "10.0.0.1" {
		t.Errorf("expected trimmed IP, got %q", ip)
	}
}

func TestLookupPublicIP_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	_, err := lookupPublicIP(srv.URL)
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestLookupPublicIP_Unreachable(t *testing.T) {
	_, err := lookupPublicIP("http://192.0.2.1:1") // non-routable
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}

func TestLookupPublicIP_ValidIP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("93.184.216.34\n"))
	}))
	defer srv.Close()

	ip, err := lookupPublicIP(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if net.ParseIP(ip) == nil {
		t.Errorf("expected valid IP, got %q", ip)
	}
}
