package encode

import escapeCmd "github.com/agejevasv/swk/cmd/escape"

// Register url and html under encode, since users commonly
// think of percent-encoding and HTML entity encoding as "encoding".
func init() {
	Cmd.AddCommand(escapeCmd.NewURLCmd())
	Cmd.AddCommand(escapeCmd.NewHTMLCmd())
}
