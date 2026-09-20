package convert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/agejevasv/swk/internal/jsonx"
)

func JSONToYAML(input []byte) ([]byte, error) {
	var data any
	if err := jsonx.Decode(input, &data); err != nil {
		return nil, err
	}

	data = numbersToNodes(data)

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// numbersToNodes replaces json.Number values with YAML scalar nodes so the
// original literal is emitted unquoted, instead of being encoded as a string
// (json.Number) or rounded through float64.
func numbersToNodes(v any) any {
	switch v := v.(type) {
	case json.Number:
		return numberNode(v)
	case map[string]any:
		for key, val := range v {
			v[key] = numbersToNodes(val)
		}
		return v
	case []any:
		for i, val := range v {
			v[i] = numbersToNodes(val)
		}
		return v
	default:
		return v
	}
}

func numberNode(n json.Number) *yaml.Node {
	tag := "!!int"
	if strings.ContainsAny(string(n), ".eE") {
		tag = "!!float"
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: tag, Value: string(n)}
}

// YAMLToJSON converts YAML to JSON. indent is the number of spaces per level;
// 0 produces compact output. Negative values are rejected.
func YAMLToJSON(input []byte, indent int) ([]byte, error) {
	if indent < 0 {
		return nil, fmt.Errorf("indent must be >= 0, got %d", indent)
	}

	var data any
	if err := yaml.Unmarshal(input, &data); err != nil {
		return nil, err
	}

	// yaml.v3 unmarshals map keys as string, but we need to ensure
	// nested maps are map[string]any for JSON marshaling.
	data = convertMaps(data)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", strings.Repeat(" ", indent))
	enc.SetEscapeHTML(false)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func convertMaps(v any) any {
	switch v := v.(type) {
	case map[string]any:
		for key, val := range v {
			v[key] = convertMaps(val)
		}
		return v
	case map[any]any:
		m := make(map[string]any, len(v))
		for key, val := range v {
			k, ok := key.(string)
			if !ok {
				k = fmt.Sprintf("%v", key)
			}
			m[k] = convertMaps(val)
		}
		return m
	case []any:
		for i, val := range v {
			v[i] = convertMaps(val)
		}
		return v
	default:
		return v
	}
}
