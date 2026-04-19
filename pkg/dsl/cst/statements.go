// CST build/execute for statements: assignments, bare booleans,
// break/continue/return, etc.

package cst

import (
	"fmt"

	"github.com/johnkerl/pgpg/go/lib/pkg/asts"
)

func (root *RootNode) BuildStatementNode(
	astNode *asts.ASTNode,
) (IExecutable, error) {

	var statement IExecutable = nil
	var err error
	switch astNode.Type {

	case asts.NodeType(NodeTypeAssignment):
		statement, err = root.BuildAssignmentNode(astNode)
		if err != nil {
			return nil, err
		}

	case asts.NodeType(NodeTypeCompoundAssignment):
		statement, err = root.BuildCompoundAssignmentNode(astNode)
		if err != nil {
			return nil, err
		}

	case asts.NodeType(NodeTypeUnset):
		statement, err = root.BuildUnsetNode(astNode)
		if err != nil {
			return nil, err
		}

	// E.g. 'NR > 10' without if or '{...}' body.  For put, these are no-ops
	// except side-effects (like regex-captures); for filter, they set the
	// filter condition only if they're the last statement in the main block.
	case asts.NodeType(NodeTypeBareBoolean):
		return root.BuildBareBooleanStatementNode(astNode)
	// E.g. 'filter NR > 10'.
	case asts.NodeType(NodeTypeFilterStatement):
		return root.BuildFilterStatementNode(astNode)

	case asts.NodeType(NodeTypePrintStatement):
		return root.BuildPrintStatementNode(astNode)
	case asts.NodeType(NodeTypePrintnStatement):
		return root.BuildPrintnStatementNode(astNode)
	case asts.NodeType(NodeTypeEprintStatement):
		return root.BuildEprintStatementNode(astNode)
	case asts.NodeType(NodeTypeEprintnStatement):
		return root.BuildEprintnStatementNode(astNode)

	case asts.NodeType(NodeTypeDumpStatement):
		return root.BuildDumpStatementNode(astNode)
	case asts.NodeType(NodeTypeEdumpStatement):
		return root.BuildEdumpStatementNode(astNode)

	case asts.NodeType(NodeTypeTeeStatement):
		return root.BuildTeeStatementNode(astNode)
	case asts.NodeType(NodeTypeEmit1Statement):
		return root.BuildEmit1StatementNode(astNode)
	case asts.NodeType(NodeTypeEmitStatement):
		return root.BuildEmitStatementNode(astNode)
	case asts.NodeType(NodeTypeEmitFStatement):
		return root.BuildEmitFStatementNode(astNode)
	case asts.NodeType(NodeTypeEmitPStatement):
		return root.BuildEmitPStatementNode(astNode)

	case asts.NodeType(NodeTypeBeginBlock):
		return nil, fmt.Errorf("begin blocks may only be declared at top level")
	case asts.NodeType(NodeTypeEndBlock):
		return nil, fmt.Errorf("end blocks may only be declared at top level")

	case asts.NodeType(NodeTypeIfChain):
		return root.BuildIfChainNode(astNode)
	case asts.NodeType(NodeTypeCondBlock):
		return root.BuildCondBlockNode(astNode)
	case asts.NodeType(NodeTypeWhileLoop):
		return root.BuildWhileLoopNode(astNode)
	case asts.NodeType(NodeTypeDoWhileLoop):
		return root.BuildDoWhileLoopNode(astNode)
	case asts.NodeType(NodeTypeForLoopOneVariable):
		return root.BuildForLoopOneVariableNode(astNode)
	case asts.NodeType(NodeTypeForLoopTwoVariable):
		return root.BuildForLoopTwoVariableNode(astNode)
	case asts.NodeType(NodeTypeForLoopMultivariable):
		return root.BuildForLoopMultivariableNode(astNode)
	case asts.NodeType(NodeTypeTripleForLoop):
		return root.BuildTripleForLoopNode(astNode)

	case asts.NodeType(NodeTypeNamedFunctionDefinition):
		return nil, fmt.Errorf("functions may only be declared at top level")
	case asts.NodeType(NodeTypeSubroutineDefinition):
		return nil, fmt.Errorf("subroutines may only be declared at top level")
	case asts.NodeType(NodeTypeSubroutineCallsite):
		return root.BuildSubroutineCallsiteNode(astNode)

	case asts.NodeType(NodeTypeBreakStatement):
		return root.BuildBreakNode(astNode)
	case asts.NodeType(NodeTypeContinueStatement):
		return root.BuildContinueNode(astNode)
	case asts.NodeType(NodeTypeReturnStatement):
		return root.BuildReturnNode(astNode)

	default:
		return nil, fmt.Errorf("at CST BuildStatementNode: unhandled AST node %s", string(astNode.Type))
	}
	return statement, nil
}
