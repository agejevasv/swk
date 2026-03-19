package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var bytesDecimal bool

var bytesCmd = &cobra.Command{
	Use:   "bytes [input]",
	Short: "Convert between byte sizes and human-readable formats",
	Long:  "Convert between raw byte counts and human-readable sizes (KB, MB, GB, etc.).",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		result, err := convLib.BytesConvert(input, bytesDecimal)
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	bytesCmd.Flags().BoolVarP(&bytesDecimal, "decimal", "d", false, "use decimal units (1000-based)")
	Cmd.AddCommand(bytesCmd)
}
