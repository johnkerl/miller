// ================================================================
// Checks for things that are syntax errors but not done in the AST for
// pragmatic reasons. For example, $anything in begin/end blocks;
// begin/end/func not at top level; etc.
// ================================================================

package cst

import (
	"errors"
	"fmt"

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
)

// ----------------------------------------------------------------
func ValidateAST(
	ast *dsl.AST,
	dslInstanceType DSLInstanceType, // mlr put, mlr filter, mlr repl
) error {
	atTopLevel := true
	inLoop := false
	inBeginOrEnd := false
	inUDF := false
	inUDS := false
	isMainBlockLastStatement := false
	isAssignmentLHS := false
	isUnset := false

	// They can do mlr put '': there are simply zero statements.
	// But filter '' is an error.
	if ast.RootNode.Children == nil || len(ast.RootNode.Children) == 0 {
		if dslInstanceType == DSLInstanceTypeFilter {
			return errors.New("mlr: filter statement must not be empty.")
		}
	}

	if ast.RootNode.Children != nil {
		for _, astChild := range ast.RootNode.Children {
			err := validateASTAux(
				astChild,
				dslInstanceType,
				atTopLevel,
				inLoop,
				inBeginOrEnd,
				inUDF,
				inUDS,
				isMainBlockLastStatement,
				isAssignmentLHS,
				isUnset,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ----------------------------------------------------------------
func validateASTAux(
	astNode *dsl.ASTNode,
	dslInstanceType DSLInstanceType, // mlr put, mlr filter, mlr repl
	atTopLevel bool,
	inLoop bool,
	inBeginOrEnd bool,
	inUDF bool,
	inUDS bool,
	isMainBlockLastStatement bool, // TODO -- keep this or not ...
	isAssignmentLHS bool,
	isUnset bool,
) error {
	nextLevelAtTopLevel := false
	nextLevelInLoop := inLoop
	nextLevelInBeginOrEnd := inBeginOrEnd
	nextLevelInUDF := inUDF
	nextLevelInUDS := inUDS
	nextLevelIsAssignmentLHS := isAssignmentLHS
	nextLevelIsUnset := isUnset

	if astNode.Type == dsl.NodeTypeFilterStatement {
		if dslInstanceType == DSLInstanceTypeFilter {
			return errors.New(
				"mlr: filter expressions must not also contain the \"filter\" keyword.",
			)
		}
	}

	// Check: begin/end/func/subr must be at top-level
	if astNode.Type == dsl.NodeTypeBeginBlock {
		if !atTopLevel {
			return errors.New(
				"mlr: begin blocks can only be at top level.",
			)
		}
		nextLevelInBeginOrEnd = true
	} else if astNode.Type == dsl.NodeTypeEndBlock {
		if !atTopLevel {
			return errors.New(
				"mlr: end blocks can only be at top level.",
			)
		}
		nextLevelInBeginOrEnd = true
	} else if astNode.Type == dsl.NodeTypeNamedFunctionDefinition {
		if !atTopLevel {
			return errors.New(
				"mlr: func blocks can only be at top level.",
			)
		}
		nextLevelInUDF = true
	} else if astNode.Type == dsl.NodeTypeUnnamedFunctionDefinition {
		nextLevelInUDF = true
	} else if astNode.Type == dsl.NodeTypeSubroutineDefinition {
		if !atTopLevel {
			return errors.New(
				"mlr: subr blocks can only be at top level.",
			)
		}
		nextLevelInUDS = true
	} else if astNode.Type == dsl.NodeTypeForLoopTwoVariable {
		err := validateForLoopTwoVariableUniqueNames(astNode)
		if err != nil {
			return err
		}
	} else if astNode.Type == dsl.NodeTypeForLoopMultivariable {
		err := validateForLoopMultivariableUniqueNames(astNode)
		if err != nil {
			return err
		}
	}

	// Check: $-anything cannot be in begin/end
	if inBeginOrEnd {
		if astNode.Type == dsl.NodeTypeDirectFieldValue ||
			astNode.Type == dsl.NodeTypeIndirectFieldValue ||
			astNode.Type == dsl.NodeTypeFullSrec {
			return errors.New(
				"mlr: begin/end blocks cannot refer to records via $x, $*, etc.",
			)
		}
	}

	// Check: break/continue outside of loop
	if !inLoop {
		if astNode.Type == dsl.NodeTypeBreak {
			return errors.New(
				"mlr: break statements are only valid within for/do/while loops.",
			)
		}
	}

	if !inLoop {
		if astNode.Type == dsl.NodeTypeContinue {
			return errors.New(
				"mlr: break statements are only valid within for/do/while loops.",
			)
		}
	}

	if astNode.Type == dsl.NodeTypeWhileLoop ||
		astNode.Type == dsl.NodeTypeDoWhileLoop ||
		astNode.Type == dsl.NodeTypeForLoopOneVariable ||
		astNode.Type == dsl.NodeTypeForLoopTwoVariable ||
		astNode.Type == dsl.NodeTypeForLoopMultivariable ||
		astNode.Type == dsl.NodeTypeTripleForLoop {
		nextLevelInLoop = true
	}

	// Check: return outside of func/subr
	if !inUDF && !inUDS {
		if astNode.Type == dsl.NodeTypeReturn {
			return errors.New(
				"mlr: return statements are only valid within func/subr blocks.",
			)
		}
	}

	// Check: enforce return-value iff in a function; return-void iff in a subroutine
	if astNode.Type == dsl.NodeTypeReturn {
		if inUDF {
			if len(astNode.Children) != 1 {
				return errors.New(
					"mlr: return statements in func blocks must return a value.",
				)
			}
		}
		if inUDS {
			if len(astNode.Children) != 0 {
				return errors.New(
					"mlr: return statements in subr blocks must not return a value.",
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
					"mlr: %s is not valid on the left-hand side of an assignment.",
					astNode.Type,
				),
			)
		}
	}

	// Check: prohibit NR etc at LHS; 1+2=3+4; etc
	if isUnset {
		ok := VALID_LHS_NODE_TYPES[astNode.Type]
		if !ok {
			return errors.New(
				fmt.Sprintf(
					"mlr: %s is not valid for unset statement.",
					astNode.Type,
				),
			)
		}
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Treewalk

	if astNode.Children != nil {
		for i, astChild := range astNode.Children {
			nextLevelIsAssignmentLHS = astNode.Type == dsl.NodeTypeAssignment && i == 0
			nextLevelIsUnset = astNode.Type == dsl.NodeTypeUnset
			err := validateASTAux(
				astChild,
				dslInstanceType,
				nextLevelAtTopLevel,
				nextLevelInLoop,
				nextLevelInBeginOrEnd,
				nextLevelInUDF,
				nextLevelInUDS,
				isMainBlockLastStatement,
				nextLevelIsAssignmentLHS,
				nextLevelIsUnset,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Check against 'for (a, a in $*)' -- repeated 'a'.
// AST:
// * statement block
//   * double-variable for-loop "for"
//     * local variable "a"
//     * local variable "a"
//     * full record "$*"
//     * statement block

func validateForLoopTwoVariableUniqueNames(astNode *dsl.ASTNode) error {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeForLoopTwoVariable)
	lib.InternalCodingErrorIf(len(astNode.Children) != 4)
	keyVarNode := astNode.Children[0]
	valVarNode := astNode.Children[1]
	lib.InternalCodingErrorIf(keyVarNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(valVarNode.Type != dsl.NodeTypeLocalVariable)
	keyVarName := string(keyVarNode.Token.Lit)
	valVarName := string(valVarNode.Token.Lit)
	if keyVarName == valVarName {
		return errors.New(
			fmt.Sprintf(
				"%s: redefinition of variable %s in the same scope.",
				"mlr",
				keyVarName,
			),
		)
	} else {
		return nil
	}
}

// Check against 'for ((a,a), b in $*)' or 'for ((a,b), a in $*)' -- repeated 'a'.
// AST:
// * statement block
//   * multi-variable for-loop "for"
//     * parameter list
//       * local variable "a"
//       * local variable "b"
//     * local variable "a"
//     * full record "$*"
//     * statement block
func validateForLoopMultivariableUniqueNames(astNode *dsl.ASTNode) error {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeForLoopMultivariable)
	keyVarsNode := astNode.Children[0]
	valVarNode := astNode.Children[1]
	lib.InternalCodingErrorIf(keyVarsNode.Type != dsl.NodeTypeParameterList)
	lib.InternalCodingErrorIf(valVarNode.Type != dsl.NodeTypeLocalVariable)

	seen := make(map[string]bool)

	for _, keyVarNode := range keyVarsNode.Children {
		lib.InternalCodingErrorIf(keyVarNode.Type != dsl.NodeTypeLocalVariable)
		name := string(keyVarNode.Token.Lit)
		_, present := seen[name]
		if present {
			return errors.New(
				fmt.Sprintf(
					"%s: redefinition of variable %s in the same scope.",
					"mlr",
					name,
				),
			)
		}
		seen[name] = true
	}

	valVarName := string(valVarNode.Token.Lit)
	if seen[valVarName] {
		return errors.New(
			fmt.Sprintf(
				"%s: redefinition of variable %s in the same scope.",
				"mlr",
				valVarName,
			),
		)
	}

	return nil
}

// ================================================================

var VALID_LHS_NODE_TYPES = map[dsl.TNodeType]bool{
	dsl.NodeTypeArrayOrMapIndexAccess:           true,
	dsl.NodeTypeDotOperator:                     true,
	dsl.NodeTypeArrayOrMapPositionalNameAccess:  true,
	dsl.NodeTypeArrayOrMapPositionalValueAccess: true,
	dsl.NodeTypeArraySliceAccess:                true,
	dsl.NodeTypeDirectFieldValue:                true,
	dsl.NodeTypeDirectOosvarValue:               true,
	dsl.NodeTypeEnvironmentVariable:             true,
	dsl.NodeTypeFullOosvar:                      true,
	dsl.NodeTypeFullSrec:                        true,
	dsl.NodeTypeIndirectFieldValue:              true,
	dsl.NodeTypeIndirectOosvarValue:             true,
	dsl.NodeTypeLocalVariable:                   true,
	dsl.NodeTypePositionalFieldName:             true,
	dsl.NodeTypePositionalFieldValue:            true,
}
