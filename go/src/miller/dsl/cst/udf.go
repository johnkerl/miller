package cst

// ================================================================
// Support for user-defined functions and subroutines
// ================================================================

type Signature struct {
	functionName   string
	arity          int // Computable from len(parameterNames) at callee, not at caller
	parameterNames []string

	// TODO: parameter typedecls
	// TODO: return-value typedecls
}
