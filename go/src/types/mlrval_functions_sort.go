// ================================================================
// For sorting
// ================================================================

package types

// Lexical sort: just stringify everything.
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
func LexicalDescendingComparator(input1 *Mlrval, input2 *Mlrval) int {
	return LexicalAscendingComparator(input2, input1)
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

// ----------------------------------------------------------------
// Sort rules (same for min, max, and comparator):
// * NUMERICS < BOOL < STRINGS < ERROR < NULL < ABSENT
// * error == error (singleton type)
// * absent == absent (singleton type)
// * null == null (singleton type)
// * string compares on strings
// * numeric compares on numbers
// * false < true

var num_cmp_dispositions = [MT_DIM][MT_DIM]ComparatorFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT    FLOAT  BOOL   ARRAY  MAP
	/*ERROR  */ {_zero, _neg1, _neg1, _pos1, _pos1, _pos1, _pos1, _pos1, _zero, _zero},
	/*ABSENT */ {_pos1, _zero, _pos1, _pos1, _pos1, _pos1, _pos1, _pos1, _zero, _zero},
	/*NULL   */ {_pos1, _neg1, _zero, _pos1, _pos1, _pos1, _pos1, _pos1, _pos1, _pos1},
	/*VOID   */ {_neg1, _neg1, _neg1, _scmp, _scmp, _pos1, _pos1, _pos1, _zero, _zero},
	/*STRING */ {_neg1, _neg1, _neg1, _scmp, _scmp, _pos1, _pos1, _pos1, _zero, _zero},
	/*INT    */ {_neg1, _neg1, _neg1, _neg1, _neg1, iicmp, ifcmp, _neg1, _zero, _zero},
	/*FLOAT  */ {_neg1, _neg1, _neg1, _neg1, _neg1, ficmp, ffcmp, _neg1, _zero, _zero},
	/*BOOL   */ {_neg1, _neg1, _neg1, _neg1, _neg1, _pos1, _pos1, bbcmp, _zero, _zero},
	/*ARRAY  */ {_zero, _zero, _neg1, _zero, _zero, _zero, _zero, _zero, _zero, _zero},
	/*MAP    */ {_zero, _zero, _neg1, _zero, _zero, _zero, _zero, _zero, _zero, _zero},
}

func NumericAscendingComparator(input1 *Mlrval, input2 *Mlrval) int {
	return num_cmp_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}
func NumericDescendingComparator(input1 *Mlrval, input2 *Mlrval) int {
	return NumericAscendingComparator(input2, input1)
}
