// ================================================================
// For sorting
// ================================================================

package types

// Lexical sort: just stringify everything.
func LexicalAscendingComparator(ma *Mlrval, mb *Mlrval) int {
	sa := ma.String()
	sb := mb.String()
	if sa < sb {
		return -1
	} else if sa > sb {
		return 1
	} else {
		return 0
	}
}
func LexicalDescendingComparator(ma *Mlrval, mb *Mlrval) int {
	return LexicalAscendingComparator(mb, ma)
}

// ----------------------------------------------------------------
// Sort rules (same for min, max, and comparator):
// * NUMERICS < BOOL < STRINGS < ERROR < ABSENT
// * error == error (singleton type)
// * absent == absent (singleton type)
// * string compares on strings
// * numeric compares on numbers
// * false < true

func _neg1(ma, mb *Mlrval) int {
	return -1
}
func _zero(ma, mb *Mlrval) int {
	return 0
}
func _pos1(ma, mb *Mlrval) int {
	return 1
}

func _scmp(ma, mb *Mlrval) int {
	if ma.printrep < mb.printrep {
		return -1
	} else if ma.printrep > mb.printrep {
		return 1
	} else {
		return 0
	}
}

func iicmp(ma, mb *Mlrval) int {
	ca := ma.intval
	cb := mb.intval
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}
func ifcmp(ma, mb *Mlrval) int {
	ca := float64(ma.intval)
	cb := mb.floatval
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}
func ficmp(ma, mb *Mlrval) int {
	ca := ma.floatval
	cb := float64(mb.intval)
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}
func ffcmp(ma, mb *Mlrval) int {
	ca := ma.floatval
	cb := mb.floatval
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}

func bbcmp(ma, mb *Mlrval) int {
	a := ma.boolval
	b := mb.boolval
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
// * NUMERICS < BOOL < STRINGS < ERROR < ABSENT
// * error == error (singleton type)
// * absent == absent (singleton type)
// * string compares on strings
// * numeric compares on numbers
// * false < true

var num_cmp_dispositions = [MT_DIM][MT_DIM]ComparatorFunc{
	//       .  ERROR   ABSENT VOID   STRING INT    FLOAT  BOOL    ARRAY MAP
	/*ERROR  */ {_zero, _neg1, _pos1, _pos1, _pos1, _pos1, _pos1, _zero, _zero},
	/*ABSENT */ {_pos1, _zero, _pos1, _pos1, _pos1, _pos1, _pos1, _zero, _zero},
	/*VOID   */ {_neg1, _neg1, _scmp, _scmp, _pos1, _pos1, _pos1, _zero, _zero},
	/*STRING */ {_neg1, _neg1, _scmp, _scmp, _pos1, _pos1, _pos1, _zero, _zero},
	/*INT    */ {_neg1, _neg1, _neg1, _neg1, iicmp, ifcmp, _neg1, _zero, _zero},
	/*FLOAT  */ {_neg1, _neg1, _neg1, _neg1, ficmp, ffcmp, _neg1, _zero, _zero},
	/*BOOL   */ {_neg1, _neg1, _neg1, _neg1, _pos1, _pos1, bbcmp, _zero, _zero},
	/*ARRAY  */ {_zero, _zero, _zero, _zero, _zero, _zero, _zero, _zero, _zero},
	/*MAP    */ {_zero, _zero, _zero, _zero, _zero, _zero, _zero, _zero, _zero},
}

func NumericAscendingComparator(ma *Mlrval, mb *Mlrval) int {
	return num_cmp_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}
func NumericDescendingComparator(ma *Mlrval, mb *Mlrval) int {
	return NumericAscendingComparator(mb, ma)
}
