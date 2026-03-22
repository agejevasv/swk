package inspect

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// DomainInfo holds parsed RDAP and DNS information for a domain.
type DomainInfo struct {
	Domain      string   `json:"domain"`
	Registrar   string   `json:"registrar,omitempty"`
	Created     string   `json:"created,omitempty"`
	Expires     string   `json:"expires,omitempty"`
	Updated     string   `json:"updated,omitempty"`
	Status      []string `json:"status,omitempty"`
	Nameservers []string `json:"nameservers,omitempty"`
	A           []string `json:"a,omitempty"`
	AAAA        []string `json:"aaaa,omitempty"`
	CNAME       string   `json:"cname,omitempty"`
	TXT         []string `json:"txt,omitempty"`
}

const rdapBaseURL = "https://rdap.org"

// LookupDomain queries RDAP and DNS for domain information.
func LookupDomain(name string) (*DomainInfo, error) {
	return lookupDomain(rdapBaseURL, name)
}

func lookupDomain(baseURL string, name string) (*DomainInfo, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return nil, fmt.Errorf("domain name required")
	}
	if !strings.Contains(name, ".") {
		return nil, fmt.Errorf("invalid domain %q", name)
	}

	info, err := queryRDAP(baseURL, name)
	if err != nil {
		return nil, err
	}

	// DNS enrichment (non-fatal)
	if result, err := LookupDNS(name, "A"); err == nil {
		for _, r := range result.Records {
			info.A = append(info.A, r.Value)
		}
	}
	if result, err := LookupDNS(name, "AAAA"); err == nil {
		for _, r := range result.Records {
			info.AAAA = append(info.AAAA, r.Value)
		}
	}
	if result, err := LookupDNS(name, "CNAME"); err == nil && len(result.Records) > 0 {
		info.CNAME = result.Records[0].Value
	}
	if result, err := LookupDNS(name, "TXT"); err == nil {
		for _, r := range result.Records {
			info.TXT = append(info.TXT, r.Value)
		}
	}

	return info, nil
}

// rdapResponse represents the relevant parts of an RDAP JSON response.
type rdapResponse struct {
	LdhName     string           `json:"ldhName"`
	Status      []string         `json:"status"`
	Entities    []rdapEntity     `json:"entities"`
	Events      []rdapEvent      `json:"events"`
	Nameservers []rdapNameserver `json:"nameservers"`
}

type rdapEntity struct {
	Roles      []string `json:"roles"`
	VcardArray []any    `json:"vcardArray"`
}

type rdapEvent struct {
	EventAction string `json:"eventAction"`
	EventDate   string `json:"eventDate"`
}

type rdapNameserver struct {
	LdhName string `json:"ldhName"`
}

func queryRDAP(baseURL string, name string) (*DomainInfo, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	url := baseURL + "/domain/" + name
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/rdap+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("RDAP query failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RDAP query failed: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading RDAP response: %w", err)
	}

	var rdap rdapResponse
	if err := json.Unmarshal(body, &rdap); err != nil {
		return nil, fmt.Errorf("parsing RDAP response: %w", err)
	}

	info := &DomainInfo{
		Domain: strings.ToLower(rdap.LdhName),
		Status: rdap.Status,
	}

	// Registrar
	info.Registrar = extractRegistrar(rdap.Entities)

	// Events
	for _, ev := range rdap.Events {
		date := formatEventDate(ev.EventDate)
		switch ev.EventAction {
		case "registration":
			info.Created = date
		case "expiration":
			info.Expires = date
		case "last changed":
			info.Updated = date
		}
	}

	// Nameservers
	for _, ns := range rdap.Nameservers {
		info.Nameservers = append(info.Nameservers, strings.ToLower(ns.LdhName))
	}

	return info, nil
}

func extractRegistrar(entities []rdapEntity) string {
	for _, e := range entities {
		for _, role := range e.Roles {
			if role == "registrar" {
				return extractVcardFN(e.VcardArray)
			}
		}
	}
	return ""
}

func extractVcardFN(vcardArray []any) string {
	if len(vcardArray) < 2 {
		return ""
	}
	entries, ok := vcardArray[1].([]any)
	if !ok {
		return ""
	}
	for _, entry := range entries {
		arr, ok := entry.([]any)
		if !ok || len(arr) < 4 {
			continue
		}
		if name, ok := arr[0].(string); ok && name == "fn" {
			if val, ok := arr[3].(string); ok {
				return val
			}
		}
	}
	return ""
}

func formatEventDate(s string) string {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.Format("2006-01-02")
	}
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}

// DomainInfoJSON returns JSON-encoded output.
func DomainInfoJSON(info *DomainInfo) ([]byte, error) {
	return json.MarshalIndent(info, "", "  ")
}
