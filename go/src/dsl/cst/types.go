// ================================================================
// Main type definitions for CST build/execute
// ================================================================

package cst

import (
	"container/list"

	"mlr/src/cliutil"
	"mlr/src/dsl"
	"mlr/src/runtime"
	"mlr/src/types"
)

// ----------------------------------------------------------------
// Please see root.go for context and comments.
type RootNode struct {
	beginBlocks                   []*StatementBlockNode
	mainBlock                     *StatementBlockNode
	replImmediateBlock            *StatementBlockNode
	endBlocks                     []*StatementBlockNode
	udfManager                    *UDFManager
	udsManager                    *UDSManager
	allowUDFUDSRedefinitions      bool
	unresolvedFunctionCallsites   *list.List
	unresolvedSubroutineCallsites *list.List
	outputHandlerManagers         *list.List
	recordWriterOptions           *cliutil.TWriterOptions
}

// ----------------------------------------------------------------
// Many functions have this signature. This type-alias is for function-name
// lookup tables.
type NodeBuilder func(astNode *dsl.ASTNode) (IEvaluable, error)

// ----------------------------------------------------------------
// This is for all statements and statemnt blocks within the CST.
type IExecutable interface {
	Execute(state *runtime.State) (*BlockExitPayload, error)
}

type Executor func(state *runtime.State) (*BlockExitPayload, error)

// ================================================================
// This is for any left-hand side (LHS or Lvalue) of an assignment statement.
type IAssignable interface {
	Assign(rvalue *types.Mlrval, state *runtime.State) error

	// 'foo = "bar"' or 'foo[3]["abc"] = "bar"'
	// For non-indexed assignment, which is the normal case, indices can be
	// zero-length or nil.
	AssignIndexed(rvalue *types.Mlrval, indices []*types.Mlrval, state *runtime.State) error

	Unassign(state *runtime.State)

	UnassignIndexed(indices []*types.Mlrval, state *runtime.State)
}

// ================================================================
// This is for any right-hand side (RHS or Rvalue) of an assignment statement.
// Also, for computed field names on the left-hand side, like '$a . $b' in mlr
// put '$[$a . $b]' = $x + $y'.
type IEvaluable interface {
	Evaluate(state *runtime.State) *types.Mlrval
}

// ================================================================
// For blocks of statements: the main put/filter block; begin/end blocks;
// for/while-loop bodies; user-defined functions/subroutines.

// ----------------------------------------------------------------
// Also implements IExecutable
type StatementBlockNode struct {
	executables []IExecutable
}

// ----------------------------------------------------------------
// Things a block of statements can do:
// * execute all the way to the end without a return
// * break
// * continue
// * (throw an exception if the Miller DSL were to support that)
// * return void
// * return a value
type BlockExitStatus int

const (
	// BLOCK_EXIT_RUN_TO_END is implemented as *BlockExitPayload being nil
	BLOCK_EXIT_BREAK        BlockExitStatus = 1
	BLOCK_EXIT_CONTINUE                     = 2
	BLOCK_EXIT_RETURN_VOID                  = 3
	BLOCK_EXIT_RETURN_VALUE                 = 4
)

type BlockExitPayload struct {
	blockExitStatus BlockExitStatus
	// No multiple return yet in the Miller DSL -- if there were, this would be
	// an array.
	blockReturnValue *types.Mlrval // TODO: TypeGatedMlrval
}
