package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	i10 := FromInt(10)
	i2 := FromInt(2)

	bfalse := FromBool(false)
	btrue := FromBool(true)

	sabc := FromString("abc")
	sdef := FromString("def")

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Within-type lexical comparisons
	assert.Equal(t, -1, LexicalAscendingComparator(i10, i2))
	assert.Equal(t, 0, LexicalAscendingComparator(bfalse, bfalse))
	assert.Equal(t, -1, LexicalAscendingComparator(bfalse, btrue))
	assert.Equal(t, -1, LexicalAscendingComparator(sabc, sdef))
	assert.Equal(t, 0, LexicalAscendingComparator(FromErrorString("foo"), FromErrorString("foo")))
	assert.Equal(t, 0, LexicalAscendingComparator(ABSENT, ABSENT))

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Within-type numeric comparisons
	assert.Equal(t, 1, NumericAscendingComparator(i10, i2))
	assert.Equal(t, 0, NumericAscendingComparator(sabc, sabc))
	assert.Equal(t, -1, NumericAscendingComparator(sabc, sdef))

	assert.Equal(t, 1, NumericAscendingComparator(btrue, bfalse))

	assert.Equal(t, 0, NumericAscendingComparator(FromErrorString("foo"), FromErrorString("foo")))
	assert.Equal(t, 0, NumericAscendingComparator(ABSENT, ABSENT))

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Across-type lexical comparisons

	assert.Equal(t, -1, LexicalAscendingComparator(i10, btrue))                 // "10" < "true"
	assert.Equal(t, -1, LexicalAscendingComparator(i10, sabc))                  // "10" < "abc"
	assert.Equal(t, 1, LexicalAscendingComparator(i10, FromErrorString("foo"))) // "10" > "(error)"

	assert.Equal(t, 1, LexicalAscendingComparator(bfalse, sabc))                   // "false" > "abc"
	assert.Equal(t, 1, LexicalAscendingComparator(bfalse, FromErrorString("foo"))) // "false" > "(error)"

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Across-type numeric comparisons

	assert.Equal(t, -1, NumericAscendingComparator(i10, btrue))
	assert.Equal(t, -1, NumericAscendingComparator(i10, sabc))
	assert.Equal(t, -1, NumericAscendingComparator(i10, FromErrorString("foo")))
	assert.Equal(t, -1, NumericAscendingComparator(i10, ABSENT))

	assert.Equal(t, -1, NumericAscendingComparator(bfalse, sabc))
	assert.Equal(t, -1, NumericAscendingComparator(bfalse, FromErrorString("foo")))
	assert.Equal(t, -1, NumericAscendingComparator(bfalse, ABSENT))

	assert.Equal(t, -1, NumericAscendingComparator(FromErrorString("foo"), ABSENT))
}
