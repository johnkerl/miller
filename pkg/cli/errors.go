package cli

import "errors"

// ErrHelpRequested is returned by verb ParseCLIFunc when -h or --help is used.
// The caller (CLI layer) should exit with code 0 after the verb has printed its
// usage to stdout.
var ErrHelpRequested = errors.New("help requested")

// ErrUsagePrinted is returned by verb ParseCLIFunc when validation fails. The
// verb has already printed its usage to stderr. The caller should exit 1
// without printing the error again.
var ErrUsagePrinted = errors.New("usage printed")
