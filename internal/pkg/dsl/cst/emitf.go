// ================================================================
// This handles emitf statements. This produces new records (in addition to $*)
// into the output record stream.
// ================================================================

package cst

import (
	"errors"
	"fmt"

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/output"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
)

// ================================================================
// Examples:
//   emitf @a
//   emitf @a, @b
//   emitf > "foo.dat", @a, @b
//
// Each argument must be a non-indexed oosvar/localvar/fieldname, so we can use
// their names as keys in the emitted record.  These restrictions are enforced
// in the CST logic, to keep this parser/AST logic simpler.

type tEmitFToRedirectFunc func(
	newrec *types.Mlrmap,
	state *runtime.State,
) error

type EmitFStatementNode struct {
	emitfNames                []string
	emitfEvaluables           []IEvaluable
	emitfToRedirectFunc       tEmitFToRedirectFunc
	redirectorTargetEvaluable IEvaluable                  // for file/pipe targets
	outputHandlerManager      output.OutputHandlerManager // for file/pipe targets
}

// ----------------------------------------------------------------
// $ mlr -n put -v 'emitf a,$b,@c'
// DSL EXPRESSION:
// emitf a,$b,@c
// AST:
// * statement block
//     * emitf statement "emitf"
//         * emittable list
//             * local variable "a"
//             * direct field value "b"
//             * direct oosvar value "c"
//         * no-op

func (root *RootNode) BuildEmitFStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEmitFStatement)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)
	expressionsNode := astNode.Children[0]
	redirectorNode := astNode.Children[1]

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Things to be emitted, e.g. @a and @b in 'emitf > "foo.dat", @a, @b'.

	n := len(expressionsNode.Children)
	lib.InternalCodingErrorIf(n < 1)
	emitfNames := make([]string, n)
	emitfEvaluables := make([]IEvaluable, n)
	for i, childNode := range expressionsNode.Children {
		emitfName, err := getNameFromNamedNode(childNode, "emitf")
		if err != nil {
			return nil, err
		}
		emitfEvaluable, err := root.BuildEvaluableNode(childNode)
		if err != nil {
			return nil, err
		}
		emitfNames[i] = emitfName
		emitfEvaluables[i] = emitfEvaluable
	}

	retval := &EmitFStatementNode{
		emitfNames:      emitfNames,
		emitfEvaluables: emitfEvaluables,
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Redirection targets (the thing after > >> |, if any).

	if redirectorNode.Type == dsl.NodeTypeNoOp {
		// No > >> or | was provided.
		retval.emitfToRedirectFunc = retval.emitfToRecordStream
	} else {
		// There is > >> or | provided.
		lib.InternalCodingErrorIf(redirectorNode.Children == nil)
		lib.InternalCodingErrorIf(len(redirectorNode.Children) != 1)
		redirectorTargetNode := redirectorNode.Children[0]
		var err error = nil

		if redirectorTargetNode.Type == dsl.NodeTypeRedirectTargetStdout {
			retval.emitfToRedirectFunc = retval.emitfToFileOrPipe
			retval.outputHandlerManager = output.NewStdoutWriteHandlerManager(root.recordWriterOptions)
			retval.redirectorTargetEvaluable = root.BuildStringLiteralNode("(stdout)")
		} else if redirectorTargetNode.Type == dsl.NodeTypeRedirectTargetStderr {
			retval.emitfToRedirectFunc = retval.emitfToFileOrPipe
			retval.outputHandlerManager = output.NewStderrWriteHandlerManager(root.recordWriterOptions)
			retval.redirectorTargetEvaluable = root.BuildStringLiteralNode("(stderr)")
		} else {
			retval.emitfToRedirectFunc = retval.emitfToFileOrPipe

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

func (node *EmitFStatementNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	newrec := types.NewMlrmapAsRecord()
	for i, emitfEvaluable := range node.emitfEvaluables {
		emitfName := node.emitfNames[i]
		emitfValue := emitfEvaluable.Evaluate(state)

		if !emitfValue.IsAbsent() {
			newrec.PutCopy(emitfName, emitfValue)
		}
	}

	err := node.emitfToRedirectFunc(newrec, state)

	return nil, err
}

// ----------------------------------------------------------------
// Gets the name of a non-indexed oosvar, localvar, or field name; otherwise,
// returns error.
//
// TODO: support indirects like 'emitf @[x."_sum"]'

func getNameFromNamedNode(astNode *dsl.ASTNode, description string) (string, error) {
	if astNode.Type == dsl.NodeTypeDirectOosvarValue {
		return string(astNode.Token.Lit), nil
	} else if astNode.Type == dsl.NodeTypeLocalVariable {
		return string(astNode.Token.Lit), nil
	} else if astNode.Type == dsl.NodeTypeDirectFieldValue {
		return string(astNode.Token.Lit), nil
	}
	return "", errors.New(
		fmt.Sprintf(
			"%s: can't get name of node type \"%s\" for %s.",
			"mlr", string(astNode.Type), description,
		),
	)
}

// ----------------------------------------------------------------
func (node *EmitFStatementNode) emitfToRecordStream(
	outrec *types.Mlrmap,
	state *runtime.State,
) error {
	// The output channel is always non-nil, except for the Miller REPL.
	if state.OutputChannel != nil {
		state.OutputChannel <- types.NewRecordAndContext(outrec, state.Context)
	} else {
		fmt.Println(outrec.String())
	}
	return nil
}

// ----------------------------------------------------------------
func (node *EmitFStatementNode) emitfToFileOrPipe(
	outrec *types.Mlrmap,
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

	return node.outputHandlerManager.WriteRecordAndContext(
		types.NewRecordAndContext(outrec, state.Context),
		outputFileName,
	)
}
