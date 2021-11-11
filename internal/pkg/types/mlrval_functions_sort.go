// ================================================================
// For sorting
// ================================================================

package types

import (
	"fmt"
	"os"
	"strings"
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
	if input1.mvtype == MT_STRING {
		sa = strings.ToLower(sa)
	}
	if input2.mvtype == MT_STRING {
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

// ----------------------------------------------------------------
// Sort rules (same for min, max, and comparator):
// * NUMERICS < BOOL < STRINGS < ERROR < ABSENT
// * error == error (singleton type)
// * absent == absent (singleton type)
// * string compares on strings
// * numeric compares on numbers
// * false < true

func _neg1(input1, input2 *Mlrval) int {
	return -1
}
func _zero(input1, input2 *Mlrval) int {
	return 0
}
func _pos1(input1, input2 *Mlrval) int {
	return 1
}

func _scmp(input1, input2 *Mlrval) int {
	if input1.printrep < input2.printrep {
		return -1
	} else if input1.printrep > input2.printrep {
		return 1
	} else {
		return 0
	}
}

func iicmp(input1, input2 *Mlrval) int {
	ca := input1.intval
	cb := input2.intval
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}
func ifcmp(input1, input2 *Mlrval) int {
	ca := float64(input1.intval)
	cb := input2.floatval
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}
func ficmp(input1, input2 *Mlrval) int {
	ca := input1.floatval
	cb := float64(input2.intval)
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}
func ffcmp(input1, input2 *Mlrval) int {
	ca := input1.floatval
	cb := input2.floatval
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}

func bbcmp(input1, input2 *Mlrval) int {
	a := input1.boolval
	b := input2.boolval
	if a == false {
		if b == false {
			return 0
		} else {
			return -1
		}
	} else {
		if b == false {
			return 1
		} else {
			return 0
		}
	}
}

func _xcmp(input1, input2 *Mlrval) int {
	fmt.Fprintf(os.Stderr, "mlr: functions cannot be sorted.\n")
	os.Exit(1)
	return 0
}

// ----------------------------------------------------------------
// Sort rules (same for min, max, and comparator):
// * NUMERICS < BOOL < STRINGS < ERROR < NULL < ABSENT
// * error == error (singleton type)
// * absent == absent (singleton type)
// * null == null (singleton type)
// * string compares on strings
// * numeric compares on numbers
// * false < true

// typed_cmp_dispositions is the disposition matrix for numerical sorting of Mlrvals.
var typed_cmp_typedpositions = [MT_DIM][MT_DIM]ComparatorFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT    FLOAT  BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_zero, _neg1, _neg1, _pos1, _pos1, _pos1, _pos1, _pos1, _zero, _zero, _xcmp},
	/*ABSENT */ {_pos1, _zero, _pos1, _pos1, _pos1, _pos1, _pos1, _pos1, _zero, _zero, _xcmp},
	/*NULL   */ {_pos1, _neg1, _zero, _pos1, _pos1, _pos1, _pos1, _pos1, _pos1, _pos1, _xcmp},
	/*VOID   */ {_neg1, _neg1, _neg1, _scmp, _scmp, _pos1, _pos1, _pos1, _zero, _zero, _xcmp},
	/*STRING */ {_neg1, _neg1, _neg1, _scmp, _scmp, _pos1, _pos1, _pos1, _zero, _zero, _xcmp},
	/*INT    */ {_neg1, _neg1, _neg1, _neg1, _neg1, iicmp, ifcmp, _neg1, _zero, _zero, _xcmp},
	/*FLOAT  */ {_neg1, _neg1, _neg1, _neg1, _neg1, ficmp, ffcmp, _neg1, _zero, _zero, _xcmp},
	/*BOOL   */ {_neg1, _neg1, _neg1, _neg1, _neg1, _pos1, _pos1, bbcmp, _zero, _zero, _xcmp},
	/*ARRAY  */ {_zero, _zero, _neg1, _zero, _zero, _zero, _zero, _zero, _zero, _zero, _xcmp},
	/*MAP    */ {_zero, _zero, _neg1, _zero, _zero, _zero, _zero, _zero, _zero, _zero, _xcmp},
	/*FUNC    */ {_xcmp, _xcmp, _xcmp, _xcmp, _xcmp, _xcmp, _xcmp, _xcmp, _xcmp, _xcmp, _xcmp},
}

// NumericAscendingComparator is for "numerical" sort: it uses Mlrval sorting
// rules by type, including numeric sort for numeric types.
func NumericAscendingComparator(input1 *Mlrval, input2 *Mlrval) int {
	return typed_cmp_typedpositions[input1.mvtype][input2.mvtype](input1, input2)
}

// NumericDescendingComparator is for "numerical" sort: it uses Mlrval sorting
// rules by type, including numeric sort for numeric types.
func NumericDescendingComparator(input1 *Mlrval, input2 *Mlrval) int {
	return NumericAscendingComparator(input2, input1)
}
