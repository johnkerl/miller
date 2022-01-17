package mlrval

import (
	"fmt"
	"os"

	"github.com/johnkerl/miller/internal/pkg/lib"
)

// It's essential that we use mv.Type() not mv.mvtype, or use an Is...()
// predicate, or another Get...(), since types are JIT-computed on first access
// for most data-file values. See type.go for more information.

func (mv *Mlrval) GetTypeBit() int {
	return 1 << mv.Type()
}

func (mv *Mlrval) GetStringValue() (stringValue string, isString bool) {
	if mv.Type() == MT_STRING || mv.Type() == MT_VOID {
		return mv.printrep, true
	} else {
		return "", false
	}
}

func (mv *Mlrval) GetIntValue() (intValue int, isInt bool) {
	if mv.Type() == MT_INT {
		return mv.intval, true
	} else {
		return -999, false
	}
}

func (mv *Mlrval) GetFloatValue() (floatValue float64, isFloat bool) {
	if mv.Type() == MT_FLOAT {
		return mv.floatval, true
	} else {
		return -777.0, false
	}
}

func (mv *Mlrval) GetNumericToFloatValue() (floatValue float64, isFloat bool) {
	if mv.Type() == MT_FLOAT {
		return mv.floatval, true
	} else if mv.Type() == MT_INT {
		return float64(mv.intval), true
	} else {
		return -888.0, false
	}
}

func (mv *Mlrval) GetNumericNegativeorDie() bool {
	floatValue, ok := mv.GetNumericToFloatValue()
	lib.InternalCodingErrorIf(!ok)
	return floatValue < 0.0
}

func (mv *Mlrval) GetBoolValue() (boolValue bool, isBool bool) {
	if mv.Type() == MT_BOOL {
		return mv.boolval, true
	} else {
		return false, false
	}
}

func (mv *Mlrval) GetArray() []*Mlrval {
	if mv.IsArray() {
		return mv.arrayval
	} else {
		return nil
	}
}

func (mv *Mlrval) GetMap() *Mlrmap {
	if mv.IsMap() {
		return mv.mapval
	} else {
		return nil
	}
}

func (mv *Mlrval) GetFunction() interface{} {
	if mv.Type() == MT_FUNC {
		return mv.funcval
	} else {
		return nil
	}
}

func (mv *Mlrval) GetTypeName() string {
	return TYPE_NAMES[mv.Type()]
}

func GetTypeName(mvtype MVType) string {
	return TYPE_NAMES[mvtype]
}

// These are for built-in functions operating within type-keyed
// disposition-vector/disposition-matrix context. They've already computed
// mv.Type() -- it's a fatal error if they haven't -- and they need the typed
// value.

func (mv *Mlrval) AcquireStringValue() string {
	lib.InternalCodingErrorIf(mv.mvtype != MT_STRING && mv.mvtype != MT_VOID)
	return mv.printrep
}

func (mv *Mlrval) AcquireIntValue() int {
	lib.InternalCodingErrorIf(mv.mvtype != MT_INT)
	return mv.intval
}

func (mv *Mlrval) AcquireFloatValue() float64 {
	lib.InternalCodingErrorIf(mv.mvtype != MT_FLOAT)
	return mv.floatval
}

func (mv *Mlrval) AcquireBoolValue() bool {
	lib.InternalCodingErrorIf(mv.mvtype != MT_BOOL)
	return mv.boolval
}

func (mv *Mlrval) AcquireArrayValue() []*Mlrval {
	lib.InternalCodingErrorIf(mv.mvtype != MT_ARRAY)
	return mv.arrayval
}

func (mv *Mlrval) AcquireMapValue() *Mlrmap {
	lib.InternalCodingErrorIf(mv.mvtype != MT_MAP)
	return mv.mapval
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

func (mv *Mlrval) AssertNumeric() {
	_ = mv.GetNumericToFloatValueOrDie()
}
