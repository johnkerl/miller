package lib

import "fmt"

// ExitRequest is a sentinel error requesting process termination with the
// given exit code. Code paths that have already produced all the output they
// want (e.g. --version, or 'mlr put --explain' on a valid expression) return
// an ExitRequest instead of calling os.Exit mid-stack; the entrypoint unwraps
// it with errors.As and exits with the requested code. This is control flow
// in the io.EOF style: an "error" that isn't a failure. See plans/exit.md.
type ExitRequest struct {
	Code int
}

func (e *ExitRequest) Error() string {
	return fmt.Sprintf("exit status %d", e.Code)
}
