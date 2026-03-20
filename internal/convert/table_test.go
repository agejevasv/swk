package convert

import (
	"strings"
	"testing"
)

func TestToTable_JSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		style    string
		contains []string
		wantErr  bool
	}{
		// Basic JSON array of objects.
		{
			name:  "basic_json_box_style",
			input: `[{"name":"alice","age":30},{"name":"bob","age":25}]`,
			style: "box",
			contains: []string{
				"name", "age",
				"alice", "30",
				"bob", "25",
				"┌", "┐", "└", "┘", "│", "─",
			},
		},
		{
			name:  "basic_json_simple_style",
			input: `[{"name":"alice","age":30}]`,
			style: "simple",
			contains: []string{
				"name", "age",
				"alice", "30",
				"+", "-", "|",
			},
		},
		{
			name:  "basic_json_plain_style",
			input: `[{"name":"alice","age":30}]`,
			style: "plain",
			contains: []string{
				"name", "age",
				"alice", "30",
			},
		},
		{
			name:  "plain_style_no_borders",
			input: `[{"x":"1"}]`,
			style: "plain",
			contains: []string{
				"x", "1",
			},
		},
		// Single field, single row.
		{
			name:     "single_field_single_row",
			input:    `[{"id":1}]`,
			style:    "box",
			contains: []string{"id", "1"},
		},
		// Multiple rows.
		{
			name:  "multiple_rows",
			input: `[{"k":"a"},{"k":"b"},{"k":"c"}]`,
			style: "box",
			contains: []string{
				"k", "a", "b", "c",
			},
		},
		// Nested JSON values should be stringified.
		{
			name:  "nested_object_stringified",
			input: `[{"name":"alice","meta":{"role":"admin"}}]`,
			style: "box",
			contains: []string{
				"name", "meta",
				"alice", `{"role":"admin"}`,
			},
		},
		{
			name:  "nested_array_stringified",
			input: `[{"name":"alice","tags":["go","rust"]}]`,
			style: "box",
			contains: []string{
				"name", "tags",
				"alice", `["go","rust"]`,
			},
		},
		// Missing fields across objects: second object lacks "age".
		{
			name:  "missing_fields_across_objects",
			input: `[{"name":"alice","age":30},{"name":"bob"}]`,
			style: "box",
			contains: []string{
				"name", "age",
				"alice", "30",
				"bob",
			},
		},
		// Headers union from multiple objects with different keys.
		{
			name:  "disjoint_keys_union",
			input: `[{"a":"1"},{"b":"2"}]`,
			style: "box",
			contains: []string{
				"a", "b",
				"1", "2",
			},
		},
		// Boolean and null values.
		{
			name:  "bool_and_null_values",
			input: `[{"flag":true,"nothing":null}]`,
			style: "box",
			contains: []string{
				"flag", "nothing",
				"true",
			},
		},

		// Error cases.
		{
			name:    "empty_input",
			input:   "",
			style:   "box",
			wantErr: true,
		},
		{
			name:    "whitespace_only_input",
			input:   "   ",
			style:   "box",
			wantErr: true,
		},
		{
			name:    "non_array_json_object",
			input:   `{"name":"alice"}`,
			style:   "box",
			wantErr: true,
		},
		{
			name:    "invalid_json",
			input:   `[{"name": broken}]`,
			style:   "box",
			wantErr: true,
		},
		{
			name:    "empty_json_array",
			input:   `[]`,
			style:   "box",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToTable([]byte(tt.input), tt.style, "json", ',')
			if (err != nil) != tt.wantErr {
				t.Fatalf("ToTable() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("ToTable() output missing %q\ngot:\n%s", want, got)
				}
			}
		})
	}
}

func TestToTable_CSV(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		style     string
		delimiter rune
		contains  []string
		wantErr   bool
	}{
		{
			name:      "basic_csv_comma",
			input:     "name,age\nalice,30",
			style:     "box",
			delimiter: ',',
			contains:  []string{"name", "age", "alice", "30", "│"},
		},
		{
			name:      "csv_multiple_rows",
			input:     "name,age\nalice,30\nbob,25",
			style:     "box",
			delimiter: ',',
			contains:  []string{"name", "age", "alice", "30", "bob", "25"},
		},
		{
			name:      "csv_tab_delimiter",
			input:     "name\tage\nalice\t30",
			style:     "box",
			delimiter: '\t',
			contains:  []string{"name", "age", "alice", "30"},
		},
		{
			name:      "csv_semicolon_delimiter",
			input:     "name;age\nalice;30",
			style:     "box",
			delimiter: ';',
			contains:  []string{"name", "age", "alice", "30"},
		},
		{
			name:      "csv_simple_style",
			input:     "x,y\n1,2",
			style:     "simple",
			delimiter: ',',
			contains:  []string{"x", "y", "1", "2", "+", "-", "|"},
		},
		{
			name:      "csv_plain_style",
			input:     "x,y\n1,2",
			style:     "plain",
			delimiter: ',',
			contains:  []string{"x", "y", "1", "2"},
		},
		{
			name:      "csv_headers_only_no_data_rows",
			input:     "name,age",
			style:     "box",
			delimiter: ',',
			contains:  []string{"name", "age"},
		},

		// Error cases.
		{
			name:      "empty_csv",
			input:     "",
			style:     "box",
			delimiter: ',',
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToTable([]byte(tt.input), tt.style, "csv", tt.delimiter)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ToTable() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("ToTable() output missing %q\ngot:\n%s", want, got)
				}
			}
		})
	}
}

func TestToTable_BoxStyleBorders(t *testing.T) {
	input := `[{"a":"1","b":"2"}]`
	got, err := ToTable([]byte(input), "box", "json", ',')
	if err != nil {
		t.Fatalf("ToTable() error = %v", err)
	}

	// Box style must contain box-drawing characters.
	for _, ch := range []string{"┌", "┐", "└", "┘", "├", "┤", "┬", "┴", "┼", "│", "─"} {
		if !strings.Contains(got, ch) {
			t.Errorf("box style output missing border character %q\ngot:\n%s", ch, got)
		}
	}

	// Plain style must NOT contain any box-drawing characters.
	plain, err := ToTable([]byte(input), "plain", "json", ',')
	if err != nil {
		t.Fatalf("ToTable() plain error = %v", err)
	}
	for _, ch := range []string{"┌", "┐", "└", "┘", "│", "─", "+", "-", "|"} {
		if strings.Contains(plain, ch) {
			t.Errorf("plain style output should not contain %q\ngot:\n%s", ch, plain)
		}
	}
}

func TestToTable_SimpleStyleBorders(t *testing.T) {
	input := `[{"a":"1"}]`
	got, err := ToTable([]byte(input), "simple", "json", ',')
	if err != nil {
		t.Fatalf("ToTable() error = %v", err)
	}

	// Simple style uses ASCII characters, not box-drawing.
	for _, ch := range []string{"+", "-", "|"} {
		if !strings.Contains(got, ch) {
			t.Errorf("simple style output missing %q\ngot:\n%s", ch, got)
		}
	}
	// Should NOT contain box-drawing characters.
	for _, ch := range []string{"┌", "┐", "└", "┘", "│", "─"} {
		if strings.Contains(got, ch) {
			t.Errorf("simple style output should not contain box-drawing char %q\ngot:\n%s", ch, got)
		}
	}
}

func TestToTable_DefaultStyleIsBox(t *testing.T) {
	input := `[{"x":"1"}]`
	got, err := ToTable([]byte(input), "", "json", ',')
	if err != nil {
		t.Fatalf("ToTable() error = %v", err)
	}
	if !strings.Contains(got, "┌") {
		t.Errorf("default style should be box, but missing box-drawing chars\ngot:\n%s", got)
	}
}
