package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/agejevasv/swk/internal/ioutil"
	textLib "github.com/agejevasv/swk/internal/text"
)

// ReadTwoInputs reads two file-or-stdin inputs. "-" means stdin (at most one).
func ReadTwoInputs(args []string, stdin io.Reader) ([]byte, []byte, error) {
	if len(args) < 2 {
		return nil, nil, fmt.Errorf("requires two arguments: <file1> <file2>")
	}
	if args[0] == "-" && args[1] == "-" {
		return nil, nil, fmt.Errorf("only one argument can be stdin (-)")
	}

	a, err := readArg(args[0], stdin)
	if err != nil {
		return nil, nil, fmt.Errorf("reading first input: %w", err)
	}
	b, err := readArg(args[1], stdin)
	if err != nil {
		return nil, nil, fmt.Errorf("reading second input: %w", err)
	}
	return a, b, nil
}

func readArg(arg string, stdin io.Reader) ([]byte, error) {
	if arg == "-" {
		return ioutil.ReadStdin(stdin)
	}
	return os.ReadFile(arg)
}

// DiffText returns a unified diff of two text inputs.
func DiffText(a, b []byte, contextLines int) string {
	return textLib.Diff(string(a), string(b), contextLines)
}

// DiffJSON normalizes both inputs (sorted keys, consistent indent) then diffs.
func DiffJSON(a, b []byte, contextLines int) (string, error) {
	normA, err := normalizeJSON(a)
	if err != nil {
		return "", fmt.Errorf("first input: %w", err)
	}
	normB, err := normalizeJSON(b)
	if err != nil {
		return "", fmt.Errorf("second input: %w", err)
	}
	return textLib.Diff(normA, normB, contextLines), nil
}

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
	colorBold  = "\033[1m"
)

// Colorize adds ANSI colors to unified diff output.
func Colorize(diff string) string {
	if diff == "" {
		return ""
	}
	lines := strings.Split(diff, "\n")
	// strings.Split on a trailing \n produces an empty final element; drop it.
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	var sb strings.Builder
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++"):
			sb.WriteString(colorBold + line + colorReset)
		case strings.HasPrefix(line, "@@"):
			sb.WriteString(colorCyan + line + colorReset)
		case strings.HasPrefix(line, "-"):
			sb.WriteString(colorRed + line + colorReset)
		case strings.HasPrefix(line, "+"):
			sb.WriteString(colorGreen + line + colorReset)
		default:
			sb.WriteString(line)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func normalizeJSON(data []byte) (string, error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out) + "\n", nil
}
