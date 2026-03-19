package convert

import (
	"bytes"
	"encoding/json"

	"gopkg.in/yaml.v3"
)

func JSONToYAML(input []byte) ([]byte, error) {
	var data interface{}
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
	var data interface{}
	if err := yaml.Unmarshal(input, &data); err != nil {
		return nil, err
	}

	// yaml.v3 unmarshals map keys as string, but we need to ensure
	// nested maps are map[string]interface{} for JSON marshaling.
	data = convertMaps(data)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", spaces(indent))
	enc.SetEscapeHTML(false)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func spaces(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}

func convertMaps(v interface{}) interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		for key, val := range v {
			v[key] = convertMaps(val)
		}
		return v
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(v))
		for key, val := range v {
			k, _ := key.(string)
			m[k] = convertMaps(val)
		}
		return m
	case []interface{}:
		for i, val := range v {
			v[i] = convertMaps(val)
		}
		return v
	default:
		return v
	}
}
