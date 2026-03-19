package generate

import (
	"fmt"

	"github.com/spf13/cobra"

	genLib "github.com/agejevasv/swk/internal/gen"
)

var (
	pwLength    int
	pwCount     int
	pwNoUpper   bool
	pwNoLower   bool
	pwNoDigits  bool
	pwNoSymbols bool
	pwExclude   string
)

var passwordCmd = &cobra.Command{
	Use:     "password",
	Aliases: []string{"pw"},
	Short:   "Generate random passwords",
	Long:    "Generate cryptographically secure random passwords.",
	RunE: func(cmd *cobra.Command, args []string) error {
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
	passwordCmd.Flags().IntVarP(&pwLength, "length", "l", 16, "password length")
	passwordCmd.Flags().IntVarP(&pwCount, "count", "n", 1, "number of passwords to generate")
	passwordCmd.Flags().BoolVar(&pwNoUpper, "no-upper", false, "exclude uppercase letters")
	passwordCmd.Flags().BoolVar(&pwNoLower, "no-lower", false, "exclude lowercase letters")
	passwordCmd.Flags().BoolVar(&pwNoDigits, "no-digits", false, "exclude digits")
	passwordCmd.Flags().BoolVar(&pwNoSymbols, "no-symbols", false, "exclude symbols")
	passwordCmd.Flags().StringVar(&pwExclude, "exclude", "", "specific characters to exclude")
	Cmd.AddCommand(passwordCmd)
}
