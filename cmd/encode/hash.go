package encode

import (
	"fmt"

	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

var hashCmd = &cobra.Command{
	Use:     "hash [input]",
	Aliases: []string{"sum"},
	Short:   "Generate hash/checksum",
	Long:    "Compute hash of input using md5, sha1, sha256, sha384, or sha512.",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInput(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		hashAlgo := ioutil.MustGetString(cmd, "algo")
		hashVerify := ioutil.MustGetString(cmd, "verify")

		if hashVerify != "" {
			ok, err := encLib.HashVerify(input, hashAlgo, hashVerify)
			if err != nil {
				return err
			}
			if ok {
				fmt.Fprintln(cmd.OutOrStdout(), "OK")
			} else {
				computed, _ := encLib.Hash(input, hashAlgo)
				return fmt.Errorf("hash mismatch\nexpected: %s\ngot:      %s", hashVerify, computed)
			}
			return nil
		}

		hash, err := encLib.Hash(input, hashAlgo)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), hash)
		return nil
	},
}

func init() {
	hashCmd.Flags().StringP("algo", "a", "sha256", "hash algorithm (md5, sha1, sha256, sha384, sha512)")
	hashCmd.Flags().StringP("verify", "V", "", "hash to verify against")
	Cmd.AddCommand(hashCmd)
}
