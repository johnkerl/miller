// ================================================================
// Most Miller tests (thousands of them) are command-line-driven via
// mlr regtest. Here are some cases needing special focus.
// ================================================================

package types

import (
	"testing"
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

	i10 := MlrvalFromInt(10)
	i2 := MlrvalFromInt(2)

	bfalse := MlrvalFromBool(false)
	btrue := MlrvalFromBool(true)

	sabc := MlrvalFromString("abc")
	sdef := MlrvalFromString("def")

	e := *MLRVAL_ERROR

	a := *MLRVAL_ABSENT

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Within-type lexical comparisons
	if LexicalAscendingComparator(i10, i2) != -1 {
		t.Fatal()
	}
	if LexicalAscendingComparator(bfalse, bfalse) != 0 {
		t.Fatal()
	}
	if LexicalAscendingComparator(bfalse, btrue) != -1 {
		t.Fatal()
	}
	if LexicalAscendingComparator(sabc, sdef) != -1 {
		t.Fatal()
	}
	if LexicalAscendingComparator(&e, &e) != 0 {
		t.Fatal()
	}
	if LexicalAscendingComparator(&a, &a) != 0 {
		t.Fatal()
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Within-type numeric comparisons
	if NumericAscendingComparator(i10, i2) != 1 {
		t.Fatal()
	}
	if NumericAscendingComparator(sabc, sabc) != 0 {
		t.Fatal()
	}
	if NumericAscendingComparator(sabc, sdef) != -1 {
		t.Fatal()
	}

	if NumericAscendingComparator(btrue, bfalse) != 1 {
		t.Fatal()
	}

	if NumericAscendingComparator(&e, &e) != 0 {
		t.Fatal()
	}
	if NumericAscendingComparator(&a, &a) != 0 {
		t.Fatal()
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Across-type lexical comparisons

	if LexicalAscendingComparator(i10, btrue) != -1 { // "10" < "true"
		t.Fatal()
	}
	if LexicalAscendingComparator(i10, sabc) != -1 { // "10" < "abc"
		t.Fatal()
	}
	if LexicalAscendingComparator(i10, &e) != 1 { // "10" > "(error)"
		t.Fatal()
	}

	if LexicalAscendingComparator(bfalse, sabc) != 1 { // "false" > "abc"
		t.Fatal()
	}
	if LexicalAscendingComparator(bfalse, &e) != 1 { // "false" > "(error)"
		t.Fatal()
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Across-type numeric comparisons

	if NumericAscendingComparator(i10, btrue) != -1 {
		t.Fatal()
	}
	if NumericAscendingComparator(i10, sabc) != -1 {
		t.Fatal()
	}
	if NumericAscendingComparator(i10, &e) != -1 {
		t.Fatal()
	}
	if NumericAscendingComparator(i10, &a) != -1 {
		t.Fatal()
	}

	if NumericAscendingComparator(bfalse, sabc) != -1 {
		t.Fatal()
	}
	if NumericAscendingComparator(bfalse, &e) != -1 {
		t.Fatal()
	}
	if NumericAscendingComparator(bfalse, &a) != -1 {
		t.Fatal()
	}

	if NumericAscendingComparator(&e, &a) != -1 {
		t.Fatal()
	}
}
