package lib

import (
	"strconv"
)

// Constructors

func MlrvalFromPending() Mlrval {
	return Mlrval{
		mvtype:        MT_PENDING,
		printrep:      "(bug-if-you-see-this-pending-type)",
		printrepValid: true,
		intval:        0,
		floatval:      0.0,
		boolval:       false,
		arrayval:      nil,
		mapval:        nil,
	}
}

func MlrvalFromError() Mlrval {
	return Mlrval{
		mvtype:        MT_ERROR,
		printrep:      "(error)", // xxx const somewhere
		printrepValid: true,
		intval:        0,
		floatval:      0.0,
		boolval:       false,
		arrayval:      nil,
		mapval:        nil,
	}
}

func MlrvalFromAbsent() Mlrval {
	return Mlrval{
		mvtype:        MT_ABSENT,
		printrep:      "(absent)",
		printrepValid: true,
		intval:        0,
		floatval:      0.0,
		boolval:       false,
		arrayval:      nil,
		mapval:        nil,
	}
}

func MlrvalFromVoid() Mlrval {
	return Mlrval{
		mvtype:        MT_VOID,
		printrep:      "",
		printrepValid: true,
		intval:        0,
		floatval:      0.0,
		boolval:       false,
		arrayval:      nil,
		mapval:        nil,
	}
}

func MlrvalFromString(input string) Mlrval {
	return Mlrval{
		mvtype:        MT_STRING,
		printrep:      input,
		printrepValid: true,
		intval:        0,
		floatval:      0.0,
		boolval:       false,
		arrayval:      nil,
		mapval:        nil,
	}
}

// xxx comment why two -- one for from parsed user data; other for from math ops
func MlrvalFromInt64String(input string) Mlrval {
	ival, ok := tryInt64FromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	InternalCodingErrorIf(!ok)
	return Mlrval{
		mvtype:        MT_INT,
		printrep:      input,
		printrepValid: true,
		intval:        ival,
		floatval:      0.0,
		boolval:       false,
		arrayval:      nil,
		mapval:        nil,
	}
}

func MlrvalFromInt64(input int64) Mlrval {
	return Mlrval{
		mvtype:        MT_INT,
		printrep:      "(bug-if-you-see-this-int-type)",
		printrepValid: false,
		intval:        input,
		floatval:      0.0,
		boolval:       false,
		arrayval:      nil,
		mapval:        nil,
	}
}

// Tries decimal, hex, octal, and binary.
func tryInt64FromString(input string) (int64, bool) {
	ival, err := strconv.ParseInt(input, 0 /* check all*/, 64)
	if err == nil {
		return ival, true
	} else {
		return 0, false
	}
}

// xxx comment why two -- one for from parsed user data; other for from math ops
// xxx comment assummption is input-string already deemed parseable so no error return
func MlrvalFromFloat64String(input string) Mlrval {
	fval, ok := tryFloat64FromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	InternalCodingErrorIf(!ok)
	return Mlrval{
		mvtype:        MT_FLOAT,
		printrep:      input,
		printrepValid: true,
		intval:        0,
		floatval:      fval,
		boolval:       false,
		arrayval:      nil,
		mapval:        nil,
	}
}

func MlrvalFromFloat64(input float64) Mlrval {
	return Mlrval{
		mvtype:        MT_FLOAT,
		printrep:      "(bug-if-you-see-this-float-type)",
		printrepValid: false,
		intval:        0,
		floatval:      input,
		boolval:       false,
		arrayval:      nil,
		mapval:        nil,
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

func MlrvalFromTrue() Mlrval {
	return Mlrval{
		mvtype:        MT_BOOL,
		printrep:      "true",
		printrepValid: true,
		intval:        0,
		floatval:      0.0,
		boolval:       true,
		arrayval:      nil,
		mapval:        nil,
	}
}

func MlrvalFromFalse() Mlrval {
	return Mlrval{
		mvtype:        MT_BOOL,
		printrep:      "false",
		printrepValid: true,
		intval:        0,
		floatval:      0.0,
		boolval:       false,
		arrayval:      nil,
		mapval:        nil,
	}
}

func MlrvalFromBool(input bool) Mlrval {
	if input == true {
		return MlrvalFromTrue()
	} else {
		return MlrvalFromFalse()
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

// xxx copy or no? needs a Mlrval.Copy() (deep) if so.
func MlrvalFromArrayLiteral(input []Mlrval) Mlrval {
	return Mlrval{
		mvtype:        MT_ARRAY,
		printrep:      "(bug-if-you-see-this-array-type)",
		printrepValid: false,
		intval:        0,
		floatval:      0.0,
		boolval:       false,
		arrayval:      input,
		mapval:        nil,
	}
}

func MlrvalEmptyArray() Mlrval {
	return Mlrval{
		mvtype:        MT_ARRAY,
		printrep:      "(bug-if-you-see-this-array-type)",
		printrepValid: false,
		intval:        0,
		floatval:      0.0,
		boolval:       false,
		arrayval:      make([]Mlrval, 0),
		mapval:        nil,
	}
}

func MlrvalEmptyMap() Mlrval {
	return Mlrval{
		mvtype:        MT_MAP,
		printrep:      "(bug-if-you-see-this-map-type)",
		printrepValid: false,
		intval:        0,
		floatval:      0.0,
		boolval:       false,
		arrayval:      nil,
		mapval:        NewMlrmap(),
	}
}
