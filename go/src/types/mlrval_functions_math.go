package types

// Return error (unary math-library func)
func _math_unary_erro1(output, input1 *Mlrval, f mathLibUnaryFunc) {
	output.SetFromError()
}

// Return absent (unary math-library func)
func _math_unary_absn1(output, input1 *Mlrval, f mathLibUnaryFunc) {
	output.SetFromAbsent()
}

// Return void (unary math-library func)
func _math_unary_void1(output, input1 *Mlrval, f mathLibUnaryFunc) {
	output.SetFromAbsent()
}
