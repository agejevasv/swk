package fmt

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JSONOptions struct {
	Indent int
	Minify bool
}

func FormatJSON(input []byte, opts JSONOptions) ([]byte, error) {
	var data any
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if opts.Minify {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(data); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	indent := "  "
	if opts.Indent > 0 {
		indent = ""
		for i := 0; i < opts.Indent; i++ {
			indent += " "
		}
	}

	result, err := json.MarshalIndent(data, "", indent)
	if err != nil {
		return nil, err
	}

	return append(result, '\n'), nil
}
