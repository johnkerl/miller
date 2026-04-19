// This handles anything on the right-hand sides of assignment statements.
// (Also, computed field names on the left-hand sides of assignment
// statements.)

package cst

import (
	"fmt"
	"os"

	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/runtime"

	"github.com/johnkerl/pgpg/go/lib/pkg/asts"
)

func (root *RootNode) BuildEvaluableNode(astNode *asts.ASTNode) (IEvaluable, error) {
	// Try BuildLeafNode first for terminals
	if astNode.Children == nil || len(astNode.Children) == 0 {
		if leaf, err := root.BuildLeafNode(astNode); err == nil {
			return leaf, nil
		}
		// Fall through to switch for leaf types BuildLeafNode doesn't know
	}

	switch astNode.Type {

	case asts.NodeType(NodeTypeArrayLiteral): // [...]
		return root.BuildArrayLiteralNode(astNode)

	case asts.NodeType(NodeTypeMapLiteral): // {...}
		return root.BuildMapLiteralNode(astNode)

	case asts.NodeType(NodeTypeArrayOrMapIndexAccess): // x[...]
		return root.BuildArrayOrMapIndexAccessNode(astNode)

	case asts.NodeType(NodeTypeArraySliceLoHi), asts.NodeType(NodeTypeArraySliceHiOnly),
		asts.NodeType(NodeTypeArraySliceLoOnly), asts.NodeType(NodeTypeArraySliceFull): // myarray[lo:hi]
		return root.BuildArraySliceAccessNode(astNode)

	case asts.NodeType(NodeTypeIndirectFieldValue): // $[...] (includes $[[n]] and $[[[n]]])
		return root.BuildIndirectFieldValueNode(astNode)
	case asts.NodeType(NodeTypeIndirectOosvarValue): // $[...]
		return root.BuildIndirectOosvarValueNode(astNode)

	case asts.NodeType(NodeTypeEnvironmentVariable): // ENV["NAME"]
		return root.BuildEnvironmentVariableNode(astNode)

	// Operators are just functions with infix syntax so we treat them like
	// functions in the CST. (The distinction between infix syntax, e.g.
	// '1+2', and prefix syntax, e.g. 'plus(1,2)' disappears post-parse -- both
	// parse to the same-shape AST.)
	case asts.NodeType(NodeTypeOperator):
		return root.BuildFunctionCallsiteNode(astNode)
	case asts.NodeType(NodeTypeFunctionCallsite):
		return root.BuildFunctionCallsiteNode(astNode)

	// The dot operator is a little different from other operators since it's
	// type-dependent: for strings/int/bools etc it's just concatenation of
	// string representations, but if the left-hand side is a map, it's a
	// key-lookup with an unquoted literal on the right. E.g. mymap.foo is the
	// same as mymap["foo"].
	case asts.NodeType(NodeTypeDotOperator):
		return root.BuildDotCallsiteNode(astNode)

	case asts.NodeType(NodeTypeParenthesized):
		// (expr) — unwrap and build the inner expression
		lib.InternalCodingErrorIf(astNode.Children == nil || len(astNode.Children) != 1)
		return root.BuildEvaluableNode(astNode.Children[0])

	// Function literals like 'func (a,b) { return b - a }'
	case asts.NodeType(NodeTypeUnnamedFunctionDefinition):
		return root.BuildUnnamedUDFNode(astNode)

	// Leaf/terminal types (PGPG may give them non-nil empty Children)
	case asts.NodeType(NodeTypeDirectFieldValue), asts.NodeType(NodeTypeBracedFieldValue),
		asts.NodeType(NodeTypeFullSrec), asts.NodeType(NodeTypeDirectOosvarValue),
		asts.NodeType(NodeTypeBracedOosvarValue), asts.NodeType(NodeTypeFullOosvar),
		asts.NodeType(NodeTypeLocalVariable),
		asts.NodeType(NodeTypeIntLiteral), asts.NodeType(NodeTypeFloatLiteral),
		asts.NodeType(NodeTypeStringLiteral), asts.NodeType(NodeTypeBoolLiteral),
		asts.NodeType(NodeTypeNullLiteral), asts.NodeType(NodeTypeRegex):
		return root.BuildLeafNode(astNode)
	}

	// Parenthesized: (expr) — unwrap and build the inner expression.
	// Use string comparison in case asts.NodeType differs from astNode.Type's type.
	if string(astNode.Type) == NodeTypeParenthesized && astNode.Children != nil && len(astNode.Children) == 1 {
		return root.BuildEvaluableNode(astNode.Children[0])
	}

	// EnvironmentVariable: ENV["FOO"] or ENV.FOO
	if string(astNode.Type) == NodeTypeEnvironmentVariable && astNode.Children != nil && len(astNode.Children) == 1 {
		return root.BuildEnvironmentVariableNode(astNode)
	}

	// Fallback: try BuildLeafNode for unhandled types (e.g. DirectFieldValue, IntLiteral).
	// Only for leaf-like nodes (0 or 1 child); nodes with 2+ children are not leaves.
	if astNode.Children == nil || len(astNode.Children) <= 1 {
		if leaf, err := root.BuildLeafNode(astNode); err == nil {
			return leaf, nil
		}
	}

	return nil, fmt.Errorf(
		"at CST BuildEvaluableNode: unhandled AST node type %q (len=%d)", string(astNode.Type), len(astNode.Children),
	)
}

type IndirectFieldValueNode struct {
	fieldNameEvaluable IEvaluable
}

func (root *RootNode) BuildIndirectFieldValueNode(
	astNode *asts.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != asts.NodeType(NodeTypeIndirectFieldValue))
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	child := astNode.Children[0]
	if child.Type == asts.NodeType(NodeTypeArrayLiteral) && len(child.Children) == 1 {
		inner := child.Children[0]
		if inner.Type == asts.NodeType(NodeTypeArrayLiteral) && len(inner.Children) == 1 {
			// $[[[n]]] → positional field value
			indexASTNode := inner.Children[0]
			syntheticAST := asts.NewASTNode(nil, asts.NodeType(NodeTypePositionalFieldValue), []*asts.ASTNode{indexASTNode})
			return root.BuildPositionalFieldValueNode(syntheticAST)
		}
		// $[[n]] → positional field name
		indexASTNode := inner
		syntheticAST := asts.NewASTNode(nil, asts.NodeType(NodeTypePositionalFieldName), []*asts.ASTNode{indexASTNode})
		return root.BuildPositionalFieldNameNode(syntheticAST)
	}

	fieldNameEvaluable, err := root.BuildEvaluableNode(child)
	if err != nil {
		return nil, err
	}
	return &IndirectFieldValueNode{
		fieldNameEvaluable: fieldNameEvaluable,
	}, nil
}

func (node *IndirectFieldValueNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval { // TODO: err
	fieldName := node.fieldNameEvaluable.Evaluate(state)
	if fieldName.IsAbsent() {
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "$[(absent)]")
	}

	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "$*")
	}

	value, err := state.Inrec.GetWithMlrvalIndex(fieldName)
	if err != nil {
		// Key isn't int or string.
		// TODO: needs error-return in the API
		fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
		os.Exit(1)
	}
	if value == nil {
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "$["+fieldName.String()+"]")
	}
	return value
}

type IndirectOosvarValueNode struct {
	oosvarNameEvaluable IEvaluable
}

func (root *RootNode) BuildIndirectOosvarValueNode(
	astNode *asts.ASTNode,
) (*IndirectOosvarValueNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != asts.NodeType(NodeTypeIndirectOosvarValue))
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
) *mlrval.Mlrval { // TODO: err
	oosvarName := node.oosvarNameEvaluable.Evaluate(state)
	if oosvarName.IsAbsent() {
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "@[(absent)]")
	}

	value := state.Oosvars.Get(oosvarName.String())
	if value == nil {
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "@["+oosvarName.String()+"]")
	}

	return value
}
