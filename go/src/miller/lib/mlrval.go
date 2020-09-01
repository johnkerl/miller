package lib

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

// ================================================================
// Requirements for mlrvals:
//
// * Keep original string-formatting even if parseable/parsed as int
//   o E.g. if 005 (octal), pass through as 005 unless math is done on it
//   o Likewise with number of decimal places -- 7.4 not 7.400 or (worse) 7.399999999
//
// * Invalidate the string-formatting as the output of a computational result
//
// * Have number-to-string formatting methods in the API/DSL which stick the string format
//
// * Final to-string method
//
// Also:
//
// * Split current C mvfuncs into mlrval-private (dispo matrices etc) and new
//   mvfuncs.go where the latter don't need access to private members
// ================================================================

// ================================================================
// There are two kinds of null: ABSENT (key not present in a record) and VOID
// (key present with empty value).  Note void is an acceptable string (empty
// string) but not an acceptable number. (In Javascript, similarly, there are
// undefined and null, respectively.)
// ================================================================

// ================================================================
type MVType int

const (
	// E.g. error encountered in one eval & it propagates up the AST at evaluation time:
	MT_ERROR MVType = 0

	// Key not present in input record, e.g. 'foo = $nosuchkey'
	MT_ABSENT = 1

	// Key present in input record with empty value, e.g. input data '$x=,$y=2'
	MT_VOID = 2

	MT_STRING = 3

	MT_INT = 4

	MT_FLOAT = 5

	MT_BOOL = 6

	// Not a type -- this is a dimension for disposition matrices
	MT_DIM = 7
)

// ================================================================
type Mlrval struct {
	mvtype        MVType
	printrep      string
	printrepValid bool
	intval        int64
	floatval      float64
	boolval       bool
}

// ================================================================
func MlrvalFromError() Mlrval {
	return Mlrval{
		MT_ERROR,
		"(error)", // xxx const somewhere
		true,
		0, 0.0, false,
	}
}

func MlrvalFromAbsent() Mlrval {
	return Mlrval{
		MT_ABSENT,
		"(absent)",
		true,
		0, 0.0, false,
	}
}

func MlrvalFromVoid() Mlrval {
	return Mlrval{
		MT_VOID,
		"(void)",
		true,
		0, 0.0, false,
	}
}

// ----------------------------------------------------------------
func MlrvalFromString(input string) Mlrval {
	return Mlrval{
		MT_STRING,
		input,
		true,
		0, 0.0, false,
	}
}

// ----------------------------------------------------------------
// xxx comment why two -- one for from parsed user data; other for from math ops
func MlrvalFromInt64String(input string) Mlrval {
	ival, ok := tryInt64FromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	if !ok {
		// xxx get file/line info here .......
		fmt.Fprintf(os.Stderr, "Internal coding error detected\n")
		os.Exit(1)
	}
	return Mlrval{
		MT_INT,
		input,
		true,
		ival,
		0.0,
		false,
	}
}

func MlrvalFromInt64(input int64) Mlrval {
	return Mlrval{
		MT_INT,
		"(bug-if-you-see-this)",
		false,
		input,
		0.0,
		false,
	}
}

func tryInt64FromString(input string) (int64, bool) {
	// xxx need to handle octal, hex, ......
	ival, err := strconv.ParseInt(input, 10, 64)
	if err == nil {
		return ival, true
	} else {
		return 0, false
	}
}

// ----------------------------------------------------------------
// xxx comment why two -- one for from parsed user data; other for from math ops
// xxx comment assummption is input-string already deemed parseable so no error return

func MlrvalFromFloat64String(input string) Mlrval {
	fval, ok := tryFloat64FromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	if !ok {
		// xxx get file/line info here .......
		fmt.Fprintf(os.Stderr, "Internal coding error detected\n")
		os.Exit(1)
	}
	return Mlrval{
		MT_FLOAT,
		input,
		true,
		0,
		fval,
		false,
	}
}

func MlrvalFromFloat64(input float64) Mlrval {
	return Mlrval{
		MT_FLOAT,
		"(bug-if-you-see-this)",
		false,
		0,
		input,
		false,
	}
}

func tryFloat64FromString(input string) (float64, bool) {
	ival, err := strconv.ParseFloat(input, 64)
	if err == nil {
		return ival, true
	} else {
		return 0, false
	}
}

// ----------------------------------------------------------------
func MlrvalFromTrue() Mlrval {
	return Mlrval{
		MT_BOOL,
		"true",
		true,
		0,
		0.0,
		true,
	}
}

func MlrvalFromFalse() Mlrval {
	return Mlrval{
		MT_BOOL,
		"false",
		true,
		0,
		0.0,
		false,
	}
}

func MlrvalFromBoolString(input string) Mlrval {
	if input == "true" {
		return MlrvalFromTrue()
	} else {
		return MlrvalFromFalse()
	}
	// else panic
}

func tryBoolFromBoolString(input string) (bool, bool) {
	if input == "true" {
		return true, true
	} else if input == "false" {
		return false, true
	} else {
		return false, false
	}
}

// ----------------------------------------------------------------
func MlrvalFromInferredType(input string) Mlrval {
	// xxx the parsing has happened so stash it ...
	// xxx emphasize the invariant that a non-invalid printrep always
	// matches the nval ...
	_, iok := tryInt64FromString(input)
	if iok {
		return MlrvalFromInt64String(input)
	}

	_, fok := tryFloat64FromString(input)
	if fok {
		return MlrvalFromFloat64String(input)
	}

	_, bok := tryBoolFromBoolString(input)
	if bok {
		return MlrvalFromBoolString(input)
	}

	return MlrvalFromString(input)
}

// ================================================================
// xxx comment about JIT-parsing of string backings
func (this *Mlrval) setPrintRep() {
	if !this.printrepValid {
		// xxx do it -- disposition vector
		// xxx temp temp temp temp temp
		switch this.mvtype {
		case MT_ERROR:
			this.printrep = "(error)" // xxx constdef at top of file
			break
		case MT_ABSENT:
			// Callsites should be using absence to do non-assigns, so flag
			// this clearly visually if it should (buggily) slip through to
			// user-level visibility.
			this.printrep = "(bug-if-you-see-this)" // xxx constdef at top of file
			break
		case MT_VOID:
			this.printrep = "" // xxx constdef at top of file
			break
		case MT_STRING:
			// panic i suppose
			break
		case MT_INT:
			this.printrep = strconv.FormatInt(this.intval, 10)
			break
		case MT_FLOAT:
			// xxx temp -- OFMT etc ...
			this.printrep = strconv.FormatFloat(this.floatval, 'g', -1, 64)
			break
		case MT_BOOL:
			if this.boolval == true {
				this.printrep = "true"
			} else {
				this.printrep = "false"
			}
			break
		}
		this.printrepValid = true
	}
}

func (this *Mlrval) String() string {
	this.setPrintRep()
	return this.printrep
}

// For JSON output. Second return value is true if the mlrval should be
// double-quoted.
func (this *Mlrval) StringWithQuoteInfo() (string, bool) {
	this.setPrintRep()
	quoteless := (this.mvtype == MT_INT || this.mvtype == MT_FLOAT || this.mvtype == MT_BOOL)
	return this.printrep, !quoteless
}

// ================================================================
func (this *Mlrval) IsAbsent() bool {
	return this.mvtype == MT_ABSENT
}

// ================================================================
// xxx comment why short names
func _erro(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromError()
}
func _absn(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromAbsent()
}
func _void(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromVoid()
}

func _1___(val1, val2 *Mlrval) Mlrval {
	return *val1
}
func _2___(val1, val2 *Mlrval) Mlrval {
	return *val2
}

func _s1__(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromString(val1.String())
}
func _s2__(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromString(val2.String())
}

func _i0__(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromInt64(0)
}
func _f0__(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(0.0)
}

// xxx comment
type dyadicFunc func(*Mlrval, *Mlrval) Mlrval

// ================================================================
func dot_s_xx(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromString(val1.String() + val2.String())
}

var dotDispositions = [MT_DIM][MT_DIM]dyadicFunc{
	//       ERROR ABSENT  EMPTY  STRING INT       FLOAT     BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _void, _2___, _s2__, _s2__, _s2__},
	/*EMPTY  */ {_erro, _void, _void, _2___, _s2__, _s2__, _s2__},
	/*STRING */ {_erro, _1___, _1___, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*INT    */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*FLOAT  */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*BOOL   */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
}

func MlrvalDot(val1, val2 *Mlrval) Mlrval {
	return dotDispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ================================================================
func plus_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval + float64(val2.intval))
}
func plus_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) + val2.floatval)
}
func plus_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval + val2.floatval)
}

// Auto-overflows up to float.  Additions & subtractions overflow by at most
// one bit so it suffices to check sign-changes.
func plus_n_ii(val1, val2 *Mlrval) Mlrval {
	a := val1.intval
	b := val2.intval
	c := a + b

	overflowed := false
	if a > 0 {
		if b > 0 && c < 0 {
			overflowed = true
		}
	} else if a < 0 {
		if b < 0 && c > 0 {
			overflowed = true
		}
	}

	if overflowed {
		return MlrvalFromFloat64(float64(a) + float64(b))
	} else {
		return MlrvalFromInt64(c)
	}
}

var plusDispositions = [MT_DIM][MT_DIM]dyadicFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, plus_n_ii, plus_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, plus_f_fi, plus_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalPlus(val1, val2 *Mlrval) Mlrval {
	return plusDispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ================================================================
func minus_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval - val2.floatval)
}
func minus_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval - float64(val2.intval))
}
func minus_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) - val2.floatval)
}

// Adds & subtracts overflow by at most one bit so it suffices to check
// sign-changes.
func minus_n_ii(val1, val2 *Mlrval) Mlrval {
	a := val1.intval
	b := val2.intval
	c := a - b

	overflowed := false
	if a > 0 {
		if b < 0 && c < 0 {
			overflowed = true
		}
	} else if a < 0 {
		if b > 0 && c > 0 {
			overflowed = true
		}
	}

	if overflowed {
		return MlrvalFromFloat64(float64(a) - float64(b))
	} else {
		return MlrvalFromInt64(c)
	}
}

var minusDispositions = [MT_DIM][MT_DIM]dyadicFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, minus_n_ii, minus_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, minus_f_fi, minus_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalMinus(val1, val2 *Mlrval) Mlrval {
	return minusDispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ================================================================
func times_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval * float64(val2.intval))
}
func times_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) * val2.floatval)
}
func times_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval * val2.floatval)
}

// Auto-overflows up to float.
//
// Unlike adds & subtracts which overflow by at most one bit, multiplies can
// overflow by a word size. Thus detecting sign-changes does not suffice to
// detect overflow. Instead we test whether the floating-point product exceeds
// the representable integer range. Now 64-bit integers have 64-bit precision
// while IEEE-doubles have only 52-bit mantissas -- so, 53 bits including
// implicit leading one.
//
// The following experiment explicitly demonstrates the resolution at this range:
//
//    64-bit integer     64-bit integer     Casted to double           Back to 64-bit
//        in hex           in decimal                                    integer
// 0x7ffffffffffff9ff 9223372036854774271 9223372036854773760.000000 0x7ffffffffffff800
// 0x7ffffffffffffa00 9223372036854774272 9223372036854773760.000000 0x7ffffffffffff800
// 0x7ffffffffffffbff 9223372036854774783 9223372036854774784.000000 0x7ffffffffffffc00
// 0x7ffffffffffffc00 9223372036854774784 9223372036854774784.000000 0x7ffffffffffffc00
// 0x7ffffffffffffdff 9223372036854775295 9223372036854774784.000000 0x7ffffffffffffc00
// 0x7ffffffffffffe00 9223372036854775296 9223372036854775808.000000 0x8000000000000000
// 0x7ffffffffffffffe 9223372036854775806 9223372036854775808.000000 0x8000000000000000
// 0x7fffffffffffffff 9223372036854775807 9223372036854775808.000000 0x8000000000000000
//
// That is, we cannot check an integer product to see if it is greater than
// 2**63-1 (or is less than -2**63) using integer arithmetic (it may have
// already overflowed) *or* using double-precision (granularity). Instead we
// check if the absolute value of the product exceeds the largest representable
// double less than 2**63. (An alterative would be to do all integer multiplies
// using handcrafted multi-word 128-bit arithmetic).

func times_n_ii(val1, val2 *Mlrval) Mlrval {
	a := val1.intval
	b := val2.intval
	c := float64(a) * float64(b)

	if math.Abs(c) > 9223372036854774784.0 {
		return MlrvalFromFloat64(c)
	} else {
		return MlrvalFromInt64(a * b)
	}
}

var timesDispositions = [MT_DIM][MT_DIM]dyadicFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, times_n_ii, times_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, times_f_fi, times_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalTimes(val1, val2 *Mlrval) Mlrval {
	return timesDispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ================================================================
func divide_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval / float64(val2.intval))
}
func divide_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) / val2.floatval)
}
func divide_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval / val2.floatval)
}

func divide_n_ii(val1, val2 *Mlrval) Mlrval {
	a := val1.intval
	b := val2.intval

	if b == 0 {
		// Compute inf/nan as with floats rather than fatal runtime FPE on integer divide by zero
		return MlrvalFromFloat64(float64(a) / float64(b))
	}

	// Pythonic division, not C division.
	if a%b == 0 {
		return MlrvalFromInt64(a / b)
	} else {
		return MlrvalFromFloat64(float64(a) / float64(b))
	}

	c := float64(a) * float64(b)

	if math.Abs(c) > 9223372036854774784.0 {
		return MlrvalFromFloat64(c)
	} else {
		return MlrvalFromInt64(a * b)
	}
}

var divideDispositions = [MT_DIM][MT_DIM]dyadicFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, divide_n_ii, divide_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, divide_f_fi, divide_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalDivide(val1, val2 *Mlrval) Mlrval {
	return divideDispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ================================================================
func int_divide_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Floor(val1.floatval / float64(val2.intval)))
}
func int_divide_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Floor(float64(val1.intval) / val2.floatval))
}
func int_divide_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Floor(val1.floatval / val2.floatval))
}

func int_divide_n_ii(val1, val2 *Mlrval) Mlrval {
	a := val1.intval
	b := val2.intval

	if b == 0 {
		// Compute inf/nan as with floats rather than fatal runtime FPE on integer divide by zero
		return MlrvalFromFloat64(float64(a) / float64(b))
	}

	// Pythonic division, not C division.
	q := a / b
	r := a % b
	if a < 0 {
		if b > 0 {
			if r != 0 {
				q--
			}
		}
	} else {
		if b < 0 {
			if r != 0 {
				q--
			}
		}
	}
	return MlrvalFromInt64(q)
}

var int_divideDispositions = [MT_DIM][MT_DIM]dyadicFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, int_divide_n_ii, int_divide_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, int_divide_f_fi, int_divide_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalIntDivide(val1, val2 *Mlrval) Mlrval {
	return int_divideDispositions[val1.mvtype][val2.mvtype](val1, val2)
}

//// ================================================================
//static mv_t oplus_f_ff(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = pb->u.fltv;
//	return mv_from_float(a + b);
//}
//static mv_t oplus_f_fi(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = (double)pb->u.intv;
//	return mv_from_float(a + b);
//}
//static mv_t oplus_f_if(mv_t* pa, mv_t* pb) {
//	double a = (double)pa->u.intv;
//	double b = pb->u.fltv;
//	return mv_from_float(a + b);
//}
//static mv_t oplus_n_ii(mv_t* pa, mv_t* pb) {
//	long long a = pa->u.intv;
//	long long b = pb->u.intv;
//	long long c = a + b;
//	return mv_from_int(c);
//}
//
//static mv_binary_func_t* oplus_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT         FLOAT       BOOL
//	/*ERROR*/  {_erro, _erro,  _erro, _erro,  _erro,       _erro,       _erro},
//	/*ABSENT*/ {_erro, _a,    _a,   _erro,  _2,         _2,         _erro},
//	/*EMPTY*/  {_erro, _a,    _void, _erro,  _void,       _void,       _erro},
//	/*STRING*/ {_erro, _erro,  _erro, _erro,  _erro,       _erro,       _erro},
//	/*INT*/    {_erro, _1,    _void, _erro,  oplus_n_ii, oplus_f_if, _erro},
//	/*FLOAT*/  {_erro, _1,    _void, _erro,  oplus_f_fi, oplus_f_ff, _erro},
//	/*BOOL*/   {_erro, _erro,  _erro, _erro,  _erro,       _erro,       _erro},
//};
//
//mv_t x_xx_oplus_func(mv_t* pval1, mv_t* pval2) { return (oplus_dispositions[pval1->type][pval2->type])(pval1,pval2); }
//
//// ----------------------------------------------------------------
//static mv_t ominus_f_ff(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = pb->u.fltv;
//	return mv_from_float(a - b);
//}
//static mv_t ominus_f_fi(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = (double)pb->u.intv;
//	return mv_from_float(a - b);
//}
//static mv_t ominus_f_if(mv_t* pa, mv_t* pb) {
//	double a = (double)pa->u.intv;
//	double b = pb->u.fltv;
//	return mv_from_float(a - b);
//}
//static mv_t ominus_n_ii(mv_t* pa, mv_t* pb) {
//	long long a = pa->u.intv;
//	long long b = pb->u.intv;
//	long long c = a - b;
//	return mv_from_int(c);
//}
//
//static mv_binary_func_t* ominus_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT          FLOAT        BOOL
//	/*ERROR*/  {_erro, _erro,  _erro, _erro,  _erro,        _erro,        _erro},
//	/*ABSENT*/ {_erro, _a,    _a,   _erro,  _2,          _2,          _erro},
//	/*EMPTY*/  {_erro, _a,    _void, _erro,  _void,        _void,        _erro},
//	/*STRING*/ {_erro, _erro,  _erro, _erro,  _erro,        _erro,        _erro},
//	/*INT*/    {_erro, _1,    _void, _erro,  ominus_n_ii, ominus_f_if, _erro},
//	/*FLOAT*/  {_erro, _1,    _void, _erro,  ominus_f_fi, ominus_f_ff, _erro},
//	/*BOOL*/   {_erro, _erro,  _erro, _erro,  _erro,        _erro,        _erro},
//};
//
//mv_t x_xx_ominus_func(mv_t* pval1, mv_t* pval2) { return (ominus_dispositions[pval1->type][pval2->type])(pval1,pval2); }
//
//// ----------------------------------------------------------------
//static mv_t otimes_f_ff(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = pb->u.fltv;
//	return mv_from_float(a * b);
//}
//static mv_t otimes_f_fi(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = (double)pb->u.intv;
//	return mv_from_float(a * b);
//}
//static mv_t otimes_f_if(mv_t* pa, mv_t* pb) {
//	double a = (double)pa->u.intv;
//	double b = pb->u.fltv;
//	return mv_from_float(a * b);
//}
//static mv_t otimes_n_ii(mv_t* pa, mv_t* pb) {
//	long long a = pa->u.intv;
//	long long b = pb->u.intv;
//	return mv_from_int(a * b);
//}
//
//static mv_binary_func_t* otimes_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT          FLOAT       BOOL
//	/*ERROR*/  {_erro, _erro,  _erro, _erro,  _erro,        _erro,       _erro},
//	/*ABSENT*/ {_erro, _a,    _a,   _erro,  _2,          _2,         _erro},
//	/*EMPTY*/  {_erro, _a,    _void, _erro,  _void,        _void,       _erro},
//	/*STRING*/ {_erro, _erro,  _erro, _erro,  _erro,        _erro,       _erro},
//	/*INT*/    {_erro, _1,    _void, _erro,  otimes_n_ii, otimes_f_if, _erro},
//	/*FLOAT*/  {_erro, _1,    _void, _erro,  otimes_f_fi, otimes_f_ff, _erro},
//	/*BOOL*/   {_erro, _erro,  _erro, _erro,  _erro,        _erro,       _erro},
//};
//
//mv_t x_xx_otimes_func(mv_t* pval1, mv_t* pval2) { return (otimes_dispositions[pval1->type][pval2->type])(pval1,pval2); }
//
//// ----------------------------------------------------------------
//static mv_t odivide_f_ff(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = pb->u.fltv;
//	return mv_from_float(a / b);
//}
//static mv_t odivide_f_fi(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = (double)pb->u.intv;
//	return mv_from_float(a / b);
//}
//static mv_t odivide_f_if(mv_t* pa, mv_t* pb) {
//	double a = (double)pa->u.intv;
//	double b = pb->u.fltv;
//	return mv_from_float(a / b);
//}
//static mv_t odivide_i_ii(mv_t* pa, mv_t* pb) {
//	long long a = pa->u.intv;
//	long long b = pb->u.intv;
//	return mv_from_int(a / b);
//}
//
//static mv_binary_func_t* odivide_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT           FLOAT         BOOL
//	/*ERROR*/  {_erro, _erro,  _erro, _erro,  _erro,         _erro,         _erro},
//	/*ABSENT*/ {_erro, _a,    _a,   _erro,  _i0,          _f0,          _erro},
//	/*EMPTY*/  {_erro, _a,    _void, _erro,  _void,         _void,         _erro},
//	/*STRING*/ {_erro, _erro,  _erro, _erro,  _erro,         _erro,         _erro},
//	/*INT*/    {_erro, _1,    _void, _erro,  odivide_i_ii, odivide_f_if, _erro},
//	/*FLOAT*/  {_erro, _1,    _void, _erro,  odivide_f_fi, odivide_f_ff, _erro},
//	/*BOOL*/   {_erro, _erro,  _erro, _erro,  _erro,         _erro,         _erro},
//};
//
//mv_t x_xx_odivide_func(mv_t* pval1, mv_t* pval2) { return (odivide_dispositions[pval1->type][pval2->type])(pval1,pval2); }
//
//// ----------------------------------------------------------------
//static mv_t oidiv_f_ff(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = pb->u.fltv;
//	return mv_from_float(floor(a / b));
//}
//static mv_t oidiv_f_fi(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = (double)pb->u.intv;
//	return mv_from_float(floor(a / b));
//}
//static mv_t oidiv_f_if(mv_t* pa, mv_t* pb) {
//	double a = (double)pa->u.intv;
//	double b = pb->u.fltv;
//	return mv_from_float(floor(a / b));
//}
//static mv_t oidiv_i_ii(mv_t* pa, mv_t* pb) {
//	long long a = pa->u.intv;
//	long long b = pb->u.intv;
//
//	// Pythonic division, not C division.
//	long long q = a / b;
//	long long r = a % b;
//	if (a < 0) {
//		if (b > 0) {
//			if (r != 0)
//				q--;
//		}
//	} else {
//		if (b < 0) {
//			if (r != 0)
//				q--;
//		}
//	}
//	return mv_from_int(q);
//}
//
//static mv_binary_func_t* oidiv_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT         FLOAT       BOOL
//	/*ERROR*/  {_erro, _erro,  _erro, _erro,  _erro,       _erro,       _erro},
//	/*ABSENT*/ {_erro, _a,    _a,   _erro,  _i0,        _f0,        _erro},
//	/*EMPTY*/  {_erro, _a,    _void, _erro,  _void,       _void,       _erro},
//	/*STRING*/ {_erro, _erro,  _erro, _erro,  _erro,       _erro,       _erro},
//	/*INT*/    {_erro, _1,    _void, _erro,  oidiv_i_ii, oidiv_f_if, _erro},
//	/*FLOAT*/  {_erro, _1,    _void, _erro,  oidiv_f_fi, oidiv_f_ff, _erro},
//	/*BOOL*/   {_erro, _erro,  _erro, _erro,  _erro,       _erro,       _erro},
//};
//
//mv_t x_xx_int_odivide_func(mv_t* pval1, mv_t* pval2) {
//	return (oidiv_dispositions[pval1->type][pval2->type])(pval1,pval2);
//}
