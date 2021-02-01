// ================================================================
// This handles anything on the right-hand sides of assignment statements.
// (Also, computed field names on the left-hand sides of assignment
// statements.)
// ================================================================

package cst

import (
	"errors"
	"fmt"
	"os"

	"miller/dsl"
	"miller/lib"
	"miller/runtime"
	"miller/types"
)

// ----------------------------------------------------------------
func (this *RootNode) BuildEvaluableNode(astNode *dsl.ASTNode) (IEvaluable, error) {

	if astNode.Children == nil {
		return this.BuildLeafNode(astNode)
	}

	switch astNode.Type {

	case dsl.NodeTypeArrayLiteral: // [...]
		return this.BuildArrayLiteralNode(astNode)

	case dsl.NodeTypeMapLiteral: // {...}
		return this.BuildMapLiteralNode(astNode)

	case dsl.NodeTypeArrayOrMapIndexAccess: // x[...]
		return this.BuildArrayOrMapIndexAccessNode(astNode)

	case dsl.NodeTypeArraySliceAccess: // myarray[lo:hi]
		return this.BuildArraySliceAccessNode(astNode)

	case dsl.NodeTypePositionalFieldName: // $[[...]]
		return this.BuildPositionalFieldNameNode(astNode)

	case dsl.NodeTypePositionalFieldValue: // $[[[...]]]
		return this.BuildPositionalFieldValueNode(astNode)

	case dsl.NodeTypeArrayOrMapPositionalNameAccess: // mymap[[...]]]
		return this.BuildArrayOrMapPositionalNameAccessNode(astNode)

	case dsl.NodeTypeArrayOrMapPositionalValueAccess: // mymap[[[...]]]
		return this.BuildArrayOrMapPositionalValueAccessNode(astNode)

	case dsl.NodeTypeIndirectFieldValue: // $[...]
		return this.BuildIndirectFieldValueNode(astNode)
	case dsl.NodeTypeIndirectOosvarValue: // $[...]
		return this.BuildIndirectOosvarValueNode(astNode)

	case dsl.NodeTypeEnvironmentVariable: // ENV["NAME"]
		return this.BuildEnvironmentVariableNode(astNode)

	// Operators are just functions with infix syntax so we treat them like
	// functions in the CST. (The distinction between infix syntax, e.g.
	// '1+2', and prefix syntax, e.g. 'plus(1,2)' disappears post-parse -- both
	// parse to the same-shape AST.)
	case dsl.NodeTypeOperator:
		return this.BuildFunctionCallsiteNode(astNode)
	case dsl.NodeTypeFunctionCallsite:
		return this.BuildFunctionCallsiteNode(astNode)

	}

	return nil, errors.New(
		"CST BuildEvaluableNode: unhandled AST node type " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
type IndirectFieldValueNode struct {
	fieldNameEvaluable IEvaluable
}

func (this *RootNode) BuildIndirectFieldValueNode(
	astNode *dsl.ASTNode,
) (*IndirectFieldValueNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectFieldValue)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	fieldNameEvaluable, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &IndirectFieldValueNode{
		fieldNameEvaluable: fieldNameEvaluable,
	}, nil
}
func (this *IndirectFieldValueNode) Evaluate(state *runtime.State) types.Mlrval { // xxx err
	fieldName := this.fieldNameEvaluable.Evaluate(state)
	if fieldName.IsAbsent() {
		return types.MlrvalFromAbsent()
	}

	value, err := state.Inrec.GetWithMlrvalIndex(&fieldName)
	if err != nil {
		// Key isn't int or string.
		// TODO: needs error-return in the API
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if value == nil {
		return types.MlrvalFromAbsent()
	}
	return *value
}

// ----------------------------------------------------------------
type IndirectOosvarValueNode struct {
	oosvarNameEvaluable IEvaluable
}

func (this *RootNode) BuildIndirectOosvarValueNode(
	astNode *dsl.ASTNode,
) (*IndirectOosvarValueNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectOosvarValue)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	oosvarNameEvaluable, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &IndirectOosvarValueNode{
		oosvarNameEvaluable: oosvarNameEvaluable,
	}, nil
}

func (this *IndirectOosvarValueNode) Evaluate(state *runtime.State) types.Mlrval { // xxx err
	oosvarName := this.oosvarNameEvaluable.Evaluate(state)
	if oosvarName.IsAbsent() {
		return types.MlrvalFromAbsent()
	}

	value := state.Oosvars.Get(oosvarName.String())
	if value == nil {
		return types.MlrvalFromAbsent()
	}
	return *value
}
