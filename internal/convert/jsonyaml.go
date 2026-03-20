package convert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func JSONToYAML(input []byte) ([]byte, error) {
	var data any
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func YAMLToJSON(input []byte, indent int) ([]byte, error) {
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
