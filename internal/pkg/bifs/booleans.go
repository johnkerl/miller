// ================================================================
// Boolean expressions for ==, !=, >, >=, <, <=
// ================================================================

package bifs

import (
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
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

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_ss(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireStringValue() == input2.AcquireStringValue())
}
func ne_b_ss(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireStringValue() != input2.AcquireStringValue())
}
func gt_b_ss(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireStringValue() > input2.AcquireStringValue())
}
func ge_b_ss(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireStringValue() >= input2.AcquireStringValue())
}
func lt_b_ss(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireStringValue() < input2.AcquireStringValue())
}
func le_b_ss(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireStringValue() <= input2.AcquireStringValue())
}
func cmp_b_ss(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(string_cmp(input1.AcquireStringValue(), input2.AcquireStringValue()))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_xs(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.String() == input2.AcquireStringValue())
}
func ne_b_xs(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.String() != input2.AcquireStringValue())
}
func gt_b_xs(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.String() > input2.AcquireStringValue())
}
func ge_b_xs(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.String() >= input2.AcquireStringValue())
}
func lt_b_xs(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.String() < input2.AcquireStringValue())
}
func le_b_xs(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.String() <= input2.AcquireStringValue())
}
func cmp_b_xs(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(string_cmp(input1.String(), input2.AcquireStringValue()))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_sx(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireStringValue() == input2.String())
}
func ne_b_sx(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireStringValue() != input2.String())
}
func gt_b_sx(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireStringValue() > input2.String())
}
func ge_b_sx(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireStringValue() >= input2.String())
}
func lt_b_sx(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireStringValue() < input2.String())
}
func le_b_sx(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireStringValue() <= input2.String())
}
func cmp_b_sx(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(string_cmp(input1.AcquireStringValue(), input2.String()))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireIntValue() == input2.AcquireIntValue())
}
func ne_b_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireIntValue() != input2.AcquireIntValue())
}
func gt_b_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireIntValue() > input2.AcquireIntValue())
}
func ge_b_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireIntValue() >= input2.AcquireIntValue())
}
func lt_b_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireIntValue() < input2.AcquireIntValue())
}
func le_b_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireIntValue() <= input2.AcquireIntValue())
}
func cmp_b_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(int_cmp(input1.AcquireIntValue(), input2.AcquireIntValue()))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(float64(input1.AcquireIntValue()) == input2.AcquireFloatValue())
}
func ne_b_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(float64(input1.AcquireIntValue()) != input2.AcquireFloatValue())
}
func gt_b_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(float64(input1.AcquireIntValue()) > input2.AcquireFloatValue())
}
func ge_b_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(float64(input1.AcquireIntValue()) >= input2.AcquireFloatValue())
}
func lt_b_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(float64(input1.AcquireIntValue()) < input2.AcquireFloatValue())
}
func le_b_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(float64(input1.AcquireIntValue()) <= input2.AcquireFloatValue())
}
func cmp_b_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(float_cmp(float64(input1.AcquireIntValue()), input2.AcquireFloatValue()))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireFloatValue() == float64(input2.AcquireIntValue()))
}
func ne_b_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireFloatValue() != float64(input2.AcquireIntValue()))
}
func gt_b_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireFloatValue() > float64(input2.AcquireIntValue()))
}
func ge_b_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireFloatValue() >= float64(input2.AcquireIntValue()))
}
func lt_b_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireFloatValue() < float64(input2.AcquireIntValue()))
}
func le_b_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireFloatValue() <= float64(input2.AcquireIntValue()))
}
func cmp_b_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(float_cmp(input1.AcquireFloatValue(), float64(input2.AcquireIntValue())))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireFloatValue() == input2.AcquireFloatValue())
}
func ne_b_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireFloatValue() != input2.AcquireFloatValue())
}
func gt_b_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireFloatValue() > input2.AcquireFloatValue())
}
func ge_b_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireFloatValue() >= input2.AcquireFloatValue())
}
func lt_b_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireFloatValue() < input2.AcquireFloatValue())
}
func le_b_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireFloatValue() <= input2.AcquireFloatValue())
}
func cmp_b_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(float_cmp(input1.AcquireFloatValue(), input2.AcquireFloatValue()))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_bb(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireBoolValue() == input2.AcquireBoolValue())
}
func ne_b_bb(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireBoolValue() != input2.AcquireBoolValue())
}

// We could say ordering on bool is error, but, Miller allows
// sorting on bool so it should allow ordering on bool.

func gt_b_bb(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(lib.BoolToInt(input1.AcquireBoolValue()) > lib.BoolToInt(input2.AcquireBoolValue()))
}
func ge_b_bb(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(lib.BoolToInt(input1.AcquireBoolValue()) >= lib.BoolToInt(input2.AcquireBoolValue()))
}
func lt_b_bb(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(lib.BoolToInt(input1.AcquireBoolValue()) < lib.BoolToInt(input2.AcquireBoolValue()))
}
func le_b_bb(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(lib.BoolToInt(input1.AcquireBoolValue()) <= lib.BoolToInt(input2.AcquireBoolValue()))
}
func cmp_b_bb(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(int_cmp(lib.BoolToInt(input1.AcquireBoolValue()), lib.BoolToInt(input2.AcquireBoolValue())))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_aa(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	a := input1.AcquireArrayValue()
	b := input2.AcquireArrayValue()

	// Different-length arrays are not equal
	if len(a) != len(b) {
		return mlrval.FALSE
	}

	// Same-length arrays: return false if any slot is not equal, else true.
	for i := range a {
		eq := BIF_equals(&a[i], &b[i])
		lib.InternalCodingErrorIf(eq.Type() != mlrval.MT_BOOL)
		if eq.AcquireBoolValue() == false {
			return mlrval.FALSE
		}
	}

	return mlrval.TRUE
}
func ne_b_aa(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	output := eq_b_aa(input1, input2)
	return mlrval.FromBool(!output.AcquireBoolValue())
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_mm(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireMapValue().Equals(input2.AcquireMapValue()))
}
func ne_b_mm(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(!input1.AcquireMapValue().Equals(input2.AcquireMapValue()))
}

// We get a Golang "initialization loop" due to recursive depth computation
// if this is defined statically. So, we use a "package init" function.
var eq_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{}

func init() {
	eq_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
		//       .  INT      FLOAT    BOOL     NULL   STRING   ARRAY    MAP      FUNC    ERROR   ABSENT NULL
		/*INT    */ {eq_b_ii, eq_b_if, _fals, _fals, eq_b_xs, _fals, _fals, _erro, _erro, _absn, _fals},
		/*FLOAT  */ {eq_b_fi, eq_b_ff, _fals, _fals, eq_b_xs, _fals, _fals, _erro, _erro, _absn, _fals},
		/*BOOL   */ {_fals, _fals, eq_b_bb, _fals, _fals, _fals, _fals, _erro, _erro, _absn, _fals},
		/*VOID   */ {eq_b_sx, eq_b_sx, _fals, _fals, eq_b_ss, _fals, _fals, _erro, _erro, _absn, _fals},
		/*STRING */ {eq_b_sx, eq_b_sx, _fals, _fals, eq_b_ss, _fals, _fals, _erro, _erro, _absn, _fals},
		/*ARRAY  */ {_fals, _fals, _fals, _fals, _fals, eq_b_aa, _fals, _erro, _erro, _absn, _fals},
		/*MAP    */ {_fals, _fals, _fals, _fals, _fals, _fals, eq_b_mm, _erro, _erro, _absn, _fals},
		/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
		/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
		/*ABSENT */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _erro, _absn, _absn},
		/*NULL   */ {_fals, _fals, _fals, _true, _fals, _fals, _fals, _erro, _erro, _absn, _true},
	}
}

var ne_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT      FLOAT    BOOL     NULL   STRING   ARRAY    MAP      FUNC    ERROR   ABSENT NULL
	/*INT    */ {ne_b_ii, ne_b_if, _true, _true, ne_b_xs, _true, _true, _erro, _erro, _absn, _true},
	/*FLOAT  */ {ne_b_fi, ne_b_ff, _true, _true, ne_b_xs, _true, _true, _erro, _erro, _absn, _true},
	/*BOOL   */ {_true, _true, ne_b_bb, _true, _true, _true, _true, _erro, _erro, _absn, _true},
	/*VOID   */ {ne_b_sx, ne_b_sx, _true, _true, ne_b_ss, _true, _true, _erro, _erro, _absn, _true},
	/*STRING */ {ne_b_sx, ne_b_sx, _true, _true, ne_b_ss, _true, _true, _erro, _erro, _absn, _true},
	/*ARRAY  */ {_true, _true, _true, _true, _true, ne_b_aa, _true, _erro, _erro, _absn, _true},
	/*MAP    */ {_true, _true, _true, _true, _true, _true, ne_b_mm, _erro, _erro, _absn, _true},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _erro, _absn, _absn},
	/*NULL   */ {_true, _true, _true, _fals, _true, _true, _true, _erro, _erro, _absn, _fals},
}
var gt_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT      FLOAT    BOOL     NULL   STRING   ARRAY  MAP    FUNC    ERROR   ABSENT NULL
	/*INT    */ {gt_b_ii, gt_b_if, _fals, _fals, gt_b_xs, _fals, _fals, _erro, _erro, _absn, _fals},
	/*FLOAT  */ {gt_b_fi, gt_b_ff, _fals, _fals, gt_b_xs, _fals, _fals, _erro, _erro, _absn, _fals},
	/*BOOL   */ {_fals, _fals, gt_b_bb, _fals, _fals, _fals, _fals, _erro, _erro, _absn, _fals},
	/*VOID   */ {gt_b_sx, gt_b_sx, _fals, _fals, gt_b_ss, _fals, _fals, _erro, _erro, _absn, _fals},
	/*STRING */ {gt_b_sx, gt_b_sx, _fals, _fals, gt_b_ss, _fals, _fals, _erro, _erro, _absn, _fals},
	/*ARRAY  */ {_fals, _fals, _fals, _fals, _fals, _erro, _fals, _erro, _erro, _absn, _fals},
	/*MAP    */ {_fals, _fals, _fals, _fals, _fals, _fals, _erro, _erro, _erro, _absn, _fals},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _fals, _erro, _erro, _erro, _erro, _erro, _erro, _fals},
	/*ABSENT */ {_absn, _absn, _absn, _true, _absn, _absn, _absn, _erro, _erro, _absn, _true},
	/*NULL   */ {_true, _true, _true, _fals, _true, _absn, _absn, _erro, _true, _fals, _fals},
}
var ge_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT      FLOAT    BOOL     NULL   STRING   ARRAY  MAP    FUNC    ERROR   ABSENT NULL
	/*INT    */ {ge_b_ii, ge_b_if, _fals, _fals, ge_b_xs, _fals, _fals, _erro, _erro, _absn, _fals},
	/*FLOAT  */ {ge_b_fi, ge_b_ff, _fals, _fals, ge_b_xs, _fals, _fals, _erro, _erro, _absn, _fals},
	/*BOOL   */ {_fals, _fals, ge_b_bb, _fals, _fals, _fals, _fals, _erro, _erro, _absn, _fals},
	/*VOID   */ {ge_b_sx, ge_b_sx, _fals, _fals, ge_b_ss, _fals, _fals, _erro, _erro, _absn, _fals},
	/*STRING */ {ge_b_sx, ge_b_sx, _fals, _fals, ge_b_ss, _fals, _fals, _erro, _erro, _absn, _fals},
	/*ARRAY  */ {_fals, _fals, _fals, _fals, _fals, _erro, _fals, _erro, _erro, _absn, _fals},
	/*MAP    */ {_fals, _fals, _fals, _fals, _fals, _fals, _erro, _erro, _erro, _absn, _fals},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _fals, _erro, _erro, _erro, _erro, _erro, _erro, _fals},
	/*ABSENT */ {_absn, _absn, _absn, _true, _absn, _absn, _absn, _erro, _erro, _absn, _true},
	/*NULL   */ {_true, _true, _true, _true, _true, _absn, _absn, _erro, _true, _fals, _true},
}
var lt_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT      FLOAT    BOOL     NULL   STRING   ARRAY  MAP    FUNC    ERROR   ABSENT NULL
	/*INT    */ {lt_b_ii, lt_b_if, _fals, _true, lt_b_xs, _fals, _fals, _erro, _erro, _absn, _true},
	/*FLOAT  */ {lt_b_fi, lt_b_ff, _fals, _true, lt_b_xs, _fals, _fals, _erro, _erro, _absn, _true},
	/*BOOL   */ {_fals, _fals, lt_b_bb, _true, _fals, _fals, _fals, _erro, _erro, _absn, _true},
	/*VOID   */ {lt_b_sx, lt_b_sx, _fals, _true, lt_b_ss, _fals, _fals, _erro, _erro, _absn, _true},
	/*STRING */ {lt_b_sx, lt_b_sx, _fals, _true, lt_b_ss, _fals, _fals, _erro, _erro, _absn, _true},
	/*ARRAY  */ {_fals, _fals, _fals, _absn, _fals, _erro, _fals, _erro, _erro, _absn, _absn},
	/*MAP    */ {_fals, _fals, _fals, _absn, _fals, _fals, _erro, _erro, _erro, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _true, _erro, _erro, _erro, _erro, _erro, _erro, _true},
	/*ABSENT */ {_absn, _absn, _absn, _fals, _absn, _absn, _absn, _erro, _erro, _absn, _fals},
	/*NULL   */ {_fals, _fals, _fals, _fals, _fals, _absn, _absn, _erro, _fals, _true, _fals},
}
var le_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT      FLOAT    BOOL     NULL   STRING   ARRAY  MAP    FUNC    ERROR   ABSENT NULL
	/*INT    */ {le_b_ii, le_b_if, _fals, _true, le_b_xs, _fals, _fals, _erro, _erro, _absn, _true},
	/*FLOAT  */ {le_b_fi, le_b_ff, _fals, _true, le_b_xs, _fals, _fals, _erro, _erro, _absn, _true},
	/*BOOL   */ {_fals, _fals, le_b_bb, _true, _fals, _fals, _fals, _erro, _erro, _absn, _true},
	/*VOID   */ {le_b_sx, le_b_sx, _fals, _true, le_b_ss, _fals, _fals, _erro, _erro, _absn, _true},
	/*STRING */ {le_b_sx, le_b_sx, _fals, _true, le_b_ss, _fals, _fals, _erro, _erro, _absn, _true},
	/*ARRAY  */ {_fals, _fals, _fals, _absn, _fals, _erro, _fals, _erro, _erro, _absn, _absn},
	/*MAP    */ {_fals, _fals, _fals, _absn, _fals, _fals, _erro, _erro, _erro, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _true, _erro, _erro, _erro, _erro, _erro, _erro, _true},
	/*ABSENT */ {_absn, _absn, _absn, _fals, _absn, _absn, _absn, _erro, _erro, _absn, _fals},
	/*NULL   */ {_fals, _fals, _fals, _true, _fals, _absn, _absn, _erro, _fals, _true, _true},
}
var cmp_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT       FLOAT     BOOL      NULL   STRING    ARRAY  MAP    FUNC    ERROR   ABSENT NULL
	/*INT    */ {cmp_b_ii, cmp_b_if, _fals, _true, cmp_b_xs, _fals, _fals, _erro, _erro, _absn, _true},
	/*FLOAT  */ {cmp_b_fi, cmp_b_ff, _fals, _true, cmp_b_xs, _fals, _fals, _erro, _erro, _absn, _true},
	/*BOOL   */ {_fals, _fals, cmp_b_bb, _true, _fals, _fals, _fals, _erro, _erro, _absn, _true},
	/*VOID   */ {cmp_b_sx, cmp_b_sx, _fals, _true, cmp_b_ss, _fals, _fals, _erro, _erro, _absn, _true},
	/*STRING */ {cmp_b_sx, cmp_b_sx, _fals, _true, cmp_b_ss, _fals, _fals, _erro, _erro, _absn, _true},
	/*ARRAY  */ {_fals, _fals, _fals, _absn, _fals, _erro, _fals, _erro, _erro, _absn, _absn},
	/*MAP    */ {_fals, _fals, _fals, _absn, _fals, _fals, _erro, _erro, _erro, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _true, _erro, _erro, _erro, _erro, _erro, _erro, _true},
	/*ABSENT */ {_absn, _absn, _absn, _fals, _absn, _absn, _absn, _erro, _erro, _absn, _fals},
	/*NULL   */ {_fals, _fals, _fals, _true, _fals, _absn, _absn, _erro, _fals, _true, _true},
}

func BIF_equals(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return eq_dispositions[input1.Type()][input2.Type()](input1, input2)
}
func BIF_not_equals(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return ne_dispositions[input1.Type()][input2.Type()](input1, input2)
}
func BIF_greater_than(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return gt_dispositions[input1.Type()][input2.Type()](input1, input2)
}
func BIF_greater_than_or_equals(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return ge_dispositions[input1.Type()][input2.Type()](input1, input2)
}
func BIF_less_than(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return lt_dispositions[input1.Type()][input2.Type()](input1, input2)
}
func BIF_less_than_or_equals(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return le_dispositions[input1.Type()][input2.Type()](input1, input2)
}
func BIF_cmp(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return cmp_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
func BIF_logical_NOT(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsBool() {
		return mlrval.FromBool(!input1.AcquireBoolValue())
	} else {
		return mlrval.ERROR
	}
}

func BIF_logical_AND(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsBool() && input2.IsBool() {
		return mlrval.FromBool(input1.AcquireBoolValue() && input2.AcquireBoolValue())
	} else {
		return mlrval.ERROR
	}
}

func BIF_logical_OR(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsBool() && input2.IsBool() {
		return mlrval.FromBool(input1.AcquireBoolValue() || input2.AcquireBoolValue())
	} else {
		return mlrval.ERROR
	}
}

func BIF_logical_XOR(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsBool() && input2.IsBool() {
		return mlrval.FromBool(input1.AcquireBoolValue() != input2.AcquireBoolValue())
	} else {
		return mlrval.ERROR
	}
}
