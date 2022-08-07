// ================================================================
// For sorting
//
// Sort rules (same for min, max, and comparator):
// * NUMERICS < BOOL < STRINGS < ERROR < ABSENT
// * error == error (singleton type)
// * absent == absent (singleton type)
// * string compares on strings
// * numeric compares on numbers
// * false < true
// ================================================================
// ================================================================

package mlrval

import (
	"strings"

	"github.com/facette/natsort"
)

// LexicalAscendingComparator is for lexical sort: it stringifies
// everything.
func LexicalAscendingComparator(input1 *Mlrval, input2 *Mlrval) int {
	sa := input1.String()
	sb := input2.String()
	if sa < sb {
		return -1
	} else if sa > sb {
		return 1
	} else {
		return 0
	}
}

// LexicalDescendingComparator is for reverse-lexical sort: it stringifies
// everything.
func LexicalDescendingComparator(input1 *Mlrval, input2 *Mlrval) int {
	return LexicalAscendingComparator(input2, input1)
}

// CaseFoldAscendingComparator is for case-folded lexical sort: it stringifies
// everything.
func CaseFoldAscendingComparator(input1 *Mlrval, input2 *Mlrval) int {
	sa := input1.String()
	sb := input2.String()
	if input1.IsString() {
		sa = strings.ToLower(sa)
	}
	if input2.IsString() {
		sb = strings.ToLower(sb)
	}
	if sa < sb {
		return -1
	} else if sa > sb {
		return 1
	} else {
		return 0
	}
}

// CaseFoldDescendingComparator is for case-folded lexical sort: it stringifies
// everything.
func CaseFoldDescendingComparator(input1 *Mlrval, input2 *Mlrval) int {
	return CaseFoldAscendingComparator(input2, input1)
}

func NumericAscendingComparator(input1 *Mlrval, input2 *Mlrval) int {
	return Cmp(input1, input2)
}

// NumericDescendingComparator is for "numerical" sort: it uses Mlrval sorting
// rules by type, including numeric sort for numeric types.
func NumericDescendingComparator(input1 *Mlrval, input2 *Mlrval) int {
	return -Cmp(input1, input2)
}

func NaturalAscendingComparator(input1, input2 *Mlrval) int {
	sa := input1.String()
	sb := input2.String()
	if sa == sb {
		return 0
	}

	// natsort.Compare puts empty strings in random places
	if sa == "" {
		return 1
	}
	if sb == "" {
		return -1
	}

	if natsort.Compare(input1.String(), input2.String()) {
		return 1
	} else {
		return -1
	}
}

func NaturalDescendingComparator(input1, input2 *Mlrval) int {
	return NaturalAscendingComparator(input2, input1)
}
