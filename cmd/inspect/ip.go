package inspect

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	inspectLib "github.com/agejevasv/swk/internal/inspect"
	"github.com/agejevasv/swk/internal/ioutil"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Show your public or local IP addresses",
	Long: `Show the public IP address of this machine, or with --local the addresses
bound to its network interfaces.

--local hides loopback, interfaces that are down, and link-local addresses;
--all includes them.`,
	Example: `  swk inspect ip
  swk inspect ip --local
  swk inspect ip --local --all
  swk inspect ip --local --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		local := ioutil.MustGetBool(cmd, "local")
		all := ioutil.MustGetBool(cmd, "all")
		asJSON := ioutil.MustGetBool(cmd, "json")

		if all && !local {
			return fmt.Errorf("--all requires --local")
		}

		if !local {
			return runPublicIP(cmd, asJSON)
		}
		return runLocalIP(cmd, all, asJSON)
	},
}

func runPublicIP(cmd *cobra.Command, asJSON bool) error {
	ip, err := inspectLib.LookupPublicIP()
	if err != nil {
		return err
	}

	if asJSON {
		out, err := json.MarshalIndent(map[string]string{"public": ip}, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(out))
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), ip)
	return nil
}

func runLocalIP(cmd *cobra.Command, all, asJSON bool) error {
	addrs, err := inspectLib.LocalAddresses(inspectLib.LocalAddressOptions{All: all})
	if err != nil {
		return err
	}

	if asJSON {
		out, err := inspectLib.LocalAddressesJSON(addrs)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(out))
	} else {
		writeLocalTable(cmd, addrs)
	}

	if len(addrs) == 0 {
		fmt.Fprintln(cmd.ErrOrStderr(), "no local addresses found (try --all)")
		return ioutil.NoMatchError{}
	}

	return nil
}

// writeLocalTable prints one row per interface with its addresses joined.
func writeLocalTable(cmd *cobra.Command, addrs []inspectLib.LocalAddress) {
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)

	var current string
	var cidrs []string
	var up bool

	flush := func() {
		if current == "" {
			return
		}
		name := current
		if !up {
			name += " (down)"
		}
		fmt.Fprintf(w, "%s\t%s\n", name, strings.Join(cidrs, "  "))
	}

	for _, a := range addrs {
		if a.Interface != current {
			flush()
			current, cidrs, up = a.Interface, nil, a.Up
		}
		cidrs = append(cidrs, a.CIDR())
	}
	flush()

	w.Flush()
}

func init() {
	ipCmd.Flags().BoolP("local", "l", false, "show local interface addresses instead of the public IP")
	ipCmd.Flags().BoolP("all", "a", false, "with --local, include loopback, down interfaces and link-local addresses")
	ipCmd.Flags().Bool("json", false, "output as JSON")
	Cmd.AddCommand(ipCmd)
}
