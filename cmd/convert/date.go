package convert

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var dateCmd = &cobra.Command{
	Use:     "date [input]",
	Aliases: []string{"dt"},
	Short:   "Convert between date/time formats",
	Long:    "Convert between unix timestamps, ISO 8601, RFC 2822, and human-readable formats.",
	RunE: func(cmd *cobra.Command, args []string) error {
		dtFrom := ioutil.MustGetString(cmd, "from")
		dtTo := ioutil.MustGetString(cmd, "to")
		dtTz := ioutil.MustGetString(cmd, "tz")

		var inputStr string
		var err error

		if len(args) > 0 && strings.EqualFold(args[0], "now") {
			inputStr = fmt.Sprintf("%d", time.Now().Unix())
			dtFrom = "unix"
		} else {
			inputStr, err = ioutil.ReadInputString(args, cmd.InOrStdin())
			if err != nil {
				return err
			}
		}

		result, err := convLib.ConvertDateTime(inputStr, dtFrom, dtTo, dtTz)
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	dateCmd.Flags().StringP("from", "f", "auto", "input format (unix, unixms, iso, rfc2822, human, auto)")
	dateCmd.Flags().StringP("to", "t", "iso", "output format (unix, unixms, iso, rfc2822, human)")
	dateCmd.Flags().String("tz", "Local", "target timezone (e.g. UTC, America/New_York)")
	Cmd.AddCommand(dateCmd)
}
