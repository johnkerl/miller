package types

// ----------------------------------------------------------------
// Boolean expressions for ==, !=, >, >=, <, <=

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
var eq_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*VOID   */ {_erro, _absn, eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _erro, _absn, _absn},
	/*STRING */ {_erro, _absn, eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _erro, _absn, _absn},
	/*INT    */ {_erro, _absn, eq_b_xs, eq_b_xs, eq_b_ii, eq_b_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _absn, eq_b_xs, eq_b_xs, eq_b_fi, eq_b_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _absn, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

var ne_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*VOID   */ {_erro, _absn, ne_b_ss, ne_b_ss, ne_b_sx, ne_b_sx, _erro, _absn, _absn},
	/*STRING */ {_erro, _absn, ne_b_ss, ne_b_ss, ne_b_sx, ne_b_sx, _erro, _absn, _absn},
	/*INT    */ {_erro, _absn, ne_b_xs, ne_b_xs, ne_b_ii, ne_b_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _absn, ne_b_xs, ne_b_xs, ne_b_fi, ne_b_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _absn, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

var gt_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*VOID   */ {_erro, _absn, gt_b_ss, gt_b_ss, gt_b_sx, gt_b_sx, _erro, _absn, _absn},
	/*STRING */ {_erro, _absn, gt_b_ss, gt_b_ss, gt_b_sx, gt_b_sx, _erro, _absn, _absn},
	/*INT    */ {_erro, _absn, gt_b_xs, gt_b_xs, gt_b_ii, gt_b_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _absn, gt_b_xs, gt_b_xs, gt_b_fi, gt_b_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _absn, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

var ge_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//          ERROR   ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*VOID   */ {_erro, _absn, ge_b_ss, ge_b_ss, ge_b_sx, ge_b_sx, _erro, _absn, _absn},
	/*STRING */ {_erro, _absn, ge_b_ss, ge_b_ss, ge_b_sx, ge_b_sx, _erro, _absn, _absn},
	/*INT    */ {_erro, _absn, ge_b_xs, ge_b_xs, ge_b_ii, ge_b_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _absn, ge_b_xs, ge_b_xs, ge_b_fi, ge_b_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _absn, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

var lt_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//          ERROR   ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*VOID   */ {_erro, _absn, lt_b_ss, lt_b_ss, lt_b_sx, lt_b_sx, _erro, _absn, _absn},
	/*STRING */ {_erro, _absn, lt_b_ss, lt_b_ss, lt_b_sx, lt_b_sx, _erro, _absn, _absn},
	/*INT    */ {_erro, _absn, lt_b_xs, lt_b_xs, lt_b_ii, lt_b_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _absn, lt_b_xs, lt_b_xs, lt_b_fi, lt_b_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _absn, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

var le_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//          ERROR   ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*VOID   */ {_erro, _absn, le_b_ss, le_b_ss, le_b_sx, le_b_sx, _erro, _absn, _absn},
	/*STRING */ {_erro, _absn, le_b_ss, le_b_ss, le_b_sx, le_b_sx, _erro, _absn, _absn},
	/*INT    */ {_erro, _absn, le_b_xs, le_b_xs, le_b_ii, le_b_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _absn, le_b_xs, le_b_xs, le_b_fi, le_b_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _absn, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
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
