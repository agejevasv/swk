package query

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestJSONPathQuery(t *testing.T) {
	bookstore := `{"store":{"book":[{"category":"reference","author":"Nigel Rees","title":"Sayings","price":8.95},{"category":"fiction","author":"Evelyn Waugh","title":"Sword","price":12.99}]}}`

	tests := []struct {
		name      string
		input     string
		query     string
		wantErr   bool
		checkFunc func(t *testing.T, output []byte)
	}{
		{
			name:  "root $ returns full document wrapped in array",
			input: `{"name":"Alice"}`,
			query: "$",
			checkFunc: func(t *testing.T, output []byte) {
				var v []any
				if err := json.Unmarshal(output, &v); err != nil {
					t.Fatalf("failed to unmarshal output: %v", err)
				}
				if len(v) != 1 {
					t.Fatalf("expected 1 result, got %d", len(v))
				}
				obj, ok := v[0].(map[string]any)
				if !ok {
					t.Fatalf("expected object, got %T", v[0])
				}
				if obj["name"] != "Alice" {
					t.Errorf("expected name=Alice, got %v", obj["name"])
				}
			},
		},
		{
			name:  "nested $.name",
			input: `{"name":"Alice"}`,
			query: "$.name",
			checkFunc: func(t *testing.T, output []byte) {
				var v []any
				if err := json.Unmarshal(output, &v); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(v) != 1 || v[0] != "Alice" {
					t.Errorf("got %s, want [\"Alice\"]", strings.TrimSpace(string(output)))
				}
			},
		},
		{
			name:  "deep nested $.a.b.c",
			input: `{"a":{"b":{"c":42}}}`,
			query: "$.a.b.c",
			checkFunc: func(t *testing.T, output []byte) {
				var v []any
				if err := json.Unmarshal(output, &v); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(v) != 1 || v[0] != float64(42) {
					t.Errorf("got %s, want [42]", strings.TrimSpace(string(output)))
				}
			},
		},
		{
			name:  "array index $.items[0]",
			input: `{"items":["a","b"]}`,
			query: "$.items[0]",
			checkFunc: func(t *testing.T, output []byte) {
				var v []any
				if err := json.Unmarshal(output, &v); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(v) != 1 || v[0] != "a" {
					t.Errorf("got %s, want [\"a\"]", strings.TrimSpace(string(output)))
				}
			},
		},
		{
			name:  "wildcard $.items[*]",
			input: `{"items":["a","b","c"]}`,
			query: "$.items[*]",
			checkFunc: func(t *testing.T, output []byte) {
				var v []any
				if err := json.Unmarshal(output, &v); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(v) != 3 {
					t.Errorf("expected 3 items, got %d", len(v))
				}
			},
		},
		{
			name:  "recursive descent $..name",
			input: `{"a":{"name":"x"},"b":{"name":"y"}}`,
			query: "$..name",
			checkFunc: func(t *testing.T, output []byte) {
				var v []any
				if err := json.Unmarshal(output, &v); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(v) != 2 {
					t.Errorf("expected 2 results, got %d", len(v))
				}
			},
		},
		{
			name:  "non-existent path returns nil",
			input: `{"name":"Alice"}`,
			query: "$.nonexistent",
			checkFunc: func(t *testing.T, output []byte) {
				if output != nil {
					t.Errorf("expected nil result for no match, got %s", output)
				}
			},
		},
		{
			name:    "invalid JSON returns error",
			input:   `{not json}`,
			query:   "$.name",
			wantErr: true,
		},
		{
			name:    "invalid JSONPath returns error",
			input:   `{"name":"Alice"}`,
			query:   "$[invalid[[[",
			wantErr: true,
		},
		{
			name:  "bookstore - get all authors via recursive descent",
			input: bookstore,
			query: "$..author",
			checkFunc: func(t *testing.T, output []byte) {
				var v []any
				if err := json.Unmarshal(output, &v); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(v) != 2 {
					t.Errorf("expected 2 authors, got %d", len(v))
				}
			},
		},
		{
			name:  "bookstore - first book title",
			input: bookstore,
			query: "$.store.book[0].title",
			checkFunc: func(t *testing.T, output []byte) {
				var v []any
				if err := json.Unmarshal(output, &v); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(v) != 1 || v[0] != "Sayings" {
					t.Errorf("got %s, want [\"Sayings\"]", strings.TrimSpace(string(output)))
				}
			},
		},
		{
			name:  "bookstore - all prices",
			input: bookstore,
			query: "$..price",
			checkFunc: func(t *testing.T, output []byte) {
				var v []any
				if err := json.Unmarshal(output, &v); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(v) != 2 {
					t.Errorf("expected 2 prices, got %d", len(v))
				}
			},
		},
		{
			name:  "bookstore - second book category",
			input: bookstore,
			query: "$.store.book[1].category",
			checkFunc: func(t *testing.T, output []byte) {
				var v []any
				if err := json.Unmarshal(output, &v); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(v) != 1 || v[0] != "fiction" {
					t.Errorf("got %s, want [\"fiction\"]", strings.TrimSpace(string(output)))
				}
			},
		},
		{
			name:  "large JSON object - access last element",
			input: generateLargeJSON(100),
			query: "$.items[99].id",
			checkFunc: func(t *testing.T, output []byte) {
				var v []any
				if err := json.Unmarshal(output, &v); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(v) != 1 || v[0] != float64(99) {
					t.Errorf("got %s, want [99]", strings.TrimSpace(string(output)))
				}
			},
		},
		{
			name:  "large JSON object - access first element",
			input: generateLargeJSON(100),
			query: "$.items[0].name",
			checkFunc: func(t *testing.T, output []byte) {
				var v []any
				if err := json.Unmarshal(output, &v); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(v) != 1 || v[0] != "item0" {
					t.Errorf("got %s, want [\"item0\"]", strings.TrimSpace(string(output)))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := JSONPathQuery([]byte(tt.input), tt.query)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, output)
			}
		})
	}
}

func generateLargeJSON(n int) string {
	var b strings.Builder
	b.WriteString(`{"items":[`)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(fmt.Sprintf(`{"id":%d,"name":"item%d"}`, i, i))
	}
	b.WriteString(`]}`)
	return b.String()
}
