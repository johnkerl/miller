package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
)

// This is for Lvalues, i.e. things on the left-hand-side of an assignment
// statement.

// ----------------------------------------------------------------
func BuildAssignableNode(
	astNode *dsl.ASTNode,
) (IAssignable, error) {

	switch astNode.Type {

	case dsl.NodeTypeDirectFieldValue:
		return BuildDirectFieldValueLvalueNode(astNode)
		break
	case dsl.NodeTypeIndirectFieldValue:
		return BuildIndirectFieldValueLvalueNode(astNode)
		break
	case dsl.NodeTypeFullSrec:
		return BuildFullSrecLvalueNode(astNode)
		break

	case dsl.NodeTypeDirectOosvarValue:
		return BuildDirectOosvarValueLvalueNode(astNode)
		break
	case dsl.NodeTypeIndirectOosvarValue:
		return BuildIndirectOosvarValueLvalueNode(astNode)
		break
	case dsl.NodeTypeFullOosvar:
		return BuildFullOosvarLvalueNode(astNode)
		break

	case dsl.NodeTypeIndexedLvalue:
		return BuildIndexedLvalueNode(astNode)
		break
	}

	return nil, errors.New(
		"CST BuildAssignableNode: unhandled AST node " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
type DirectFieldValueLvalueNode struct {
	lhsFieldName string
}

func BuildDirectFieldValueLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDirectFieldValue)

	lhsFieldName := string(astNode.Token.Lit)
	return NewDirectFieldValueLvalueNode(lhsFieldName), nil
}

func NewDirectFieldValueLvalueNode(lhsFieldName string) *DirectFieldValueLvalueNode {
	return &DirectFieldValueLvalueNode{
		lhsFieldName: lhsFieldName,
	}
}

func (this *DirectFieldValueLvalueNode) Assign(
	rvalue *lib.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *DirectFieldValueLvalueNode) AssignIndexed(
	rvalue *lib.Mlrval,
	indices []*lib.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	if indices == nil {
		state.Inrec.PutCopy(&this.lhsFieldName, rvalue)
		return nil
	} else {
		return state.Inrec.PutIndexed(&this.lhsFieldName, indices, rvalue)
	}
}

// ----------------------------------------------------------------
type IndirectFieldValueLvalueNode struct {
	lhsFieldNameExpression IEvaluable
}

func BuildIndirectFieldValueLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectFieldValue)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	lhsFieldNameExpression, err := BuildEvaluableNode(astNode.Children[0])
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
	rvalue *lib.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *IndirectFieldValueLvalueNode) AssignIndexed(
	rvalue *lib.Mlrval,
	indices []*lib.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	lhsFieldName := this.lhsFieldNameExpression.Evaluate(state)

	if !lhsFieldName.IsString() {
		return errors.New(
			"Miller DSL: computed field name [%s] should be a string but was " +
				lhsFieldName.GetTypeName() +
				".",
		)
	}

	sval := lhsFieldName.String()
	if indices == nil {
		state.Inrec.PutCopy(&sval, rvalue)
		return nil
	} else {
		return state.Inrec.PutIndexed(&sval, indices, rvalue)
	}
}

// ----------------------------------------------------------------
type FullSrecLvalueNode struct {
}

func BuildFullSrecLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeFullSrec)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(astNode.Children != nil)
	return NewFullSrecLvalueNode(), nil
}

func NewFullSrecLvalueNode() *FullSrecLvalueNode {
	return &FullSrecLvalueNode{}
}

func (this *FullSrecLvalueNode) Assign(
	rvalue *lib.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *FullSrecLvalueNode) AssignIndexed(
	rvalue *lib.Mlrval,
	indices []*lib.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	// The input record is a *Mlrmap so just invoke its PutIndexedKeyless.
	err := state.Inrec.PutIndexedKeyless(indices, rvalue)
	if err != nil {
		return err
	}
	return nil
}

// ----------------------------------------------------------------
type DirectOosvarValueLvalueNode struct {
	lhsOosvarName string
}

func BuildDirectOosvarValueLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDirectOosvarValue)

	lhsOosvarName := string(astNode.Token.Lit)
	return NewDirectOosvarValueLvalueNode(lhsOosvarName), nil
}

func NewDirectOosvarValueLvalueNode(lhsOosvarName string) *DirectOosvarValueLvalueNode {
	return &DirectOosvarValueLvalueNode{
		lhsOosvarName: lhsOosvarName,
	}
}

func (this *DirectOosvarValueLvalueNode) Assign(
	rvalue *lib.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *DirectOosvarValueLvalueNode) AssignIndexed(
	rvalue *lib.Mlrval,
	indices []*lib.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	if indices == nil {
		state.Oosvars.PutCopy(&this.lhsOosvarName, rvalue)
		return nil
	} else {
		return state.Oosvars.PutIndexed(&this.lhsOosvarName, indices, rvalue)
	}
}

// ----------------------------------------------------------------
type IndirectOosvarValueLvalueNode struct {
	lhsOosvarNameExpression IEvaluable
}

func BuildIndirectOosvarValueLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectOosvarValue)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	lhsOosvarNameExpression, err := BuildEvaluableNode(astNode.Children[0])
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
	rvalue *lib.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *IndirectOosvarValueLvalueNode) AssignIndexed(
	rvalue *lib.Mlrval,
	indices []*lib.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	lhsOosvarName := this.lhsOosvarNameExpression.Evaluate(state)

	if !lhsOosvarName.IsString() {
		return errors.New(
			"Miller DSL: computed field name [%s] should be a string but was " +
				lhsOosvarName.GetTypeName() +
				".",
		)
	}

	sval := lhsOosvarName.String()
	if indices == nil {
		state.Oosvars.PutCopy(&sval, rvalue)
		return nil
	} else {
		return state.Oosvars.PutIndexed(&sval, indices, rvalue)
	}
}

// ----------------------------------------------------------------
type FullOosvarLvalueNode struct {
}

func BuildFullOosvarLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeFullOosvar)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(astNode.Children != nil)
	return NewFullOosvarLvalueNode(), nil
}

func NewFullOosvarLvalueNode() *FullOosvarLvalueNode {
	return &FullOosvarLvalueNode{}
}

func (this *FullOosvarLvalueNode) Assign(
	rvalue *lib.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *FullOosvarLvalueNode) AssignIndexed(
	rvalue *lib.Mlrval,
	indices []*lib.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	// The input record is a *Mlrmap so just invoke its PutIndexedKeyless.
	err := state.Oosvars.PutIndexedKeyless(indices, rvalue)
	if err != nil {
		return err
	}
	return nil
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

func BuildIndexedLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndexedLvalue)
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
	//         * IndexedLvalue "[]"
	//             * IndexedLvalue "[]"
	//                 * DirectFieldValue "x"
	//                 * IntLiteral "1"
	//             * IntLiteral "2"
	//         * IntLiteral "3"

	// In the AST, the indices come in last-shallowest, down to first-deepest,
	// then the base Lvalue.
	walkerNode := astNode
	for {
		if walkerNode.Type == dsl.NodeTypeIndexedLvalue {
			lib.InternalCodingErrorIf(walkerNode == nil)
			lib.InternalCodingErrorIf(len(walkerNode.Children) != 2)
			indexEvaluable, err := BuildEvaluableNode(walkerNode.Children[1])
			if err != nil {
				return nil, err
			}
			indexEvaluables = append([]IEvaluable{indexEvaluable}, indexEvaluables...)
			walkerNode = walkerNode.Children[0]
		} else {
			baseLvalue, err = BuildAssignableNode(walkerNode)
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
	rvalue *lib.Mlrval,
	state *State,
) error {
	indices := make([]*lib.Mlrval, len(this.indexEvaluables))

	for i, indexEvaluable := range this.indexEvaluables {
		index := indexEvaluable.Evaluate(state)
		indices[i] = &index
	}
	return this.baseLvalue.AssignIndexed(rvalue, indices, state)
}

func (this *IndexedLvalueNode) AssignIndexed(
	rvalue *lib.Mlrval,
	indices []*lib.Mlrval,
	state *State,
) error {
	// We are the delegator, not the delegatee
	lib.InternalCodingErrorIf(true)
	return nil // not reached
}
