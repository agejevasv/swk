package ioutil

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadInput returns raw bytes from args or stdin.
// Use for "value" commands where args[0] is literal data.
// Supports "-" as explicit stdin.
func ReadInput(args []string, stdin io.Reader) ([]byte, error) {
	if len(args) > 0 {
		if args[0] == "-" {
			return ReadStdin(stdin)
		}
		return []byte(args[0]), nil
	}
	return ReadStdin(stdin)
}

// ReadInputString is ReadInput with trailing newline trimming.
// Use for text "value" commands where args[0] is literal data.
func ReadInputString(args []string, stdin io.Reader) (string, error) {
	b, err := ReadInput(args, stdin)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(b), "\r\n"), nil
}

// ReadFileInput returns raw bytes from a file path arg, literal content, or stdin.
// Use for "document" commands where args[0] is typically a file path.
// Priority: no args → stdin, "-" → stdin, existing regular file → read it, otherwise → literal.
func ReadFileInput(args []string, stdin io.Reader) ([]byte, error) {
	if len(args) > 0 {
		return resolveFileArg(args[0], stdin)
	}
	return ReadStdin(stdin)
}

// ReadFileInputString is ReadFileInput with trailing newline trimming.
// Use for text "document" commands where args[0] is typically a file path.
func ReadFileInputString(args []string, stdin io.Reader) (string, error) {
	b, err := ReadFileInput(args, stdin)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(b), "\r\n"), nil
}

// resolveFileArg checks if arg is "-" (stdin), an existing regular file (read it),
// or literal content (return as-is).
func resolveFileArg(arg string, stdin io.Reader) ([]byte, error) {
	if arg == "-" {
		return ReadStdin(stdin)
	}
	info, err := os.Stat(arg)
	if err == nil && info.Mode().IsRegular() {
		return os.ReadFile(arg)
	}
	return []byte(arg), nil
}

// MaxInputSize is the maximum bytes ReadStdin will read (64 MiB).
// Prevents accidental OOM from unbounded pipes.
const MaxInputSize = 64 << 20

func ReadStdin(stdin io.Reader) ([]byte, error) {
	if stdin == nil {
		return nil, fmt.Errorf("no input provided")
	}
	limited := io.LimitReader(stdin, MaxInputSize+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if len(data) > MaxInputSize {
		return nil, fmt.Errorf("input exceeds maximum size (%d MiB)", MaxInputSize>>20)
	}
	return data, nil
}
