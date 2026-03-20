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
