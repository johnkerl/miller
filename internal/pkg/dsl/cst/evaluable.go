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

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
func (root *RootNode) BuildEvaluableNode(astNode *dsl.ASTNode) (IEvaluable, error) {

	if astNode.Children == nil {
		return root.BuildLeafNode(astNode)
	}

	switch astNode.Type {

	case dsl.NodeTypeArrayLiteral: // [...]
		return root.BuildArrayLiteralNode(astNode)

	case dsl.NodeTypeMapLiteral: // {...}
		return root.BuildMapLiteralNode(astNode)

	case dsl.NodeTypeArrayOrMapIndexAccess: // x[...]
		return root.BuildArrayOrMapIndexAccessNode(astNode)

	case dsl.NodeTypeArraySliceAccess: // myarray[lo:hi]
		return root.BuildArraySliceAccessNode(astNode)

	case dsl.NodeTypePositionalFieldName: // $[[...]]
		return root.BuildPositionalFieldNameNode(astNode)

	case dsl.NodeTypePositionalFieldValue: // $[[[...]]]
		return root.BuildPositionalFieldValueNode(astNode)

	case dsl.NodeTypeArrayOrMapPositionalNameAccess: // mymap[[...]]]
		return root.BuildArrayOrMapPositionalNameAccessNode(astNode)

	case dsl.NodeTypeArrayOrMapPositionalValueAccess: // mymap[[[...]]]
		return root.BuildArrayOrMapPositionalValueAccessNode(astNode)

	case dsl.NodeTypeIndirectFieldValue: // $[...]
		return root.BuildIndirectFieldValueNode(astNode)
	case dsl.NodeTypeIndirectOosvarValue: // $[...]
		return root.BuildIndirectOosvarValueNode(astNode)

	case dsl.NodeTypeEnvironmentVariable: // ENV["NAME"]
		return root.BuildEnvironmentVariableNode(astNode)

	// Operators are just functions with infix syntax so we treat them like
	// functions in the CST. (The distinction between infix syntax, e.g.
	// '1+2', and prefix syntax, e.g. 'plus(1,2)' disappears post-parse -- both
	// parse to the same-shape AST.)
	case dsl.NodeTypeOperator:
		return root.BuildFunctionCallsiteNode(astNode)
	case dsl.NodeTypeFunctionCallsite:
		return root.BuildFunctionCallsiteNode(astNode)

		// The dot operator is a little different from other operators since it's
		// type-dependent: for strings/int/bools etc it's just concatenation of
		// string representations, but if the left-hand side is a map, it's a
		// key-lookup with an unquoted literal on the right. E.g. mymap.foo is the
		// same as mymap["foo"].
	case dsl.NodeTypeDotOperator:
		return root.BuildDotCallsiteNode(astNode)

	// Function literals like 'func (a,b) { return b - a }'
	case dsl.NodeTypeUnnamedFunctionDefinition:
		return root.BuildUnnamedUDFNode(astNode)
	}

	return nil, errors.New(
		"CST BuildEvaluableNode: unhandled AST node type " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
type IndirectFieldValueNode struct {
	fieldNameEvaluable IEvaluable
}

func (root *RootNode) BuildIndirectFieldValueNode(
	astNode *dsl.ASTNode,
) (*IndirectFieldValueNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectFieldValue)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	fieldNameEvaluable, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &IndirectFieldValueNode{
		fieldNameEvaluable: fieldNameEvaluable,
	}, nil
}

func (node *IndirectFieldValueNode) Evaluate(
	state *runtime.State,
) *types.Mlrval { // TODO: err
	fieldName := node.fieldNameEvaluable.Evaluate(state)
	if fieldName.IsAbsent() {
		return types.MLRVAL_ABSENT
	}

	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return types.MLRVAL_ABSENT
	}

	value, err := state.Inrec.GetWithMlrvalIndex(fieldName)
	if err != nil {
		// Key isn't int or string.
		// TODO: needs error-return in the API
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if value == nil {
		return types.MLRVAL_ABSENT
	}
	return value
}

// ----------------------------------------------------------------
type IndirectOosvarValueNode struct {
	oosvarNameEvaluable IEvaluable
}

func (root *RootNode) BuildIndirectOosvarValueNode(
	astNode *dsl.ASTNode,
) (*IndirectOosvarValueNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectOosvarValue)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	oosvarNameEvaluable, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &IndirectOosvarValueNode{
		oosvarNameEvaluable: oosvarNameEvaluable,
	}, nil
}

func (node *IndirectOosvarValueNode) Evaluate(
	state *runtime.State,
) *types.Mlrval { // TODO: err
	oosvarName := node.oosvarNameEvaluable.Evaluate(state)
	if oosvarName.IsAbsent() {
		return types.MLRVAL_ABSENT
	}

	value := state.Oosvars.Get(oosvarName.String())
	if value == nil {
		return types.MLRVAL_ABSENT
	}

	return value
}
