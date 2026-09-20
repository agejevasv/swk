package ioutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
// or literal content (return as-is). An argument that looks like a path but is
// not readable is reported rather than hashed or parsed as literal text.
func resolveFileArg(arg string, stdin io.Reader) ([]byte, error) {
	if arg == "-" {
		return ReadStdin(stdin)
	}

	info, err := os.Stat(arg)
	switch {
	case err == nil && info.Mode().IsRegular():
		// No size cap here: the caller named this file deliberately, and
		// hashing or converting a large one is a normal thing to ask for.
		// The cap on stdin exists because a pipe's size cannot be known.
		return os.ReadFile(arg)
	case err == nil && info.IsDir() && looksLikePath(arg):
		return nil, fmt.Errorf("%s is a directory", arg)
	case err != nil && looksLikePath(arg):
		return nil, err
	}

	return []byte(arg), nil
}

// looksLikePath reports whether an argument was meant as a file rather than as
// literal text, so a typo'd filename is reported instead of being parsed as
// data. Documents carry separators too, so anything that opens like structured
// content, or spans lines, is treated as literal.
func looksLikePath(arg string) bool {
	if arg == "" || strings.ContainsAny(arg, "\n\r") {
		return false
	}
	switch arg[0] {
	case '{', '[', '<', '"', '-':
		return false
	}
	if strings.ContainsRune(arg, '/') || strings.ContainsRune(arg, os.PathSeparator) {
		return true
	}

	// A bare name ending in a data-file extension is a path as well, so a
	// mistyped "notes.md" is reported rather than hashed as literal text.
	return !strings.ContainsAny(arg, " \t") && dataFileExtensions[strings.ToLower(filepath.Ext(arg))]
}

// dataFileExtensions lists suffixes that mean "this argument names a file".
// The list is deliberately narrow: anything not on it stays literal text.
var dataFileExtensions = map[string]bool{
	".json": true, ".yaml": true, ".yml": true, ".xml": true, ".csv": true, ".tsv": true,
	".md": true, ".markdown": true, ".txt": true, ".log": true, ".html": true, ".htm": true,
	".pem": true, ".crt": true, ".cer": true, ".key": true, ".jwt": true, ".toml": true,
	".ini": true, ".conf": true, ".cfg": true, ".bin": true, ".png": true, ".jpg": true,
	".jpeg": true, ".gif": true, ".iso": true, ".tar": true, ".gz": true, ".zip": true,
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
