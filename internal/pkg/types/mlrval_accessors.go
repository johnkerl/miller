package types

import (
	"fmt"
	"os"
	"strconv"

	"mlr/internal/pkg/lib"
)

// ----------------------------------------------------------------
func (mv *Mlrval) GetType() MVType {
	return mv.mvtype
}

func (mv *Mlrval) GetTypeName() string {
	return TYPE_NAMES[mv.mvtype]
}

func GetTypeName(mvtype MVType) string {
	return TYPE_NAMES[mvtype]
}

// ----------------------------------------------------------------
func (mv *Mlrval) IsLegit() bool {
	return mv.mvtype >= MT_VOID
}

func (mv *Mlrval) IsErrorOrAbsent() bool {
	return mv.mvtype == MT_ERROR || mv.mvtype == MT_ABSENT
}

func (mv *Mlrval) IsError() bool {
	return mv.mvtype == MT_ERROR
}

func (mv *Mlrval) IsAbsent() bool {
	return mv.mvtype == MT_ABSENT
}

func (mv *Mlrval) IsNull() bool {
	return mv.mvtype == MT_NULL
}

func (mv *Mlrval) IsVoid() bool {
	return mv.mvtype == MT_VOID
}

func (mv *Mlrval) IsErrorOrVoid() bool {
	return mv.mvtype == MT_ERROR || mv.mvtype == MT_VOID
}

// * Error is non-empty
// * Absent is non-empty (shouldn't have been assigned in the first place; error
//   should be surfaced)
// * Void is empty
// * Empty string is empty
// * Int/float/bool/array/map are all non-empty
func (mv *Mlrval) IsEmpty() bool {
	if mv.mvtype == MT_VOID {
		return true
	} else if mv.mvtype == MT_STRING {
		return mv.printrep == ""
	} else {
		return false
	}
}

func (mv *Mlrval) IsString() bool {
	return mv.mvtype == MT_STRING
}

func (mv *Mlrval) IsStringOrVoid() bool {
	return mv.mvtype == MT_STRING || mv.mvtype == MT_VOID
}

func (mv *Mlrval) IsInt() bool {
	return mv.mvtype == MT_INT
}

func (mv *Mlrval) IsFloat() bool {
	return mv.mvtype == MT_FLOAT
}

func (mv *Mlrval) IsNumeric() bool {
	return mv.mvtype == MT_INT || mv.mvtype == MT_FLOAT
}

func (mv *Mlrval) IsIntZero() bool {
	return mv.mvtype == MT_INT && mv.intval == 0
}

func (mv *Mlrval) IsBool() bool {
	return mv.mvtype == MT_BOOL
}

func (mv *Mlrval) IsTrue() bool {
	return mv.mvtype == MT_BOOL && mv.boolval == true
}
func (mv *Mlrval) IsFalse() bool {
	return mv.mvtype == MT_BOOL && mv.boolval == false
}

func (mv *Mlrval) IsArray() bool {
	return mv.mvtype == MT_ARRAY
}
func (mv *Mlrval) IsMap() bool {
	return mv.mvtype == MT_MAP
}
func (mv *Mlrval) IsArrayOrMap() bool {
	return mv.mvtype == MT_ARRAY || mv.mvtype == MT_MAP
}

func (mv *Mlrval) IsFunction() bool {
	return mv.mvtype == MT_FUNC
}

// ----------------------------------------------------------------
func (mv *Mlrval) GetIntValue() (intValue int, isInt bool) {
	if mv.mvtype == MT_INT {
		return mv.intval, true
	} else {
		return -999, false
	}
}

func (mv *Mlrval) GetFloatValue() (floatValue float64, isFloat bool) {
	if mv.mvtype == MT_FLOAT {
		return mv.floatval, true
	} else {
		return -777.0, false
	}
}

func (mv *Mlrval) GetNumericToFloatValue() (floatValue float64, isFloat bool) {
	if mv.mvtype == MT_FLOAT {
		return mv.floatval, true
	} else if mv.mvtype == MT_INT {
		return float64(mv.intval), true
	} else {
		return -888.0, false
	}
}

func (mv *Mlrval) GetNumericToFloatValueOrDie() (floatValue float64) {
	floatValue, ok := mv.GetNumericToFloatValue()
	if !ok {
		fmt.Fprintf(
			os.Stderr,
			"%s: couldn't parse \"%s\" as number.",
			"mlr", mv.String(),
		)
		os.Exit(1)
	}
	return floatValue
}

func (mv *Mlrval) GetNumericNegativeorDie() bool {
	floatValue, ok := mv.GetNumericToFloatValue()
	lib.InternalCodingErrorIf(!ok)
	return floatValue < 0.0
}

func (mv *Mlrval) AssertNumeric() {
	_ = mv.GetNumericToFloatValueOrDie()
}

func (mv *Mlrval) GetString() (stringValue string, isString bool) {
	if mv.mvtype == MT_STRING || mv.mvtype == MT_VOID {
		return mv.printrep, true
	} else {
		return "", false
	}
}

func (mv *Mlrval) GetBoolValue() (boolValue bool, isBool bool) {
	if mv.mvtype == MT_BOOL {
		return mv.boolval, true
	} else {
		return false, false
	}
}

func (mv *Mlrval) GetArray() []Mlrval {
	if mv.mvtype == MT_ARRAY {
		return mv.arrayval
	} else {
		return nil
	}
}

func (mv *Mlrval) GetArrayLength() (int, bool) {
	if mv.mvtype == MT_ARRAY {
		return len(mv.arrayval), true
	} else {
		return -999, false
	}
}

func (mv *Mlrval) GetMap() *Mlrmap {
	if mv.mvtype == MT_MAP {
		return mv.mapval
	} else {
		return nil
	}
}

func (mv *Mlrval) GetFunction() interface{} {
	if mv.mvtype == MT_FUNC {
		return mv.funcval
	} else {
		return nil
	}
}

// ----------------------------------------------------------------
func (mv *Mlrval) Copy() *Mlrval {
	other := *mv
	if mv.mvtype == MT_MAP {
		other.mapval = mv.mapval.Copy()
	} else if mv.mvtype == MT_ARRAY {
		other.arrayval = CopyMlrvalArray(mv.arrayval)
	}
	return &other
}

func CopyMlrvalArray(input []Mlrval) []Mlrval {
	output := make([]Mlrval, len(input))
	for i, element := range input {
		output[i] = *element.Copy()
	}
	return output
}

func CopyMlrvalPointerArray(input []*Mlrval) []*Mlrval {
	output := make([]*Mlrval, len(input))
	for i, element := range input {
		output[i] = element.Copy()
	}
	return output
}

// ---------------------------------------------------------------
// For the flatten verb and DSL function.

func (mv *Mlrval) FlattenToMap(prefix string, delimiter string) Mlrval {
	retval := NewMlrmap()

	if mv.mvtype == MT_MAP {
		// Without this, the for-loop below is zero-pass and fields with "{}"
		// values would disappear entirely in a JSON-to-CSV conversion.
		if mv.mapval.IsEmpty() {
			if prefix != "" {
				retval.PutCopy(prefix, MlrvalFromString("{}"))
			}
		}

		for pe := mv.mapval.Head; pe != nil; pe = pe.Next {
			nextPrefix := pe.Key
			if prefix != "" {
				nextPrefix = prefix + delimiter + nextPrefix
			}
			if pe.Value.mvtype == MT_MAP || pe.Value.mvtype == MT_ARRAY {
				nextResult := pe.Value.FlattenToMap(nextPrefix, delimiter)
				lib.InternalCodingErrorIf(nextResult.mvtype != MT_MAP)
				for pf := nextResult.mapval.Head; pf != nil; pf = pf.Next {
					retval.PutCopy(pf.Key, pf.Value.Copy())
				}
			} else {
				retval.PutCopy(nextPrefix, pe.Value.Copy())
			}
		}

	} else if mv.mvtype == MT_ARRAY {
		// Without this, the for-loop below is zero-pass and fields with "[]"
		// values would disappear entirely in a JSON-to-CSV conversion.
		if len(mv.arrayval) == 0 {
			if prefix != "" {
				retval.PutCopy(prefix, MlrvalFromString("[]"))
			}
		}

		for zindex, value := range mv.arrayval {
			nextPrefix := strconv.Itoa(zindex + 1) // Miller user-space indices are 1-up
			if prefix != "" {
				nextPrefix = prefix + delimiter + nextPrefix
			}
			if value.mvtype == MT_MAP || value.mvtype == MT_ARRAY {
				nextResult := value.FlattenToMap(nextPrefix, delimiter)
				lib.InternalCodingErrorIf(nextResult.mvtype != MT_MAP)
				for pf := nextResult.mapval.Head; pf != nil; pf = pf.Next {
					retval.PutCopy(pf.Key, pf.Value.Copy())
				}
			} else {
				retval.PutCopy(nextPrefix, value.Copy())
			}
		}

	} else {
		retval.PutCopy(prefix, mv.Copy())
	}

	return *MlrvalFromMapReferenced(retval)
}

// ----------------------------------------------------------------
// Used by stats1.

func (mv *Mlrval) Increment() {
	if mv.mvtype == MT_INT {
		mv.intval++
	} else if mv.mvtype == MT_FLOAT {
		mv.floatval++
	}
}
