package inspect

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	inspectLib "github.com/agejevasv/swk/internal/inspect"
	"github.com/agejevasv/swk/internal/ioutil"
)

var dnsCmd = &cobra.Command{
	Use:   "dns [hostname|ip]",
	Short: "DNS lookups",
	Long:  "Resolve DNS records for a hostname, or perform reverse lookup for an IP address.",
	Example: `  swk inspect dns example.com
  swk inspect dns example.com --type MX
  swk inspect dns 8.8.8.8`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		recordType := ioutil.MustGetString(cmd, "type")

		result, err := inspectLib.LookupDNS(input, recordType)
		if err != nil {
			return err
		}

		if ioutil.MustGetBool(cmd, "json") {
			out, err := inspectLib.DNSResultJSON(result)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(out))
			return nil
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		for _, r := range result.Records {
			fmt.Fprintf(w, "%s:\t%s\n", r.Type, r.Value)
		}
		w.Flush()

		return nil
	},
}

func init() {
	dnsCmd.Flags().StringP("type", "t", "", "record type (A, AAAA, MX, NS, TXT, CNAME, PTR)")
	dnsCmd.Flags().Bool("json", false, "output as JSON")
	Cmd.AddCommand(dnsCmd)
}
