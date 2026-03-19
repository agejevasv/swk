package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"

	convertCmd "github.com/agejevasv/swk/cmd/convert"
	encodeCmd "github.com/agejevasv/swk/cmd/encode"
	escapeCmd "github.com/agejevasv/swk/cmd/escape"
	generateCmd "github.com/agejevasv/swk/cmd/generate"
	inspectCmd "github.com/agejevasv/swk/cmd/inspect"
	queryCmd "github.com/agejevasv/swk/cmd/query"
)

var showVersion bool

// Version is set via ldflags: -X github.com/agejevasv/swk/cmd.Version=v1.0.0
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:           "swk",
	Short:         "Developer's Swiss Army Knife",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			commit, date, dirty := vcsInfo()
			dirtyMark := ""
			if dirty {
				dirtyMark = "-dirty"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "swk %s (commit: %s%s, built: %s)\n", Version, commit, dirtyMark, date)
			return
		}
		cmd.Help()
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "print version")

	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.AddCommand(convertCmd.Cmd)
	rootCmd.AddCommand(encodeCmd.Cmd)
	rootCmd.AddCommand(escapeCmd.Cmd)
	rootCmd.AddCommand(generateCmd.Cmd)
	rootCmd.AddCommand(inspectCmd.Cmd)
	rootCmd.AddCommand(queryCmd.Cmd)
}

func vcsInfo() (commit, date string, dirty bool) {
	commit = "unknown"
	date = "unknown"

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			if len(s.Value) > 7 {
				commit = s.Value[:7]
			} else {
				commit = s.Value
			}
		case "vcs.time":
			date = s.Value
		case "vcs.modified":
			dirty = s.Value == "true"
		}
	}

	return
}

func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	return err
}
