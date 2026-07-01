// Boolean expressions for ==, !=, >, >=, <, <=

package bifs

import (
	"bytes"

	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

//   - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//
// string_cmp implements the spaceship operator for strings.
func string_cmp(a, b string) int64 {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// int_cmp implements the spaceship operator for ints.
func int_cmp(a, b int64) int64 {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// float_cmp implements the spaceship operator for floats.
func float_cmp(a, b float64) int64 {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
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
	return mlrval.FromInt(int64(string_cmp(input1.AcquireStringValue(), input2.AcquireStringValue())))
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
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
	return mlrval.FromInt(int64(string_cmp(input1.String(), input2.AcquireStringValue())))
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
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

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
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

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
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

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
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

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
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

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
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

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_aa(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	a := input1.AcquireArrayValue()
	b := input2.AcquireArrayValue()

	// Different-length arrays are not equal
	if len(a) != len(b) {
		return mlrval.FALSE
	}

	// Same-length arrays: return false if any slot is not equal, else true.
	for i := range a {
		eq := BIF_equals(a[i], b[i])
		// Treat invalid comparison as false
		if eq.Type() == mlrval.MT_ABSENT {
			return mlrval.FALSE
		}
		lib.InternalCodingErrorIf(eq.Type() != mlrval.MT_BOOL)
		if !eq.AcquireBoolValue() {
			return mlrval.FALSE
		}
	}

	return mlrval.TRUE
}
func ne_b_aa(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	output := eq_b_aa(input1, input2)
	return mlrval.FromBool(!output.AcquireBoolValue())
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// 'y' is for bytes ('b' is for boolean)
func eq_b_yy(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(bytes.Equal(input1.AcquireBytesValue(), input2.AcquireBytesValue()))
}
func ne_b_yy(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(!bytes.Equal(input1.AcquireBytesValue(), input2.AcquireBytesValue()))
}
func gt_b_yy(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(bytes.Compare(input1.AcquireBytesValue(), input2.AcquireBytesValue()) > 0)
}
func ge_b_yy(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(bytes.Compare(input1.AcquireBytesValue(), input2.AcquireBytesValue()) >= 0)
}
func lt_b_yy(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(bytes.Compare(input1.AcquireBytesValue(), input2.AcquireBytesValue()) < 0)
}
func le_b_yy(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(bytes.Compare(input1.AcquireBytesValue(), input2.AcquireBytesValue()) <= 0)
}
func cmp_b_yy(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(int64(bytes.Compare(input1.AcquireBytesValue(), input2.AcquireBytesValue())))
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_mm(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireMapValue().Equals(input2.AcquireMapValue()))
}
func ne_b_mm(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(!input1.AcquireMapValue().Equals(input2.AcquireMapValue()))
}

// We get a Golang "initialization loop" due to recursive depth computation
// if this is defined statically. So, we use a "package init" function.
var eq_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{}

func eqte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary("==", input1, input2)
}

func init() {
	eq_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
		//       .  INT       FLOAT    BOOL     VOID     STRING   ARRAY    MAP      FUNC   ERROR  NULL   ABSENT BYTES
		/*INT    */ {eq_b_ii, eq_b_if, _fals, eq_b_xs, eq_b_xs, _fals, _fals, eqte, eqte, _fals, _absn, _fals},
		/*FLOAT  */ {eq_b_fi, eq_b_ff, _fals, eq_b_xs, eq_b_xs, _fals, _fals, eqte, eqte, _fals, _absn, _fals},
		/*BOOL   */ {_fals, _fals, eq_b_bb, _fals, _fals, _fals, _fals, eqte, eqte, _fals, _absn, _fals},
		/*VOID   */ {eq_b_sx, eq_b_sx, _fals, eq_b_ss, eq_b_ss, _fals, _fals, eqte, eqte, _fals, _absn, _fals},
		/*STRING */ {eq_b_sx, eq_b_sx, _fals, eq_b_ss, eq_b_ss, _fals, _fals, eqte, eqte, _fals, _absn, _fals},
		/*ARRAY  */ {_fals, _fals, _fals, _fals, _fals, eq_b_aa, _fals, eqte, eqte, _fals, _absn, _fals},
		/*MAP    */ {_fals, _fals, _fals, _fals, _fals, _fals, eq_b_mm, eqte, eqte, _fals, _absn, _fals},
		/*FUNC   */ {eqte, eqte, eqte, eqte, eqte, eqte, eqte, eqte, eqte, eqte, eqte, eqte},
		/*ERROR  */ {eqte, eqte, eqte, eqte, eqte, eqte, eqte, eqte, eqte, eqte, eqte, eqte},
		/*NULL   */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, eqte, eqte, _true, _absn, _fals},
		/*ABSENT */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, eqte, eqte, _absn, _absn, _absn},
		/*BYTES  */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, eqte, eqte, _fals, _absn, eq_b_yy},
	}
}

func nete(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary("!=", input1, input2)
}

var ne_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT       FLOAT    BOOL     VOID     STRING   ARRAY    MAP      FUNC   ERROR  NULL   ABSENT BYTES
	/*INT    */ {ne_b_ii, ne_b_if, _true, ne_b_xs, ne_b_xs, _true, _true, nete, nete, _true, _absn, _true},
	/*FLOAT  */ {ne_b_fi, ne_b_ff, _true, ne_b_xs, ne_b_xs, _true, _true, nete, nete, _true, _absn, _true},
	/*BOOL   */ {_true, _true, ne_b_bb, _true, _true, _true, _true, nete, nete, _true, _absn, _true},
	/*VOID   */ {ne_b_sx, ne_b_sx, _true, ne_b_ss, ne_b_ss, _true, _true, nete, nete, _true, _absn, _true},
	/*STRING */ {ne_b_sx, ne_b_sx, _true, ne_b_ss, ne_b_ss, _true, _true, nete, nete, _true, _absn, _true},
	/*ARRAY  */ {_true, _true, _true, _true, _true, ne_b_aa, _true, nete, nete, _true, _absn, _true},
	/*MAP    */ {_true, _true, _true, _true, _true, _true, ne_b_mm, nete, nete, _true, _absn, _true},
	/*FUNC   */ {nete, nete, nete, nete, nete, nete, nete, nete, nete, nete, nete, nete},
	/*ERROR  */ {nete, nete, nete, nete, nete, nete, nete, nete, nete, nete, nete, nete},
	/*NULL   */ {_true, _true, _true, _true, _true, _true, _true, nete, nete, _fals, _absn, _true},
	/*ABSENT */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, nete, nete, _absn, _absn, _absn},
	/*BYTES  */ {_true, _true, _true, _true, _true, _true, _true, nete, nete, _true, _absn, ne_b_yy},
}

func gtte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary(">", input1, input2)
}

var gt_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT       FLOAT    BOOL     VOID     STRING   ARRAY  MAP    FUNC   ERROR  NULL   ABSENT BYTES
	/*INT    */ {gt_b_ii, gt_b_if, _fals, gt_b_xs, gt_b_xs, _fals, _fals, gtte, gtte, _fals, _absn, _fals},
	/*FLOAT  */ {gt_b_fi, gt_b_ff, _fals, gt_b_xs, gt_b_xs, _fals, _fals, gtte, gtte, _fals, _absn, _fals},
	/*BOOL   */ {_fals, _fals, gt_b_bb, _fals, _fals, _fals, _fals, gtte, gtte, _fals, _absn, _fals},
	/*VOID   */ {gt_b_sx, gt_b_sx, _fals, gt_b_ss, gt_b_ss, _fals, _fals, gtte, gtte, _fals, _absn, _fals},
	/*STRING */ {gt_b_sx, gt_b_sx, _fals, gt_b_ss, gt_b_ss, _fals, _fals, gtte, gtte, _fals, _absn, _fals},
	/*ARRAY  */ {_fals, _fals, _fals, _fals, _fals, gtte, _fals, gtte, gtte, _fals, _absn, _fals},
	/*MAP    */ {_fals, _fals, _fals, _fals, _fals, _fals, gtte, gtte, gtte, _fals, _absn, _fals},
	/*FUNC   */ {gtte, gtte, gtte, gtte, gtte, gtte, gtte, gtte, gtte, gtte, gtte, gtte},
	/*ERROR  */ {gtte, gtte, gtte, gtte, gtte, gtte, gtte, gtte, gtte, _fals, gtte, gtte},
	/*NULL   */ {_true, _true, _true, _true, _true, _absn, _absn, gtte, _true, _fals, _fals, _true},
	/*ABSENT */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, gtte, gtte, _true, _absn, _absn},
	/*BYTES  */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, gtte, gtte, _fals, _absn, gt_b_yy},
}

func gete(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary(">=", input1, input2)
}

var ge_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT       FLOAT    BOOL     VOID     STRING   ARRAY  MAP    FUNC   ERROR  NULL   ABSENT BYTES
	/*INT    */ {ge_b_ii, ge_b_if, _fals, ge_b_xs, ge_b_xs, _fals, _fals, gete, gete, _fals, _absn, _fals},
	/*FLOAT  */ {ge_b_fi, ge_b_ff, _fals, ge_b_xs, ge_b_xs, _fals, _fals, gete, gete, _fals, _absn, _fals},
	/*BOOL   */ {_fals, _fals, ge_b_bb, _fals, _fals, _fals, _fals, gete, gete, _fals, _absn, _fals},
	/*VOID   */ {ge_b_sx, ge_b_sx, _fals, ge_b_ss, ge_b_ss, _fals, _fals, gete, gete, _fals, _absn, _fals},
	/*STRING */ {ge_b_sx, ge_b_sx, _fals, ge_b_ss, ge_b_ss, _fals, _fals, gete, gete, _fals, _absn, _fals},
	/*ARRAY  */ {_fals, _fals, _fals, _fals, _fals, gete, _fals, gete, gete, _fals, _absn, _fals},
	/*MAP    */ {_fals, _fals, _fals, _fals, _fals, _fals, gete, gete, gete, _fals, _absn, _fals},
	/*FUNC   */ {gete, gete, gete, gete, gete, gete, gete, gete, gete, gete, gete, gete},
	/*ERROR  */ {gete, gete, gete, gete, gete, gete, gete, gete, gete, _fals, gete, gete},
	/*NULL   */ {_true, _true, _true, _true, _true, _absn, _absn, gete, _true, _true, _fals, _true},
	/*ABSENT */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, gete, gete, _true, _absn, _absn},
	/*BYTES  */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, gete, gete, _fals, _absn, ge_b_yy},
}

func ltte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary("<", input1, input2)
}

var lt_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT       FLOAT    BOOL     VOID     STRING   ARRAY  MAP    FUNC   ERROR  NULL   ABSENT BYTES
	/*INT    */ {lt_b_ii, lt_b_if, _fals, lt_b_xs, lt_b_xs, _fals, _fals, ltte, ltte, _true, _absn, _fals},
	/*FLOAT  */ {lt_b_fi, lt_b_ff, _fals, lt_b_xs, lt_b_xs, _fals, _fals, ltte, ltte, _true, _absn, _fals},
	/*BOOL   */ {_fals, _fals, lt_b_bb, _fals, _fals, _fals, _fals, ltte, ltte, _true, _absn, _fals},
	/*VOID   */ {lt_b_sx, lt_b_sx, _fals, lt_b_ss, lt_b_ss, _fals, _fals, ltte, ltte, _true, _absn, _fals},
	/*STRING */ {lt_b_sx, lt_b_sx, _fals, lt_b_ss, lt_b_ss, _fals, _fals, ltte, ltte, _true, _absn, _fals},
	/*ARRAY  */ {_fals, _fals, _fals, _fals, _fals, ltte, _fals, ltte, ltte, _absn, _absn, _fals},
	/*MAP    */ {_fals, _fals, _fals, _fals, _fals, _fals, ltte, ltte, ltte, _absn, _absn, _fals},
	/*FUNC   */ {ltte, ltte, ltte, ltte, ltte, ltte, ltte, ltte, ltte, ltte, ltte, ltte},
	/*ERROR  */ {ltte, ltte, ltte, ltte, ltte, ltte, ltte, ltte, ltte, _true, ltte, ltte},
	/*NULL   */ {_fals, _fals, _fals, _fals, _fals, _absn, _absn, ltte, _fals, _fals, _true, _fals},
	/*ABSENT */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, ltte, ltte, _fals, _absn, _absn},
	/*BYTES  */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, ltte, ltte, _fals, _absn, lt_b_yy},
}

func lete(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary("<=", input1, input2)
}

var le_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT       FLOAT    BOOL     VOID     STRING   ARRAY  MAP    FUNC   ERROR  NULL   ABSENT BYTES
	/*INT    */ {le_b_ii, le_b_if, _fals, le_b_xs, le_b_xs, _fals, _fals, lete, lete, _true, _absn, _fals},
	/*FLOAT  */ {le_b_fi, le_b_ff, _fals, le_b_xs, le_b_xs, _fals, _fals, lete, lete, _true, _absn, _fals},
	/*BOOL   */ {_fals, _fals, le_b_bb, _fals, _fals, _fals, _fals, lete, lete, _true, _absn, _fals},
	/*VOID   */ {le_b_sx, le_b_sx, _fals, le_b_ss, le_b_ss, _fals, _fals, lete, lete, _true, _absn, _fals},
	/*STRING */ {le_b_sx, le_b_sx, _fals, le_b_ss, le_b_ss, _fals, _fals, lete, lete, _true, _absn, _fals},
	/*ARRAY  */ {_fals, _fals, _fals, _fals, _fals, lete, _fals, lete, lete, _absn, _absn, _fals},
	/*MAP    */ {_fals, _fals, _fals, _fals, _fals, _fals, lete, lete, lete, _absn, _absn, _fals},
	/*FUNC   */ {lete, lete, lete, lete, lete, lete, lete, lete, lete, lete, lete, lete},
	/*ERROR  */ {lete, lete, lete, lete, lete, lete, lete, lete, lete, _true, lete, lete},
	/*NULL   */ {_fals, _fals, _fals, _fals, _fals, _absn, _absn, lete, _fals, _true, _true, _fals},
	/*ABSENT */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, lete, lete, _fals, _absn, _absn},
	/*BYTES  */ {_fals, _fals, _fals, _fals, _fals, _fals, _fals, lete, lete, _fals, _absn, le_b_yy},
}

func cmpte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary("<=>", input1, input2)
}

var cmp_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT        FLOAT     BOOL      VOID      STRING    ARRAY  MAP    FUNC   ERROR  NULL   ABSENT BYTES
	/*INT    */ {cmp_b_ii, cmp_b_if, _less, cmp_b_xs, cmp_b_xs, _less, _less, cmpte, cmpte, _true, _absn, _less},
	/*FLOAT  */ {cmp_b_fi, cmp_b_ff, _less, cmp_b_xs, cmp_b_xs, _less, _less, cmpte, cmpte, _true, _absn, _less},
	/*BOOL   */ {_more, _more, cmp_b_bb, _less, _less, _less, _less, cmpte, cmpte, _true, _absn, _less},
	/*VOID   */ {cmp_b_sx, cmp_b_sx, _more, cmp_b_ss, cmp_b_ss, _less, _less, cmpte, cmpte, _true, _absn, _less},
	/*STRING */ {cmp_b_sx, cmp_b_sx, _more, cmp_b_ss, cmp_b_ss, _less, _less, cmpte, cmpte, _true, _absn, _less},
	/*ARRAY  */ {_more, _more, _more, _more, _more, cmpte, _less, cmpte, cmpte, _absn, _absn, _more},
	/*MAP    */ {_more, _more, _more, _more, _more, _more, cmpte, cmpte, cmpte, _absn, _absn, _more},
	/*FUNC   */ {cmpte, cmpte, cmpte, cmpte, cmpte, cmpte, cmpte, cmpte, cmpte, cmpte, cmpte, cmpte},
	/*ERROR  */ {cmpte, cmpte, cmpte, cmpte, cmpte, cmpte, cmpte, cmpte, cmpte, _true, cmpte, cmpte},
	/*NULL   */ {_more, _more, _more, _more, _more, _absn, _absn, cmpte, _more, _same, _true, _more},
	/*ABSENT */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, cmpte, cmpte, _more, _absn, _absn},
	/*BYTES  */ {_more, _more, _more, _more, _more, _less, _less, cmpte, cmpte, _less, _absn, cmp_b_yy},
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
