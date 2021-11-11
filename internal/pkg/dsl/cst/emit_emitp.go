// ================================================================
// This handles emit and emitp statements. These produce new records (in
// addition to the current record, $*) into the output record stream.
//
// Some complications here are due to legacy.
//
// Emit statements existed in the Miller DSL before there were for-loops. As a
// result, some of the side-by-side emit syntaxes were invented (and supported)
// to allow things that might have been more easily done with simpler emit
// syntax.  Nonetheless, those syntaxes have been introduced into general use
// and we need to continue to support them.
//
// Examples for emit and emitp:
//
//   emit @a
//   emit (@a, @b)
//   emit {"a": @a, "b": @b}
//   emit @*
//
//   emit @a, "x", "y"
//   emit (@a, @b), "x", "y"
//   emit {"a": @a, "b": @b}, "x", "y"
//   emit @*, "x", "y"
//
// The first argument must be non-indexed oosvars/localvars/fieldnames, or
// a list thereof, so that we can use their names as keys in the emitted record
// -- or they must be maps. So the first complexity in this code is, do we have
// a named variable (or list thereof), or a map.
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

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/output"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
)

// ================================================================
// Shared by emit and emitp

type tEmitToRedirectFunc func(
	newrec *types.Mlrmap,
	state *runtime.State,
) error

type tEmitExecutorFunc func(
	names []string,
	values []*types.Mlrval,
	state *runtime.State,
) error

type EmitXStatementNode struct {
	// For 'emit @a' and 'emit (@a, @b)'
	topLevelNameList      []string
	topLevelEvaluableList []IEvaluable
	// For 'emit @*', 'emit @*', 'emit {...}'
	topLevelEvaluableMap IEvaluable

	// The "x","y" parts.
	indexEvaluables []IEvaluable

	// Appropriate function to evaluate statements, depending on indexed or not.
	executorFunc tEmitExecutorFunc

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
	isEmitP  bool
	isLashed bool

	// TODO: comment
	// root writerOptions AutoFlatten FLATSEP
	autoFlatten bool
	flatsep     string
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

var EMITX_NAMED_NODE_TYPES = map[dsl.TNodeType]bool{
	dsl.NodeTypeLocalVariable:       true,
	dsl.NodeTypeDirectOosvarValue:   true,
	dsl.NodeTypeIndirectOosvarValue: true,
	dsl.NodeTypeDirectFieldValue:    true,
	dsl.NodeTypeIndirectFieldValue:  true,
}

var EMITX_NAMELESS_NODE_TYPES = map[dsl.TNodeType]bool{
	dsl.NodeTypeFullSrec:   true,
	dsl.NodeTypeFullOosvar: true,
	dsl.NodeTypeMapLiteral: true,
}

// ----------------------------------------------------------------
// EMIT AND EMITP

func (root *RootNode) buildEmitXStatementNode(
	astNode *dsl.ASTNode,
	isEmitP bool,
) (IExecutable, error) {
	lib.InternalCodingErrorIf(len(astNode.Children) != 3)

	emittablesNode := astNode.Children[0]
	keysNode := astNode.Children[1]
	redirectorNode := astNode.Children[2]

	retval := &EmitXStatementNode{
		isEmitP:     isEmitP,
		isLashed:    false, // will be determined below
		autoFlatten: cli.DecideFinalFlatten(root.recordWriterOptions),
		flatsep:     root.recordWriterOptions.FLATSEP,
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Things to be emitted, e.g. $a and $b in 'emit > "foo.dat", ($a, $b), "x", "y"'.
	// Non-lashed: 'emit @a'
	// Lashed: 'emit (@a, @b)'

	lib.InternalCodingErrorIf(len(emittablesNode.Children) < 1)
	if len(emittablesNode.Children) == 1 {
		childNode := emittablesNode.Children[0]

		if EMITX_NAMED_NODE_TYPES[childNode.Type] {
			retval.topLevelNameList = make([]string, 1)
			retval.topLevelNameList[0] = string(childNode.Token.Lit)

			retval.topLevelEvaluableList = make([]IEvaluable, 1)
			evaluable, err := root.BuildEvaluableNode(childNode)
			if err != nil {
				return nil, err
			}
			retval.topLevelEvaluableList[0] = evaluable

		} else if EMITX_NAMELESS_NODE_TYPES[childNode.Type] {
			evaluable, err := root.BuildEvaluableNode(childNode)
			if err != nil {
				return nil, err
			}
			retval.topLevelEvaluableMap = evaluable

		} else {
			return nil, errors.New(
				fmt.Sprintf(
					"mlr: unlashe-demit node types must be local variables, field names, oosvars, or maps; got %s.",
					childNode.Type,
				),
			)
		}

	} else {
		retval.isLashed = true
		for _, childNode := range emittablesNode.Children {
			if !EMITX_NAMED_NODE_TYPES[childNode.Type] {
				return nil, errors.New(
					fmt.Sprintf(
						"mlr: lashed-emit node types must be local variables, field names, or oosvars; got %s.",
						childNode.Type,
					),
				)
			}
		}

		retval.topLevelNameList = make([]string, len(emittablesNode.Children))
		retval.topLevelEvaluableList = make([]IEvaluable, len(emittablesNode.Children))
		for i, childNode := range emittablesNode.Children {
			retval.topLevelNameList[i] = string(childNode.Token.Lit)
			evaluable, err := root.BuildEvaluableNode(childNode)
			if err != nil {
				return nil, err
			}
			retval.topLevelEvaluableList[i] = evaluable
		}
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Indices (if any) on the emittables

	isIndexed := false
	if keysNode.Type != dsl.NodeTypeNoOp { // There are "x","y" present
		lib.InternalCodingErrorIf(keysNode.Type != dsl.NodeTypeEmitKeys)
		isIndexed = true
		numKeys := len(keysNode.Children)
		retval.indexEvaluables = make([]IEvaluable, numKeys)
		for i, keyNode := range keysNode.Children {
			indexEvaluable, err := root.BuildEvaluableNode(keyNode)
			if err != nil {
				return nil, err
			}
			retval.indexEvaluables[i] = indexEvaluable
		}
	}

	if !isIndexed {
		if !retval.isLashed {
			if !isEmitP {
				retval.executorFunc = retval.executeNonIndexedNonLashedEmit
			} else {
				retval.executorFunc = retval.executeNonIndexedNonLashedEmitP
			}
		} else {
			if !isEmitP {
				retval.executorFunc = retval.executeNonIndexedLashedEmit
			} else {
				retval.executorFunc = retval.executeNonIndexedLashedEmitP
			}
		}
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

	// Register this with the CST root node so that open file descriptors can be
	// closed, etc at end of stream.
	if retval.outputHandlerManager != nil {
		root.RegisterOutputHandlerManager(retval.outputHandlerManager)
	}

	return retval, nil
}

// ================================================================
func (node *EmitXStatementNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	if node.topLevelEvaluableMap == nil {
		// 'emit @a', 'emit (@a, @b)', etc.
		names := node.topLevelNameList
		values := make([]*types.Mlrval, len(names))
		for i, evaluable := range node.topLevelEvaluableList {
			values[i] = evaluable.Evaluate(state)
		}
		return nil, node.executorFunc(names, values, state)

	} else {
		// 'emit @*', 'emit {...}', etc.
		parentValue := node.topLevelEvaluableMap.Evaluate(state)
		parentMapValue := parentValue.GetMap()
		if parentMapValue == nil {
			// TODO: what else to do if the should-be-a-map evaluates to:
			// * absent -- clearly returning is the right thing
			// * error -- what to emit?
			// * anything else other than a map -- ?
			return nil, nil
		}
		names := make([]string, parentMapValue.FieldCount)
		values := make([]*types.Mlrval, parentMapValue.FieldCount)

		i := 0
		for pe := parentMapValue.Head; pe != nil; pe = pe.Next {
			names[i] = pe.Key
			values[i] = pe.Value
			i++
		}

		return nil, node.executorFunc(names, values, state)
	}
}

// ----------------------------------------------------------------
// emit @* (supposing @a and @b exist) means @a and @b material are
// emitted in separate records.
//
// Example:
// DSL expression: @sum[$a][$b] += $n; end { dump; emit @sum }
// Name: "sum"
// Values: single array contaning the map
//   {
//     "sum": {
//       "vee": {
//         "wye": 2,
//         "zee": 4
//       },
//       "eks": {
//         "wye": 6,
//         "zee": 8
//       }
//     }
//   }
// Desired output:
//   {
//     "wye": 2,
//     "zee": 4
//   }
//   {
//     "wye": 6,
//     "zee": 8
//   }

func (node *EmitXStatementNode) executeNonIndexedNonLashedEmit(
	names []string,
	values []*types.Mlrval,
	state *runtime.State,
) error {
	for i, value := range values {
		if value.IsAbsent() {
			continue
		}

		valueAsMap := value.GetMap() // nil if not a map

		if valueAsMap == nil {
			newrec := types.NewMlrmapAsRecord()
			newrec.PutCopy(names[i], value)
			err := node.emitToRedirectFunc(newrec, state)
			if err != nil {
				return err
			}

		} else {
			recurse := valueAsMap.IsNested()
			if !recurse {
				newrec := types.NewMlrmapAsRecord()
				for pe := valueAsMap.Head; pe != nil; pe = pe.Next {
					newrec.PutCopy(pe.Key, pe.Value)
				}
				err := node.emitToRedirectFunc(newrec, state)
				if err != nil {
					return err
				}

			} else { // recurse
				nextLevelNames := make([]string, 0)
				nextLevelValues := make([]*types.Mlrval, 0)
				for pe := value.GetMap().Head; pe != nil; pe = pe.Next {
					nextLevelNames = append(nextLevelNames, pe.Key)
					nextLevelValues = append(nextLevelValues, pe.Value.Copy())
				}
				node.executeNonIndexedNonLashedEmit(nextLevelNames, nextLevelValues, state)
			}
		}
	}

	return nil
}

func (node *EmitXStatementNode) executeNonIndexedNonLashedEmitP(
	names []string,
	values []*types.Mlrval,
	state *runtime.State,
) error {
	for i, value := range values {
		if value.IsAbsent() {
			continue
		}
		newrec := types.NewMlrmapAsRecord()
		newrec.PutCopy(names[i], value)
		err := node.emitToRedirectFunc(newrec, state)
		if err != nil {
			return err
		}
	}

	return nil
}

// ----------------------------------------------------------------
// emit (@a, $b) means @a and @b material are lashed together in the
// same record.

func (node *EmitXStatementNode) executeNonIndexedLashedEmit(
	names []string,
	values []*types.Mlrval,
	state *runtime.State,
) error {
	lib.InternalCodingErrorIf(len(values) < 1)
	leadingValueAsMap := values[0].GetMap()
	if leadingValueAsMap == nil {
		// Emit a record like a=1,b=2
		newrec := types.NewMlrmapAsRecord()
		for i, value := range values {
			if value.IsAbsent() {
				continue
			}
			if value.IsMap() {
				newrec.Merge(value.GetMap())
			} else {
				newrec.PutCopy(names[i], value)
			}
		}
		return node.emitToRedirectFunc(newrec, state)

	} else {
		for i, value := range values {
			if value.IsAbsent() {
				continue
			}

			valueAsMap := value.GetMap() // nil if not a map

			if valueAsMap == nil {
				newrec := types.NewMlrmapAsRecord()
				newrec.PutCopy(names[i], value)
				err := node.emitToRedirectFunc(newrec, state)
				if err != nil {
					return err
				}

			} else {
				recurse := valueAsMap.IsNested()
				if !recurse {
					newrec := types.NewMlrmapAsRecord()
					for pe := valueAsMap.Head; pe != nil; pe = pe.Next {
						newrec.PutCopy(pe.Key, pe.Value)
					}
					err := node.emitToRedirectFunc(newrec, state)
					if err != nil {
						return err
					}

				} else { // recurse
					nextLevelNames := make([]string, 0)
					nextLevelValues := make([]*types.Mlrval, 0)
					for pe := value.GetMap().Head; pe != nil; pe = pe.Next {
						nextLevelNames = append(nextLevelNames, pe.Key)
						nextLevelValues = append(nextLevelValues, pe.Value.Copy())
					}
					node.executeNonIndexedNonLashedEmit(nextLevelNames, nextLevelValues, state)
				}
			}
		}

		return nil
	}
}

func (node *EmitXStatementNode) executeNonIndexedLashedEmitP(
	names []string,
	values []*types.Mlrval,
	state *runtime.State,
) error {
	newrec := types.NewMlrmapAsRecord()

	for i, value := range values {
		if value.IsAbsent() {
			continue
		}

		newrec.PutCopy(names[i], value)
	}

	return node.emitToRedirectFunc(newrec, state)
}

// ----------------------------------------------------------------
func (node *EmitXStatementNode) executeIndexed(
	names []string,
	values []*types.Mlrval,
	state *runtime.State,
) error {

	emittableMaps := make([]*types.Mlrmap, len(values))
	for i, value := range values {
		if value.IsAbsent() {
			return nil
		}
		mapValue := value.GetMap()
		if mapValue == nil {
			return nil
		}
		emittableMaps[i] = mapValue
	}

	// TODO: libify this
	indices := make([]*types.Mlrval, len(node.indexEvaluables))
	for i := range node.indexEvaluables {
		indices[i] = node.indexEvaluables[i].Evaluate(state)
		if indices[i].IsAbsent() {
			return nil
		}
		if indices[i].IsError() {
			// TODO: surface this more highly
			return nil
		}
	}

	if !node.isEmitP {
		if !node.isLashed {
			return node.executeIndexedNonLashedEmitAux(
				types.NewMlrmapAsRecord(),
				names,
				emittableMaps,
				indices,
				state,
			)
		} else {
			return node.executeIndexedLashedEmitAux(
				types.NewMlrmapAsRecord(),
				names,
				emittableMaps,
				indices,
				state,
			)
		}
	} else {
		if !node.isLashed {
			return node.executeIndexedNonLashedEmitPAux(
				types.NewMlrmapAsRecord(),
				names,
				emittableMaps,
				indices,
				state,
			)
		} else {
			return node.executeIndexedLashedEmitPAux(
				types.NewMlrmapAsRecord(),
				names,
				emittableMaps,
				indices,
				state,
			)
		}
	}
}

// ----------------------------------------------------------------
// Recurses over indices.
//
// Example:
// DSL expression:
//   @sum[$a][$b] += $n; end { dump @sum; emit @sum, "a", "b" }
// Input data (DKVP):
//   a=vee,b=wye,n=2
//   a=vee,b=zee,n=4
//   a=eks,b=wye,n=6
//   a=eks,b=zee,n=8
// @sum data structure at end:
//   {
//     "vee": {
//       "wye": 2,
//       "zee": 4
//     },
//     "eks": {
//       "wye": 6,
//       "zee": 8
//     }
//   }
//
// Output data (JSON):
//   {
//     "a": "vee",
//     "b": "wye",
//     "sum": 2
//   }
//   {
//     "a": "vee",
//     "b": "zee",
//     "sum": 4
//   }
//   {
//     "a": "eks",
//     "b": "wye",
//     "sum": 6
//   }
//   {
//     "a": "eks",
//     "b": "zee",
//     "sum": 8
//   }
//
// Outer call:
// * names = ["sum"]
// * templateRecord is empty
// * emittableMaps = [ {"vee": { "wye": 2, "zee": 4 }, "eks": { "wye": 6, "zee": 8 }} ]
// * indices = ["a", "b"]
//
// Inner call 1:
// * names = ["sum"]
// * templateRecord is {"a":"vee"}
// * emittableMaps = [{ "wye": 2, "zee": 4 }]
// * indices = ["b"]
//
// Inner call 1:
// * names = ["sum"]
// * templateRecord is {"a":"eks"}
// * emittableMaps = [{ "wye": 6, "zee": 8 }]
// * indices = ["b"]

func (node *EmitXStatementNode) executeIndexedNonLashedEmitAux(
	templateRecord *types.Mlrmap,
	names []string,
	emittableMaps []*types.Mlrmap,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	lib.InternalCodingErrorIf(len(indices) < 1)
	index := indices[0]
	indexString := index.String()

	for i, emittableMap := range emittableMaps {
		for pe := emittableMap.Head; pe != nil; pe = pe.Next {
			newrec := templateRecord.Copy()
			newrec.PutCopy(indexString, types.MlrvalFromString(pe.Key))

			if len(indices) == 1 {
				valueAsMap := pe.Value.GetMap()
				if valueAsMap == nil {
					newrec.PutCopy(names[i], pe.Value)
				} else {
					for pe := valueAsMap.Head; pe != nil; pe = pe.Next {
						newrec.PutCopy(pe.Key, pe.Value)
					}
				}
				err := node.emitToRedirectFunc(newrec, state)
				if err != nil {
					return err
				}
			} else { // recurse
				valueAsMap := pe.Value.GetMap()
				if valueAsMap == nil {
					newrec.PutCopy(names[i], pe.Value)
					err := node.emitToRedirectFunc(newrec, state)
					if err != nil {
						return err
					}
				} else {
					node.executeIndexedNonLashedEmitPAux(
						newrec,
						[]string{names[i]},
						[]*types.Mlrmap{valueAsMap},
						indices[1:],
						state,
					)
				}
			}
		}
	}

	return nil
}

// ----------------------------------------------------------------
// Example:
//
// DSL expression: @count[$a][$b] += 1; @sum[$a][$b] += $n; end { emit (@count, @sum), "a", "b" }
//
// @count and @sum maps:
//   {
//     "count": {
//       "vee": {
//         "wye": 1,
//         "zee": 1
//       },
//       "eks": {
//         "wye": 1,
//         "zee": 1
//       }
//     },
//     "sum": {
//       "vee": {
//         "wye": 2,
//         "zee": 4
//       },
//       "eks": {
//         "wye": 6,
//         "zee": 8
//       }
//     }
//   }
//
// Desired output:
//   {
//     "a": "vee",
//     "b": "wye",
//     "count": 1,
//     "sum": 2
//   }
//   {
//     "a": "vee",
//     "b": "zee",
//     "count": 1,
//     "sum": 4
//   }
//   {
//     "a": "eks",
//     "b": "wye",
//     "count": 1,
//     "sum": 6
//   }
//   {
//     "a": "eks",
//     "b": "zee",
//     "count": 1,
//     "sum": 8
//   }
//
// First call:
// * templateRecord is empty
// * names ["count", "sum"]
// * emittableMaps
//     {                {
//       "vee": {         "vee": {
//         "wye": 1,        "wye": 2,
//         "zee": 1         "zee": 4
//       },               },
//       "eks": {         "eks": {
//         "wye": 1,        "wye": 6,
//         "zee": 1         "zee": 8
//       }                }
//     }                }
// * indices ["a", "b"]
//
// * Loop over first-level keys of the leading map which is the "count" map: "vee" and "eks"
//
// * Recurse with:
//   o templateRecord {"a":"vee"}
//   o names ["count", "sum"]
//   o emittableMaps
//       {                {
//         "wye": 1,        "wye": 2,
//         "zee": 1         "zee": 4
//       }                }
//   o indices ["b"]
//
// * Recurse with:
//   o templateRecord {"a":"eks"}
//   o names ["count", "sum"]
//   o emittableMaps
//       {                {
//         "wye": 1,        "wye": 6,
//         "zee": 1         "zee": 8
//       }                }
//   o indices ["b"]

func (node *EmitXStatementNode) executeIndexedLashedEmitAux(
	templateRecord *types.Mlrmap,
	names []string,
	emittableMaps []*types.Mlrmap,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	lib.InternalCodingErrorIf(len(indices) < 1)
	index := indices[0]
	indexString := index.String()

	leadingMap := emittableMaps[0]

	for pe := leadingMap.Head; pe != nil; pe = pe.Next {
		newrec := templateRecord.Copy()
		indexValue := types.MlrvalFromString(pe.Key)
		newrec.PutCopy(indexString, indexValue)

		nextLevelValues := make([]*types.Mlrval, len(emittableMaps))
		nextLevelMaps := make([]*types.Mlrmap, len(emittableMaps))
		for i, emittableMap := range emittableMaps {
			if emittableMap != nil {
				nextLevelValues[i] = emittableMaps[i].Get(pe.Key)
				// Can be nil for lashed indexing with heterogeneous data: e.g.
				// @x={"a":1}; @y={"b":2}; emit (@x, @y), "a"
				if nextLevelValues[i] != nil && nextLevelValues[i].IsMap() {
					nextLevelMaps[i] = nextLevelValues[i].GetMap().Copy()
				} else {
					nextLevelMaps[i] = nil
				}
			} else {
				nextLevelMaps = append(nextLevelMaps, nil)
			}
		}

		if len(indices) > 1 && nextLevelMaps[0] != nil {
			// Recurse.  The leading map drives the iteration; we don't
			// continue even if other maps aren't empty
			node.executeIndexedLashedEmitAux(
				newrec,
				names,
				nextLevelMaps,
				indices[1:],
				state,
			)
		} else { // end of recursion
			for i, nextLevelValue := range nextLevelValues {
				if nextLevelValue != nil {
					if nextLevelMaps[i] != nil {
						newrec.Merge(nextLevelMaps[i])
					} else {
						newrec.PutCopy(names[i], nextLevelValue)
					}
				}
			}

			err := node.emitToRedirectFunc(newrec, state)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ----------------------------------------------------------------
// Recurses over indices.

func (node *EmitXStatementNode) executeIndexedNonLashedEmitPAux(
	templateRecord *types.Mlrmap,
	names []string,
	emittableMaps []*types.Mlrmap,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	lib.InternalCodingErrorIf(len(indices) < 1)
	index := indices[0]
	indexString := index.String()

	for i, emittableMap := range emittableMaps {
		for pe := emittableMap.Head; pe != nil; pe = pe.Next {
			newrec := templateRecord.Copy()
			newrec.PutCopy(indexString, types.MlrvalFromString(pe.Key))

			if len(indices) == 1 {
				newrec.PutCopy(names[i], pe.Value)
				err := node.emitToRedirectFunc(newrec, state)
				if err != nil {
					return err
				}
			} else { // recurse
				valueAsMap := pe.Value.GetMap()
				if valueAsMap == nil {
					newrec.PutCopy(names[i], pe.Value)
					err := node.emitToRedirectFunc(newrec, state)
					if err != nil {
						return err
					}
				} else {
					node.executeIndexedNonLashedEmitPAux(
						newrec,
						[]string{names[i]},
						[]*types.Mlrmap{valueAsMap},
						indices[1:],
						state,
					)
				}
			}
		}
	}

	return nil
}

func (node *EmitXStatementNode) executeIndexedLashedEmitPAux(
	templateRecord *types.Mlrmap,
	names []string,
	emittableMaps []*types.Mlrmap,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	lib.InternalCodingErrorIf(len(indices) < 1)
	index := indices[0]
	indexString := index.String()

	leadingMap := emittableMaps[0]

	for pe := leadingMap.Head; pe != nil; pe = pe.Next {
		newrec := templateRecord.Copy()

		indexValue := types.MlrvalFromString(pe.Key)
		newrec.PutCopy(indexString, indexValue)
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
					// xxx need to put names[i] and the accumulator value
				} else {
					nextLevelMaps[i] = nil
				}
			} else {
				nextLevelMaps[i] = nil
			}
		}

		if nextLevelMaps[0] != nil && len(indices) >= 2 {
			// recurse
			node.executeIndexedLashedEmitPAux(
				newrec,
				names,
				nextLevelMaps,
				indices[1:],
				state,
			)
		} else {
			// end of recursion
			for i, nextLevel := range nextLevels {
				if nextLevel != nil {
					newrec.PutCopy(names[i], nextLevel)
				}
			}

			err := node.emitToRedirectFunc(newrec, state)
			if err != nil {
				return err
			}
		}

	}

	return nil
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

	//fmt.Println("PRE")
	//outrec.Dump()
	//fmt.Println("POST")
	if node.autoFlatten {
		outrec.Flatten(node.flatsep)
	}
	return node.outputHandlerManager.WriteRecordAndContext(
		types.NewRecordAndContext(outrec, state.Context),
		outputFileName,
	)
}
