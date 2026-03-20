package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var durationCmd = &cobra.Command{
	Use:   "duration [input]",
	Short:   "Convert between seconds and human-readable durations",
	Long:    "Convert between seconds and human-readable duration formats (e.g., 86400 <-> 1d).",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		durationTo := ioutil.MustGetString(cmd, "to")
		result, err := convLib.DurationConvert(input, durationTo)
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	durationCmd.Flags().String("to", "", "target format: human, seconds, minutes, hours (default: auto)")
	Cmd.AddCommand(durationCmd)
}
