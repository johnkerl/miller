package types

import (
	"fmt"
	"os"

	"mlr/src/lib"
)

// ================================================================
func MlrvalTypeof(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromString(input1.GetTypeName())
}

// ----------------------------------------------------------------
func string_to_int(input1 *Mlrval) *Mlrval {
	i, ok := lib.TryIntFromString(input1.printrep)
	if ok {
		return MlrvalPointerFromInt(i)
	} else {
		return MLRVAL_ERROR
	}
}

func float_to_int(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromInt(int(input1.floatval))
}

func bool_to_int(input1 *Mlrval) *Mlrval {
	if input1.boolval == true {
		return MlrvalPointerFromInt(1)
	} else {
		return MlrvalPointerFromInt(0)
	}
}

var to_int_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*NULL   */ _null1,
	/*VOID   */ _void1,
	/*STRING */ string_to_int,
	/*INT    */ _1u___,
	/*FLOAT  */ float_to_int,
	/*BOOL   */ bool_to_int,
	/*ARRAY  */ _erro1,
	/*MAP    */ _erro1,
}

func MlrvalToInt(input1 *Mlrval) *Mlrval {
	return to_int_dispositions[input1.mvtype](input1)
}

// ----------------------------------------------------------------
func string_to_float(input1 *Mlrval) *Mlrval {
	f, ok := lib.TryFloat64FromString(input1.printrep)
	if ok {
		return MlrvalPointerFromFloat64(f)
	} else {
		return MLRVAL_ERROR
	}
}

func int_to_float(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromFloat64(float64(input1.intval))
}

func bool_to_float(input1 *Mlrval) *Mlrval {
	if input1.boolval == true {
		return MlrvalPointerFromFloat64(1.0)
	} else {
		return MlrvalPointerFromFloat64(0.0)
	}
}

var to_float_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*NULL   */ _null1,
	/*VOID   */ _void1,
	/*STRING */ string_to_float,
	/*INT    */ int_to_float,
	/*FLOAT  */ _1u___,
	/*BOOL   */ bool_to_float,
	/*ARRAY  */ _erro1,
	/*MAP    */ _erro1,
}

func MlrvalToFloat(input1 *Mlrval) *Mlrval {
	return to_float_dispositions[input1.mvtype](input1)
}

// ----------------------------------------------------------------
func string_to_boolean(input1 *Mlrval) *Mlrval {
	b, ok := lib.TryBoolFromBoolString(input1.printrep)
	if ok {
		return MlrvalPointerFromBool(b)
	} else {
		return MLRVAL_ERROR
	}
}

func int_to_bool(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.intval != 0)
}

func float_to_bool(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.floatval != 0.0)
}

var to_boolean_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*NULL   */ _null1,
	/*VOID   */ _void1,
	/*STRING */ string_to_boolean,
	/*INT    */ int_to_bool,
	/*FLOAT  */ float_to_bool,
	/*BOOL   */ _1u___,
	/*ARRAY  */ _erro1,
	/*MAP    */ _erro1,
}

func MlrvalToBoolean(input1 *Mlrval) *Mlrval {
	return to_boolean_dispositions[input1.mvtype](input1)
}

// ----------------------------------------------------------------
func MlrvalIsAbsent(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype == MT_ABSENT)
}
func MlrvalIsError(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype == MT_ERROR)
}
func MlrvalIsBool(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype == MT_BOOL)
}
func MlrvalIsBoolean(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype == MT_BOOL)
}
func MlrvalIsEmpty(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_VOID {
		return MLRVAL_TRUE
	} else if input1.mvtype == MT_STRING {
		if input1.printrep == "" {
			return MLRVAL_TRUE
		} else {
			return MLRVAL_FALSE
		}
	} else {
		return MLRVAL_FALSE
	}
}
func MlrvalIsEmptyMap(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype == MT_MAP && input1.mapval.FieldCount == 0)
}
func MlrvalIsFloat(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype == MT_FLOAT)
}
func MlrvalIsInt(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype == MT_INT)
}
func MlrvalIsMap(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype == MT_MAP)
}
func MlrvalIsArray(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype == MT_ARRAY)
}
func MlrvalIsNonEmptyMap(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype == MT_MAP && input1.mapval.FieldCount != 0)
}
func MlrvalIsNotEmpty(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_VOID {
		return MLRVAL_FALSE
	} else if input1.mvtype == MT_STRING {
		if input1.printrep == "" {
			return MLRVAL_FALSE
		} else {
			return MLRVAL_TRUE
		}
	} else {
		return MLRVAL_TRUE
	}
}
func MlrvalIsNotMap(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype != MT_MAP)
}
func MlrvalIsNotArray(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype != MT_ARRAY)
}
func MlrvalIsNotNull(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype != MT_ABSENT && input1.mvtype != MT_VOID)
}
func MlrvalIsNull(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype == MT_ABSENT || input1.mvtype == MT_VOID)
}
func MlrvalIsNumeric(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype == MT_INT || input1.mvtype == MT_FLOAT)
}
func MlrvalIsPresent(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype != MT_ABSENT)
}
func MlrvalIsString(input1 *Mlrval) *Mlrval {
	return MlrvalPointerFromBool(input1.mvtype == MT_STRING || input1.mvtype == MT_VOID)
}

// ----------------------------------------------------------------
func assertingCommon(input1, check *Mlrval, description string, context *Context) *Mlrval {
	if check.IsFalse() {
		// TODO: get context as in the C impl
		//fprintf(stderr, "%s: %s type-assertion failed at NR=%lld FNR=%lld FILENAME=%s\n",
		//MLR_GLOBALS.bargv0, pstate->desc, pvars->pctx->nr, pvars->pctx->fnr, pvars->pctx->filename);
		//exit(1);
		fmt.Fprintf(
			os.Stderr,
			"mlr: %s type-assertion failed at NR=%d FNR=%d FILENAME=%s\n",
			description,
			context.NR,
			context.FNR,
			context.FILENAME,
		)
		os.Exit(1)
	}
	return input1
}

func MlrvalAssertingAbsent(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsAbsent(input1), "is_absent", context)
}
func MlrvalAssertingError(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsError(input1), "is_error", context)
}
func MlrvalAssertingBool(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsBool(input1), "is_bool", context)
}
func MlrvalAssertingBoolean(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsBoolean(input1), "is_boolean", context)
}
func MlrvalAssertingEmpty(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsEmpty(input1), "is_empty", context)
}
func MlrvalAssertingEmptyMap(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsEmptyMap(input1), "is_empty_map", context)
}
func MlrvalAssertingFloat(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsFloat(input1), "is_float", context)
}
func MlrvalAssertingInt(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsInt(input1), "is_int", context)
}
func MlrvalAssertingMap(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsMap(input1), "is_map", context)
}
func MlrvalAssertingArray(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsArray(input1), "is_array", context)
}
func MlrvalAssertingNonEmptyMap(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsNonEmptyMap(input1), "is_non_empty_map", context)
}
func MlrvalAssertingNotEmpty(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsNotEmpty(input1), "is_not_empty", context)
}
func MlrvalAssertingNotMap(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsNotMap(input1), "is_not_map", context)
}
func MlrvalAssertingNotArray(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsNotArray(input1), "is_not_array", context)
}
func MlrvalAssertingNotNull(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsNotNull(input1), "is_not_null", context)
}
func MlrvalAssertingNull(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsNull(input1), "is_null", context)
}
func MlrvalAssertingNumeric(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsNumeric(input1), "is_numeric", context)
}
func MlrvalAssertingPresent(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsPresent(input1), "is_present", context)
}
func MlrvalAssertingString(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, MlrvalIsString(input1), "is_string", context)
}
