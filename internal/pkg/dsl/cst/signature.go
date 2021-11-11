// ================================================================
// Signatures for user-defined functions and user-defined subroutines
// ("UDFs" and "UDSs").
// ================================================================

package cst

import (
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
type Signature struct {
	funcOrSubrName          string
	arity                   int // Computable from len(typeGatedParameterNames) at callee, not at caller
	typeGatedParameterNames []*types.TypeGatedMlrvalName
	typeGatedReturnValue    *types.TypeGatedMlrvalName
}

func NewSignature(
	funcOrSubrName string,
	arity int,
	typeGatedParameterNames []*types.TypeGatedMlrvalName,
	typeGatedReturnValue *types.TypeGatedMlrvalName,
) *Signature {
	return &Signature{
		funcOrSubrName:          funcOrSubrName,
		arity:                   arity,
		typeGatedParameterNames: typeGatedParameterNames,
		typeGatedReturnValue:    typeGatedReturnValue,
	}
}
