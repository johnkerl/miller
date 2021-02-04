// ================================================================
// Just playing around -- nothing serious here.
// ================================================================

package repl

import (
	"fmt"
	"strings"

	"miller/types"
)

// ----------------------------------------------------------------
func (this *Repl) handleDSLStringImmediate(dslString string) error {
	return this.handleDSLStringAux(dslString, true)
}
func (this *Repl) handleDSLStringBulk(dslString string) error {
	return this.handleDSLStringAux(dslString, false)
}

func (this *Repl) handleDSLStringAux(dslString string, isReplImmediate bool) error {
	if strings.TrimSpace(dslString) == "" {
		return nil
	}

	astRootNode, err := this.BuildASTFromStringWithMessage(dslString)
	if err != nil {
		// Error message already printed out
		return err
	}

	this.cstRootNode.ResetForREPL()

	// For load-from-file / multi-line, each statement not in begin/end, and
	// not a user-defined function/subroutine, is recorded in the "main block"
	// to be executed later when the user asks to do so. For
	// single-line/interactive mode, begin/end statements and UDF/UDS are
	// recorded, but any other statements are executed immediately.
	err = this.cstRootNode.IngestAST(
		astRootNode,
		false, /*isFilter*/
		isReplImmediate,
	)
	if err != nil {
		return err
	}

	err = this.cstRootNode.Resolve()
	if err != nil {
		return err
	}

	if isReplImmediate {
		outrec, err := this.cstRootNode.ExecuteREPLImmediate(this.runtimeState)
		if err != nil {
			return err
		}
		this.runtimeState.Inrec = outrec

		// The filter expression for the main Miller DSL is any non-assignment
		// statment like 'true' or '$x > 0.5' etc. For the REPL, we re-use this for
		// interactive expressions to be printed to the terminal. For the main DSL,
		// the default is types.MlrvalFromTrue(); for the REPL, the default is
		// types.MlrvalFromVoid().
		filterExpression := this.runtimeState.FilterExpression
		if filterExpression.IsVoid() {
			// nothing to print
		} else {
			fmt.Println(filterExpression.String())
		}
		this.runtimeState.FilterExpression = types.MlrvalFromVoid()
	}

	return nil
}
