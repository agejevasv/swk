package encode

import (
	"fmt"

	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

var certCheckExpiry bool

var certCmd = &cobra.Command{
	Use:   "cert [input]",
	Short: "Decode an X.509 PEM certificate",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		info, err := encLib.CertDecode([]byte(input))
		if err != nil {
			return err
		}

		output, err := encLib.CertInfoJSON(info)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(output))

		if certCheckExpiry && info.IsExpired {
			return fmt.Errorf("certificate is expired (expired at %s)", info.NotAfter.Format("2006-01-02T15:04:05Z07:00"))
		}

		return nil
	},
}

func init() {
	certCmd.Flags().BoolVar(&certCheckExpiry, "check-expiry", false, "exit with code 1 if certificate is expired")

	Cmd.AddCommand(certCmd)
}
