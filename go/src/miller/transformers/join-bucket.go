// ================================================================
// Helper data structure for the join verb
// ================================================================

package transformers

import (
	"container/list"

	"miller/types"
)

// ----------------------------------------------------------------
type tJoinBucket struct {
	leftFieldValues    []types.Mlrval
	recordsAndContexts *list.List
	wasPaired          bool
}

func newJoinBucket(
	leftFieldValues []types.Mlrval,
) *tJoinBucket {
	return &tJoinBucket{
		leftFieldValues:    leftFieldValues,
		recordsAndContexts: list.New(),
		wasPaired:          false,
	}
}
