package inspect

import (
	"fmt"

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

		output, err := inspectLib.CertInfoJSON(info)
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
	certCmd.Flags().Bool("check-expiry", false, "exit with code 1 if certificate is expired")
	Cmd.AddCommand(certCmd)
}
