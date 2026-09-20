package text

import (
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func Diff(a, b string, contextLines int) string {
	if a == b {
		return ""
	}

	dmp := diffmatchpatch.New()
	aLines, bLines, lineArray := dmp.DiffLinesToChars(a, b)
	diffs := dmp.DiffMain(aLines, bLines, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	return formatUnifiedDiff(diffs, contextLines,
		a != "" && !strings.HasSuffix(a, "\n"),
		b != "" && !strings.HasSuffix(b, "\n"))
}

// hunkRange renders one side of a hunk header. A zero-length range starts at
// the line before it, and a single-line range omits the count, both of which
// diff(1) requires.
func hunkRange(start, count int) string {
	switch count {
	case 0:
		return fmt.Sprintf("%d,0", start-1)
	case 1:
		return fmt.Sprintf("%d", start)
	default:
		return fmt.Sprintf("%d,%d", start, count)
	}
}

const noNewlineMarker = "\\ No newline at end of file"

func formatUnifiedDiff(diffs []diffmatchpatch.Diff, contextLines int, aNoNewline, bNoNewline bool) string {
	type diffLine struct {
		op   diffmatchpatch.Operation
		text string
	}

	var lines []diffLine
	for _, d := range diffs {
		parts := strings.Split(d.Text, "\n")
		// If the last element is empty (trailing newline), don't add an empty line
		for i, p := range parts {
			if i == len(parts)-1 && p == "" {
				continue
			}
			lines = append(lines, diffLine{op: d.Type, text: p})
		}
	}

	if len(lines) == 0 {
		return ""
	}

	// The last line of each side needs the "no newline" marker when that input
	// did not end with one.
	lastA, lastB := -1, -1
	for i, l := range lines {
		if l.op != diffmatchpatch.DiffInsert {
			lastA = i
		}
		if l.op != diffmatchpatch.DiffDelete {
			lastB = i
		}
	}

	var out strings.Builder
	out.WriteString("--- a\n+++ b\n")

	show := make([]bool, len(lines))
	for i, l := range lines {
		if l.op != diffmatchpatch.DiffEqual {
			// Show this line and surrounding context
			for j := max(0, i-contextLines); j <= min(len(lines)-1, i+contextLines); j++ {
				show[j] = true
			}
		}
	}

	aLine, bLine := 1, 1
	i := 0
	for i < len(lines) {
		if !show[i] {
			switch lines[i].op {
			case diffmatchpatch.DiffEqual:
				aLine++
				bLine++
			case diffmatchpatch.DiffDelete:
				aLine++
			case diffmatchpatch.DiffInsert:
				bLine++
			}
			i++
			continue
		}

		hunkStart := i
		hunkEnd := i
		for hunkEnd < len(lines) && show[hunkEnd] {
			hunkEnd++
		}

		aCount, bCount := 0, 0
		aStart, bStart := aLine, bLine
		for j := hunkStart; j < hunkEnd; j++ {
			switch lines[j].op {
			case diffmatchpatch.DiffEqual:
				aCount++
				bCount++
			case diffmatchpatch.DiffDelete:
				aCount++
			case diffmatchpatch.DiffInsert:
				bCount++
			}
		}

		out.WriteString(fmt.Sprintf("@@ -%s +%s @@\n", hunkRange(aStart, aCount), hunkRange(bStart, bCount)))

		for j := hunkStart; j < hunkEnd; j++ {
			switch lines[j].op {
			case diffmatchpatch.DiffEqual:
				out.WriteString(" " + lines[j].text + "\n")
				aLine++
				bLine++
			case diffmatchpatch.DiffDelete:
				out.WriteString("-" + lines[j].text + "\n")
				aLine++
			case diffmatchpatch.DiffInsert:
				out.WriteString("+" + lines[j].text + "\n")
				bLine++
			}

			if (j == lastA && aNoNewline && lines[j].op != diffmatchpatch.DiffInsert) ||
				(j == lastB && bNoNewline && lines[j].op != diffmatchpatch.DiffDelete) {
				out.WriteString(noNewlineMarker + "\n")
			}
		}

		i = hunkEnd
	}

	return out.String()
}
