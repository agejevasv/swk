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
	Long: `Show the public IP address of this machine, its local interface addresses,
or both.

With no flags, only the public address is looked up. --local switches to the
local interfaces, hiding loopback, interfaces that are down and link-local
addresses. --all stops hiding those, and on its own also includes the public
address.

  swk inspect ip                 public address
  swk inspect ip --local         local addresses, filtered
  swk inspect ip --local --all   local addresses, nothing hidden (no network)
  swk inspect ip --all           public and local, nothing hidden`,
	Example: `  swk inspect ip
  swk inspect ip --local
  swk inspect ip --all
  swk inspect ip --local --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showPublic, showLocal, includeHidden := resolveIPMode(
			ioutil.MustGetBool(cmd, "local"),
			ioutil.MustGetBool(cmd, "all"),
		)
		return runIP(cmd, showPublic, showLocal, includeHidden, ioutil.MustGetBool(cmd, "json"))
	},
}

// lookupPublicIP is a seam so the offline path can be exercised in tests.
var lookupPublicIP = inspectLib.LookupPublicIP

// resolveIPMode turns the two view flags into what should be shown.
// --local pins the output to the local view; --all on its own covers both.
func resolveIPMode(local, all bool) (showPublic, showLocal, includeHidden bool) {
	return !local, local || all, all
}

func runIP(cmd *cobra.Command, showPublic, showLocal, includeHidden, asJSON bool) error {
	var public string
	var publicErr error
	if showPublic {
		public, publicErr = lookupPublicIP()
		// Losing the network must not throw away the local addresses too.
		if publicErr != nil && !showLocal {
			return publicErr
		}
	}

	var addrs []inspectLib.LocalAddress
	if showLocal {
		var err error
		addrs, err = inspectLib.LocalAddresses(inspectLib.LocalAddressOptions{All: includeHidden})
		if err != nil {
			return err
		}
	}

	if asJSON {
		if err := writeIPJSON(cmd, public, addrs, showPublic, showLocal); err != nil {
			return err
		}
	} else {
		writeIPTable(cmd, public, addrs)
	}

	if publicErr != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "public IP unavailable: %v\n", publicErr)
	}

	if public == "" && len(addrs) == 0 {
		return ioutil.NoMatchError{}
	}

	return nil
}

// writeIPJSON shapes the output to the view that was asked for: a bare list for
// the local view, an object once the public address is part of the answer.
func writeIPJSON(cmd *cobra.Command, public string, addrs []inspectLib.LocalAddress, showPublic, showLocal bool) error {
	var out []byte
	var err error

	switch {
	case showPublic && showLocal:
		if addrs == nil {
			addrs = []inspectLib.LocalAddress{}
		}
		out, err = json.MarshalIndent(struct {
			Public string                    `json:"public,omitempty"`
			Local  []inspectLib.LocalAddress `json:"local"`
		}{public, addrs}, "", "  ")
	case showLocal:
		out, err = inspectLib.LocalAddressesJSON(addrs)
	default:
		out, err = json.MarshalIndent(map[string]string{"public": public}, "", "  ")
	}
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), string(out))
	return nil
}

// writeIPTable prints the public address, when present, as the first row, then
// one row per interface with its addresses joined.
func writeIPTable(cmd *cobra.Command, public string, addrs []inspectLib.LocalAddress) {
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)

	// Alone, the public address is the whole answer and stays bare so it can
	// be piped straight into other commands.
	if public != "" && len(addrs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), public)
		return
	}
	if public != "" {
		fmt.Fprintf(w, "public\t%s\n", public)
	}

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
	ipCmd.Flags().BoolP("all", "a", false, "include loopback, down interfaces and link-local addresses; on its own, also show the public IP")
	ipCmd.Flags().Bool("json", false, "output as JSON")
	Cmd.AddCommand(ipCmd)
}
