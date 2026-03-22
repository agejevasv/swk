package inspect

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	inspectLib "github.com/agejevasv/swk/internal/inspect"
	"github.com/agejevasv/swk/internal/ioutil"
)

var certCmd = &cobra.Command{
	Use:   "cert [input]",
	Short: "Inspect an X.509 PEM certificate",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		certCheckExpiry := ioutil.MustGetBool(cmd, "check-expiry")

		info, err := inspectLib.CertDecode([]byte(input))
		if err != nil {
			return err
		}

		if ioutil.MustGetBool(cmd, "json") {
			output, err := inspectLib.CertInfoJSON(info)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(output))
		} else {
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "Subject:\t%s\n", info.Subject)
			fmt.Fprintf(w, "Issuer:\t%s\n", info.Issuer)
			fmt.Fprintf(w, "Not Before:\t%s\n", info.NotBefore.Format("2006-01-02 15:04:05 UTC"))
			fmt.Fprintf(w, "Not After:\t%s\n", info.NotAfter.Format("2006-01-02 15:04:05 UTC"))
			fmt.Fprintf(w, "Serial:\t%s\n", info.SerialNumber)
			fmt.Fprintf(w, "Algorithm:\t%s\n", info.SignatureAlgorithm)
			if len(info.DNSNames) > 0 {
				for _, name := range info.DNSNames {
					fmt.Fprintf(w, "DNS Name:\t%s\n", name)
				}
			}
			fmt.Fprintf(w, "Expired:\t%v\n", info.IsExpired)
			w.Flush()
		}

		if certCheckExpiry && info.IsExpired {
			return ioutil.CheckFailedError{}
		}

		return nil
	},
}

func init() {
	certCmd.Flags().Bool("check-expiry", false, "exit with code 1 if certificate is expired")
	certCmd.Flags().Bool("json", false, "output as JSON")
	Cmd.AddCommand(certCmd)
}
