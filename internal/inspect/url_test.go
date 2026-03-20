package inspect

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		checkFn func(t *testing.T, info *URLInfo)
	}{
		{
			name:  "full_url_all_components",
			input: "https://user@example.com:8443/api/v1?key=value&foo=bar#section",
			checkFn: func(t *testing.T, info *URLInfo) {
				if info.Scheme != "https" {
					t.Errorf("Scheme = %q, want https", info.Scheme)
				}
				if info.Host != "example.com" {
					t.Errorf("Host = %q, want example.com", info.Host)
				}
				if info.Port != "8443" {
					t.Errorf("Port = %q, want 8443", info.Port)
				}
				if info.Path != "/api/v1" {
					t.Errorf("Path = %q, want /api/v1", info.Path)
				}
				if info.Fragment != "section" {
					t.Errorf("Fragment = %q, want section", info.Fragment)
				}
				if info.User != "user" {
					t.Errorf("User = %q, want user", info.User)
				}
				if len(info.Query["key"]) != 1 || info.Query["key"][0] != "value" {
					t.Errorf("Query[key] = %v, want [value]", info.Query["key"])
				}
				if len(info.Query["foo"]) != 1 || info.Query["foo"][0] != "bar" {
					t.Errorf("Query[foo] = %v, want [bar]", info.Query["foo"])
				}
			},
		},
		{
			name:  "no_scheme_adds_https",
			input: "example.com/path",
			checkFn: func(t *testing.T, info *URLInfo) {
				if info.Scheme != "https" {
					t.Errorf("Scheme = %q, want https", info.Scheme)
				}
				if info.Host != "example.com" {
					t.Errorf("Host = %q, want example.com", info.Host)
				}
				if info.Path != "/path" {
					t.Errorf("Path = %q, want /path", info.Path)
				}
			},
		},
		{
			name:  "http_scheme_preserved",
			input: "http://example.com",
			checkFn: func(t *testing.T, info *URLInfo) {
				if info.Scheme != "http" {
					t.Errorf("Scheme = %q, want http", info.Scheme)
				}
			},
		},
		{
			name:  "url_with_port",
			input: "https://localhost:3000",
			checkFn: func(t *testing.T, info *URLInfo) {
				if info.Host != "localhost" {
					t.Errorf("Host = %q, want localhost", info.Host)
				}
				if info.Port != "3000" {
					t.Errorf("Port = %q, want 3000", info.Port)
				}
			},
		},
		{
			name:  "url_without_port",
			input: "https://example.com",
			checkFn: func(t *testing.T, info *URLInfo) {
				if info.Port != "" {
					t.Errorf("Port = %q, want empty", info.Port)
				}
			},
		},
		{
			name:  "url_with_query_params",
			input: "https://example.com/search?q=hello&lang=en",
			checkFn: func(t *testing.T, info *URLInfo) {
				if len(info.Query) != 2 {
					t.Fatalf("expected 2 query params, got %d", len(info.Query))
				}
				if len(info.Query["q"]) != 1 || info.Query["q"][0] != "hello" {
					t.Errorf("Query[q] = %v, want [hello]", info.Query["q"])
				}
				if len(info.Query["lang"]) != 1 || info.Query["lang"][0] != "en" {
					t.Errorf("Query[lang] = %v, want [en]", info.Query["lang"])
				}
			},
		},
		{
			name:  "url_with_fragment",
			input: "https://example.com/page#top",
			checkFn: func(t *testing.T, info *URLInfo) {
				if info.Fragment != "top" {
					t.Errorf("Fragment = %q, want top", info.Fragment)
				}
			},
		},
		{
			name:  "url_with_user_info",
			input: "https://admin@example.com",
			checkFn: func(t *testing.T, info *URLInfo) {
				if info.User != "admin" {
					t.Errorf("User = %q, want admin", info.User)
				}
			},
		},
		{
			name:  "minimal_url_just_host",
			input: "example.com",
			checkFn: func(t *testing.T, info *URLInfo) {
				if info.Scheme != "https" {
					t.Errorf("Scheme = %q, want https", info.Scheme)
				}
				if info.Host != "example.com" {
					t.Errorf("Host = %q, want example.com", info.Host)
				}
				if info.Port != "" {
					t.Errorf("Port = %q, want empty", info.Port)
				}
				if info.Fragment != "" {
					t.Errorf("Fragment = %q, want empty", info.Fragment)
				}
				if info.User != "" {
					t.Errorf("User = %q, want empty", info.User)
				}
				if len(info.Query) != 0 {
					t.Errorf("Query = %v, want empty", info.Query)
				}
			},
		},
		{
			name:  "url_with_encoded_characters",
			input: "https://example.com/path%20with%20spaces?name=hello%20world",
			checkFn: func(t *testing.T, info *URLInfo) {
				// net/url.Parse decodes percent-encoded path segments
				if info.Path != "/path with spaces" {
					t.Errorf("Path = %q, want '/path with spaces'", info.Path)
				}
				if len(info.Query["name"]) != 1 || info.Query["name"][0] != "hello world" {
					t.Errorf("Query[name] = %v, want [hello world]", info.Query["name"])
				}
			},
		},
		{
			name:  "url_with_multiple_query_values_same_key",
			input: "https://example.com?tag=a&tag=b",
			checkFn: func(t *testing.T, info *URLInfo) {
				if len(info.Query["tag"]) != 2 || info.Query["tag"][0] != "a" || info.Query["tag"][1] != "b" {
					t.Errorf("Query[tag] = %v, want [a b]", info.Query["tag"])
				}
			},
		},
		{
			name:  "ftp_scheme",
			input: "ftp://files.example.com/pub/readme.txt",
			checkFn: func(t *testing.T, info *URLInfo) {
				if info.Scheme != "ftp" {
					t.Errorf("Scheme = %q, want ftp", info.Scheme)
				}
				if info.Host != "files.example.com" {
					t.Errorf("Host = %q, want files.example.com", info.Host)
				}
				if info.Path != "/pub/readme.txt" {
					t.Errorf("Path = %q, want /pub/readme.txt", info.Path)
				}
			},
		},
		{
			name:  "input_with_whitespace",
			input: "  https://example.com  ",
			checkFn: func(t *testing.T, info *URLInfo) {
				if info.Host != "example.com" {
					t.Errorf("Host = %q, want example.com", info.Host)
				}
			},
		},
		{
			name:  "url_with_deep_path",
			input: "https://example.com/a/b/c/d/e",
			checkFn: func(t *testing.T, info *URLInfo) {
				if info.Path != "/a/b/c/d/e" {
					t.Errorf("Path = %q, want /a/b/c/d/e", info.Path)
				}
			},
		},
		{
			name:  "url_no_query_has_nil_query_map",
			input: "https://example.com/path",
			checkFn: func(t *testing.T, info *URLInfo) {
				if info.Query != nil {
					t.Errorf("Query = %v, want nil", info.Query)
				}
			},
		},
		{
			name:    "empty_input",
			input:   "",
			wantErr: false, // net/url.Parse accepts empty strings with a prepended scheme
			checkFn: func(t *testing.T, info *URLInfo) {
				// After prepending https://, it parses but host is empty
				if info.Scheme != "https" {
					t.Errorf("Scheme = %q, want https", info.Scheme)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ParseURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseURL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if tt.checkFn != nil && info != nil {
				tt.checkFn(t, info)
			}
		})
	}
}

func TestURLInfoJSON(t *testing.T) {
	tests := []struct {
		name     string
		info     *URLInfo
		contains []string
	}{
		{
			name: "full_info",
			info: &URLInfo{
				Scheme:   "https",
				Host:     "example.com",
				Port:     "443",
				Path:     "/api",
				Query:    map[string][]string{"key": {"value"}},
				Fragment: "section",
				User:     "admin",
			},
			contains: []string{
				`"scheme": "https"`,
				`"host": "example.com"`,
				`"port": "443"`,
				`"path": "/api"`,
				`"fragment": "section"`,
				`"user": "admin"`,
				`"key"`,
				`"value"`,
			},
		},
		{
			name: "omits_empty_optional_fields",
			info: &URLInfo{
				Scheme: "https",
				Host:   "example.com",
			},
			contains: []string{
				`"scheme": "https"`,
				`"host": "example.com"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := URLInfoJSON(tt.info)
			if err != nil {
				t.Fatalf("URLInfoJSON() error = %v", err)
			}
			if len(out) == 0 {
				t.Fatal("expected non-empty JSON output")
			}
			s := string(out)
			for _, want := range tt.contains {
				if !strings.Contains(s, want) {
					t.Errorf("JSON output missing %q, got:\n%s", want, s)
				}
			}
		})
	}
}

func TestURLInfoJSON_OmitsEmpty(t *testing.T) {
	info := &URLInfo{
		Scheme: "https",
		Host:   "example.com",
	}
	out, err := URLInfoJSON(info)
	if err != nil {
		t.Fatalf("URLInfoJSON() error = %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	for _, key := range []string{"port", "path", "query", "fragment", "user"} {
		if _, ok := m[key]; ok {
			t.Errorf("expected key %q to be omitted from JSON, but it was present", key)
		}
	}
}

func TestURLInfoTable(t *testing.T) {
	tests := []struct {
		name     string
		info     *URLInfo
		contains []string
		excludes []string
	}{
		{
			name: "full_info",
			info: &URLInfo{
				Scheme:   "https",
				Host:     "example.com",
				Port:     "8080",
				Path:     "/api/v1",
				Query:    map[string][]string{"key": {"value"}},
				Fragment: "top",
				User:     "admin",
			},
			contains: []string{
				"Scheme:",
				"https",
				"Host:",
				"example.com",
				"Port:",
				"8080",
				"Path:",
				"/api/v1",
				"Query:",
				"key=value",
				"Fragment:",
				"top",
				"User:",
				"admin",
			},
		},
		{
			name: "minimal_info_omits_optional_labels",
			info: &URLInfo{
				Scheme: "https",
				Host:   "example.com",
			},
			contains: []string{
				"Scheme:",
				"https",
				"Host:",
				"example.com",
			},
			excludes: []string{
				"Port:",
				"Path:",
				"Query:",
				"Fragment:",
				"User:",
			},
		},
		{
			name: "with_port_only",
			info: &URLInfo{
				Scheme: "http",
				Host:   "localhost",
				Port:   "3000",
			},
			contains: []string{
				"Scheme:",
				"http",
				"Host:",
				"localhost",
				"Port:",
				"3000",
			},
			excludes: []string{
				"Path:",
				"Query:",
				"Fragment:",
				"User:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := URLInfoTable(tt.info)
			if out == "" {
				t.Fatal("expected non-empty table output")
			}
			for _, want := range tt.contains {
				if !strings.Contains(out, want) {
					t.Errorf("table output missing %q, got:\n%s", want, out)
				}
			}
			for _, exclude := range tt.excludes {
				if strings.Contains(out, exclude) {
					t.Errorf("table output should not contain %q, got:\n%s", exclude, out)
				}
			}
		})
	}
}
