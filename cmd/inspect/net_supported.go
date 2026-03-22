//go:build linux || darwin

package inspect

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	inspectLib "github.com/agejevasv/swk/internal/inspect"
	"github.com/agejevasv/swk/internal/ioutil"
)

var netCmd = &cobra.Command{
	Use:   "net",
	Short: "List processes listening on network ports",
	Long:  "List processes listening on TCP/UDP ports.",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := inspectLib.NetFilterOptions{
			All:  ioutil.MustGetBool(cmd, "all"),
			TCP:  ioutil.MustGetBool(cmd, "tcp"),
			UDP:  ioutil.MustGetBool(cmd, "udp"),
			Port: ioutil.MustGetInt(cmd, "port"),
		}

		entries, err := inspectLib.ListSockets(opts)
		if err != nil {
			return err
		}

		if ioutil.MustGetBool(cmd, "json") {
			out, err := inspectLib.NetSocketsJSON(entries)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(out))
			return nil
		}

		hasContainers := false
		for _, e := range entries {
			if e.Container != "" {
				hasContainers = true
				break
			}
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		if opts.All {
			if hasContainers {
				fmt.Fprintf(w, "PROTO\tLOCAL ADDRESS\tREMOTE ADDRESS\tSTATE\tPID\tPROCESS\tUSER\tSERVICE\tCONTAINER\n")
			} else {
				fmt.Fprintf(w, "PROTO\tLOCAL ADDRESS\tREMOTE ADDRESS\tSTATE\tPID\tPROCESS\tUSER\tSERVICE\n")
			}
			for _, e := range entries {
				if hasContainers {
					fmt.Fprintf(w, "%s\t%s:%d\t%s:%d\t%s\t%d\t%s\t%s\t%s\t%s\n",
						e.Proto,
						e.LocalIP, e.LocalPort,
						e.RemoteIP, e.RemotePort,
						e.State,
						e.PID, e.Process, e.User,
						e.Service, e.Container,
					)
				} else {
					fmt.Fprintf(w, "%s\t%s:%d\t%s:%d\t%s\t%d\t%s\t%s\t%s\n",
						e.Proto,
						e.LocalIP, e.LocalPort,
						e.RemoteIP, e.RemotePort,
						e.State,
						e.PID, e.Process, e.User,
						e.Service,
					)
				}
			}
		} else {
			if hasContainers {
				fmt.Fprintf(w, "PROTO\tLOCAL ADDRESS\tSTATE\tPID\tPROCESS\tUSER\tSERVICE\tCONTAINER\n")
			} else {
				fmt.Fprintf(w, "PROTO\tLOCAL ADDRESS\tSTATE\tPID\tPROCESS\tUSER\tSERVICE\n")
			}
			for _, e := range entries {
				if hasContainers {
					fmt.Fprintf(w, "%s\t%s:%d\t%s\t%d\t%s\t%s\t%s\t%s\n",
						e.Proto,
						e.LocalIP, e.LocalPort,
						e.State,
						e.PID, e.Process, e.User,
						e.Service, e.Container,
					)
				} else {
					fmt.Fprintf(w, "%s\t%s:%d\t%s\t%d\t%s\t%s\t%s\n",
						e.Proto,
						e.LocalIP, e.LocalPort,
						e.State,
						e.PID, e.Process, e.User,
						e.Service,
					)
				}
			}
		}
		w.Flush()

		return nil
	},
}

func init() {
	netCmd.Flags().BoolP("all", "a", false, "show all socket states, not just LISTEN")
	netCmd.Flags().Bool("tcp", false, "show only TCP sockets")
	netCmd.Flags().Bool("udp", false, "show only UDP sockets")
	netCmd.Flags().IntP("port", "p", 0, "filter by local port number")
	netCmd.Flags().Bool("json", false, "output as JSON")
	Cmd.AddCommand(netCmd)
}
