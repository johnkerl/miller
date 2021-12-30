package bifs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

// ----------------------------------------------------------------
// SORTING
//
// Lexical compare is just string-sort on stringify of mlrvals:
// e.g. "hello" < "true".
//
// Numerical sort rules (same for min, max, and comparator):
// * NUMERICS < BOOL < STRINGS < ERROR < ABSENT
// * error == error (singleton type)
// * absent == absent (singleton type)
// * string compares on strings
// * numeric compares on numbers
// * false < true

func TestComparators(t *testing.T) {

	i10 := mlrval.FromInt(10)
	i2 := mlrval.FromInt(2)

	bfalse := mlrval.FromBool(false)
	btrue := mlrval.FromBool(true)

	sabc := mlrval.FromString("abc")
	sdef := mlrval.FromString("def")

	e := *mlrval.ERROR

	a := *mlrval.ABSENT

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Within-type lexical comparisons
	assert.Equal(t, -1, LexicalAscendingComparator(i10, i2))
	assert.Equal(t, 0, LexicalAscendingComparator(bfalse, bfalse))
	assert.Equal(t, -1, LexicalAscendingComparator(bfalse, btrue))
	assert.Equal(t, -1, LexicalAscendingComparator(sabc, sdef))
	assert.Equal(t, 0, LexicalAscendingComparator(&e, &e))
	assert.Equal(t, 0, LexicalAscendingComparator(&a, &a))

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Within-type numeric comparisons
	assert.Equal(t, 1, NumericAscendingComparator(i10, i2))
	assert.Equal(t, 0, NumericAscendingComparator(sabc, sabc))
	assert.Equal(t, -1, NumericAscendingComparator(sabc, sdef))

	assert.Equal(t, 1, NumericAscendingComparator(btrue, bfalse))

	assert.Equal(t, 0, NumericAscendingComparator(&e, &e))
	assert.Equal(t, 0, NumericAscendingComparator(&a, &a))

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Across-type lexical comparisons

	assert.Equal(t, -1, LexicalAscendingComparator(i10, btrue)) // "10" < "true"
	assert.Equal(t, -1, LexicalAscendingComparator(i10, sabc))   // "10" < "abc"
	assert.Equal(t, 1, LexicalAscendingComparator(i10, &e))     // "10" > "(error)"

	assert.Equal(t, 1, LexicalAscendingComparator(bfalse, sabc)) // "false" > "abc"
	assert.Equal(t, 1, LexicalAscendingComparator(bfalse, &e))   // "false" > "(error)"

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Across-type numeric comparisons

	assert.Equal(t, -1, NumericAscendingComparator(i10, btrue))
	assert.Equal(t, -1, NumericAscendingComparator(i10, sabc))
	assert.Equal(t, -1, NumericAscendingComparator(i10, &e))
	assert.Equal(t, -1, NumericAscendingComparator(i10, &a))

	assert.Equal(t, -1, NumericAscendingComparator(bfalse, sabc))
	assert.Equal(t, -1, NumericAscendingComparator(bfalse, &e))
	assert.Equal(t, -1, NumericAscendingComparator(bfalse, &a))

	assert.Equal(t, -1, NumericAscendingComparator(&e, &a))
}
