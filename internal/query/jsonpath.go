package query

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

func JSONPathQuery(input []byte, query string) ([]byte, error) {
	// oj.Parse accepts empty input and returns a nil document, which then
	// panics inside the evaluator for filter expressions.
	if len(bytes.TrimSpace(input)) == 0 {
		return nil, fmt.Errorf("invalid JSON: empty input")
	}

	obj, err := oj.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	expr, err := jp.ParseString(query)
	if err != nil {
		return nil, fmt.Errorf("invalid JSONPath expression: %w", err)
	}

	results, err := evalExpr(expr, obj)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}

	out, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	return append(out, '\n'), nil
}

// evalExpr isolates the third-party evaluator, which panics on some
// combinations of expression and document rather than returning an error.
func evalExpr(expr jp.Expr, obj any) (results []any, err error) {
	defer func() {
		if r := recover(); r != nil {
			results = nil
			err = fmt.Errorf("cannot evaluate JSONPath expression against this document: %v", r)
		}
	}()
	return expr.Get(obj), nil
}
