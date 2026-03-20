package generate

import (
	"fmt"

	"github.com/spf13/cobra"

	genLib "github.com/agejevasv/swk/internal/gen"
)

var passwordCmd = &cobra.Command{
	Use:     "password",
	Aliases: []string{"pw"},
	Short:   "Generate random passwords",
	Long:    "Generate cryptographically secure random passwords.",
	RunE: func(cmd *cobra.Command, args []string) error {
		pwLength, _ := cmd.Flags().GetInt("length")
		pwCount, _ := cmd.Flags().GetInt("count")
		pwNoUpper, _ := cmd.Flags().GetBool("no-upper")
		pwNoLower, _ := cmd.Flags().GetBool("no-lower")
		pwNoDigits, _ := cmd.Flags().GetBool("no-digits")
		pwNoSymbols, _ := cmd.Flags().GetBool("no-symbols")
		pwExclude, _ := cmd.Flags().GetString("exclude")

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
