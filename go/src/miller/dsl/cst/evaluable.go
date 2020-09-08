package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
// This handles anything on the right-hand sides of assignment statements.
// (Also, computed field names on the left-hand sides of assignment
// statements.)
// ================================================================

// ----------------------------------------------------------------
func BuildEvaluableNode(astNode *dsl.ASTNode) (IEvaluable, error) {

	if astNode.Children == nil {
		return BuildLeafNode(astNode)
	}

	switch astNode.Type {

	case dsl.NodeTypeOperator:
		return BuildOperatorNode(astNode)

	case dsl.NodeTypeArrayLiteral:
		return BuildArrayLiteralNode(astNode)

	case dsl.NodeTypeMapLiteral:
		return BuildMapLiteralNode(astNode)

	case dsl.NodeTypeArrayOrMapIndexAccess:
		return BuildArrayOrMapIndexAccessNode(astNode)

	case dsl.NodeTypeArraySliceAccess:
		return BuildArraySliceAccessNode(astNode)

	case dsl.NodeTypeIndirectFieldName:
		return BuildIndirectFieldNameNode(astNode)
	}

	// xxx if/while/etc
	// xxx function
	// xxx more

	return nil, errors.New(
		"CST BuildEvaluableNode: unhandled AST node type " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
type IndirectFieldNameNode struct {
	fieldNameEvaluable IEvaluable
}

func BuildIndirectFieldNameNode(astNode *dsl.ASTNode) (*IndirectFieldNameNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectFieldName)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	fieldNameEvaluable, err := BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &IndirectFieldNameNode{
		fieldNameEvaluable: fieldNameEvaluable,
	}, nil
}
func (this *IndirectFieldNameNode) Evaluate(state *State) lib.Mlrval { // xxx err
	fieldName := this.fieldNameEvaluable.Evaluate(state)
	// xxx handle int-index too. needs a centralized place for that.
	if fieldName.IsAbsent() {
		return lib.MlrvalFromAbsent()
	}
	if !fieldName.IsString() {
		return lib.MlrvalFromError() // xxx needs err-return?
	}
	skey := fieldName.String()
	value := state.Inrec.Get(&skey)
	if value == nil {
		return lib.MlrvalFromAbsent()
	} else {
		return *value
	}
}
