package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
	"miller/types"
)

// This is for Lvalues, i.e. things on the left-hand-side of an assignment
// statement.

// ----------------------------------------------------------------
func (this *RootNode) BuildAssignableNode(
	astNode *dsl.ASTNode,
) (IAssignable, error) {

	switch astNode.Type {

	case dsl.NodeTypeDirectFieldValue:
		return this.BuildDirectFieldValueLvalueNode(astNode)
		break
	case dsl.NodeTypeIndirectFieldValue:
		return this.BuildIndirectFieldValueLvalueNode(astNode)
		break
	case dsl.NodeTypeFullSrec:
		return this.BuildFullSrecLvalueNode(astNode)
		break

	case dsl.NodeTypeDirectOosvarValue:
		return this.BuildDirectOosvarValueLvalueNode(astNode)
		break
	case dsl.NodeTypeIndirectOosvarValue:
		return this.BuildIndirectOosvarValueLvalueNode(astNode)
		break
	case dsl.NodeTypeFullOosvar:
		return this.BuildFullOosvarLvalueNode(astNode)
		break
	case dsl.NodeTypeLocalVariable:
		return this.BuildLocalVariableLvalueNode(astNode)
		break

	case dsl.NodeTypeArrayOrMapIndexAccess:
		return this.BuildIndexedLvalueNode(astNode)
		break
	}

	return nil, errors.New(
		"CST BuildAssignableNode: unhandled AST node " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
type DirectFieldValueLvalueNode struct {
	lhsFieldName *types.Mlrval
}

func (this *RootNode) BuildDirectFieldValueLvalueNode(
	astNode *dsl.ASTNode,
) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDirectFieldValue)

	lhsFieldName := types.MlrvalFromString(string(astNode.Token.Lit))
	return NewDirectFieldValueLvalueNode(&lhsFieldName), nil
}

func NewDirectFieldValueLvalueNode(lhsFieldName *types.Mlrval) *DirectFieldValueLvalueNode {
	return &DirectFieldValueLvalueNode{
		lhsFieldName: lhsFieldName,
	}
}

func (this *DirectFieldValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *DirectFieldValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	if indices == nil {
		err := state.Inrec.PutCopyWithMlrvalIndex(this.lhsFieldName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Inrec.PutIndexed(
			append([]*types.Mlrval{this.lhsFieldName}, indices...),
			rvalue,
		)
	}
}

func (this *DirectFieldValueLvalueNode) Unset(
	state *State,
) {
	this.UnsetIndexed(nil, state)
}

func (this *DirectFieldValueLvalueNode) UnsetIndexed(
	indices []*types.Mlrval,
	state *State,
) {
	if indices == nil {
		lib.InternalCodingErrorIf(!this.lhsFieldName.IsString())
		name := this.lhsFieldName.String()
		state.Inrec.Remove(&name)
	} else {
		state.Inrec.UnsetIndexed(
			append([]*types.Mlrval{this.lhsFieldName}, indices...),
		)
	}
}

// ----------------------------------------------------------------
type IndirectFieldValueLvalueNode struct {
	lhsFieldNameExpression IEvaluable
}

func (this *RootNode) BuildIndirectFieldValueLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectFieldValue)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	lhsFieldNameExpression, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return NewIndirectFieldValueLvalueNode(lhsFieldNameExpression), nil
}

func NewIndirectFieldValueLvalueNode(
	lhsFieldNameExpression IEvaluable,
) *IndirectFieldValueLvalueNode {
	return &IndirectFieldValueLvalueNode{
		lhsFieldNameExpression: lhsFieldNameExpression,
	}
}

func (this *IndirectFieldValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *IndirectFieldValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	lhsFieldName := this.lhsFieldNameExpression.Evaluate(state)

	if indices == nil {
		err := state.Inrec.PutCopyWithMlrvalIndex(&lhsFieldName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Inrec.PutIndexed(
			append([]*types.Mlrval{&lhsFieldName}, indices...),
			rvalue,
		)
	}
}

func (this *IndirectFieldValueLvalueNode) Unset(
	state *State,
) {
	this.UnsetIndexed(nil, state)
}

func (this *IndirectFieldValueLvalueNode) UnsetIndexed(
	indices []*types.Mlrval,
	state *State,
) {
	lhsFieldName := this.lhsFieldNameExpression.Evaluate(state)
	if indices == nil {
		name := lhsFieldName.String()
		state.Inrec.Remove(&name)
	} else {
		state.Inrec.UnsetIndexed(
			append([]*types.Mlrval{&lhsFieldName}, indices...),
		)
	}
}

// ----------------------------------------------------------------
type FullSrecLvalueNode struct {
}

func (this *RootNode) BuildFullSrecLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeFullSrec)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(astNode.Children != nil)
	return NewFullSrecLvalueNode(), nil
}

func NewFullSrecLvalueNode() *FullSrecLvalueNode {
	return &FullSrecLvalueNode{}
}

func (this *FullSrecLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *FullSrecLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	// The input record is a *Mlrmap so just invoke its PutIndexed.
	err := state.Inrec.PutIndexed(indices, rvalue)
	if err != nil {
		return err
	}
	return nil
}

func (this *FullSrecLvalueNode) Unset(
	state *State,
) {
	this.UnsetIndexed(nil, state)
}

func (this *FullSrecLvalueNode) UnsetIndexed(
	indices []*types.Mlrval,
	state *State,
) {
	if indices == nil {
		state.Inrec.Clear()
	} else {
		state.Inrec.UnsetIndexed(indices)
	}
}

// ----------------------------------------------------------------
type DirectOosvarValueLvalueNode struct {
	lhsOosvarName *types.Mlrval
}

func (this *RootNode) BuildDirectOosvarValueLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDirectOosvarValue)

	lhsOosvarName := types.MlrvalFromString(string(astNode.Token.Lit))
	return NewDirectOosvarValueLvalueNode(&lhsOosvarName), nil
}

func NewDirectOosvarValueLvalueNode(lhsOosvarName *types.Mlrval) *DirectOosvarValueLvalueNode {
	return &DirectOosvarValueLvalueNode{
		lhsOosvarName: lhsOosvarName,
	}
}

func (this *DirectOosvarValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *DirectOosvarValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	if indices == nil {
		err := state.Oosvars.PutCopyWithMlrvalIndex(this.lhsOosvarName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Oosvars.PutIndexed(
			append([]*types.Mlrval{this.lhsOosvarName}, indices...),
			rvalue,
		)
	}
}

func (this *DirectOosvarValueLvalueNode) Unset(
	state *State,
) {
	this.UnsetIndexed(nil, state)
}

func (this *DirectOosvarValueLvalueNode) UnsetIndexed(
	indices []*types.Mlrval,
	state *State,
) {
	if indices == nil {
		name := this.lhsOosvarName.String()
		state.Oosvars.Remove(&name)
	} else {
		state.Oosvars.UnsetIndexed(
			append([]*types.Mlrval{this.lhsOosvarName}, indices...),
		)
	}
}

// ----------------------------------------------------------------
type IndirectOosvarValueLvalueNode struct {
	lhsOosvarNameExpression IEvaluable
}

func (this *RootNode) BuildIndirectOosvarValueLvalueNode(
	astNode *dsl.ASTNode,
) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectOosvarValue)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	lhsOosvarNameExpression, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return NewIndirectOosvarValueLvalueNode(lhsOosvarNameExpression), nil
}

func NewIndirectOosvarValueLvalueNode(
	lhsOosvarNameExpression IEvaluable,
) *IndirectOosvarValueLvalueNode {
	return &IndirectOosvarValueLvalueNode{
		lhsOosvarNameExpression: lhsOosvarNameExpression,
	}
}

func (this *IndirectOosvarValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *IndirectOosvarValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	lhsOosvarName := this.lhsOosvarNameExpression.Evaluate(state)

	if indices == nil {
		err := state.Oosvars.PutCopyWithMlrvalIndex(&lhsOosvarName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Oosvars.PutIndexed(
			append([]*types.Mlrval{&lhsOosvarName}, indices...),
			rvalue,
		)
	}
}

func (this *IndirectOosvarValueLvalueNode) Unset(
	state *State,
) {
	this.UnsetIndexed(nil, state)
}

func (this *IndirectOosvarValueLvalueNode) UnsetIndexed(
	indices []*types.Mlrval,
	state *State,
) {
	name := this.lhsOosvarNameExpression.Evaluate(state)

	if indices == nil {
		sname := name.String()
		state.Oosvars.Remove(&sname)
	} else {
		state.Oosvars.UnsetIndexed(
			append([]*types.Mlrval{&name}, indices...),
		)
	}
}

// ----------------------------------------------------------------
type FullOosvarLvalueNode struct {
}

func (this *RootNode) BuildFullOosvarLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeFullOosvar)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(astNode.Children != nil)
	return NewFullOosvarLvalueNode(), nil
}

func NewFullOosvarLvalueNode() *FullOosvarLvalueNode {
	return &FullOosvarLvalueNode{}
}

func (this *FullOosvarLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *FullOosvarLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	// The input record is a *Mlrmap so just invoke its PutIndexed.
	err := state.Oosvars.PutIndexed(indices, rvalue)
	if err != nil {
		return err
	}
	return nil
}

func (this *FullOosvarLvalueNode) Unset(
	state *State,
) {
	this.UnsetIndexed(nil, state)
}

func (this *FullOosvarLvalueNode) UnsetIndexed(
	indices []*types.Mlrval,
	state *State,
) {
	if indices == nil {
		state.Oosvars.Clear()
	} else {
		state.Oosvars.UnsetIndexed(indices)
	}
}

// ----------------------------------------------------------------
type LocalVariableLvalueNode struct {
	typeGatedMlrvalName *types.TypeGatedMlrvalName

	// a = 1;
	// b = 1;
	// if (true) {
	//   a = 3;     <-- frameBind is false; updates outer a
	//   var b = 4; <-- frameBind is true;  creates new inner b
	// }
	frameBind bool
}

func (this *RootNode) BuildLocalVariableLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeLocalVariable)

	variableName := string(astNode.Token.Lit)
	typeName := "any"
	frameBind := false
	if astNode.Children != nil { // typed, like 'num x = 3'
		typeNode := astNode.Children[0]
		lib.InternalCodingErrorIf(typeNode.Type != dsl.NodeTypeTypedecl)
		typeName = string(typeNode.Token.Lit)
		frameBind = true
	}
	typeGatedMlrvalName, err := types.NewTypeGatedMlrvalName(
		variableName,
		typeName,
	)
	if err != nil {
		return nil, err
	}
	// TODO: type-gated mlrval
	return NewLocalVariableLvalueNode(
		typeGatedMlrvalName,
		frameBind,
	), nil
}

func NewLocalVariableLvalueNode(
	typeGatedMlrvalName *types.TypeGatedMlrvalName,
	frameBind bool,
) *LocalVariableLvalueNode {
	return &LocalVariableLvalueNode{
		typeGatedMlrvalName: typeGatedMlrvalName,
		frameBind:           frameBind,
	}
}

func (this *LocalVariableLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *LocalVariableLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	if indices == nil {
		err := this.typeGatedMlrvalName.Check(rvalue)
		if err != nil {
			return err
		}

		if this.frameBind {
			state.stack.BindVariable(this.typeGatedMlrvalName.Name, rvalue)
		} else {
			state.stack.SetVariable(this.typeGatedMlrvalName.Name, rvalue)
		}
		return nil
	} else {
		// TODO: propagate error return
		if this.frameBind {
			state.stack.BindVariableIndexed(this.typeGatedMlrvalName.Name, indices, rvalue)
		} else {
			state.stack.SetVariableIndexed(this.typeGatedMlrvalName.Name, indices, rvalue)
		}
		return nil
	}
}

func (this *LocalVariableLvalueNode) Unset(
	state *State,
) {
	this.UnsetIndexed(nil, state)
}

func (this *LocalVariableLvalueNode) UnsetIndexed(
	indices []*types.Mlrval,
	state *State,
) {
	if indices == nil {
		state.stack.UnsetVariable(this.typeGatedMlrvalName.Name)
	} else {
		state.stack.UnsetVariableIndexed(this.typeGatedMlrvalName.Name, indices)
	}
}

// ----------------------------------------------------------------
// IndexedValueNode is a delegator to base-lvalue types.
// * The baseLvalue is some IAssignable
// * The indexEvaluables are an array of IEvaluables
// * Each needs to evaluate to int or string
// * Assignment needs to walk each level:
//   o error if ith mlrval is int and that level isn't an array
//   o error if ith mlrval is string and that level isn't a map
//   o error for any other types -- maybe absent-handling for heterogeneity ...

// ----------------------------------------------------------------
type IndexedLvalueNode struct {
	baseLvalue      IAssignable
	indexEvaluables []IEvaluable
}

func (this *RootNode) BuildIndexedLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayOrMapIndexAccess)
	lib.InternalCodingErrorIf(astNode == nil)

	var baseLvalue IAssignable = nil
	indexEvaluables := make([]IEvaluable, 0)
	var err error = nil

	// $ mlr -n put -v '$x[1][2]=3'
	// DSL EXPRESSION:
	// $x[1][2]=3
	// RAW AST:
	// * StatementBlock
	//     * Assignment "="
	//         * ArrayOrMapIndexAccess "[]"
	//             * ArrayOrMapIndexAccess "[]"
	//                 * DirectFieldValue "x"
	//                 * IntLiteral "1"
	//             * IntLiteral "2"
	//         * IntLiteral "3"

	// In the AST, the indices come in last-shallowest, down to first-deepest,
	// then the base Lvalue.
	walkerNode := astNode
	for {
		if walkerNode.Type == dsl.NodeTypeArrayOrMapIndexAccess {
			lib.InternalCodingErrorIf(walkerNode == nil)
			lib.InternalCodingErrorIf(len(walkerNode.Children) != 2)
			indexEvaluable, err := this.BuildEvaluableNode(walkerNode.Children[1])
			if err != nil {
				return nil, err
			}
			indexEvaluables = append([]IEvaluable{indexEvaluable}, indexEvaluables...)
			walkerNode = walkerNode.Children[0]
		} else {
			baseLvalue, err = this.BuildAssignableNode(walkerNode)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	return NewIndexedLvalueNode(baseLvalue, indexEvaluables), nil
}

func NewIndexedLvalueNode(
	baseLvalue IAssignable,
	indexEvaluables []IEvaluable,
) *IndexedLvalueNode {
	return &IndexedLvalueNode{
		baseLvalue:      baseLvalue,
		indexEvaluables: indexEvaluables,
	}
}

// Computes Lvalue indices and then delegates to the baseLvalue.  E.g. for
// '$x[1][2] = 3' or '@x[1][2] = 3', the indices are [1,2], and the baseLvalue
// is '$x' or '@x' respectively.
func (this *IndexedLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	indices := make([]*types.Mlrval, len(this.indexEvaluables))

	for i, indexEvaluable := range this.indexEvaluables {
		index := indexEvaluable.Evaluate(state)
		indices[i] = &index
	}
	return this.baseLvalue.AssignIndexed(rvalue, indices, state)
}

func (this *IndexedLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// We are the delegator, not the delegatee
	lib.InternalCodingErrorIf(true)
	return nil // not reached
}

func (this *IndexedLvalueNode) Unset(
	state *State,
) {
	indices := make([]*types.Mlrval, len(this.indexEvaluables))

	for i, indexEvaluable := range this.indexEvaluables {
		index := indexEvaluable.Evaluate(state)
		indices[i] = &index
	}

	this.baseLvalue.UnsetIndexed(indices, state)
}

func (this *IndexedLvalueNode) UnsetIndexed(
	indices []*types.Mlrval,
	state *State,
) {
	// We are the delegator, not the delegatee
	lib.InternalCodingErrorIf(true)
}
