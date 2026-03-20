package query

import (
	"encoding/json"
	"fmt"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

func JSONPathQuery(input []byte, query string) ([]byte, error) {
	obj, err := oj.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	expr, err := jp.ParseString(query)
	if err != nil {
		return nil, fmt.Errorf("invalid JSONPath expression: %w", err)
	}

	results := expr.Get(obj)
	if len(results) == 0 {
		return nil, nil
	}

	var output any
	if len(results) == 1 {
		output = results[0]
	} else {
		output = results
	}

	out, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	return append(out, '\n'), nil
}
