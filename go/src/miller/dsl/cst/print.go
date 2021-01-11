// ================================================================ This
// handles print and dump statements.
// ================================================================

package cst

import (
	"errors"
	"fmt"
	"os"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
type PrintStatementNode struct {
	outputHandlerManager OutputHandlerManager // TODO: comments
	terminator           string
	expressions          []IEvaluable
	redirectorTarget     IEvaluable
}

// ----------------------------------------------------------------
func (this *RootNode) BuildPrintStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypePrintStatement)
	return this.buildPrintxStatementNode(
		astNode,
		os.Stdout,
		"\n",
	)
}

func (this *RootNode) BuildPrintnStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypePrintnStatement)
	return this.buildPrintxStatementNode(
		astNode,
		os.Stdout,
		"",
	)
}

func (this *RootNode) BuildEprintStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEprintStatement)
	return this.buildPrintxStatementNode(
		astNode,
		os.Stderr,
		"\n",
	)
}

func (this *RootNode) BuildEprintnStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEprintnStatement)
	return this.buildPrintxStatementNode(
		astNode,
		os.Stderr,
		"",
	)
}

// ----------------------------------------------------------------
// Common code for building print/eprint/printn/eprintn nodes
//
// Example ASTs:
//
// $ mlr -n put -v 'print 1, 2'
// DSL EXPRESSION:
// print 1, 2
// RAW AST:
// * statement block
//     * print statement "print"
//         * function callsite
//             * int literal "1"
//             * int literal "2"
//         * no-op
//
// $ mlr -n put -v 'print > "foo", 1, 2'
// DSL EXPRESSION:
// print > "foo", 1, 2
// RAW AST:
// * statement block
//     * print statement "print"
//         * function callsite
//             * int literal "1"
//             * int literal "2"
//         * redirect write ">"
//             * string literal "foo"
//
// $ mlr -n put -v 'print >> "foo", 1, 2'
// DSL EXPRESSION:
// print >> "foo", 1, 2
// RAW AST:
// * statement block
//     * print statement "print"
//         * function callsite
//             * int literal "1"
//             * int literal "2"
//         * redirect append ">>"
//             * string literal "foo"

func (this *RootNode) buildPrintxStatementNode(
	astNode *dsl.ASTNode,
	defaultOutputStream *os.File,
	terminator string,
) (IExecutable, error) {
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	expressionsNode := astNode.Children[0]
	redirectNode := astNode.Children[1]

	expressions := make([]IEvaluable, len(expressionsNode.Children))
	for i, childNode := range expressionsNode.Children {
		expression, err := this.BuildEvaluableNode(childNode)
		if err != nil {
			return nil, err
		}
		expressions[i] = expression
	}

	// Without explicit redirect, the redirect AST node comes in as a no-op
	// node from the parser.
	var outputHandlerManager OutputHandlerManager = nil
	if redirectNode.Type == dsl.NodeTypeNoOp {
		if defaultOutputStream == os.Stdout {
			outputHandlerManager = NewStdoutOutputHandlerManager()
		} else if defaultOutputStream == os.Stderr {
			outputHandlerManager = NewStderrOutputHandlerManager()
		} else {
			lib.InternalCodingErrorIf(true)
		}
	} else if redirectNode.Type == dsl.NodeTypeRedirectWrite {
		outputHandlerManager = NewFileWritetHandlerManager()
	} else if redirectNode.Type == dsl.NodeTypeRedirectAppend {
		outputHandlerManager = NewFileAppendHandlerManager()
	} else if redirectNode.Type == dsl.NodeTypeRedirectPipe {
		outputHandlerManager = NewPipeWriteHandlerManager()
	} else {
		return nil, errors.New(
			fmt.Sprintf(
				"%s: unhandled redirection node type %s.",
				os.Args[0], string(redirectNode.Type),
			),
		)
	}

	var redirectorTarget IEvaluable = nil
	if redirectNode.Type != dsl.NodeTypeNoOp {
		lib.InternalCodingErrorIf(redirectNode.Children == nil)
		lib.InternalCodingErrorIf(len(redirectNode.Children) != 1)
		redirectorTargetNode := redirectNode.Children[0]
		var err error = nil
		redirectorTarget, err = this.BuildEvaluableNode(redirectorTargetNode)
		if err != nil {
			return nil, err
		}
	}

	// TODO: root node register oututHandlerManager to add to close-handles list

	return &PrintStatementNode{
		outputHandlerManager: outputHandlerManager,
		terminator:           terminator,
		expressions:          expressions,
		redirectorTarget:     redirectorTarget,
	}, nil
}

// ----------------------------------------------------------------
// TODO: maybe get rid of the stdout/stderr managers ... just a
// function-pointer here would suffice.
func (this *PrintStatementNode) Execute(state *State) (*BlockExitPayload, error) {
	outputFileName := "_miller_bug_if_you_see_this_"
	if this.redirectorTarget != nil {
		redirectorEvaluation := this.redirectorTarget.Evaluate(state)
		if !redirectorEvaluation.IsString() {
			return nil, errors.New(
				fmt.Sprintf(
					"%s: output redirection yielded %s, not string.",
					os.Args[0], redirectorEvaluation.GetTypeName(),
				),
			)
		}
		outputFileName = redirectorEvaluation.String()
	}

	if len(this.expressions) == 0 {
		this.outputHandlerManager.Print(this.terminator, outputFileName)
	} else {
		// TODO: buffered bytes like record-writers do
		for i, expression := range this.expressions {
			if i > 0 {
				this.outputHandlerManager.Print(" ", outputFileName)
			}
			evaluation := expression.Evaluate(state)
			if !evaluation.IsAbsent() {
				this.outputHandlerManager.Print(evaluation.String(), outputFileName)
			}
		}
		this.outputHandlerManager.Print(this.terminator, outputFileName)
	}
	return nil, nil
}
