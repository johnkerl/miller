package types

import (
	"fmt"
	"os"

	"mlr/internal/pkg/lib"
)

// ================================================================
func BIF_typeof(input1 *Mlrval) *Mlrval {
	return MlrvalFromString(input1.GetTypeName())
}

// ----------------------------------------------------------------
func string_to_int(input1 *Mlrval) *Mlrval {
	i, ok := lib.TryIntFromString(input1.printrep)
	if ok {
		return MlrvalFromInt(i)
	} else {
		return MLRVAL_ERROR
	}
}

func float_to_int(input1 *Mlrval) *Mlrval {
	return MlrvalFromInt(int(input1.floatval))
}

func bool_to_int(input1 *Mlrval) *Mlrval {
	if input1.boolval == true {
		return MlrvalFromInt(1)
	} else {
		return MlrvalFromInt(0)
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
	/*FUNC   */ _erro1,
}

func BIF_int(input1 *Mlrval) *Mlrval {
	return to_int_dispositions[input1.mvtype](input1)
}

// ----------------------------------------------------------------
func string_to_float(input1 *Mlrval) *Mlrval {
	f, ok := lib.TryFloat64FromString(input1.printrep)
	if ok {
		return MlrvalFromFloat64(f)
	} else {
		return MLRVAL_ERROR
	}
}

func int_to_float(input1 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(float64(input1.intval))
}

func bool_to_float(input1 *Mlrval) *Mlrval {
	if input1.boolval == true {
		return MlrvalFromFloat64(1.0)
	} else {
		return MlrvalFromFloat64(0.0)
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
	/*FUNC   */ _erro1,
}

func BIF_float(input1 *Mlrval) *Mlrval {
	return to_float_dispositions[input1.mvtype](input1)
}

// ----------------------------------------------------------------
func string_to_boolean(input1 *Mlrval) *Mlrval {
	b, ok := lib.TryBoolFromBoolString(input1.printrep)
	if ok {
		return MlrvalFromBool(b)
	} else {
		return MLRVAL_ERROR
	}
}

func int_to_bool(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.intval != 0)
}

func float_to_bool(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.floatval != 0.0)
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
	/*FUNC   */ _erro1,
}

func BIF_boolean(input1 *Mlrval) *Mlrval {
	return to_boolean_dispositions[input1.mvtype](input1)
}

// ----------------------------------------------------------------
func BIF_is_absent(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_ABSENT)
}
func BIF_is_error(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_ERROR)
}
func BIF_is_bool(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_BOOL)
}
func BIF_is_boolean(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_BOOL)
}
func BIF_is_empty(input1 *Mlrval) *Mlrval {
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
func BIF_is_emptymap(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_MAP && input1.mapval.IsEmpty())
}
func BIF_is_float(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_FLOAT)
}
func BIF_is_int(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_INT)
}
func BIF_is_map(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_MAP)
}
func BIF_is_array(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_ARRAY)
}
func BIF_is_nonemptymap(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_MAP && input1.mapval.FieldCount != 0)
}
func BIF_is_notempty(input1 *Mlrval) *Mlrval {
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
func BIF_is_notmap(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype != MT_MAP)
}
func BIF_is_notarray(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype != MT_ARRAY)
}
func BIF_is_notnull(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype != MT_ABSENT && input1.mvtype != MT_VOID)
}
func BIF_is_null(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_ABSENT || input1.mvtype == MT_VOID)
}
func BIF_is_numeric(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_INT || input1.mvtype == MT_FLOAT)
}
func BIF_is_present(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype != MT_ABSENT)
}
func BIF_is_string(input1 *Mlrval) *Mlrval {
	return MlrvalFromBool(input1.mvtype == MT_STRING || input1.mvtype == MT_VOID)
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

func BIF_asserting_absent(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_absent(input1), "is_absent", context)
}
func BIF_asserting_error(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_error(input1), "is_error", context)
}
func BIF_asserting_bool(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_bool(input1), "is_bool", context)
}
func BIF_asserting_boolean(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_boolean(input1), "is_boolean", context)
}
func BIF_asserting_empty(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_empty(input1), "is_empty", context)
}
func BIF_asserting_emptyMap(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_emptymap(input1), "is_empty_map", context)
}
func BIF_asserting_float(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_float(input1), "is_float", context)
}
func BIF_asserting_int(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_int(input1), "is_int", context)
}
func BIF_asserting_map(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_map(input1), "is_map", context)
}
func BIF_asserting_array(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_array(input1), "is_array", context)
}
func BIF_asserting_nonempty_map(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_nonemptymap(input1), "is_non_empty_map", context)
}
func BIF_asserting_not_empty(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_notempty(input1), "is_not_empty", context)
}
func BIF_asserting_not_map(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_notmap(input1), "is_not_map", context)
}
func BIF_asserting_not_array(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_notarray(input1), "is_not_array", context)
}
func BIF_asserting_not_null(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_notnull(input1), "is_not_null", context)
}
func BIF_asserting_null(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_null(input1), "is_null", context)
}
func BIF_asserting_numeric(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_numeric(input1), "is_numeric", context)
}
func BIF_asserting_present(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_present(input1), "is_present", context)
}
func BIF_asserting_string(input1 *Mlrval, context *Context) *Mlrval {
	return assertingCommon(input1, BIF_is_string(input1), "is_string", context)
}
