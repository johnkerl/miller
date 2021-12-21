package mlrval

import (
	"github.com/johnkerl/miller/internal/pkg/lib"
)

// It's essential that we use mv.Type() not mv.mvtype since types are
// JIT-computed on first access for most data-file values. See type.go for more
// information.

func (mv *Mlrval) IsLegit() bool {
	return mv.Type() >= MT_VOID
}

// TODO: comment no JIT-infer here -- absent is non-inferrable and we needn't take the expense of JIT.
func (mv *Mlrval) IsErrorOrAbsent() bool {
	t := mv.mvtype
	return t == MT_ERROR || t == MT_ABSENT
}

func (mv *Mlrval) IsError() bool {
	return mv.Type() == MT_ERROR
}

// TODO: comment no JIT-infer here -- absent is non-inferrable and we needn't take the expense of JIT.
func (mv *Mlrval) IsAbsent() bool {
	return mv.mvtype == MT_ABSENT
}

// TODO: comment no JIT-infer here -- NULL is non-inferrable and we needn't take the expense of JIT.
// This is a literal in JSON files, or else explicitly set to NULL.
func (mv *Mlrval) IsNull() bool {
	return mv.mvtype == MT_NULL
}

func (mv *Mlrval) IsVoid() bool {
	if mv.mvtype == MT_VOID {
		return true
	}
	if mv.mvtype == MT_PENDING && mv.printrep == "" {
		lib.InternalCodingErrorIf(!mv.printrepValid)
		return true
	}
	return false
}

func (mv *Mlrval) IsErrorOrVoid() bool {
	return mv.IsError() || mv.IsVoid()
}

// * Error is non-empty
// * Absent is non-empty (shouldn't have been assigned in the first place; error should be surfaced)
// * Void is empty
// * Empty string is empty
// * Int/float/bool/array/map are all non-empty
func (mv *Mlrval) IsEmptyString() bool {
	if mv.mvtype == MT_VOID {
		return true
	}
	if mv.mvtype == MT_STRING && mv.printrep == "" {
		return true
	}
	if mv.mvtype == MT_PENDING && mv.printrep == "" {
		lib.InternalCodingErrorIf(!mv.printrepValid)
		return true
	}
	return false
}

func (mv *Mlrval) IsString() bool {
	return mv.Type() == MT_STRING
}

func (mv *Mlrval) IsStringOrVoid() bool {
	t := mv.Type()
	return t == MT_STRING || t == MT_VOID
}

func (mv *Mlrval) IsStringOrInt() bool {
	t := mv.Type()
	return t == MT_STRING || t == MT_VOID || t == MT_INT
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
	return mv.mvtype == MT_FUNC
}
