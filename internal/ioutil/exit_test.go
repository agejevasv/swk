package ioutil

import "testing"

func TestNoMatchError(t *testing.T) {
	var err error = NoMatchError{}
	if err.Error() != "" {
		t.Errorf("expected empty error string, got %q", err.Error())
	}
	ec, ok := err.(ExitCoder)
	if !ok {
		t.Fatal("NoMatchError does not implement ExitCoder")
	}
	if ec.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", ec.ExitCode())
	}
}

func TestCheckFailedError(t *testing.T) {
	var err error = CheckFailedError{}
	if err.Error() != "" {
		t.Errorf("expected empty error string, got %q", err.Error())
	}
	ec, ok := err.(ExitCoder)
	if !ok {
		t.Fatal("CheckFailedError does not implement ExitCoder")
	}
	if ec.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", ec.ExitCode())
	}
}
