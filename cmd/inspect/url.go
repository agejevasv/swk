package inspect

import (
	"fmt"

	"github.com/spf13/cobra"

	inspLib "github.com/agejevasv/swk/internal/inspect"
	"github.com/agejevasv/swk/internal/ioutil"
)

var urlCmd = &cobra.Command{
	Use:   "url [input]",
	Short: "Parse and inspect a URL",
	Long:  "Parse a URL into its components (scheme, host, port, path, query, fragment).",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		info, err := inspLib.ParseURL(input)
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), inspLib.URLInfoTable(info))
		return nil
	},
}

func init() {
	Cmd.AddCommand(urlCmd)
}
