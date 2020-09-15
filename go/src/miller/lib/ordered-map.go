// ================================================================
// ORDERED MAP FROM STRING TO INTERFACE{}
//
// Quite like types.OrderedMap but only with interface{} keys. See orderedMap.go for
// more information.
// ================================================================

package lib

// ----------------------------------------------------------------
type OrderedMap struct {
	FieldCount    int64
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
func (this *OrderedMap) Has(key string) bool {
	return this.findEntry(&key) != nil
}

func (this *OrderedMap) findEntry(key *string) *orderedMapEntry {
	if this.keysToEntries != nil {
		return this.keysToEntries[*key]
	} else {
		for pe := this.Head; pe != nil; pe = pe.Next {
			if pe.Key == *key {
				return pe
			}
		}
		return nil
	}
}

// ----------------------------------------------------------------
func (this *OrderedMap) Put(key string, value interface{}) {
	pe := this.findEntry(&key)
	if pe == nil {
		pe = newOrderedMapEntry(&key, value)
		if this.Head == nil {
			this.Head = pe
			this.Tail = pe
		} else {
			pe.Prev = this.Tail
			pe.Next = nil
			this.Tail.Next = pe
			this.Tail = pe
		}
		if this.keysToEntries != nil {
			this.keysToEntries[key] = pe
		}
		this.FieldCount++
	} else {
		pe.Value = value
	}
}

// ----------------------------------------------------------------
func (this *OrderedMap) Get(key string) interface{} {
	pe := this.findEntry(&key)
	if pe == nil {
		return nil
	} else {
		return pe.Value
	}
	return nil
}

// ----------------------------------------------------------------
func (this *OrderedMap) Clear() {
	this.FieldCount = 0
	this.Head = nil
	this.Tail = nil
}

// ----------------------------------------------------------------
// Returns true if it was found and removed
func (this *OrderedMap) Remove(key string) bool {
	pe := this.findEntry(&key)
	if pe == nil {
		return false
	} else {
		this.unlink(pe)
		return true
	}
}

// ----------------------------------------------------------------
func (this *OrderedMap) unlink(pe *orderedMapEntry) {
	if pe == this.Head {
		if pe == this.Tail {
			this.Head = nil
			this.Tail = nil
		} else {
			this.Head = pe.Next
			pe.Next.Prev = nil
		}
	} else {
		pe.Prev.Next = pe.Next
		if pe == this.Tail {
			this.Tail = pe.Prev
		} else {
			pe.Next.Prev = pe.Prev
		}
	}
	if this.keysToEntries != nil {
		delete(this.keysToEntries, pe.Key)
	}
	this.FieldCount--
}
