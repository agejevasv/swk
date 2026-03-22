package inspect

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var cronCmd = &cobra.Command{
	Use:   "cron [expression]",
	Short: "Explain cron expressions",
	Example: `  swk inspect cron '*/5 * * * *'
  swk inspect cron --next 3 '0 9 * * MON'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		cronNext := ioutil.MustGetInt(cmd, "next")
		cronExplain := ioutil.MustGetBool(cmd, "explain")

		showExplain := cronExplain || !cmd.Flags().Changed("next")
		showNext := cmd.Flags().Changed("next") || !cronExplain

		if showExplain {
			explanation, err := convLib.CronExplain(input)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), explanation)
			if showNext {
				fmt.Fprintln(cmd.OutOrStdout())
			}
		}

		if showNext {
			times, err := convLib.CronNext(input, cronNext, time.Now())
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Next %d runs:\n", cronNext)
			for _, t := range times {
				fmt.Fprintln(cmd.OutOrStdout(), " ", t.Format(time.RFC3339))
			}
		}

		return nil
	},
}

func init() {
	cronCmd.Flags().IntP("next", "n", 5, "show next N scheduled runs")
	cronCmd.Flags().BoolP("explain", "e", false, "show human-readable explanation only")
	Cmd.AddCommand(cronCmd)
}
