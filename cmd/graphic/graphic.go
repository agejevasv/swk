package graphic

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "graphic",
	Aliases: []string{"gfx"},
	Short:   "Graphic tools",
}
