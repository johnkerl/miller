// ================================================================
// Just playing around -- nothing serious here.
// ================================================================

package repl

import (
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
	inputIsTerminal     bool
	prompt1             string
	prompt2             string
	doingMultilineInput bool

	astPrintMode ASTPrintMode
	isFilter     bool
	cstRootNode  *cst.RootNode

	options *cliutil.TOptions

	inputChannel chan *types.RecordAndContext
	errorChannel chan error
	recordReader input.IRecordReader
	recordWriter output.IRecordWriter

	runtimeState *runtime.State
}
