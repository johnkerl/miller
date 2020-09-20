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

	case dsl.NodeTypeBeginBlock:
		return nil, nil // xxx temp

	case dsl.NodeTypeEndBlock:
		return nil, nil // xxx temp

	case dsl.NodeTypeFilterStatement:
		return BuildFilterStatementNode(astNode)

	default:
		return nil, errors.New(
			"CST BuildStatementNode: unhandled AST node " + string(astNode.Type),
		)
		break
	}
	return statement, nil
}
