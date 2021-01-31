package types

import (
	"strconv"

	"miller/lib"
)

// ----------------------------------------------------------------
func (this *Mlrval) GetType() MVType {
	return this.mvtype
}

func (this *Mlrval) GetTypeName() string {
	return TYPE_NAMES[this.mvtype]
}

func GetTypeName(mvtype MVType) string {
	return TYPE_NAMES[mvtype]
}

// ----------------------------------------------------------------
func (this *Mlrval) IsLegit() bool {
	return this.mvtype > MT_VOID
}

func (this *Mlrval) IsErrorOrAbsent() bool {
	return this.mvtype == MT_ERROR || this.mvtype == MT_ABSENT
}

func (this *Mlrval) IsError() bool {
	return this.mvtype == MT_ERROR
}

func (this *Mlrval) IsAbsent() bool {
	return this.mvtype == MT_ABSENT
}

func (this *Mlrval) IsVoid() bool {
	return this.mvtype == MT_VOID
}

func (this *Mlrval) IsErrorOrVoid() bool {
	return this.mvtype == MT_ERROR || this.mvtype == MT_VOID
}

// * Error is non-empty
// * Absent is non-empty (shouldn't have been assigned in the first place; error
//   should be surfaced)
// * Void is empty
// * Empty string is empty
// * Int/float/bool/array/map are all non-empty
func (this *Mlrval) IsEmpty() bool {
	if this.mvtype == MT_VOID {
		return true
	} else if this.mvtype == MT_STRING {
		return this.printrep == ""
	} else {
		return false
	}
}

func (this *Mlrval) IsString() bool {
	return this.mvtype == MT_STRING
}

func (this *Mlrval) IsStringOrVoid() bool {
	return this.mvtype == MT_STRING || this.mvtype == MT_VOID
}

func (this *Mlrval) IsInt() bool {
	return this.mvtype == MT_INT
}

func (this *Mlrval) IsFloat() bool {
	return this.mvtype == MT_FLOAT
}

func (this *Mlrval) IsNumeric() bool {
	return this.mvtype == MT_INT || this.mvtype == MT_FLOAT
}

func (this *Mlrval) IsIntZero() bool {
	return this.mvtype == MT_INT && this.intval == 0
}

func (this *Mlrval) IsBool() bool {
	return this.mvtype == MT_BOOL
}

func (this *Mlrval) IsTrue() bool {
	return this.mvtype == MT_BOOL && this.boolval == true
}
func (this *Mlrval) IsFalse() bool {
	return this.mvtype == MT_BOOL && this.boolval == false
}

func (this *Mlrval) IsArray() bool {
	return this.mvtype == MT_ARRAY
}
func (this *Mlrval) IsMap() bool {
	return this.mvtype == MT_MAP
}
func (this *Mlrval) IsArrayOrMap() bool {
	return this.mvtype == MT_ARRAY || this.mvtype == MT_MAP
}

// ----------------------------------------------------------------
func (this *Mlrval) GetIntValue() (intValue int, isInt bool) {
	if this.mvtype == MT_INT {
		return this.intval, true
	} else {
		return -999, false
	}
}

func (this *Mlrval) GetFloatValue() (floatValue float64, isFloat bool) {
	if this.mvtype == MT_FLOAT {
		return this.floatval, true
	} else {
		return -777.0, false
	}
}

func (this *Mlrval) GetNumericToFloatValue() (floatValue float64, isFloat bool) {
	if this.mvtype == MT_FLOAT {
		return this.floatval, true
	} else if this.mvtype == MT_INT {
		return float64(this.intval), true
	} else {
		return -888.0, false
	}
}

func (this *Mlrval) GetBoolValue() (boolValue bool, isBool bool) {
	if this.mvtype == MT_BOOL {
		return this.boolval, true
	} else {
		return false, false
	}
}

func (this *Mlrval) GetArray() []Mlrval {
	if this.mvtype == MT_ARRAY {
		return this.arrayval
	} else {
		return nil
	}
}

func (this *Mlrval) GetArrayLength() (int, bool) {
	if this.mvtype == MT_ARRAY {
		return len(this.arrayval), true
	} else {
		return -999, false
	}
}

func (this *Mlrval) GetMap() *Mlrmap {
	if this.mvtype == MT_MAP {
		return this.mapval
	} else {
		return nil
	}
}

// ----------------------------------------------------------------
func (this *Mlrval) Copy() *Mlrval {
	that := *this
	if this.mvtype == MT_MAP {
		that.mapval = this.mapval.Copy()
	} else if this.mvtype == MT_ARRAY {
		that.arrayval = CopyMlrvalArray(this.arrayval)
	}
	return &that
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

func (this *Mlrval) FlattenToMap(prefix string, delimiter string) Mlrval {
	retval := NewMlrmap()

	if this.mvtype == MT_MAP {
		// Without this, the for-loop below is zero-pass and fields with "{}"
		// values would disappear entirely in a JSON-to-CSV conversion.
		if this.mapval.FieldCount == 0 {
			if prefix != "" {
				retval.PutCopy(prefix, MlrvalPointerFromString("{}"))
			}
		}

		for pe := this.mapval.Head; pe != nil; pe = pe.Next {
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

	} else if this.mvtype == MT_ARRAY {
		// Without this, the for-loop below is zero-pass and fields with "[]"
		// values would disappear entirely in a JSON-to-CSV conversion.
		if len(this.arrayval) == 0 {
			if prefix != "" {
				retval.PutCopy(prefix, MlrvalPointerFromString("[]"))
			}
		}

		for zindex, value := range this.arrayval {
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
		retval.PutCopy(prefix, this.Copy())
	}

	return MlrvalFromMapReferenced(retval)
}

// ----------------------------------------------------------------
// Used by stats1.

func (this *Mlrval) Increment() {
	if this.mvtype == MT_INT {
		this.intval++
	} else if this.mvtype == MT_FLOAT {
		this.floatval++
	}
}
