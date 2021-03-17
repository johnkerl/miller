// ================================================================
// Data structure for mlr top: just a decorated array.
// ================================================================

package utils

import (
	"miller/src/types"
)

// ----------------------------------------------------------------
type TopKeeper struct {
	TopValues             []*types.Mlrval
	TopRecordsAndContexts []*types.RecordAndContext
	size                  int
	capacity              int
	bsearchFunc           types.BsearchMlrvalArrayFunc
}

// ----------------------------------------------------------------
func NewTopKeeper(capacity int, doMax bool) *TopKeeper {
	this := &TopKeeper{
		TopValues:             make([]*types.Mlrval, capacity),
		TopRecordsAndContexts: make([]*types.RecordAndContext, capacity),
		size:                  0,
		capacity:              capacity,
	}
	if doMax {
		this.bsearchFunc = types.BsearchMlrvalArrayForDescendingInsert
	} else {
		this.bsearchFunc = types.BsearchMlrvalArrayForAscendingInsert
	}
	return this
}

func (this *TopKeeper) GetSize() int {
	return this.size
}

// ----------------------------------------------------------------
// Cases:
// 1. array size <  capacity
//    * find destidx
//    * if destidx == size
//        put it there
//      else
//        shift down & insert
//      increment size
//
// 2. array size == capacity
//    * find destidx
//    * if destidx == size
//        discard
//      else
//        shift down & insert
//
// capacity = 10, size = 6, destidx = 3     capacity = 10, size = 10, destidx = 3
// [0 #]   [0 #]                            [0 #]   [0 #]
// [1 #]   [1 #]                            [1 #]   [1 #]
// [2 #]   [2 #]                            [2 #]   [2 #]
// [3 #]*  [3 X]                            [3 #]*  [3 X]
// [4 #]   [4 #]                            [4 #]   [4 #]
// [5 #]   [5 #]                            [5 #]   [5 #]
// [6  ]   [6 #]                            [6 #]   [6 #]
// [7  ]   [7  ]                            [7 #]   [7 #]
// [8  ]   [8  ]                            [8 #]   [8 #]
// [9  ]   [9  ]                            [9 #]   [9 #]
//
// Our caller, the 'top' verb, feeds us records. We keep them or not; in the
// latter case, the Go runtime GCs them.

func (this *TopKeeper) Add(value *types.Mlrval, recordAndContext *types.RecordAndContext) {
	destidx := this.bsearchFunc(&this.TopValues, this.size, value)

	if this.size < this.capacity {
		for i := this.size - 1; i >= destidx; i-- {
			this.TopValues[i+1] = this.TopValues[i]
			this.TopRecordsAndContexts[i+1] = this.TopRecordsAndContexts[i]
		}
		this.TopValues[destidx] = value.Copy()
		this.TopRecordsAndContexts[destidx] = recordAndContext.Copy() // might be nil
		this.size++
	} else {
		if destidx >= this.capacity {
			return
		}
		for i := this.size - 2; i >= destidx; i-- {
			this.TopValues[i+1] = this.TopValues[i]
			this.TopRecordsAndContexts[i+1] = this.TopRecordsAndContexts[i]
		}
		this.TopValues[destidx] = value.Copy()
		this.TopRecordsAndContexts[destidx] = recordAndContext.Copy()
	}
}
