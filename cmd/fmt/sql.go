package fmt

import (
	goFmt "fmt"

	"github.com/spf13/cobra"

	fmtLib "github.com/agejevasv/swk/internal/fmt"
	"github.com/agejevasv/swk/internal/ioutil"
)

var sqlUppercase bool

var sqlCmd = &cobra.Command{
	Use:   "sql [input]",
	Short: "Format SQL",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		opts := fmtLib.SQLOptions{
			Uppercase: sqlUppercase,
		}

		result, err := fmtLib.FormatSQL([]byte(input), opts)
		if err != nil {
			return err
		}

		goFmt.Fprintln(cmd.OutOrStdout(), string(result))
		return nil
	},
}

func init() {
	sqlCmd.Flags().BoolVarP(&sqlUppercase, "uppercase", "u", false, "uppercase SQL keywords")

	Cmd.AddCommand(sqlCmd)
}
