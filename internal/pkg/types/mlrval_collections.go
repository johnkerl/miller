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
//   the UnaliasArrayIndex() function.
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
	"fmt"
	"os"

	"mlr/internal/pkg/lib"
)

// ================================================================
// TODO: copy-reduction refactor
func (mv *Mlrval) ArrayGet(mindex *Mlrval) Mlrval {
	if mv.mvtype != MT_ARRAY {
		return *MLRVAL_ERROR
	}
	if mindex.mvtype != MT_INT {
		return *MLRVAL_ERROR
	}
	value := arrayGetAliased(&mv.arrayval, mindex.intval)
	if value == nil {
		return *MLRVAL_ABSENT
	} else {
		return *value
	}
}

// ----------------------------------------------------------------
// TODO: make this return error so caller can do 'if err == nil { ... }'
func (mv *Mlrval) ArrayPut(mindex *Mlrval, value *Mlrval) {
	if mv.mvtype != MT_ARRAY {
		fmt.Fprintf(
			os.Stderr,
			"mlr: expected array as indexed item in ArrayPut; got %s\n",
			mv.GetTypeName(),
		)
		os.Exit(1)
	}
	if mindex.mvtype != MT_INT {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
		fmt.Fprintf(
			os.Stderr,
			"mlr: expected int as index item ArrayPut; got %s\n",
			mindex.GetTypeName(),
		)
		os.Exit(1)
	}

	ok := arrayPutAliased(&mv.arrayval, mindex.intval, value)
	if !ok {
		fmt.Fprintf(
			os.Stderr,
			"mlr: array index %d out of bounds %d..%d\n",
			mindex.intval, 1, len(mv.arrayval),
		)
		os.Exit(1)
	}
}

// ----------------------------------------------------------------
func arrayGetAliased(array *[]Mlrval, mindex int) *Mlrval {
	zindex, ok := UnaliasArrayIndex(array, mindex)
	if ok {
		return &(*array)[zindex]
	} else {
		return nil
	}
}

func arrayPutAliased(array *[]Mlrval, mindex int, value *Mlrval) bool {
	zindex, ok := UnaliasArrayIndex(array, mindex)
	if ok {
		clone := value.Copy()
		(*array)[zindex] = *clone
		return true
	} else {
		return false
	}
}

func UnaliasArrayIndex(array *[]Mlrval, mindex int) (int, bool) {
	n := int(len(*array))
	return UnaliasArrayLengthIndex(n, mindex)
}

// Input "mindex" is a Miller DSL array index. These are 1-up, so 1..n where n
// is the length of the array. Also, -n..-1 are aliases to 1..n. 0 is never a
// valid index.
//
// Output "zindex" is a Golang array index. These are 0-up, so 0..(n-1).
//
// The second return value indicates whether the Miller index is in-bounds.
// Even if it's out of bounds, while the second return value is false, the
// first return is correctly de-aliased. E.g. if the array has length 5 and the
// mindex is 8, zindex is 7 and valid=false. This is so in array-slice
// operations like 'v = myarray[2:8]' the callsite can hand back slots 2-5 of
// the array (which is the same way Python handles beyond-the-end indexing).

// Examples with n = 5:
//
// mindex zindex ok
// -7    -2      false
// -6    -1      false
// -5     0      true
// -4     1      true
// -3     2      true
// -2     3      true
// -1     4      true
//  0    -1      false
//  1     0      true
//  2     1      true
//  3     2      true
//  4     3      true
//  5     4      true
//  6     5      false
//  7     6      false

func UnaliasArrayLengthIndex(n int, mindex int) (int, bool) {
	if 1 <= mindex {
		zindex := mindex - 1
		if mindex <= n { // in bounds
			return zindex, true
		} else { // out of bounds
			return zindex, false
		}
	} else if mindex <= -1 {
		zindex := mindex + n
		if -n <= mindex { // in bounds
			return zindex, true
		} else { // out of bounds
			return zindex, false
		}
	} else {
		// mindex is 0
		return -1, false
	}
}

// ----------------------------------------------------------------
// TODO: thinking about capacity-resizing
func (mv *Mlrval) ArrayAppend(value *Mlrval) {
	if mv.mvtype != MT_ARRAY {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
		return
	}
	mv.arrayval = append(mv.arrayval, *value)
}

// ================================================================
func (mv *Mlrval) MapGet(key *Mlrval) Mlrval {
	if mv.mvtype != MT_MAP {
		return *MLRVAL_ERROR
	}

	mval, err := mv.mapval.GetWithMlrvalIndex(key)
	if err != nil { // xxx maybe error-return in the API
		return *MLRVAL_ERROR
	}
	if mval == nil {
		return *MLRVAL_ABSENT
	}
	// This returns a reference, not a (deep) copy. In general in Miller, we
	// copy only on write/put.
	return *mval
}

// ----------------------------------------------------------------
func (mv *Mlrval) MapPut(key *Mlrval, value *Mlrval) {
	if mv.mvtype != MT_MAP {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
		return
	}

	if key.mvtype == MT_STRING {
		mv.mapval.PutCopy(key.printrep, value)
	} else if key.mvtype == MT_INT {
		mv.mapval.PutCopy(key.String(), value)
	}
	// TODO: need to be careful about semantics here.
	// Silent no-ops are not good UX ...
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
//   o Integers are stringified, then interpreted as map keys.
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

func (mv *Mlrval) PutIndexed(indices []*Mlrval, rvalue *Mlrval) error {
	lib.InternalCodingErrorIf(len(indices) < 1)

	if mv.mvtype == MT_MAP {
		return putIndexedOnMap(mv.mapval, indices, rvalue)

	} else if mv.mvtype == MT_ARRAY {
		return putIndexedOnArray(&mv.arrayval, indices, rvalue)

	} else {
		baseIndex := indices[0]
		if baseIndex.mvtype == MT_STRING {
			*mv = *MlrvalFromEmptyMap()
			return putIndexedOnMap(mv.mapval, indices, rvalue)
		} else if baseIndex.mvtype == MT_INT {
			*mv = MlrvalEmptyArray()
			return putIndexedOnArray(&mv.arrayval, indices, rvalue)
		} else {
			return errors.New(
				"mlr: only maps and arrays are indexable; got " + mv.GetTypeName(),
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
	if baseIndex.mvtype != MT_STRING && baseIndex.mvtype != MT_INT {
		// Base is map, index is invalid type
		return errors.New(
			"mlr: map indices must be string, int, or array thereof; got " + baseIndex.GetTypeName(),
		)
	}

	baseValue := baseMap.Get(baseIndex.String())
	if baseValue == nil {
		// Create a new level in order to recurse from
		nextIndex := indices[1]

		var err error = nil
		baseValue, err = NewMlrvalForAutoDeepen(nextIndex.mvtype)
		if err != nil {
			return err
		}
		baseMap.PutReference(baseIndex.String(), baseValue)
	}
	return baseValue.PutIndexed(indices[1:], rvalue)
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
	zindex, inBounds := UnaliasArrayIndex(baseArray, mindex.intval)

	if numIndices == 1 {
		// If last index, then assign.
		if inBounds {
			(*baseArray)[zindex] = *rvalue.Copy()
		} else if mindex.intval == 0 {
			return errors.New("mlr: zero indices are not supported. Indices are 1-up.")
		} else if mindex.intval < 0 {
			return errors.New("mlr: Cannot use negative indices to auto-lengthen arrays.")
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
					(*baseArray)[zindex] = *MlrvalFromEmptyMap()
				}
			} else if nextIndex.mvtype == MT_INT {
				if (*baseArray)[zindex].mvtype != MT_ARRAY {
					(*baseArray)[zindex] = MlrvalEmptyArray()
				}
			} else {
				return errors.New(
					"mlr: indices must be string, int, or array thereof; got " + nextIndex.GetTypeName(),
				)
			}

			return (*baseArray)[zindex].PutIndexed(indices[1:], rvalue)

		} else if mindex.intval == 0 {
			return errors.New("mlr: zero indices are not supported. Indices are 1-up.")
		} else if mindex.intval < 0 {
			return errors.New("mlr: Cannot use negative indices to auto-lengthen arrays.")
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
func (mv *Mlrval) RemoveIndexed(indices []*Mlrval) error {
	lib.InternalCodingErrorIf(len(indices) < 1)

	if mv.mvtype == MT_MAP {
		return removeIndexedOnMap(mv.mapval, indices)

	} else if mv.mvtype == MT_ARRAY {
		return removeIndexedOnArray(&mv.arrayval, indices)

	} else {
		return errors.New(
			"mlr: cannot unset index variable which is neither map nor array.",
		)
	}
}

// ----------------------------------------------------------------
// Helper function for Mlrval.RemoveIndexed, for mlrvals of map type.
func removeIndexedOnMap(baseMap *Mlrmap, indices []*Mlrval) error {
	numIndices := len(indices)
	lib.InternalCodingErrorIf(numIndices < 1)
	baseIndex := indices[0]

	// If last index, then unset.
	if numIndices == 1 {
		if baseIndex.mvtype == MT_STRING || baseIndex.mvtype == MT_INT {
			baseMap.Remove(baseIndex.String())
			return nil
		} else {
			return errors.New(
				"mlr: map indices must be string, int, or array thereof; got " +
					baseIndex.GetTypeName(),
			)
		}
	}

	// If not last index, then recurse.
	if baseIndex.mvtype == MT_STRING || baseIndex.mvtype == MT_INT {
		// Base is map, index is string
		baseValue := baseMap.Get(baseIndex.String())
		if baseValue != nil {
			return baseValue.RemoveIndexed(indices[1:])
		}

	} else {
		// Base is map, index is invalid type
		return errors.New(
			"mlr: map indices must be string, int, or array thereof; got " + baseIndex.GetTypeName(),
		)
	}

	return nil
}

// ----------------------------------------------------------------
// Helper function for Mlrval.PutIndexed, for mlrvals of array type.
func removeIndexedOnArray(
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
	zindex, inBounds := UnaliasArrayIndex(baseArray, mindex.intval)

	// If last index, then unset.
	if numIndices == 1 {
		if inBounds {
			leftSlice := (*baseArray)[0:zindex]
			rightSlice := (*baseArray)[zindex+1 : len((*baseArray))]
			*baseArray = append(leftSlice, rightSlice...)
		} else if mindex.intval == 0 {
			return errors.New("mlr: zero indices are not supported. Indices are 1-up.")
		} else {
			// TODO: improve wording
			return errors.New("mlr: array index out of bounds for unset.")
		}
	} else {
		// More indices remain; recurse
		if inBounds {
			return (*baseArray)[zindex].RemoveIndexed(indices[1:])
		} else if mindex.intval == 0 {
			return errors.New("mlr: zero indices are not supported. Indices are 1-up.")
		} else {
			// TODO: improve wording
			return errors.New("mlr: array index out of bounds for unset.")
		}

	}

	return nil
}

// ----------------------------------------------------------------
// Used for API-matching in multi-index contexts.
func MakePointerArray(
	valueArray []Mlrval,
) (pointerArray []*Mlrval) {
	pointerArray = make([]*Mlrval, len(valueArray))
	for i := range valueArray {
		pointerArray[i] = &valueArray[i]
	}
	return pointerArray
}

// ----------------------------------------------------------------
// Nominally for TopKeeper

type BsearchMlrvalArrayFunc func(
	array *[]*Mlrval,
	size int, // maybe less than len(array)
	value *Mlrval,
) int

func BsearchMlrvalArrayForDescendingInsert(
	array *[]*Mlrval,
	size int, // maybe less than len(array)
	value *Mlrval,
) int {
	lo := 0
	hi := size - 1
	mid := (hi + lo) / 2
	var newmid int

	if size == 0 {
		return 0
	}

	if MlrvalGreaterThanAsBool(value, (*array)[0]) {
		return 0
	}
	if MlrvalLessThanAsBool(value, (*array)[hi]) {
		return size
	}

	for lo < hi {
		middleElement := (*array)[mid]
		if MlrvalEqualsAsBool(value, middleElement) {
			return mid
		} else if MlrvalGreaterThanAsBool(value, middleElement) {
			hi = mid
			newmid = (hi + lo) / 2
		} else {
			lo = mid
			newmid = (hi + lo) / 2
		}
		if mid == newmid {
			if MlrvalGreaterThanOrEqualsAsBool(value, (*array)[lo]) {
				return lo
			} else if MlrvalGreaterThanOrEqualsAsBool(value, (*array)[hi]) {
				return hi
			} else {
				return hi + 1
			}
		}
		mid = newmid
	}

	return lo
}

func BsearchMlrvalArrayForAscendingInsert(
	array *[]*Mlrval,
	size int, // maybe less than len(array)
	value *Mlrval,
) int {
	lo := 0
	hi := size - 1
	mid := (hi + lo) / 2
	var newmid int

	if size == 0 {
		return 0
	}

	if MlrvalLessThanAsBool(value, (*array)[0]) {
		return 0
	}
	if MlrvalGreaterThanAsBool(value, (*array)[hi]) {
		return size
	}

	for lo < hi {
		middleElement := (*array)[mid]
		if MlrvalEqualsAsBool(value, middleElement) {
			return mid
		} else if MlrvalLessThanAsBool(value, middleElement) {
			hi = mid
			newmid = (hi + lo) / 2
		} else {
			lo = mid
			newmid = (hi + lo) / 2
		}
		if mid == newmid {
			if MlrvalLessThanOrEqualsAsBool(value, (*array)[lo]) {
				return lo
			} else if MlrvalLessThanOrEqualsAsBool(value, (*array)[hi]) {
				return hi
			} else {
				return hi + 1
			}
		}
		mid = newmid
	}

	return lo
}
