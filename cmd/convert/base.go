package convert

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var baseNameMap = map[string]int{
	"bin": 2,
	"oct": 8,
	"dec": 10,
	"hex": 16,
}

var baseCmd = &cobra.Command{
	Use:   "base [input]",
	Short: "Convert numbers between bases",
	Long:  "Convert numbers between binary, octal, decimal, and hexadecimal.",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		nbFrom := ioutil.MustGetString(cmd, "from")
		nbTo := ioutil.MustGetString(cmd, "to")

		fromBase, ok := baseNameMap[strings.ToLower(nbFrom)]
		if !ok {
			return fmt.Errorf("unknown base name %q (use bin, oct, dec, hex)", nbFrom)
		}
		toBase, ok := baseNameMap[strings.ToLower(nbTo)]
		if !ok {
			return fmt.Errorf("unknown base name %q (use bin, oct, dec, hex)", nbTo)
		}

		result, err := convLib.ConvertBase(input, fromBase, toBase)
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	baseCmd.Flags().StringP("from", "f", "dec", "source base (bin, oct, dec, hex)")
	baseCmd.Flags().StringP("to", "t", "hex", "target base (bin, oct, dec, hex)")
	Cmd.AddCommand(baseCmd)
}
