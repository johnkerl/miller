// ================================================================
// This handles dump and edump statements.
// See print.go for comments; this is similar.
//
// Differences between print and dump:
//
// * 'print $x' and 'dump $x' are the same.
//
// * 'print' and 'dump' with no specific value: print outputs a newline; dump
//   outputs a JSON representation of all out-of-stream variables.
//
// * 'print $x,$y,$z' prints all items on one line; 'dump $x,$y,$z' prints each on
//   its own line.
// ================================================================

package cst

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/output"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
)

// ================================================================
type tDumpToRedirectFunc func(
	outputString string,
	state *runtime.State,
) error

type DumpStatementNode struct {
	expressionEvaluables      []IEvaluable
	dumpToRedirectFunc        tDumpToRedirectFunc
	redirectorTargetEvaluable IEvaluable                  // for file/pipe targets
	outputHandlerManager      output.OutputHandlerManager // for file/pipe targets
}

// ----------------------------------------------------------------
func (root *RootNode) BuildDumpStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDumpStatement)
	return root.buildDumpxStatementNode(
		astNode,
		os.Stdout,
	)
}

func (root *RootNode) BuildEdumpStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEdumpStatement)
	return root.buildDumpxStatementNode(
		astNode,
		os.Stderr,
	)
}

// ----------------------------------------------------------------
// Common code for building dump/edump nodes

func (root *RootNode) buildDumpxStatementNode(
	astNode *dsl.ASTNode,
	defaultOutputStream *os.File,
) (IExecutable, error) {
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)
	expressionsNode := astNode.Children[0]
	redirectorNode := astNode.Children[1]

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Things to be dumped, e.g. $a and $b in 'dump > "foo.dat", $a, $b'.

	var expressionEvaluables []IEvaluable = nil

	if expressionsNode.Type == dsl.NodeTypeNoOp {
		// Just 'dump' without 'dump $something'
		expressionEvaluables = make([]IEvaluable, 1)
		expressionEvaluable := root.BuildFullOosvarRvalueNode()
		expressionEvaluables[0] = expressionEvaluable
	} else if expressionsNode.Type == dsl.NodeTypeFunctionCallsite {
		expressionEvaluables = make([]IEvaluable, len(expressionsNode.Children))
		for i, childNode := range expressionsNode.Children {
			expressionEvaluable, err := root.BuildEvaluableNode(childNode)
			if err != nil {
				return nil, err
			}
			expressionEvaluables[i] = expressionEvaluable
		}
	} else {
		lib.InternalCodingErrorIf(true)
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Redirection targets (the thing after > >> |, if any).

	retval := &DumpStatementNode{
		expressionEvaluables:      expressionEvaluables,
		dumpToRedirectFunc:        nil,
		redirectorTargetEvaluable: nil,
		outputHandlerManager:      nil,
	}

	if redirectorNode.Type == dsl.NodeTypeNoOp {
		// No > >> or | was provided.
		if defaultOutputStream == os.Stdout {
			retval.dumpToRedirectFunc = retval.dumpToStdout
		} else if defaultOutputStream == os.Stderr {
			retval.dumpToRedirectFunc = retval.dumpToStderr
		} else {
			lib.InternalCodingErrorIf(true)
		}
	} else {
		// There is > >> or | provided.
		lib.InternalCodingErrorIf(redirectorNode.Children == nil)
		lib.InternalCodingErrorIf(len(redirectorNode.Children) != 1)
		redirectorTargetNode := redirectorNode.Children[0]
		var err error = nil

		if redirectorTargetNode.Type == dsl.NodeTypeRedirectTargetStdout {
			retval.dumpToRedirectFunc = retval.dumpToStdout
		} else if redirectorTargetNode.Type == dsl.NodeTypeRedirectTargetStderr {
			retval.dumpToRedirectFunc = retval.dumpToStderr
		} else {
			retval.dumpToRedirectFunc = retval.dumpToFileOrPipe

			retval.redirectorTargetEvaluable, err = root.BuildEvaluableNode(redirectorTargetNode)
			if err != nil {
				return nil, err
			}

			if redirectorNode.Type == dsl.NodeTypeRedirectWrite {
				retval.outputHandlerManager = output.NewFileWritetHandlerManager(root.recordWriterOptions)
			} else if redirectorNode.Type == dsl.NodeTypeRedirectAppend {
				retval.outputHandlerManager = output.NewFileAppendHandlerManager(root.recordWriterOptions)
			} else if redirectorNode.Type == dsl.NodeTypeRedirectPipe {
				retval.outputHandlerManager = output.NewPipeWriteHandlerManager(root.recordWriterOptions)
			} else {
				return nil, errors.New(
					fmt.Sprintf(
						"%s: unhandled redirector node type %s.",
						"mlr", string(redirectorNode.Type),
					),
				)
			}
		}
	}

	// Register this with the CST root node so that open file descriptrs can be
	// closed, etc at end of stream.
	if retval.outputHandlerManager != nil {
		root.RegisterOutputHandlerManager(retval.outputHandlerManager)
	}

	return retval, nil
}

// ----------------------------------------------------------------
func (node *DumpStatementNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	// 5x faster than fmt.Dump() separately: note that os.Stdout is
	// non-buffered in Go whereas stdout is buffered in C.
	//
	// Minus: we need to do our own buffering for performance.
	//
	// Plus: we never have to worry about forgetting to do fflush(). :)
	var buffer bytes.Buffer

	for _, expressionEvaluable := range node.expressionEvaluables {
		evaluation := expressionEvaluable.Evaluate(state)
		if !evaluation.IsAbsent() {
			s := evaluation.String()
			buffer.WriteString(s)
			if !strings.HasSuffix(s, "\n") {
				buffer.WriteString("\n")
			}
		}
	}
	outputString := buffer.String()
	node.dumpToRedirectFunc(outputString, state)
	return nil, nil
}

// ----------------------------------------------------------------
func (node *DumpStatementNode) dumpToStdout(
	outputString string,
	state *runtime.State,
) error {
	// Insert the string into the record-output stream, so that goroutine can
	// print it, resulting in deterministic output-ordering.
	//
	// The output channel is always non-nil, except for the Miller REPL.
	if state.OutputChannel != nil {
		state.OutputChannel <- types.NewOutputString(outputString, state.Context)
	} else {
		fmt.Println(outputString)
	}

	return nil
}

// ----------------------------------------------------------------
func (node *DumpStatementNode) dumpToStderr(
	outputString string,
	state *runtime.State,
) error {
	fmt.Fprintf(os.Stderr, outputString)
	return nil
}

// ----------------------------------------------------------------
func (node *DumpStatementNode) dumpToFileOrPipe(
	outputString string,
	state *runtime.State,
) error {
	redirectorTarget := node.redirectorTargetEvaluable.Evaluate(state)
	if !redirectorTarget.IsString() {
		return errors.New(
			fmt.Sprintf(
				"%s: output redirection yielded %s, not string.",
				"mlr", redirectorTarget.GetTypeName(),
			),
		)
	}
	outputFileName := redirectorTarget.String()

	node.outputHandlerManager.WriteString(outputString, outputFileName)
	return nil
}
