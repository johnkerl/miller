package cst

import (
	"errors"

	"miller/dsl"
)

// ================================================================
// CST build/execute for statements: assignments, bare booleans,
// break/continue/return, etc.
// ================================================================

// ----------------------------------------------------------------
func (this *RootNode) BuildStatementNode(
	astNode *dsl.ASTNode,
) (IExecutable, error) {

	var statement IExecutable = nil
	var err error = nil
	// xxx more to do
	switch astNode.Type {

	case dsl.NodeTypeAssignment:
		statement, err = this.BuildAssignmentNode(astNode)
		if err != nil {
			return nil, err
		}

	case dsl.NodeTypeUnset:
		statement, err = this.BuildUnsetNode(astNode)
		if err != nil {
			return nil, err
		}

	case dsl.NodeTypeFilterStatement:
		return this.BuildFilterStatementNode(astNode)
	case dsl.NodeTypeBareBoolean:
		return this.BuildFilterStatementNode(astNode)

	case dsl.NodeTypeEmitStatement:
		return this.BuildEmitStatementNode(astNode)
	case dsl.NodeTypeDumpStatement:
		return this.BuildDumpStatementNode(astNode)
	case dsl.NodeTypeEdumpStatement:
		return this.BuildEdumpStatementNode(astNode)
	case dsl.NodeTypePrintStatement:
		return this.BuildPrintStatementNode(astNode)
	case dsl.NodeTypeEprintStatement:
		return this.BuildEprintStatementNode(astNode)
	case dsl.NodeTypePrintnStatement:
		return this.BuildPrintnStatementNode(astNode)
	case dsl.NodeTypeEprintnStatement:
		return this.BuildEprintnStatementNode(astNode)

	case dsl.NodeTypeBeginBlock:
		return nil, errors.New(
			"Miller: begin blocks may only be declared at top level.",
		)
	case dsl.NodeTypeEndBlock:
		return nil, errors.New(
			"Miller: end blocks may only be declared at top level.",
		)

	case dsl.NodeTypeIfChain:
		return this.BuildIfChainNode(astNode)
	case dsl.NodeTypeCondBlock:
		return this.BuildCondBlockNode(astNode)
	case dsl.NodeTypeWhileLoop:
		return this.BuildWhileLoopNode(astNode)
	case dsl.NodeTypeDoWhileLoop:
		return this.BuildDoWhileLoopNode(astNode)
	case dsl.NodeTypeForLoopOneVariable:
		return this.BuildForLoopOneVariableNode(astNode)
	case dsl.NodeTypeForLoopTwoVariable:
		return this.BuildForLoopTwoVariableNode(astNode)
	case dsl.NodeTypeTripleForLoop:
		return this.BuildTripleForLoopNode(astNode)

	case dsl.NodeTypeFunctionDefinition:
		return nil, errors.New(
			"Miller: functions may only be declared at top level.",
		)

	case dsl.NodeTypeBreak:
		return this.BuildBreakNode(astNode)
	case dsl.NodeTypeContinue:
		return this.BuildContinueNode(astNode)
	case dsl.NodeTypeReturn:
		return this.BuildReturnNode(astNode)

	default:
		return nil, errors.New(
			"CST BuildStatementNode: unhandled AST node " + string(astNode.Type),
		)
		break
	}
	return statement, nil
}
