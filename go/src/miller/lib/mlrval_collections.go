package lib

import (
	"errors"
)

// ----------------------------------------------------------------
func (this *Mlrval) ArrayGet(index *Mlrval) Mlrval {
	if this.mvtype != MT_ARRAY {
		return MlrvalFromError()
	}
	if index.mvtype != MT_INT {
		return MlrvalFromError()
	}
	i := index.intval
	n := int64(len(this.arrayval))

	// TODO: document this (pythonic)
	if i < 0 && i > -n {
		i += n
	}

	if i < 0 || i >= n {
		return MlrvalFromError()
	}
	return this.arrayval[i]
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
	i := index.intval
	n := int64(len(this.arrayval))

	// TODO: document this (pythonic)
	if i < 0 && i > -n {
		i += n
	}
	if i < 0 || i >= n {
		// TODO: need to be careful about semantics here.
		// Silent no-ops are not good UX ...
		return
	}
	this.arrayval[i] = *value.Copy()
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
	if key.mvtype != MT_STRING {
		return MlrvalFromError()
	}

	mval := this.mapval.Get(&key.printrep)
	if mval == nil {
		return MlrvalFromAbsent()
	} else {
		// This returns a reference, not a copy. In general in Miller, we copy
		// only on write/put.
		return *mval
	}
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
// See also indexed-lvalues.md
func (this *Mlrval) PutIndexed(indices []*Mlrval, rvalue *Mlrval) error {
	n := len(indices)
	InternalCodingErrorIf(n < 1)

	levelMlrval := this

	// xxx temp -- at very first just do strings.
	for i, index := range indices {
		if !levelMlrval.IsMap() {
			return errors.New("indexed level not map") // xxx needs better messaging
		}
		levelMlrmap := levelMlrval.mapval

		if !index.IsString() {
			return errors.New("string-only indices for now, sorry!")
		}
		key := index.printrep

		nextLevelMlrval := levelMlrmap.Get(&key)
		if nextLevelMlrval == nil {
			if i < n-1 {
				next := MlrvalEmptyMap()
				nextLevelMlrval = &next
				levelMlrmap.PutCopy(&key, nextLevelMlrval)
			} else {
				levelMlrmap.PutCopy(&key, rvalue)
				break
			}
		}
		levelMlrval = nextLevelMlrval
	}

	return nil
}
