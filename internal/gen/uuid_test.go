package gen

import (
	"regexp"
	"testing"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func TestGenerateUUID(t *testing.T) {
	tests := []struct {
		name      string
		version   int
		namespace string
		uuidName  string
		wantErr   bool
		wantExact string
	}{
		{
			name:    "v1_valid_format",
			version: 1,
		},
		{
			name:    "v4_valid_format",
			version: 4,
		},
		{
			name:      "v5_dns_known_vector",
			version:   5,
			namespace: "dns",
			uuidName:  "example.com",
			wantExact: "cfbff0d1-9375-5685-968c-48ce8b15ae17",
		},
		{
			name:      "v5_url_namespace",
			version:   5,
			namespace: "url",
			uuidName:  "https://example.com",
		},
		{
			name:      "v5_oid_namespace",
			version:   5,
			namespace: "oid",
			uuidName:  "1.2.3",
		},
		{
			name:      "v5_x500_namespace",
			version:   5,
			namespace: "x500",
			uuidName:  "CN=test",
		},
		{
			name:      "v5_custom_uuid_namespace",
			version:   5,
			namespace: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			uuidName:  "test",
		},
		{
			name:    "v7_valid_format",
			version: 7,
		},
		{
			name:    "invalid_version_3",
			version: 3,
			wantErr: true,
		},
		{
			name:    "invalid_version_0",
			version: 0,
			wantErr: true,
		},
		{
			name:    "invalid_version_99",
			version: 99,
			wantErr: true,
		},
		{
			name:      "v5_invalid_namespace",
			version:   5,
			namespace: "invalid",
			uuidName:  "test",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateUUID(tt.version, tt.namespace, tt.uuidName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !uuidRegex.MatchString(got) {
				t.Errorf("GenerateUUID() = %q, not a valid UUID format", got)
			}
			if tt.wantExact != "" && got != tt.wantExact {
				t.Errorf("GenerateUUID() = %q, want %q", got, tt.wantExact)
			}
		})
	}
}

func TestGenerateUUID_V4VersionNibble(t *testing.T) {
	for i := 0; i < 10; i++ {
		got, err := GenerateUUID(4, "", "")
		if err != nil {
			t.Fatalf("GenerateUUID(4): %v", err)
		}
		// The version nibble is the first character of the 3rd group.
		// Format: xxxxxxxx-xxxx-Vxxx-xxxx-xxxxxxxxxxxx
		if got[14] != '4' {
			t.Errorf("v4 UUID version nibble = %c, want '4': %s", got[14], got)
		}
	}
}

func TestGenerateUUID_V1TwoCallsDiffer(t *testing.T) {
	a, err := GenerateUUID(1, "", "")
	if err != nil {
		t.Fatalf("GenerateUUID(1): %v", err)
	}
	b, err := GenerateUUID(1, "", "")
	if err != nil {
		t.Fatalf("GenerateUUID(1): %v", err)
	}
	if a == b {
		t.Errorf("two v1 UUIDs should differ: %s == %s", a, b)
	}
}

func TestGenerateUUID_V7TwoCallsDiffer(t *testing.T) {
	a, err := GenerateUUID(7, "", "")
	if err != nil {
		t.Fatalf("GenerateUUID(7): %v", err)
	}
	b, err := GenerateUUID(7, "", "")
	if err != nil {
		t.Fatalf("GenerateUUID(7): %v", err)
	}
	if a == b {
		t.Errorf("two v7 UUIDs should differ: %s == %s", a, b)
	}
}

func TestGenerateUUID_V5Deterministic(t *testing.T) {
	a, _ := GenerateUUID(5, "dns", "example.com")
	b, _ := GenerateUUID(5, "dns", "example.com")
	if a != b {
		t.Errorf("v5 UUID should be deterministic: %s != %s", a, b)
	}
}

func TestGenerateUUID_V5DifferentNames(t *testing.T) {
	a, _ := GenerateUUID(5, "dns", "example.com")
	b, _ := GenerateUUID(5, "dns", "other.com")
	if a == b {
		t.Errorf("v5 UUIDs with different names should differ: %s == %s", a, b)
	}
}

func TestGenerateUUID_V4AllUnique(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		got, err := GenerateUUID(4, "", "")
		if err != nil {
			t.Fatalf("GenerateUUID(4): %v", err)
		}
		if seen[got] {
			t.Fatalf("duplicate v4 UUID found: %s", got)
		}
		seen[got] = true
	}
}
