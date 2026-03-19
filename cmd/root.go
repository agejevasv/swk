package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	convertCmd "github.com/agejevasv/swk/cmd/convert"
	encodeCmd "github.com/agejevasv/swk/cmd/encode"
	fmtCmd "github.com/agejevasv/swk/cmd/fmt"
	genCmd "github.com/agejevasv/swk/cmd/gen"
	graphicCmd "github.com/agejevasv/swk/cmd/graphic"
	testCmd "github.com/agejevasv/swk/cmd/test"
	textCmd "github.com/agejevasv/swk/cmd/text"
)

var rootCmd = &cobra.Command{
	Use:   "swk",
	Short: "Developer's Swiss Army Knife",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(convertCmd.Cmd)
	rootCmd.AddCommand(encodeCmd.Cmd)
	rootCmd.AddCommand(fmtCmd.Cmd)
	rootCmd.AddCommand(genCmd.Cmd)
	rootCmd.AddCommand(testCmd.Cmd)
	rootCmd.AddCommand(textCmd.Cmd)
	rootCmd.AddCommand(graphicCmd.Cmd)
	rootCmd.AddCommand(versionCmd)
}

func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	return err
}
