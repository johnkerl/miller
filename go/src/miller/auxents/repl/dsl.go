// ================================================================
// This is the handler for taking DSL statements typed in interactively by the
// user, parsing them to an AST, building a CST from the AST, and executing the
// CST. It also handles DSL statements invoked using ':load' or multiline '<'
// ... '>' wherein statements are built into the AST without being executed
// right away.
//
// Specifically:
//
// * begin/end statements are parsed and stored into the AST regardless.
//
// * func/subr definitions are parsed and stored into the AST regardless.
//
// * For anything else in an interactive command besides begin/end/func/subr
//   blocks, the statement(s) is/are executed immediately for interactive mode,
//   else populated into the CST's main-statements block for load-from-file
//   multiline mode.
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

func (this *Repl) handleDSLStringAux(
	dslString string,
	isReplImmediate bool, // False for load-from-file or multiline; true otherwise
) error {
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
