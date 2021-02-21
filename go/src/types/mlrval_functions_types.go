package types

import (
	"fmt"
	"os"

	"miller/src/lib"
)

// ================================================================
func MlrvalTypeof(output, input1 *Mlrval) {
	output.SetFromString(input1.GetTypeName())
}

// ----------------------------------------------------------------
func string_to_int(output, input1 *Mlrval) {
	i, ok := lib.TryIntFromString(input1.printrep)
	if ok {
		output.SetFromInt(i)
	} else {
		output.SetFromError()
	}
}

func float_to_int(output, input1 *Mlrval) {
	output.SetFromInt(int(input1.floatval))
}

func bool_to_int(output, input1 *Mlrval) {
	if input1.boolval == true {
		output.SetFromInt(1)
	} else {
		output.SetFromInt(0)
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

func MlrvalToInt(output, input1 *Mlrval) {
	to_int_dispositions[input1.mvtype](output, input1)
}

// ----------------------------------------------------------------
func string_to_float(output, input1 *Mlrval) {
	f, ok := lib.TryFloat64FromString(input1.printrep)
	if ok {
		output.SetFromFloat64(f)
	} else {
		output.SetFromError()
	}
}

func int_to_float(output, input1 *Mlrval) {
	output.SetFromFloat64(float64(input1.intval))
}

func bool_to_float(output, input1 *Mlrval) {
	if input1.boolval == true {
		output.SetFromFloat64(1.0)
	} else {
		output.SetFromFloat64(0.0)
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

func MlrvalToFloat(output, input1 *Mlrval) {
	to_float_dispositions[input1.mvtype](output, input1)
}

// ----------------------------------------------------------------
func string_to_boolean(output, input1 *Mlrval) {
	b, ok := lib.TryBoolFromBoolString(input1.printrep)
	if ok {
		output.SetFromBool(b)
	} else {
		output.SetFromError()
	}
}

func int_to_bool(output, input1 *Mlrval) {
	output.SetFromBool(input1.intval != 0)
}

func float_to_bool(output, input1 *Mlrval) {
	output.SetFromBool(input1.floatval != 0.0)
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

func MlrvalToBoolean(output, input1 *Mlrval) {
	to_boolean_dispositions[input1.mvtype](output, input1)
}

// ----------------------------------------------------------------
func MlrvalIsAbsent(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype == MT_ABSENT)
}
func MlrvalIsError(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype == MT_ERROR)
}
func MlrvalIsBool(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype == MT_BOOL)
}
func MlrvalIsBoolean(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype == MT_BOOL)
}
func MlrvalIsEmpty(output, input1 *Mlrval) {
	if input1.mvtype == MT_VOID {
		output.SetFromTrue()
	} else if input1.mvtype == MT_STRING {
		if input1.printrep == "" {
			output.SetFromTrue()
		} else {
			output.SetFromFalse()
		}
	} else {
		output.SetFromFalse()
	}
}
func MlrvalIsEmptyMap(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype == MT_MAP && input1.mapval.FieldCount == 0)
}
func MlrvalIsFloat(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype == MT_FLOAT)
}
func MlrvalIsInt(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype == MT_INT)
}
func MlrvalIsMap(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype == MT_MAP)
}
func MlrvalIsArray(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype == MT_ARRAY)
}
func MlrvalIsNonEmptyMap(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype == MT_MAP && input1.mapval.FieldCount != 0)
}
func MlrvalIsNotEmpty(output, input1 *Mlrval) {
	if input1.mvtype == MT_VOID {
		output.SetFromFalse()
	} else if input1.mvtype == MT_STRING {
		if input1.printrep == "" {
			output.SetFromFalse()
		} else {
			output.SetFromTrue()
		}
	} else {
		output.SetFromTrue()
	}
}
func MlrvalIsNotMap(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype != MT_MAP)
}
func MlrvalIsNotArray(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype != MT_ARRAY)
}
func MlrvalIsNotNull(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype != MT_ABSENT && input1.mvtype != MT_VOID)
}
func MlrvalIsNull(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype == MT_ABSENT || input1.mvtype == MT_VOID)
}
func MlrvalIsNumeric(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype == MT_INT || input1.mvtype == MT_FLOAT)
}
func MlrvalIsPresent(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype != MT_ABSENT)
}
func MlrvalIsString(output, input1 *Mlrval) {
	output.SetFromBool(input1.mvtype == MT_STRING || input1.mvtype == MT_VOID)
}

// ----------------------------------------------------------------
func assertingCommon(output, input1, check *Mlrval, description string, context *Context) {
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
	output.CopyFrom(input1)
}

func MlrvalAssertingAbsent(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsAbsent(&check, input1)
	assertingCommon(output, input1, &check, "is_absent", context)
}
func MlrvalAssertingError(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsError(&check, input1)
	assertingCommon(output, input1, &check, "is_error", context)
}
func MlrvalAssertingBool(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsBool(&check, input1)
	assertingCommon(output, input1, &check, "is_bool", context)
}
func MlrvalAssertingBoolean(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsBoolean(&check, input1)
	assertingCommon(output, input1, &check, "is_boolean", context)
}
func MlrvalAssertingEmpty(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsEmpty(&check, input1)
	assertingCommon(output, input1, &check, "is_empty", context)
}
func MlrvalAssertingEmptyMap(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsEmptyMap(&check, input1)
	assertingCommon(output, input1, &check, "is_empty_map", context)
}
func MlrvalAssertingFloat(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsFloat(&check, input1)
	assertingCommon(output, input1, &check, "is_float", context)
}
func MlrvalAssertingInt(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsInt(&check, input1)
	assertingCommon(output, input1, &check, "is_int", context)
}
func MlrvalAssertingMap(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsMap(&check, input1)
	assertingCommon(output, input1, &check, "is_map", context)
}
func MlrvalAssertingArray(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsArray(&check, input1)
	assertingCommon(output, input1, &check, "is_array", context)
}
func MlrvalAssertingNonEmptyMap(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsNonEmptyMap(&check, input1)
	assertingCommon(output, input1, &check, "is_non_empty_map", context)
}
func MlrvalAssertingNotEmpty(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsNotEmpty(&check, input1)
	assertingCommon(output, input1, &check, "is_not_empty", context)
}
func MlrvalAssertingNotMap(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsNotMap(&check, input1)
	assertingCommon(output, input1, &check, "is_not_map", context)
}
func MlrvalAssertingNotArray(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsNotArray(&check, input1)
	assertingCommon(output, input1, &check, "is_not_array", context)
}
func MlrvalAssertingNotNull(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsNotNull(&check, input1)
	assertingCommon(output, input1, &check, "is_not_null", context)
}
func MlrvalAssertingNull(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsNull(&check, input1)
	assertingCommon(output, input1, &check, "is_null", context)
}
func MlrvalAssertingNumeric(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsNumeric(&check, input1)
	assertingCommon(output, input1, &check, "is_numeric", context)
}
func MlrvalAssertingPresent(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsPresent(&check, input1)
	assertingCommon(output, input1, &check, "is_present", context)
}
func MlrvalAssertingString(output, input1 *Mlrval, context *Context) {
	check := MlrvalFromAbsent()
	MlrvalIsString(&check, input1)
	assertingCommon(output, input1, &check, "is_string", context)
}
