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

	default:
		return nil, errors.New("Non-assignment AST node unhandled")
		break
	}
	return statement, nil
}
