package mlrval

import "strconv"

// arenaChunkFields is the slab size (in fields) used when a RecordArena needs
// to grow. It is large enough that a typical record batch draws from one or two
// slabs, amortizing allocation overhead across thousands of fields.
const arenaChunkFields = 4096

// RecordArena is a batch (slab) allocator for record fields. A record reader
// allocates one arena per batch of records and draws each field's MlrmapEntry
// and Mlrval from contiguous slabs, turning two heap allocations per field into
// roughly two allocations per slab.
//
// Lifetime: a slab stays alive as long as any entry or value drawn from it is
// reachable. For streaming verbs the whole batch is processed and released
// together, so its slabs are freed as units. For accumulating verbs (e.g. sort)
// the slabs are retained for the duration -- the same bytes as before, but as a
// few large objects rather than millions of tiny ones, which also lowers
// resident-set size via reduced fragmentation.
//
// The arena grows on demand (allocating a fresh slab when the current one is
// exhausted), so the size hint passed to NewRecordArena need not be exact: it
// only sizes the first slab.
type RecordArena struct {
	entries []MlrmapEntry
	values  []Mlrval
	ei      int
	vi      int
	chunk   int

	// records is a slab of Mlrmap structs vended by NewRecord, so the record
	// container itself is batch-allocated alongside its fields.
	records []Mlrmap
	ri      int
}

// arenaChunkRecords is the slab size (in records) for NewRecord. It tracks the
// nominal records-per-batch so a batch typically draws from one slab.
const arenaChunkRecords = 512

// NewRecordArena returns an arena whose first slabs are sized to nfieldsHint
// (clamped to a sane range). Subsequent slabs, if needed, use arenaChunkFields.
func NewRecordArena(nfieldsHint int) *RecordArena {
	chunk := nfieldsHint
	if chunk < 1 {
		chunk = 1
	}
	if chunk > arenaChunkFields {
		chunk = arenaChunkFields
	}
	return &RecordArena{chunk: chunk}
}

// NewRecord vends a fresh empty record (lazily-hashed, like NewMlrmapAsRecord)
// from the arena's record slab, batch-allocating the Mlrmap struct itself. It
// respects the global hashRecords setting: with --no-hash-records it returns an
// unhashed map, matching NewMlrmapAsRecord.
func (a *RecordArena) NewRecord() *Mlrmap {
	if !hashRecords {
		return newMlrmapUnhashed()
	}
	if a.ri >= len(a.records) {
		a.records = make([]Mlrmap, arenaChunkRecords)
		a.ri = 0
	}
	m := &a.records[a.ri]
	a.ri++
	m.autoHash = true
	return m
}

// PutDeferred appends a field to mlrmap, drawing the entry and its deferred-type
// value from the arena slabs. It mirrors PutReferenceMaybeDedupe's semantics for
// duplicate keys. The value is built from the raw input string with type
// inference deferred (MT_PENDING), exactly as FromDeferredType does.
func (a *RecordArena) PutDeferred(mlrmap *Mlrmap, key string, input string, dedupe bool) {
	pe := mlrmap.findEntry(key)
	if pe == nil {
		mlrmap.linkNewEntry(a.newEntry(key, input))
		return
	}
	if !dedupe {
		pe.Value = a.newValue(input)
		return
	}
	for i := 2; ; i++ {
		newKey := key + "_" + strconv.Itoa(i)
		if mlrmap.findEntry(newKey) == nil {
			mlrmap.linkNewEntry(a.newEntry(newKey, input))
			return
		}
	}
}

func (a *RecordArena) newValue(input string) *Mlrval {
	if a.vi >= len(a.values) {
		a.values = make([]Mlrval, a.chunk)
		a.vi = 0
	}
	v := &a.values[a.vi]
	a.vi++
	v.mvtype = MT_PENDING
	v.printrep = input
	v.printrepValid = true
	return v
}

func (a *RecordArena) newEntry(key string, input string) *MlrmapEntry {
	if a.ei >= len(a.entries) {
		a.entries = make([]MlrmapEntry, a.chunk)
		a.ei = 0
	}
	e := &a.entries[a.ei]
	a.ei++
	e.Key = key
	e.Value = a.newValue(input)
	return e
}
