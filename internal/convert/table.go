package convert

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

type TableStyle struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
	LeftT       string
	RightT      string
	TopT        string
	BottomT     string
	Cross       string
}

var BoxStyle = TableStyle{
	TopLeft: "┌", TopRight: "┐", BottomLeft: "└", BottomRight: "┘",
	Horizontal: "─", Vertical: "│",
	LeftT: "├", RightT: "┤", TopT: "┬", BottomT: "┴", Cross: "┼",
}

var SimpleStyle = TableStyle{
	TopLeft: "+", TopRight: "+", BottomLeft: "+", BottomRight: "+",
	Horizontal: "-", Vertical: "|",
	LeftT: "+", RightT: "+", TopT: "+", BottomT: "+", Cross: "+",
}

// ToTable converts JSON or CSV input to a formatted table string.
func ToTable(input []byte, style string, inputFormat string, delimiter rune) (string, error) {
	var headers []string
	var rows [][]string
	var err error

	switch inputFormat {
	case "csv":
		headers, rows, err = parseCSVData(input, delimiter)
	default:
		headers, rows, err = parseJSONData(input)
	}
	if err != nil {
		return "", err
	}

	if len(headers) == 0 {
		return "", fmt.Errorf("no data to display")
	}

	s := getStyle(style)
	return renderTable(headers, rows, s), nil
}

func getStyle(name string) TableStyle {
	switch name {
	case "simple":
		return SimpleStyle
	case "plain":
		return TableStyle{}
	default:
		return BoxStyle
	}
}

func parseJSONData(input []byte) ([]string, [][]string, error) {
	// Try array first
	var data []map[string]interface{}
	if err := json.Unmarshal(input, &data); err != nil {
		// Try single object, wrap in array
		var single map[string]interface{}
		if err2 := json.Unmarshal(input, &single); err2 == nil {
			data = []map[string]interface{}{single}
		} else {
			return nil, nil, fmt.Errorf("expected JSON array or object")
		}
	}
	if len(data) == 0 {
		return nil, nil, fmt.Errorf("empty JSON array")
	}

	// Collect headers preserving insertion order from first object
	seen := map[string]bool{}
	var headers []string
	for _, obj := range data {
		for k := range obj {
			if !seen[k] {
				seen[k] = true
				headers = append(headers, k)
			}
		}
	}

	var rows [][]string
	for _, obj := range data {
		row := make([]string, len(headers))
		for i, h := range headers {
			if v, ok := obj[h]; ok {
				row[i] = fmt.Sprintf("%v", v)
			}
		}
		rows = append(rows, row)
	}

	return headers, rows, nil
}

func parseCSVData(input []byte, delimiter rune) ([]string, [][]string, error) {
	r := csv.NewReader(bytes.NewReader(input))
	r.Comma = delimiter

	records, err := r.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("invalid CSV: %w", err)
	}
	if len(records) < 1 {
		return nil, nil, fmt.Errorf("empty CSV")
	}

	return records[0], records[1:], nil
}

func renderTable(headers []string, rows [][]string, s TableStyle) string {
	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = utf8.RuneCountInString(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				w := utf8.RuneCountInString(cell)
				if w > widths[i] {
					widths[i] = w
				}
			}
		}
	}

	var buf strings.Builder
	plain := s.Horizontal == ""

	if !plain {
		buf.WriteString(borderLine(widths, s.TopLeft, s.TopT, s.TopRight, s.Horizontal))
		buf.WriteByte('\n')
	}

	// Header row
	buf.WriteString(dataLine(headers, widths, s.Vertical, plain))
	buf.WriteByte('\n')

	if !plain {
		buf.WriteString(borderLine(widths, s.LeftT, s.Cross, s.RightT, s.Horizontal))
		buf.WriteByte('\n')
	}

	// Data rows
	for _, row := range rows {
		padded := make([]string, len(headers))
		for i := range headers {
			if i < len(row) {
				padded[i] = row[i]
			}
		}
		buf.WriteString(dataLine(padded, widths, s.Vertical, plain))
		buf.WriteByte('\n')
	}

	if !plain {
		buf.WriteString(borderLine(widths, s.BottomLeft, s.BottomT, s.BottomRight, s.Horizontal))
		buf.WriteByte('\n')
	}

	return buf.String()
}

func borderLine(widths []int, left, mid, right, h string) string {
	var buf strings.Builder
	buf.WriteString(left)
	for i, w := range widths {
		buf.WriteString(strings.Repeat(h, w+2))
		if i < len(widths)-1 {
			buf.WriteString(mid)
		}
	}
	buf.WriteString(right)
	return buf.String()
}

func dataLine(cells []string, widths []int, sep string, plain bool) string {
	var buf strings.Builder
	if !plain {
		buf.WriteString(sep)
		buf.WriteByte(' ')
	}
	for i, cell := range cells {
		w := utf8.RuneCountInString(cell)
		buf.WriteString(cell)
		buf.WriteString(strings.Repeat(" ", widths[i]-w))
		if !plain {
			buf.WriteByte(' ')
			buf.WriteString(sep)
			if i < len(cells)-1 {
				buf.WriteByte(' ')
			}
		} else if i < len(cells)-1 {
			buf.WriteString("  ")
		}
	}
	return buf.String()
}
