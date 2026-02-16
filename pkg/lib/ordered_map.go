// ================================================================
// ORDERED MAP FROM STRING TO GENERIC VALUE TYPE
//
// Quite like types.OrderedMap but only with string keys. See orderedMap.go for
// more information.
// ================================================================

package lib

type OrderedMap[V any] struct {
	FieldCount    int64
	Head          *orderedMapEntry[V]
	Tail          *orderedMapEntry[V]
	keysToEntries map[string]*orderedMapEntry[V]
}

type orderedMapEntry[V any] struct {
	Key   string
	Value V
	Prev  *orderedMapEntry[V]
	Next  *orderedMapEntry[V]
}

func NewOrderedMap[V any]() *OrderedMap[V] {
	return &OrderedMap[V]{
		FieldCount:    0,
		Head:          nil,
		Tail:          nil,
		keysToEntries: make(map[string]*orderedMapEntry[V]),
	}
}

// ----------------------------------------------------------------
// Value-copy is up to the caller -- PutReference and PutCopy
// are in the public OrderedMap API.
func newOrderedMapEntry[V any](key *string, value V) *orderedMapEntry[V] {
	return &orderedMapEntry[V]{
		*key,
		value,
		nil,
		nil,
	}
}

func (omap *OrderedMap[V]) IsEmpty() bool {
	return omap.FieldCount == 0
}

func (omap *OrderedMap[V]) Has(key string) bool {
	return omap.findEntry(&key) != nil
}

func (omap *OrderedMap[V]) findEntry(key *string) *orderedMapEntry[V] {
	if omap.keysToEntries != nil {
		return omap.keysToEntries[*key]
	} else {
		for pe := omap.Head; pe != nil; pe = pe.Next {
			if pe.Key == *key {
				return pe
			}
		}
		return nil
	}
}

func (omap *OrderedMap[V]) Put(key string, value V) {
	pe := omap.findEntry(&key)
	if pe == nil {
		pe = newOrderedMapEntry(&key, value)
		if omap.Head == nil {
			omap.Head = pe
			omap.Tail = pe
		} else {
			pe.Prev = omap.Tail
			pe.Next = nil
			omap.Tail.Next = pe
			omap.Tail = pe
		}
		if omap.keysToEntries != nil {
			omap.keysToEntries[key] = pe
		}
		omap.FieldCount++
	} else {
		pe.Value = value
	}
}

func (omap *OrderedMap[V]) Get(key string) V {
	pe := omap.findEntry(&key)
	if pe == nil {
		var zero V
		return zero
	}
	return pe.Value
}

// The Get is sufficient for pointer values -- the caller can check if the
// return value is nil. For int/string values (which are non-nullable) we have
// this method.
func (omap *OrderedMap[V]) GetWithCheck(key string) (V, bool) {
	pe := omap.findEntry(&key)
	if pe == nil {
		var zero V
		return zero, false
	}
	return pe.Value, true
}

// ----------------------------------------------------------------
// Returns true if it was found and removed
func (omap *OrderedMap[V]) Remove(key string) bool {
	pe := omap.findEntry(&key)
	if pe == nil {
		return false
	} else {
		omap.unlink(pe)
		return true
	}
}

func (omap *OrderedMap[V]) unlink(pe *orderedMapEntry[V]) {
	if pe == omap.Head {
		if pe == omap.Tail {
			omap.Head = nil
			omap.Tail = nil
		} else {
			omap.Head = pe.Next
			pe.Next.Prev = nil
		}
	} else {
		pe.Prev.Next = pe.Next
		if pe == omap.Tail {
			omap.Tail = pe.Prev
		} else {
			pe.Next.Prev = pe.Prev
		}
	}
	if omap.keysToEntries != nil {
		delete(omap.keysToEntries, pe.Key)
	}
	omap.FieldCount--
}
