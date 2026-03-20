package encode

import (
	"fmt"

	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

var jwtCmd = &cobra.Command{
	Use:   "jwt [input]",
	Short: "Encode or decode JWT tokens",
	Long: `Encode: pass a JSON payload to create a signed JWT (requires --secret).
Decode: pass a JWT token with -d to inspect header and payload.`,
	Example: `  # Create a JWT
  swk encode jwt --secret mykey '{"sub":"user1","role":"admin"}'

  # Decode a JWT
  swk encode jwt -d 'eyJhbGciOiJIUzI1NiIs...'

  # Decode and verify signature
  swk encode jwt -d --secret mykey 'eyJhbGciOiJIUzI1NiIs...'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		decode, _ := cmd.Flags().GetBool("decode")
		secret, _ := cmd.Flags().GetString("secret")
		algo, _ := cmd.Flags().GetString("algo")

		if decode {
			return jwtDecodeRun(cmd, input, secret)
		}
		return jwtEncodeRun(cmd, input, secret, algo)
	},
}

func jwtEncodeRun(cmd *cobra.Command, payload, secret, algo string) error {
	if secret == "" {
		return fmt.Errorf("--secret is required to create a JWT")
	}
	token, err := encLib.JWTEncode(payload, secret, algo)
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), token)
	return nil
}

func jwtDecodeRun(cmd *cobra.Command, tokenStr, secret string) error {
	var info *encLib.JWTInfo
	var err error

	if secret != "" {
		info, err = encLib.JWTVerify(tokenStr, secret)
		if err != nil {
			if info != nil {
				output, jsonErr := encLib.JWTInfoJSON(info)
				if jsonErr == nil {
					fmt.Fprintln(cmd.OutOrStdout(), string(output))
				}
			}
			return err
		}
	} else {
		info, err = encLib.JWTDecode(tokenStr)
		if err != nil {
			return err
		}
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
	jwtCmd.Flags().StringP("algo", "a", "HS256", "signing algorithm (HS256, HS384, HS512)")
	Cmd.AddCommand(jwtCmd)
}
