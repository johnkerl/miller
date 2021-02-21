package types

// Return error (unary math-library func)
func _math_unary_erro1(input1 *Mlrval, f mathLibUnaryFunc) Mlrval {
	return MlrvalFromError()
}

// Return absent (unary math-library func)
func _math_unary_absn1(input1 *Mlrval, f mathLibUnaryFunc) Mlrval {
	return MlrvalFromAbsent()
}

// Return void (unary math-library func)
func _math_unary_void1(input1 *Mlrval, f mathLibUnaryFunc) Mlrval {
	return MlrvalFromAbsent()
}
