package bifs

import (
	"testing"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

// Disposition matrices/vectors are positional array literals sized by
// mlrval.MT_DIM. When a new mlrval type is added and MT_DIM is bumped, Go
// zero-fills the new slots in any table that hasn't been extended, leaving
// nil function pointers which panic at dispatch time rather than at compile
// time. These sweeps turn that hazard into a test failure.

func TestNoNilCellsInUnaryDispositionVectors(t *testing.T) {
	vectors := map[string][mlrval.MT_DIM]UnaryFunc{
		"upos_dispositions":        upos_dispositions,
		"uneg_dispositions":        uneg_dispositions,
		"min_unary_dispositions":   min_unary_dispositions,
		"max_unary_dispositions":   max_unary_dispositions,
		"bitwise_not_dispositions": bitwise_not_dispositions,
		"bitcount_dispositions":    bitcount_dispositions,
		"depth_dispositions":       depth_dispositions,
		"leafcount_dispositions":   leafcount_dispositions,
		"to_int_dispositions":      to_int_dispositions,
		"to_float_dispositions":    to_float_dispositions,
		"to_boolean_dispositions":  to_boolean_dispositions,
	}
	for name, vector := range vectors {
		for i := 0; i < int(mlrval.MT_DIM); i++ {
			if vector[i] == nil {
				t.Errorf("nil cell in %s[%d]", name, i)
			}
		}
	}
}

func TestNoNilCellsInUnaryMathLibDispositionVectors(t *testing.T) {
	vectors := map[string][mlrval.MT_DIM]mathLibUnaryFuncWrapper{
		"mudispo":  mudispo,
		"imudispo": imudispo,
	}
	for name, vector := range vectors {
		for i := 0; i < int(mlrval.MT_DIM); i++ {
			if vector[i] == nil {
				t.Errorf("nil cell in %s[%d]", name, i)
			}
		}
	}
}

func TestNoNilCellsInBinaryDispositionVectors(t *testing.T) {
	vectors := map[string][mlrval.MT_DIM]BinaryFunc{
		"to_int_with_base_dispositions": to_int_with_base_dispositions,
	}
	for name, vector := range vectors {
		for i := 0; i < int(mlrval.MT_DIM); i++ {
			if vector[i] == nil {
				t.Errorf("nil cell in %s[%d]", name, i)
			}
		}
	}
}

func TestNoNilCellsInBinaryDispositionMatrices(t *testing.T) {
	matrices := map[string][mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
		"plus_dispositions":                 plus_dispositions,
		"minus_dispositions":                minus_dispositions,
		"times_dispositions":                times_dispositions,
		"divide_dispositions":               divide_dispositions,
		"int_divide_dispositions":           int_divide_dispositions,
		"dot_plus_dispositions":             dot_plus_dispositions,
		"dotminus_dispositions":             dotminus_dispositions,
		"dottimes_dispositions":             dottimes_dispositions,
		"dotdivide_dispositions":            dotdivide_dispositions,
		"dotidivide_dispositions":           dotidivide_dispositions,
		"modulus_dispositions":              modulus_dispositions,
		"min_dispositions":                  min_dispositions,
		"max_dispositions":                  max_dispositions,
		"bitwise_and_dispositions":          bitwise_and_dispositions,
		"bitwise_or_dispositions":           bitwise_or_dispositions,
		"bitwise_xor_dispositions":          bitwise_xor_dispositions,
		"left_shift_dispositions":           left_shift_dispositions,
		"signed_right_shift_dispositions":   signed_right_shift_dispositions,
		"unsigned_right_shift_dispositions": unsigned_right_shift_dispositions,
		"eq_dispositions":                   eq_dispositions,
		"ne_dispositions":                   ne_dispositions,
		"gt_dispositions":                   gt_dispositions,
		"ge_dispositions":                   ge_dispositions,
		"lt_dispositions":                   lt_dispositions,
		"le_dispositions":                   le_dispositions,
		"cmp_dispositions":                  cmp_dispositions,
		"pow_dispositions":                  pow_dispositions,
		"atan2_dispositions":                atan2_dispositions,
		"roundm_dispositions":               roundm_dispositions,
		"dot_dispositions":                  dot_dispositions,
		"fmtnum_dispositions":               fmtnum_dispositions,
	}
	for name, matrix := range matrices {
		for i := 0; i < int(mlrval.MT_DIM); i++ {
			for j := 0; j < int(mlrval.MT_DIM); j++ {
				if matrix[i][j] == nil {
					t.Errorf("nil cell in %s[%d][%d]", name, i, j)
				}
			}
		}
	}
}
