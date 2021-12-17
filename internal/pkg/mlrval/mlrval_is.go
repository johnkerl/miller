package mlrval

// It's essential that we use mv.Type() not mv.mvtype since types are
// JIT-computed on first access for most data-file values. See type.go for more
// information.

func (mv *Mlrval) IsLegit() bool {
	return mv.Type() >= MT_VOID
}

func (mv *Mlrval) IsErrorOrAbsent() bool {
	t := mv.Type()
	return t == MT_ERROR || t == MT_ABSENT
}

func (mv *Mlrval) IsError() bool {
	return mv.Type() == MT_ERROR
}

func (mv *Mlrval) IsAbsent() bool {
	return mv.Type() == MT_ABSENT
}

func (mv *Mlrval) IsNull() bool {
	return mv.Type() == MT_NULL
}

func (mv *Mlrval) IsVoid() bool {
	return mv.Type() == MT_VOID
}

func (mv *Mlrval) IsErrorOrVoid() bool {
	t := mv.Type()
	return t == MT_ERROR || t == MT_VOID
}

// * Error is non-empty
// * Absent is non-empty (shouldn't have been assigned in the first place; error should be surfaced)
// * Void is empty
// * Empty string is empty
// * Int/float/bool/array/map are all non-empty
func (mv *Mlrval) IsEmptyString() bool {
	t := mv.Type()
	if t == MT_VOID {
		return true
	} else if t == MT_STRING {
		return mv.printrep == ""
	} else {
		return false
	}
}

func (mv *Mlrval) IsString() bool {
	return mv.Type() == MT_STRING
}

func (mv *Mlrval) IsStringOrVoid() bool {
	t := mv.Type()
	return t == MT_STRING || t == MT_VOID
}

func (mv *Mlrval) IsInt() bool {
	return mv.Type() == MT_INT
}

func (mv *Mlrval) IsFloat() bool {
	return mv.Type() == MT_FLOAT
}

func (mv *Mlrval) IsNumeric() bool {
	t := mv.Type()
	return t == MT_INT || t == MT_FLOAT
}

func (mv *Mlrval) IsIntZero() bool {
	return mv.Type() == MT_INT && mv.intval == 0
}

func (mv *Mlrval) IsBool() bool {
	return mv.Type() == MT_BOOL
}

func (mv *Mlrval) IsTrue() bool {
	return mv.Type() == MT_BOOL && mv.boolval == true
}
func (mv *Mlrval) IsFalse() bool {
	return mv.Type() == MT_BOOL && mv.boolval == false
}

func (mv *Mlrval) IsArray() bool {
	// TODO: comment non-deferrable type -- don't force a (potentially
	// expensive in bulk) JIT-infer of other types
	// return mv.Type() == MT_ARRAY
	return mv.mvtype == MT_ARRAY
}
func (mv *Mlrval) IsMap() bool {
	// TODO: comment non-deferrable type -- don't force a (potentially
	// expensive in bulk) JIT-infer of other types
	// return mv.Type() == MT_ARRAY
	return mv.mvtype == MT_MAP
}
func (mv *Mlrval) IsArrayOrMap() bool {
	// TODO: comment why not
	// In flatten we don't want to type-infer things that don't need to be jitted.
	// Arrays & maps are never from deferred type.
	// t := mv.Type()
	t := mv.mvtype
	return t == MT_ARRAY || t == MT_MAP
}

func (mv *Mlrval) IsFunction() bool {
	return mv.Type() == MT_FUNC
}
