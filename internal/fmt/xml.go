package fmt

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// XMLOptions holds options for XML formatting.
type XMLOptions struct {
	Indent int
	Minify bool
}

func FormatXML(input []byte, opts XMLOptions) ([]byte, error) {
	decoder := xml.NewDecoder(bytes.NewReader(input))
	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)

	if opts.Minify {
		encoder.Indent("", "")
	} else {
		indent := strings.Repeat(" ", opts.Indent)
		encoder.Indent("", indent)
	}

	foundToken := false
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("invalid XML: %w", err)
		}
		foundToken = true

		// Skip whitespace-only CharData so the encoder can re-indent properly
		if cd, ok := token.(xml.CharData); ok {
			if strings.TrimSpace(string(cd)) == "" {
				continue
			}
		}

		if err := encoder.EncodeToken(token); err != nil {
			return nil, err
		}
	}

	if !foundToken {
		return nil, fmt.Errorf("invalid XML: empty input")
	}

	if err := encoder.Flush(); err != nil {
		return nil, err
	}

	result := buf.Bytes()

	if !opts.Minify && len(result) > 0 && result[len(result)-1] != '\n' {
		result = append(result, '\n')
	}

	return result, nil
}
