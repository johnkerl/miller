package cst

import (
	"errors"
	"fmt"
	"os"

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

	case dsl.NodeTypeIndirectFieldValue:
		return BuildIndirectFieldValueNode(astNode)
	}

	// xxx if/while/etc
	// xxx function
	// xxx more

	return nil, errors.New(
		"CST BuildEvaluableNode: unhandled AST node type " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
type IndirectFieldValueNode struct {
	fieldNameEvaluable IEvaluable
}

func BuildIndirectFieldValueNode(astNode *dsl.ASTNode) (*IndirectFieldValueNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectFieldValue)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	fieldNameEvaluable, err := BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &IndirectFieldValueNode{
		fieldNameEvaluable: fieldNameEvaluable,
	}, nil
}
func (this *IndirectFieldValueNode) Evaluate(state *State) lib.Mlrval { // xxx err
	fieldName := this.fieldNameEvaluable.Evaluate(state)
	if fieldName.IsAbsent() {
		return lib.MlrvalFromAbsent()
	}

	// Positional indices are supported, e.g. $[3] is the third field in the record.
	value, err := state.Inrec.GetWithMlrvalIndex(&fieldName)
	if err != nil {
		// Key isn't int or string.
		// xxx needs error-return in the API
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if value == nil {
		// E.g. $[7] but there aren't 7 fields in this record.
		return lib.MlrvalFromAbsent()
	}
	return *value
}
