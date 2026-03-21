//go:build !linux

package inspect

import (
	"fmt"

	"github.com/spf13/cobra"
)

var netCmd = &cobra.Command{
	Use:    "net",
	Short:  "List processes listening on network ports (Linux only)",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("inspect net is only available on Linux")
	},
}

func init() {
	Cmd.AddCommand(netCmd)
}
