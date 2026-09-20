// Package jsonx provides JSON decoding that preserves number literals.
//
// encoding/json decodes JSON numbers into float64 by default. That silently
// corrupts integers beyond 2^53 (9007199254740993 becomes 9007199254740992)
// and re-renders exact values in scientific notation (1000000 becomes 1e+06)
// once they are printed with %v. Decode keeps every number as a json.Number,
// so the original literal survives a decode/encode round trip.
package jsonx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Decode unmarshals data into v like json.Unmarshal, but keeps numbers as
// json.Number instead of float64. Like json.Unmarshal, it rejects trailing
// content after the top-level value.
func Decode(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	if err := dec.Decode(v); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("unexpected end of JSON input")
		}
		return err
	}

	if dec.More() {
		return fmt.Errorf("invalid character after top-level value")
	}

	return nil
}
