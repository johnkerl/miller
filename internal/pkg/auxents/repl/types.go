// ================================================================
// Data types for the Miller REPL.
// ================================================================

package repl

import (
	"os"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/dsl/cst"
	"mlr/internal/pkg/input"
	"mlr/internal/pkg/output"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
)

// ================================================================
type ASTPrintMode int

const (
	ASTPrintNone ASTPrintMode = iota
	ASTPrintParex
	ASTPrintParexOneLine
	ASTPrintIndent
)

// ================================================================
type Repl struct {
	// From os.Args[] as we were invoked. These are for printing error messages.
	exeName  string
	replName string

	// Prompt1 is the main prompt, like $PS1. Prompt2 is for
	// multi-line-input mode with "<" ... ">" or "<<" ... ">>".
	inputIsTerminal   bool
	showStartupBanner bool
	showPrompts       bool
	prompt1           string
	prompt2           string

	astPrintMode ASTPrintMode
	doWarnings   bool
	cstRootNode  *cst.RootNode

	options *cli.TOptions

	inputChannel          chan *types.RecordAndContext
	errorChannel          chan error
	downstreamDoneChannel chan bool
	recordReader          input.IRecordReader
	recordWriter          output.IRecordWriter
	outputStream          *os.File

	runtimeState *runtime.State

	// For control-C handling
	sysToSignalHandlerChannel    chan os.Signal // Our signal handler reads system notification here
	appSignalNotificationChannel chan bool      // Our signal handler writes this for our app to poll
}
