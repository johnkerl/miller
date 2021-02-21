// ================================================================
// Boolean expressions for ==, !=, >, >=, <, <=
// ================================================================

package types

import (
	"miller/src/lib"
)

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_ss(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.printrep == input2.printrep)
}
func ne_b_ss(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.printrep != input2.printrep)
}
func gt_b_ss(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.printrep > input2.printrep)
}
func ge_b_ss(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.printrep >= input2.printrep)
}
func lt_b_ss(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.printrep < input2.printrep)
}
func le_b_ss(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.printrep <= input2.printrep)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_xs(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.String() == input2.printrep)
}
func ne_b_xs(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.String() != input2.printrep)
}
func gt_b_xs(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.String() > input2.printrep)
}
func ge_b_xs(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.String() >= input2.printrep)
}
func lt_b_xs(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.String() < input2.printrep)
}
func le_b_xs(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.String() <= input2.printrep)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_sx(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.printrep == input2.String())
}
func ne_b_sx(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.printrep != input2.String())
}
func gt_b_sx(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.printrep > input2.String())
}
func ge_b_sx(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.printrep >= input2.String())
}
func lt_b_sx(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.printrep < input2.String())
}
func le_b_sx(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.printrep <= input2.String())
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_ii(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.intval == input2.intval)
}
func ne_b_ii(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.intval != input2.intval)
}
func gt_b_ii(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.intval > input2.intval)
}
func ge_b_ii(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.intval >= input2.intval)
}
func lt_b_ii(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.intval < input2.intval)
}
func le_b_ii(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.intval <= input2.intval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_if(output, input1, input2 *Mlrval) {
	output.SetFromBool(float64(input1.intval) == input2.floatval)
}
func ne_b_if(output, input1, input2 *Mlrval) {
	output.SetFromBool(float64(input1.intval) != input2.floatval)
}
func gt_b_if(output, input1, input2 *Mlrval) {
	output.SetFromBool(float64(input1.intval) > input2.floatval)
}
func ge_b_if(output, input1, input2 *Mlrval) {
	output.SetFromBool(float64(input1.intval) >= input2.floatval)
}
func lt_b_if(output, input1, input2 *Mlrval) {
	output.SetFromBool(float64(input1.intval) < input2.floatval)
}
func le_b_if(output, input1, input2 *Mlrval) {
	output.SetFromBool(float64(input1.intval) <= input2.floatval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_fi(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.floatval == float64(input2.intval))
}
func ne_b_fi(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.floatval != float64(input2.intval))
}
func gt_b_fi(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.floatval > float64(input2.intval))
}
func ge_b_fi(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.floatval >= float64(input2.intval))
}
func lt_b_fi(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.floatval < float64(input2.intval))
}
func le_b_fi(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.floatval <= float64(input2.intval))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_ff(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.floatval == input2.floatval)
}
func ne_b_ff(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.floatval != input2.floatval)
}
func gt_b_ff(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.floatval > input2.floatval)
}
func ge_b_ff(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.floatval >= input2.floatval)
}
func lt_b_ff(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.floatval < input2.floatval)
}
func le_b_ff(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.floatval <= input2.floatval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_bb(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.boolval == input2.boolval)
}
func ne_b_bb(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.boolval != input2.boolval)
}

// We could say ordering on bool is error, but, Miller allows
// sorting on bool so it should allow ordering on bool.

func gt_b_bb(output, input1, input2 *Mlrval) {
	output.SetFromBool(lib.BoolToInt(input1.boolval) > lib.BoolToInt(input2.boolval))
}
func ge_b_bb(output, input1, input2 *Mlrval) {
	output.SetFromBool(lib.BoolToInt(input1.boolval) >= lib.BoolToInt(input2.boolval))
}
func lt_b_bb(output, input1, input2 *Mlrval) {
	output.SetFromBool(lib.BoolToInt(input1.boolval) < lib.BoolToInt(input2.boolval))
}
func le_b_bb(output, input1, input2 *Mlrval) {
	output.SetFromBool(lib.BoolToInt(input1.boolval) <= lib.BoolToInt(input2.boolval))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_aa(output, input1, input2 *Mlrval) {
	a := input1.arrayval
	b := input2.arrayval

	// Different-length arrays are not equal
	if len(a) != len(b) {
		output.SetFromBool(false) // stub
		return
	}

	// Same-length arrays: return false if any slot is not equal, else true.
	for i := range a {
		eq := MlrvalFromError()
		MlrvalEquals(&eq, &a[i], &b[i])
		lib.InternalCodingErrorIf(eq.mvtype != MT_BOOL)
		if eq.boolval == false {
			output.SetFromBool(false)
			return
		}
	}

	output.SetFromBool(true)
}
func ne_b_aa(output, input1, input2 *Mlrval) {
	eq_b_aa(output, input1, input2)
	output.SetFromBool(!output.boolval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_mm(output, input1, input2 *Mlrval) {
	output.SetFromBool(input1.mapval.Equals(input2.mapval))
}
func ne_b_mm(output, input1, input2 *Mlrval) {
	output.SetFromBool(!input1.mapval.Equals(input2.mapval))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//var eq_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
//	//       .  ERROR   ABSENT VOID     STRING   INT      FLOAT    BOOL     ARRAY    MAP
//	/*ERROR  */ {_erro, _erro, _erro,   _erro,   _erro,   _erro,   _erro,   _erro,   _erro},
//	/*ABSENT */ {_erro, _absn, _absn,   _absn,   _absn,   _absn,   _absn,   _absn,   _absn},
//	/*VOID   */ {_erro, _absn, eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _fals,   _fals,   _fals},
//	/*STRING */ {_erro, _absn, eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _fals,   _fals,   _fals},
//	/*INT    */ {_erro, _absn, eq_b_xs, eq_b_xs, eq_b_ii, eq_b_if, _fals,   _fals,   _fals},
//	/*FLOAT  */ {_erro, _absn, eq_b_xs, eq_b_xs, eq_b_fi, eq_b_ff, _fals,   _fals,   _fals},
//	/*BOOL   */ {_erro, _absn, _fals,   _fals,   _fals,   _fals,   eq_b_bb, _fals,   _fals},
//	/*ARRAY  */ {_erro, _absn, _fals,   _fals,   _fals,   _fals,   _fals,   eq_b_aa, _fals},
//	/*MAP    */ {_erro, _absn, _fals,   _fals,   _fals,   _fals,   _fals,   _fals,   eq_b_mm},
//}

var eq_dispositions = [MT_DIM][MT_DIM]BinaryFunc{}

func init() {
	eq_dispositions = [MT_DIM][MT_DIM]BinaryFunc{

		//       .  ERROR   ABSENT VOID     STRING   INT      FLOAT    BOOL   ARRAY    MAP
		/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
		/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
		/*VOID   */ {_erro, _absn, eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _fals, _fals, _fals},
		/*STRING */ {_erro, _absn, eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _fals, _fals, _fals},
		/*INT    */ {_erro, _absn, eq_b_xs, eq_b_xs, eq_b_ii, eq_b_if, _fals, _fals, _fals},
		/*FLOAT  */ {_erro, _absn, eq_b_xs, eq_b_xs, eq_b_fi, eq_b_ff, _fals, _fals, _fals},
		/*BOOL   */ {_erro, _absn, _fals, _fals, _fals, _fals, eq_b_bb, _fals, _fals},
		/*ARRAY  */ {_erro, _absn, _fals, _fals, _fals, _fals, _fals, eq_b_aa, _fals},
		/*MAP    */ {_erro, _absn, _fals, _fals, _fals, _fals, _fals, _fals, eq_b_mm},
	}

}

var ne_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT VOID     STRING   INT      FLOAT    BOOL   ARRAY    MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*VOID   */ {_erro, _absn, ne_b_ss, ne_b_ss, ne_b_sx, ne_b_sx, _true, _true, _true},
	/*STRING */ {_erro, _absn, ne_b_ss, ne_b_ss, ne_b_sx, ne_b_sx, _true, _true, _true},
	/*INT    */ {_erro, _absn, ne_b_xs, ne_b_xs, ne_b_ii, ne_b_if, _true, _true, _true},
	/*FLOAT  */ {_erro, _absn, ne_b_xs, ne_b_xs, ne_b_fi, ne_b_ff, _true, _true, _true},
	/*BOOL   */ {_erro, _absn, _true, _true, _true, _true, ne_b_bb, _true, _true},
	/*ARRAY  */ {_erro, _absn, _true, _true, _true, _true, _true, ne_b_aa, _true},
	/*MAP    */ {_erro, _absn, _true, _true, _true, _true, _true, _true, ne_b_mm},
}

var gt_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT VOID     STRING   INT      FLOAT    BOOL     ARRAY  MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*VOID   */ {_erro, _absn, gt_b_ss, gt_b_ss, gt_b_sx, gt_b_sx, _fals, _fals, _fals},
	/*STRING */ {_erro, _absn, gt_b_ss, gt_b_ss, gt_b_sx, gt_b_sx, _fals, _fals, _fals},
	/*INT    */ {_erro, _absn, gt_b_xs, gt_b_xs, gt_b_ii, gt_b_if, _fals, _fals, _fals},
	/*FLOAT  */ {_erro, _absn, gt_b_xs, gt_b_xs, gt_b_fi, gt_b_ff, _fals, _fals, _fals},
	/*BOOL   */ {_erro, _absn, _fals, _fals, _fals, _fals, gt_b_bb, _fals, _fals},
	/*ARRAY  */ {_erro, _absn, _fals, _fals, _fals, _fals, _fals, _erro, _fals},
	/*MAP    */ {_erro, _absn, _fals, _fals, _fals, _fals, _fals, _fals, _erro},
}

var ge_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT VOID     STRING   INT      FLOAT    BOOL   ARRAY  MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*VOID   */ {_erro, _absn, ge_b_ss, ge_b_ss, ge_b_sx, ge_b_sx, _fals, _fals, _fals},
	/*STRING */ {_erro, _absn, ge_b_ss, ge_b_ss, ge_b_sx, ge_b_sx, _fals, _fals, _fals},
	/*INT    */ {_erro, _absn, ge_b_xs, ge_b_xs, ge_b_ii, ge_b_if, _fals, _fals, _fals},
	/*FLOAT  */ {_erro, _absn, ge_b_xs, ge_b_xs, ge_b_fi, ge_b_ff, _fals, _fals, _fals},
	/*BOOL   */ {_erro, _absn, _fals, _fals, _fals, _fals, ge_b_bb, _fals, _fals},
	/*ARRAY  */ {_erro, _absn, _fals, _fals, _fals, _fals, _fals, _erro, _fals},
	/*MAP    */ {_erro, _absn, _fals, _fals, _fals, _fals, _fals, _fals, _erro},
}

var lt_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT VOID     STRING   INT      FLOAT    BOOL     ARRAY  MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*VOID   */ {_erro, _absn, lt_b_ss, lt_b_ss, lt_b_sx, lt_b_sx, _fals, _fals, _fals},
	/*STRING */ {_erro, _absn, lt_b_ss, lt_b_ss, lt_b_sx, lt_b_sx, _fals, _fals, _fals},
	/*INT    */ {_erro, _absn, lt_b_xs, lt_b_xs, lt_b_ii, lt_b_if, _fals, _fals, _fals},
	/*FLOAT  */ {_erro, _absn, lt_b_xs, lt_b_xs, lt_b_fi, lt_b_ff, _fals, _fals, _fals},
	/*BOOL   */ {_erro, _absn, _fals, _fals, _fals, _fals, lt_b_bb, _fals, _fals},
	/*ARRAY  */ {_erro, _absn, _fals, _fals, _fals, _fals, _fals, _erro, _fals},
	/*MAP    */ {_erro, _absn, _fals, _fals, _fals, _fals, _fals, _fals, _erro},
}

var le_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT VOID     STRING   INT      FLOAT    BOOL     ARRAY  MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*VOID   */ {_erro, _absn, le_b_ss, le_b_ss, le_b_sx, le_b_sx, _fals, _fals, _fals},
	/*STRING */ {_erro, _absn, le_b_ss, le_b_ss, le_b_sx, le_b_sx, _fals, _fals, _fals},
	/*INT    */ {_erro, _absn, le_b_xs, le_b_xs, le_b_ii, le_b_if, _fals, _fals, _fals},
	/*FLOAT  */ {_erro, _absn, le_b_xs, le_b_xs, le_b_fi, le_b_ff, _fals, _fals, _fals},
	/*BOOL   */ {_erro, _absn, _fals, _fals, _fals, _fals, le_b_bb, _fals, _fals},
	/*ARRAY  */ {_erro, _absn, _fals, _fals, _fals, _fals, _fals, _erro, _fals},
	/*MAP    */ {_erro, _absn, _fals, _fals, _fals, _fals, _fals, _fals, _erro},
}

func MlrvalEquals(output, input1, input2 *Mlrval) {
	eq_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}
func MlrvalNotEquals(output, input1, input2 *Mlrval) {
	ne_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}
func MlrvalGreaterThan(output, input1, input2 *Mlrval) {
	gt_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}
func MlrvalGreaterThanOrEquals(output, input1, input2 *Mlrval) {
	ge_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}
func MlrvalLessThan(output, input1, input2 *Mlrval) {
	lt_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}
func MlrvalLessThanOrEquals(output, input1, input2 *Mlrval) {
	le_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}

// For Go's sort.Slice
func MlrvalLessThanForSort(input1, input2 *Mlrval) bool {
	// xxx refactor to avoid copy
	mretval := MlrvalFromFalse()
	lt_dispositions[input1.mvtype][input2.mvtype](&mretval, input1, input2)
	retval, ok := mretval.GetBoolValue()
	lib.InternalCodingErrorIf(!ok)
	return retval
}

// ----------------------------------------------------------------
func MlrvalLogicalAND(output, input1, input2 *Mlrval) {
	if input1.mvtype == MT_BOOL && input2.mvtype == MT_BOOL {
		output.SetFromBool(input1.boolval && input2.boolval)
	} else {
		output.SetFromError()
	}
}

func MlrvalLogicalOR(output, input1, input2 *Mlrval) {
	if input1.mvtype == MT_BOOL && input2.mvtype == MT_BOOL {
		output.SetFromBool(input1.boolval || input2.boolval)
	} else {
		output.SetFromError()
	}
}

func MlrvalLogicalXOR(output, input1, input2 *Mlrval) {
	if input1.mvtype == MT_BOOL && input2.mvtype == MT_BOOL {
		output.SetFromBool(input1.boolval != input2.boolval)
	} else {
		output.SetFromError()
	}
}
