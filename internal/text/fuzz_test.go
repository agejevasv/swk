package text

import "testing"

func FuzzConvertCase(f *testing.F) {
	for _, s := range []string{"helloWorld", "hello world", "HTTPServer", "", "_", "ABC", "a-b_c.d/e", "日本語テスト"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		for _, to := range []string{"camel", "pascal", "snake", "kebab", "upper", "lower", "title", "sentence", "dot", "path", "bogus"} {
			_, _ = ConvertCase(s, to)
		}
	})
}

func FuzzEscapeRoundTrip(f *testing.F) {
	for _, s := range []string{"", "a", `"quoted"`, "<tag>", "it's", "a\nb", "\x00", "日本語", `\`} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		for _, mode := range []string{"json", "xml", "html", "shell"} {
			escaped, err := Escape(s, mode)
			if err != nil {
				continue
			}
			back, err := Unescape(escaped, mode)
			if err != nil {
				t.Fatalf("mode %s: escaping %q produced %q, which does not unescape: %v", mode, s, escaped, err)
			}
			if mode == "json" || mode == "shell" {
				// These two are exact inverses; xml/html entity handling is not.
				if back != s {
					t.Fatalf("mode %s: round trip of %q gave %q (via %q)", mode, s, back, escaped)
				}
			}
		}
	})
}

func FuzzStripMarkdown(f *testing.F) {
	for _, s := range []string{"# h", "**b**", "_i_", "![a](b)", "[a](b)", "snake_case", "```\ncode\n```", "", "> q", "---"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = RenderMarkdown([]byte(s), false, false, "")
	})
}

func FuzzInspect(f *testing.F) {
	for _, s := range []string{"", "hello world", "a\nb\nc", "日本語", "\x00\xff", "one. two! three?"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		info := Inspect(s)
		if info.Characters < 0 || info.Words < 0 || info.Lines < 0 || info.Bytes != len(s) {
			t.Fatalf("inconsistent counts for %q: %+v", s, info)
		}
	})
}

func FuzzDiff(f *testing.F) {
	f.Add("a\nb\n", "a\nc\n")
	f.Add("", "")
	f.Add("x", "")
	f.Add("a\n\n\nb", "b\n\n\na")
	f.Fuzz(func(t *testing.T, a, b string) {
		_ = Diff(a, b, 3)
		_ = Diff(a, b, 0)
	})
}
