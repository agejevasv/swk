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
		dtFrom, _ := cmd.Flags().GetString("from")
		dtTo, _ := cmd.Flags().GetString("to")
		dtTz, _ := cmd.Flags().GetString("tz")

		var inputStr string
		var err error

		if len(args) > 0 && strings.EqualFold(args[0], "now") {
			now := time.Now()
			if dtTz != "" && !strings.EqualFold(dtTz, "Local") {
				loc, err := time.LoadLocation(dtTz)
				if err != nil {
					return fmt.Errorf("invalid timezone %q: %w", dtTz, err)
				}
				now = now.In(loc)
			}
			inputStr = fmt.Sprintf("%d", now.Unix())
			result, err := convLib.ConvertDateTime(inputStr, "unix", dtTo, dtTz)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), result)
			return nil
		}

		inputStr, err = ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
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
