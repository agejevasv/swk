package convert

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	"github.com/agejevasv/swk/internal/ioutil"
)

var (
	cronNext    int
	cronExplain bool
)

var cronCmd = &cobra.Command{
	Use:     "cron [expression]",
	Aliases: []string{"cr"},
	Short:   "Parse and explain cron expressions",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

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
	cronCmd.Flags().IntVarP(&cronNext, "next", "n", 5, "show next N scheduled runs")
	cronCmd.Flags().BoolVarP(&cronExplain, "explain", "e", false, "show human-readable explanation only")
	Cmd.AddCommand(cronCmd)
}
