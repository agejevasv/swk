package ioutil

import (
	"strconv"
	"strings"
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

func TestIntInRange(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Int("n", 0, "")

	tests := []struct {
		name      string
		value     string
		low, high int
		wantErr   bool
	}{
		{"inside range", "5", 0, 10, false},
		{"at lower bound", "0", 0, 10, false},
		{"at upper bound", "10", 0, 10, false},
		{"below range", "-1", 0, 10, true},
		{"above range", "11", 0, 10, true},
		{"far below", "-999999", 0, 10, true},
		{"negative range allows negatives", "-5", -64, 64, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cmd.Flags().Set("n", tt.value); err != nil {
				t.Fatalf("Set: %v", err)
			}
			got, err := IntInRange(cmd, "n", tt.low, tt.high)
			if (err != nil) != tt.wantErr {
				t.Fatalf("IntInRange(%s) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
			if err == nil && got != mustAtoi(t, tt.value) {
				t.Errorf("IntInRange returned %d, want %s", got, tt.value)
			}
			if err != nil && !strings.Contains(err.Error(), "--n") {
				t.Errorf("error should name the flag, got %v", err)
			}
		})
	}
}

func TestIntAtLeast(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Int("n", 0, "")

	for _, tt := range []struct {
		value   string
		low     int
		wantErr bool
	}{
		{"0", 0, false},
		{"1000000", 0, false},
		{"-1", 0, true},
		{"-1", -5, false},
	} {
		if err := cmd.Flags().Set("n", tt.value); err != nil {
			t.Fatalf("Set: %v", err)
		}
		if _, err := IntAtLeast(cmd, "n", tt.low); (err != nil) != tt.wantErr {
			t.Errorf("IntAtLeast(%s, low=%d) error = %v, wantErr %v", tt.value, tt.low, err, tt.wantErr)
		}
	}
}

func mustAtoi(t *testing.T, s string) int {
	t.Helper()
	n, err := strconv.Atoi(s)
	if err != nil {
		t.Fatalf("Atoi(%q): %v", s, err)
	}
	return n
}
