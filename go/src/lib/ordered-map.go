// ================================================================
// ORDERED MAP FROM STRING TO INTERFACE{}
//
// Quite like types.OrderedMap but only with interface{} keys. See orderedMap.go for
// more information.
// ================================================================

package lib

// ----------------------------------------------------------------
type OrderedMap struct {
	FieldCount    int
	Head          *orderedMapEntry
	Tail          *orderedMapEntry
	keysToEntries map[string]*orderedMapEntry
}

type orderedMapEntry struct {
	Key   string
	Value interface{}
	Prev  *orderedMapEntry
	Next  *orderedMapEntry
}

// ----------------------------------------------------------------
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		FieldCount:    0,
		Head:          nil,
		Tail:          nil,
		keysToEntries: make(map[string]*orderedMapEntry),
	}
}

// ----------------------------------------------------------------
// Value-copy is up to the caller -- PutReference and PutCopy
// are in the public OrderedMap API.
func newOrderedMapEntry(key *string, value interface{}) *orderedMapEntry {
	return &orderedMapEntry{
		*key,
		value,
		nil,
		nil,
	}
}

// ----------------------------------------------------------------
func (omap *OrderedMap) IsEmpty() bool {
	return omap.FieldCount == 0
}

func (omap *OrderedMap) Has(key string) bool {
	return omap.findEntry(&key) != nil
}

func (omap *OrderedMap) findEntry(key *string) *orderedMapEntry {
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

// ----------------------------------------------------------------
func (omap *OrderedMap) Put(key string, value interface{}) {
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

// ----------------------------------------------------------------
func (omap *OrderedMap) Get(key string) interface{} {
	pe := omap.findEntry(&key)
	if pe == nil {
		return nil
	} else {
		return pe.Value
	}
	return nil
}

// The Get is sufficient for pointer values -- the caller can check if the
// return value is nil. For int/string values (which are non-nullable) we have
// this method.
func (omap *OrderedMap) GetWithCheck(key string) (interface{}, bool) {
	pe := omap.findEntry(&key)
	if pe == nil {
		return nil, false
	} else {
		return pe.Value, true
	}
	return nil, false
}

// ----------------------------------------------------------------
func (omap *OrderedMap) Clear() {
	omap.FieldCount = 0
	omap.Head = nil
	omap.Tail = nil
}

// ----------------------------------------------------------------
// Returns true if it was found and removed
func (omap *OrderedMap) Remove(key string) bool {
	pe := omap.findEntry(&key)
	if pe == nil {
		return false
	} else {
		omap.unlink(pe)
		return true
	}
}

// ----------------------------------------------------------------
func (omap *OrderedMap) unlink(pe *orderedMapEntry) {
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
