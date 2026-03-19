package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	textLib "github.com/agejevasv/swk/internal/text"
)

var caseTo string

var caseCmd = &cobra.Command{
	Use:   "case [text]",
	Short: "Convert text between case conventions",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		result, err := textLib.ConvertCase(input, caseTo)
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), result)
		return nil
	},
}

func init() {
	caseCmd.Flags().StringVarP(&caseTo, "to", "t", "", "Target case (camel, pascal, snake, kebab, upper, lower, title, sentence, dot, path)")
	caseCmd.MarkFlagRequired("to")
	Cmd.AddCommand(caseCmd)
}
