// ================================================================
// FLATTEN/UNFLATTEN
//
// These are used by the flatten/unflatten verbs and DSL functions.  They are
// crucial to the operation of Miller 6 wherein records have full Mlrval
// values, i.e. they can be arrays/maps as well as int/float/string.
//
// When we read JSON and write (say) CSV, we have two choices for handling the
// fact that JSON handles multi-level data and CSV does not:
//
// (1) JSON-stringify values, using the json-stringify verb or json_stringify
// DSL function. For example, the array of ints [1,2,3] becomes the string
// "[1,2,3]" which works fine as a CSV field.
//
// (2) Flatten them by key-spreading. For example, the single field with key
// "x" with value {"a":1,"b":2} flattens to the *pair* of fields x:a=1 and
// x:b=2.
//
// The former are used implicitly (i.e. unless the user explicitly requests
// otherwise) when we convert to/from JSON.
// ================================================================

package types

import (
	"strings"

	"miller/lib"
)

// ----------------------------------------------------------------
// Flattens all field values in the record. This is a special case of
// FlattenFields but it's worth its own special case (to avoid iffing on the
// nullity of the fieldNameSet) since the flatten/unflatten check is done by
// default on ALL Miller records whenever we convert to/from JSON. So, the
// default path should be fast.

func (this *Mlrmap) Flatten(separator string) {
	if !this.isFlattenable() { // fast path: don't modify the record at all
		return
	}

	that := NewMlrmapAsRecord()

	for pe := this.Head; pe != nil; pe = pe.Next {
		if pe.Value.IsArrayOrMap() {
			pieces := pe.Value.FlattenToMap(pe.Key, separator)
			for pf := pieces.GetMap().Head; pf != nil; pf = pf.Next {
				that.PutReference(pf.Key, pf.Value)
			}
		} else {
			that.PutReference(pe.Key, pe.Value)
		}
	}

	*this = *that
}

// ----------------------------------------------------------------
// For mlr flatten -f.

func (this *Mlrmap) FlattenFields(
	fieldNameSet map[string]bool,
	separator string,
) {
	if !this.isFlattenable() { // fast path
		return
	}

	that := NewMlrmapAsRecord()

	for pe := this.Head; pe != nil; pe = pe.Next {
		if pe.Value.IsArrayOrMap() && fieldNameSet[pe.Key] {
			pieces := pe.Value.FlattenToMap(pe.Key, separator)
			for pf := pieces.GetMap().Head; pf != nil; pf = pf.Next {
				that.PutReference(pf.Key, pf.Value)
			}
		} else {
			that.PutReference(pe.Key, pe.Value)
		}
	}

	*this = *that
}

// ----------------------------------------------------------------
// Optimization for Flatten, to avoid needless data motion in the case
// where all field values are non-collections.

func (this *Mlrmap) isFlattenable() bool {
	for pe := this.Head; pe != nil; pe = pe.Next {
		if pe.Value.IsArrayOrMap() {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------
func (this *Mlrmap) Unflatten(separator string) {
	that := NewMlrmapAsRecord()

	for pe := this.Head; pe != nil; pe = pe.Next {
		if strings.Contains(pe.Key, separator) {
			arrayOfIndices := mlrvalSplitAXHelper(pe.Key, separator)
			that.PutIndexed(
				MakePointerArray(arrayOfIndices.arrayval),
				unflattenTerminal(pe.Value).Copy(),
			)
		} else {
			that.PutReference(pe.Key, unflattenTerminal(pe.Value))
		}
	}

	*this = *that
}

// ----------------------------------------------------------------
// For mlr unflatten -f.
func (this *Mlrmap) UnflattenFields(
	fieldNameSet map[string]bool,
	separator string,
) {
	that := NewMlrmapAsRecord()

	for pe := this.Head; pe != nil; pe = pe.Next {
		if strings.Contains(pe.Key, separator) {
			arrayOfIndices := mlrvalSplitAXHelper(pe.Key, separator)
			lib.InternalCodingErrorIf(len(arrayOfIndices.arrayval) < 1)
			baseIndex := arrayOfIndices.arrayval[0].String()
			if fieldNameSet[baseIndex] {
				that.PutIndexed(
					MakePointerArray(arrayOfIndices.arrayval),
					unflattenTerminal(pe.Value).Copy(),
				)
			} else {
				that.PutReference(pe.Key, unflattenTerminal(pe.Value))
			}
		} else {
			that.PutReference(pe.Key, unflattenTerminal(pe.Value))
		}
	}

	*this = *that
}

// ----------------------------------------------------------------
// Flatten of empty map and empty array produce "{}" and "[]" as special cases.
// (Without this, key-spreading would cause such fields to disappear entirely:
// the field "x" -> {"a": 1, "b": 2} would spread to the pair of fields "x:a"
// -> 1 and "x:b" -> 2, and the field "x" -> {"a": 1} would spread to the
// single field "x:a" -> 1, so the field "x" -> {} would spread to zero
// fields.) Here we reverse that special case of the flatten operation.

func unflattenTerminal(input *Mlrval) *Mlrval {
	if !input.IsString() {
		return input
	}
	if input.printrep == "{}" {
		retval := MlrvalFromMapReferenced(NewMlrmap())
		return &retval
	}
	if input.printrep == "[]" {
		retval := MlrvalFromArrayLiteralReference(make([]Mlrval, 0))
		return &retval
	}
	return input
}
