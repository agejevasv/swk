package fmt

import "testing"

func FuzzFormatJSON(f *testing.F) {
	for _, s := range []string{`{"a":1}`, `[]`, ``, `{"a":12345678901234567890}`, `1e400`, `{"a":`, `null`, `"x"`} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		for _, indent := range []int{0, 2, 64} {
			_, _ = FormatJSON([]byte(s), JSONOptions{Indent: indent})
		}
		_, _ = FormatJSON([]byte(s), JSONOptions{Minify: true})
	})
}

func FuzzFormatXML(f *testing.F) {
	for _, s := range []string{"<a/>", "<a><b>1</b></a>", "", "<a>", "<!-- c -->", "<?xml version=\"1.0\"?><a/>", "<a b=\"&amp;\"/>"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		for _, indent := range []int{0, 2, 64} {
			_, _ = FormatXML([]byte(s), XMLOptions{Indent: indent})
		}
		_, _ = FormatXML([]byte(s), XMLOptions{Minify: true})
	})
}
