package cst

import (
	"errors"
	"fmt"

	"miller/dsl"
)

// ================================================================
// Checks for things that are syntax errors but not done in the AST for
// pragmatic reasons. For example, $anything in begin/end blocks;
// begin/end/func not at top level; etc.
// ================================================================

// ----------------------------------------------------------------
func ValidateAST(ast *dsl.AST) error {
	atTopLevel := true
	inLoop := false
	inBeginOrEnd := false
	inUDF := false
	inUDS := false
	isMainBlockLastStatement := false
	isAssignmentLHS := false

	// They can do mlr put '': there are simply zero statements.
	if ast.RootNode.Type == dsl.NodeTypeEmptyStatement {
		return nil
	}

	if ast.RootNode.Children != nil {
		for _, astChild := range ast.RootNode.Children {
			return validateASTAux(
				astChild,
				atTopLevel,
				inLoop,
				inBeginOrEnd,
				inUDF,
				inUDS,
				isMainBlockLastStatement,
				isAssignmentLHS,
			)
		}
	}

	return nil
}

// ----------------------------------------------------------------
func validateASTAux(
	astNode *dsl.ASTNode,
	atTopLevel bool,
	inLoop bool,
	inBeginOrEnd bool,
	inUDF bool,
	inUDS bool,
	isMainBlockLastStatement bool, // TODO -- keep this or not ...
	isAssignmentLHS bool,
) error {
	nextAtTopLevel := false
	nextInLoop := inLoop
	nextInBeginOrEnd := inBeginOrEnd
	nextInUDF := inUDF
	nextInUDS := inUDS
	nextIsAssignmentLHS := isAssignmentLHS

	// Check: begin/end/func/subr must be at top-level
	if astNode.Type == dsl.NodeTypeBeginBlock {
		if !atTopLevel {
			return errors.New(
				"Miller: begin blocks can only be at top level.",
			)
		}
		nextInBeginOrEnd = true
	}
	if astNode.Type == dsl.NodeTypeEndBlock {
		if !atTopLevel {
			return errors.New(
				"Miller: end blocks can only be at top level.",
			)
		}
		nextInBeginOrEnd = true
	}
	if astNode.Type == dsl.NodeTypeFunctionDefinition {
		if !atTopLevel {
			return errors.New(
				"Miller: func blocks can only be at top level.",
			)
		}
		nextInUDF = true
	}
	if astNode.Type == dsl.NodeTypeSubroutineDefinition {
		if !atTopLevel {
			return errors.New(
				"Miller: subr blocks can only be at top level.",
			)
		}
		nextInUDS = true
	}

	// Check: $-anything cannot be in begin/end
	if inBeginOrEnd {
		if astNode.Type == dsl.NodeTypeDirectFieldValue ||
			astNode.Type == dsl.NodeTypeIndirectFieldValue ||
			astNode.Type == dsl.NodeTypeFullSrec {
			return errors.New(
				"Miller: begin/end blocks cannot refer to records via $x, $*, etc.",
			)
		}
	}

	// Check: filter statements cannot be in begin/end
	if inBeginOrEnd {
		if astNode.Type == dsl.NodeTypeBareBoolean {
			return errors.New(
				"Miller: begin/end blocks cannot have bare booleans: there is no record to filter",
			)
		}
	}

	// Check: break/continue outside of loop
	if !inLoop {
		if astNode.Type == dsl.NodeTypeBreak {
			return errors.New(
				"Miller: break statements are only valid within for/do/while loops.",
			)
		}
	}

	if !inLoop {
		if astNode.Type == dsl.NodeTypeContinue {
			return errors.New(
				"Miller: break statements are only valid within for/do/while loops.",
			)
		}
	}

	if astNode.Type == dsl.NodeTypeWhileLoop ||
		astNode.Type == dsl.NodeTypeDoWhileLoop ||
		astNode.Type == dsl.NodeTypeForLoopKeyOnly ||
		astNode.Type == dsl.NodeTypeForLoopKeyValue ||
		astNode.Type == dsl.NodeTypeTripleForLoop {
		nextInLoop = true
	}

	// Check: return outside of func/subr
	if !inUDF && !inUDS {
		if astNode.Type == dsl.NodeTypeReturn {
			return errors.New(
				"Miller: return statements are only valid within func/subr blocks.",
			)
		}
	}

	// Check: enforce return-value iff in a function; return-void iff in a subroutine
	if astNode.Type == dsl.NodeTypeReturn {
		if inUDF {
			if len(astNode.Children) != 1 {
				return errors.New(
					"Miller: return statements in func blocks must return a value.",
				)
			}
		}
		if inUDS {
			if len(astNode.Children) != 0 {
				return errors.New(
					"Miller: return statements in subr blocks must not return a value.",
				)
			}
		}
	}

	// Check: prohibit NR etc at LHS; 1+2=3+4; etc
	if isAssignmentLHS {
		ok := VALID_LHS_NODE_TYPES[astNode.Type]
		if !ok {
			return errors.New(
				fmt.Sprintf(
					"Miller: %s is not valid on the left-hand side of an assignment.",
					astNode.Type,
				),
			)
		}
	}

	// Check: ENV disallowed at left-hand side of assignment? or nah?
	// Needs to be supported for system() ...
	// * TODO

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Treewalk

	if astNode.Children != nil {
		for i, astChild := range astNode.Children {
			nextIsAssignmentLHS = astNode.Type == dsl.NodeTypeAssignment && i == 0
			err := validateASTAux(
				astChild,
				nextAtTopLevel,
				nextInLoop,
				nextInBeginOrEnd,
				nextInUDF,
				nextInUDS,
				isMainBlockLastStatement,
				nextIsAssignmentLHS,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

var VALID_LHS_NODE_TYPES = map[dsl.TNodeType]bool{
	dsl.NodeTypeArrayOrMapIndexAccess: true,
	dsl.NodeTypeArraySliceAccess:      true,
	dsl.NodeTypeDirectFieldValue:      true,
	dsl.NodeTypeIndirectFieldValue:    true,
	dsl.NodeTypeFullSrec:              true,
	dsl.NodeTypeDirectOosvarValue:     true,
	dsl.NodeTypeIndirectOosvarValue:   true,
	dsl.NodeTypeFullOosvar:            true,
	dsl.NodeTypeLocalVariable:         true,
}
