package convert

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
)

func JSONToCSV(input []byte, delimiter rune) ([]byte, error) {
	var data []map[string]interface{}
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, fmt.Errorf("input must be a JSON array of objects: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("empty JSON array")
	}

	headers := make([]string, 0, len(data[0]))
	for k := range data[0] {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.Comma = delimiter

	if err := w.Write(headers); err != nil {
		return nil, err
	}

	for _, obj := range data {
		row := make([]string, len(headers))
		for i, h := range headers {
			row[i] = fmt.Sprintf("%v", obj[h])
		}
		if err := w.Write(row); err != nil {
			return nil, err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func CSVToJSON(input []byte, delimiter rune) ([]byte, error) {
	r := csv.NewReader(bytes.NewReader(input))
	r.Comma = delimiter

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV must have at least a header row and one data row")
	}

	headers := records[0]
	var result []map[string]interface{}

	for _, row := range records[1:] {
		obj := make(map[string]interface{}, len(headers))
		for i, h := range headers {
			if i < len(row) {
				obj[h] = row[i]
			} else {
				obj[h] = ""
			}
		}
		result = append(result, obj)
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}
	out = append(out, '\n')
	return out, nil
}
