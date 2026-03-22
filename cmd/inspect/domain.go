package inspect

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	inspectLib "github.com/agejevasv/swk/internal/inspect"
	"github.com/agejevasv/swk/internal/ioutil"
)

var domainCmd = &cobra.Command{
	Use:   "domain [name]",
	Short: "Domain registration info via RDAP",
	Long:  "Query domain registration data (registrar, dates, nameservers) via RDAP and resolve DNS records.",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		info, err := inspectLib.LookupDomain(input)
		if err != nil {
			return err
		}

		if ioutil.MustGetBool(cmd, "json") {
			out, err := inspectLib.DomainInfoJSON(info)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(out))
			return nil
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "Domain:\t%s\n", info.Domain)
		if info.Registrar != "" {
			fmt.Fprintf(w, "Registrar:\t%s\n", info.Registrar)
		}
		if info.Created != "" {
			fmt.Fprintf(w, "Created:\t%s\n", info.Created)
		}
		if info.Expires != "" {
			fmt.Fprintf(w, "Expires:\t%s\n", info.Expires)
		}
		if info.Updated != "" {
			fmt.Fprintf(w, "Updated:\t%s\n", info.Updated)
		}
		if len(info.Status) > 0 {
			fmt.Fprintf(w, "Status:\t%s\n", strings.Join(info.Status, ", "))
		}
		if len(info.Nameservers) > 0 {
			fmt.Fprintf(w, "Nameservers:\t%s\n", strings.Join(info.Nameservers, ", "))
		}
		if len(info.A) > 0 {
			fmt.Fprintf(w, "A:\t%s\n", strings.Join(info.A, ", "))
		}
		if len(info.AAAA) > 0 {
			fmt.Fprintf(w, "AAAA:\t%s\n", strings.Join(info.AAAA, ", "))
		}
		if info.CNAME != "" {
			fmt.Fprintf(w, "CNAME:\t%s\n", info.CNAME)
		}
		if len(info.TXT) > 0 {
			for _, txt := range info.TXT {
				fmt.Fprintf(w, "TXT:\t%s\n", txt)
			}
		}
		w.Flush()

		return nil
	},
}

func init() {
	domainCmd.Flags().Bool("json", false, "output as JSON")
	Cmd.AddCommand(domainCmd)
}
