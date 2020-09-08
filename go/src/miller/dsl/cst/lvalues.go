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
	case dsl.NodeTypeDirectFieldName:
		return BuildDirectFieldNameLvalueNode(astNode)
		break
	case dsl.NodeTypeIndirectFieldName:
		return BuildIndirectFieldNameLvalueNode(astNode)
		break
	case dsl.NodeTypeFullSrec:
		return BuildFullSrecLvalueNode(astNode)
		break
	}

	// xxx temp
	return nil, errors.New(
		"CST BuildAssignableNode: unhandled AST node " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
type DirectFieldNameLvalueNode struct {
	lhsFieldName string
}

func BuildDirectFieldNameLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDirectFieldName)

	lhsFieldName := string(astNode.Token.Lit)
	return NewDirectFieldNameLvalueNode(lhsFieldName), nil
}

func NewDirectFieldNameLvalueNode(lhsFieldName string) *DirectFieldNameLvalueNode {
	return &DirectFieldNameLvalueNode{
		lhsFieldName: lhsFieldName,
	}
}

func (this *DirectFieldNameLvalueNode) Assign(
	rvalue *lib.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	state.Inrec.Put(&this.lhsFieldName, rvalue)
	return nil
}

// ----------------------------------------------------------------
type IndirectFieldNameLvalueNode struct {
	lhsFieldNameExpression IEvaluable
}

func BuildIndirectFieldNameLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectFieldName)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	lhsFieldNameExpression, err := BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return NewIndirectFieldNameLvalueNode(lhsFieldNameExpression), nil
}

func NewIndirectFieldNameLvalueNode(
	lhsFieldNameExpression IEvaluable,
) *IndirectFieldNameLvalueNode {
	return &IndirectFieldNameLvalueNode{
		lhsFieldNameExpression: lhsFieldNameExpression,
	}
}

func (this *IndirectFieldNameLvalueNode) Assign(
	rvalue *lib.Mlrval,
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

	state.Inrec.Put(&sval, rvalue)
	return nil
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
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	if !rvalue.IsMap() {
		// need 2nd-arg error in the API ... maybe
	}

	// xxx deepcopy!
	state.Inrec = rvalue.GetMap()

	return nil
}
