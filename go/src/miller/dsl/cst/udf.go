package cst

// ================================================================
// Support for user-defined functions and subroutines
// ================================================================

type Signature struct {
	functionName   string
	parameterNames []string
	// todo: parameter typedecls
	// todo: return-value typedecls
}
