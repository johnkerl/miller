// ================================================================
// ORDERED MAP FROM STRING TO MLRVAL
//
// This is an implementation of insertion-ordered key-value pairs for Miller's
// fundamental record data structure. It's also an ordered-map data structure,
// suitable for Miller JSON decode/encode.
//
// ----------------------------------------------------------------
// DESIGN
//
// * It keeps a doubly-linked list of key-value pairs.
//
// * With hash-records set to false, no hash functions are computed when the map
//   is written to or read from.
//
// * Gets are implemented by sequential scan through the list: given a key,
//   the key-value pairs are scanned through until a match is (or is not) found.
//
// * Performance improvement of 25% percent over hash-mapping from key to entry
//   was found in the Go implementation. Test data was million-line CSV and
//   DKVP, with a dozen columns or so.
//
// * However, with higher column-count (see https://github.com/johnkerl/miller/issues/1506
//   and https://github.com/johnkerl/miller/pull/1507), non-hashing becomes
//   a substantial penalty.
//
// Note however that an auxiliary constructor is provided which does use
// a key-to-entry hashmap in place of linear search for get/put/has/delete.
// This may be useful in certain contexts, even though it's not the default
// chosen for stream-records.
//
// ----------------------------------------------------------------
// MOTIVATION
//
// * The use case for records in Miller is that *all* fields are read from
//   strings & written to strings (split/join), while only *some* fields are
//   operated on.
//
// * Meanwhile there are few repeated accesses to a given record: the
//   access-to-construct ratio is quite low for Miller data records.  Miller
//   instantiates thousands, millions, billions of records (depending on the
//   input data) but accesses each record only once per transforming operation.
//   (This is in contrast to accumulator hashmaps which are repeatedly accessed
//   during a stats run.)
//
// * The hashed impl computes hashsums for *all* fields whether operated on or not,
//   for the benefit of the *few* fields looked up during the transforming operation.
//
// * The hashless impl only keeps string pointers.  Lookups are done at runtime
//   doing prefix search on the key names. Assuming field names are distinct,
//   this is just a few char-ptr accesses which (in experiments) turn out to
//   offer about a 10-15% performance improvement.
//
// * Added benefit: the field-rename operation (preserving field order) becomes
//   trivial.
// ================================================================

package mlrval

// As noted above, hashing has a minor penalty for low column count: computing
// hashmaps takes more time than is saved later on. But for higher column-count,
// non-hashing has a huge penalty.  Therefore we default to on. And users can
// use `mlr --no-hash-records` or `mlr --hash-records` to flip the behavior.
var hashRecords = true

func HashRecords(onOff bool) {
	hashRecords = onOff
}

// ----------------------------------------------------------------
type Mlrmap struct {
	FieldCount int64
	Head       *MlrmapEntry
	Tail       *MlrmapEntry

	// This can be nil if hashRecords is off.
	keysToEntries map[string]*MlrmapEntry
}

type MlrmapEntry struct {
	Key   string
	Value *Mlrval
	Prev  *MlrmapEntry
	Next  *MlrmapEntry
}

// MlrmapEntryForArray is for use by sorting routines where the Prev/Next pointers
// are irrelevant as well as ephemeral
type MlrmapEntryForArray struct {
	Key   string
	Value *Mlrval
}

// Only used for sorting, map-to-pairs-array and pairs-array-to-map contexts.
type MlrmapPair struct {
	Key   string
	Value *Mlrval
}

// ----------------------------------------------------------------
func NewMlrmapAsRecord() *Mlrmap {
	if hashRecords {
		return newMlrmapHashed()
	} else {
		return newMlrmapUnhashed()
	}
}
func NewMlrmap() *Mlrmap {
	return newMlrmapHashed()
}

// Faster on record-stream data as noted above.
func newMlrmapUnhashed() *Mlrmap {
	return &Mlrmap{
		FieldCount:    0,
		Head:          nil,
		Tail:          nil,
		keysToEntries: nil,
	}
}

// Intended for use in DSL expressions wherein the access-to-construct ratio
// might be higher (although this needs profiling over a variety of use-cases).
func newMlrmapHashed() *Mlrmap {
	return &Mlrmap{
		FieldCount:    0,
		Head:          nil,
		Tail:          nil,
		keysToEntries: make(map[string]*MlrmapEntry),
	}
}

func NewMlrmapMaybeHashed(wantHashing bool) *Mlrmap {
	if wantHashing {
		return newMlrmapHashed()
	} else {
		return newMlrmapUnhashed()
	}
}

func (mlrmap *Mlrmap) isHashed() bool {
	return mlrmap.keysToEntries != nil
}

// ----------------------------------------------------------------
// Value-copy is up to the caller -- PutReference and PutCopy
// are in the public Mlrmap API.
func newMlrmapEntry(key string, value *Mlrval) *MlrmapEntry {
	return &MlrmapEntry{
		key,
		value,
		nil,
		nil,
	}
}
