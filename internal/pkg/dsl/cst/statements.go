// ================================================================
// CST build/execute for statements: assignments, bare booleans,
// break/continue/return, etc.
// ================================================================

package cst

import (
	"errors"

	"mlr/internal/pkg/dsl"
)

// ----------------------------------------------------------------
func (root *RootNode) BuildStatementNode(
	astNode *dsl.ASTNode,
) (IExecutable, error) {

	var statement IExecutable = nil
	var err error = nil
	switch astNode.Type {

	case dsl.NodeTypeAssignment:
		statement, err = root.BuildAssignmentNode(astNode)
		if err != nil {
			return nil, err
		}

	case dsl.NodeTypeUnset:
		statement, err = root.BuildUnsetNode(astNode)
		if err != nil {
			return nil, err
		}

	// E.g. 'NR > 10' without if or '{...}' body.  For put, these are no-ops
	// except side-effects (like regex-captures); for filter, they set the
	// filter condition only if they're the last statement in the main block.
	case dsl.NodeTypeBareBoolean:
		return root.BuildBareBooleanStatementNode(astNode)
	// E.g. 'filter NR > 10'.
	case dsl.NodeTypeFilterStatement:
		return root.BuildFilterStatementNode(astNode)

	case dsl.NodeTypePrintStatement:
		return root.BuildPrintStatementNode(astNode)
	case dsl.NodeTypePrintnStatement:
		return root.BuildPrintnStatementNode(astNode)
	case dsl.NodeTypeEprintStatement:
		return root.BuildEprintStatementNode(astNode)
	case dsl.NodeTypeEprintnStatement:
		return root.BuildEprintnStatementNode(astNode)

	case dsl.NodeTypeDumpStatement:
		return root.BuildDumpStatementNode(astNode)
	case dsl.NodeTypeEdumpStatement:
		return root.BuildEdumpStatementNode(astNode)

	case dsl.NodeTypeTeeStatement:
		return root.BuildTeeStatementNode(astNode)
	case dsl.NodeTypeEmit1Statement:
		return root.BuildEmit1StatementNode(astNode)
	case dsl.NodeTypeEmitStatement:
		return root.BuildEmitStatementNode(astNode)
	case dsl.NodeTypeEmitFStatement:
		return root.BuildEmitFStatementNode(astNode)
	case dsl.NodeTypeEmitPStatement:
		return root.BuildEmitPStatementNode(astNode)

	case dsl.NodeTypeBeginBlock:
		return nil, errors.New(
			"mlr: begin blocks may only be declared at top level.",
		)
	case dsl.NodeTypeEndBlock:
		return nil, errors.New(
			"mlr: end blocks may only be declared at top level.",
		)

	case dsl.NodeTypeIfChain:
		return root.BuildIfChainNode(astNode)
	case dsl.NodeTypeCondBlock:
		return root.BuildCondBlockNode(astNode)
	case dsl.NodeTypeWhileLoop:
		return root.BuildWhileLoopNode(astNode)
	case dsl.NodeTypeDoWhileLoop:
		return root.BuildDoWhileLoopNode(astNode)
	case dsl.NodeTypeForLoopOneVariable:
		return root.BuildForLoopOneVariableNode(astNode)
	case dsl.NodeTypeForLoopTwoVariable:
		return root.BuildForLoopTwoVariableNode(astNode)
	case dsl.NodeTypeForLoopMultivariable:
		return root.BuildForLoopMultivariableNode(astNode)
	case dsl.NodeTypeTripleForLoop:
		return root.BuildTripleForLoopNode(astNode)

	case dsl.NodeTypeNamedFunctionDefinition:
		return nil, errors.New(
			"mlr: functions may only be declared at top level.",
		)
	case dsl.NodeTypeSubroutineDefinition:
		return nil, errors.New(
			"mlr: subroutines may only be declared at top level.",
		)
	case dsl.NodeTypeSubroutineCallsite:
		return root.BuildSubroutineCallsiteNode(astNode)

	case dsl.NodeTypeBreak:
		return root.BuildBreakNode(astNode)
	case dsl.NodeTypeContinue:
		return root.BuildContinueNode(astNode)
	case dsl.NodeTypeReturn:
		return root.BuildReturnNode(astNode)

	default:
		return nil, errors.New(
			"CST BuildStatementNode: unhandled AST node " + string(astNode.Type),
		)
		break
	}
	return statement, nil
}
