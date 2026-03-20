package query

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

func ValidateXML(input []byte) error {
	decoder := xml.NewDecoder(strings.NewReader(string(input)))
	hasElement := false
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("invalid XML: %w", err)
		}
		if _, ok := tok.(xml.StartElement); ok {
			hasElement = true
		}
	}
	if !hasElement {
		return fmt.Errorf("invalid XML: no XML elements found")
	}
	return nil
}
