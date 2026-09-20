package fmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/agejevasv/swk/internal/jsonx"
)

type JSONOptions struct {
	// Indent is the number of spaces per level. Values <= 0 use the default of 2.
	Indent int
	Minify bool
}

const defaultIndent = 2

func FormatJSON(input []byte, opts JSONOptions) ([]byte, error) {
	var data any
	if err := jsonx.Decode(input, &data); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	// Keep string contents verbatim; <, > and & are valid inside JSON strings.
	enc.SetEscapeHTML(false)
	if !opts.Minify {
		enc.SetIndent("", indentString(opts.Indent))
	}

	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func indentString(n int) string {
	if n <= 0 {
		n = defaultIndent
	}
	return strings.Repeat(" ", n)
}
