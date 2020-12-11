package types

import (
	"bytes"
	"errors"

	"miller/lib"
)

// ----------------------------------------------------------------
func (this *Mlrmap) Has(key *string) bool {
	return this.findEntry(key) != nil
}

// ----------------------------------------------------------------
// Copies the key and value (deep-copying in case the value is array/map).
// This is safe for DSL use. See also PutReference.
func (this *Mlrmap) PutCopy(key *string, value *Mlrval) {
	pe := this.findEntry(key)
	if pe == nil {
		pe = newMlrmapEntry(key, value.Copy())
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
			this.keysToEntries[*key] = pe
		}
		this.FieldCount++
	} else {
		pe.Value = value.Copy()
	}
}

// Copies the key but not the value. This is not safe for DSL use, where we
// could create undesired references between different objects.  Only intended
// to be used at callsites which allocate a mlrval solely for the purpose of
// putting into a map, e.g. input-record readers.
func (this *Mlrmap) PutReference(key *string, value *Mlrval) {
	pe := this.findEntry(key)
	if pe == nil {
		pe = newMlrmapEntry(key, value)
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
			this.keysToEntries[*key] = pe
		}
		this.FieldCount++
	} else {
		pe.Value = value
	}
}

func (this *Mlrmap) PutCopyWithMlrvalIndex(key *Mlrval, value *Mlrval) error {
	if key.mvtype == MT_STRING {
		this.PutCopy(&key.printrep, value)
		return nil
	} else if key.mvtype == MT_INT {
		s := key.String()
		this.PutCopy(&s, value)
		return nil
	} else {
		return errors.New(
			"Miller: record/map indices must be string or int; got " + key.GetTypeName(),
		)
	}
}

// ----------------------------------------------------------------
func (this *Mlrmap) PrependCopy(key *string, value *Mlrval) {
	pe := this.findEntry(key)
	if pe == nil {
		pe = newMlrmapEntry(key, value)
		if this.Tail == nil {
			this.Head = pe
			this.Tail = pe
		} else {
			pe.Prev = nil
			pe.Next = this.Head
			this.Head.Prev = pe
			this.Head = pe
		}
		if this.keysToEntries != nil {
			this.keysToEntries[*key] = pe
		}
		this.FieldCount++
	} else {
		pe.Value = value.Copy()
	}
}

// ----------------------------------------------------------------
func (this *Mlrmap) Get(key *string) *Mlrval {
	pe := this.findEntry(key)
	if pe == nil {
		return nil
	} else {
		return pe.Value
	}
	return nil
}

// ----------------------------------------------------------------
func (this *Mlrmap) GetKeys() []string {
	keys := make([]string, this.FieldCount)
	i := 0
	for pe := this.Head; pe != nil; pe = pe.Next {
		keys[i] = *pe.Key
		i++
	}
	return keys
}

// ----------------------------------------------------------------
// For '$[[[1]]]' etc. in the DSL.
//
// Notes:
// * This is a linear search.
// * Indices are 1-up not 0-up
// * Indices -n..-1 are aliases for 1..n. In particular, it will be faster to
//   get the -1st field than the nth.
// * Returns 0 on invalid index: 0, or < -n, or > n where n is the number of
//   fields.

func (this *Mlrmap) GetWithPositionalIndex(position int64) *Mlrval {
	mapEntry := this.findEntryByPositionalIndex(position)
	if mapEntry == nil {
		return nil
	}
	return mapEntry.Value
}

func (this *Mlrmap) GetWithMlrvalIndex(index *Mlrval) (*Mlrval, error) {
	if index.mvtype == MT_STRING {
		return this.Get(&index.printrep), nil
	} else if index.mvtype == MT_INT {
		s := index.String()
		return this.Get(&s), nil
	} else {
		return nil, errors.New(
			"Record/map indices must be string or int; got " + index.GetTypeName(),
		)
	}
}

// For '$[[1]]' etc. in the DSL.
//
// Notes:
// * This is a linear search.
// * Indices are 1-up not 0-up
// * Indices -n..-1 are aliases for 1..n. In particular, it will be faster to
//   get the -1st field than the nth.
// * Returns 0 on invalid index: 0, or < -n, or > n where n is the number of
//   fields.

func (this *Mlrmap) GetNameAtPositionalIndex(position int64) *string {
	mapEntry := this.findEntryByPositionalIndex(position)
	if mapEntry == nil {
		return nil
	}
	return mapEntry.Key
}

// ----------------------------------------------------------------
// TODO: put error-return into this API
func (this *Mlrmap) PutNameWithPositionalIndex(position int64, name *Mlrval) {
	mapEntry := this.findEntryByPositionalIndex(position)

	if mapEntry != nil {
		// TODO: rekey the hashmap
		if name.mvtype == MT_STRING {
			s := name.printrep
			mapEntry.Key = &s
		} else if name.mvtype == MT_INT {
			s := name.String()
			mapEntry.Key = &s
		} else {
			// TODO: return MlrvalFromError()
		}
	} else {
		// TODO: handle out-of-bounds accesses
		return
	}
}

// ----------------------------------------------------------------
// Copies the key and value (deep-copying in case the value is array/map).
// This is safe for DSL use. See also PutReference.

// TODO: put error-return into this API
func (this *Mlrmap) PutCopyWithPositionalIndex(position int64, value *Mlrval) {
	mapEntry := this.findEntryByPositionalIndex(position)

	if mapEntry != nil {
		mapEntry.Value = value.Copy()
	} else {
		return
	}
}

func (this *Mlrmap) RemoveWithPositionalIndex(position int64) {
	mapEntry := this.findEntryByPositionalIndex(position)
	if mapEntry != nil {
		this.unlink(mapEntry)
	}
}

// ----------------------------------------------------------------
func (this *Mlrmap) Equals(that *Mlrmap) bool {
	if this.FieldCount != this.FieldCount {
		return false
	}
	if !this.Contains(that) {
		return false
	}
	if !that.Contains(this) {
		return false
	}
	return true
}

// True if this contains that, i.e. if that is contained by this.
// * If any key of that is not a key of this, return false.
// * If any key of that has a value unequal to this' value at the same key, return false.
// * Else return true
func (this *Mlrmap) Contains(that *Mlrmap) bool {
	for pe := that.Head; pe != nil; pe = pe.Next {
		if !this.Has(pe.Key) {
			return false
		}
		thisval := this.Get(pe.Key)
		thatval := pe.Value
		meq := MlrvalEquals(thisval, thatval)
		eq, _ := meq.GetBoolValue()
		if !eq {
			return false
		}
	}
	return true
}

// ----------------------------------------------------------------
func (this *Mlrmap) Clear() {
	this.FieldCount = 0
	// Assuming everything unreferenced is getting GC'ed by the Go runtime
	this.Head = nil
	this.Tail = nil
}

// ----------------------------------------------------------------
func (this *Mlrmap) Copy() *Mlrmap {
	that := NewMlrmapMaybeHashed(this.isHashed())
	for pe := this.Head; pe != nil; pe = pe.Next {
		that.PutCopy(pe.Key, pe.Value)
	}
	return that
}

// ----------------------------------------------------------------
// Returns true if it was found and removed
func (this *Mlrmap) Remove(key *string) bool {
	pe := this.findEntry(key)
	if pe == nil {
		return false
	} else {
		this.unlink(pe)
		return true
	}
}

// ----------------------------------------------------------------
func (this *Mlrmap) MoveToHead(key *string) {
	pe := this.findEntry(key)
	if pe != nil {
		this.unlink(pe)
		this.linkAtHead(pe)
	}
}

// ----------------------------------------------------------------
func (this *Mlrmap) MoveToTail(key *string) {
	pe := this.findEntry(key)
	if pe != nil {
		this.unlink(pe)
		this.linkAtTail(pe)
	}
}

// ----------------------------------------------------------------
// E.g. '$name[1]["foo"] = "bar"' or '$*["foo"][1] = "bar"'
// In the former case the indices are ["name", 1, "foo"] and in the latter case
// the indices are ["foo", 1]. See also indexed-lvalues.md.
//
// This is a Mlrmap (from string to Mlrval) so we handle the first level of
// indexing here, then pass the remaining indices to the Mlrval at the desired
// slot.
func (this *Mlrmap) PutIndexed(indices []*Mlrval, rvalue *Mlrval) error {
	return putIndexedOnMap(this, indices, rvalue)
}

func (this *Mlrmap) UnsetIndexed(indices []*Mlrval) error {
	return unsetIndexedOnMap(this, indices)
}

// ----------------------------------------------------------------
func (this *Mlrmap) GetKeysJoined() string {
	var buffer bytes.Buffer
	i := 0
	for pe := this.Head; pe != nil; pe = pe.Next {
		if i > 0 {
			buffer.WriteString(",")
		}
		i++
		buffer.WriteString(*pe.Key)
	}
	return buffer.String()
}

// ----------------------------------------------------------------
// For group-by in several transformers.  If the record is 'a=x,b=y,c=3,d=4,e=5' and
// selectedFieldNames is 'a,b,c' then values are 'x,y,3'. This is returned as a
// comma-joined string.  The boolean ok is false if not all selected field
// names were present in the record.
//
// It's OK for the selected-field-namees list to be empty. This happens for
// transformers which support a -g option but are invoked without it (e.g. 'mlr tail
// -n 1' vs 'mlr tail -n 1 -g a,b,c'). In this case the return value is simply
// the empty string.
func (this *Mlrmap) GetSelectedValuesJoined(selectedFieldNames []string) (string, bool) {
	if len(selectedFieldNames) == 0 {
		// The fall-through is functionally correct, but this is quicker with
		// skipping setting up an empty bytes-buffer and stringifying it. The
		// non-grouped case is quite normal and is worth optimizing for.
		return "", true
	}

	var buffer bytes.Buffer
	for i, selectedFieldName := range selectedFieldNames {
		entry := this.findEntry(&selectedFieldName)
		if entry == nil {
			return "", false
		}
		if i > 0 {
			buffer.WriteString(",")
		}
		// This may be an array or map, or just a string/int/etc. Regardless we
		// stringify it.
		buffer.WriteString(entry.Value.String())
	}
	return buffer.String(), true
}

// As with GetSelectedValuesJoined but also returning the array of mlrvals.
// For sort.
func (this *Mlrmap) GetSelectedValuesAndJoined(selectedFieldNames []string) (
	string,
	[]Mlrval,
	bool,
) {
	mlrvals := make([]Mlrval, 0, len(selectedFieldNames))

	if len(selectedFieldNames) == 0 {
		// The fall-through is functionally correct, but this is quicker with
		// skipping setting up an empty bytes-buffer and stringifying it. The
		// non-grouped case is quite normal and is worth optimizing for.
		return "", mlrvals, true
	}

	var buffer bytes.Buffer
	for i, selectedFieldName := range selectedFieldNames {
		entry := this.findEntry(&selectedFieldName)
		if entry == nil {
			return "", mlrvals, false
		}
		if i > 0 {
			buffer.WriteString(",")
		}
		// This may be an array or map, or just a string/int/etc. Regardless we
		// stringify it.
		buffer.WriteString(entry.Value.String())
		mlrvals = append(mlrvals, *entry.Value.Copy())
	}
	return buffer.String(), mlrvals, true
}

// As above but only returns the array. Also, these are references, NOT copies.
// For step and join.
func (this *Mlrmap) ReferenceSelectedValues(selectedFieldNames []string) ([]*Mlrval, bool) {
	allFound := true
	mlrvals := make([]*Mlrval, 0, len(selectedFieldNames))

	for _, selectedFieldName := range selectedFieldNames {
		entry := this.findEntry(&selectedFieldName)
		if entry != nil {
			mlrvals = append(mlrvals, entry.Value)
		} else {
			mlrvals = append(mlrvals, nil)
			allFound = false
		}
	}
	return mlrvals, allFound
}

// Similar to the above but only checks availability. For join.
func (this *Mlrmap) HasSelectedKeys(selectedFieldNames []string) bool {
	for _, selectedFieldName := range selectedFieldNames {
		entry := this.findEntry(&selectedFieldName)
		if entry == nil {
			return false
		}
	}
	return true
}

// ----------------------------------------------------------------
func (this *Mlrmap) Rename(oldKey *string, newKey *string) bool {
	entry := this.findEntry(oldKey)
	if entry == nil {
		// Rename field from 'a' to 'b' where there is no 'a': no-op
		return false
	}

	existing := this.findEntry(newKey)
	if existing == nil {
		// Rename field from 'a' to 'b' where there is no 'b': simple update
		copy := *newKey
		entry.Key = &copy

		if this.keysToEntries != nil {
			delete(this.keysToEntries, *oldKey)
			this.keysToEntries[*newKey] = entry
		}
	} else {
		// Rename field from 'a' to 'b' where there are both 'a' and 'b':
		// remove old 'a' and put its value into the slot of 'b'.
		existing.Value = entry.Value
		delete(this.keysToEntries, *oldKey)
		this.unlink(entry)
	}

	return true
}

// ----------------------------------------------------------------
func (this *Mlrmap) Label(newNames []string) {
	that := NewMlrmapAsRecord()

	i := 0
	numNewNames := len(newNames)
	for {
		if i >= numNewNames {
			break
		}
		pe := this.pop()
		if pe == nil {
			break
		}
		// Old record will be GC'ed: just move pointers
		that.PutReference(&newNames[i], pe.Value)
		i++
	}

	for {
		pe := this.pop()
		if pe == nil {
			break
		}
		that.PutReference(pe.Key, pe.Value)
	}

	*this = *that
}

// ----------------------------------------------------------------
func (this *Mlrmap) SortByKey() {
	keys := this.GetKeys()

	lib.SortStrings(keys)

	that := NewMlrmapAsRecord()

	for _, key := range keys {
		// Old record will be GC'ed: just move pointers
		that.PutReference(&key, this.Get(&key))
	}

	*this = *that
}

// ================================================================
// PRIVATE METHODS

// ----------------------------------------------------------------
func (this *Mlrmap) findEntry(key *string) *mlrmapEntry {
	if this.keysToEntries != nil {
		return this.keysToEntries[*key]
	} else {
		for pe := this.Head; pe != nil; pe = pe.Next {
			if *pe.Key == *key {
				return pe
			}
		}
		return nil
	}
}

// ----------------------------------------------------------------
// For '$[1]' etc. in the DSL.
//
// Notes:
// * This is a linear search.
// * Indices are 1-up not 0-up
// * Indices -n..-1 are aliases for 1..n. In particular, it will be faster to
//   get the -1st field than the nth.
// * Returns 0 on invalid index: 0, or < -n, or > n where n is the number of
//   fields.
func (this *Mlrmap) findEntryByPositionalIndex(position int64) *mlrmapEntry {
	if position > this.FieldCount || position < -this.FieldCount || position == 0 {
		return nil
	}
	if position > 0 {
		var i int64 = 1
		for pe := this.Head; pe != nil; pe = pe.Next {
			if i == position {
				return pe
			}
			i++
		}
		lib.InternalCodingErrorIf(true)
	} else {
		var i int64 = -1
		for pe := this.Tail; pe != nil; pe = pe.Prev {
			if i == position {
				return pe
			}
			i--
		}
		lib.InternalCodingErrorIf(true)
	}
	lib.InternalCodingErrorIf(true)
	return nil
}

// ----------------------------------------------------------------
func (this *Mlrmap) unlink(pe *mlrmapEntry) {
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
		delete(this.keysToEntries, *pe.Key)
	}
	this.FieldCount--
}

// ----------------------------------------------------------------
// Does not check for duplicate keys
func (this *Mlrmap) linkAtHead(pe *mlrmapEntry) {
	if this.Head == nil {
		pe.Prev = nil
		pe.Next = nil
		this.Head = pe
		this.Tail = pe
	} else {
		pe.Prev = nil
		pe.Next = this.Head
		this.Head.Prev = pe
		this.Head = pe
	}
	this.FieldCount++
}

// Does not check for duplicate keys
func (this *Mlrmap) linkAtTail(pe *mlrmapEntry) {
	if this.Head == nil {
		pe.Prev = nil
		pe.Next = nil
		this.Head = pe
		this.Tail = pe
	} else {
		pe.Prev = this.Tail
		pe.Next = nil
		this.Tail.Next = pe
		this.Tail = pe
	}
	this.FieldCount++
}

// ----------------------------------------------------------------
func (this *Mlrmap) pop() *mlrmapEntry {
	if this.Head == nil {
		return nil
	} else {
		pe := this.Head
		this.unlink(pe)
		return pe
	}
}
