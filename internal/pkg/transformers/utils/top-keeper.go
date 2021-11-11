// ================================================================
// Data structure for mlr top: just a decorated array.
// ================================================================

package utils

import (
	"mlr/internal/pkg/types"
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
	keeper := &TopKeeper{
		TopValues:             make([]*types.Mlrval, capacity),
		TopRecordsAndContexts: make([]*types.RecordAndContext, capacity),
		size:                  0,
		capacity:              capacity,
	}
	if doMax {
		keeper.bsearchFunc = types.BsearchMlrvalArrayForDescendingInsert
	} else {
		keeper.bsearchFunc = types.BsearchMlrvalArrayForAscendingInsert
	}
	return keeper
}

func (keeper *TopKeeper) GetSize() int {
	return keeper.size
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

func (keeper *TopKeeper) Add(value *types.Mlrval, recordAndContext *types.RecordAndContext) {
	destidx := keeper.bsearchFunc(&keeper.TopValues, keeper.size, value)

	if keeper.size < keeper.capacity {
		for i := keeper.size - 1; i >= destidx; i-- {
			keeper.TopValues[i+1] = keeper.TopValues[i]
			keeper.TopRecordsAndContexts[i+1] = keeper.TopRecordsAndContexts[i]
		}
		keeper.TopValues[destidx] = value.Copy()
		keeper.TopRecordsAndContexts[destidx] = recordAndContext.Copy() // might be nil
		keeper.size++
	} else {
		if destidx >= keeper.capacity {
			return
		}
		for i := keeper.size - 2; i >= destidx; i-- {
			keeper.TopValues[i+1] = keeper.TopValues[i]
			keeper.TopRecordsAndContexts[i+1] = keeper.TopRecordsAndContexts[i]
		}
		keeper.TopValues[destidx] = value.Copy()
		keeper.TopRecordsAndContexts[destidx] = recordAndContext.Copy()
	}
}
