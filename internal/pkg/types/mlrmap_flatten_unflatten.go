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

	"mlr/internal/pkg/lib"
)

// ----------------------------------------------------------------
// Flattens all field values in the record. This is a special case of
// FlattenFields but it's worth its own special case (to avoid iffing on the
// nullity of the fieldNameSet) since the flatten/unflatten check is done by
// default on ALL Miller records whenever we convert to/from JSON. So, the
// default path should be fast.
//
// Examples:
// * The single field x = {"a": 7, "b": 8, "c": 9} becomes the three fields
//   x.a = 7, x.b = 8, x.c = 9.
// * The single field x = [7,8,9] becomes the three fields
//   x.1 = 7, x.2 = 8, x.3 = 9.

func (mlrmap *Mlrmap) Flatten(separator string) {
	if !mlrmap.isFlattenable() { // fast path: don't modify the record at all
		return
	}

	other := NewMlrmapAsRecord()

	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		if pe.Value.IsArrayOrMap() {
			pieces := pe.Value.FlattenToMap(pe.Key, separator)
			for pf := pieces.GetMap().Head; pf != nil; pf = pf.Next {
				other.PutReference(pf.Key, pf.Value)
			}
		} else {
			other.PutReference(pe.Key, pe.Value)
		}
	}

	*mlrmap = *other
}

// ----------------------------------------------------------------
// For mlr flatten -f.

func (mlrmap *Mlrmap) FlattenFields(
	fieldNameSet map[string]bool,
	separator string,
) {
	if !mlrmap.isFlattenable() { // fast path
		return
	}

	other := NewMlrmapAsRecord()

	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		if pe.Value.IsArrayOrMap() && fieldNameSet[pe.Key] {
			pieces := pe.Value.FlattenToMap(pe.Key, separator)
			for pf := pieces.GetMap().Head; pf != nil; pf = pf.Next {
				other.PutReference(pf.Key, pf.Value)
			}
		} else {
			other.PutReference(pe.Key, pe.Value)
		}
	}

	*mlrmap = *other
}

// ----------------------------------------------------------------
// Optimization for Flatten, to avoid needless data motion in the case
// where all field values are non-collections.

func (mlrmap *Mlrmap) isFlattenable() bool {
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		if pe.Value.IsArrayOrMap() {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------
// For mlr unflatten without -f. This undoes Unflatten.  This is for conversion
// from non-JSON to JSON.  If there are fields x.a, x.b, x.c, etc. they're put
// into a single field x with map-valued value keyed by "a", "b", "c".

// Special case: if the resulting string keys are string representations of 1,
// 2, 3, etc -- without gaps -- then the map is converted to an array.
//
// Examples:
//
// * The three fields x.a = 7, x.b = 8, x.c = 9  become
//   the single field x = {"a": 7, "b": 8, "c": 9}.
//
// * The three fields x.1 = 7, x.2 = 8, x.3 = 9 become
//   the single field x = [7,8,9].
//
// * The two fields x.1 = 7, x.3 = 9 become
//   the single field x = {"1": 7, "3": 9}
func (mlrmap *Mlrmap) Unflatten(
	separator string,
) {
	*mlrmap = *(mlrmap.CopyUnflattened(separator))
}

func (mlrmap *Mlrmap) CopyUnflattened(
	separator string,
) *Mlrmap {
	other := NewMlrmapAsRecord()
	affectedBaseIndices := make(map[string]bool)

	// We'll come through this loop once for x.a, another for x.b, etc.
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		// Is the field name something dot something?
		if strings.Contains(pe.Key, separator) {
			arrayOfIndices := mlrvalSplitAXHelper(pe.Key, separator)
			lib.InternalCodingErrorIf(len(arrayOfIndices.arrayval) < 1)
			// If the input field name was "x.a" then remember the "x".
			baseIndex := arrayOfIndices.arrayval[0].String()
			affectedBaseIndices[baseIndex] = true
			// Use PutIndexed to assign $x["a"] = 7, or $x["b"] = 8, etc.
			other.PutIndexed(
				MakePointerArray(arrayOfIndices.arrayval),
				unflattenTerminal(pe.Value).Copy(),
			)
		} else {
			other.PutReference(pe.Key, unflattenTerminal(pe.Value))
		}
	}

	// Go through all the field names which were turned into maps -- e.g.  "x"
	// in the example above -- and see if the keys were like "1", "2", etc and
	// if so then convert to array. This undoes how Flatten flattens arrays.
	for baseIndex := range affectedBaseIndices {
		oldValue := other.Get(baseIndex)
		lib.InternalCodingErrorIf(oldValue == nil)
		newValue := BIF_arrayify(oldValue)
		other.PutReference(baseIndex, newValue)
	}

	return other
}

// ----------------------------------------------------------------
// For mlr unflatten -f. See comments on Unflatten. Largely copypasta of
// Unflatten, but split out separately since Flatten needn't check a
// fieldNameSet.
func (mlrmap *Mlrmap) UnflattenFields(
	fieldNameSet map[string]bool,
	separator string,
) {
	*mlrmap = *(mlrmap.CopyUnflattenFields(fieldNameSet, separator))
}

func (mlrmap *Mlrmap) CopyUnflattenFields(
	fieldNameSet map[string]bool,
	separator string,
) *Mlrmap {
	other := NewMlrmapAsRecord()
	affectedBaseIndices := make(map[string]bool)

	// We'll come through this loop once for x.a, another for x.b, etc.
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		// Is the field name something dot something?
		if strings.Contains(pe.Key, separator) {
			arrayOfIndices := mlrvalSplitAXHelper(pe.Key, separator)
			lib.InternalCodingErrorIf(len(arrayOfIndices.arrayval) < 1)
			// If the input field name was "x.a" then remember the "x".
			baseIndex := arrayOfIndices.arrayval[0].String()
			if fieldNameSet[baseIndex] {
				// Use PutIndexed to assign $x["a"] = 7, or $x["b"] = 8, etc.
				other.PutIndexed(
					MakePointerArray(arrayOfIndices.arrayval),
					unflattenTerminal(pe.Value).Copy(),
				)
				affectedBaseIndices[baseIndex] = true
			} else {
				other.PutReference(pe.Key, unflattenTerminal(pe.Value))
			}
		} else {
			other.PutReference(pe.Key, unflattenTerminal(pe.Value))
		}
	}

	// Go through all the field names which were turned into maps -- e.g.  "x"
	// in the example above -- and see if the keys were like "1", "2", etc and
	// if so then convert to array. This undoes how Flatten flattens arrays.
	for baseIndex := range affectedBaseIndices {
		oldValue := other.Get(baseIndex)
		lib.InternalCodingErrorIf(oldValue == nil)
		newValue := BIF_arrayify(oldValue)
		other.PutReference(baseIndex, newValue)
	}

	return other
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
		return MlrvalFromMapReferenced(NewMlrmap())
	}
	if input.printrep == "[]" {
		return MlrvalFromArrayReference(make([]Mlrval, 0))
	}
	return input
}
