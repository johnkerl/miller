// ================================================================
// ABOUT ARRAY/MAP INDEXING
//
// Arrays:
//
// * Within the Miller domain-specific language:
//   o Arrays are indexed 1-up within the DSL
//   o Negative-index aliases are supported: -1 means the last element
//     of the array (if any), -2 means second-to-last, and in general
//     -n..-1 are aliases for 1..n.
// * Within the Go implementation:
//   o Array indices are of course 0-up.
//
// Maps:
//
// * Map keys are strings
//
// * They can also be indexed positionally (which is a linear search
//   through the linked-hash-map struct). One of Miller's primary reasons for
//   existing is to allow string-keyed access, but, sometimes people really
//   want to access the nth field in a record. These are also 1-up.
//
// Why 1-up indices?
//
// * This is an odd choice in the broader programming-language context.
//   There are a few languages such as Fortran, Julia, Matlab, and R which are
//   1-up but the overal trend is decidedly toward 0-up. This means that
//   if Miller does 1-up it should do so for a good reason.
//
// * Reasons: so many other things are already 1-up, mostly inherited
//   from AWK:
//
//   o NR and FNR, the record-counter context variables start with 1.
//     When we retain records, '@records[NR] = $*' is the natural thing
//     to write, and the empty 0 slot would be a perpetual nuiscance.
//
//   o Field indices $[1], ..., $[NF] match AWK as well as NIDX format.
//     The 1-up indexing for NIDX format was in turn devised specifically
//     to match behavior of Unix tools like cut, sort, and so on -- all
//     of which are 1-up.
//
//   o No zero index: AWK uses $0 like Miller uses $*, to refer
//     to the entire record. In Miller, 0 is never a valid array index but one
//     can use string(index) to get insertion-ordered hash maps:
//     "0" is as valid a key as "1" or "100".
//
// Why integer indices at all?
//
// * AWK indices are always stringified and "arrays" are always associative,
//   i.e. hashmaps. The Miller DSL is, essentially, a name-indexed AWK-ish
//   processor. And, Miller versions up to Miller 5 stringified int indices
//   awkishly.
//
// * But this is a new era; JSON is now widespread; people want arrays
//   (per se) in their JSON files to be passed through as such.
//
// Naming conventions:
//
// * All userspace indexing is done in this file, ultimately through
//   the unaliasArrayIndex() function.
//
// * Outside of this file I simply say 'index'.
//
// * Inside this file I say 'zindex' for the 0-up Go indices and 'mindex'
//   for 1-up Miller indices.
//
// ================================================================

package types

import (
	"errors"
	"strconv"

	"miller/lib"
)

// ================================================================
func (this *Mlrval) ArrayGet(mindex *Mlrval) Mlrval {
	if this.mvtype != MT_ARRAY {
		return MlrvalFromError()
	}
	if mindex.mvtype != MT_INT {
		return MlrvalFromError()
	}
	value := arrayGetAliased(&this.arrayval, mindex.intval)
	if value == nil {
		return MlrvalFromAbsent()
	} else {
		return *value
	}
}

// ----------------------------------------------------------------
func (this *Mlrval) ArrayPut(mindex *Mlrval, value *Mlrval) {
	if this.mvtype != MT_ARRAY {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
		return
	}
	if mindex.mvtype != MT_INT {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
		return
	}

	ok := arrayPutAliased(&this.arrayval, mindex.intval, value)
	if ok {
	} else {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
	}
}

// ----------------------------------------------------------------
func arrayGetAliased(array *[]Mlrval, mindex int64) *Mlrval {
	zindex, ok := unaliasArrayIndex(array, mindex)
	if ok {
		return &(*array)[zindex]
	} else {
		return nil
	}
}

func arrayPutAliased(array *[]Mlrval, mindex int64, value *Mlrval) bool {
	zindex, ok := unaliasArrayIndex(array, mindex)
	if ok {
		clone := value.Copy()
		(*array)[zindex] = *clone
		return true
	} else {
		return false
	}
}

func unaliasArrayIndex(array *[]Mlrval, mindex int64) (int64, bool) {
	n := int64(len(*array))
	return UnaliasArrayLengthIndex(n, mindex)
}

func UnaliasArrayLengthIndex(n int64, mindex int64) (int64, bool) {
	if 1 <= mindex && mindex <= n {
		zindex := mindex - 1
		return zindex, true
	} else if -n <= mindex && mindex <= -1 {
		zindex := mindex + n
		return zindex, true
	} else {
		return -999, false
	}
}

// ----------------------------------------------------------------
// TODO: thinking about capacity-resizing
func (this *Mlrval) ArrayAppend(value *Mlrval) {
	if this.mvtype != MT_ARRAY {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
		return
	}
	this.arrayval = append(this.arrayval, *value)
}

// ================================================================
func (this *Mlrval) MapGet(key *Mlrval) Mlrval {
	if this.mvtype != MT_MAP {
		return MlrvalFromError()
	}

	// Support positional indices, e.g. '$*[3]' is the same as '$[3]'.
	mval, err := this.mapval.GetWithMlrvalIndex(key)
	if err != nil { // xxx maybe error-return in the API
		return MlrvalFromError()
	}
	if mval == nil {
		return MlrvalFromAbsent()
	}
	// This returns a reference, not a (deep) copy. In general in Miller, we
	// copy only on write/put.
	return *mval
}

// ----------------------------------------------------------------
func (this *Mlrval) MapPut(key *Mlrval, value *Mlrval) {
	if this.mvtype != MT_MAP {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
		return
	}
	if key.mvtype != MT_STRING {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
		return
	}
	this.mapval.PutCopy(&key.printrep, value)
}

// ----------------------------------------------------------------
// This is a multi-level map/array put.
//
// E.g. '$name[1]["foo"] = "bar"' or '$*["foo"][1] = "bar"' In the former case
// the indices are ["name", 1, "foo"] and in the latter case the indices are
// ["foo", 1]. See also indexed-lvalues.md.
//
// There is auto-create for not-yet-populated levels, so for example if '@a' is
// already of type map, then on assignment of '@a["b"][2]["c"] = "d"' we'll
// create a map key "b", pointing to an array whose slot 2 is a map from "c" to
// "d".
//
// * If it's a map-type mlrval then:
//
//   o Strings are map keys.
//   o Integers are interpreted as positional indices, only into existing
//     fields. (We don't auto-deepen nested maps by positional indices.)
//
// * If it's an array-type mlrval then:
//
//   o Integers are array indices
//   o We auto-create/auto-lengthen/auto-deepen.
//   o If '@foo[]' has indices 0..3, then on '@foo[4] = "new"' we lengthen
//     the array by one.
//   o If '@foo[]' has indices 0..3, then on '@foo[6] = "new"'
//     we lengthen the array and absent-fill the intervenings.
//   o If '@foo["bar"]' does not exist, then on '@foo["bar"][0] = "new"'
//     we create an array and populate the 0th slot.
//   o If '@foo["bar"]' does not exist, then on '@foo["bar"][6] = "new"'
//     we create an array and populate the 6th slot, and absent-fill
//     the intervenings.
//
// * If this is a non-collection mlrval like string/int/float/etc.  then
//   it's non-indexable.
//
// See also indexed-lvalues.md.

func (this *Mlrval) PutIndexed(indices []*Mlrval, rvalue *Mlrval) error {
	lib.InternalCodingErrorIf(len(indices) < 1)

	if this.mvtype == MT_MAP {
		return putIndexedOnMap(this.mapval, indices, rvalue)

	} else if this.mvtype == MT_ARRAY {
		return putIndexedOnArray(&this.arrayval, indices, rvalue)

	} else {
		baseIndex := indices[0]
		if baseIndex.mvtype == MT_STRING {
			*this = MlrvalEmptyMap()
			return putIndexedOnMap(this.mapval, indices, rvalue)
		} else if baseIndex.mvtype == MT_INT {
			*this = MlrvalEmptyArray()
			return putIndexedOnArray(&this.arrayval, indices, rvalue)
		} else {
			return errors.New(
				"Miller: only maps and arrays are indexable; got " + this.GetTypeName(),
			)
		}
	}

	return nil
}

// ----------------------------------------------------------------
// Helper function for Mlrval.PutIndexed, for mlrvals of map type.
func putIndexedOnMap(baseMap *Mlrmap, indices []*Mlrval, rvalue *Mlrval) error {
	numIndices := len(indices)

	if numIndices == 0 {
		// E.g. mlr put '$* = {"a":1, "b":2}'
		if !rvalue.IsMap() {
			return errors.New(
				"Cannot assign non-map to existing map; got " +
					rvalue.GetTypeName() +
					".",
			)
		}
		*baseMap = *rvalue.mapval.Copy()
		return nil
	}

	baseIndex := indices[0]
	// If last index, then assign.
	if numIndices == 1 {
		// E.g. mlr put '$*["a"] = 3'
		return baseMap.PutCopyWithMlrvalIndex(baseIndex, rvalue)
	}

	// If not last index, then recurse.
	if baseIndex.mvtype == MT_STRING {
		// Base is map, index is string
		baseValue := baseMap.Get(&baseIndex.printrep)
		if baseValue == nil {
			// Create a new level in order to recurse from
			nextIndex := indices[1]

			var err error = nil
			baseValue, err = NewMlrvalForAutoDeepen(nextIndex.mvtype)
			if err != nil {
				return err
			}
			baseMap.PutReference(&baseIndex.printrep, baseValue)
		}
		return baseValue.PutIndexed(indices[1:], rvalue)

	} else if baseIndex.mvtype == MT_INT {
		// Base is map, index is int
		baseValue := baseMap.GetWithPositionalIndex(baseIndex.intval)
		if baseValue == nil {
			// There is no auto-deepen for positional indices on maps
			return errors.New(
				"Miller: positional index " +
					strconv.Itoa(int(baseIndex.intval)) +
					" not found.",
			)
		}
		baseValue.PutIndexed(indices[1:], rvalue)

	} else {
		// Base is map, index is invalid type
		return errors.New(
			"Miller: map indices must be string or positional int; got " + baseIndex.GetTypeName(),
		)
	}

	return nil
}

// ----------------------------------------------------------------
// Helper function for Mlrval.PutIndexed, for mlrvals of array type.
func putIndexedOnArray(
	baseArray *[]Mlrval,
	indices []*Mlrval,
	rvalue *Mlrval,
) error {

	numIndices := len(indices)
	lib.InternalCodingErrorIf(numIndices < 1)
	mindex := indices[0]

	if mindex.mvtype != MT_INT {
		return errors.New(
			"Array index must be int, but was " +
				mindex.GetTypeName() +
				".",
		)
	}
	zindex, inBounds := unaliasArrayIndex(baseArray, mindex.intval)

	if numIndices == 1 {
		// If last index, then assign.
		if inBounds {
			(*baseArray)[zindex] = *rvalue.Copy()
		} else if mindex.intval == 0 {
			return errors.New("Miller: zero indices are not supported. Indices are 1-up.")
		} else if mindex.intval < 0 {
			return errors.New("Miller: Cannot use negative indices to auto-lengthen arrays.")
		} else {
			// Array is [a,b,c] with mindices 1,2,3. Length is 3. Zindices are 0,1,2.
			// Given mindex is 4.
			LengthenMlrvalArray(baseArray, mindex.intval)
			zindex := mindex.intval - 1
			(*baseArray)[zindex] = *rvalue.Copy()
		}
		return nil

	} else {
		// More indices remain; recurse
		if inBounds {
			nextIndex := indices[1]

			// Overwrite what's in this slot if it's the wrong type
			if nextIndex.mvtype == MT_STRING {
				if (*baseArray)[zindex].mvtype != MT_MAP {
					(*baseArray)[zindex] = MlrvalEmptyMap()
				}
			} else if nextIndex.mvtype == MT_INT {
				if (*baseArray)[zindex].mvtype != MT_ARRAY {
					(*baseArray)[zindex] = MlrvalEmptyArray()
				}
			} else {
				return errors.New(
					"Miller: indices must be string or int; got " + nextIndex.GetTypeName(),
				)
			}

			return (*baseArray)[zindex].PutIndexed(indices[1:], rvalue)

		} else if mindex.intval == 0 {
			return errors.New("Miller: zero indices are not supported. Indices are 1-up.")
		} else if mindex.intval < 0 {
			return errors.New("Miller: Cannot use negative indices to auto-lengthen arrays.")
		} else {
			// Already allocated but needs to be longer
			LengthenMlrvalArray(baseArray, mindex.intval)
			zindex := mindex.intval - 1
			return (*baseArray)[zindex].PutIndexed(indices[1:], rvalue)
		}

	}

	return nil
}

// ----------------------------------------------------------------
func (this *Mlrval) UnsetIndexed(indices []*Mlrval) error {
	lib.InternalCodingErrorIf(len(indices) < 1)

	if this.mvtype == MT_MAP {
		return unsetIndexedOnMap(this.mapval, indices)

	} else if this.mvtype == MT_ARRAY {
		return unsetIndexedOnArray(&this.arrayval, indices)

	} else {
		return errors.New(
			"Miller: cannot unset index variable which is neither map nor array.",
		)
	}
}

// ----------------------------------------------------------------
// Helper function for Mlrval.UnsetIndexed, for mlrvals of map type.
func unsetIndexedOnMap(baseMap *Mlrmap, indices []*Mlrval) error {
	numIndices := len(indices)
	lib.InternalCodingErrorIf(numIndices < 1)
	baseIndex := indices[0]

	// If last index, then unset.
	if numIndices == 1 {
		if baseIndex.mvtype == MT_STRING {
			baseMap.Remove(&baseIndex.printrep)
			return nil
		} else if baseIndex.mvtype == MT_INT {
			baseMap.RemoveWithPositionalIndex(baseIndex.intval)
			return nil
		} else {
			return errors.New(
				"Miller: map indices must be string or positional int; got " +
					baseIndex.GetTypeName(),
			)
		}
	}

	// If not last index, then recurse.
	if baseIndex.mvtype == MT_STRING {
		// Base is map, index is string
		baseValue := baseMap.Get(&baseIndex.printrep)
		return baseValue.UnsetIndexed(indices[1:])

	} else if baseIndex.mvtype == MT_INT {
		// Base is map, index is int
		baseValue := baseMap.GetWithPositionalIndex(baseIndex.intval)
		baseValue.UnsetIndexed(indices[1:])

	} else {
		// Base is map, index is invalid type
		return errors.New(
			"Miller: map indices must be string or positional int; got " + baseIndex.GetTypeName(),
		)
	}

	return nil
}

// ----------------------------------------------------------------
// Helper function for Mlrval.PutIndexed, for mlrvals of array type.
func unsetIndexedOnArray(
	baseArray *[]Mlrval,
	indices []*Mlrval,
) error {
	numIndices := len(indices)
	lib.InternalCodingErrorIf(numIndices < 1)
	mindex := indices[0]

	if mindex.mvtype != MT_INT {
		return errors.New(
			"Array index must be int, but was " +
				mindex.GetTypeName() +
				".",
		)
	}
	zindex, inBounds := unaliasArrayIndex(baseArray, mindex.intval)

	// If last index, then unset.
	if numIndices == 1 {
		if inBounds {
			leftSlice := (*baseArray)[0:zindex]
			rightSlice := (*baseArray)[zindex+1 : len((*baseArray))]
			*baseArray = append(leftSlice, rightSlice...)
		} else if mindex.intval == 0 {
			return errors.New("Miller: zero indices are not supported. Indices are 1-up.")
		} else {
			// TODO: improve wording
			return errors.New("Miller: array index out of bounds for unset.")
		}
	} else {
		// More indices remain; recurse
		if inBounds {
			return (*baseArray)[zindex].UnsetIndexed(indices[1:])
		} else if mindex.intval == 0 {
			return errors.New("Miller: zero indices are not supported. Indices are 1-up.")
		} else {
			// TODO: improve wording
			return errors.New("Miller: array index out of bounds for unset.")
		}

	}

	return nil
}
