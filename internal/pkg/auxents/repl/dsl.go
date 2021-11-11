// ================================================================
// This is the handler for taking DSL statements typed in interactively by the
// user, parsing them to an AST, building a CST from the AST, and executing the
// CST. It also handles DSL statements invoked using ':load' or multi-line '<<'
// ... '>>' wherein statements are built into the AST without being executed
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
//   multi-line mode.
// ================================================================

package repl

import (
	"fmt"
	"strings"

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/dsl/cst"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
func (repl *Repl) handleDSLStringImmediate(dslString string, doWarnings bool) error {
	return repl.handleDSLStringAux(dslString, true, doWarnings)
}
func (repl *Repl) handleDSLStringBulk(dslString string, doWarnings bool) error {
	return repl.handleDSLStringAux(dslString, false, doWarnings)
}

func (repl *Repl) handleDSLStringAux(
	dslString string,
	isReplImmediate bool, // False for load-from-file or non-immediate multi-line; true otherwise
	doWarnings bool,
) error {
	if strings.TrimSpace(dslString) == "" {
		return nil
	}

	repl.cstRootNode.ResetForREPL()

	err := repl.cstRootNode.Build(
		[]string{dslString},
		cst.DSLInstanceTypeREPL,
		isReplImmediate,
		doWarnings,
		false, // warningsAreFatal
		func(dslString string, astNode *dsl.AST) {
			if repl.astPrintMode == ASTPrintParex {
				astNode.PrintParex()
			} else if repl.astPrintMode == ASTPrintParexOneLine {
				astNode.PrintParexOneLine()
			} else if repl.astPrintMode == ASTPrintIndent {
				astNode.Print()
			}
		},
	)
	if err != nil {
		return err
	}

	// For load-from-file / non-immediate multi-line, each statement not in
	// begin/end, and not a user-defined function/subroutine, is recorded in
	// the "main block" to be executed later when the user asks to do so. For
	// single-line/interactive mode, begin/end statements and UDF/UDS are
	// recorded, but any other statements are executed immediately.

	if isReplImmediate {
		outrec, err := repl.cstRootNode.ExecuteREPLImmediate(repl.runtimeState)
		if err != nil {
			return err
		}
		repl.runtimeState.Inrec = outrec

		// The filter expression for the main Miller DSL is any non-assignment
		// statment like 'true' or '$x > 0.5' etc. For the REPL, we re-use this for
		// interactive expressions to be printed to the terminal. For the main DSL,
		// the default is types.MlrvalFromTrue(); for the REPL, the default is
		// types.MLRVAL_VOID.
		filterExpression := repl.runtimeState.FilterExpression
		if filterExpression.IsVoid() {
			// nothing to print
		} else if filterExpression.IsString() {
			fmt.Printf("\"%s\"\n", filterExpression.String())
		} else {
			fmt.Println(filterExpression.String())
		}
		repl.runtimeState.FilterExpression = types.MLRVAL_VOID
	}

	return nil
}
