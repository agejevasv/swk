package encode

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

var jwtCmd = &cobra.Command{
	Use:   "jwt [input]",
	Short: "Encode or decode JWT tokens",
	Long: `Encode: pass a JSON payload to create a signed JWT.
Decode: pass a JWT token with -d to inspect header and payload (works with any algorithm).
Verify: pass -d with --secret (HMAC) or --key (RSA/EC/Ed25519) to verify the signature.

Supported algorithms: HS256, HS384, HS512, RS256, RS384, RS512, ES256, ES384, ES512, EdDSA.`,
	Example: `  # Create a JWT with HMAC
  swk encode jwt --secret mykey '{"sub":"user1","role":"admin"}'

  # Create a JWT with RSA
  swk encode jwt --algo RS256 --key private.pem '{"sub":"user1"}'

  # Decode a JWT (no verification, works with any algorithm)
  swk encode jwt -d 'eyJhbGciOiJIUzI1NiIs...'

  # Verify with HMAC secret
  swk encode jwt -d --secret mykey 'eyJhbGciOiJIUzI1NiIs...'

  # Verify with public key
  swk encode jwt -d --key public.pem 'eyJhbGciOiJSUzI1NiIs...'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		decode := ioutil.MustGetBool(cmd, "decode")
		secret := ioutil.MustGetString(cmd, "secret")
		algo := ioutil.MustGetString(cmd, "algo")
		keyPath := ioutil.MustGetString(cmd, "key")

		var keyPEM []byte
		if keyPath != "" {
			keyPEM, err = os.ReadFile(keyPath)
			if err != nil {
				return fmt.Errorf("failed to read key file: %w", err)
			}
		}

		if decode {
			return jwtDecodeRun(cmd, input, secret, keyPEM)
		}
		return jwtEncodeRun(cmd, input, secret, keyPEM, algo)
	},
}

func jwtEncodeRun(cmd *cobra.Command, payload, secret string, keyPEM []byte, algo string) error {
	if secret == "" && len(keyPEM) == 0 {
		return fmt.Errorf("--secret or --key is required to create a JWT")
	}
	token, err := encLib.JWTEncode(payload, secret, keyPEM, algo)
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), token)
	return nil
}

func jwtDecodeRun(cmd *cobra.Command, tokenStr, secret string, keyPEM []byte) error {
	var info *encLib.JWTInfo
	var err error

	if secret != "" || len(keyPEM) > 0 {
		info, err = encLib.JWTVerify(tokenStr, secret, keyPEM)
	} else {
		info, err = encLib.JWTDecode(tokenStr)
	}
	if err != nil {
		return err
	}

	output, err := encLib.JWTInfoJSON(info)
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(output))
	return nil
}

func init() {
	jwtCmd.Flags().BoolP("decode", "d", false, "decode/inspect a JWT token")
	jwtCmd.Flags().StringP("secret", "s", "", "HMAC secret for signing or verification")
	jwtCmd.Flags().StringP("key", "k", "", "path to PEM key file (private for sign, public for verify)")
	jwtCmd.Flags().StringP("algo", "a", "HS256", "signing algorithm (HS256, RS256, ES256, EdDSA, etc.)")
	Cmd.AddCommand(jwtCmd)
}
