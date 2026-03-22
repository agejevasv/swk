package ioutil

// ExitCoder allows commands to signal a specific exit code.
// Exit code 1 = negative result (no match, check failed), not a bug.
// Exit code 2 = actual error (bad input, invalid flags, etc.).
//
// Convention: Error() should return "" for exit code 1 types. The root
// command prints "Error: <msg>" to stderr for non-empty messages, so
// returning "" suppresses that output for expected negative results.
type ExitCoder interface {
	error
	ExitCode() int
}

// NoMatchError signals that a query found no results (exit code 1, like grep).
type NoMatchError struct{}

func (NoMatchError) Error() string { return "" }
func (NoMatchError) ExitCode() int { return 1 }

// CheckFailedError signals that a verification check failed (exit code 1).
// Output is already on stdout; no stderr diagnostic needed.
type CheckFailedError struct{}

func (CheckFailedError) Error() string { return "" }
func (CheckFailedError) ExitCode() int { return 1 }
