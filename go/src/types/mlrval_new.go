// ================================================================
// Constructors
// ================================================================

package types

import (
	"errors"

	"miller/src/lib"
)

// ----------------------------------------------------------------
// TODO: comment how these are part of the copy-reduction project.

var value_MLRVAL_ERROR = mlrvalFromError()
var value_MLRVAL_ABSENT = mlrvalFromAbsent()
var value_MLRVAL_VOID = mlrvalFromVoid()
var value_MLRVAL_TRUE = mlrvalFromTrue()
var value_MLRVAL_FALSE = mlrvalFromFalse()

var MLRVAL_ERROR = &value_MLRVAL_ERROR
var MLRVAL_ABSENT = &value_MLRVAL_ABSENT
var MLRVAL_VOID = &value_MLRVAL_VOID
var MLRVAL_TRUE = &value_MLRVAL_TRUE
var MLRVAL_FALSE = &value_MLRVAL_FALSE

// ----------------------------------------------------------------
func MlrvalPointerFromString(input string) *Mlrval {
	if input == "" {
		return MLRVAL_VOID
	}
	var this Mlrval
	this.mvtype = MT_STRING
	this.printrep = input
	this.printrepValid = true
	return &this
}

// xxx comment why two -- one for from parsed user data; other for from math ops
func MlrvalPointerFromIntString(input string) *Mlrval {
	ival, ok := lib.TryIntFromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	lib.InternalCodingErrorIf(!ok)
	var this Mlrval
	this.mvtype = MT_INT
	this.printrep = input
	this.printrepValid = true
	this.intval = ival
	return &this
}

func MlrvalPointerFromInt(input int) *Mlrval {
	var this Mlrval
	this.mvtype = MT_INT
	this.printrepValid = false
	this.intval = input
	return &this
}

// xxx comment why two -- one for from parsed user data; other for from math ops
// xxx comment assummption is input-string already deemed parseable so no error return
func MlrvalPointerFromFloat64String(input string) *Mlrval {
	fval, ok := lib.TryFloat64FromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	lib.InternalCodingErrorIf(!ok)
	var this Mlrval
	this.mvtype = MT_FLOAT
	this.printrep = input
	this.printrepValid = true
	this.floatval = fval
	return &this
}

func MlrvalPointerFromFloat64(input float64) *Mlrval {
	var this Mlrval
	this.mvtype = MT_FLOAT
	this.printrepValid = false
	this.floatval = input
	return &this
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

func MlrvalPointerFromInferredType(input string) *Mlrval {
	// xxx the parsing has happened so stash it ...
	// xxx emphasize the invariant that a non-invalid printrep always
	// matches the nval ...
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
	var this Mlrval
	this.mvtype = MT_MAP
	this.printrepValid = false
	this.mapval = NewMlrmap()
	return &this
}

// ----------------------------------------------------------------
// Does not copy the data. We can make a SetFromArrayLiteralCopy if needed
// using values.CopyMlrvalArray().
func MlrvalPointerFromArrayLiteralReference(input []Mlrval) *Mlrval {
	var this Mlrval
	this.mvtype = MT_ARRAY
	this.printrepValid = false
	this.arrayval = input
	return &this
}

func MlrvalPointerFromMap(that *Mlrmap) *Mlrval {
	this := MlrvalPointerFromEmptyMap()
	if that == nil {
		// TODO maybe return 2nd-arg error in the API
		return MLRVAL_ERROR
	}

	for pe := that.Head; pe != nil; pe = pe.Next {
		this.mapval.PutCopy(pe.Key, pe.Value)
	}
	return this
}

// Like previous but doesn't copy. Only safe when the argument's sole purpose
// is to be passed into here.
func MlrvalPointerFromMapReferenced(that *Mlrmap) *Mlrval {
	this := MlrvalPointerFromEmptyMap()
	if that == nil {
		// xxx maybe return 2nd-arg error in the API
		return MLRVAL_ERROR
	}

	for pe := that.Head; pe != nil; pe = pe.Next {
		this.mapval.PutReference(pe.Key, pe.Value)
	}
	return this
}

// TODO: comment not MLRVAL_PENDING constants since this intended to be mutated
// by the JSON parser.
func MlrvalPointerFromPending() *Mlrval {
	var this Mlrval
	this.mvtype = MT_PENDING
	this.printrepValid = false
	return &this
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
			// TODO: comment why not MT_ABSENT
			newArray[zindex] = MlrvalFromString("")
		}
		*array = newArray
	} else {
		newArray := make([]Mlrval, newLength, 2*newLength)
		zindex := 0
		for zindex = 0; zindex < len(*array); zindex++ {
			newArray[zindex] = (*array)[zindex]
		}
		for zindex = len(*array); zindex < newLength; zindex++ {
			newArray[zindex] = MlrvalFromString("")
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
func MlrvalFromMapReferenced(that *Mlrmap) Mlrval {
	this := MlrvalEmptyMap()
	if that == nil {
		// xxx maybe return 2nd-arg error in the API
		return *MLRVAL_ERROR
	}

	for pe := that.Head; pe != nil; pe = pe.Next {
		this.mapval.PutReference(pe.Key, pe.Value)
	}

	return this
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
