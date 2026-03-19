package test

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	testLib "github.com/agejevasv/swk/internal/test"
)

var xmlvalCmd = &cobra.Command{
	Use:     "xmlval [xml]",
	Aliases: []string{"xv"},
	Short:   "Validate XML well-formedness",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		if err := testLib.ValidateXML([]byte(input)); err != nil {
			return fmt.Errorf("invalid: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Valid")
		return nil
	},
}

func init() {
	Cmd.AddCommand(xmlvalCmd)
}
