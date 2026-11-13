package bifs

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

func TestBIF_plus_unary(t *testing.T) {
	input := mlrval.FromDeferredType("123")
	output := BIF_plus_unary(input)
	intval, ok := output.GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, int64(123), intval)

	input = mlrval.FromDeferredType("-123.5")
	output = BIF_plus_unary(input)
	floatval, ok := output.GetFloatValue()
	assert.True(t, ok)
	assert.Equal(t, -123.5, floatval)
}

func TestBIF_minus_unary(t *testing.T) {
	input := mlrval.FromDeferredType("123")
	output := BIF_minus_unary(input)
	intval, ok := output.GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, int64(-123), intval)

	input = mlrval.FromDeferredType("-123.5")
	output = BIF_minus_unary(input)
	floatval, ok := output.GetFloatValue()
	assert.True(t, ok)
	assert.Equal(t, 123.5, floatval)
}

func TestBIF_plus_binary(t *testing.T) {
	input1 := mlrval.FromDeferredType("123")
	input2 := mlrval.FromDeferredType("456")
	output := BIF_plus_binary(input1, input2)
	intval, ok := output.GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, int64(579), intval)

	input1 = mlrval.FromDeferredType("123.5")
	input2 = mlrval.FromDeferredType("456")
	output = BIF_plus_binary(input1, input2)
	floatval, ok := output.GetFloatValue()
	assert.True(t, ok)
	assert.Equal(t, 579.5, floatval)
}

func TestBIF_plus_binary_overflow(t *testing.T) {
	input1 := mlrval.FromInt(0x07ffffffffffffff)
	input2 := mlrval.FromInt(0x07fffffffffffffe)
	output := BIF_plus_binary(input1, input2)
	intval, ok := output.GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, int64(0x0ffffffffffffffd), intval)

	input1 = mlrval.FromInt(0x7fffffffffffffff)
	input2 = mlrval.FromInt(0x7ffffffffffffffe)
	output = BIF_plus_binary(input1, input2)
	floatval, ok := output.GetFloatValue()
	assert.True(t, ok)
	assert.Equal(t, 18446744073709552000.0, floatval)
}

func TestBIF_pow_int_int(t *testing.T) {
	// Exact results, including ones that a float64 can't hold
	cases := []struct {
		base, exponent, expected int64
	}{
		{2, 10, 1024},
		{-3, 3, -27},
		{0, 0, 1},
		{5, 0, 1},
		{-9223372036854775808, 1, -9223372036854775808},
		{3, 39, 4052555153018976267},
		{7, 22, 3909821048582988049},
		{-2, 63, -9223372036854775808},
		{-1, 1000000000001, -1},
	}
	for _, c := range cases {
		output := BIF_pow(mlrval.FromInt(c.base), mlrval.FromInt(c.exponent))
		intval, ok := output.GetIntValue()
		assert.True(t, ok, "%d ** %d", c.base, c.exponent)
		assert.Equal(t, c.expected, intval, "%d ** %d", c.base, c.exponent)
	}

	// Overflow goes to float
	for _, c := range [][2]int64{{2, 63}, {8, 21}, {-2, 64}, {3, 40}} {
		output := BIF_pow(mlrval.FromInt(c[0]), mlrval.FromInt(c[1]))
		_, ok := output.GetFloatValue()
		assert.True(t, ok, "%d ** %d", c[0], c[1])
	}

	// Negative exponents
	output := BIF_pow(mlrval.FromInt(2), mlrval.FromInt(-1))
	floatval, ok := output.GetFloatValue()
	assert.True(t, ok)
	assert.Equal(t, 0.5, floatval)
}

// TODO: copy in more unit-test cases from existing regression-test data

//func BIF_minus_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_times(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_divide(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_int_divide(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_dot_plus(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_dot_minus(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_dot_times(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_dot_divide(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_modulus(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_mod_add(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_mod_sub(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_mod_mul(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_mod_exp(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_min_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_min_variadic(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval
//func BIF_max_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
//func BIF_max_variadic(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval

func TestBIF_dot_divide_int_by_zero(t *testing.T) {
	output := BIF_dot_divide(mlrval.FromInt(5), mlrval.FromInt(0))
	floatval, ok := output.GetFloatValue()
	assert.True(t, ok)
	assert.True(t, math.IsInf(floatval, 1))

	output = BIF_dot_divide(mlrval.FromInt(-5), mlrval.FromInt(0))
	floatval, ok = output.GetFloatValue()
	assert.True(t, ok)
	assert.True(t, math.IsInf(floatval, -1))

	output = BIF_dot_divide(mlrval.FromInt(7), mlrval.FromInt(2))
	intval, ok := output.GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, int64(3), intval)
}

func TestBIF_mod_exp_small_exponents(t *testing.T) {
	cases := []struct {
		a, e, m, want int64
	}{
		{10, 0, 7, 1},
		{10, 1, 7, 3},
		{-3, 1, 7, 4},
		{5, 0, 1, 0},
		{5, 1, 1, 0},
		{10, 2, 7, 2},
		{3, 5, 7, 5},
	}
	for _, c := range cases {
		output := BIF_mod_exp(mlrval.FromInt(c.a), mlrval.FromInt(c.e), mlrval.FromInt(c.m))
		intval, ok := output.GetIntValue()
		assert.True(t, ok)
		assert.Equal(t, c.want, intval, "mexp(%d, %d, %d)", c.a, c.e, c.m)
	}
}

func TestBIF_mod_ops_zero_modulus(t *testing.T) {
	a := mlrval.FromInt(5)
	b := mlrval.FromInt(3)
	m := mlrval.FromInt(0)
	assert.True(t, BIF_mod_add(a, b, m).IsError())
	assert.True(t, BIF_mod_sub(a, b, m).IsError())
	assert.True(t, BIF_mod_mul(a, b, m).IsError())
	assert.True(t, BIF_mod_exp(a, b, m).IsError())
}

func TestBIF_int_preserving_unary_math_large_ints(t *testing.T) {
	// 2**53 + 1 is not representable as a float64.
	big := int64(9007199254740993)
	for _, f := range []func(*mlrval.Mlrval) *mlrval.Mlrval{BIF_abs, BIF_ceil, BIF_floor, BIF_round} {
		intval, ok := f(mlrval.FromInt(big)).GetIntValue()
		assert.True(t, ok)
		assert.Equal(t, big, intval)
	}

	intval, ok := BIF_abs(mlrval.FromInt(-big)).GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, big, intval)

	// -math.MinInt64 does not fit in an int, so it overflows to float.
	floatval, ok := BIF_abs(mlrval.FromInt(math.MinInt64)).GetFloatValue()
	assert.True(t, ok)
	assert.Equal(t, 9223372036854775808.0, floatval)

	intval, ok = BIF_sgn(mlrval.FromInt(-big)).GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, int64(-1), intval)
	intval, ok = BIF_sgn(mlrval.FromInt(0)).GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, int64(0), intval)
}
