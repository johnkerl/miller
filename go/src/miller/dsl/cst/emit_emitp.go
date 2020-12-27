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

	// For shared code use between emit and emitp.
	isEmitP bool

	executorFunc Executor
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
	nchild := len(astNode.Children)
	lib.InternalCodingErrorIf(nchild != 1 && nchild != 2)

	var names []string = nil
	var emitEvaluables []IEvaluable = nil
	var indexEvaluables []IEvaluable = nil

	emittablesNode := astNode.Children[0]
	if emittablesNode.Type == dsl.NodeTypeEmittableList {
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

	} else {
		// Non-lashed: emit @a, "x"
		names = make([]string, 1)
		emitEvaluables = make([]IEvaluable, 1)
		name, emitEvaluable, err := this.buildEmittableNode(emittablesNode)
		if err != nil {
			return nil, err
		}
		names[0] = name
		emitEvaluables[0] = emitEvaluable
	}

	if nchild == 2 { // There are "x","y" present
		keysNode := astNode.Children[1]
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

	if len(emitEvaluables) == 1 {
		if nchild == 1 {
			emitxStatementNode.executorFunc = emitxStatementNode.executeNonLashedNonIndexed
		} else {
			emitxStatementNode.executorFunc = emitxStatementNode.executeNonLashedIndexed
		}
	} else {
		if nchild == 1 {
			emitxStatementNode.executorFunc = emitxStatementNode.executeLashedNonIndexed
		} else {
			emitxStatementNode.executorFunc = emitxStatementNode.executeLashedIndexed
		}
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
// Just one mqp or named variable being emitted.
func (this *EmitXStatementNode) executeNonLashedNonIndexed(
	state *State,
) (*BlockExitPayload, error) {
	emittable := this.emitEvaluables[0].Evaluate(state)
	if emittable.IsAbsent() {
		return nil, nil
	}

	if emittable.IsMap() {
		// Emittable is map
		var newrec *types.Mlrmap = nil
		if this.isEmitP {
			newrec = types.NewMlrmapAsRecord()
			newrec.PutCopy(&this.names[0], &emittable)
		} else {
			newrec = emittable.Copy().GetMap()
		}
		state.OutputChannel <- types.NewRecordAndContext(
			newrec,
			state.Context.Copy(),
		)

	} else {
		// Emittable is named variable
		newrec := types.NewMlrmapAsRecord()
		newrec.PutCopy(&this.names[0], &emittable)
		state.OutputChannel <- types.NewRecordAndContext(
			newrec,
			state.Context.Copy(),
		)
	}

	return nil, nil
}

// ----------------------------------------------------------------
func (this *EmitXStatementNode) executeNonLashedIndexed(
	state *State,
) (*BlockExitPayload, error) {

	primaryEmittable := this.emitEvaluables[0].Evaluate(state)
	if primaryEmittable.IsAbsent() {
		return nil, nil
	}

	return nil, nil
}

// ----------------------------------------------------------------
func (this *EmitXStatementNode) executeLashedNonIndexed(
	state *State,
) (*BlockExitPayload, error) {

	newrec := types.NewMlrmapAsRecord()

	for i, emitEvaluable := range this.emitEvaluables {
		emittable := emitEvaluable.Evaluate(state)
		if emittable.IsAbsent() {
			continue
		}

		if emittable.IsMap() {
			// Emittable is map
			for pe := emittable.GetMap().Head; pe != nil; pe = pe.Next {
				newrec.PutCopy(pe.Key, pe.Value)
			}

		} else {
			// Emittable is named variable
			newrec.PutCopy(&this.names[i], &emittable)
		}
	}

	state.OutputChannel <- types.NewRecordAndContext(
		newrec,
		state.Context.Copy(),
	)

	return nil, nil
}

// ----------------------------------------------------------------
func (this *EmitXStatementNode) executeLashedIndexed(
	state *State,
) (*BlockExitPayload, error) {

	primaryEmittable := this.emitEvaluables[0].Evaluate(state)
	if primaryEmittable.IsAbsent() {
		return nil, nil
	}

	return nil, nil
}
