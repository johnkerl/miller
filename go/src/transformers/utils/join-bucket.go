// ================================================================
// Helper data structure for the join verb
// ================================================================

package utils

import (
	"container/list"

	"miller/src/types"
)

// ----------------------------------------------------------------
type JoinBucket struct {
	leftFieldValues    []*types.Mlrval
	RecordsAndContexts *list.List
	WasPaired          bool
}

func NewJoinBucket(
	leftFieldValues []*types.Mlrval,
) *JoinBucket {
	return &JoinBucket{
		leftFieldValues:    leftFieldValues,
		RecordsAndContexts: list.New(),
		WasPaired:          false,
	}
}
