package convert

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONToYAML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "simple_object",
			input: `{"name":"alice","age":30}`,
			want:  "age: 30\nname: alice\n",
		},
		{
			name:  "nested_objects_with_arrays",
			input: `{"user":{"name":"bob","tags":["admin","editor"]}}`,
			want:  "user:\n  name: bob\n  tags:\n    - admin\n    - editor\n",
		},
		{
			name:  "multi_line_pretty_json_input",
			input: "{\n  \"a\": 1,\n  \"b\": 2\n}\n",
			want:  "a: 1\nb: 2\n",
		},
		{
			name:  "simple_string_values",
			input: `{"key":"value","name":"alice"}`,
			want:  "key: value\nname: alice\n",
		},
		{
			name:  "empty_object",
			input: `{}`,
			want:  "{}\n",
		},
		{
			name:  "array_of_numbers",
			input: `[1,2,3]`,
			want:  "- 1\n- 2\n- 3\n",
		},
		{
			name:  "all_json_types",
			input: `{"a":null,"b":true,"c":42,"d":"hello","e":[1],"f":{}}`,
			want:  "a: null\nb: true\nc: 42\nd: hello\ne:\n  - 1\nf: {}\n",
		},
		{
			name:  "deeply_nested",
			input: `{"l1":{"l2":{"l3":{"l4":"deep"}}}}`,
			want:  "l1:\n  l2:\n    l3:\n      l4: deep\n",
		},
		{
			name:    "invalid_json",
			input:   `{bad json`,
			wantErr: true,
		},
		{
			name:    "invalid_json_trailing_comma",
			input:   `{"a":1,}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONToYAML([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Fatalf("JSONToYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && string(got) != tt.want {
				t.Errorf("JSONToYAML() =\n%q\nwant\n%q", string(got), tt.want)
			}
		})
	}
}

func TestYAMLToJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		indent  int
		want    string
		wantErr bool
	}{
		{
			name:   "simple_object",
			input:  "name: alice\nage: 30\n",
			indent: 2,
			want:   "{\n  \"age\": 30,\n  \"name\": \"alice\"\n}\n",
		},
		{
			name:   "array_of_numbers",
			input:  "- 1\n- 2\n- 3\n",
			indent: 2,
			want:   "[\n  1,\n  2,\n  3\n]\n",
		},
		{
			name:   "nested_object",
			input:  "user:\n  name: bob\n",
			indent: 2,
			want:   "{\n  \"user\": {\n    \"name\": \"bob\"\n  }\n}\n",
		},
		{
			name:   "indent_4",
			input:  "x: 1\n",
			indent: 4,
			want:   "{\n    \"x\": 1\n}\n",
		},
		{
			name:   "indent_0_compact",
			input:  "x: 1\n",
			indent: 0,
			want:   "{\"x\":1}\n",
		},
		{
			name:   "empty_object_yaml",
			input:  "{}\n",
			indent: 2,
			want:   "{}\n",
		},
		{
			name:   "boolean_and_null",
			input:  "flag: true\nval: null\n",
			indent: 2,
			want:   "{\n  \"flag\": true,\n  \"val\": null\n}\n",
		},
		{
			name:    "invalid_yaml",
			input:   ":\n  :\n    - :\n      invalid: [",
			indent:  2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := YAMLToJSON([]byte(tt.input), tt.indent)
			if (err != nil) != tt.wantErr {
				t.Fatalf("YAMLToJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.want != "" {
				if string(got) != tt.want {
					t.Errorf("YAMLToJSON() =\n%q\nwant\n%q", string(got), tt.want)
				}
			}
		})
	}
}

func TestJSONYAMLRoundtrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "object_with_array",
			input: `{"items":[{"id":1,"name":"first"},{"id":2,"name":"second"}]}`,
		},
		{
			name:  "simple_string_value",
			input: `{"greeting":"hello world"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yamlBytes, err := JSONToYAML([]byte(tt.input))
			if err != nil {
				t.Fatalf("JSONToYAML: %v", err)
			}
			jsonBytes, err := YAMLToJSON(yamlBytes, 0)
			if err != nil {
				t.Fatalf("YAMLToJSON: %v", err)
			}

			// Parse both original and roundtripped to compare structurally.
			var orig, roundtripped any
			if err := json.Unmarshal([]byte(tt.input), &orig); err != nil {
				t.Fatalf("unmarshal original: %v", err)
			}
			if err := json.Unmarshal(jsonBytes, &roundtripped); err != nil {
				t.Fatalf("unmarshal roundtripped: %v", err)
			}

			origJSON, _ := json.Marshal(orig)
			rtJSON, _ := json.Marshal(roundtripped)
			if string(origJSON) != string(rtJSON) {
				t.Errorf("roundtrip mismatch:\norig: %s\ngot:  %s", origJSON, rtJSON)
			}
		})
	}
}

func TestConvertMaps(t *testing.T) {
	// Test the map[any]any branch directly since yaml.v3 doesn't produce it
	input := map[any]any{
		"key":    "value",
		42:       "numeric key",
		"nested": map[any]any{"inner": "data"},
	}
	result := convertMaps(input)
	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}
	if m["key"] != "value" {
		t.Errorf("expected key=value, got %v", m["key"])
	}
	if m["42"] != "numeric key" {
		t.Errorf("expected 42=numeric key, got %v", m["42"])
	}
	nested, ok := m["nested"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested map[string]any, got %T", m["nested"])
	}
	if nested["inner"] != "data" {
		t.Errorf("expected inner=data, got %v", nested["inner"])
	}
}

func TestConvertMaps_Slice(t *testing.T) {
	input := []any{"a", map[any]any{"k": "v"}}
	result := convertMaps(input)
	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", result)
	}
	if arr[0] != "a" {
		t.Errorf("expected arr[0]=a, got %v", arr[0])
	}
	m, ok := arr[1].(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", arr[1])
	}
	if m["k"] != "v" {
		t.Errorf("expected k=v, got %v", m["k"])
	}
}

func TestConvertMaps_Scalar(t *testing.T) {
	if convertMaps(42) != 42 {
		t.Error("scalar should pass through unchanged")
	}
	if convertMaps("hello") != "hello" {
		t.Error("string should pass through unchanged")
	}
}

func TestYAMLToJSON_HTMLNotEscaped(t *testing.T) {
	// Verify SetEscapeHTML(false) works: & < > should not be escaped.
	input := "url: https://example.com?a=1&b=2\n"
	got, err := YAMLToJSON([]byte(input), 2)
	if err != nil {
		t.Fatalf("YAMLToJSON: %v", err)
	}
	if strings.Contains(string(got), `\u0026`) {
		t.Errorf("YAMLToJSON should not escape HTML, got: %s", got)
	}
}
