package query

import "testing"

func FuzzJSONPathQuery(f *testing.F) {
	f.Add(`{"a":[1,2]}`, "$.a[*]")
	f.Add(`{}`, "$")
	f.Add(``, "")
	f.Add(`[1]`, "$[?(@>0)]")
	f.Fuzz(func(t *testing.T, doc, expr string) {
		_, _ = JSONPathQuery([]byte(doc), expr)
	})
}

func FuzzRegexTest(f *testing.F) {
	f.Add("hello", "l+")
	f.Add("", "")
	f.Add("a", "(")
	f.Add("abc", "(a)(b)?(z)?")
	f.Fuzz(func(t *testing.T, input, pattern string) {
		res, err := RegexTest(input, pattern, true)
		if err != nil {
			return
		}
		// Reported offsets must be valid slice bounds into the input.
		for _, m := range res.Matches {
			if m.Start < 0 || m.End > len(input) || m.Start > m.End {
				t.Fatalf("match %+v out of bounds for input of length %d", m, len(input))
			}
			if input[m.Start:m.End] != m.Value {
				t.Fatalf("match value %q does not match offsets %d:%d in %q", m.Value, m.Start, m.End, input)
			}
		}
	})
}

func FuzzHTMLQuery(f *testing.F) {
	f.Add("<p>x</p>", "p")
	f.Add("", "")
	f.Add("<a href='x'>y</a>", "a")
	f.Fuzz(func(t *testing.T, doc, selector string) {
		_, _ = HTMLQuery(doc, selector, "")
		_, _ = HTMLQuery(doc, selector, "href")
	})
}
