package inspect

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

// DNSRecord represents a single DNS record.
type DNSRecord struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// DNSResult holds all DNS records for a query.
type DNSResult struct {
	Name    string      `json:"name"`
	Records []DNSRecord `json:"records"`
}

var validTypes = map[string]bool{
	"A": true, "AAAA": true, "MX": true, "NS": true,
	"TXT": true, "CNAME": true, "PTR": true,
}

// LookupDNS performs DNS lookups for the given name and optional record type.
func LookupDNS(name string, recordType string) (*DNSResult, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("hostname or IP address required")
	}
	recordType = strings.ToUpper(strings.TrimSpace(recordType))

	if recordType != "" && !validTypes[recordType] {
		return nil, fmt.Errorf("unknown record type %q: use A, AAAA, MX, NS, TXT, CNAME, or PTR", recordType)
	}

	result := &DNSResult{Name: name}

	// Reverse lookup if input is an IP
	if ip := net.ParseIP(name); ip != nil {
		return reverseLookup(name, ip)
	}

	if recordType != "" {
		records, err := lookupType(name, recordType)
		if err != nil {
			return nil, err
		}
		result.Records = records
	} else {
		result.Records = lookupAll(name)
	}

	return result, nil
}

func reverseLookup(name string, ip net.IP) (*DNSResult, error) {
	names, err := net.LookupAddr(ip.String())
	if err != nil {
		return nil, fmt.Errorf("reverse lookup failed for %s: %w", name, err)
	}
	result := &DNSResult{Name: name}
	for _, n := range names {
		result.Records = append(result.Records, DNSRecord{Type: "PTR", Value: n})
	}
	return result, nil
}

func lookupAll(name string) []DNSRecord {
	var records []DNSRecord
	for _, typ := range []string{"A", "AAAA", "CNAME", "MX", "NS", "TXT"} {
		if recs, err := lookupType(name, typ); err == nil {
			records = append(records, recs...)
		}
	}
	return records
}

func lookupType(name string, typ string) ([]DNSRecord, error) {
	switch typ {
	case "A":
		return lookupA(name)
	case "AAAA":
		return lookupAAAA(name)
	case "MX":
		return lookupMX(name)
	case "NS":
		return lookupNS(name)
	case "TXT":
		return lookupTXT(name)
	case "CNAME":
		return lookupCNAME(name)
	case "PTR":
		return nil, fmt.Errorf("PTR requires an IP address, not a hostname")
	default:
		return nil, fmt.Errorf("unknown record type %q", typ)
	}
}

func lookupA(name string) ([]DNSRecord, error) {
	addrs, err := net.LookupHost(name)
	if err != nil {
		return nil, err
	}
	var records []DNSRecord
	for _, addr := range addrs {
		if ip := net.ParseIP(addr); ip != nil && ip.To4() != nil {
			records = append(records, DNSRecord{Type: "A", Value: addr})
		}
	}
	return records, nil
}

func lookupAAAA(name string) ([]DNSRecord, error) {
	addrs, err := net.LookupHost(name)
	if err != nil {
		return nil, err
	}
	var records []DNSRecord
	for _, addr := range addrs {
		if ip := net.ParseIP(addr); ip != nil && ip.To4() == nil {
			records = append(records, DNSRecord{Type: "AAAA", Value: addr})
		}
	}
	return records, nil
}

func lookupMX(name string) ([]DNSRecord, error) {
	mxs, err := net.LookupMX(name)
	if err != nil {
		return nil, err
	}
	var records []DNSRecord
	for _, mx := range mxs {
		records = append(records, DNSRecord{Type: "MX", Value: fmt.Sprintf("%d %s", mx.Pref, mx.Host)})
	}
	return records, nil
}

func lookupNS(name string) ([]DNSRecord, error) {
	nss, err := net.LookupNS(name)
	if err != nil {
		return nil, err
	}
	var records []DNSRecord
	for _, ns := range nss {
		records = append(records, DNSRecord{Type: "NS", Value: ns.Host})
	}
	return records, nil
}

func lookupTXT(name string) ([]DNSRecord, error) {
	txts, err := net.LookupTXT(name)
	if err != nil {
		return nil, err
	}
	var records []DNSRecord
	for _, txt := range txts {
		records = append(records, DNSRecord{Type: "TXT", Value: txt})
	}
	return records, nil
}

func lookupCNAME(name string) ([]DNSRecord, error) {
	cname, err := net.LookupCNAME(name)
	if err != nil {
		return nil, err
	}
	// LookupCNAME returns the name itself if there's no CNAME — skip that
	if cname == name || cname == name+"." {
		return nil, nil
	}
	return []DNSRecord{{Type: "CNAME", Value: cname}}, nil
}

// DNSResultJSON returns JSON-encoded output.
func DNSResultJSON(result *DNSResult) ([]byte, error) {
	return json.MarshalIndent(result, "", "  ")
}
