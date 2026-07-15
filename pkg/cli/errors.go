package cli

import (
	"errors"
	"fmt"
)

// ErrHelpRequested is returned by verb ParseCLIFunc when -h or --help is used.
// The caller (CLI layer) should exit with code 0 after the verb has printed its
// usage to stdout.
var ErrHelpRequested = errors.New("help requested")

// ErrUsagePrinted is returned by verb ParseCLIFunc when validation fails. The
// verb has already printed its usage to stderr. The caller should exit 1
// without printing the error again.
var ErrUsagePrinted = errors.New("usage printed")

// FlagErrorf is for user-facing main-flag errors (the counterpart of
// VerbErrorf for verb flags). The message is complete prose, ready for the
// entrypoint layer to print as-is: it should carry the "mlr: " prefix and may
// span multiple lines.
func FlagErrorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
