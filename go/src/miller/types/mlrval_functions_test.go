package types

// ================================================================
// Most Miller tests (thousands of them) are command-line-driven via
// reg_test/run. Here are some cases needing special focus.
// ================================================================

import (
	"testing"
	//"types"
)

// ----------------------------------------------------------------
func TestLexicalAscendingComparator(t *testing.T) {
	var a = MlrvalFromInt64(10)
	var b = MlrvalFromInt64(2)
	if LexicalAscendingComparator(&a, &b) != -1 {
		t.Fatal()
	}

	a = MlrvalFromString("abc")
	b = MlrvalFromString("def")
	if LexicalAscendingComparator(&a, &b) != -1 {
		t.Fatal()
	}

	a = MlrvalFromInt64(3)
	b = MlrvalFromBool(true)
	if LexicalAscendingComparator(&a, &b) != -1 {
		t.Fatal()
	}
}

// ----------------------------------------------------------------
// SORTING
// Sort rules (same for min, max, and comparator):
// * NUMERICS < BOOL < STRINGS < ERROR < ABSENT
// * error == error (singleton type)
// * absent == absent (singleton type)
// * string compares on strings
// * numeric compares on numbers
// * false < true

func TestNumericAscendingComparator(t *testing.T) {
	i10 := MlrvalFromInt64(10)
	i2 := MlrvalFromInt64(2)

	bfalse := MlrvalFromBool(false)
	btrue := MlrvalFromBool(true)

	sabc := MlrvalFromString("abc")
	sdef := MlrvalFromString("def")

	e := MlrvalFromError()

	a := MlrvalFromAbsent()

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Within-type comparisons
	if NumericAscendingComparator(&i10, &i2) != 1 {
		t.Fatal()
	}
	if NumericAscendingComparator(&sabc, &sabc) != 0 {
		t.Fatal()
	}
	if NumericAscendingComparator(&sabc, &sdef) != -1 {
		t.Fatal()
	}

	if NumericAscendingComparator(&btrue, &bfalse) != 1 {
		t.Fatal()
	}

	if NumericAscendingComparator(&e, &e) != 0 {
		t.Fatal()
	}
	if NumericAscendingComparator(&a, &a) != 0 {
		t.Fatal()
	}
	// xxx more

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Across-type comparisons

	if NumericAscendingComparator(&i10, &btrue) != -1 {
		t.Fatal()
	}
	if NumericAscendingComparator(&e, &a) != -1 {
		t.Fatal()
	}

	// xxx more
}
