package ioutil

import (
	"fmt"
	"io"
	"strings"
)

// Use for binary commands (gzip, hash, image) where every byte matters.
func ReadInput(args []string, stdin io.Reader) ([]byte, error) {
	if len(args) > 0 {
		return []byte(args[0]), nil
	}
	return ReadStdin(stdin)
}

// Trims trailing newline that stdin pipes add. Use for text-processing commands.
func ReadInputString(args []string, stdin io.Reader) (string, error) {
	b, err := ReadInput(args, stdin)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(b), "\r\n"), nil
}

func ReadStdin(stdin io.Reader) ([]byte, error) {
	if stdin == nil {
		return nil, fmt.Errorf("no input provided")
	}
	return io.ReadAll(stdin)
}
