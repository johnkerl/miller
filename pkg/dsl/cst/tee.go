// This handles tee statements.

package cst

import (
	"fmt"

	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/output"
	"github.com/johnkerl/miller/v6/pkg/runtime"
	"github.com/johnkerl/miller/v6/pkg/types"
	"github.com/johnkerl/pgpg/go/lib/pkg/asts"
)

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
// AST:
// * statement block
//     * tee statement "tee"
//         * full record "$*"
//         * redirect write ">"
//             * stdout redirect target "stdout"
//
// $ mlr -n put -v 'tee > "foo.dat", $*'
// DSL EXPRESSION:
// tee > "foo.dat", $*
// AST:
// * statement block
//     * tee statement "tee"
//         * full record "$*"
//         * redirect write ">"
//             * string literal "foo.dat"
//
// $ mlr -n put -v 'tee | "jq .", $*'
// DSL EXPRESSION:
// tee | "jq .", $*
// AST:
// * statement block
//     * tee statement "tee"
//         * full record "$*"
//         * redirect pipe "|"
//             * string literal "jq ."

type tTeeToRedirectFunc func(
	outrec *mlrval.Mlrmap,
	state *runtime.State,
) error

type TeeStatementNode struct {
	expressionEvaluable       IEvaluable // always the $* evaluable
	teeToRedirectFunc         tTeeToRedirectFunc
	redirectorTargetEvaluable IEvaluable                  // for file/pipe targets
	outputHandlerManager      output.OutputHandlerManager // for file/pipe targets
}

func (root *RootNode) BuildTeeStatementNode(astNode *asts.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != asts.NodeType(NodeTypeTeeStatement))
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)
	// PGPG: kw_tee Redirector comma FullSrec -> children [1,3] = [Redirector, FullSrec]
	redirectorNode := astNode.Children[0]
	expressionNode := astNode.Children[1]

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Expression to be teed, which is $*.

	lib.InternalCodingErrorIf(expressionNode.Type != asts.NodeType(NodeTypeFullSrec))
	expressionEvaluable := root.BuildFullSrecRvalueNode()

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
	var err error

	if redirectorTargetNode.Type == asts.NodeType(NodeTypeRedirectTargetStdout) {
		retval.teeToRedirectFunc = retval.teeToFileOrPipe
		retval.outputHandlerManager = output.NewStdoutWriteHandlerManager(root.recordWriterOptions)
		retval.redirectorTargetEvaluable = root.BuildStringLiteralNode("(stdout)")
	} else if redirectorTargetNode.Type == asts.NodeType(NodeTypeRedirectTargetStderr) {
		retval.teeToRedirectFunc = retval.teeToFileOrPipe
		retval.outputHandlerManager = output.NewStderrWriteHandlerManager(root.recordWriterOptions)
		retval.redirectorTargetEvaluable = root.BuildStringLiteralNode("(stderr)")
	} else {
		retval.teeToRedirectFunc = retval.teeToFileOrPipe
		targetNode := redirectorTargetNode
		if redirectorTargetNode.Type == asts.NodeType(NodeTypeRedirectTargetRvalue) &&
			redirectorTargetNode.Children != nil && len(redirectorTargetNode.Children) > 0 {
			targetNode = redirectorTargetNode.Children[0]
		}
		retval.redirectorTargetEvaluable, err = root.BuildEvaluableNode(targetNode)
		if err != nil {
			return nil, err
		}

		if redirectorNode.Type == asts.NodeType(NodeTypeRedirectWrite) {
			retval.outputHandlerManager = output.NewFileWritetHandlerManager(root.recordWriterOptions)
		} else if redirectorNode.Type == asts.NodeType(NodeTypeRedirectAppend) {
			retval.outputHandlerManager = output.NewFileAppendHandlerManager(root.recordWriterOptions)
		} else if redirectorNode.Type == asts.NodeType(NodeTypeRedirectPipe) {
			retval.outputHandlerManager = output.NewPipeWriteHandlerManager(root.recordWriterOptions)
		} else {
			return nil, fmt.Errorf("unhandled redirector node type %s", string(redirectorNode.Type))
		}
	}

	// Register this with the CST root node so that open file descriptrs can be
	// closed, etc at end of stream.
	if retval.outputHandlerManager != nil {
		root.RegisterOutputHandlerManager(retval.outputHandlerManager)
	}

	return retval, nil
}

func (node *TeeStatementNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	expression := node.expressionEvaluable.Evaluate(state)
	if !expression.IsMap() {
		return nil, fmt.Errorf("tee-evaluation yielded %s, not map", expression.GetTypeName())
	}
	err := node.teeToRedirectFunc(expression.GetMap(), state)
	return nil, err
}

func (node *TeeStatementNode) teeToFileOrPipe(
	outrec *mlrval.Mlrmap,
	state *runtime.State,
) error {
	redirectorTarget := node.redirectorTargetEvaluable.Evaluate(state)
	if !redirectorTarget.IsString() {
		return fmt.Errorf("output redirection yielded %s, not string", redirectorTarget.GetTypeName())
	}
	outputFileName := redirectorTarget.String()

	return node.outputHandlerManager.WriteRecordAndContext(
		types.NewRecordAndContext(outrec, state.Context),
		outputFileName,
	)
}
