package test

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"check", "t"},
	Short:   "Data testers and validators",
}
