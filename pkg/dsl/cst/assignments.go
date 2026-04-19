// CST build/execute for assignment and unset statements.

package cst

import (
	"fmt"

	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/runtime"
	"github.com/johnkerl/pgpg/go/lib/pkg/asts"
)

func (root *RootNode) BuildAssignmentNode(
	astNode *asts.ASTNode,
) (*AssignmentNode, error) {

	lib.InternalCodingErrorIf(astNode.Type != asts.NodeType(NodeTypeAssignment))
	err := astNode.CheckArity(2)
	if err != nil {
		return nil, err
	}

	lhsASTNode := astNode.Children[0]
	rhsASTNode := astNode.Children[1]

	lvalueNode, err := root.BuildAssignableNode(lhsASTNode)
	if err != nil {
		return nil, err
	}

	rvalueNode, err := root.BuildEvaluableNode(rhsASTNode)
	if err != nil {
		return nil, err
	}

	return &AssignmentNode{
		lvalueNode: lvalueNode,
		rvalueNode: rvalueNode,
	}, nil
}

func (root *RootNode) BuildCompoundAssignmentNode(
	astNode *asts.ASTNode,
) (*AssignmentNode, error) {

	lib.InternalCodingErrorIf(astNode.Type != asts.NodeType(NodeTypeCompoundAssignment))
	err := astNode.CheckArity(3)
	if err != nil {
		return nil, err
	}

	lhsASTNode := astNode.Children[0]
	opASTNode := astNode.Children[1]
	rhsASTNode := astNode.Children[2]

	compoundOp := tokenLit(opASTNode)
	baseOp := compoundOpToBaseOp(compoundOp)
	if baseOp == "" {
		return nil, fmt.Errorf("unknown compound assignment operator: %s", compoundOp)
	}

	lvalueNode, err := root.BuildAssignableNode(lhsASTNode)
	if err != nil {
		return nil, err
	}

	lvalueAsRvalue, err := root.BuildEvaluableNode(lhsASTNode)
	if err != nil {
		return nil, err
	}

	rvalueNode, err := root.BuildEvaluableNode(rhsASTNode)
	if err != nil {
		return nil, err
	}

	compoundRvalueNode, err := root.buildBinaryOperatorFromEvaluables(baseOp, lvalueAsRvalue, rvalueNode, rhsASTNode)
	if err != nil {
		return nil, err
	}

	return &AssignmentNode{
		lvalueNode: lvalueNode,
		rvalueNode: compoundRvalueNode,
	}, nil
}

// compoundOpToBaseOp maps compound assignment operators to their base operator.
// E.g. "+=" -> "+", "**=" -> "**".
func compoundOpToBaseOp(compound string) string {
	switch compound {
	case "||=":
		return "||"
	case "^^=":
		return "^^"
	case "&&=":
		return "&&"
	case "??=":
		return "??"
	case "???=":
		return "???"
	case "|=":
		return "|"
	case "&=":
		return "&"
	case "^=":
		return "^"
	case "<<=":
		return "<<"
	case ">>=":
		return ">>"
	case ">>>=":
		return ">>>"
	case "+=":
		return "+"
	case ".=":
		return "."
	case "-=":
		return "-"
	case "*=":
		return "*"
	case "/=":
		return "/"
	case "//=":
		return "//"
	case "%=":
		return "%"
	case "**=":
		return "**"
	default:
		return ""
	}
}

type AssignmentNode struct {
	lvalueNode IAssignable
	rvalueNode IEvaluable
}

func (node *AssignmentNode) Execute(
	state *runtime.State,
) (*BlockExitPayload, error) {
	rvalue := node.rvalueNode.Evaluate(state)
	if !rvalue.IsAbsent() {
		err := node.lvalueNode.Assign(rvalue, state)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (root *RootNode) BuildUnsetNode(
	astNode *asts.ASTNode,
) (*UnsetNode, error) {

	lib.InternalCodingErrorIf(astNode.Type != asts.NodeType(NodeTypeUnset))

	lvalueNodes := make([]IAssignable, len(astNode.Children))

	for i, lhsASTNode := range astNode.Children {
		var lvalueNode IAssignable
		var err error
		// "all" is a synonym for @* (full oosvar) in unset context.
		if lhsASTNode.Type == asts.NodeType(NodeTypeLocalVariable) && tokenLit(lhsASTNode) == "all" {
			lvalueNode = NewFullOosvarLvalueNode()
		} else {
			lvalueNode, err = root.BuildAssignableNode(lhsASTNode)
			if err != nil {
				return nil, err
			}
		}
		lvalueNodes[i] = lvalueNode
	}

	return &UnsetNode{
		lvalueNodes: lvalueNodes,
	}, nil
}

type UnsetNode struct {
	lvalueNodes []IAssignable
}

func (node *UnsetNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	for _, lvalueNode := range node.lvalueNodes {
		lvalueNode.Unassign(state)
	}
	return nil, nil
}
