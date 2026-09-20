package inspect

import (
	"encoding/json"
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

func TestLocalAddress_CIDR(t *testing.T) {
	tests := []struct {
		addr LocalAddress
		want string
	}{
		{LocalAddress{IP: "192.168.1.42", Prefix: 24}, "192.168.1.42/24"},
		{LocalAddress{IP: "10.0.0.1", Prefix: 32}, "10.0.0.1/32"},
		{LocalAddress{IP: "::1", Prefix: 128}, "::1/128"},
	}

	for _, tt := range tests {
		if got := tt.addr.CIDR(); got != tt.want {
			t.Errorf("CIDR() = %q, want %q", got, tt.want)
		}
	}
}

// The default view hides loopback, down interfaces and link-local noise.
func TestLocalAddresses_DefaultExcludesNoise(t *testing.T) {
	addrs, err := LocalAddresses(LocalAddressOptions{})
	if err != nil {
		t.Fatalf("LocalAddresses: %v", err)
	}

	for _, a := range addrs {
		if a.Loopback {
			t.Errorf("loopback address %s on %s should be hidden by default", a.IP, a.Interface)
		}
		if a.LinkLocal {
			t.Errorf("link-local address %s on %s should be hidden by default", a.IP, a.Interface)
		}
		if !a.Up {
			t.Errorf("address %s on down interface %s should be hidden by default", a.IP, a.Interface)
		}
	}
}

func TestLocalAddresses_AllIncludesLoopback(t *testing.T) {
	addrs, err := LocalAddresses(LocalAddressOptions{All: true})
	if err != nil {
		t.Fatalf("LocalAddresses: %v", err)
	}

	var foundLoopback bool
	for _, a := range addrs {
		if a.Loopback {
			foundLoopback = true
			break
		}
	}
	if !foundLoopback {
		t.Error("--all should include the loopback interface")
	}

	if len(addrs) < len(mustLocalAddresses(t, LocalAddressOptions{})) {
		t.Error("--all should return at least as many addresses as the default view")
	}
}

func TestLocalAddresses_FieldsAreConsistent(t *testing.T) {
	addrs, err := LocalAddresses(LocalAddressOptions{All: true})
	if err != nil {
		t.Fatalf("LocalAddresses: %v", err)
	}

	for _, a := range addrs {
		ip := net.ParseIP(a.IP)
		if ip == nil {
			t.Errorf("interface %s: %q is not a valid IP", a.Interface, a.IP)
			continue
		}

		wantFamily := "ipv6"
		maxPrefix := 128
		if ip.To4() != nil {
			wantFamily = "ipv4"
			maxPrefix = 32
		}
		if a.Family != wantFamily {
			t.Errorf("%s: family = %q, want %q", a.IP, a.Family, wantFamily)
		}
		if a.Prefix < 0 || a.Prefix > maxPrefix {
			t.Errorf("%s: prefix %d out of range for %s", a.IP, a.Prefix, wantFamily)
		}
		if a.Interface == "" {
			t.Errorf("%s: missing interface name", a.IP)
		}
	}
}

// IPv4 sorts before IPv6 within an interface, so runs are comparable.
func TestLocalAddresses_StableOrdering(t *testing.T) {
	first := mustLocalAddresses(t, LocalAddressOptions{All: true})
	second := mustLocalAddresses(t, LocalAddressOptions{All: true})

	if len(first) != len(second) {
		t.Fatalf("two calls returned %d and %d addresses", len(first), len(second))
	}
	for i := range first {
		if first[i] != second[i] {
			t.Errorf("entry %d differs between calls: %+v vs %+v", i, first[i], second[i])
		}
	}

	seen := map[string]bool{}
	for i, a := range first {
		if i > 0 && first[i-1].Interface == a.Interface {
			if first[i-1].Family == "ipv6" && a.Family == "ipv4" {
				t.Errorf("interface %s lists ipv6 before ipv4", a.Interface)
			}
			continue
		}
		if seen[a.Interface] {
			t.Errorf("interface %s appears in more than one group", a.Interface)
		}
		seen[a.Interface] = true
	}
}

func TestLocalAddressesJSON(t *testing.T) {
	out, err := LocalAddressesJSON(nil)
	if err != nil {
		t.Fatalf("LocalAddressesJSON: %v", err)
	}
	if string(out) != "[]" {
		t.Errorf("empty list should encode as [], got %s", out)
	}

	out, err = LocalAddressesJSON([]LocalAddress{{Interface: "eth0", IP: "10.0.0.1", Prefix: 8, Family: "ipv4", Up: true}})
	if err != nil {
		t.Fatalf("LocalAddressesJSON: %v", err)
	}

	var round []LocalAddress
	if err := json.Unmarshal(out, &round); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(round) != 1 || round[0].Interface != "eth0" || round[0].IP != "10.0.0.1" {
		t.Errorf("round trip lost data: %+v", round)
	}
}

func mustLocalAddresses(t *testing.T, opts LocalAddressOptions) []LocalAddress {
	t.Helper()
	addrs, err := LocalAddresses(opts)
	if err != nil {
		t.Fatalf("LocalAddresses: %v", err)
	}
	return addrs
}
