// ================================================================
// Boolean expressions for ==, !=, >, >=, <, <=, <=> on Mlrvals.
//
// Note that in bifs/boolean.go we have similar functions which take pairs of
// Mlrvals as input and return Mlrval as output. Those are for use in the
// Miller DSL. The functions here are primarily for 'mlr sort'. Their benefit
// is they don't allocate memory, and so are more efficient for sort we don't
// want to trigger lots of allocations, nor garbage collection, if we can avoid
// it.
// ================================================================

// TODO: comment about mvtype; deferral; copying of deferrence.

package mlrval

import (
	"github.com/johnkerl/miller/internal/pkg/lib"
)

type CmpFuncBool func(input1, input2 *Mlrval) bool
type CmpFuncInt func(input1, input2 *Mlrval) int // -1, 0, 1 for <=>

// ----------------------------------------------------------------
// Exported methods

func Equals(input1, input2 *Mlrval) bool {
	return cmp_dispositions[input1.Type()][input2.Type()](input1, input2) == 0
}
func NotEquals(input1, input2 *Mlrval) bool {
	return cmp_dispositions[input1.Type()][input2.Type()](input1, input2) != 0
}
func GreaterThan(input1, input2 *Mlrval) bool {
	return cmp_dispositions[input1.Type()][input2.Type()](input1, input2) > 0
}
func GreaterThanOrEquals(input1, input2 *Mlrval) bool {
	return cmp_dispositions[input1.Type()][input2.Type()](input1, input2) >= 0
}
func LessThan(input1, input2 *Mlrval) bool {
	return cmp_dispositions[input1.Type()][input2.Type()](input1, input2) < 0
}
func LessThanOrEquals(input1, input2 *Mlrval) bool {
	return cmp_dispositions[input1.Type()][input2.Type()](input1, input2) <= 0
}
func Cmp(input1, input2 *Mlrval) int {
	return cmp_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// Support routines for disposition-matrix entries

// _same returns int 0 as a binary-input function
func _same(input1, input2 *Mlrval) int {
	return 0
}

// _less returns int -1 as a binary-input function
func _less(input1, input2 *Mlrval) int {
	return -1
}

// _more returns int 1 as a binary-input function
func _more(input1, input2 *Mlrval) int {
	return 1
}

// int_cmp implements the spaceship operator for ints.
func int_cmp(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// float_cmp implements the spaceship operator for floats.
func float_cmp(a, b float64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// string_cmp implements the spaceship operator for strings.
func string_cmp(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ----------------------------------------------------------------
// Disposition-matrix entries

func cmp_b_ss(input1, input2 *Mlrval) int {
	return string_cmp(input1.printrep, input2.printrep)
}
func cmp_b_xs(input1, input2 *Mlrval) int {
	return string_cmp(input1.String(), input2.printrep)
}
func cmp_b_sx(input1, input2 *Mlrval) int {
	return string_cmp(input1.printrep, input2.String())
}
func cmp_b_ii(input1, input2 *Mlrval) int {
	return int_cmp(input1.intval, input2.intval)
}
func cmp_b_if(input1, input2 *Mlrval) int {
	return float_cmp(float64(input1.intval), input2.floatval)
}
func cmp_b_fi(input1, input2 *Mlrval) int {
	return float_cmp(input1.floatval, float64(input2.intval))
}
func cmp_b_ff(input1, input2 *Mlrval) int {
	return float_cmp(input1.floatval, input2.floatval)
}
func cmp_b_bb(input1, input2 *Mlrval) int {
	return int_cmp(lib.BoolToInt(input1.boolval), lib.BoolToInt(input2.boolval))
}

// TODO: cmp on array & map
//func eq_b_aa(input1, input2 *Mlrval) bool {
//	a := input1.arrayval
//	b := input2.arrayval
//
//	// Different-length arrays are not equal
//	if len(a) != len(b) {
//		return false
//	}
//
//	// Same-length arrays: return false if any slot is not equal, else true.
//	for i := range a {
//		if !Equals(&a[i], &b[i]) {
//			return false
//		}
//	}
//
//	return true
//}

//func eq_b_mm(input1, input2 *Mlrval) bool {
//	return input1.mapval.Equals(input2.mapval)
//}

var cmp_dispositions = [MT_DIM][MT_DIM]CmpFuncInt{
	//       .  INT        FLOAT     BOOL      VOID      STRING    ARRAY  MAP    FUNC   ERROR  NULL   ABSENT
	/*INT    */ {cmp_b_ii, cmp_b_if, _less, _less, _less, _less, _less, _less, _less, _less, _less},
	/*FLOAT  */ {cmp_b_fi, cmp_b_ff, _less, _less, _less, _less, _less, _less, _less, _less, _less},
	/*BOOL   */ {_more, _more, cmp_b_bb, _less, _less, _less, _less, _less, _less, _less, _less},
	/*VOID   */ {_more, _more, _more, cmp_b_ss, cmp_b_ss, _less, _less, _less, _less, _less, _less},
	/*STRING */ {_more, _more, _more, cmp_b_ss, cmp_b_ss, _less, _less, _less, _less, _less, _less},
	/*ARRAY  */ {_more, _more, _more, _more, _more, _same, _less, _less, _less, _less, _less},
	/*MAP    */ {_more, _more, _more, _more, _more, _more, _same, _less, _less, _less, _less},
	/*func   */ {_more, _more, _more, _more, _more, _more, _more, _same, _less, _less, _less},
	/*ERROR  */ {_more, _more, _more, _more, _more, _more, _more, _more, _same, _less, _less},
	/*NULL   */ {_more, _more, _more, _more, _more, _more, _more, _more, _more, _same, _less},
	/*ABSENT */ {_more, _more, _more, _more, _more, _more, _more, _more, _more, _more, _same},
}
