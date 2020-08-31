package lib

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
//	/*ABSENT*/ {_err, _a,    _a,   _err,  _2,        _2,        _err},
//	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,      _emt,      _err},
//	/*STRING*/ {_err, _err,  _err, _err,  _err,      _err,      _err},
//	/*INT*/    {_err, _1,    _emt, _err,  plus_n_ii, plus_f_if, _err},
//	/*FLOAT*/  {_err, _1,    _emt, _err,  plus_f_fi, plus_f_ff, _err},
//	/*BOOL*/   {_err, _err,  _err, _err,  _err,      _err,      _err},
//};

// ----------------------------------------------------------------
type MVType int

const (
	MT_ERROR MVType = iota
	MT_ABSENT
	MT_VOID
	MT_STRING
	MT_INT
	MT_FLOAT
	MT_BOOL
	// Not a type -- this is a dimension for disposition matrices
	MT_DIM
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
	return Mlrval{
		MT_INT,
		input,
		true,
		0, // xxx parse
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
func MlrvalFromFloat64String(input string) Mlrval {
	return Mlrval{
		MT_FLOAT,
		input,
		true,
		0,
		0.0, // xxx parse
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
		this.printrepValid = true
	}
}
func (this *Mlrval) String() string {
	this.setPrintRep()
	return this.printrep
}
