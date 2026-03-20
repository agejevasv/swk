package encode

import escapeCmd "github.com/agejevasv/swk/cmd/escape"

// Register escape commands as aliases under encode,
// so users who look for "encode url" or "encode html" find them.
func init() {
	Cmd.AddCommand(escapeCmd.NewURLCmd())
	Cmd.AddCommand(escapeCmd.NewHTMLCmd())
	Cmd.AddCommand(escapeCmd.NewJSONCmd())
	Cmd.AddCommand(escapeCmd.NewXMLCmd())
	Cmd.AddCommand(escapeCmd.NewShellCmd())
}
