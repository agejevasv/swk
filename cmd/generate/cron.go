package generate

import (
	"fmt"

	"github.com/spf13/cobra"

	genLib "github.com/agejevasv/swk/internal/gen"
	"github.com/agejevasv/swk/internal/ioutil"
)

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Generate cron expressions from flags",
	Long:  "Build a cron expression from human-readable scheduling flags.",
	Example: `  swk generate cron --every 5m
  swk generate cron --daily --at 9:00
  swk generate cron --weekdays --at 9:00
  swk generate cron --weekly --day MON --at 9:00
  swk generate cron --monthly --day 15
  swk generate cron --yearly --month JUN --day 1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := genLib.CronOptions{
			Every:    ioutil.MustGetString(cmd, "every"),
			Daily:    ioutil.MustGetBool(cmd, "daily"),
			Weekly:   ioutil.MustGetBool(cmd, "weekly"),
			Monthly:  ioutil.MustGetBool(cmd, "monthly"),
			Yearly:   ioutil.MustGetBool(cmd, "yearly"),
			Weekdays: ioutil.MustGetBool(cmd, "weekdays"),
			At:       ioutil.MustGetString(cmd, "at"),
			Day:      ioutil.MustGetString(cmd, "day"),
			Month:    ioutil.MustGetString(cmd, "month"),
		}

		expr, err := genLib.GenerateCron(opts)
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), expr)
		return nil
	},
}

func init() {
	cronCmd.Flags().String("every", "", "repeating interval (e.g. 5m, 2h)")
	cronCmd.Flags().Bool("daily", false, "run once per day")
	cronCmd.Flags().Bool("weekly", false, "run once per week")
	cronCmd.Flags().Bool("monthly", false, "run once per month")
	cronCmd.Flags().Bool("yearly", false, "run once per year")
	cronCmd.Flags().Bool("weekdays", false, "run on weekdays (Mon-Fri)")
	cronCmd.Flags().String("at", "", "time of day (HH:MM)")
	cronCmd.Flags().String("day", "", "day of week (MON-SUN) or day of month (1-31)")
	cronCmd.Flags().String("month", "", "month (1-12 or JAN-DEC)")
	Cmd.AddCommand(cronCmd)
}
