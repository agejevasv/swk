package inspect

import (
	"encoding/json"
	"testing"
)

func TestParseSubnet_24(t *testing.T) {
	info, err := ParseSubnet("192.168.1.0/24")
	if err != nil {
		t.Fatal(err)
	}
	if info.Network != "192.168.1.0/24" {
		t.Errorf("expected network 192.168.1.0/24, got %s", info.Network)
	}
	if info.Netmask != "255.255.255.0" {
		t.Errorf("expected netmask 255.255.255.0, got %s", info.Netmask)
	}
	if info.Broadcast != "192.168.1.255" {
		t.Errorf("expected broadcast 192.168.1.255, got %s", info.Broadcast)
	}
	if info.First != "192.168.1.1" {
		t.Errorf("expected first 192.168.1.1, got %s", info.First)
	}
	if info.Last != "192.168.1.254" {
		t.Errorf("expected last 192.168.1.254, got %s", info.Last)
	}
	if info.Hosts != 254 {
		t.Errorf("expected 254 hosts, got %d", info.Hosts)
	}
}

func TestParseSubnet_16(t *testing.T) {
	info, err := ParseSubnet("10.0.0.0/16")
	if err != nil {
		t.Fatal(err)
	}
	if info.Network != "10.0.0.0/16" {
		t.Errorf("expected 10.0.0.0/16, got %s", info.Network)
	}
	if info.Netmask != "255.255.0.0" {
		t.Errorf("expected 255.255.0.0, got %s", info.Netmask)
	}
	if info.Broadcast != "10.0.255.255" {
		t.Errorf("expected 10.0.255.255, got %s", info.Broadcast)
	}
	if info.Hosts != 65534 {
		t.Errorf("expected 65534 hosts, got %d", info.Hosts)
	}
}

func TestParseSubnet_8(t *testing.T) {
	info, err := ParseSubnet("10.0.0.0/8")
	if err != nil {
		t.Fatal(err)
	}
	if info.Hosts != 16777214 {
		t.Errorf("expected 16777214 hosts, got %d", info.Hosts)
	}
}

func TestParseSubnet_32(t *testing.T) {
	info, err := ParseSubnet("10.0.0.1/32")
	if err != nil {
		t.Fatal(err)
	}
	if info.Hosts != 1 {
		t.Errorf("expected 1 host, got %d", info.Hosts)
	}
	if info.First != "10.0.0.1" {
		t.Errorf("expected first 10.0.0.1, got %s", info.First)
	}
	if info.Last != "10.0.0.1" {
		t.Errorf("expected last 10.0.0.1, got %s", info.Last)
	}
}

func TestParseSubnet_31(t *testing.T) {
	info, err := ParseSubnet("10.0.0.0/31")
	if err != nil {
		t.Fatal(err)
	}
	if info.Hosts != 2 {
		t.Errorf("expected 2 hosts, got %d", info.Hosts)
	}
	if info.First != "10.0.0.0" {
		t.Errorf("expected first 10.0.0.0, got %s", info.First)
	}
	if info.Last != "10.0.0.1" {
		t.Errorf("expected last 10.0.0.1, got %s", info.Last)
	}
}

func TestParseSubnet_HostBits(t *testing.T) {
	// Input has host bits set — net.ParseCIDR normalizes to network address
	info, err := ParseSubnet("192.168.1.50/24")
	if err != nil {
		t.Fatal(err)
	}
	if info.Network != "192.168.1.0/24" {
		t.Errorf("expected normalized network 192.168.1.0/24, got %s", info.Network)
	}
}

func TestParseSubnet_Invalid(t *testing.T) {
	_, err := ParseSubnet("not-a-cidr")
	if err == nil {
		t.Fatal("expected error for invalid CIDR")
	}
}

func TestParseSubnet_IPv6(t *testing.T) {
	_, err := ParseSubnet("2001:db8::/32")
	if err == nil {
		t.Fatal("expected error for IPv6")
	}
}

func TestParseSubnet_Empty(t *testing.T) {
	_, err := ParseSubnet("")
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestSubnetInfoJSON(t *testing.T) {
	info := &SubnetInfo{
		Network: "192.168.1.0/24",
		Netmask: "255.255.255.0",
		Hosts:   254,
	}
	out, err := SubnetInfoJSON(info)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if parsed["network"] != "192.168.1.0/24" {
		t.Errorf("expected network in JSON, got %v", parsed["network"])
	}
}
