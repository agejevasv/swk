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
