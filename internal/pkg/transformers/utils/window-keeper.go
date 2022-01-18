package utils

import (
	"github.com/johnkerl/miller/internal/pkg/lib"
)

// WindowKeeper is a sliding-window container, nominally for use by mlr step,
// for holding a number of records before the current one, the current one, and
// a number of records after. The payload is interface{}, not *mlrval.Mlrmap,
// for ease of unit-testing -- as well as since nothing here inspects the
// payload, so this code could be repurposed.
type WindowKeeper struct {
	numBackward int
	numForward  int

	recordsBackward []interface{}
	currentRecord   interface{}
	recordsForward  []interface{}
}

func NewWindowKeeper(
	numBackward int,
	numForward int,
) *WindowKeeper {
	return &WindowKeeper{
		numBackward: numBackward,
		numForward:  numForward,

		recordsBackward: make([]interface{}, numBackward),
		currentRecord:   nil,
		recordsForward:  make([]interface{}, numForward),
	}
}

func (wk *WindowKeeper) IngestRecord(
	inrec interface{},
) {
	for i := wk.numBackward - 1; i > 0; i-- {
		wk.recordsBackward[i] = wk.recordsBackward[i-1]
	}
	if wk.numBackward > 0 {
		wk.recordsBackward[0] = wk.currentRecord
	}
	if wk.numForward > 0 {
		wk.currentRecord = wk.recordsForward[0]
		for i := 0; i < wk.numForward-1; i++ {
			wk.recordsForward[i] = wk.recordsForward[i+1]
		}
		wk.recordsForward[wk.numForward-1] = inrec
	} else {
		wk.currentRecord = inrec
	}
}

// GetRecord maps a user-visible indexing ..., -3, -2, -1, 0, 1, 2, 3, ...
// into this struct's zero-index array storage.
func (wk *WindowKeeper) GetRecord(
	index int,
) interface{} {
	if index == 0 {
		return wk.currentRecord
	} else if index > 0 {
		lib.InternalCodingErrorIf(index > wk.numForward)
		return wk.recordsForward[index-1]
	} else {
		index = -index
		lib.InternalCodingErrorIf(index > wk.numBackward)
		return wk.recordsBackward[index-1]
	}
}
