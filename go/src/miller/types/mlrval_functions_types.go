package types

import (
	"fmt"
	"os"

	"miller/lib"
)

// ================================================================
func MlrvalTypeof(ma *Mlrval) Mlrval {
	return MlrvalFromString(ma.GetTypeName())
}

// ----------------------------------------------------------------
func string_to_int(ma *Mlrval) Mlrval {
	i, ok := lib.TryInt64FromString(ma.printrep)
	if ok {
		return MlrvalFromInt64(i)
	} else {
		return MlrvalFromError()
	}
}

func float_to_int(ma *Mlrval) Mlrval {
	return MlrvalFromInt64(int64(ma.floatval))
}

func bool_to_int(ma *Mlrval) Mlrval {
	if ma.boolval == true {
		return MlrvalFromInt64(1)
	} else {
		return MlrvalFromInt64(0)
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

func MlrvalToInt(ma *Mlrval) Mlrval {
	return to_int_dispositions[ma.mvtype](ma)
}

// ----------------------------------------------------------------
func string_to_float(ma *Mlrval) Mlrval {
	f, ok := lib.TryFloat64FromString(ma.printrep)
	if ok {
		return MlrvalFromFloat64(f)
	} else {
		return MlrvalFromError()
	}
}

func int_to_float(ma *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(ma.intval))
}

func bool_to_float(ma *Mlrval) Mlrval {
	if ma.boolval == true {
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

func MlrvalToFloat(ma *Mlrval) Mlrval {
	return to_float_dispositions[ma.mvtype](ma)
}

// ----------------------------------------------------------------
func string_to_boolean(ma *Mlrval) Mlrval {
	b, ok := lib.TryBoolFromBoolString(ma.printrep)
	if ok {
		return MlrvalFromBool(b)
	} else {
		return MlrvalFromError()
	}
}

func int_to_bool(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval != 0)
}

func float_to_bool(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval != 0.0)
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

func MlrvalToBoolean(ma *Mlrval) Mlrval {
	return to_boolean_dispositions[ma.mvtype](ma)
}

// ----------------------------------------------------------------
func MlrvalIsAbsent(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype == MT_ABSENT)
}
func MlrvalIsError(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype == MT_ERROR)
}
func MlrvalIsBool(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype == MT_BOOL)
}
func MlrvalIsBoolean(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype == MT_BOOL)
}
func MlrvalIsEmpty(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_VOID {
		return MlrvalFromTrue()
	}
	if ma.mvtype == MT_STRING {
		if ma.printrep == "" {
			return MlrvalFromTrue()
		}
	}
	return MlrvalFromFalse()
}
func MlrvalIsEmptyMap(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype == MT_MAP && ma.mapval.FieldCount == 0)
}
func MlrvalIsFloat(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype == MT_FLOAT)
}
func MlrvalIsInt(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype == MT_INT)
}
func MlrvalIsMap(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype == MT_MAP)
}
func MlrvalIsArray(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype == MT_ARRAY)
}
func MlrvalIsNonEmptyMap(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype == MT_MAP && ma.mapval.FieldCount != 0)
}
func MlrvalIsNotEmpty(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_VOID {
		return MlrvalFromFalse()
	}
	if ma.mvtype == MT_STRING {
		if ma.printrep == "" {
			return MlrvalFromFalse()
		}
	}
	return MlrvalFromTrue()
}
func MlrvalIsNotMap(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype != MT_MAP)
}
func MlrvalIsNotArray(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype != MT_ARRAY)
}
func MlrvalIsNotNull(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype != MT_ABSENT && ma.mvtype != MT_VOID)
}
func MlrvalIsNull(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype == MT_ABSENT || ma.mvtype == MT_VOID)
}
func MlrvalIsNumeric(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype == MT_INT || ma.mvtype == MT_FLOAT)
}
func MlrvalIsPresent(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype != MT_ABSENT)
}
func MlrvalIsString(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.mvtype == MT_STRING || ma.mvtype == MT_VOID)
}

// ----------------------------------------------------------------
func assertingCommon(ma, check *Mlrval, description string, context *Context) Mlrval {
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
	return *ma
}

func MlrvalAssertingAbsent(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsAbsent(ma)
	return assertingCommon(ma, &check, "is_absent", context)
}
func MlrvalAssertingError(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsError(ma)
	return assertingCommon(ma, &check, "is_error", context)
}
func MlrvalAssertingBool(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsBool(ma)
	return assertingCommon(ma, &check, "is_bool", context)
}
func MlrvalAssertingBoolean(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsBoolean(ma)
	return assertingCommon(ma, &check, "is_boolean", context)
}
func MlrvalAssertingEmpty(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsEmpty(ma)
	return assertingCommon(ma, &check, "is_empty", context)
}
func MlrvalAssertingEmptyMap(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsEmptyMap(ma)
	return assertingCommon(ma, &check, "is_empty_map", context)
}
func MlrvalAssertingFloat(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsFloat(ma)
	return assertingCommon(ma, &check, "is_float", context)
}
func MlrvalAssertingInt(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsInt(ma)
	return assertingCommon(ma, &check, "is_int", context)
}
func MlrvalAssertingMap(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsMap(ma)
	return assertingCommon(ma, &check, "is_map", context)
}
func MlrvalAssertingArray(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsArray(ma)
	return assertingCommon(ma, &check, "is_array", context)
}
func MlrvalAssertingNonEmptyMap(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNonEmptyMap(ma)
	return assertingCommon(ma, &check, "is_non_empty_map", context)
}
func MlrvalAssertingNotEmpty(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNotEmpty(ma)
	return assertingCommon(ma, &check, "is_not_empty", context)
}
func MlrvalAssertingNotMap(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNotMap(ma)
	return assertingCommon(ma, &check, "is_not_map", context)
}
func MlrvalAssertingNotArray(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNotArray(ma)
	return assertingCommon(ma, &check, "is_not_array", context)
}
func MlrvalAssertingNotNull(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNotNull(ma)
	return assertingCommon(ma, &check, "is_not_null", context)
}
func MlrvalAssertingNull(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNull(ma)
	return assertingCommon(ma, &check, "is_null", context)
}
func MlrvalAssertingNumeric(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsNumeric(ma)
	return assertingCommon(ma, &check, "is_numeric", context)
}
func MlrvalAssertingPresent(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsPresent(ma)
	return assertingCommon(ma, &check, "is_present", context)
}
func MlrvalAssertingString(ma *Mlrval, context *Context) Mlrval {
	check := MlrvalIsString(ma)
	return assertingCommon(ma, &check, "is_string", context)
}
