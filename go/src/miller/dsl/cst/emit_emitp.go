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
	"miller/dsl"
	"miller/lib"
	"miller/types"
)

// ================================================================
// Shared by emit and emitp

type EmitXStatementNode struct {
	// These are "_" for maps like in 'emit {...}'; "x" for named variables like
	// in 'emit @x'.
	names []string

	// Maps or named variables: the @a, @b parts.
	emitEvaluables []IEvaluable

	// The "x","y" parts.
	indexEvaluables []IEvaluable

	// Appropriate function to evaluate statements, depending on indexed or not.
	executorFunc Executor

	// For code-reuse between executors.
	isEmitP bool
}

func (this *RootNode) BuildEmitStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEmitStatement)
	return this.buildEmitXStatementNode(astNode, false)
}
func (this *RootNode) BuildEmitPStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEmitPStatement)
	return this.buildEmitXStatementNode(astNode, true)
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

func (this *RootNode) buildEmitXStatementNode(
	astNode *dsl.ASTNode,
	isEmitP bool,
) (IExecutable, error) {
	lib.InternalCodingErrorIf(len(astNode.Children) != 3)
	emittablesNode := astNode.Children[0]
	keysNode := astNode.Children[1]

	var names []string = nil
	var emitEvaluables []IEvaluable = nil
	var indexEvaluables []IEvaluable = nil

	// Lashed: emit (@a, @b), "x"
	numEmittables := len(emittablesNode.Children)
	names = make([]string, numEmittables)
	emitEvaluables = make([]IEvaluable, numEmittables)
	for i, emittableNode := range emittablesNode.Children {
		name, emitEvaluable, err := this.buildEmittableNode(emittableNode)
		if err != nil {
			return nil, err
		}
		names[i] = name
		emitEvaluables[i] = emitEvaluable
	}

	// xxx temp
	isIndexed := false
	if keysNode.Type != dsl.NodeTypeNoOp { // There are "x","y" present
		isIndexed = true
		lib.InternalCodingErrorIf(keysNode.Type != dsl.NodeTypeEmitKeys)
		numKeys := len(keysNode.Children)
		indexEvaluables = make([]IEvaluable, numKeys)
		for i, keyNode := range keysNode.Children {
			indexEvaluable, err := this.BuildEvaluableNode(keyNode)
			if err != nil {
				return nil, err
			}
			indexEvaluables[i] = indexEvaluable
		}
	}

	emitxStatementNode := &EmitXStatementNode{
		names:           names,
		emitEvaluables:  emitEvaluables,
		indexEvaluables: indexEvaluables,
		isEmitP:         isEmitP,
	}

	if !isIndexed {
		emitxStatementNode.executorFunc = emitxStatementNode.executeNonIndexed
	} else {
		emitxStatementNode.executorFunc = emitxStatementNode.executeIndexed
	}

	return emitxStatementNode, nil
}

// ----------------------------------------------------------------
// This is a helper method for deciding whether an emittable node is a named
// variable or a map.

func (this *RootNode) buildEmittableNode(
	astNode *dsl.ASTNode,
) (name string, emitEvaluable IEvaluable, err error) {
	name = "_"
	emitEvaluable = nil
	err = nil

	if astNode.Type == dsl.NodeTypeDirectOosvarValue {
		name = string(astNode.Token.Lit)
	} else if astNode.Type == dsl.NodeTypeLocalVariable {
		name = string(astNode.Token.Lit)
	} else if astNode.Type == dsl.NodeTypeDirectFieldValue {
		name = string(astNode.Token.Lit)
	}

	emitEvaluable, err = this.BuildEvaluableNode(astNode)

	return name, emitEvaluable, err
}

// ================================================================
func (this *EmitXStatementNode) Execute(state *State) (*BlockExitPayload, error) {
	return this.executorFunc(state)
}

// ----------------------------------------------------------------
func (this *EmitXStatementNode) executeNonIndexed(
	state *State,
) (*BlockExitPayload, error) {

	newrec := types.NewMlrmapAsRecord()

	for i, emitEvaluable := range this.emitEvaluables {
		emittable := emitEvaluable.Evaluate(state)
		if emittable.IsAbsent() {
			continue
		}

		if this.isEmitP {
			newrec.PutCopy(&this.names[i], &emittable)
		} else {
			if emittable.IsMap() {
				newrec.Merge(emittable.GetMap())
			} else {
				newrec.PutCopy(&this.names[i], &emittable)
			}
		}
	}

	state.OutputChannel <- types.NewRecordAndContext(
		newrec,
		state.Context.Copy(),
	)

	return nil, nil
}

// ----------------------------------------------------------------
func (this *EmitXStatementNode) executeIndexed(
	state *State,
) (*BlockExitPayload, error) {
	emittableMaps := make([]*types.Mlrmap, len(this.emitEvaluables))
	for i, emitEvaluable := range this.emitEvaluables {

		emittable := emitEvaluable.Evaluate(state)
		if emittable.IsAbsent() {
			return nil, nil
		}
		if !emittable.IsMap() {
			return nil, nil
		}
		emittableMaps[i] = emittable.GetMap()
	}

	// TODO: libify this
	indices := make([]types.Mlrval, len(this.indexEvaluables))
	for i, indexEvaluable := range this.indexEvaluables {
		index := indexEvaluable.Evaluate(state)
		if index.IsAbsent() {
			return nil, nil
		}
		if index.IsError() {
			// TODO: surface this more highly
			return nil, nil
		}
		indices[i] = index
	}

	return this.executeIndexedAux(
		this.names,
		types.NewMlrmapAsRecord(),
		emittableMaps,
		indices,
		state,
	)
}

// Recurses over indices.
func (this *EmitXStatementNode) executeIndexedAux(
	mapNames []string,
	templateRecord *types.Mlrmap,
	emittableMaps []*types.Mlrmap,
	indices []types.Mlrval,
	state *State,
) (*BlockExitPayload, error) {
	lib.InternalCodingErrorIf(len(indices) < 1)
	index := indices[0]
	indexString := index.String()

	for pe := emittableMaps[0].Head; pe != nil; pe = pe.Next {
		newrec := templateRecord.Copy()

		indexValue := types.MlrvalFromString(*pe.Key)
		newrec.PutCopy(&indexString, &indexValue)
		indexValueString := indexValue.String()

		nextLevels := make([]*types.Mlrval, len(emittableMaps))
		nextLevelMaps := make([]*types.Mlrmap, len(emittableMaps))
		for i, _ := range emittableMaps {
			nextLevel := emittableMaps[i].Get(&indexValueString)
			nextLevels[i] = nextLevel
			if nextLevel.IsMap() {
				nextLevelMaps[i] = nextLevel.GetMap()
			} else {
				nextLevelMaps[i] = nil
			}
		}

		if nextLevelMaps[0] != nil && len(indices) >= 2 {
			// recurse
			this.executeIndexedAux(
				mapNames,
				newrec,
				nextLevelMaps,
				indices[1:],
				state,
			)
		} else {
			// end of recursion
			if this.isEmitP {
				for i, nextLevel := range nextLevels {
					newrec.PutCopy(&mapNames[i], nextLevel)
				}
			} else {
				for i, nextLevel := range nextLevels {
					if nextLevel.IsMap() {
						newrec.Merge(nextLevelMaps[i])
					} else {
						newrec.PutCopy(&mapNames[i], nextLevel)
					}
				}
			}

			state.OutputChannel <- types.NewRecordAndContext(
				newrec,
				state.Context.Copy(),
			)
		}
	}

	return nil, nil
}
