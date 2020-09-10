package lib

import (
	"errors"
	"strconv"
)

// ----------------------------------------------------------------
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
		return 0, false
	} else {
		return i, true
	}
}

// ----------------------------------------------------------------
// TODO: thinking about capacity-resizing
func (this *Mlrval) ArrayExtend(value *Mlrval) {
	if this.mvtype != MT_ARRAY {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
		return
	}
	this.arrayval = append(this.arrayval, *value)
}

// ----------------------------------------------------------------
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
// This is a multi-level map/array put. There is auto-create for
// not-yet-populated levels, so for example if '@a' is already of type map,
// then on assignment of '@a["b"][2]["c"] = "d"' we'll create a map key "b",
// pointing to an array whose slot 2 is a map from "c" to "d".
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

	if this.mvtype == MT_MAP {
		return putIndexedOnMap(this.mapval, indices, rvalue)

	} else if this.mvtype == MT_ARRAY {
		return putIndexedOnArray(&this.arrayval, indices, rvalue)

	} else {
		return errors.New( // xxx libify
			"Only maps and arrays are indexable; got " + this.GetTypeName(),
		)
	}

	return nil
}

// ----------------------------------------------------------------
// Helper function for Mlrval.PutIndexed, for mlrvals of map type.
func putIndexedOnMap(baseMap *Mlrmap, indices []*Mlrval, rvalue *Mlrval) error {
	n := len(indices)
	InternalCodingErrorIf(n < 1)
	index := indices[0]

	// If last index, then assign.
	if n == 1 {
		return baseMap.PutCopyWithMlrvalIndex(index, rvalue)
	}

	// If not last index, then recurse.
	if index.mvtype == MT_STRING {
		// Base is map, index is string
		value := baseMap.Get(&index.printrep)
		if value == nil {
			// Create/auto-extend a new map level in order to recurse from
			empty := MlrvalEmptyMap()
			value = &empty
			baseMap.PutCopyWithMlrvalIndex(index, value)
		}
		value.PutIndexed(indices[1:], rvalue)

	} else if index.mvtype == MT_INT {
		// Base is map, index is int
		value := baseMap.GetWithPositionalIndex(index.intval)
		if value == nil {
			// There is no auto-extend for positional indices on maps
			return errors.New( // xxx libify
				"Positional index " +
					strconv.Itoa(int(index.intval)) +
					" not found.",
			)
		}
		value.PutIndexed(indices[1:], rvalue)

	} else {
		// Base is map, index is invalid type
		return errors.New( // xxx libify
			"Map indices must be string or positional int; got " + index.GetTypeName(),
		)
	}

	return nil
}

// ----------------------------------------------------------------
// Helper function for Mlrval.PutIndexed, for mlrvals of array type.
func putIndexedOnArray(baseArray *[]Mlrval, indices []*Mlrval, rvalue *Mlrval) error {
	n := len(indices)
	InternalCodingErrorIf(n < 1)
	index := indices[0]

	if index.mvtype != MT_INT {
		return errors.New(
			"Array index must be int, but was " +
				index.GetTypeName() +
				".",
		)
	}
	intIndex, inBounds := unaliasArrayIndex(baseArray, index.intval)

	if n == 1 {
		// If last index, then assign.
		if inBounds {
			clone := rvalue.Copy()
			(*baseArray)[intIndex] = *clone
		} else if index.intval < 0 {
			return errors.New("Cannot use negative indices to auto-extend arrays")
		} else{
			// TODO: auto-extend ...
			return errors.New("array auto-extend not yet implemented")
		}
		return nil

	} else {
		// More indices remain; recurse
		if inBounds {
			return (*baseArray)[intIndex].PutIndexed(indices[1:], rvalue)
		} else if index.intval < 0 {
			return errors.New("Cannot use negative indices to auto-extend arrays")
		} else {
			// TODO: auto-extend ...
			return errors.New("array auto-extend not yet implemented")
		}
	}

	return nil
}
