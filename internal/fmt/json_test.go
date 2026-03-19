package fmt

import (
	"testing"
)

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		opts    JSONOptions
		want    string
		wantErr bool
	}{
		// Pretty-print variations.
		{
			name:  "pretty_print_default_indent_2",
			input: `{"name":"John","age":30}`,
			opts:  JSONOptions{Indent: 2},
			want:  "{\n  \"age\": 30,\n  \"name\": \"John\"\n}",
		},
		{
			name:  "pretty_print_indent_4",
			input: `{"a":1}`,
			opts:  JSONOptions{Indent: 4},
			want:  "{\n    \"a\": 1\n}",
		},

		// Minify.
		{
			name:  "minify_removes_all_whitespace",
			input: "{\n  \"name\": \"John\",\n  \"age\": 30\n}",
			opts:  JSONOptions{Minify: true},
			want:  `{"age":30,"name":"John"}`,
		},
		{
			name:  "minify_array",
			input: `[ 1 , 2 , 3 ]`,
			opts:  JSONOptions{Minify: true},
			want:  `[1,2,3]`,
		},

		// Sort keys (default behavior since json.Marshal sorts by default).
		{
			name:  "sort_keys_alphabetical",
			input: `{"z":1,"a":2,"m":3}`,
			want:  "{\n  \"a\": 2,\n  \"m\": 3,\n  \"z\": 1\n}",
		},

		// All JSON types.
		{
			name:  "string_type",
			input: `{"k":"hello"}`,
			opts:  JSONOptions{Indent: 2},
			want:  "{\n  \"k\": \"hello\"\n}",
		},
		{
			name:  "integer_number",
			input: `{"k":42}`,
			opts:  JSONOptions{Indent: 2},
			want:  "{\n  \"k\": 42\n}",
		},
		{
			name:  "float_number",
			input: `{"k":3.14}`,
			opts:  JSONOptions{Indent: 2},
			want:  "{\n  \"k\": 3.14\n}",
		},
		{
			name:  "negative_number",
			input: `{"k":-100}`,
			opts:  JSONOptions{Indent: 2},
			want:  "{\n  \"k\": -100\n}",
		},
		{
			name:  "scientific_number",
			input: `{"k":1e10}`,
			opts:  JSONOptions{Minify: true},
			want:  `{"k":10000000000}`,
		},
		{
			name:  "bool_true",
			input: `{"k":true}`,
			opts:  JSONOptions{Indent: 2},
			want:  "{\n  \"k\": true\n}",
		},
		{
			name:  "bool_false",
			input: `{"k":false}`,
			opts:  JSONOptions{Indent: 2},
			want:  "{\n  \"k\": false\n}",
		},
		{
			name:  "null_value",
			input: `{"k":null}`,
			opts:  JSONOptions{Indent: 2},
			want:  "{\n  \"k\": null\n}",
		},
		{
			name:  "array_type",
			input: `{"k":[1,2,3]}`,
			opts:  JSONOptions{Indent: 2},
			want:  "{\n  \"k\": [\n    1,\n    2,\n    3\n  ]\n}",
		},
		{
			name:  "nested_object_type",
			input: `{"k":{"inner":"val"}}`,
			opts:  JSONOptions{Indent: 2},
			want:  "{\n  \"k\": {\n    \"inner\": \"val\"\n  }\n}",
		},

		// Deeply nested (5 levels).
		{
			name:  "deeply_nested_5_levels",
			input: `{"l1":{"l2":{"l3":{"l4":{"l5":"deep"}}}}}`,
			opts:  JSONOptions{Indent: 2},
			want:  "{\n  \"l1\": {\n    \"l2\": {\n      \"l3\": {\n        \"l4\": {\n          \"l5\": \"deep\"\n        }\n      }\n    }\n  }\n}",
		},

		// Idempotent operations.
		{
			name:  "already_pretty_is_idempotent",
			input: "{\n  \"a\": 1,\n  \"b\": 2\n}",
			opts:  JSONOptions{Indent: 2},
			want:  "{\n  \"a\": 1,\n  \"b\": 2\n}",
		},
		{
			name:  "already_minified_is_idempotent",
			input: `{"a":1,"b":2}`,
			opts:  JSONOptions{Minify: true},
			want:  `{"a":1,"b":2}`,
		},

		// Top-level array.
		{
			name:  "top_level_array_pretty",
			input: `[1,2,3]`,
			opts:  JSONOptions{Indent: 2},
			want:  "[\n  1,\n  2,\n  3\n]",
		},

		// Error cases.
		{
			name:    "invalid_json",
			input:   `{invalid}`,
			opts:    JSONOptions{Indent: 2},
			wantErr: true,
		},
		{
			name:    "invalid_json_trailing_comma",
			input:   `{"a":1,}`,
			opts:    JSONOptions{Indent: 2},
			wantErr: true,
		},
		{
			name:    "empty_string",
			input:   "",
			opts:    JSONOptions{Indent: 2},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatJSON([]byte(tt.input), tt.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FormatJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && string(got) != tt.want {
				t.Errorf("FormatJSON() =\n%s\nwant:\n%s", string(got), tt.want)
			}
		})
	}
}
