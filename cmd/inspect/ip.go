package inspect

import (
	"fmt"

	"github.com/spf13/cobra"

	inspectLib "github.com/agejevasv/swk/internal/inspect"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Show your public IP address",
	RunE: func(cmd *cobra.Command, args []string) error {
		ip, err := inspectLib.LookupPublicIP()
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), ip)
		return nil
	},
}

func init() {
	Cmd.AddCommand(ipCmd)
}
