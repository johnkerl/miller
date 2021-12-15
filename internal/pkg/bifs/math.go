package bifs

import (
	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

// Return error (unary math-library func)
func _math_unary_erro1(input1 *mlrval.Mlrval, f mathLibUnaryFunc) *mlrval.Mlrval {
	return mlrval.ERROR
}

// Return absent (unary math-library func)
func _math_unary_absn1(input1 *mlrval.Mlrval, f mathLibUnaryFunc) *mlrval.Mlrval {
	return mlrval.ABSENT
}

// Return null (unary math-library func)
func _math_unary_null1(input1 *mlrval.Mlrval, f mathLibUnaryFunc) *mlrval.Mlrval {
	return mlrval.NULL
}

// Return void (unary math-library func)
func _math_unary_void1(input1 *mlrval.Mlrval, f mathLibUnaryFunc) *mlrval.Mlrval {
	return mlrval.VOID
}
