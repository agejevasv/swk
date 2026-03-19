package fmt

import (
	"regexp"
	"sort"
	"strings"
)

// SQLOptions holds options for SQL formatting.
type SQLOptions struct {
	Uppercase bool
}

var majorKeywords = []string{
	"SELECT", "FROM", "WHERE", "INNER JOIN", "LEFT JOIN", "RIGHT JOIN",
	"OUTER JOIN", "CROSS JOIN", "JOIN", "ORDER BY", "GROUP BY",
	"HAVING", "LIMIT", "INSERT INTO", "INSERT", "UPDATE", "DELETE FROM",
	"DELETE", "CREATE", "ALTER", "DROP", "SET", "VALUES", "ON",
	"AND", "OR", "UNION ALL", "UNION", "CASE", "WHEN", "THEN", "ELSE", "END",
}

var subClauseKeywords = map[string]bool{
	"AND": true, "OR": true, "ON": true,
	"WHEN": true, "THEN": true, "ELSE": true, "END": true,
}

func FormatSQL(input []byte, opts SQLOptions) ([]byte, error) {
	sql := string(input)

	spaceRe := regexp.MustCompile(`\s+`)
	sql = spaceRe.ReplaceAllString(sql, " ")

	sorted := make([]string, len(majorKeywords))
	copy(sorted, majorKeywords)
	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i]) > len(sorted[j])
	})

	patterns := make([]string, len(sorted))
	for i, kw := range sorted {
		patterns[i] = `\b` + regexp.QuoteMeta(kw) + `\b`
	}
	kwRe := regexp.MustCompile(`(?i)(` + strings.Join(patterns, "|") + `)`)

	type segment struct {
		keyword string
		text    string
	}

	indices := kwRe.FindAllStringIndex(sql, -1)
	var segments []segment

	if len(indices) == 0 {
		result := sql
		if opts.Uppercase {
			result = uppercaseKeywords(result, kwRe)
		}
		return []byte(result), nil
	}

	if indices[0][0] > 0 {
		pre := strings.TrimSpace(sql[:indices[0][0]])
		if pre != "" {
			segments = append(segments, segment{keyword: "", text: pre})
		}
	}

	for i, idx := range indices {
		kw := strings.ToUpper(sql[idx[0]:idx[1]])
		var text string
		if i+1 < len(indices) {
			text = strings.TrimSpace(sql[idx[1]:indices[i+1][0]])
		} else {
			text = strings.TrimSpace(sql[idx[1]:])
		}
		segments = append(segments, segment{keyword: kw, text: text})
	}

	var sb strings.Builder
	for _, seg := range segments {
		kw := seg.keyword
		if !opts.Uppercase {
			kw = strings.ToLower(kw)
		}

		if seg.keyword == "" {
			sb.WriteString(seg.text)
			continue
		}

		if sb.Len() > 0 {
			sb.WriteString("\n")
		}

		if subClauseKeywords[seg.keyword] {
			sb.WriteString("  ")
		}

		sb.WriteString(kw)
		if seg.text != "" {
			sb.WriteString(" ")
			sb.WriteString(seg.text)
		}
	}

	return []byte(sb.String()), nil
}

func uppercaseKeywords(sql string, kwRe *regexp.Regexp) string {
	return kwRe.ReplaceAllStringFunc(sql, func(match string) string {
		return strings.ToUpper(match)
	})
}
