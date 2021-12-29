// ================================================================
// Tests mlrval comparator functions
// ================================================================

package mlrval

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Documented contract:
// NUMERICS < BOOL < VOID < STRING < ERROR < NULL < ABSENT

var orderedMlrvals = []*Mlrval{

	FromInt(1),
	FromFloat(1.1),
	FromInt(2),
	FromFloat(2.2),

	FromBool(false),
	FromBool(true),

	FromString(""),
	FromString("abc"),
	FromString("defgh"),

	// TODO:
	// FromArray([]Mlrval{FromInt(1), FromInt(2)}),
	// FromArray([]Mlrval{FromInt(1), FromInt(3)}),
	// FromMap(NewMlrmap()),

	// TODO:
	// ERROR,
	// ABSENT,
	// NULL,
}

func TestEqual(t *testing.T) {
	for i := range orderedMlrvals {
		mvi := orderedMlrvals[i]
		for j := range orderedMlrvals {
			mvj := orderedMlrvals[j]
			assert.Equal(t, i == j, Equals(mvi, mvj), fmt.Sprintf(
				"slots i=%d type=%s value=%s, j=%d type=%s value=%s",
				i, mvi.GetTypeName(), mvi.String(),
				j, mvj.GetTypeName(), mvj.String(),
			))
		}
	}
}

func TestNotEquals(t *testing.T) {
	for i := range orderedMlrvals {
		mvi := orderedMlrvals[i]
		for j := range orderedMlrvals {
			mvj := orderedMlrvals[j]
			assert.Equal(t, i != j, NotEquals(mvi, mvj), fmt.Sprintf(
				"slots i=%d type=%s value=%s, j=%d type=%s value=%s",
				i, mvi.GetTypeName(), mvi.String(),
				j, mvj.GetTypeName(), mvj.String(),
			))
		}
	}
}

func TestLessThan(t *testing.T) {
	for i := range orderedMlrvals {
		mvi := orderedMlrvals[i]
		for j := range orderedMlrvals {
			mvj := orderedMlrvals[j]
			assert.Equal(t, i < j, LessThan(mvi, mvj), fmt.Sprintf(
				"slots i=%d type=%s value=%s, j=%d type=%s value=%s",
				i, mvi.GetTypeName(), mvi.String(),
				j, mvj.GetTypeName(), mvj.String(),
			))
		}
	}
}

//func TestLessThanOrEquals(t *testing.T) {
//	for i := range orderedMlrvals {
//		mvi := orderedMlrvals[i]
//		for j := range orderedMlrvals {
//			mvj := orderedMlrvals[j]
//			assert.Equal(t, i <= j, LessThan(mvi, mvj), fmt.Sprintf(
//				"slots i=%d type=%s value=%s, j=%d type=%s value=%s",
//				i, mvi.GetTypeName(), mvi.String(),
//				j, mvj.GetTypeName(), mvj.String(),
//			))
//		}
//	}
//}

//func TestGreaterThan(t *testing.T) {
//	for i := range orderedMlrvals {
//		mvi := orderedMlrvals[i]
//		for j := range orderedMlrvals {
//			mvj := orderedMlrvals[j]
//			assert.Equal(t, i > j, GreaterThan(mvi, mvj), fmt.Sprintf(
//				"slots i=%d type=%s value=%s, j=%d type=%s value=%s",
//				i, mvi.GetTypeName(), mvi.String(),
//				j, mvj.GetTypeName(), mvj.String(),
//			))
//		}
//	}
//}

//func TestGreaterThanOrEquals(t *testing.T) {
//	for i := range orderedMlrvals {
//		mvi := orderedMlrvals[i]
//		for j := range orderedMlrvals {
//			mvj := orderedMlrvals[j]
//			assert.Equal(t, i >= j, GreaterThan(mvi, mvj), fmt.Sprintf(
//				"slots i=%d type=%s value=%s, j=%d type=%s value=%s",
//				i, mvi.GetTypeName(), mvi.String(),
//				j, mvj.GetTypeName(), mvj.String(),
//			))
//		}
//	}
//}

//func TestCmp(t *testing.T) {
//	for i := range orderedMlrvals {
//		mvi := orderedMlrvals[i]
//		for j := range orderedMlrvals {
//			mvj := orderedMlrvals[j]
//			expect := 0
//			if i < j {
//				expect = -1
//			} else if i > j {
//				expect = 1
//			}
//			assert.Equal(t, expect, Cmp(mvi, mvj), fmt.Sprintf(
//				"slots i=%d type=%s value=%s, j=%d type=%s value=%s",
//				i, mvi.GetTypeName(), mvi.String(),
//				j, mvj.GetTypeName(), mvj.String(),
//			))
//		}
//	}
//}
