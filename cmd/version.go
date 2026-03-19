package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// Version is set via ldflags: -X github.com/agejevasv/swk/cmd.Version=v1.0.0
// Commit and date come automatically from git via debug.ReadBuildInfo().
var Version = "dev"

var shortVersion bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		commit, date, dirty := vcsInfo()

		if shortVersion {
			fmt.Fprintln(cmd.OutOrStdout(), Version)
			return
		}

		dirtyMark := ""
		if dirty {
			dirtyMark = "-dirty"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "swk %s (commit: %s%s, built: %s)\n", Version, commit, dirtyMark, date)
	},
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

func init() {
	versionCmd.Flags().BoolVar(&shortVersion, "short", false, "print version number only")
}
