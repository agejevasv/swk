package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// walkCommands visits every command in the tree.
func walkCommands(c *cobra.Command, fn func(*cobra.Command)) {
	fn(c)
	for _, sub := range c.Commands() {
		walkCommands(sub, fn)
	}
}

// exampleLines returns the runnable lines of a command's Example text.
func exampleLines(c *cobra.Command) []string {
	var out []string
	for _, line := range strings.Split(c.Example, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, line)
	}
	return out
}

// Example text is printed verbatim by --help, so whatever it shows has to work
// when pasted into a shell.
func TestExamplesArePasteable(t *testing.T) {
	walkCommands(rootCmd, func(c *cobra.Command) {
		if c.Example == "" {
			return
		}

		for _, line := range exampleLines(c) {
			path := c.CommandPath()

			// Example is a raw string literal: `\\n` reaches the terminal as
			// two characters and printf then emits a literal backslash.
			if strings.Contains(line, `\\`) {
				t.Errorf("%s: example contains an over-escaped backslash, which is printed verbatim:\n  %s", path, line)
			}

			// bash's builtin echo does not interpret escapes, so an example
			// relying on that silently produces one garbled line.
			if strings.HasPrefix(line, "echo ") && strings.Contains(line, `\n`) {
				t.Errorf("%s: echo does not interpret \\n in bash; use printf:\n  %s", path, line)
			}

			if !strings.Contains(line, "swk") {
				t.Errorf("%s: example line does not invoke swk:\n  %s", path, line)
			}
		}
	})
}

// Every leaf command should say how it is used.
func TestLeafCommandsAreDocumented(t *testing.T) {
	walkCommands(rootCmd, func(c *cobra.Command) {
		if c.HasSubCommands() || c.Hidden || c.Name() == "help" || c.Name() == "completion" {
			return
		}
		if c.Short == "" {
			t.Errorf("%s: missing Short description", c.CommandPath())
		}
		if c.Args == nil {
			t.Errorf("%s: missing an Args validator, so extra arguments are ignored", c.CommandPath())
		}
	})
}
