// ================================================================
// Helper data structure for the join verb
// ================================================================

package utils

import (
	"container/list"

	"github.com/johnkerl/miller/pkg/mlrval"
)

// ----------------------------------------------------------------
type JoinBucket struct {
	leftFieldValues    []*mlrval.Mlrval
	RecordsAndContexts *list.List
	WasPaired          bool
}

func NewJoinBucket(
	leftFieldValues []*mlrval.Mlrval,
) *JoinBucket {
	return &JoinBucket{
		leftFieldValues:    leftFieldValues,
		RecordsAndContexts: list.New(),
		WasPaired:          false,
	}
}
