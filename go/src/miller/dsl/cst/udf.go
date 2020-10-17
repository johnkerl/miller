package cst

import (
	"miller/types"
)

// ================================================================
// Support for user-defined functions and subroutines
// ================================================================

// ----------------------------------------------------------------
type Signature struct {
	functionName   string
	arity          int // Computable from len(parameterNames) at callee, not at caller
	parameterNames []string

	// TODO: parameter typedecls
	// TODO: return-value typedecls
}

// ----------------------------------------------------------------
type UDFManager struct {
}

func NewUDFManager() *UDFManager {
	return &UDFManager{}
}

func (this *UDFManager) LookUp(functionName string, callsiteArity int) *UDF {
	return nil
}

// ----------------------------------------------------------------
type UDF struct {
	signature Signature
}

func (this *UDF) BuildEvaluableNode() IEvaluable {
	return nil
}

func (this *UDF) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromInt64(999)
}
