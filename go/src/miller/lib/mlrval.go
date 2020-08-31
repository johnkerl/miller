package lib

import (
	"fmt"
	"os"
	"strconv"
)

// ================================================================
// Requirements for mlrvals:
//
// * Keep original string-formatting even if parseable/parsed as int
//   o E.g. if 005 (octal), pass through as 005 unless math is done on it
//   o Likewise with number of decimal places -- 7.4 not 7.400 or (worse) 7.399999999
// * Invalidate the string-formatting as the output of a computational result
// * Have number-to-string formatting methods in the API/DSL which stick the string format
// * Final to-string method
//
// Also:
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
	} else if input == "false"{
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
