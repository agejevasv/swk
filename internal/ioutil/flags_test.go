package ioutil

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestMustGetString(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "default", "test flag")
	got := MustGetString(cmd, "name")
	if got != "default" {
		t.Errorf("expected 'default', got %q", got)
	}
}

func TestMustGetString_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unregistered flag")
		}
	}()
	cmd := &cobra.Command{}
	MustGetString(cmd, "nonexistent")
}

func TestMustGetBool(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("verbose", true, "test flag")
	got := MustGetBool(cmd, "verbose")
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

func TestMustGetBool_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unregistered flag")
		}
	}()
	cmd := &cobra.Command{}
	MustGetBool(cmd, "nonexistent")
}

func TestMustGetInt(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Int("count", 42, "test flag")
	got := MustGetInt(cmd, "count")
	if got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
}

func TestMustGetInt_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unregistered flag")
		}
	}()
	cmd := &cobra.Command{}
	MustGetInt(cmd, "nonexistent")
}

func TestParseDelimiter(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    rune
		wantErr bool
	}{
		{"empty defaults to comma", "", ',', false},
		{"comma", ",", ',', false},
		{"semicolon", ";", ';', false},
		{"tab", "\t", '\t', false},
		{"pipe", "|", '|', false},
		{"multi-byte rune", "§", '§', false},
		{"multi-byte rune 2", "→", '→', false},
		{"two characters", ",,", 0, true},
		{"word", "sep", 0, true},
		{"invalid utf-8", "\xff", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDelimiter(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseDelimiter(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseDelimiter(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
