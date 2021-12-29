// ================================================================
// For sorting
// ================================================================

package bifs

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

// LexicalAscendingComparator is for lexical sort: it stringifies
// everything.
func LexicalAscendingComparator(input1 *mlrval.Mlrval, input2 *mlrval.Mlrval) int {
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
func LexicalDescendingComparator(input1 *mlrval.Mlrval, input2 *mlrval.Mlrval) int {
	return LexicalAscendingComparator(input2, input1)
}

// CaseFoldAscendingComparator is for case-folded lexical sort: it stringifies
// everything.
func CaseFoldAscendingComparator(input1 *mlrval.Mlrval, input2 *mlrval.Mlrval) int {
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
func CaseFoldDescendingComparator(input1 *mlrval.Mlrval, input2 *mlrval.Mlrval) int {
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

func _neg1(input1, input2 *mlrval.Mlrval) int {
	return -1
}
func _zero(input1, input2 *mlrval.Mlrval) int {
	return 0
}
func _pos1(input1, input2 *mlrval.Mlrval) int {
	return 1
}

func _scmp(input1, input2 *mlrval.Mlrval) int {
	if input1.AcquireStringValue() < input2.AcquireStringValue() {
		return -1
	} else if input1.AcquireStringValue() > input2.AcquireStringValue() {
		return 1
	} else {
		return 0
	}
}

func iicmp(input1, input2 *mlrval.Mlrval) int {
	ca := input1.AcquireIntValue()
	cb := input2.AcquireIntValue()
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}
func ifcmp(input1, input2 *mlrval.Mlrval) int {
	ca := float64(input1.AcquireIntValue())
	cb := input2.AcquireFloatValue()
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}
func ficmp(input1, input2 *mlrval.Mlrval) int {
	ca := input1.AcquireFloatValue()
	cb := float64(input2.AcquireIntValue())
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}
func ffcmp(input1, input2 *mlrval.Mlrval) int {
	ca := input1.AcquireFloatValue()
	cb := input2.AcquireFloatValue()
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}

func bbcmp(input1, input2 *mlrval.Mlrval) int {
	a := input1.AcquireBoolValue()
	b := input2.AcquireBoolValue()
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

func _xcmp(input1, input2 *mlrval.Mlrval) int {
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
var typed_cmp_typedpositions = [mlrval.MT_DIM][mlrval.MT_DIM]ComparatorFunc{
//       .  INT    FLOAT  BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
/*INT    */ {iicmp, ifcmp, _neg1, _neg1, _neg1, _zero, _zero, _xcmp, _neg1, _neg1, _neg1},
/*FLOAT  */ {ficmp, ffcmp, _neg1, _neg1, _neg1, _zero, _zero, _xcmp, _neg1, _neg1, _neg1},
/*BOOL   */ {_pos1, _pos1, bbcmp, _neg1, _neg1, _zero, _zero, _xcmp, _neg1, _neg1, _neg1},
/*VOID   */ {_pos1, _pos1, _pos1, _scmp, _scmp, _zero, _zero, _xcmp, _neg1, _neg1, _neg1},
/*STRING */ {_pos1, _pos1, _pos1, _scmp, _scmp, _zero, _zero, _xcmp, _neg1, _neg1, _neg1},
/*ARRAY  */ {_zero, _zero, _zero, _zero, _zero, _zero, _zero, _xcmp, _zero, _neg1, _zero},
/*MAP    */ {_zero, _zero, _zero, _zero, _zero, _zero, _zero, _xcmp, _zero, _neg1, _zero},
/*FUNC   */ {_xcmp, _xcmp, _xcmp, _xcmp, _xcmp, _xcmp, _xcmp, _xcmp, _xcmp, _xcmp, _xcmp},
/*ERROR  */ {_pos1, _pos1, _pos1, _pos1, _pos1, _zero, _zero, _xcmp, _zero, _neg1, _neg1},
/*NULL   */ {_pos1, _pos1, _pos1, _pos1, _pos1, _pos1, _pos1, _xcmp, _pos1, _zero, _neg1},
/*ABSENT */ {_pos1, _pos1, _pos1, _pos1, _pos1, _zero, _zero, _xcmp, _pos1, _pos1, _zero},
}

// NumericAscendingComparator is for "numerical" sort: it uses Mlrval sorting
// rules by type, including numeric sort for numeric types.
func NumericAscendingComparator(input1 *mlrval.Mlrval, input2 *mlrval.Mlrval) int {
	return typed_cmp_typedpositions[input1.Type()][input2.Type()](input1, input2)
}

// NumericDescendingComparator is for "numerical" sort: it uses Mlrval sorting
// rules by type, including numeric sort for numeric types.
func NumericDescendingComparator(input1 *mlrval.Mlrval, input2 *mlrval.Mlrval) int {
	return NumericAscendingComparator(input2, input1)
}
