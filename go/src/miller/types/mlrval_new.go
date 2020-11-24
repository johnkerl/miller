// ================================================================
// Constructors
// ================================================================

package types

import (
	"errors"

	"miller/lib"
)

// ----------------------------------------------------------------
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
	if input == "" {
		return MlrvalFromVoid()
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

// xxx comment why two -- one for from parsed user data; other for from math ops
func MlrvalFromInt64String(input string) Mlrval {
	ival, ok := lib.TryInt64FromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	lib.InternalCodingErrorIf(!ok)
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

// xxx comment why two -- one for from parsed user data; other for from math ops
// xxx comment assummption is input-string already deemed parseable so no error return
func MlrvalFromFloat64String(input string) Mlrval {
	fval, ok := lib.TryFloat64FromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	lib.InternalCodingErrorIf(!ok)
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

func MlrvalFromInferredType(input string) Mlrval {
	// xxx the parsing has happened so stash it ...
	// xxx emphasize the invariant that a non-invalid printrep always
	// matches the nval ...
	if input == "" {
		return MlrvalFromVoid()
	}

	_, iok := lib.TryInt64FromString(input)
	if iok {
		return MlrvalFromInt64String(input)
	}

	_, fok := lib.TryFloat64FromString(input)
	if fok {
		return MlrvalFromFloat64String(input)
	}

	_, bok := lib.TryBoolFromBoolString(input)
	if bok {
		return MlrvalFromBoolString(input)
	}

	return MlrvalFromString(input)
}

// ----------------------------------------------------------------
// Does not copy the data. We can make a MlrvalFromArrayLiteralCopy if needed,
// using values.CopyMlrvalArray().
func MlrvalFromArrayLiteralReference(input []Mlrval) Mlrval {
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
		arrayval:      make([]Mlrval, 0, 10),
		mapval:        nil,
	}
}

// Users can do things like '$new[1][2][3] = 4' even if '$new' isn't already
// allocated. This function supports that.
func NewSizedMlrvalArray(length int64) *Mlrval {
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

func LengthenMlrvalArray(array *[]Mlrval, newLength64 int64) {
	newLength := int(newLength64)
	lib.InternalCodingErrorIf(newLength <= len(*array))

	if newLength <= cap(*array) {
		newArray := (*array)[:newLength]
		for zindex := len(*array); zindex < newLength; zindex++ {
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

func MlrvalFromMap(that *Mlrmap) Mlrval {
	this := MlrvalEmptyMap()
	if that == nil {
		// xxx maybe return 2nd-arg error in the API
		return MlrvalFromError()
	}

	for pe := that.Head; pe != nil; pe = pe.Next {
		this.mapval.PutCopy(pe.Key, pe.Value)
	}

	return this
}

// ----------------------------------------------------------------
// This is for auto-deepen of nested arrays/maps in things like
// '$foo[1]["a"][2]["b"] = 3' It takes the type of the next index-slot to be
// created, returing string for map, int for array, error otherwise.

func NewMlrvalForAutoDeepen(mvtype MVType) (*Mlrval, error) {
	if mvtype == MT_STRING {
		empty := MlrvalEmptyMap()
		return &empty, nil
	} else if mvtype == MT_INT {
		empty := MlrvalEmptyArray()
		return &empty, nil
	} else {
		return nil, errors.New(
			"Miller: indices must be string or int; got " + GetTypeName(mvtype),
		)
	}
}
