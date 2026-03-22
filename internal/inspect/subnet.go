package inspect

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

// SubnetInfo holds parsed CIDR subnet information.
type SubnetInfo struct {
	Network   string `json:"network"`
	Netmask   string `json:"netmask"`
	Broadcast string `json:"broadcast"`
	First     string `json:"first"`
	Last      string `json:"last"`
	Hosts     uint32 `json:"hosts"`
}

// ParseSubnet parses a CIDR string and returns subnet information.
func ParseSubnet(cidr string) (*SubnetInfo, error) {
	cidr = strings.TrimSpace(cidr)

	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR %q: %w", cidr, err)
	}

	// Only support IPv4 for now
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, fmt.Errorf("IPv6 subnets not supported")
	}

	mask := ipNet.Mask
	ones, _ := mask.Size()

	network := ipNet.IP.To4()
	networkInt := binary.BigEndian.Uint32(network)
	maskInt := binary.BigEndian.Uint32(mask)
	broadcastInt := networkInt | ^maskInt

	var firstInt, lastInt uint32
	var hosts uint32
	if ones == 32 {
		firstInt = networkInt
		lastInt = networkInt
		hosts = 1
	} else if ones == 31 {
		firstInt = networkInt
		lastInt = broadcastInt
		hosts = 2
	} else {
		firstInt = networkInt + 1
		lastInt = broadcastInt - 1
		hosts = (1 << (32 - ones)) - 2
	}

	return &SubnetInfo{
		Network:   fmt.Sprintf("%s/%d", network, ones),
		Netmask:   fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3]),
		Broadcast: uint32ToIP(broadcastInt),
		First:     uint32ToIP(firstInt),
		Last:      uint32ToIP(lastInt),
		Hosts:     hosts,
	}, nil
}

func uint32ToIP(n uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", n>>24, (n>>16)&0xff, (n>>8)&0xff, n&0xff)
}

// SubnetInfoJSON returns JSON-encoded output.
func SubnetInfoJSON(info *SubnetInfo) ([]byte, error) {
	return json.MarshalIndent(info, "", "  ")
}
