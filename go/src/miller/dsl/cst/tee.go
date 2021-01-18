// ================================================================
// This handles tee statements.
// ================================================================

package cst

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"miller/dsl"
	"miller/lib"
	"miller/types"
)

// ----------------------------------------------------------------
// Examples:
//   tee > "foo.dat", $*
//   tee > stderr, $*
//   tee > stdout, $*
//   tee | "jq .", $*
//
// The item being teed can only be $*. This is a special case of emit.  (This
// doesn't do anything emit can't do.)
//
// $ mlr -n put -v 'tee > stdout, $*'
// DSL EXPRESSION:
// tee > stdout, $*
// RAW AST:
// * statement block
//     * tee statement "tee"
//         * full record "$*"
//         * redirect write ">"
//             * stdout redirect target "stdout"
//
// $ mlr -n put -v 'tee > "foo.dat", $*'
// DSL EXPRESSION:
// tee > "foo.dat", $*
// RAW AST:
// * statement block
//     * tee statement "tee"
//         * full record "$*"
//         * redirect write ">"
//             * string literal "foo.dat"
//
// $ mlr -n put -v 'tee | "jq .", $*'
// DSL EXPRESSION:
// tee | "jq .", $*
// RAW AST:
// * statement block
//     * tee statement "tee"
//         * full record "$*"
//         * redirect pipe "|"
//             * string literal "jq ."
// ----------------------------------------------------------------

// ================================================================
type tTeeToRedirectFunc func(
	outputString string,
	state *State,
) error

type TeeStatementNode struct {
	expressionEvaluable       IEvaluable // always the $* evaluable
	teeToRedirectFunc         tTeeToRedirectFunc
	redirectorTargetEvaluable IEvaluable           // for file/pipe targets
	outputHandlerManager      OutputHandlerManager // for file/pipe targets
}

// ----------------------------------------------------------------
func (this *RootNode) BuildTeeStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeTeeStatement)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)
	expressionNode := astNode.Children[0]
	redirectorNode := astNode.Children[1]

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Expresosin to be teed, which is $*.

	lib.InternalCodingErrorIf(expressionNode.Type != dsl.NodeTypeFullSrec)
	expressionEvaluable := this.BuildFullSrecRvalueNode()

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Redirection targets (the thing after > >> |, if any).

	retval := &TeeStatementNode{
		expressionEvaluable:       expressionEvaluable,
		teeToRedirectFunc:         nil,
		redirectorTargetEvaluable: nil,
		outputHandlerManager:      nil,
	}

	// There is > >> or | provided.
	lib.InternalCodingErrorIf(redirectorNode.Children == nil)
	lib.InternalCodingErrorIf(len(redirectorNode.Children) != 1)
	redirectorTargetNode := redirectorNode.Children[0]
	var err error = nil

	if redirectorTargetNode.Type == dsl.NodeTypeRedirectTargetStdout {
		retval.teeToRedirectFunc = retval.teeToStdout
	} else if redirectorTargetNode.Type == dsl.NodeTypeRedirectTargetStderr {
		retval.teeToRedirectFunc = retval.teeToStderr
	} else {
		retval.teeToRedirectFunc = retval.teeToFileOrPipe

		retval.redirectorTargetEvaluable, err = this.BuildEvaluableNode(redirectorTargetNode)
		if err != nil {
			return nil, err
		}

		if redirectorNode.Type == dsl.NodeTypeRedirectWrite {
			retval.outputHandlerManager = NewFileWritetHandlerManager()
		} else if redirectorNode.Type == dsl.NodeTypeRedirectAppend {
			retval.outputHandlerManager = NewFileAppendHandlerManager()
		} else if redirectorNode.Type == dsl.NodeTypeRedirectPipe {
			retval.outputHandlerManager = NewPipeWriteHandlerManager()
		} else {
			return nil, errors.New(
				fmt.Sprintf(
					"%s: unhandled redirector node type %s.",
					os.Args[0], string(redirectorNode.Type),
				),
			)
		}
	}

	// TODO: root node register outputHandlerManager to add to close-handles list

	return retval, nil
}

// ----------------------------------------------------------------
func (this *TeeStatementNode) Execute(state *State) (*BlockExitPayload, error) {
	evaluation := this.expressionEvaluable.Evaluate(state)
	outputString := evaluation.String()
	if !strings.HasSuffix(outputString, "\n") {
		outputString += "\n"
	}
	this.teeToRedirectFunc(outputString, state)
	return nil, nil
}

// ----------------------------------------------------------------
func (this *TeeStatementNode) teeToStdout(
	outputString string,
	state *State,
) error {
	// Insert the string into the record-output stream, so that goroutine can
	// print it, resulting in deterministic output-ordering.
	state.OutputChannel <- types.NewOutputString(outputString, state.Context)
	return nil
}

// ----------------------------------------------------------------
func (this *TeeStatementNode) teeToStderr(
	outputString string,
	state *State,
) error {
	fmt.Fprintf(os.Stderr, outputString)
	return nil
}

// ----------------------------------------------------------------
func (this *TeeStatementNode) teeToFileOrPipe(
	outputString string,
	state *State,
) error {
	redirectorTarget := this.redirectorTargetEvaluable.Evaluate(state)
	if !redirectorTarget.IsString() {
		return errors.New(
			fmt.Sprintf(
				"%s: output redirection yielded %s, not string.",
				os.Args[0], redirectorTarget.GetTypeName(),
			),
		)
	}
	outputFileName := redirectorTarget.String()

	this.outputHandlerManager.Print(outputString, outputFileName)
	return nil
}
