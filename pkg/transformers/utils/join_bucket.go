// ================================================================
// Helper data structure for the join verb
// ================================================================

package utils

import (
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

// ----------------------------------------------------------------
type JoinBucket struct {
	leftFieldValues    []*mlrval.Mlrval
	RecordsAndContexts []*types.RecordAndContext
	WasPaired          bool
}

func NewJoinBucket(
	leftFieldValues []*mlrval.Mlrval,
) *JoinBucket {
	return &JoinBucket{
		leftFieldValues:    leftFieldValues,
		RecordsAndContexts: make([]*types.RecordAndContext, 0),
		WasPaired:          false,
	}
}
