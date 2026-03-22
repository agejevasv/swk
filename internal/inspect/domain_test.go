package inspect

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const fakeRDAPResponse = `{
  "ldhName": "EXAMPLE.COM",
  "status": ["client transfer prohibited", "server delete prohibited"],
  "entities": [
    {
      "roles": ["registrar"],
      "vcardArray": ["vcard", [
        ["version", {}, "text", "4.0"],
        ["fn", {}, "text", "Test Registrar Inc."]
      ]]
    }
  ],
  "events": [
    {"eventAction": "registration", "eventDate": "1997-09-15T04:00:00Z"},
    {"eventAction": "expiration", "eventDate": "2028-09-14T04:00:00Z"},
    {"eventAction": "last changed", "eventDate": "2019-09-09T15:39:04Z"}
  ],
  "nameservers": [
    {"ldhName": "NS1.EXAMPLE.COM"},
    {"ldhName": "NS2.EXAMPLE.COM"}
  ]
}`

func fakeRDAPServer(response string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rdap+json")
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	}))
}

func TestLookupDomain_Valid(t *testing.T) {
	srv := fakeRDAPServer(fakeRDAPResponse, 200)
	defer srv.Close()

	info, err := lookupDomain(srv.URL, "example.com")
	if err != nil {
		t.Fatal(err)
	}

	if info.Domain != "example.com" {
		t.Errorf("expected domain example.com, got %s", info.Domain)
	}
	if info.Registrar != "Test Registrar Inc." {
		t.Errorf("expected registrar 'Test Registrar Inc.', got %q", info.Registrar)
	}
	if info.Created != "1997-09-15" {
		t.Errorf("expected created 1997-09-15, got %s", info.Created)
	}
	if info.Expires != "2028-09-14" {
		t.Errorf("expected expires 2028-09-14, got %s", info.Expires)
	}
	if info.Updated != "2019-09-09" {
		t.Errorf("expected updated 2019-09-09, got %s", info.Updated)
	}
	if len(info.Status) != 2 {
		t.Errorf("expected 2 status entries, got %d", len(info.Status))
	}
	if len(info.Nameservers) != 2 {
		t.Errorf("expected 2 nameservers, got %d", len(info.Nameservers))
	}
	if info.Nameservers[0] != "ns1.example.com" {
		t.Errorf("expected lowercase nameserver, got %s", info.Nameservers[0])
	}
}

func TestLookupDomain_RegistrarExtraction(t *testing.T) {
	// Entity without registrar role should be skipped
	resp := `{
		"ldhName": "TEST.COM",
		"entities": [
			{"roles": ["technical"], "vcardArray": ["vcard", [["version",{},"text","4.0"],["fn",{},"text","Tech Contact"]]]},
			{"roles": ["registrar"], "vcardArray": ["vcard", [["version",{},"text","4.0"],["fn",{},"text","My Registrar"]]]}
		],
		"events": [],
		"nameservers": []
	}`

	srv := fakeRDAPServer(resp, 200)
	defer srv.Close()

	info, err := lookupDomain(srv.URL, "test.com")
	if err != nil {
		t.Fatal(err)
	}
	if info.Registrar != "My Registrar" {
		t.Errorf("expected 'My Registrar', got %q", info.Registrar)
	}
}

func TestLookupDomain_MissingFields(t *testing.T) {
	resp := `{"ldhName": "MINIMAL.COM"}`

	srv := fakeRDAPServer(resp, 200)
	defer srv.Close()

	info, err := lookupDomain(srv.URL, "minimal.com")
	if err != nil {
		t.Fatal(err)
	}
	if info.Domain != "minimal.com" {
		t.Errorf("expected minimal.com, got %s", info.Domain)
	}
	if info.Registrar != "" {
		t.Errorf("expected empty registrar, got %q", info.Registrar)
	}
	if info.Created != "" {
		t.Errorf("expected empty created, got %q", info.Created)
	}
}

func TestLookupDomain_Non200(t *testing.T) {
	srv := fakeRDAPServer("not found", 404)
	defer srv.Close()

	_, err := lookupDomain(srv.URL, "nonexistent.com")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestLookupDomain_InvalidJSON(t *testing.T) {
	srv := fakeRDAPServer("not json", 200)
	defer srv.Close()

	_, err := lookupDomain(srv.URL, "bad.com")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLookupDomain_EmptyName(t *testing.T) {
	_, err := lookupDomain("http://unused", "")
	if err == nil {
		t.Fatal("expected error for empty domain")
	}
}

func TestLookupDomain_NoDot(t *testing.T) {
	_, err := lookupDomain("http://unused", "localhost")
	if err == nil {
		t.Fatal("expected error for domain without dot")
	}
}

func TestLookupDomain_NameserversLowercase(t *testing.T) {
	resp := `{
		"ldhName": "TEST.COM",
		"nameservers": [
			{"ldhName": "NS1.CLOUDFLARE.COM"},
			{"ldhName": "NS2.CLOUDFLARE.COM"}
		]
	}`

	srv := fakeRDAPServer(resp, 200)
	defer srv.Close()

	info, err := lookupDomain(srv.URL, "test.com")
	if err != nil {
		t.Fatal(err)
	}
	if info.Nameservers[0] != "ns1.cloudflare.com" {
		t.Errorf("expected lowercase, got %s", info.Nameservers[0])
	}
}

func TestExtractVcardFN_EmptyArray(t *testing.T) {
	result := extractVcardFN(nil)
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestExtractVcardFN_NoFN(t *testing.T) {
	vcard := []any{"vcard", []any{
		[]any{"version", map[string]any{}, "text", "4.0"},
	}}
	result := extractVcardFN(vcard)
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestFormatEventDate_RFC3339(t *testing.T) {
	got := formatEventDate("2024-01-15T09:00:00Z")
	if got != "2024-01-15" {
		t.Errorf("expected 2024-01-15, got %s", got)
	}
}

func TestFormatEventDate_Truncate(t *testing.T) {
	got := formatEventDate("2024-01-15 some extra stuff")
	if got != "2024-01-15" {
		t.Errorf("expected 2024-01-15, got %s", got)
	}
}

func TestFormatEventDate_Short(t *testing.T) {
	got := formatEventDate("short")
	if got != "short" {
		t.Errorf("expected 'short', got %s", got)
	}
}

func TestDomainInfoJSON(t *testing.T) {
	info := &DomainInfo{
		Domain:    "example.com",
		Registrar: "Test Inc.",
	}
	out, err := DomainInfoJSON(info)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if parsed["domain"] != "example.com" {
		t.Errorf("expected example.com in JSON, got %v", parsed["domain"])
	}
}
