package utils

import (
	"github.com/johnkerl/miller/pkg/lib"
)

// WindowKeeper is a sliding-window container, nominally for use by mlr step,
// for holding a number of records before the current one, the current one, and
// a number of records after. The payload is interface{}, not *mlrval.Mlrmap,
// for ease of unit-testing -- and also, since nothing here inspects the
// payload, so that this code could be repurposed.
type TWindowKeeper struct {
	numBackward int
	numForward  int

	itemsBackward []interface{}
	currentItem   interface{}
	itemsForward  []interface{}
}

func NewWindowKeeper(
	numBackward int,
	numForward int,
) *TWindowKeeper {
	return &TWindowKeeper{
		numBackward: numBackward,
		numForward:  numForward,

		itemsBackward: make([]interface{}, numBackward),
		currentItem:   nil,
		itemsForward:  make([]interface{}, numForward),
	}
}

func (wk *TWindowKeeper) Ingest(
	inrec interface{},
) {
	for i := wk.numBackward - 1; i > 0; i-- {
		wk.itemsBackward[i] = wk.itemsBackward[i-1]
	}
	if wk.numBackward > 0 {
		wk.itemsBackward[0] = wk.currentItem
	}
	if wk.numForward > 0 {
		wk.currentItem = wk.itemsForward[0]
		for i := 0; i < wk.numForward-1; i++ {
			wk.itemsForward[i] = wk.itemsForward[i+1]
		}
		wk.itemsForward[wk.numForward-1] = inrec
	} else {
		wk.currentItem = inrec
	}
}

// Get maps a user-visible indexing ..., -3, -2, -1, 0, 1, 2, 3, ...
// into this struct's zero-index array storage.
func (wk *TWindowKeeper) Get(
	index int,
) interface{} {
	if index == 0 {
		return wk.currentItem
	} else if index > 0 {
		lib.InternalCodingErrorIf(index > wk.numForward)
		return wk.itemsForward[index-1]
	} else {
		index = -index
		lib.InternalCodingErrorIf(index > wk.numBackward)
		return wk.itemsBackward[index-1]
	}
}
