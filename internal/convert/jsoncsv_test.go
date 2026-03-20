package convert

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONToCSV(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter rune
		want      string
		wantErr   bool
	}{
		{
			name:      "basic_array_of_objects",
			input:     `[{"name":"alice","age":"30"},{"name":"bob","age":"25"}]`,
			delimiter: ',',
			want:      "age,name\n30,alice\n25,bob\n",
		},
		{
			name:      "csv_with_quoted_fields_containing_commas",
			input:     `[{"desc":"hello, world","id":"1"}]`,
			delimiter: ',',
			want:      "desc,id\n\"hello, world\",1\n",
		},
		{
			name:      "tab_delimiter",
			input:     `[{"a":"1","b":"2"}]`,
			delimiter: '\t',
			want:      "a\tb\n1\t2\n",
		},
		{
			name:      "single_column",
			input:     `[{"name":"alice"},{"name":"bob"}]`,
			delimiter: ',',
			want:      "name\nalice\nbob\n",
		},
		{
			name:      "numeric_values_formatted_via_sprintf",
			input:     `[{"val":42}]`,
			delimiter: ',',
			want:      "val\n42\n",
		},
		{
			name:      "empty_array",
			input:     `[]`,
			delimiter: ',',
			wantErr:   true,
		},
		{
			name:      "invalid_json",
			input:     `not json at all`,
			delimiter: ',',
			wantErr:   true,
		},
		{
			name:      "not_array_of_objects",
			input:     `{"key":"value"}`,
			delimiter: ',',
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONToCSV([]byte(tt.input), tt.delimiter)
			if (err != nil) != tt.wantErr {
				t.Fatalf("JSONToCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && string(got) != tt.want {
				t.Errorf("JSONToCSV() =\n%q\nwant\n%q", string(got), tt.want)
			}
		})
	}
}

func TestCSVToJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter rune
		contains  []string
		wantErr   bool
	}{
		{
			name:      "simple_csv",
			input:     "name,age\nalice,30\nbob,25\n",
			delimiter: ',',
			contains:  []string{`"name": "alice"`, `"age": "30"`, `"name": "bob"`, `"age": "25"`},
		},
		{
			name:      "tab_delimiter",
			input:     "a\tb\n1\t2\n",
			delimiter: '\t',
			contains:  []string{`"a": "1"`, `"b": "2"`},
		},
		{
			name:      "csv_with_quoted_fields",
			input:     "desc,id\n\"hello, world\",1\n",
			delimiter: ',',
			contains:  []string{`"desc": "hello, world"`, `"id": "1"`},
		},
		{
			name:      "single_column",
			input:     "name\nalice\nbob\n",
			delimiter: ',',
			contains:  []string{`"name": "alice"`, `"name": "bob"`},
		},
		{
			name:      "header_only_no_data",
			input:     "name,age\n",
			delimiter: ',',
			wantErr:   true,
		},
		{
			name:      "empty_input",
			input:     "",
			delimiter: ',',
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CSVToJSON([]byte(tt.input), tt.delimiter)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CSVToJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				gotStr := string(got)
				for _, want := range tt.contains {
					if !strings.Contains(gotStr, want) {
						t.Errorf("CSVToJSON() output missing %q\ngot:\n%s", want, gotStr)
					}
				}
			}
		})
	}
}

func TestCSVToJSON_OutputIsValidJSON(t *testing.T) {
	input := "name,age\nalice,30\nbob,25\n"
	got, err := CSVToJSON([]byte(input), ',')
	if err != nil {
		t.Fatalf("CSVToJSON: %v", err)
	}
	var result []map[string]any
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\ngot: %s", err, got)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 objects, got %d", len(result))
	}
}

func TestCSVToJSON_TrailingNewline(t *testing.T) {
	// CSVToJSON appends a trailing newline to output.
	input := "x\n1\n"
	got, err := CSVToJSON([]byte(input), ',')
	if err != nil {
		t.Fatalf("CSVToJSON: %v", err)
	}
	if got[len(got)-1] != '\n' {
		t.Errorf("expected trailing newline in output")
	}
}

func TestJSONCSVRoundtrip(t *testing.T) {
	// Note: roundtrip is lossy for non-string types since CSV values become strings.
	originalJSON := `[{"age":"30","name":"alice"},{"age":"25","name":"bob"}]`
	csv, err := JSONToCSV([]byte(originalJSON), ',')
	if err != nil {
		t.Fatalf("JSONToCSV: %v", err)
	}
	jsonBack, err := CSVToJSON(csv, ',')
	if err != nil {
		t.Fatalf("CSVToJSON: %v", err)
	}

	var orig, roundtripped []map[string]any
	json.Unmarshal([]byte(originalJSON), &orig)
	json.Unmarshal(jsonBack, &roundtripped)

	if len(orig) != len(roundtripped) {
		t.Fatalf("roundtrip length mismatch: %d vs %d", len(orig), len(roundtripped))
	}
	for i := range orig {
		for k, v := range orig[i] {
			if roundtripped[i][k] != v {
				t.Errorf("roundtrip mismatch at [%d][%s]: %v vs %v", i, k, v, roundtripped[i][k])
			}
		}
	}
}
