// ================================================================
// Constructors
// ================================================================

package types

import (
	"miller/src/lib"
)

// ----------------------------------------------------------------
func (this *Mlrval) CopyFrom(that *Mlrval) {
	*this = *that
	if this.mvtype == MT_MAP {
		this.mapval = that.mapval.Copy()
	} else if this.mvtype == MT_ARRAY {
		this.arrayval = CopyMlrvalArray(that.arrayval)
	}
}

// ----------------------------------------------------------------
func (this *Mlrval) SetFromPending() {
	this.mvtype = MT_PENDING
	this.printrep = "(bug-if-you-see-this-pending-type)"
	this.printrepValid = false
}

func (this *Mlrval) SetFromError() {
	this.mvtype = MT_ERROR
	this.printrep = "(error)" // xxx const somewhere
	this.printrepValid = true
}

func (this *Mlrval) SetFromAbsent() {
	this.mvtype = MT_ABSENT
	this.printrep = "(absent)"
	this.printrepValid = true
}

func (this *Mlrval) SetFromVoid() {
	this.mvtype = MT_VOID
	this.printrep = ""
	this.printrepValid = true
}

func (this *Mlrval) SetFromString(input string) {
	if input == "" {
		this.SetFromVoid()
	} else {
		this.mvtype = MT_STRING
		this.printrep = input
		this.printrepValid = true
	}
}

// xxx comment why two -- one for from parsed user data; other for from math ops
func (this *Mlrval) SetFromIntString(input string) {
	ival, ok := lib.TryIntFromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	lib.InternalCodingErrorIf(!ok)
	this.mvtype = MT_INT
	this.printrep = input
	this.printrepValid = true
	this.intval = ival
}

func (this *Mlrval) SetFromInt(input int) {
	this.mvtype = MT_INT
	this.printrepValid = false
	this.intval = input
}

// xxx comment why two -- one for from parsed user data; other for from math ops
// xxx comment assummption is input-string already deemed parseable so no error return
func (this *Mlrval) SetFromFloat64String(input string) {
	fval, ok := lib.TryFloat64FromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	lib.InternalCodingErrorIf(!ok)
	this.mvtype = MT_FLOAT
	this.printrep = input
	this.printrepValid = true
	this.floatval = fval
}

func (this *Mlrval) SetFromFloat64(input float64) {
	this.mvtype = MT_FLOAT
	this.printrepValid = false
	this.floatval = input
}

func (this *Mlrval) SetFromTrue() {
	this.mvtype = MT_BOOL
	this.printrep = "true"
	this.printrepValid = true
	this.boolval = true
}

func (this *Mlrval) SetFromFalse() {
	this.mvtype = MT_BOOL
	this.printrep = "false"
	this.printrepValid = true
	this.boolval = false
}

func (this *Mlrval) SetFromBool(input bool) {
	if input == true {
		this.SetFromTrue()
	} else {
		this.SetFromFalse()
	}
}

func (this *Mlrval) SetFromBoolString(input string) {
	if input == "true" {
		this.SetFromTrue()
	} else {
		this.SetFromFalse()
	}
	// else panic
}

func (this *Mlrval) SetFromInferredType(input string) {
	// xxx the parsing has happened so stash it ...
	// xxx emphasize the invariant that a non-invalid printrep always
	// matches the nval ...
	if input == "" {
		this.SetFromVoid()
	}

	_, iok := lib.TryIntFromString(input)
	if iok {
		this.SetFromIntString(input)
	}

	_, fok := lib.TryFloat64FromString(input)
	if fok {
		this.SetFromFloat64String(input)
	}

	_, bok := lib.TryBoolFromBoolString(input)
	if bok {
		this.SetFromBoolString(input)
	}

	this.SetFromString(input)
}

func (this *Mlrval) SetFromEmptyMap() {
	this.mvtype = MT_MAP
	this.printrepValid = false
	this.mapval = NewMlrmap()
}

// ----------------------------------------------------------------
// Does not copy the data. We can make a (this *Mlrval) SetFromArrayLiteralCopy if needed
// using values.CopyMlrvalArray().
func (this *Mlrval) SetFromArrayLiteralReference(input []Mlrval) {
	this.mvtype = MT_ARRAY
	this.printrepValid = false
	this.arrayval = input
}

func (this *Mlrval) SetFromMap(that *Mlrmap) {
	this.SetFromEmptyMap()
	if that == nil {
		// xxx maybe return 2nd-arg error in the API
		this.SetFromError()
		return
	}

	for pe := that.Head; pe != nil; pe = pe.Next {
		this.mapval.PutCopy(pe.Key, pe.Value)
	}
}

// Like previous but doesn't copy. Only safe when the argument's sole purpose
// is to be passed into here.
func (this *Mlrval) SetFromMapReferenced(that *Mlrmap) {
	this.SetFromEmptyMap()
	if that == nil {
		// xxx maybe return 2nd-arg error in the API
		this.SetFromError()
		return
	}

	for pe := that.Head; pe != nil; pe = pe.Next {
		this.mapval.PutReference(pe.Key, pe.Value)
	}
}
