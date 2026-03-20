package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var chmodCmd = &cobra.Command{
	Use:   "chmod [input]",
	Short: "Convert between numeric and symbolic file permissions",
	Long:  "Convert between numeric (755) and symbolic (rwxr-xr-x) file permissions.",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		chmodTo := ioutil.MustGetString(cmd, "to")
		switch chmodTo {
		case "numeric":
			result, err := convLib.ChmodToNumeric(input)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), result)
		case "symbolic":
			result, err := convLib.ChmodToSymbolic(input)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), result)
		default:
			result, err := convLib.ChmodExplain(input)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), result)
		}

		return nil
	},
}

func init() {
	chmodCmd.Flags().String("to", "", "output format: numeric or symbolic")
	Cmd.AddCommand(chmodCmd)
}
