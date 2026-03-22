package inspect

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	inspectLib "github.com/agejevasv/swk/internal/inspect"
	"github.com/agejevasv/swk/internal/ioutil"
)

var subnetCmd = &cobra.Command{
	Use:     "subnet [cidr]",
	Short:   "Calculate subnet information from CIDR notation",
	Example: "  swk inspect subnet 192.168.0.0/24",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		info, err := inspectLib.ParseSubnet(input)
		if err != nil {
			return err
		}

		if ioutil.MustGetBool(cmd, "json") {
			out, err := inspectLib.SubnetInfoJSON(info)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(out))
			return nil
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "Network:\t%s\n", info.Network)
		fmt.Fprintf(w, "Netmask:\t%s\n", info.Netmask)
		fmt.Fprintf(w, "Broadcast:\t%s\n", info.Broadcast)
		fmt.Fprintf(w, "First:\t%s\n", info.First)
		fmt.Fprintf(w, "Last:\t%s\n", info.Last)
		fmt.Fprintf(w, "Hosts:\t%d\n", info.Hosts)
		w.Flush()

		return nil
	},
}

func init() {
	subnetCmd.Flags().Bool("json", false, "output as JSON")
	Cmd.AddCommand(subnetCmd)
}
