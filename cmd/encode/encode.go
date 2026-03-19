package encode

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "encode",
	Aliases: []string{"enc", "e"},
	Short:   "Encoders and decoders",
}
