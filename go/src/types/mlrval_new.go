// ================================================================
// Constructors
// ================================================================

package types

import (
	"errors"
	"strings"

	"mlr/src/lib"
)

// ----------------------------------------------------------------
// TODO: comment how these are part of the copy-reduction project.

var value_MLRVAL_ERROR = mlrvalFromError()
var value_MLRVAL_ABSENT = mlrvalFromAbsent()
var value_MLRVAL_NULL = mlrvalFromNull()
var value_MLRVAL_VOID = mlrvalFromVoid()
var value_MLRVAL_TRUE = mlrvalFromTrue()
var value_MLRVAL_FALSE = mlrvalFromFalse()

var MLRVAL_ERROR = &value_MLRVAL_ERROR
var MLRVAL_ABSENT = &value_MLRVAL_ABSENT
var MLRVAL_NULL = &value_MLRVAL_NULL
var MLRVAL_VOID = &value_MLRVAL_VOID
var MLRVAL_TRUE = &value_MLRVAL_TRUE
var MLRVAL_FALSE = &value_MLRVAL_FALSE

// ----------------------------------------------------------------
func MlrvalPointerFromString(input string) *Mlrval {
	if input == "" {
		return MLRVAL_VOID
	}
	var mv Mlrval
	mv.mvtype = MT_STRING
	mv.printrep = input
	mv.printrepValid = true
	return &mv
}

// xxx comment why two -- one for from parsed user data; other for from math ops
func MlrvalPointerFromIntString(input string) *Mlrval {
	ival, ok := lib.TryIntFromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	lib.InternalCodingErrorIf(!ok)
	var mv Mlrval
	mv.mvtype = MT_INT
	mv.printrep = input
	mv.printrepValid = true
	mv.intval = ival
	return &mv
}

func MlrvalPointerFromInt(input int) *Mlrval {
	var mv Mlrval
	mv.mvtype = MT_INT
	mv.printrepValid = false
	mv.intval = input
	return &mv
}

// xxx comment why two -- one for from parsed user data; other for from math ops
// xxx comment assummption is input-string already deemed parseable so no error return
func MlrvalPointerFromFloat64String(input string) *Mlrval {
	fval, ok := lib.TryFloat64FromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	lib.InternalCodingErrorIf(!ok)
	var mv Mlrval
	mv.mvtype = MT_FLOAT
	mv.printrep = input
	mv.printrepValid = true
	mv.floatval = fval
	return &mv
}

func MlrvalPointerFromFloat64(input float64) *Mlrval {
	var mv Mlrval
	mv.mvtype = MT_FLOAT
	mv.printrepValid = false
	mv.floatval = input
	return &mv
}

func MlrvalPointerFromBool(input bool) *Mlrval {
	if input == true {
		return MLRVAL_TRUE
	} else {
		return MLRVAL_FALSE
	}
}

func MlrvalPointerFromBoolString(input string) *Mlrval {
	if input == "true" {
		return MLRVAL_TRUE
	} else if input == "false" {
		return MLRVAL_FALSE
	} else {
		lib.InternalCodingErrorIf(true)
		return MLRVAL_ERROR // not reached
	}
}

var downcasedFloatNamesToNotInfer = map[string]bool{
	"inf":       true,
	"+inf":      true,
	"-inf":      true,
	"infinity":  true,
	"+infinity": true,
	"-infinity": true,
	"nan":       true,
}

// MlrvalPointerFromInferredTypeForDataFiles is for parsing field values from
// data files (except JSON, which is typed -- "true" and true are distinct).
// Mostly the same as MlrvalPointerFromInferredType, except it doesn't
// auto-infer true/false to bool; don't auto-infer NaN/Inf to float; etc.
func MlrvalPointerFromInferredTypeForDataFiles(input string) *Mlrval {
	if input == "" {
		return MLRVAL_VOID
	}

	_, iok := lib.TryIntFromString(input)
	if iok {
		return MlrvalPointerFromIntString(input)
	}

	if downcasedFloatNamesToNotInfer[strings.ToLower(input)] == false {
		_, fok := lib.TryFloat64FromString(input)
		if fok {
			return MlrvalPointerFromFloat64String(input)
		}
	}

	return MlrvalPointerFromString(input)
}

func MlrvalPointerFromInferredType(input string) *Mlrval {
	if input == "" {
		return MLRVAL_VOID
	}

	_, iok := lib.TryIntFromString(input)
	if iok {
		return MlrvalPointerFromIntString(input)
	}

	_, fok := lib.TryFloat64FromString(input)
	if fok {
		return MlrvalPointerFromFloat64String(input)
	}

	_, bok := lib.TryBoolFromBoolString(input)
	if bok {
		return MlrvalPointerFromBoolString(input)
	}

	return MlrvalPointerFromString(input)
}

func MlrvalPointerFromEmptyMap() *Mlrval {
	var mv Mlrval
	mv.mvtype = MT_MAP
	mv.printrepValid = false
	mv.mapval = NewMlrmap()
	return &mv
}

// ----------------------------------------------------------------
// Used by MlrvalFormatter (fmtnum DSL function, format-values verb, etc).
// Each mlrval has printrep and a printrepValid for its original string, then a
// type-code like MT_INT or MT_FLOAT, and type-specific storage like intval or
// floatval.
//
// If the user has taken a mlrval with original string "3.14" and formatted it
// with "%.4f" then its printrep will be "3.1400" but its type should still be
// MT_FLOAT.
//
// If on the other hand the user has formatted the same mlrval with "[[%.4f]]"
// then its printrep will be "[[3.1400]]" and it will be MT_STRING.
// This method supports that.

func MlrvalTryPointerFromFloatString(input string) *Mlrval {
	_, fok := lib.TryFloat64FromString(input)
	if fok {
		return MlrvalPointerFromFloat64String(input)
	} else {
		return MlrvalPointerFromString(input)
	}
}

func MlrvalTryPointerFromIntString(input string) *Mlrval {
	_, iok := lib.TryIntFromString(input)
	if iok {
		return MlrvalPointerFromIntString(input)
	} else {
		return MlrvalPointerFromString(input)
	}
}

// ----------------------------------------------------------------
// Does not copy the data. We can make a SetFromArrayLiteralCopy if needed
// using values.CopyMlrvalArray().
func MlrvalPointerFromArrayReference(input []Mlrval) *Mlrval {
	var mv Mlrval
	mv.mvtype = MT_ARRAY
	mv.printrepValid = false
	mv.arrayval = input
	return &mv
}

func MlrvalPointerFromMap(mlrmap *Mlrmap) *Mlrval {
	mv := MlrvalPointerFromEmptyMap()
	if mlrmap == nil {
		// TODO maybe return 2nd-arg error in the API
		return MLRVAL_ERROR
	}

	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		mv.mapval.PutCopy(pe.Key, pe.Value)
	}
	return mv
}

// Like previous but doesn't copy. Only safe when the argument's sole purpose
// is to be passed into here.
func MlrvalPointerFromMapReferenced(mlrmap *Mlrmap) *Mlrval {
	mv := MlrvalPointerFromEmptyMap()
	if mlrmap == nil {
		// xxx maybe return 2nd-arg error in the API
		return MLRVAL_ERROR
	}

	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		mv.mapval.PutReference(pe.Key, pe.Value)
	}
	return mv
}

// TODO: comment not MLRVAL_PENDING constants since this intended to be mutated
// by the JSON parser.
func MlrvalPointerFromPending() *Mlrval {
	var mv Mlrval
	mv.mvtype = MT_PENDING
	mv.printrepValid = false
	return &mv
}

// ----------------------------------------------------------------
// TODO: comment about being designed to be mutated for JSON API.
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

func mlrvalFromError() Mlrval {
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

func mlrvalFromAbsent() Mlrval {
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

func mlrvalFromNull() Mlrval {
	return Mlrval{
		mvtype:        MT_NULL,
		printrep:      "null",
		printrepValid: true,
		intval:        0,
		floatval:      0.0,
		boolval:       false,
		arrayval:      nil,
		mapval:        nil,
	}
}

func mlrvalFromVoid() Mlrval {
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
	if input == "" {
		return mlrvalFromVoid()
	} else {
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
}

func MlrvalFromInt(input int) Mlrval {
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

func mlrvalFromTrue() Mlrval {
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

func mlrvalFromFalse() Mlrval {
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
		return mlrvalFromTrue()
	} else {
		return mlrvalFromFalse()
	}
}

func MlrvalFromBoolString(input string) Mlrval {
	if input == "true" {
		return mlrvalFromTrue()
	} else {
		return mlrvalFromFalse()
	}
	// else panic
}

// ----------------------------------------------------------------
// Does not copy the data. We can make a MlrvalFromArrayLiteralCopy if needed,
// using values.CopyMlrvalArray().

func MlrvalEmptyArray() Mlrval {
	return Mlrval{
		mvtype:        MT_ARRAY,
		printrep:      "(bug-if-you-see-this-array-type)",
		printrepValid: false,
		intval:        0,
		floatval:      0.0,
		boolval:       false,
		arrayval:      make([]Mlrval, 0, 10),
		mapval:        nil,
	}
}

// Users can do things like '$new[1][2][3] = 4' even if '$new' isn't already
// allocated. This function supports that.
func NewSizedMlrvalArray(length int) *Mlrval {
	arrayval := make([]Mlrval, length, 2*length)

	for i := 0; i < int(length); i++ {
		arrayval[i] = MlrvalFromString("")
	}

	return &Mlrval{
		mvtype:        MT_ARRAY,
		printrep:      "(bug-if-you-see-this-array-type)",
		printrepValid: false,
		intval:        0,
		floatval:      0.0,
		boolval:       false,
		arrayval:      arrayval,
		mapval:        nil,
	}
}

func LengthenMlrvalArray(array *[]Mlrval, newLength64 int) {
	newLength := int(newLength64)
	lib.InternalCodingErrorIf(newLength <= len(*array))

	if newLength <= cap(*array) {
		newArray := (*array)[:newLength]
		for zindex := len(*array); zindex < newLength; zindex++ {
			// TODO: comment why not MT_ABSENT or MT_VOID
			newArray[zindex] = *MLRVAL_NULL
		}
		*array = newArray
	} else {
		newArray := make([]Mlrval, newLength, 2*newLength)
		zindex := 0
		for zindex = 0; zindex < len(*array); zindex++ {
			newArray[zindex] = (*array)[zindex]
		}
		for zindex = len(*array); zindex < newLength; zindex++ {
			// TODO: comment why not MT_ABSENT or MT_VOID
			newArray[zindex] = *MLRVAL_NULL
		}
		*array = newArray
	}
}

// ----------------------------------------------------------------
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

// Like previous but doesn't copy. Only safe when the argument's sole purpose
// is to be passed into here.
func MlrvalFromMapReferenced(mlrmap *Mlrmap) Mlrval {
	mv := MlrvalEmptyMap()
	if mlrmap == nil {
		// xxx maybe return 2nd-arg error in the API
		return *MLRVAL_ERROR
	}

	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		mv.mapval.PutReference(pe.Key, pe.Value)
	}

	return mv
}

// ----------------------------------------------------------------
// This is for auto-deepen of nested maps in things like
//
//   $foo[1]["a"][2]["b"] = 3
//
// Autocreated levels are maps.  Array levels can be explicitly created e.g.
//
//   $foo[1]["a"] ??= []
//   $foo[1]["a"][2]["b"] = 3

func NewMlrvalForAutoDeepen(mvtype MVType) (*Mlrval, error) {
	if mvtype == MT_STRING || mvtype == MT_INT {
		empty := MlrvalEmptyMap()
		return &empty, nil
	} else {
		return nil, errors.New(
			"Miller: indices must be string or int; got " + GetTypeName(mvtype),
		)
	}
}
