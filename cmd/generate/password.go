package generate

import (
	"fmt"

	"github.com/spf13/cobra"

	genLib "github.com/agejevasv/swk/internal/gen"
	"github.com/agejevasv/swk/internal/ioutil"
)

var passwordCmd = &cobra.Command{
	Use:     "password",
	Aliases: []string{"pw"},
	Short:   "Generate random passwords",
	Long:    "Generate cryptographically secure random passwords.",
	RunE: func(cmd *cobra.Command, args []string) error {
		pwLength := ioutil.MustGetInt(cmd, "length")
		pwCount := ioutil.MustGetInt(cmd, "count")
		pwNoUpper := ioutil.MustGetBool(cmd, "no-upper")
		pwNoLower := ioutil.MustGetBool(cmd, "no-lower")
		pwNoDigits := ioutil.MustGetBool(cmd, "no-digits")
		pwNoSymbols := ioutil.MustGetBool(cmd, "no-symbols")
		pwExclude := ioutil.MustGetString(cmd, "exclude")

		opts := genLib.PasswordOpts{
			Length:  pwLength,
			Upper:   !pwNoUpper,
			Lower:   !pwNoLower,
			Digits:  !pwNoDigits,
			Symbols: !pwNoSymbols,
			Exclude: pwExclude,
		}

		for i := 0; i < pwCount; i++ {
			pw, err := genLib.GeneratePassword(opts)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), pw)
		}
		return nil
	},
}

func init() {
	passwordCmd.Flags().IntP("length", "l", 16, "password length")
	passwordCmd.Flags().IntP("count", "n", 1, "number of passwords to generate")
	passwordCmd.Flags().Bool("no-upper", false, "exclude uppercase letters")
	passwordCmd.Flags().Bool("no-lower", false, "exclude lowercase letters")
	passwordCmd.Flags().Bool("no-digits", false, "exclude digits")
	passwordCmd.Flags().Bool("no-symbols", false, "exclude symbols")
	passwordCmd.Flags().String("exclude", "", "specific characters to exclude")
	Cmd.AddCommand(passwordCmd)
}
