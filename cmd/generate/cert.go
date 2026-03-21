package generate

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	genLib "github.com/agejevasv/swk/internal/gen"
	"github.com/agejevasv/swk/internal/ioutil"
)

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Generate self-signed TLS certificates",
	Long:  "Generate a self-signed X.509 certificate and private key for local development.",
	Example: `  swk generate cert
  swk generate cert --cn myapp.local --days 30
  swk generate cert --dns localhost --dns myapp.local --ip 127.0.0.1
  swk generate cert -o ./certs/server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dns := cleanSlice(cmd, "dns")
		ips := cleanSlice(cmd, "ip")

		opts := genLib.CertOptions{
			CN:      ioutil.MustGetString(cmd, "cn"),
			DNS:     dns,
			IPs:     ips,
			Days:    ioutil.MustGetInt(cmd, "days"),
			KeyType: ioutil.MustGetString(cmd, "key-type"),
		}

		result, err := genLib.GenerateCert(opts)
		if err != nil {
			return err
		}

		out := ioutil.MustGetString(cmd, "out")
		certPath := out + ".pem"
		keyPath := out + "-key.pem"

		if err := os.WriteFile(certPath, result.CertPEM, 0o644); err != nil {
			return fmt.Errorf("failed to write %s: %w", certPath, err)
		}
		if err := os.WriteFile(keyPath, result.KeyPEM, 0o600); err != nil {
			return fmt.Errorf("failed to write %s: %w", keyPath, err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), certPath)
		fmt.Fprintln(cmd.OutOrStdout(), keyPath)
		return nil
	},
}

func cleanSlice(cmd *cobra.Command, name string) []string {
	vals, _ := cmd.Flags().GetStringSlice(name)
	var result []string
	for _, v := range vals {
		if v != "" && v != "[]" {
			result = append(result, v)
		}
	}
	return result
}

func init() {
	certCmd.Flags().String("cn", "localhost", "Common Name (subject)")
	certCmd.Flags().StringSlice("dns", nil, "DNS Subject Alternative Names")
	certCmd.Flags().StringSlice("ip", nil, "IP Subject Alternative Names")
	certCmd.Flags().Int("days", 365, "validity in days")
	certCmd.Flags().String("key-type", "ec", "key type: ec or rsa")
	certCmd.Flags().StringP("out", "o", "cert", "output path prefix (<prefix>.pem, <prefix>-key.pem)")
	Cmd.AddCommand(certCmd)
}
