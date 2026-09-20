package ioutil

import (
	"fmt"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

// MustGetString returns a string flag value or panics if the flag is not registered.
func MustGetString(cmd *cobra.Command, name string) string {
	v, err := cmd.Flags().GetString(name)
	if err != nil {
		panic(fmt.Sprintf("bug: flag %q not registered: %v", name, err))
	}
	return v
}

// MustGetBool returns a bool flag value or panics if the flag is not registered.
func MustGetBool(cmd *cobra.Command, name string) bool {
	v, err := cmd.Flags().GetBool(name)
	if err != nil {
		panic(fmt.Sprintf("bug: flag %q not registered: %v", name, err))
	}
	return v
}

// MustGetInt returns an int flag value or panics if the flag is not registered.
func MustGetInt(cmd *cobra.Command, name string) int {
	v, err := cmd.Flags().GetInt(name)
	if err != nil {
		panic(fmt.Sprintf("bug: flag %q not registered: %v", name, err))
	}
	return v
}

// ParseDelimiter converts a delimiter flag value to a single rune.
// Taking s[0] instead would truncate any multi-byte delimiter.
func ParseDelimiter(s string) (rune, error) {
	if s == "" {
		return ',', nil
	}
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return 0, fmt.Errorf("invalid delimiter %q: not valid UTF-8", s)
	}
	if size != len(s) {
		return 0, fmt.Errorf("invalid delimiter %q: must be a single character", s)
	}
	return r, nil
}

// IntInRange returns an int flag value, rejecting anything outside [low, high].
// Numeric flags feed arithmetic, allocations and protocol fields, so they are
// bounded here rather than trusted downstream.
func IntInRange(cmd *cobra.Command, name string, low, high int) (int, error) {
	v := MustGetInt(cmd, name)
	if v < low || v > high {
		return 0, fmt.Errorf("--%s must be between %d and %d, got %d", name, low, high, v)
	}
	return v, nil
}

// IntAtLeast is IntInRange with no upper bound.
func IntAtLeast(cmd *cobra.Command, name string, low int) (int, error) {
	v := MustGetInt(cmd, name)
	if v < low {
		return 0, fmt.Errorf("--%s must be at least %d, got %d", name, low, v)
	}
	return v, nil
}
