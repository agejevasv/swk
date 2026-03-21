package serve

import (
	"os"
	"path/filepath"
	"testing"
)

func TestServe_InvalidDir(t *testing.T) {
	Cmd.SetArgs([]string{"/nonexistent/path"})
	Cmd.SetOut(new(devNull))
	Cmd.SetErr(new(devNull))
	err := Cmd.Execute()
	if err == nil {
		t.Fatal("expected error for nonexistent directory")
	}
}

func TestServe_FileNotDir(t *testing.T) {
	f := filepath.Join(t.TempDir(), "file.txt")
	os.WriteFile(f, []byte("hello"), 0o644)

	Cmd.SetArgs([]string{f})
	Cmd.SetOut(new(devNull))
	Cmd.SetErr(new(devNull))
	err := Cmd.Execute()
	if err == nil {
		t.Fatal("expected error when path is a file, not a directory")
	}
}

type devNull struct{}

func (devNull) Write(p []byte) (int, error) { return len(p), nil }
