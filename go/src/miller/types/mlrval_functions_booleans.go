// ================================================================
// Boolean expressions for ==, !=, >, >=, <, <=
// ================================================================

package types

import (
	"miller/lib"
)

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_ss(ma *Mlrval, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep == mb.printrep)
}
func ne_b_ss(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep != mb.printrep)
}
func gt_b_ss(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep > mb.printrep)
}
func ge_b_ss(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep >= mb.printrep)
}
func lt_b_ss(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep < mb.printrep)
}
func le_b_ss(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep <= mb.printrep)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_xs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.String() == mb.printrep)
}
func ne_b_xs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.String() != mb.printrep)
}
func gt_b_xs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.String() > mb.printrep)
}
func ge_b_xs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.String() >= mb.printrep)
}
func lt_b_xs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.String() < mb.printrep)
}
func le_b_xs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.String() <= mb.printrep)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_sx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep == mb.String())
}
func ne_b_sx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep != mb.String())
}
func gt_b_sx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep > mb.String())
}
func ge_b_sx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep >= mb.String())
}
func lt_b_sx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep < mb.String())
}
func le_b_sx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep <= mb.String())
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval == mb.intval)
}
func ne_b_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval != mb.intval)
}
func gt_b_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval > mb.intval)
}
func ge_b_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval >= mb.intval)
}
func lt_b_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval < mb.intval)
}
func le_b_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval <= mb.intval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(float64(ma.intval) == mb.floatval)
}
func ne_b_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(float64(ma.intval) != mb.floatval)
}
func gt_b_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(float64(ma.intval) > mb.floatval)
}
func ge_b_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(float64(ma.intval) >= mb.floatval)
}
func lt_b_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(float64(ma.intval) < mb.floatval)
}
func le_b_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(float64(ma.intval) <= mb.floatval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval == float64(mb.intval))
}
func ne_b_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval != float64(mb.intval))
}
func gt_b_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval > float64(mb.intval))
}
func ge_b_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval >= float64(mb.intval))
}
func lt_b_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval < float64(mb.intval))
}
func le_b_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval <= float64(mb.intval))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval == mb.floatval)
}
func ne_b_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval != mb.floatval)
}
func gt_b_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval > mb.floatval)
}
func ge_b_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval >= mb.floatval)
}
func lt_b_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval < mb.floatval)
}
func le_b_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval <= mb.floatval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_bb(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.boolval == mb.boolval)
}
func ne_b_bb(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.boolval != mb.boolval)
}

// We could say ordering on bool is error, but, Miller allows
// sorting on bool so it should allow ordering on bool.

func gt_b_bb(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(lib.BoolToInt(ma.boolval) > lib.BoolToInt(mb.boolval))
}
func ge_b_bb(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(lib.BoolToInt(ma.boolval) >= lib.BoolToInt(mb.boolval))
}
func lt_b_bb(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(lib.BoolToInt(ma.boolval) < lib.BoolToInt(mb.boolval))
}
func le_b_bb(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(lib.BoolToInt(ma.boolval) <= lib.BoolToInt(mb.boolval))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_aa(ma, mb *Mlrval) Mlrval {
	a := ma.arrayval
	b := mb.arrayval

	// Different-length arrays are not equal
	if len(a) != len(b) {
		return MlrvalFromBool(false) // stub
	}

	// Same-length arrays: return false if any slot is not equal, else true.
	for i := range a {
		eq := MlrvalEquals(&a[i], &b[i])
		if eq.boolval == false {
			return MlrvalFromBool(false)
		}
	}

	return MlrvalFromBool(true)
}
func ne_b_aa(ma, mb *Mlrval) Mlrval {
	eq := eq_b_aa(ma, mb)
	return MlrvalFromBool(!eq.boolval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_mm(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mapval.Equals(mb.mapval))
}
func ne_b_mm(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(!ma.mapval.Equals(mb.mapval))
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

func MlrvalEquals(ma, mb *Mlrval) Mlrval {
	return eq_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}
func MlrvalNotEquals(ma, mb *Mlrval) Mlrval {
	return ne_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}
func MlrvalGreaterThan(ma, mb *Mlrval) Mlrval {
	return gt_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}
func MlrvalGreaterThanOrEquals(ma, mb *Mlrval) Mlrval {
	return ge_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}
func MlrvalLessThan(ma, mb *Mlrval) Mlrval {
	return lt_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}
func MlrvalLessThanOrEquals(ma, mb *Mlrval) Mlrval {
	return le_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// For Go's sort.Slice
func MlrvalLessThanForSort(ma, mb *Mlrval) bool {
	mretval := lt_dispositions[ma.mvtype][mb.mvtype](ma, mb)
	retval, ok := mretval.GetBoolValue()
	lib.InternalCodingErrorIf(!ok)
	return retval
}

// ----------------------------------------------------------------
func MlrvalLogicalAND(ma, mb *Mlrval) Mlrval {
	if ma.mvtype == MT_BOOL && mb.mvtype == MT_BOOL {
		return MlrvalFromBool(ma.boolval && mb.boolval)
	} else {
		return MlrvalFromError()
	}
}

func MlrvalLogicalOR(ma, mb *Mlrval) Mlrval {
	if ma.mvtype == MT_BOOL && mb.mvtype == MT_BOOL {
		return MlrvalFromBool(ma.boolval || mb.boolval)
	} else {
		return MlrvalFromError()
	}
}

func MlrvalLogicalXOR(ma, mb *Mlrval) Mlrval {
	if ma.mvtype == MT_BOOL && mb.mvtype == MT_BOOL {
		return MlrvalFromBool(ma.boolval != mb.boolval)
	} else {
		return MlrvalFromError()
	}
}
