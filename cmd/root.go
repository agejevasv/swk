package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"

	convertCmd "github.com/agejevasv/swk/cmd/convert"
	diffCmd "github.com/agejevasv/swk/cmd/diff"
	encodeCmd "github.com/agejevasv/swk/cmd/encode"
	escapeCmd "github.com/agejevasv/swk/cmd/escape"
	formatCmd "github.com/agejevasv/swk/cmd/format"
	generateCmd "github.com/agejevasv/swk/cmd/generate"
	inspectCmd "github.com/agejevasv/swk/cmd/inspect"
	listenCmd "github.com/agejevasv/swk/cmd/listen"
	queryCmd "github.com/agejevasv/swk/cmd/query"
	serveCmd "github.com/agejevasv/swk/cmd/serve"
)

// Version is set via ldflags: -X github.com/agejevasv/swk/cmd.Version=v1.0.0
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:           "swk",
	Short:         "Developer's Swiss Army Knife",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func versionString() string {
	commit, date, dirty := vcsInfo()
	dirtyMark := ""
	if dirty {
		dirtyMark = "-dirty"
	}
	return fmt.Sprintf("%s (commit: %s%s, built: %s)", Version, commit, dirtyMark, date)
}

func init() {
	rootCmd.SetVersionTemplate("swk version {{.Version}}\n")
	rootCmd.Version = versionString()
	rootCmd.Flags().BoolP("version", "V", false, "print version")

	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.AddCommand(convertCmd.Cmd)
	rootCmd.AddCommand(diffCmd.Cmd)
	rootCmd.AddCommand(encodeCmd.Cmd)
	rootCmd.AddCommand(escapeCmd.Cmd)
	rootCmd.AddCommand(formatCmd.Cmd)
	rootCmd.AddCommand(generateCmd.Cmd)
	rootCmd.AddCommand(inspectCmd.Cmd)
	rootCmd.AddCommand(listenCmd.Cmd)
	rootCmd.AddCommand(queryCmd.Cmd)
	rootCmd.AddCommand(serveCmd.Cmd)
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
	if err != nil && err.Error() != "" {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	return err
}
