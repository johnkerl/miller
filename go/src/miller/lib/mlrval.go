package lib

import (
	"fmt"
	"os"
	"strconv"
)

// Two kinds of null: absent (key not present in a record) and void (key
// present with empty value).  Note void is an acceptable string (empty string)
// but not an acceptable number. (In Javascript, similarly, there are null and
// undefined.) Void-valued mlrvals have u.strv = "".

// #define MT_ERROR    0 // E.g. error encountered in one eval & it propagates up the AST.

// --> JS undefined -- rename this -- or maybe MT_ABSENT ok ...
// #define MT_ABSENT   1 // No such key, e.g. $z in 'x=,y=2'

// --> JS null -- rename this -- or maybe MT_VOID
// xxx note it seralizes to "" ... which is kinda whack ...
// #define MT_EMPTY    2 // Empty value, e.g. $x in 'x=,y=2'

// #define MT_STRING   3
// #define MT_INT      4
// #define MT_FLOAT    5
// #define MT_BOOL     6
// #define MT_DIM      7

// typedef struct _mv_t {
// 	union {
// 		char*      strv;  // MT_STRING and MT_EMPTY
// 		long long  intv;  // MT_INT, and == 0 for MT_ABSENT and MT_ERROR
// 		double     fltv;  // MT_FLOAT
// 		int        boolv; // MT_BOOL
// 	} u;
// 	unsigned char type;
// } mv_t;

// Requirements:
// * Keep original string-formatting even if parseable/parsed as int
//   o E.g. if 005 (octal), pass through as 005 unless math is done on it
//   o Likewise with number of decimal places -- 7.4 not 7.400 or (worse) 7.399999999
// * Invalidate the string-formatting as the output of a computational result
// * Have number-to-string formatting methods in the API/DSL which stick the string format
// * Final to-string method

// Also:
// * split current C mvfuncs into mlrval-private (dispo matrices etc) and new
//   mvfuncs.go where the latter don't need access to private members

// Question:
// * What to do with disposition matrices from the C impl?
//   Example:
//static mv_binary_func_t* plus_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT        FLOAT      BOOL
//	/*ERROR*/  {_err, _err,  _err, _err,  _err,      _err,      _err},
//	/*ABSENT*/ {_err, _a,    _a,   _err,  _2___,        _2___,        _err},
//	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,      _emt,      _err},
//	/*STRING*/ {_err, _err,  _err, _err,  _err,      _err,      _err},
//	/*INT*/    {_err, _1___,    _emt, _err,  plus_n_ii, plus_f_if, _err},
//	/*FLOAT*/  {_err, _1___,    _emt, _err,  plus_f_fi, plus_f_ff, _err},
//	/*BOOL*/   {_err, _err,  _err, _err,  _err,      _err,      _err},
//};

// ----------------------------------------------------------------
type MVType int

const (
	MT_ERROR MVType = 0
	MT_ABSENT = 1
	MT_VOID = 2
	MT_STRING = 3
	MT_INT = 4
	MT_FLOAT = 5
	MT_BOOL = 6
	// Not a type -- this is a dimension for disposition matrices
	MT_DIM = 7
)

// ----------------------------------------------------------------
type Mlrval struct {
	mvtype        MVType
	printrep      string
	printrepValid bool
	intval        int64
	floatval      float64
	boolval       bool
}

// ----------------------------------------------------------------
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

func MlrvalFromString(input string) Mlrval {
	return Mlrval{
		MT_STRING,
		input,
		true,
		0, 0.0, false,
	}
}

// xxx comment why two -- one for from parsed user data; other for from math ops
func MlrvalFromInt64String(input string) Mlrval {
	// xxx handle octal, hex, ......
	ival, err := strconv.ParseInt(input, 10, 64)
	// xxx comment assummption is input-string already deemed parseable so no error return
	if err != nil {
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
		"(uninit)",
		false,
		input,
		0.0,
		false,
	}
}

// xxx comment why two -- one for from parsed user data; other for from math ops
// xxx comment assummption is input-string already deemed parseable so no error return
func MlrvalFromFloat64String(input string) Mlrval {
	fval, err := strconv.ParseFloat(input, 64)
	// xxx comment assummption is input-string already deemed parseable so no error return
	if err != nil {
		// xxx panic ?
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
		"(uninit)",
		false,
		0,
		input,
		false,
	}
}

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

func MlrvalFromBoolean(input bool) Mlrval {
	if input == true {
		return MlrvalFromTrue()
	} else {
		return MlrvalFromFalse()
	}
}

// ----------------------------------------------------------------
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
			this.printrep = "(absent)" // xxx constdef at top of file
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

// ----------------------------------------------------------------
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

func plus_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval + float64(val2.intval))
}
func plus_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) + val2.floatval)
}
func plus_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval + val2.floatval)
}

//  // var pfunc func(*Mlrval, *Mlrval) Mlrval
type dyadicFunc func(*Mlrval, *Mlrval) Mlrval

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

// ----------------------------------------------------------------
//static mv_binary_func_t* plus_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT        FLOAT      BOOL
//	/*ERROR*/  {_err, _err,  _err, _err,  _err,      _err,      _err},
//	/*ABSENT*/ {_err, _a,    _a,   _err,  _2___,        _2___,        _err},
//	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,      _emt,      _err},
//	/*STRING*/ {_err, _err,  _err, _err,  _err,      _err,      _err},
//	/*INT*/    {_err, _1___,    _emt, _err,  plus_n_ii, plus_f_if, _err},
//	/*FLOAT*/  {_err, _1___,    _emt, _err,  plus_f_fi, plus_f_ff, _err},
//	/*BOOL*/   {_err, _err,  _err, _err,  _err,      _err,      _err},
//};
// mv_t x_xx_plus_func(mv_t* pval1, mv_t* pval2) { return (plus_dispositions[pval1->type][pval2->type])(pval1,pval2); }
// func MvPlus(val1, val2 *Mlrval) Mlrval {
//	return plusDispositions[val1.mvtype][val2.mvtype]
// }
