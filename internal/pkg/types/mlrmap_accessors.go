package types

import (
	"bytes"
	"errors"

	"mlr/internal/pkg/lib"
)

// IsEmpty determines if an map is empty.
func (mlrmap *Mlrmap) IsEmpty() bool {
	return mlrmap.Head == nil
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) Has(key string) bool {
	return mlrmap.findEntry(key) != nil
}

// PutCpoy copies the key and value (deep-copying in case the value is array/map).
// This is safe for DSL use. See also PutReference.
func (mlrmap *Mlrmap) PutCopy(key string, value *Mlrval) {
	pe := mlrmap.findEntry(key)
	if pe == nil {
		pe = newMlrmapEntry(key, value.Copy())
		if mlrmap.Head == nil {
			mlrmap.Head = pe
			mlrmap.Tail = pe
		} else {
			pe.Prev = mlrmap.Tail
			pe.Next = nil
			mlrmap.Tail.Next = pe
			mlrmap.Tail = pe
		}
		if mlrmap.keysToEntries != nil {
			mlrmap.keysToEntries[key] = pe
		}
		mlrmap.FieldCount++
	} else {
		pe.Value = value.Copy()
	}
}

// PutReference copies the key but not the value. This is not safe for DSL use,
// where we could create undesired references between different objects.  Only
// intended to be used at callsites which allocate a mlrval solely for the
// purpose of putting into a map, e.g. input-record readers.
func (mlrmap *Mlrmap) PutReference(key string, value *Mlrval) {
	pe := mlrmap.findEntry(key)
	if pe == nil {
		pe = newMlrmapEntry(key, value)
		if mlrmap.Head == nil {
			mlrmap.Head = pe
			mlrmap.Tail = pe
		} else {
			pe.Prev = mlrmap.Tail
			pe.Next = nil
			mlrmap.Tail.Next = pe
			mlrmap.Tail = pe
		}
		if mlrmap.keysToEntries != nil {
			mlrmap.keysToEntries[key] = pe
		}
		mlrmap.FieldCount++
	} else {
		pe.Value = value
	}
}

// TODO: COMMENT
func (mlrmap *Mlrmap) PutReferenceAfter(
	pe *MlrmapEntry,
	key string,
	value *Mlrval,
) *MlrmapEntry {
	if pe == nil || pe.Next == nil {
		// New entry is supposed to go 'after' old, but there is no such old
		// entry.  Or, the old entry exists and is at the tail. In either case,
		// add the new entry at the end of the record (new tail).

		// TODO: make a helper method for code-dedupe
		pf := newMlrmapEntry(key, value)
		if mlrmap.Head == nil { // First entry into empty map
			mlrmap.Head = pf
			mlrmap.Tail = pf
		} else {
			// Before: ... pc pd
			// After:  ... pc pd pf
			pf.Prev = mlrmap.Tail
			pf.Next = nil
			mlrmap.Tail.Next = pf
			mlrmap.Tail = pf
		}

		if mlrmap.keysToEntries != nil {
			mlrmap.keysToEntries[key] = pf
		}
		mlrmap.FieldCount++
		return pf

	} else {
		// Before: ... pe pg ...
		// After:  ... pe pf pg ...
		//
		// New entry is neither the new head (pe != nil) nor the new tail
		// (pe.Next != nil, otherwise we'd be in the if-branch above).

		pf := newMlrmapEntry(key, value)
		pg := pe.Next

		pe.Next = pf
		pf.Next = pg
		pf.Prev = pe
		if pg != nil {
			pg.Prev = pf
		}

		if mlrmap.keysToEntries != nil {
			mlrmap.keysToEntries[key] = pf
		}
		mlrmap.FieldCount++
		return pf
	}
}

func (mlrmap *Mlrmap) PutCopyWithMlrvalIndex(key *Mlrval, value *Mlrval) error {
	if key.mvtype == MT_STRING {
		mlrmap.PutCopy(key.printrep, value)
		return nil
	} else if key.mvtype == MT_INT {
		mlrmap.PutCopy(key.String(), value)
		return nil
	} else {
		return errors.New(
			"mlr: record/map indices must be string, int, or array thereof; got " + key.GetTypeName(),
		)
	}
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) PrependCopy(key string, value *Mlrval) {
	mlrmap.PrependReference(key, value.Copy())
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) PrependReference(key string, value *Mlrval) {
	pe := mlrmap.findEntry(key)
	if pe == nil {
		pe = newMlrmapEntry(key, value)
		if mlrmap.Tail == nil {
			mlrmap.Head = pe
			mlrmap.Tail = pe
		} else {
			pe.Prev = nil
			pe.Next = mlrmap.Head
			mlrmap.Head.Prev = pe
			mlrmap.Head = pe
		}
		if mlrmap.keysToEntries != nil {
			mlrmap.keysToEntries[key] = pe
		}
		mlrmap.FieldCount++
	} else {
		pe.Value = value
	}
}

// ----------------------------------------------------------------
// Merges that into mlrmap.
func (mlrmap *Mlrmap) Merge(other *Mlrmap) {
	for pe := other.Head; pe != nil; pe = pe.Next {
		mlrmap.PutCopy(pe.Key, pe.Value)
	}
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) Get(key string) *Mlrval {
	pe := mlrmap.findEntry(key)
	if pe == nil {
		return nil
	} else {
		return pe.Value
	}
	return nil
}

// ----------------------------------------------------------------
// Exposed for the 'nest' verb
func (mlrmap *Mlrmap) GetEntry(key string) *MlrmapEntry {
	return mlrmap.findEntry(key)
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) GetKeys() []string {
	keys := make([]string, mlrmap.FieldCount)
	i := 0
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		keys[i] = pe.Key
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

func (mlrmap *Mlrmap) GetWithPositionalIndex(position int) *Mlrval {
	mapEntry := mlrmap.findEntryByPositionalIndex(position)
	if mapEntry == nil {
		return nil
	}
	return mapEntry.Value
}

func (mlrmap *Mlrmap) GetWithMlrvalIndex(index *Mlrval) (*Mlrval, error) {
	if index.mvtype == MT_ARRAY {
		return mlrmap.getWithMlrvalArrayIndex(index)
	} else {
		return mlrmap.getWithMlrvalSingleIndex(index)
	}
}

// This lets the user do '$y = $x[ ["a", "b", "c"] ]' in lieu of
// '$y = $x["a"]["b"]["c"]'.
func (mlrmap *Mlrmap) getWithMlrvalArrayIndex(index *Mlrval) (*Mlrval, error) {
	current := mlrmap
	var retval *Mlrval = nil
	lib.InternalCodingErrorIf(index.mvtype != MT_ARRAY)
	array := index.arrayval
	n := len(array)
	for i, piece := range array {
		next, err := current.GetWithMlrvalIndex(&piece)
		if err != nil {
			return nil, err
		}
		if i < n-1 {
			if !next.IsMap() {
				return nil, errors.New(
					"mlr: cannot multi-index non-map.",
				)
			}
			current = next.mapval
		} else {
			retval = next.Copy()
		}
	}
	lib.InternalCodingErrorIf(retval == nil)
	return retval, nil
}

func (mlrmap *Mlrmap) getWithMlrvalSingleIndex(index *Mlrval) (*Mlrval, error) {
	if index.mvtype == MT_STRING {
		return mlrmap.Get(index.printrep), nil
	} else if index.mvtype == MT_INT {
		return mlrmap.Get(index.String()), nil
	} else {
		return nil, errors.New(
			"Record/map indices must be string, int, or array thereof; got " + index.GetTypeName(),
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

func (mlrmap *Mlrmap) GetNameAtPositionalIndex(position int) (string, bool) {
	mapEntry := mlrmap.findEntryByPositionalIndex(position)
	if mapEntry == nil {
		return "", false
	}
	return mapEntry.Key, true
}

// ----------------------------------------------------------------
// TODO: put error-return into this API
func (mlrmap *Mlrmap) PutNameWithPositionalIndex(position int, name *Mlrval) {
	positionalEntry := mlrmap.findEntryByPositionalIndex(position)

	if positionalEntry == nil {
		// TODO: handle out-of-bounds accesses
		return
	}

	// TODO: rekey the hashmap
	s := ""
	if name.mvtype == MT_STRING {
		s = name.printrep
	} else if name.mvtype == MT_INT {
		s = name.String()
	} else {
		// TODO: return MlrvalFromError()
		return
	}

	// E.g. there are fields named 'a' and 'b', as positions 1 and 2,
	// and the user does '$[[1]] = $[[2]]'. Then there would be two b's.
	mapEntry := mlrmap.findEntry(s)
	if mapEntry != nil && mapEntry != positionalEntry {
		mlrmap.Unlink(mapEntry)
	}

	lib.InternalCodingErrorIf(s == "")
	positionalEntry.Key = s
}

// ----------------------------------------------------------------
// Copies the key and value (deep-copying in case the value is array/map).
// This is safe for DSL use. See also PutReference.

// TODO: put error-return into this API
func (mlrmap *Mlrmap) PutCopyWithPositionalIndex(position int, value *Mlrval) {
	mapEntry := mlrmap.findEntryByPositionalIndex(position)

	if mapEntry != nil {
		mapEntry.Value = value.Copy()
	} else {
		return
	}
}

func (mlrmap *Mlrmap) RemoveWithPositionalIndex(position int) {
	mapEntry := mlrmap.findEntryByPositionalIndex(position)
	if mapEntry != nil {
		mlrmap.Unlink(mapEntry)
	}
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) Equals(other *Mlrmap) bool {
	if mlrmap.FieldCount != mlrmap.FieldCount {
		return false
	}
	if !mlrmap.Contains(other) {
		return false
	}
	if !other.Contains(mlrmap) {
		return false
	}
	return true
}

// True if this contains other, i.e. if other is contained by mlrmap.
// * If any key of other is not a key of this, return false.
// * If any key of other has a value unequal to this' value at the same key, return false.
// * Else return true
func (mlrmap *Mlrmap) Contains(other *Mlrmap) bool {
	for pe := other.Head; pe != nil; pe = pe.Next {
		if !mlrmap.Has(pe.Key) {
			return false
		}
		thisval := mlrmap.Get(pe.Key)
		thatval := pe.Value
		meq := BIF_equals(thisval, thatval)
		eq, ok := meq.GetBoolValue()
		lib.InternalCodingErrorIf(!ok)
		if !eq {
			return false
		}
	}
	return true
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) Clear() {
	mlrmap.FieldCount = 0
	// Assuming everything unreferenced is getting GC'ed by the Go runtime
	mlrmap.Head = nil
	mlrmap.Tail = nil
	if mlrmap.keysToEntries != nil {
		mlrmap.keysToEntries = make(map[string]*MlrmapEntry)
	}
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) Copy() *Mlrmap {
	other := NewMlrmapMaybeHashed(mlrmap.isHashed())
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		other.PutCopy(pe.Key, pe.Value)
	}
	return other
}

// ----------------------------------------------------------------
// Returns true if it was found and removed
func (mlrmap *Mlrmap) Remove(key string) bool {
	pe := mlrmap.findEntry(key)
	if pe == nil {
		return false
	} else {
		mlrmap.Unlink(pe)
		return true
	}
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) MoveToHead(key string) {
	pe := mlrmap.findEntry(key)
	if pe != nil {
		mlrmap.Unlink(pe)
		mlrmap.linkAtHead(pe)
	}
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) MoveToTail(key string) {
	pe := mlrmap.findEntry(key)
	if pe != nil {
		mlrmap.Unlink(pe)
		mlrmap.linkAtTail(pe)
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
func (mlrmap *Mlrmap) PutIndexed(indices []*Mlrval, rvalue *Mlrval) error {
	return putIndexedOnMap(mlrmap, indices, rvalue)
}

func (mlrmap *Mlrmap) RemoveIndexed(indices []*Mlrval) error {
	return removeIndexedOnMap(mlrmap, indices)
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) GetKeysJoined() string {
	var buffer bytes.Buffer
	i := 0
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		if i > 0 {
			buffer.WriteString(",")
		}
		i++
		buffer.WriteString(pe.Key)
	}
	return buffer.String()
}

// For mlr reshape
func (mlrmap *Mlrmap) GetValuesJoined() string {
	var buffer bytes.Buffer
	i := 0
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		if i > 0 {
			buffer.WriteString(",")
		}
		i++
		buffer.WriteString(pe.Value.String())
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
func (mlrmap *Mlrmap) GetSelectedValuesJoined(selectedFieldNames []string) (string, bool) {
	if len(selectedFieldNames) == 0 {
		// The fall-through is functionally correct, but this is quicker with
		// skipping setting up an empty bytes-buffer and stringifying it. The
		// non-grouped case is quite normal and is worth optimizing for.
		return "", true
	}

	var buffer bytes.Buffer
	for i, selectedFieldName := range selectedFieldNames {
		entry := mlrmap.findEntry(selectedFieldName)
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
// TODO: put 'Copy' into the method name
func (mlrmap *Mlrmap) GetSelectedValuesAndJoined(selectedFieldNames []string) (
	string,
	[]*Mlrval,
	bool,
) {
	mlrvals := make([]*Mlrval, 0, len(selectedFieldNames))

	if len(selectedFieldNames) == 0 {
		// The fall-through is functionally correct, but this is quicker with
		// skipping setting up an empty bytes-buffer and stringifying it. The
		// non-grouped case is quite normal and is worth optimizing for.
		return "", mlrvals, true
	}

	var buffer bytes.Buffer
	for i, selectedFieldName := range selectedFieldNames {
		entry := mlrmap.findEntry(selectedFieldName)
		if entry == nil {
			return "", mlrvals, false
		}
		if i > 0 {
			buffer.WriteString(",")
		}
		// This may be an array or map, or just a string/int/etc. Regardless we
		// stringify it.
		buffer.WriteString(entry.Value.String())
		mlrvals = append(mlrvals, entry.Value.Copy())
	}
	return buffer.String(), mlrvals, true
}

// As above but only returns the array. Also, these are references, NOT copies.
// For step and join.
func (mlrmap *Mlrmap) ReferenceSelectedValues(selectedFieldNames []string) ([]*Mlrval, bool) {
	allFound := true
	mlrvals := make([]*Mlrval, 0, len(selectedFieldNames))

	for _, selectedFieldName := range selectedFieldNames {
		entry := mlrmap.findEntry(selectedFieldName)
		if entry != nil {
			mlrvals = append(mlrvals, entry.Value)
		} else {
			mlrvals = append(mlrvals, nil)
			allFound = false
		}
	}
	return mlrvals, allFound
}

// TODO: rename to CopySelectedValues
// As previous but with copying. For stats1.
func (mlrmap *Mlrmap) GetSelectedValues(selectedFieldNames []string) ([]*Mlrval, bool) {
	allFound := true
	mlrvals := make([]*Mlrval, 0, len(selectedFieldNames))

	for _, selectedFieldName := range selectedFieldNames {
		entry := mlrmap.findEntry(selectedFieldName)
		if entry != nil {
			mlrvals = append(mlrvals, entry.Value.Copy())
		} else {
			mlrvals = append(mlrvals, nil)
			allFound = false
		}
	}
	return mlrvals, allFound
}

// Similar to the above but only checks availability. For join.
func (mlrmap *Mlrmap) HasSelectedKeys(selectedFieldNames []string) bool {
	for _, selectedFieldName := range selectedFieldNames {
		entry := mlrmap.findEntry(selectedFieldName)
		if entry == nil {
			return false
		}
	}
	return true
}

// ----------------------------------------------------------------
// For mlr nest implode across records.
func (mlrmap *Mlrmap) GetKeysJoinedExcept(px *MlrmapEntry) string {
	var buffer bytes.Buffer
	i := 0
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		if pe == px {
			continue
		}
		if i > 0 {
			buffer.WriteString(",")
		}
		i++
		buffer.WriteString(pe.Key)
	}
	return buffer.String()
}

// For mlr nest implode across records.
func (mlrmap *Mlrmap) GetValuesJoinedExcept(px *MlrmapEntry) string {
	var buffer bytes.Buffer
	i := 0
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		if pe == px {
			continue
		}
		if i > 0 {
			buffer.WriteString(",")
		}
		i++
		// This may be an array or map, or just a string/int/etc. Regardless we
		// stringify it.
		buffer.WriteString(pe.Value.String())
	}
	return buffer.String()
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) Rename(oldKey string, newKey string) bool {
	entry := mlrmap.findEntry(oldKey)
	if entry == nil {
		// Rename field from 'a' to 'b' where there is no 'a': no-op
		return false
	}

	existing := mlrmap.findEntry(newKey)
	if existing == nil {
		// Rename field from 'a' to 'b' where there is no 'b': simple update
		entry.Key = newKey

		if mlrmap.keysToEntries != nil {
			delete(mlrmap.keysToEntries, oldKey)
			mlrmap.keysToEntries[newKey] = entry
		}
	} else {
		// Rename field from 'a' to 'b' where there are both 'a' and 'b':
		// remove old 'a' and put its value into the slot of 'b'.
		existing.Value = entry.Value
		if mlrmap.keysToEntries != nil {
			delete(mlrmap.keysToEntries, oldKey)
		}
		mlrmap.Unlink(entry)
	}

	return true
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) Label(newNames []string) {
	other := NewMlrmapAsRecord()

	i := 0
	numNewNames := len(newNames)
	for {
		if i >= numNewNames {
			break
		}
		pe := mlrmap.pop()
		if pe == nil {
			break
		}
		// Old record will be GC'ed: just move pointers
		other.PutReference(newNames[i], pe.Value)
		i++
	}

	for {
		pe := mlrmap.pop()
		if pe == nil {
			break
		}
		// Example:
		// * Input record has keys a,b,i,x,y
		// * Requested labeling is d,x,f
		// * The first three records a,b,i should be renamed to d,x,f
		// * The old x needs to disappear (for key-uniqueness)
		// * The y field is carried through
		if other.Has(pe.Key) {
			continue
		}
		other.PutReference(pe.Key, pe.Value)
	}

	*mlrmap = *other
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) SortByKey() {
	keys := mlrmap.GetKeys()

	lib.SortStrings(keys)

	other := NewMlrmapAsRecord()

	for _, key := range keys {
		// Old record will be GC'ed: just move pointers
		other.PutReference(key, mlrmap.Get(key))
	}

	*mlrmap = *other
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) SortByKeyRecursively() {
	keys := mlrmap.GetKeys()

	lib.SortStrings(keys)

	other := NewMlrmapAsRecord()

	for _, key := range keys {
		// Old record will be GC'ed: just move pointers
		value := mlrmap.Get(key)
		if value.IsMap() {
			value.mapval.SortByKeyRecursively()
		}
		other.PutReference(key, value)
	}

	*mlrmap = *other
}

// ----------------------------------------------------------------
// Only checks to see if the first entry is a map. For emit/emitp.
func (mlrmap *Mlrmap) IsNested() bool {
	if mlrmap.Head == nil {
		return false
	} else if mlrmap.Head.Value.GetMap() == nil {
		return false
	} else {
		return true
	}
}

// ================================================================
// PRIVATE METHODS

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) findEntry(key string) *MlrmapEntry {
	if mlrmap.keysToEntries != nil {
		return mlrmap.keysToEntries[key]
	} else {
		for pe := mlrmap.Head; pe != nil; pe = pe.Next {
			if pe.Key == key {
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
func (mlrmap *Mlrmap) findEntryByPositionalIndex(position int) *MlrmapEntry {
	if position > mlrmap.FieldCount || position < -mlrmap.FieldCount || position == 0 {
		return nil
	}
	if position > 0 {
		var i int = 1
		for pe := mlrmap.Head; pe != nil; pe = pe.Next {
			if i == position {
				return pe
			}
			i++
		}
		lib.InternalCodingErrorIf(true)
	} else {
		var i int = -1
		for pe := mlrmap.Tail; pe != nil; pe = pe.Prev {
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
func (mlrmap *Mlrmap) Unlink(pe *MlrmapEntry) {
	if pe == mlrmap.Head {
		if pe == mlrmap.Tail {
			mlrmap.Head = nil
			mlrmap.Tail = nil
		} else {
			mlrmap.Head = pe.Next
			pe.Next.Prev = nil
		}
	} else {
		pe.Prev.Next = pe.Next
		if pe == mlrmap.Tail {
			mlrmap.Tail = pe.Prev
		} else {
			pe.Next.Prev = pe.Prev
		}
	}
	if mlrmap.keysToEntries != nil {
		delete(mlrmap.keysToEntries, pe.Key)
	}
	mlrmap.FieldCount--
}

// ----------------------------------------------------------------
// Does not check for duplicate keys
func (mlrmap *Mlrmap) linkAtHead(pe *MlrmapEntry) {
	if mlrmap.Head == nil {
		pe.Prev = nil
		pe.Next = nil
		mlrmap.Head = pe
		mlrmap.Tail = pe
	} else {
		pe.Prev = nil
		pe.Next = mlrmap.Head
		mlrmap.Head.Prev = pe
		mlrmap.Head = pe
	}
	if mlrmap.keysToEntries != nil {
		mlrmap.keysToEntries[pe.Key] = pe
	}
	mlrmap.FieldCount++
}

// Does not check for duplicate keys
func (mlrmap *Mlrmap) linkAtTail(pe *MlrmapEntry) {
	if mlrmap.Head == nil {
		pe.Prev = nil
		pe.Next = nil
		mlrmap.Head = pe
		mlrmap.Tail = pe
	} else {
		pe.Prev = mlrmap.Tail
		pe.Next = nil
		mlrmap.Tail.Next = pe
		mlrmap.Tail = pe
	}
	if mlrmap.keysToEntries != nil {
		mlrmap.keysToEntries[pe.Key] = pe
	}
	mlrmap.FieldCount++
}

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) pop() *MlrmapEntry {
	if mlrmap.Head == nil {
		return nil
	} else {
		pe := mlrmap.Head
		mlrmap.Unlink(pe)
		return pe
	}
}

// ----------------------------------------------------------------

// ToPairsArray is used for sorting maps by key/value/etc, e.g. the sortmf DSL function.
func (mlrmap *Mlrmap) ToPairsArray() []MlrmapPair {
	pairsArray := make([]MlrmapPair, mlrmap.FieldCount)
	i := 0
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		pairsArray[i].Key = pe.Key
		pairsArray[i].Value = pe.Value.Copy()
		i++
	}

	return pairsArray
}

// MlrmapFromPairsArray is used for sorting maps by key/value/etc, e.g. the sortmf DSL function.
func MlrmapFromPairsArray(pairsArray []MlrmapPair) *Mlrmap {
	mlrmap := NewMlrmap()
	for i := range pairsArray {
		mlrmap.PutCopy(pairsArray[i].Key, pairsArray[i].Value)
	}

	return mlrmap
}

// ----------------------------------------------------------------

// GetFirstPair returns the first key-value pair as its own map.  If the map is
// empty (i.e. there is no first pair) it returns nil.
func (mlrmap *Mlrmap) GetFirstPair() *Mlrmap {
	if mlrmap.Head == nil {
		return nil
	}
	pair := NewMlrmap()
	pair.PutCopy(mlrmap.Head.Key, mlrmap.Head.Value)
	return pair
}

func (mlrmap *Mlrmap) IsSinglePair() bool {
	return mlrmap.FieldCount == 1
}
