// ================================================================
// Boolean expressions for ==, !=, >, >=, <, <=, <=> on Mlrvals
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
	return eq_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}
func NotEquals(input1, input2 *Mlrval) bool {
	return ne_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}
func GreaterThan(input1, input2 *Mlrval) bool {
	return gt_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}
func GreaterThanOrEquals(input1, input2 *Mlrval) bool {
	return ge_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}
func LessThan(input1, input2 *Mlrval) bool {
	return lt_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}
func LessThanOrEquals(input1, input2 *Mlrval) bool {
	return le_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}
func Cmp(input1, input2 *Mlrval) int {
	return cmp_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ----------------------------------------------------------------
// Support routines for disposition-matrix entries

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

// _true return boolean true as a binary-input function
func _true(input1, input2 *Mlrval) bool {
	return true
}

// _fals return boolean false as a binary-input function
func _fals(input1, input2 *Mlrval) bool {
	return false
}

// _i0__ return int 0 as a binary-input function
func _i0__(input1, input2 *Mlrval) int {
	return 0
}

// _i1__ return int 1 as a binary-input function
func _i1__(input1, input2 *Mlrval) int {
	return 1
}

// _n1__ return int -1 as a binary-input function
func _n1__(input1, input2 *Mlrval) int {
	return -1
}

// ----------------------------------------------------------------
// Disposition-matrix entries

func eq_b_ss(input1, input2 *Mlrval) bool {
	return input1.printrep == input2.printrep
}
func ne_b_ss(input1, input2 *Mlrval) bool {
	return input1.printrep != input2.printrep
}
func gt_b_ss(input1, input2 *Mlrval) bool {
	return input1.printrep > input2.printrep
}
func ge_b_ss(input1, input2 *Mlrval) bool {
	return input1.printrep >= input2.printrep
}
func lt_b_ss(input1, input2 *Mlrval) bool {
	return input1.printrep < input2.printrep
}
func le_b_ss(input1, input2 *Mlrval) bool {
	return input1.printrep <= input2.printrep
}
func cmp_b_ss(input1, input2 *Mlrval) int {
	return string_cmp(input1.printrep, input2.printrep)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

func eq_b_xs(input1, input2 *Mlrval) bool {
	return input1.String() == input2.printrep
}
func ne_b_xs(input1, input2 *Mlrval) bool {
	return input1.String() != input2.printrep
}
func gt_b_xs(input1, input2 *Mlrval) bool {
	return input1.String() > input2.printrep
}
func ge_b_xs(input1, input2 *Mlrval) bool {
	return input1.String() >= input2.printrep
}
func lt_b_xs(input1, input2 *Mlrval) bool {
	return input1.String() < input2.printrep
}
func le_b_xs(input1, input2 *Mlrval) bool {
	return input1.String() <= input2.printrep
}
func cmp_b_xs(input1, input2 *Mlrval) int {
	return string_cmp(input1.String(), input2.printrep)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

func eq_b_sx(input1, input2 *Mlrval) bool {
	return input1.printrep == input2.String()
}
func ne_b_sx(input1, input2 *Mlrval) bool {
	return input1.printrep != input2.String()
}
func gt_b_sx(input1, input2 *Mlrval) bool {
	return input1.printrep > input2.String()
}
func ge_b_sx(input1, input2 *Mlrval) bool {
	return input1.printrep >= input2.String()
}
func lt_b_sx(input1, input2 *Mlrval) bool {
	return input1.printrep < input2.String()
}
func le_b_sx(input1, input2 *Mlrval) bool {
	return input1.printrep <= input2.String()
}
func cmp_b_sx(input1, input2 *Mlrval) int {
	return string_cmp(input1.printrep, input2.String())
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

func eq_b_ii(input1, input2 *Mlrval) bool {
	return input1.intval == input2.intval
}
func ne_b_ii(input1, input2 *Mlrval) bool {
	return input1.intval != input2.intval
}
func gt_b_ii(input1, input2 *Mlrval) bool {
	return input1.intval > input2.intval
}
func ge_b_ii(input1, input2 *Mlrval) bool {
	return input1.intval >= input2.intval
}
func lt_b_ii(input1, input2 *Mlrval) bool {
	return input1.intval < input2.intval
}
func le_b_ii(input1, input2 *Mlrval) bool {
	return input1.intval <= input2.intval
}
func cmp_b_ii(input1, input2 *Mlrval) int {
	return int_cmp(input1.intval, input2.intval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

func eq_b_if(input1, input2 *Mlrval) bool {
	return float64(input1.intval) == input2.floatval
}
func ne_b_if(input1, input2 *Mlrval) bool {
	return float64(input1.intval) != input2.floatval
}
func gt_b_if(input1, input2 *Mlrval) bool {
	return float64(input1.intval) > input2.floatval
}
func ge_b_if(input1, input2 *Mlrval) bool {
	return float64(input1.intval) >= input2.floatval
}
func lt_b_if(input1, input2 *Mlrval) bool {
	return float64(input1.intval) < input2.floatval
}
func le_b_if(input1, input2 *Mlrval) bool {
	return float64(input1.intval) <= input2.floatval
}
func cmp_b_if(input1, input2 *Mlrval) int {
	return float_cmp(float64(input1.intval), input2.floatval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

func eq_b_fi(input1, input2 *Mlrval) bool {
	return input1.floatval == float64(input2.intval)
}
func ne_b_fi(input1, input2 *Mlrval) bool {
	return input1.floatval != float64(input2.intval)
}
func gt_b_fi(input1, input2 *Mlrval) bool {
	return input1.floatval > float64(input2.intval)
}
func ge_b_fi(input1, input2 *Mlrval) bool {
	return input1.floatval >= float64(input2.intval)
}
func lt_b_fi(input1, input2 *Mlrval) bool {
	return input1.floatval < float64(input2.intval)
}
func le_b_fi(input1, input2 *Mlrval) bool {
	return input1.floatval <= float64(input2.intval)
}
func cmp_b_fi(input1, input2 *Mlrval) int {
	return float_cmp(input1.floatval, float64(input2.intval))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

func eq_b_ff(input1, input2 *Mlrval) bool {
	return input1.floatval == input2.floatval
}
func ne_b_ff(input1, input2 *Mlrval) bool {
	return input1.floatval != input2.floatval
}
func gt_b_ff(input1, input2 *Mlrval) bool {
	return input1.floatval > input2.floatval
}
func ge_b_ff(input1, input2 *Mlrval) bool {
	return input1.floatval >= input2.floatval
}
func lt_b_ff(input1, input2 *Mlrval) bool {
	return input1.floatval < input2.floatval
}
func le_b_ff(input1, input2 *Mlrval) bool {
	return input1.floatval <= input2.floatval
}
func cmp_b_ff(input1, input2 *Mlrval) int {
	return float_cmp(input1.floatval, input2.floatval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

func eq_b_bb(input1, input2 *Mlrval) bool {
	return input1.boolval == input2.boolval
}
func ne_b_bb(input1, input2 *Mlrval) bool {
	return input1.boolval != input2.boolval
}

func gt_b_bb(input1, input2 *Mlrval) bool {
	return lib.BoolToInt(input1.boolval) > lib.BoolToInt(input2.boolval)
}
func ge_b_bb(input1, input2 *Mlrval) bool {
	return lib.BoolToInt(input1.boolval) >= lib.BoolToInt(input2.boolval)
}
func lt_b_bb(input1, input2 *Mlrval) bool {
	return lib.BoolToInt(input1.boolval) < lib.BoolToInt(input2.boolval)
}
func le_b_bb(input1, input2 *Mlrval) bool {
	return lib.BoolToInt(input1.boolval) <= lib.BoolToInt(input2.boolval)
}
func cmp_b_bb(input1, input2 *Mlrval) int {
	return int_cmp(lib.BoolToInt(input1.boolval), lib.BoolToInt(input2.boolval))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

func eq_b_aa(input1, input2 *Mlrval) bool {
	a := input1.arrayval
	b := input2.arrayval

	// Different-length arrays are not equal
	if len(a) != len(b) {
		return false
	}

	// Same-length arrays: return false if any slot is not equal, else true.
	for i := range a {
		if !Equals(&a[i], &b[i]) {
			return false
		}
	}

	return true
}
func ne_b_aa(input1, input2 *Mlrval) bool {
	return !eq_b_aa(input1, input2)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

func eq_b_mm(input1, input2 *Mlrval) bool {
	return input1.mapval.Equals(input2.mapval)
}
func ne_b_mm(input1, input2 *Mlrval) bool {
	return !input1.mapval.Equals(input2.mapval)
}

// We get a Golang "initialization loop" due to recursive depth computation
// if this is defined statically. So, we use a "package init" function.
var eq_dispositions = [MT_DIM][MT_DIM]CmpFuncBool{}

func init() {
	eq_dispositions = [MT_DIM][MT_DIM]CmpFuncBool{
		//       .  ERROR   ABSENT  NULL    VOID     STRING   INT      FLOAT    BOOL     ARRAY    MAP      FUNC
		/*ERROR  */ {_true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
		/*ABSENT */ {_fals, _true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
		/*NULL   */ {_fals, _fals, _true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
		/*VOID   */ {_fals, _fals, _fals, eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _fals, _fals, _fals, _fals},
		/*STRING */ {_fals, _fals, _fals, eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _fals, _fals, _fals, _fals},
		/*INT    */ {_fals, _fals, _fals, eq_b_xs, eq_b_xs, eq_b_ii, eq_b_if, _fals, _fals, _fals, _fals},
		/*FLOAT  */ {_fals, _fals, _fals, eq_b_xs, eq_b_xs, eq_b_fi, eq_b_ff, _fals, _fals, _fals, _fals},
		/*BOOL   */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, eq_b_bb, _fals, _fals, _fals},
		/*ARRAY  */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, eq_b_aa, _fals, _fals},
		/*MAP    */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, eq_b_mm, _fals},
		/*FUNC   */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _true},
	}
}

var ne_dispositions = [MT_DIM][MT_DIM]CmpFuncBool{
	//       .  ERROR   ABSENT NULL   VOID     STRING   INT      FLOAT    BOOL     ARRAY    MAP      FUNC
	/*ERROR  */ {_true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*ABSENT */ {_fals, _true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*NULL   */ {_fals, _fals, _fals, _true, _true, _true, _true, _true, _true, _true, _fals},
	/*VOID   */ {_fals, _fals, _true, ne_b_ss, ne_b_ss, ne_b_sx, ne_b_sx, _true, _true, _true, _fals},
	/*STRING */ {_fals, _fals, _true, ne_b_ss, ne_b_ss, ne_b_sx, ne_b_sx, _true, _true, _true, _fals},
	/*INT    */ {_fals, _fals, _true, ne_b_xs, ne_b_xs, ne_b_ii, ne_b_if, _true, _true, _true, _fals},
	/*FLOAT  */ {_fals, _fals, _true, ne_b_xs, ne_b_xs, ne_b_fi, ne_b_ff, _true, _true, _true, _fals},
	/*BOOL   */ {_fals, _fals, _true, _true, _true, _true, _true, ne_b_bb, _true, _true, _fals},
	/*ARRAY  */ {_fals, _fals, _true, _true, _true, _true, _true, _true, ne_b_aa, _true, _fals},
	/*MAP    */ {_fals, _fals, _true, _true, _true, _true, _true, _true, _true, ne_b_mm, _fals},
	/*FUNC   */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _true},
}

var gt_dispositions = [MT_DIM][MT_DIM]CmpFuncBool{
	//       .  ERROR   ABSENT NULL   VOID     STRING   INT      FLOAT    BOOL     ARRAY  MAP    FUNC
	/*ERROR  */ {_true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*ABSENT */ {_fals, _true, _true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*NULL   */ {_true, _fals, _fals, _true, _true, _true, _true, _true, _fals, _fals, _fals},
	/*VOID   */ {_fals, _fals, _fals, gt_b_ss, gt_b_ss, gt_b_sx, gt_b_sx, _fals, _fals, _fals, _fals},
	/*STRING */ {_fals, _fals, _fals, gt_b_ss, gt_b_ss, gt_b_sx, gt_b_sx, _fals, _fals, _fals, _fals},
	/*INT    */ {_fals, _fals, _fals, gt_b_xs, gt_b_xs, gt_b_ii, gt_b_if, _fals, _fals, _fals, _fals},
	/*FLOAT  */ {_fals, _fals, _fals, gt_b_xs, gt_b_xs, gt_b_fi, gt_b_ff, _fals, _fals, _fals, _fals},
	/*BOOL   */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, gt_b_bb, _fals, _fals, _fals},
	/*ARRAY  */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*MAP    */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*FUNC   */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _true},
}

var ge_dispositions = [MT_DIM][MT_DIM]CmpFuncBool{
	//       .  ERROR   ABSENT NULL   VOID     STRING   INT      FLOAT    BOOL     ARRAY  MAP    FUNC
	/*ERROR  */ {_true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*ABSENT */ {_fals, _true, _true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*NULL   */ {_true, _fals, _true, _true, _true, _true, _true, _true, _fals, _fals, _fals},
	/*VOID   */ {_fals, _fals, _fals, ge_b_ss, ge_b_ss, ge_b_sx, ge_b_sx, _fals, _fals, _fals, _fals},
	/*STRING */ {_fals, _fals, _fals, ge_b_ss, ge_b_ss, ge_b_sx, ge_b_sx, _fals, _fals, _fals, _fals},
	/*INT    */ {_fals, _fals, _fals, ge_b_xs, ge_b_xs, ge_b_ii, ge_b_if, _fals, _fals, _fals, _fals},
	/*FLOAT  */ {_fals, _fals, _fals, ge_b_xs, ge_b_xs, ge_b_fi, ge_b_ff, _fals, _fals, _fals, _fals},
	/*BOOL   */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, ge_b_bb, _fals, _fals, _fals},
	/*ARRAY  */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*MAP    */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*FUNC   */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _true},
}

var lt_dispositions = [MT_DIM][MT_DIM]CmpFuncBool{
	//       .  ERROR   ABSENT NULL   VOID     STRING   INT      FLOAT    BOOL     ARRAY  MAP    FUNC
	/*ERROR  */ {_true, _fals, _true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*ABSENT */ {_fals, _true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*NULL   */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*VOID   */ {_fals, _fals, _true, lt_b_ss, lt_b_ss, lt_b_sx, lt_b_sx, _fals, _fals, _fals, _fals},
	/*STRING */ {_fals, _fals, _true, lt_b_ss, lt_b_ss, lt_b_sx, lt_b_sx, _fals, _fals, _fals, _fals},
	/*INT    */ {_fals, _fals, _true, lt_b_xs, lt_b_xs, lt_b_ii, lt_b_if, _fals, _fals, _fals, _fals},
	/*FLOAT  */ {_fals, _fals, _true, lt_b_xs, lt_b_xs, lt_b_fi, lt_b_ff, _fals, _fals, _fals, _fals},
	/*BOOL   */ {_fals, _fals, _true, _fals, _fals, _fals, _fals, lt_b_bb, _fals, _fals, _fals},
	/*ARRAY  */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*MAP    */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*FUNC   */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _true},
}

var le_dispositions = [MT_DIM][MT_DIM]CmpFuncBool{
	//       .  ERROR   ABSENT NULL   VOID     STRING   INT      FLOAT    BOOL     ARRAY  MAP    FUNC
	/*ERROR  */ {_true, _fals, _true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*ABSENT */ {_fals, _true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*NULL   */ {_fals, _fals, _true, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*VOID   */ {_fals, _fals, _true, le_b_ss, le_b_ss, le_b_sx, le_b_sx, _fals, _fals, _fals, _fals},
	/*STRING */ {_fals, _fals, _true, le_b_ss, le_b_ss, le_b_sx, le_b_sx, _fals, _fals, _fals, _fals},
	/*INT    */ {_fals, _fals, _true, le_b_xs, le_b_xs, le_b_ii, le_b_if, _fals, _fals, _fals, _fals},
	/*FLOAT  */ {_fals, _fals, _true, le_b_xs, le_b_xs, le_b_fi, le_b_ff, _fals, _fals, _fals, _fals},
	/*BOOL   */ {_fals, _fals, _true, _fals, _fals, _fals, _fals, le_b_bb, _fals, _fals, _fals},
	/*ARRAY  */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*MAP    */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals},
	/*FUNC   */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _fals, _true},
}

// TODO: flesh these out for array and map
var cmp_dispositions = [MT_DIM][MT_DIM]CmpFuncInt{
	//       .  ERROR   ABSENT NULL   VOID      STRING    INT       FLOAT     BOOL      ARRAY  MAP    FUNC
	/*ERROR  */ {_i0__, _i0__, _n1__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__},
	/*ABSENT */ {_i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__},
	/*NULL   */ {_i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__},
	/*VOID   */ {_i0__, _i0__, _n1__, cmp_b_ss, cmp_b_ss, cmp_b_sx, cmp_b_sx, _i0__, _i0__, _i0__, _i0__},
	/*STRING */ {_i0__, _i0__, _n1__, cmp_b_ss, cmp_b_ss, cmp_b_sx, cmp_b_sx, _i0__, _i0__, _i0__, _i0__},
	/*INT    */ {_i0__, _i0__, _n1__, cmp_b_xs, cmp_b_xs, cmp_b_ii, cmp_b_if, _i0__, _i0__, _i0__, _i0__},
	/*FLOAT  */ {_i0__, _i0__, _n1__, cmp_b_xs, cmp_b_xs, cmp_b_fi, cmp_b_ff, _i0__, _i0__, _i0__, _i0__},
	/*BOOL   */ {_i0__, _i0__, _n1__, _i0__, _i0__, _i0__, _i0__, cmp_b_bb, _i0__, _i0__, _i0__},
	/*ARRAY  */ {_i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__},
	/*MAP    */ {_i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__},
	/*FUNC   */ {_i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__, _i0__},
}
