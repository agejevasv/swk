package inspect

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"
)

const ipCheckURL = "https://checkip.amazonaws.com"

// LookupPublicIP returns the public IP address.
func LookupPublicIP() (string, error) {
	return lookupPublicIP(ipCheckURL)
}

func lookupPublicIP(url string) (string, error) {
	client := &http.Client{Timeout: 3 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to query public IP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to query public IP: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	return strings.TrimSpace(string(body)), nil
}

// LocalAddress is one address bound to a local network interface.
type LocalAddress struct {
	Interface string `json:"interface"`
	IP        string `json:"ip"`
	Prefix    int    `json:"prefix"`
	Family    string `json:"family"`
	Up        bool   `json:"up"`
	Loopback  bool   `json:"loopback"`
	LinkLocal bool   `json:"link_local"`
}

// CIDR renders the address in "address/prefix" form.
func (a LocalAddress) CIDR() string {
	return fmt.Sprintf("%s/%d", a.IP, a.Prefix)
}

// LocalAddressOptions controls which addresses LocalAddresses returns.
type LocalAddressOptions struct {
	// All includes loopback, interfaces that are down, and link-local
	// addresses, which are hidden by default as noise.
	All bool
}

// LocalAddresses returns the addresses bound to this machine's interfaces.
func LocalAddresses(opts LocalAddressOptions) ([]LocalAddress, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("listing network interfaces: %w", err)
	}

	var result []LocalAddress

	for _, iface := range ifaces {
		up := iface.Flags&net.FlagUp != 0
		loopback := iface.Flags&net.FlagLoopback != 0
		if !opts.All && (!up || loopback) {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			// An interface can disappear between the listing and the read.
			continue
		}

		var group []LocalAddress
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			linkLocal := ipNet.IP.IsLinkLocalUnicast() || ipNet.IP.IsLinkLocalMulticast()
			if !opts.All && linkLocal {
				continue
			}

			prefix, _ := ipNet.Mask.Size()
			family := "ipv6"
			if ipNet.IP.To4() != nil {
				family = "ipv4"
			}

			group = append(group, LocalAddress{
				Interface: iface.Name,
				IP:        ipNet.IP.String(),
				Prefix:    prefix,
				Family:    family,
				Up:        up,
				Loopback:  loopback,
				LinkLocal: linkLocal,
			})
		}

		// IPv4 first, then by address, so the output is stable between runs.
		sort.SliceStable(group, func(i, j int) bool {
			if group[i].Family != group[j].Family {
				return group[i].Family == "ipv4"
			}
			return group[i].IP < group[j].IP
		})

		result = append(result, group...)
	}

	return result, nil
}

// LocalAddressesJSON returns JSON-encoded output for a local address list.
func LocalAddressesJSON(addrs []LocalAddress) ([]byte, error) {
	if addrs == nil {
		addrs = []LocalAddress{}
	}
	return json.MarshalIndent(addrs, "", "  ")
}
