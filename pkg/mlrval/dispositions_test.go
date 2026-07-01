package mlrval

import (
	"testing"
)

// See pkg/bifs/dispositions_test.go: guards against nil cells left by
// zero-fill when MT_DIM grows.
func TestNoNilCellsInCmpDispositions(t *testing.T) {
	for i := 0; i < int(MT_DIM); i++ {
		for j := 0; j < int(MT_DIM); j++ {
			if cmp_dispositions[i][j] == nil {
				t.Errorf("nil cell in cmp_dispositions[%d][%d]", i, j)
			}
		}
	}
}
