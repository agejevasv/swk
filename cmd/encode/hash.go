package encode

import (
	"fmt"

	"github.com/spf13/cobra"

	genLib "github.com/agejevasv/swk/internal/gen"
	"github.com/agejevasv/swk/internal/ioutil"
)

var (
	hashAlgo   string
	hashVerify string
)

var hashCmd = &cobra.Command{
	Use:     "hash [input]",
	Aliases: []string{"h"},
	Short:   "Generate hash/checksum",
	Long:    "Compute hash of input using md5, sha1, sha256, sha384, or sha512.",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInput(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		if hashVerify != "" {
			ok, err := genLib.HashVerify(input, hashAlgo, hashVerify)
			if err != nil {
				return err
			}
			if ok {
				fmt.Fprintln(cmd.OutOrStdout(), "OK")
			} else {
				computed, _ := genLib.Hash(input, hashAlgo)
				return fmt.Errorf("hash mismatch\nexpected: %s\ngot:      %s", hashVerify, computed)
			}
			return nil
		}

		hash, err := genLib.Hash(input, hashAlgo)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), hash)
		return nil
	},
}

func init() {
	hashCmd.Flags().StringVarP(&hashAlgo, "algo", "a", "sha256", "hash algorithm (md5, sha1, sha256, sha384, sha512)")
	hashCmd.Flags().StringVarP(&hashVerify, "verify", "V", "", "hash to verify against")
	Cmd.AddCommand(hashCmd)
}
