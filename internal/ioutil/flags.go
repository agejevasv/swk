package ioutil

import (
	"fmt"

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
