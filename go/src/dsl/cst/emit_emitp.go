// ================================================================
// This handles emit and emitp statements. These produce new records (in
// addition to $*) into the output record stream.
//
// Some complications here are due to legacy. Emit statements existed in the
// Miller DSL before there were for-loops. As a result, some of the
// side-by-side emit syntaxes were invented (and supported) to allow things
// that might have been more easily done with simpler emit syntax.
// Nonetheless, those syntaxes are now supported and we need to support them.
//
// Examples for emit and emitp:
//   emit @a
//   emit (@a, @b)
//   emit @a, "x", "y"
//   emit (@a, @b), "x", "y"
//
// The first argument (single or in parentheses) must be non-indexed
// oosvars/localvars/fieldnames, so we can use their names as keys in the
// emitted record, or they must be maps. So the first complexity in this code
// is, do we have a named variable or a map.
//
// The second complexity here is whether we have 'emit @a' or 'emit (@a, @b)'
// -- the latter being the "lashed" variant. Here, the keys of the first
// argument are used to drive indexing of the remaining arguments.
//
// The third complexlity here is whether we have the '"x", "y"' after the
// emittables. These control how nested maps are used to generate multiple
// records (via implicit looping).
// ================================================================

package cst

import (
	"errors"
	"fmt"

	"miller/src/dsl"
	"miller/src/lib"
	"miller/src/output"
	"miller/src/runtime"
	"miller/src/types"
)

// ================================================================
// Shared by emit and emitp

type tEmitToRedirectFunc func(
	newrec *types.Mlrmap,
	state *runtime.State,
) error

type EmitXStatementNode struct {
	// These are "_" for nameless maps like in 'emit {...}'; "x" for named
	// variables like in 'emit @x'.
	names []string

	// Maps or named variables: the @a, @b parts.
	emitEvaluables []IEvaluable

	// The "x","y" parts.
	indexEvaluables []IEvaluable

	// Whether each emittable is 'all'/'@*', or not
	isFullOosvars []bool

	// Appropriate function to evaluate statements, depending on indexed or not.
	executorFunc Executor

	// Appropriate function to send record(s) to stdout, stderr, write-to-file,
	// append-to-file, pipe-to-command, or insert into the record stream.
	emitToRedirectFunc tEmitToRedirectFunc
	// For file/pipe targets: 'emit > $a . ".dat", @x' -- the
	// redirectorTargetEvaluable is the evaluable for '$a . ".dat"'.
	redirectorTargetEvaluable IEvaluable
	// For file/pipe targets: keeps track of file handles for various values of
	// the redirectorTargetEvaluable expression.
	outputHandlerManager output.OutputHandlerManager

	// For code-reuse between executors.
	isEmitP bool
}

func (root *RootNode) BuildEmitStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEmitStatement)
	return root.buildEmitXStatementNode(astNode, false)
}
func (root *RootNode) BuildEmitPStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEmitPStatement)
	return root.buildEmitXStatementNode(astNode, true)
}

// ----------------------------------------------------------------
// EMIT AND EMITP
//
// Examples:
//   emit @a
//   emit (@a, @b)
//   emit @a, "x", "y"
//   emit (@a, @b), "x", "y"
// First argument (single or in parentheses) must be non-indexed
// oosvar/localvar/fieldname/map, so we can use their names as keys in the
// emitted record.

func (root *RootNode) buildEmitXStatementNode(
	astNode *dsl.ASTNode,
	isEmitP bool,
) (IExecutable, error) {
	lib.InternalCodingErrorIf(len(astNode.Children) != 3)
	emittablesNode := astNode.Children[0]
	keysNode := astNode.Children[1]
	redirectorNode := astNode.Children[2]

	var names []string = nil
	var emitEvaluables []IEvaluable = nil
	var isFullOosvars []bool = nil
	var indexEvaluables []IEvaluable = nil

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Things to be emitted, e.g. $a and $b in 'emit > "foo.dat", $a, $b'.

	// Non-lashed: emit @a, "x"
	// Lashed: emit (@a, @b), "x"
	numEmittables := len(emittablesNode.Children)
	names = make([]string, numEmittables)
	emitEvaluables = make([]IEvaluable, numEmittables)
	isFullOosvars = make([]bool, numEmittables)
	for i, emittableNode := range emittablesNode.Children {
		name, emitEvaluable, isFullOosvar, err := root.buildEmittableNode(emittableNode)
		if err != nil {
			return nil, err
		}
		names[i] = name
		emitEvaluables[i] = emitEvaluable
		isFullOosvars[i] = isFullOosvar
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Indices (if any) on the emittables

	isIndexed := false
	if keysNode.Type != dsl.NodeTypeNoOp { // There are "x","y" present
		lib.InternalCodingErrorIf(keysNode.Type != dsl.NodeTypeEmitKeys)
		isIndexed = true
		numKeys := len(keysNode.Children)
		indexEvaluables = make([]IEvaluable, numKeys)
		for i, keyNode := range keysNode.Children {
			indexEvaluable, err := root.BuildEvaluableNode(keyNode)
			if err != nil {
				return nil, err
			}
			indexEvaluables[i] = indexEvaluable
		}
	}

	retval := &EmitXStatementNode{
		names:           names,
		emitEvaluables:  emitEvaluables,
		indexEvaluables: indexEvaluables,
		isFullOosvars:   isFullOosvars,
		isEmitP:         isEmitP,
	}

	if !isIndexed {
		retval.executorFunc = retval.executeNonIndexed
	} else {
		retval.executorFunc = retval.executeIndexed
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Redirections and redirection targets (the thing after > >> |, if any).

	if redirectorNode.Type == dsl.NodeTypeNoOp {
		// No > >> or | was provided.
		retval.emitToRedirectFunc = retval.emitRecordToRecordStream
	} else {
		// There is > >> or | provided.
		lib.InternalCodingErrorIf(redirectorNode.Children == nil)
		lib.InternalCodingErrorIf(len(redirectorNode.Children) != 1)
		redirectorTargetNode := redirectorNode.Children[0]
		var err error = nil

		if redirectorTargetNode.Type == dsl.NodeTypeRedirectTargetStdout {
			retval.emitToRedirectFunc = retval.emitRecordToFileOrPipe
			retval.outputHandlerManager = output.NewStdoutWriteHandlerManager(root.recordWriterOptions)
			retval.redirectorTargetEvaluable = root.BuildStringLiteralNode("(stdout)")
		} else if redirectorTargetNode.Type == dsl.NodeTypeRedirectTargetStderr {
			retval.emitToRedirectFunc = retval.emitRecordToFileOrPipe
			retval.outputHandlerManager = output.NewStderrWriteHandlerManager(root.recordWriterOptions)
			retval.redirectorTargetEvaluable = root.BuildStringLiteralNode("(stderr)")
		} else {
			retval.emitToRedirectFunc = retval.emitRecordToFileOrPipe

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
// This is a helper method for deciding whether an emittable node is a named
// variable or a map.

func (root *RootNode) buildEmittableNode(
	astNode *dsl.ASTNode,
) (name string, emitEvaluable IEvaluable, isFullOosvar bool, err error) {
	name = ""
	emitEvaluable = nil
	isFullOosvar = true
	err = nil

	if astNode.Type == dsl.NodeTypeLocalVariable {
		name = string(astNode.Token.Lit)
		isFullOosvar = false
	} else if astNode.Type == dsl.NodeTypeDirectOosvarValue {
		name = string(astNode.Token.Lit)
		isFullOosvar = false
	} else if astNode.Type == dsl.NodeTypeDirectFieldValue {
		name = string(astNode.Token.Lit)
		isFullOosvar = false
	} else if astNode.Type == dsl.NodeTypeMapLiteral {
		name = "_"
		isFullOosvar = false
	}

	// xxx temp
	// ----------------------------------------------------------------
	// Emittable
	//   y LocalVariable
	//
	//   n FullOosvar
	//   y DirectOosvarValue -- includes BracedOosvarValue
	//  -> IndirectOosvarValue
	//
	//   n FullSrec
	//   y DirectFieldValue -- includes BracedFieldValue
	//  -> IndirectFieldValue
	//
	//   n MapLiteral
	// ;
	// ----------------------------------------------------------------

	emitEvaluable, err = root.BuildEvaluableNode(astNode)

	return name, emitEvaluable, isFullOosvar, err
}

// ================================================================
func (node *EmitXStatementNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	return node.executorFunc(state)
}

// ----------------------------------------------------------------
func (node *EmitXStatementNode) executeNonIndexed(
	state *runtime.State,
) (*BlockExitPayload, error) {

	newrec := types.NewMlrmapAsRecord()

	for i, emitEvaluable := range node.emitEvaluables {
		emittable := emitEvaluable.Evaluate(state)
		if emittable.IsAbsent() {
			continue
		}

		if node.isEmitP {
			if node.isFullOosvars[i] {
				// The top-level map for 'all oosvars' doesn't have a name --
				// other than '@*' of course. But we need to not include that
				// 'name' as the first output slot. Without this check,
				//
				//   mlr put -q @sum[$a][$b] += $x; end { emitp @* } foo.dat
				//
				// results in something like
				//
				//   {
				//     "_": {
				//       "sum": {
				//         "eks": {
				//           "wye": 0.38139939387114097,
				//           "zee": 0.6117840605678454
				//         }
				//       },
				//       "count": {
				//         "eks": {
				//           "wye": 2,
				//           "zee": 4,
				//         }
				//       },
				//     }
				//   }
				//
				// when we really want
				//
				//   {
				//     "sum": {
				//       "eks": {
				//         "wye": 0.38139939387114097,
				//         "zee": 0.6117840605678454
				//       }
				//     },
				//     "count": {
				//       "eks": {
				//         "wye": 2,
				//         "zee": 4,
				//       }
				//     },
				//   }

				top := emittable.GetMap()
				lib.InternalCodingErrorIf(top == nil)
				for pe := top.Head; pe != nil; pe = pe.Next {
					newrec.PutCopy(pe.Key, pe.Value)
				}
			} else {
				newrec.PutCopy(node.names[i], emittable)
			}
		} else {
			if emittable.IsMap() {
				newrec.Merge(emittable.GetMap())
			} else {
				newrec.PutCopy(node.names[i], emittable)
			}
		}
	}

	err := node.emitToRedirectFunc(newrec, state)

	return nil, err
}

// ----------------------------------------------------------------
func (node *EmitXStatementNode) executeIndexed(
	state *runtime.State,
) (*BlockExitPayload, error) {
	emittableMaps := make([]*types.Mlrmap, len(node.emitEvaluables))
	for i, emitEvaluable := range node.emitEvaluables {
		emittable := emitEvaluable.Evaluate(state)
		if emittable.IsAbsent() {
			return nil, nil
		}
		if !emittable.IsMap() {
			return nil, nil
		}
		emittableMaps[i] = emittable.GetMap()
	}
	indices := make([]*types.Mlrval, len(node.indexEvaluables))

	// TODO: libify this
	for i, _ := range node.indexEvaluables {
		indices[i] = node.indexEvaluables[i].Evaluate(state)
		if indices[i].IsAbsent() {
			return nil, nil
		}
		if indices[i].IsError() {
			// TODO: surface this more highly
			return nil, nil
		}
	}

	return node.executeIndexedAux(
		node.names,
		types.NewMlrmapAsRecord(),
		emittableMaps,
		indices,
		state,
	)
}

// Recurses over indices.
func (node *EmitXStatementNode) executeIndexedAux(
	mapNames []string,
	templateRecord *types.Mlrmap,
	emittableMaps []*types.Mlrmap,
	indices []*types.Mlrval,
	state *runtime.State,
) (*BlockExitPayload, error) {
	lib.InternalCodingErrorIf(len(indices) < 1)
	index := indices[0]
	indexString := index.String()

	for pe := emittableMaps[0].Head; pe != nil; pe = pe.Next {
		newrec := templateRecord.Copy()

		indexValue := types.MlrvalFromString(pe.Key)
		newrec.PutCopy(indexString, &indexValue)
		indexValueString := indexValue.String()

		nextLevels := make([]*types.Mlrval, len(emittableMaps))
		nextLevelMaps := make([]*types.Mlrmap, len(emittableMaps))
		for i, emittableMap := range emittableMaps {
			if emittableMap != nil {
				nextLevel := emittableMap.Get(indexValueString)
				nextLevels[i] = nextLevel
				// Can be nil for lashed indexing with heterogeneous data: e.g.
				// @x={"a":1}; @y={"b":2}; emit (@x, @y), "a"
				if nextLevel != nil && nextLevel.IsMap() {
					nextLevelMaps[i] = nextLevel.GetMap()
				} else {
					nextLevelMaps[i] = nil
				}
			} else {
				nextLevelMaps[i] = nil
			}
		}

		if nextLevelMaps[0] != nil && len(indices) >= 2 {
			// recurse
			node.executeIndexedAux(
				mapNames,
				newrec,
				nextLevelMaps,
				indices[1:],
				state,
			)
		} else {
			// end of recursion
			if node.isEmitP {
				for i, nextLevel := range nextLevels {
					if nextLevel != nil {
						if node.isFullOosvars[i] {
							// See extended comment for the non-indexed case, above.
							top := nextLevel.GetMap()
							lib.InternalCodingErrorIf(top == nil)
							for pe := top.Head; pe != nil; pe = pe.Next {
								newrec.PutCopy(pe.Key, pe.Value)
							}
						} else {
							newrec.PutCopy(mapNames[i], nextLevel)
						}
					}
				}
			} else {
				for i, nextLevel := range nextLevels {
					if nextLevel != nil {
						if nextLevel.IsMap() {
							newrec.Merge(nextLevelMaps[i])
						} else {
							newrec.PutCopy(mapNames[i], nextLevel)
						}
					}
				}
			}

			err := node.emitToRedirectFunc(newrec, state)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

// ----------------------------------------------------------------
func (node *EmitXStatementNode) emitRecordToRecordStream(
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
func (node *EmitXStatementNode) emitRecordToFileOrPipe(
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
