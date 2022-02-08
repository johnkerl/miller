package bifs

import (
	"fmt"
	"math"
	"os"

	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ================================================================
func BIF_typeof(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromString(input1.GetTypeName())
}

// ----------------------------------------------------------------
func string_to_int(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	i, ok := lib.TryIntFromString(input1.AcquireStringValue())
	if ok {
		return mlrval.FromInt(i)
	} else {
		return mlrval.ERROR
	}
}

func float_to_int(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(int64(input1.AcquireFloatValue()))
}

func bool_to_int(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.AcquireBoolValue() == true {
		return mlrval.FromInt(1)
	} else {
		return mlrval.FromInt(0)
	}
}

var to_int_dispositions = [mlrval.MT_DIM]UnaryFunc{
	/*INT    */ _1u___,
	/*FLOAT  */ float_to_int,
	/*BOOL   */ bool_to_int,
	/*VOID   */ _void1,
	/*STRING */ string_to_int,
	/*ARRAY  */ _erro1,
	/*MAP    */ _erro1,
	/*FUNC   */ _erro1,
	/*ERROR  */ _erro1,
	/*NULL   */ _null1,
	/*ABSENT */ _absn1,
}

func BIF_int(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return to_int_dispositions[input1.Type()](input1)
}

// ----------------------------------------------------------------
func string_to_float(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	f, ok := lib.TryFloatFromString(input1.AcquireStringValue())
	if ok {
		return mlrval.FromFloat(f)
	} else {
		return mlrval.ERROR
	}
}

func int_to_float(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(float64(input1.AcquireIntValue()))
}

func bool_to_float(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.AcquireBoolValue() == true {
		return mlrval.FromFloat(1.0)
	} else {
		return mlrval.FromFloat(0.0)
	}
}

var to_float_dispositions = [mlrval.MT_DIM]UnaryFunc{
	/*INT    */ int_to_float,
	/*FLOAT  */ _1u___,
	/*BOOL   */ bool_to_float,
	/*VOID   */ _void1,
	/*STRING */ string_to_float,
	/*ARRAY  */ _erro1,
	/*MAP    */ _erro1,
	/*FUNC   */ _erro1,
	/*ERROR  */ _erro1,
	/*NULL   */ _null1,
	/*ABSENT */ _absn1,
}

func BIF_float(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return to_float_dispositions[input1.Type()](input1)
}

// ----------------------------------------------------------------
func string_to_boolean(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	b, ok := lib.TryBoolFromBoolString(input1.AcquireStringValue())
	if ok {
		return mlrval.FromBool(b)
	} else {
		return mlrval.ERROR
	}
}

func int_to_bool(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireIntValue() != 0)
}

func float_to_bool(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.AcquireFloatValue() != 0.0)
}

var to_boolean_dispositions = [mlrval.MT_DIM]UnaryFunc{
	/*INT    */ int_to_bool,
	/*FLOAT  */ float_to_bool,
	/*BOOL   */ _1u___,
	/*VOID   */ _void1,
	/*STRING */ string_to_boolean,
	/*ARRAY  */ _erro1,
	/*MAP    */ _erro1,
	/*FUNC   */ _erro1,
	/*ERROR  */ _erro1,
	/*NULL   */ _null1,
	/*ABSENT */ _absn1,
}

func BIF_boolean(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return to_boolean_dispositions[input1.Type()](input1)
}

// ----------------------------------------------------------------
func BIF_is_absent(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.IsAbsent())
}
func BIF_is_error(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.IsError())
}
func BIF_is_bool(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.IsBool())
}
func BIF_is_boolean(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.IsBool())
}
func BIF_is_empty(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsVoid() {
		return mlrval.TRUE
	} else if input1.IsString() {
		if input1.AcquireStringValue() == "" {
			return mlrval.TRUE
		} else {
			return mlrval.FALSE
		}
	} else {
		return mlrval.FALSE
	}
}
func BIF_is_emptymap(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.IsMap() && input1.AcquireMapValue().IsEmpty())
}
func BIF_is_float(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.IsFloat())
}
func BIF_is_int(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.IsInt())
}
func BIF_is_map(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.IsMap())
}
func BIF_is_array(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.IsArray())
}
func BIF_is_nonemptymap(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.IsMap() && input1.AcquireMapValue().FieldCount != 0)
}
func BIF_is_notempty(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsAbsent() {
		return mlrval.FALSE
	} else if input1.IsVoid() {
		return mlrval.FALSE
	} else if input1.IsString() {
		if input1.AcquireStringValue() == "" {
			return mlrval.FALSE
		} else {
			return mlrval.TRUE
		}
	} else {
		return mlrval.TRUE
	}
}
func BIF_is_notmap(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(!input1.IsMap())
}
func BIF_is_notarray(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(!input1.IsArray())
}
func BIF_is_notnull(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(!input1.IsAbsent() && !input1.IsVoid() && !input1.IsNull())
}
func BIF_is_null(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.IsAbsent() || input1.IsVoid() || input1.IsNull())
}
func BIF_is_numeric(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.IsInt() || input1.IsFloat())
}
func BIF_is_present(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(!input1.IsAbsent())
}
func BIF_is_string(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromBool(input1.IsStringOrVoid())
}
func BIF_is_nan(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	fval, ok := input1.GetFloatValue()
	if ok {
		return mlrval.FromBool(math.IsNaN(fval))
	} else {
		return mlrval.FALSE
	}
}

// ----------------------------------------------------------------
func assertingCommon(input1, check *mlrval.Mlrval, description string, context *types.Context) *mlrval.Mlrval {
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

func BIF_asserting_absent(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_absent(input1), "is_absent", context)
}
func BIF_asserting_error(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_error(input1), "is_error", context)
}
func BIF_asserting_bool(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_bool(input1), "is_bool", context)
}
func BIF_asserting_boolean(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_boolean(input1), "is_boolean", context)
}
func BIF_asserting_empty(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_empty(input1), "is_empty", context)
}
func BIF_asserting_emptyMap(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_emptymap(input1), "is_empty_map", context)
}
func BIF_asserting_float(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_float(input1), "is_float", context)
}
func BIF_asserting_int(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_int(input1), "is_int", context)
}
func BIF_asserting_map(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_map(input1), "is_map", context)
}
func BIF_asserting_array(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_array(input1), "is_array", context)
}
func BIF_asserting_nonempty_map(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_nonemptymap(input1), "is_non_empty_map", context)
}
func BIF_asserting_not_empty(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_notempty(input1), "is_not_empty", context)
}
func BIF_asserting_not_map(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_notmap(input1), "is_not_map", context)
}
func BIF_asserting_not_array(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_notarray(input1), "is_not_array", context)
}
func BIF_asserting_not_null(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_notnull(input1), "is_not_null", context)
}
func BIF_asserting_null(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_null(input1), "is_null", context)
}
func BIF_asserting_numeric(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_numeric(input1), "is_numeric", context)
}
func BIF_asserting_present(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_present(input1), "is_present", context)
}
func BIF_asserting_string(input1 *mlrval.Mlrval, context *types.Context) *mlrval.Mlrval {
	return assertingCommon(input1, BIF_is_string(input1), "is_string", context)
}
