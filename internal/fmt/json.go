package fmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type JSONOptions struct {
	// Indent is the number of spaces per level; 0 produces compact output.
	Indent int
	Minify bool
}

// FormatJSON re-indents the document without rewriting it. Decoding into Go
// values would sort object keys, drop duplicate keys and round numbers through
// float64, none of which a formatter should do.
func FormatJSON(input []byte, opts JSONOptions) ([]byte, error) {
	var buf bytes.Buffer
	var err error

	if opts.Minify || opts.Indent == 0 {
		err = json.Compact(&buf, input)
	} else {
		err = json.Indent(&buf, input, "", strings.Repeat(" ", opts.Indent))
	}
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if buf.Len() == 0 {
		return nil, fmt.Errorf("invalid JSON: empty input")
	}

	buf.WriteByte('\n')
	return buf.Bytes(), nil
}
