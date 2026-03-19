package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	graphicLib "github.com/agejevasv/swk/internal/graphic"
	"github.com/agejevasv/swk/internal/ioutil"
)

var (
	colorFrom string
	colorTo   string
)

var colorCmd = &cobra.Command{
	Use:     "color [input]",
	Aliases: []string{"col"},
	Short:   "Convert between color formats",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		result, err := graphicLib.ConvertColor(input, colorFrom, colorTo)
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	colorCmd.Flags().StringVarP(&colorFrom, "from", "f", "auto", "input format: hex, rgb, hsl, hsv, cmyk, auto")
	colorCmd.Flags().StringVarP(&colorTo, "to", "t", "all", "output format: hex, rgb, hsl, hsv, cmyk, all")
	Cmd.AddCommand(colorCmd)
}
