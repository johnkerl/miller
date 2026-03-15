// Checks for things that are syntax errors but not done in the AST for
// pragmatic reasons. For example, $anything in begin/end blocks;
// begin/end/func not at top level; etc.

package cst

import (
	"fmt"

	"github.com/johnkerl/miller/v6/pkg/lib"

	"github.com/johnkerl/pgpg/go/lib/pkg/asts"
)

func ValidateAST(
	ast *asts.AST,
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
	if len(ast.RootNode.Children) == 0 {
		if dslInstanceType == DSLInstanceTypeFilter {
			return fmt.Errorf("filter statement must not be empty")
		}
	}

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

	return nil
}

func validateASTAux(
	astNode *asts.ASTNode,
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

	if astNode.Type == asts.NodeType(NodeTypeFilterStatement) {
		if dslInstanceType == DSLInstanceTypeFilter {
			return fmt.Errorf(
				`filter expressions must not also contain the "filter" keyword`,
			)
		}
	}

	// Check: begin/end/func/subr must be at top-level
	if astNode.Type == asts.NodeType(NodeTypeBeginBlock) {
		if !atTopLevel {
			return fmt.Errorf(
				"begin blocks can only be at top level",
			)
		}
		nextLevelInBeginOrEnd = true
	} else if astNode.Type == asts.NodeType(NodeTypeEndBlock) {
		if !atTopLevel {
			return fmt.Errorf(
				"end blocks can only be at top level",
			)
		}
		nextLevelInBeginOrEnd = true
	} else if astNode.Type == asts.NodeType(NodeTypeNamedFunctionDefinition) {
		if !atTopLevel {
			return fmt.Errorf(
				"func blocks can only be at top level",
			)
		}
		nextLevelInUDF = true
	} else if astNode.Type == asts.NodeType(NodeTypeUnnamedFunctionDefinition) {
		nextLevelInUDF = true
	} else if astNode.Type == asts.NodeType(NodeTypeSubroutineDefinition) {
		if !atTopLevel {
			return fmt.Errorf(
				"subr blocks can only be at top level",
			)
		}
		nextLevelInUDS = true
	} else if astNode.Type == asts.NodeType(NodeTypeForLoopTwoVariable) {
		err := validateForLoopTwoVariableUniqueNames(astNode)
		if err != nil {
			return err
		}
	} else if astNode.Type == asts.NodeType(NodeTypeForLoopMultivariable) {
		err := validateForLoopMultivariableUniqueNames(astNode)
		if err != nil {
			return err
		}
	}

	// Check: $-anything cannot be in begin/end
	if inBeginOrEnd {
		if astNode.Type == asts.NodeType(NodeTypeDirectFieldValue) ||
			astNode.Type == asts.NodeType(NodeTypeIndirectFieldValue) ||
			astNode.Type == asts.NodeType(NodeTypeFullSrec) {
			return fmt.Errorf(
				"begin/end blocks cannot refer to records via $x, $*, etc",
			)
		}
	}

	// Check: break/continue outside of loop
	if !inLoop {
		if astNode.Type == asts.NodeType(NodeTypeBreakStatement) {
			return fmt.Errorf(
				"break statements are only valid within for/do/while loops",
			)
		}
	}

	if !inLoop {
		if astNode.Type == asts.NodeType(NodeTypeContinueStatement) {
			return fmt.Errorf(
				"break statements are only valid within for/do/while loops",
			)
		}
	}

	if astNode.Type == asts.NodeType(NodeTypeWhileLoop) ||
		astNode.Type == asts.NodeType(NodeTypeDoWhileLoop) ||
		astNode.Type == asts.NodeType(NodeTypeForLoopOneVariable) ||
		astNode.Type == asts.NodeType(NodeTypeForLoopTwoVariable) ||
		astNode.Type == asts.NodeType(NodeTypeForLoopMultivariable) ||
		astNode.Type == asts.NodeType(NodeTypeTripleForLoop) {
		nextLevelInLoop = true
	}

	// Check: return outside of func/subr
	if !inUDF && !inUDS {
		if astNode.Type == asts.NodeType(NodeTypeReturnStatement) {
			return fmt.Errorf(
				"return statements are only valid within func/subr blocks",
			)
		}
	}

	// Check: enforce return-value iff in a function; return-void iff in a subroutine
	if astNode.Type == asts.NodeType(NodeTypeReturnStatement) {
		if inUDF {
			if len(astNode.Children) != 1 {
				return fmt.Errorf(
					"return statements in func blocks must return a value",
				)
			}
		}
		if inUDS {
			if len(astNode.Children) != 0 {
				return fmt.Errorf(
					"return statements in subr blocks must not return a value",
				)
			}
		}
	}

	// Check: prohibit NR etc at LHS; 1+2=3+4; etc
	if isAssignmentLHS {
		ok := VALID_LHS_NODE_TYPES[string(astNode.Type)]
		if !ok {
			return fmt.Errorf(
				"%s is not valid on the left-hand side of an assignment",
				astNode.Type,
			)
		}
	}

	// Check: prohibit NR etc at LHS; 1+2=3+4; etc
	if isUnset {
		ok := VALID_LHS_NODE_TYPES[string(astNode.Type)]
		if !ok {
			return fmt.Errorf(
				"%s is not valid for unset statement",
				astNode.Type,
			)
		}
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Treewalk

	for i, astChild := range astNode.Children {
		nextLevelIsAssignmentLHS = (astNode.Type == asts.NodeType(NodeTypeAssignment) || astNode.Type == asts.NodeType(NodeTypeCompoundAssignment)) && i == 0
		nextLevelIsUnset = astNode.Type == asts.NodeType(NodeTypeUnset)
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

func validateForLoopTwoVariableUniqueNames(astNode *asts.ASTNode) error {
	lib.InternalCodingErrorIf(astNode.Type != asts.NodeType(NodeTypeForLoopTwoVariable))
	lib.InternalCodingErrorIf(len(astNode.Children) != 4)
	keyVarNode := astNode.Children[0]
	valVarNode := astNode.Children[1]
	lib.InternalCodingErrorIf(keyVarNode.Type != asts.NodeType(NodeTypeLocalVariable))
	lib.InternalCodingErrorIf(valVarNode.Type != asts.NodeType(NodeTypeLocalVariable))
	keyVarName := tokenLit(keyVarNode)
	valVarName := tokenLit(valVarNode)
	if keyVarName == valVarName {
		return fmt.Errorf("redefinition of variable %s in the same scope", keyVarName)
	}
	return nil
}

// Check against 'for ((a,a), b in $*)' or 'for ((a,b), a in $*)' -- repeated 'a'.
// AST:
// * statement block
//   - multi-variable for-loop "for"
//   - parameter list
//   - local variable "a"
//   - local variable "b"
//   - local variable "a"
//   - full record "$*"
//   - statement block
func validateForLoopMultivariableUniqueNames(astNode *asts.ASTNode) error {
	lib.InternalCodingErrorIf(astNode.Type != asts.NodeType(NodeTypeForLoopMultivariable))
	keyVarsNode := astNode.Children[0]
	valVarNode := astNode.Children[1]
	// PGPG produces MultiIndex; legacy produced ParameterList. Both have LocalVariable children.
	lib.InternalCodingErrorIf(keyVarsNode.Type != asts.NodeType(NodeTypeParameterList) &&
		keyVarsNode.Type != asts.NodeType(NodeTypeMultiIndex))
	lib.InternalCodingErrorIf(valVarNode.Type != asts.NodeType(NodeTypeLocalVariable))

	seen := make(map[string]bool)

	for _, keyVarNode := range keyVarsNode.Children {
		lib.InternalCodingErrorIf(keyVarNode.Type != asts.NodeType(NodeTypeLocalVariable))
		name := tokenLit(keyVarNode)
		_, present := seen[name]
		if present {
			return fmt.Errorf("redefinition of variable %s in the same scope", name)
		}
		seen[name] = true
	}

	valVarName := tokenLit(valVarNode)
	if seen[valVarName] {
		return fmt.Errorf("redefinition of variable %s in the same scope", valVarName)
	}

	return nil
}

var VALID_LHS_NODE_TYPES = map[string]bool{
	NodeTypeArrayOrMapIndexAccess: true,
	NodeTypeDotOperator:           true,
	NodeTypeArraySliceLoHi:        true,
	NodeTypeArraySliceHiOnly:      true,
	NodeTypeArraySliceLoOnly:      true,
	NodeTypeArraySliceFull:        true,
	NodeTypeDirectFieldValue:      true,
	NodeTypeBracedFieldValue:      true, // ${foo}, ${x+y}
	NodeTypeEnvironmentVariable:   true, // ENV["FOO"] = "bar"
	NodeTypeDirectOosvarValue:     true,
	NodeTypeBracedOosvarValue:     true, // @{foo}, @{variable.name}
	NodeTypeFullOosvar:            true,
	NodeTypeFullSrec:              true,
	NodeTypeIndirectFieldValue:    true, // includes $[[n]] and $[[[n]]]
	NodeTypeIndirectOosvarValue:   true,
	NodeTypeLocalVariable:         true,
	"TypedeclLocalVariable":       true, // int a = 5; etc
}
