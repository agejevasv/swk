package text

import (
	"strings"
	"testing"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		name         string
		a            string
		b            string
		context      int
		wantEmpty    bool
		wantContains []string
		wantMissing  []string
	}{
		{
			name:      "identical strings returns empty",
			a:         "hello\nworld\n",
			b:         "hello\nworld\n",
			context:   3,
			wantEmpty: true,
		},
		{
			name:         "single line added",
			a:            "hello\n",
			b:            "hello\nworld\n",
			context:      3,
			wantContains: []string{"+world"},
		},
		{
			name:         "single line deleted",
			a:            "hello\nworld\n",
			b:            "hello\n",
			context:      3,
			wantContains: []string{"-world"},
		},
		{
			name:         "single line modified",
			a:            "hello\n",
			b:            "world\n",
			context:      3,
			wantContains: []string{"-hello", "+world"},
		},
		{
			name:         "multi-line add/remove in middle",
			a:            "aaa\nbbb\nccc\n",
			b:            "aaa\nxxx\nccc\n",
			context:      3,
			wantContains: []string{"-bbb", "+xxx"},
		},
		{
			name:         "context lines appear around changes",
			a:            "line1\nline2\nline3\nline4\nline5\n",
			b:            "line1\nline2\nchanged\nline4\nline5\n",
			context:      1,
			wantContains: []string{" line2", "-line3", "+changed", " line4"},
		},
		{
			name:         "context 0 only changed lines",
			a:            "line1\nline2\nline3\n",
			b:            "line1\nchanged\nline3\n",
			context:      0,
			wantContains: []string{"-line2", "+changed"},
			wantMissing:  []string{" line1", " line3"},
		},
		{
			name:         "complete rewrite",
			a:            "aaa\nbbb\n",
			b:            "xxx\nyyy\n",
			context:      3,
			wantContains: []string{"-aaa", "-bbb", "+xxx", "+yyy"},
		},
		{
			name:         "empty vs non-empty",
			a:            "",
			b:            "hello\n",
			context:      3,
			wantContains: []string{"+hello"},
		},
		{
			name:         "has unified diff header",
			a:            "aaa\n",
			b:            "bbb\n",
			context:      3,
			wantContains: []string{"--- a", "+++ b", "@@"},
		},
		{
			name:         "multiple additions",
			a:            "line1\nline3\n",
			b:            "line1\nline2\nline3\nline4\n",
			context:      1,
			wantContains: []string{"+line2", "+line4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Diff(tt.a, tt.b, tt.context)
			if tt.wantEmpty {
				if got != "" {
					t.Errorf("expected empty diff, got:\n%s", got)
				}
				return
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("diff output missing %q.\nGot:\n%s", want, got)
				}
			}
			for _, notWant := range tt.wantMissing {
				if strings.Contains(got, notWant) {
					t.Errorf("diff output should not contain %q.\nGot:\n%s", notWant, got)
				}
			}
		})
	}
}
