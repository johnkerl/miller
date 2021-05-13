package types

// Return error (unary math-library func)
func _math_unary_erro1(input1 *Mlrval, f mathLibUnaryFunc) *Mlrval {
	return MLRVAL_ERROR
}

// Return absent (unary math-library func)
func _math_unary_absn1(input1 *Mlrval, f mathLibUnaryFunc) *Mlrval {
	return MLRVAL_ABSENT
}

// Return null (unary math-library func)
func _math_unary_null1(input1 *Mlrval, f mathLibUnaryFunc) *Mlrval {
	return MLRVAL_NULL
}

// Return void (unary math-library func)
func _math_unary_void1(input1 *Mlrval, f mathLibUnaryFunc) *Mlrval {
	return MLRVAL_VOID
}
