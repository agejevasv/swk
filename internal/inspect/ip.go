package inspect

import (
	"fmt"
	"io"
	"net/http"
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
