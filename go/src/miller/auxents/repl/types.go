// ================================================================
// Data types for the Miller REPL.
// ================================================================

package repl

import (
	"os"

	"miller/cliutil"
	"miller/dsl/cst"
	"miller/input"
	"miller/output"
	"miller/runtime"
	"miller/types"
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

	// Prompt1 is the main prompt, like $PS1. Prompt2 is for "<" ... ">"
	// multi-line-input mode.
	inputIsTerminal     bool
	prompt1             string
	prompt2             string
	doingMultilineInput bool

	astPrintMode ASTPrintMode
	cstRootNode  *cst.RootNode

	options *cliutil.TOptions

	inputChannel chan *types.RecordAndContext
	errorChannel chan error
	recordReader input.IRecordReader
	recordWriter output.IRecordWriter
	outputStream *os.File

	runtimeState *runtime.State

	// For control-C handling
	sysToSignalHandlerChannel    chan os.Signal // Our signal handler reads system notification here
	appSignalNotificationChannel chan bool      // Our signal handler writes this for our app to poll
}
