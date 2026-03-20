package ioutil

// ExitCoder allows commands to signal a specific exit code.
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
