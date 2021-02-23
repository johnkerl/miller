// ================================================================
// Constructors
// ================================================================

package types

import (
	"miller/src/lib"
)

//// ----------------------------------------------------------------
//func (this *Mlrval) CopyFrom(that *Mlrval) {
//	*this = *that
//	if this.mvtype == MT_MAP {
//		this.mapval = that.mapval.Copy()
//	} else if this.mvtype == MT_ARRAY {
//		this.arrayval = CopyMlrvalArray(that.arrayval)
//	}
//}

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
