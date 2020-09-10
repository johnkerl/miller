package lib

import (
	"errors"
	"strconv"
)

// ================================================================
func (this *Mlrval) ArrayGet(index *Mlrval) Mlrval {
	if this.mvtype != MT_ARRAY {
		return MlrvalFromError()
	}
	if index.mvtype != MT_INT {
		return MlrvalFromError()
	}
	value := arrayGetAliased(&this.arrayval, index.intval)
	if value == nil {
		return MlrvalFromError()
	} else {
		return *value
	}
}

// ----------------------------------------------------------------
func (this *Mlrval) ArrayPut(index *Mlrval, value *Mlrval) {
	if this.mvtype != MT_ARRAY {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
		return
	}
	if index.mvtype != MT_INT {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
		return
	}

	ok := arrayPutAliased(&this.arrayval, index.intval, value)
	if ok {
	} else {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
	}
}

// ----------------------------------------------------------------
func arrayGetAliased(array *[]Mlrval, i int64) *Mlrval {
	index, ok := unaliasArrayIndex(array, i)
	if ok {
		return &(*array)[index]
	} else {
		return nil
	}
}

func arrayPutAliased(array *[]Mlrval, i int64, value *Mlrval) bool {
	index, ok := unaliasArrayIndex(array, i)
	if ok {
		clone := value.Copy()
		(*array)[index] = *clone
		return true
	} else {
		return false
	}
}

func unaliasArrayIndex(array *[]Mlrval, i int64) (int64, bool) {
	n := int64(len(*array))
	// TODO: document this (pythonic)
	if i < 0 && i > -n {
		i += n
	}
	if i < 0 || i >= n {
		return -999, false
	} else {
		return i, true
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
//     fields. (We don't auto-extend maps by positional indices.)
//
// * If it's an array-type mlrval then:
//
//   o Integers are array indices
//   o We auto-create/auto-extend.
//   o If '@foo[]' has indices 0..3, then on '@foo[4] = "new"' we extend
//     the array by one.
//   o If '@foo[]' has indices 0..3, then on '@foo[6] = "new"'
//     we extend the array and absent-fill the intervenings.
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
	InternalCodingErrorIf(len(indices) < 1)

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
			baseValue, err = NewMlrvalForAutoExtend(nextIndex.mvtype)
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
			// There is no auto-extend for positional indices on maps
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
	InternalCodingErrorIf(numIndices < 1)
	index := indices[0]

	if index.mvtype != MT_INT {
		return errors.New(
			"Array index must be int, but was " +
				index.GetTypeName() +
				".",
		)
	}
	intIndex, inBounds := unaliasArrayIndex(baseArray, index.intval)

	if numIndices == 1 {
		// If last index, then assign.
		if inBounds {
			(*baseArray)[intIndex] = *rvalue.Copy()
		} else if index.intval < 0 {
			return errors.New("Cannot use negative indices to auto-extend arrays")
		} else {
			LengthenMlrvalArray(baseArray, index.intval+1)
			(*baseArray)[index.intval] = *rvalue.Copy()
		}
		return nil

	} else {
		// More indices remain; recurse
		if inBounds {
			nextIndex := indices[1]

			// Overwrite what's in this slot if it's the wrong type
			if nextIndex.mvtype == MT_STRING {
				if (*baseArray)[intIndex].mvtype != MT_MAP {
					(*baseArray)[intIndex] = MlrvalEmptyMap()
				}
			} else if nextIndex.mvtype == MT_INT {
				if (*baseArray)[intIndex].mvtype != MT_ARRAY {
					(*baseArray)[intIndex] = MlrvalEmptyArray()
				}
			} else {
				return errors.New(
					"Indices must be string or int; got " + nextIndex.GetTypeName(),
				)
			}

			return (*baseArray)[intIndex].PutIndexed(indices[1:], rvalue)

		} else if index.intval < 0 {
			return errors.New("Cannot use negative indices to auto-extend arrays")
		} else {
			// Already allocated but needs to be longer
			LengthenMlrvalArray(baseArray, index.intval+1)
			return (*baseArray)[index.intval].PutIndexed(indices[1:], rvalue)
		}

	}

	return nil
}
