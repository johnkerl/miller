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
func BuildStatementNode(
	astNode *dsl.ASTNode,
) (IExecutable, error) {

	var statement IExecutable = nil
	var err error = nil
	// xxx more to do
	switch astNode.Type {

	case dsl.NodeTypeAssignment:
		statement, err = BuildAssignmentNode(astNode)
		if err != nil {
			return nil, err
		}

	case dsl.NodeTypeFilterStatement:
		return BuildFilterStatementNode(astNode)
	case dsl.NodeTypeBareBoolean:
		return BuildFilterStatementNode(astNode)

	case dsl.NodeTypeEmitStatement:
		return BuildEmitStatementNode(astNode)
	case dsl.NodeTypeDumpStatement:
		return BuildDumpStatementNode(astNode)
	case dsl.NodeTypeEdumpStatement:
		return BuildEdumpStatementNode(astNode)

	case dsl.NodeTypeBeginBlock:
		return nil, nil // xxx temp -- only valid at top level; say so here w/ error
	case dsl.NodeTypeEndBlock:
		return nil, nil // xxx temp -- only valid at top level; say so here w/ error

	case dsl.NodeTypeIfChain:
		return BuildIfChainNode(astNode)
	case dsl.NodeTypeCondBlock:
		return BuildCondBlockNode(astNode)
	case dsl.NodeTypeWhileLoop:
		return BuildWhileLoopNode(astNode)
	case dsl.NodeTypeDoWhileLoop:
		return BuildDoWhileLoopNode(astNode)
	case dsl.NodeTypeForLoopKeyOnly:
		return BuildForLoopKeyOnlyNode(astNode)
	case dsl.NodeTypeForLoopKeyValue:
		return BuildForLoopKeyValueNode(astNode)
	case dsl.NodeTypeTripleForLoop:
		return BuildTripleForLoopNode(astNode)

	default:
		return nil, errors.New(
			"CST BuildStatementNode: unhandled AST node " + string(astNode.Type),
		)
		break
	}
	return statement, nil
}
