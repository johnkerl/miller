package types

import (
	"fmt"
	"os"

	"miller/src/lib"
)

// ================================================================
func MlrvalTypeof(input1 *Mlrval) Mlrval {
	return MlrvalFromString(input1.GetTypeName())
}

// ----------------------------------------------------------------
func string_to_int(input1 *Mlrval) Mlrval {
	i, ok := lib.TryIntFromString(input1.printrep)
	if ok {
		return MlrvalFromInt(i)
	} else {
		return MlrvalFromError()
	}
}

func float_to_int(input1 *Mlrval) Mlrval {
	return MlrvalFromInt(int(input1.floatval))
}

func bool_to_int(input1 *Mlrval) Mlrval {
	if input1.boolval == true {
		return MlrvalFromInt(1)
	} else {
		return MlrvalFromInt(0)
	}
}

var to_int_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*VOID   */ _void1,
	/*STRING */ string_to_int,
	/*INT    */ _1u___,
	/*FLOAT  */ float_to_int,
	/*BOOL   */ bool_to_int,
	/*ARRAY  */ _erro1,
	/*MAP    */ _erro1,
}

func MlrvalToInt(input1 *Mlrval) Mlrval {
	return to_int_dispositions[input1.mvtype](input1)
}

// ----------------------------------------------------------------
func string_to_float(input1 *Mlrval) Mlrval {
	f, ok := lib.TryFloat64FromString(input1.printrep)
	if ok {
		return MlrvalFromFloat64(f)
	} else {
		return MlrvalFromError()
	}
}

func int_to_float(input1 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(input1.intval))
}

func bool_to_float(input1 *Mlrval) Mlrval {
	if input1.boolval == true {
		return MlrvalFromFloat64(1.0)
	} else {
		return MlrvalFromFloat64(0.0)
	}
}

var to_float_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*VOID   */ _void1,
	/*STRING */ string_to_float,
	/*INT    */ int_to_float,
	/*FLOAT  */ _1u___,
	/*BOOL   */ bool_to_float,
	/*ARRAY  */ _erro1,
	/*MAP    */ _erro1,
}

func MlrvalToFloat(input1 *Mlrval) Mlrval {
	return to_float_dispositions[input1.mvtype](input1)
}

// ----------------------------------------------------------------
func string_to_boolean(input1 *Mlrval) Mlrval {
	b, ok := lib.TryBoolFromBoolString(input1.printrep)
	if ok {
		return MlrvalFromBool(b)
	} else {
		return MlrvalFromError()
	}
}

func int_to_bool(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.intval != 0)
}

func float_to_bool(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.floatval != 0.0)
}

var to_boolean_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*VOID   */ _void1,
	/*STRING */ string_to_boolean,
	/*INT    */ int_to_bool,
	/*FLOAT  */ float_to_bool,
	/*BOOL   */ _1u___,
	/*ARRAY  */ _erro1,
	/*MAP    */ _erro1,
}

func MlrvalToBoolean(input1 *Mlrval) Mlrval {
	return to_boolean_dispositions[input1.mvtype](input1)
}

// ----------------------------------------------------------------
func MlrvalIsAbsent(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_ABSENT)
}
func MlrvalIsError(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_ERROR)
}
func MlrvalIsBool(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_BOOL)
}
func MlrvalIsBoolean(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_BOOL)
}
func MlrvalIsEmpty(input1 *Mlrval) Mlrval {
	if input1.mvtype == MT_VOID {
		return MlrvalFromTrue()
	}
	if input1.mvtype == MT_STRING {
		if input1.printrep == "" {
			return MlrvalFromTrue()
		}
	}
	return MlrvalFromFalse()
}
func MlrvalIsEmptyMap(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_MAP && input1.mapval.FieldCount == 0)
}
func MlrvalIsFloat(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_FLOAT)
}
func MlrvalIsInt(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_INT)
}
func MlrvalIsMap(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_MAP)
}
func MlrvalIsArray(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_ARRAY)
}
func MlrvalIsNonEmptyMap(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_MAP && input1.mapval.FieldCount != 0)
}
func MlrvalIsNotEmpty(input1 *Mlrval) Mlrval {
	if input1.mvtype == MT_VOID {
		return MlrvalFromFalse()
	}
	if input1.mvtype == MT_STRING {
		if input1.printrep == "" {
			return MlrvalFromFalse()
		}
	}
	return MlrvalFromTrue()
}
func MlrvalIsNotMap(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype != MT_MAP)
}
func MlrvalIsNotArray(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype != MT_ARRAY)
}
func MlrvalIsNotNull(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype != MT_ABSENT && input1.mvtype != MT_VOID)
}
func MlrvalIsNull(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_ABSENT || input1.mvtype == MT_VOID)
}
func MlrvalIsNumeric(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_INT || input1.mvtype == MT_FLOAT)
}
func MlrvalIsPresent(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype != MT_ABSENT)
}
func MlrvalIsString(input1 *Mlrval) Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_STRING || input1.mvtype == MT_VOID)
}

// ----------------------------------------------------------------
func assertingCommon(input1, check *Mlrval, description string, context *Context) Mlrval {
	if check.IsFalse() {
		// TODO: get context as in the C impl
		//fprintf(stderr, "%s: %s type-assertion failed at NR=%lld FNR=%lld FILENAME=%s\n",
		//MLR_GLOBALS.bargv0, pstate->desc, pvars->pctx->nr, pvars->pctx->fnr, pvars->pctx->filename);
		//exit(1);
		fmt.Fprintf(
			os.Stderr,
			"Miller: %s type-assertion failed at NR=%d FNR=%d FILENAME=%s\n",
			description,
			context.NR,
			context.FNR,
			context.FILENAME,
		)
		os.Exit(1)
	}
	return *input1
}

func MlrvalAssertingAbsent(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsAbsent(input1)
	return assertingCommon(input1, &check, "is_absent", context)
}
func MlrvalAssertingError(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsError(input1)
	return assertingCommon(input1, &check, "is_error", context)
}
func MlrvalAssertingBool(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsBool(input1)
	return assertingCommon(input1, &check, "is_bool", context)
}
func MlrvalAssertingBoolean(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsBoolean(input1)
	return assertingCommon(input1, &check, "is_boolean", context)
}
func MlrvalAssertingEmpty(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsEmpty(input1)
	return assertingCommon(input1, &check, "is_empty", context)
}
func MlrvalAssertingEmptyMap(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsEmptyMap(input1)
	return assertingCommon(input1, &check, "is_empty_map", context)
}
func MlrvalAssertingFloat(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsFloat(input1)
	return assertingCommon(input1, &check, "is_float", context)
}
func MlrvalAssertingInt(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsInt(input1)
	return assertingCommon(input1, &check, "is_int", context)
}
func MlrvalAssertingMap(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsMap(input1)
	return assertingCommon(input1, &check, "is_map", context)
}
func MlrvalAssertingArray(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsArray(input1)
	return assertingCommon(input1, &check, "is_array", context)
}
func MlrvalAssertingNonEmptyMap(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNonEmptyMap(input1)
	return assertingCommon(input1, &check, "is_non_empty_map", context)
}
func MlrvalAssertingNotEmpty(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNotEmpty(input1)
	return assertingCommon(input1, &check, "is_not_empty", context)
}
func MlrvalAssertingNotMap(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNotMap(input1)
	return assertingCommon(input1, &check, "is_not_map", context)
}
func MlrvalAssertingNotArray(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNotArray(input1)
	return assertingCommon(input1, &check, "is_not_array", context)
}
func MlrvalAssertingNotNull(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNotNull(input1)
	return assertingCommon(input1, &check, "is_not_null", context)
}
func MlrvalAssertingNull(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNull(input1)
	return assertingCommon(input1, &check, "is_null", context)
}
func MlrvalAssertingNumeric(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNumeric(input1)
	return assertingCommon(input1, &check, "is_numeric", context)
}
func MlrvalAssertingPresent(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsPresent(input1)
	return assertingCommon(input1, &check, "is_present", context)
}
func MlrvalAssertingString(input1 *Mlrval, context *Context) Mlrval {
	check := MlrvalIsString(input1)
	return assertingCommon(input1, &check, "is_string", context)
}
