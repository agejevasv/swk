package inspect

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// URLInfo holds parsed URL components.
type URLInfo struct {
	Scheme   string              `json:"scheme"`
	Host     string              `json:"host"`
	Port     string              `json:"port,omitempty"`
	Path     string              `json:"path,omitempty"`
	Query    map[string][]string `json:"query,omitempty"`
	Fragment string              `json:"fragment,omitempty"`
	User     string              `json:"user,omitempty"`
}

// ParseURL parses a URL string into its components.
func ParseURL(input string) (*URLInfo, error) {
	input = strings.TrimSpace(input)

	if !strings.Contains(input, "://") {
		return nil, fmt.Errorf("invalid URL %q: missing scheme (e.g. https://)", input)
	}

	u, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	info := &URLInfo{
		Scheme: u.Scheme,
		// EscapedPath keeps %2F distinct from a real separator, which is the
		// distinction someone inspecting a URL is usually looking for.
		Path:     u.EscapedPath(),
		Fragment: u.EscapedFragment(),
	}

	// Hosts are case-insensitive; report them the way the scheme is reported.
	info.Host = strings.ToLower(u.Hostname())
	info.Port = u.Port()

	if u.RawQuery != "" {
		info.Query = make(map[string][]string)
		for key, values := range u.Query() {
			info.Query[key] = values
		}
		// Query() drops pairs it cannot parse; keep the raw form rather than
		// reporting no query at all.
		if len(info.Query) == 0 {
			info.Query = map[string][]string{"": {u.RawQuery}}
		}
	}

	// User info
	if u.User != nil {
		info.User = u.User.Username()
	}

	return info, nil
}

// URLInfoJSON returns the URL info as formatted JSON.
func URLInfoJSON(info *URLInfo) ([]byte, error) {
	return json.MarshalIndent(info, "", "  ")
}

// URLInfoTable returns a formatted table-style string.
func URLInfoTable(info *URLInfo) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Scheme:    %s\n", info.Scheme)
	fmt.Fprintf(&sb, "Host:      %s\n", info.Host)
	if info.Port != "" {
		fmt.Fprintf(&sb, "Port:      %s\n", info.Port)
	}
	if info.Path != "" {
		fmt.Fprintf(&sb, "Path:      %s\n", info.Path)
	}
	if len(info.Query) > 0 {
		keys := make([]string, 0, len(info.Query))
		for k := range info.Query {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var parts []string
		for _, k := range keys {
			for _, v := range info.Query[k] {
				parts = append(parts, k+"="+v)
			}
		}
		fmt.Fprintf(&sb, "Query:     %s\n", strings.Join(parts, "&"))
	}
	if info.Fragment != "" {
		fmt.Fprintf(&sb, "Fragment:  %s\n", info.Fragment)
	}
	if info.User != "" {
		fmt.Fprintf(&sb, "User:      %s\n", info.User)
	}

	return strings.TrimRight(sb.String(), "\n")
}
