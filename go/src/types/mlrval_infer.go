// ================================================================
// Constructors
// ================================================================

package types

import (
	"strings"

	"mlr/src/lib"
)

// XXX
// TODO: rename w/ "blessed" or "prevalidated" or somesuch
// xxx comment why two -- one for from parsed user data; other for from math ops
func MlrvalFromIntString(input string) *Mlrval {
	ival, ok := lib.TryIntFromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	lib.InternalCodingErrorIf(!ok)
	var mv Mlrval
	mv.mvtype = MT_INT
	mv.printrep = input
	mv.printrepValid = true
	mv.intval = ival
	return &mv
}

// xxx comment why two -- one for from parsed user data; other for from math ops
// xxx comment assummption is input-string already deemed parseable so no error return
// TODO: rename w/ "blessed" or "prevalidated" or somesuch
func MlrvalFromFloat64String(input string) *Mlrval {
	fval, ok := lib.TryFloat64FromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	lib.InternalCodingErrorIf(!ok)
	var mv Mlrval
	mv.mvtype = MT_FLOAT
	mv.printrep = input
	mv.printrepValid = true
	mv.floatval = fval
	return &mv
}

// Used by MlrvalFormatter (fmtnum DSL function, format-values verb, etc).
// Each mlrval has printrep and a printrepValid for its original string, then a
// type-code like MT_INT or MT_FLOAT, and type-specific storage like intval or
// floatval.
//
// If the user has taken a mlrval with original string "3.14" and formatted it
// with "%.4f" then its printrep will be "3.1400" but its type should still be
// MT_FLOAT.
//
// If on the other hand the user has formatted the same mlrval with "[[%.4f]]"
// then its printrep will be "[[3.1400]]" and it will be MT_STRING.
// This method supports that.

func MlrvalTryPointerFromFloatString(input string) *Mlrval {
	_, fok := lib.TryFloat64FromString(input)
	if fok {
		return MlrvalFromFloat64String(input)
	} else {
		return MlrvalFromString(input)
	}
}

func MlrvalTryPointerFromIntString(input string) *Mlrval {
	_, iok := lib.TryIntFromString(input)
	if iok {
		return MlrvalFromIntString(input)
	} else {
		return MlrvalFromString(input)
	}
}

// ================================================================
// MlrvalFromInferredTypeForDataFiles is for parsing field values from
// data files (except JSON, which is typed -- "true" and true are distinct).
// Mostly the same as MlrvalFromInferredType, except it doesn't
// auto-infer true/false to bool; don't auto-infer NaN/Inf to float; etc.

func MlrvalFromInferredTypeForDataFiles(input string) *Mlrval {
	return inferrer(input)
}

func MlrvalFromInferredType(input string) *Mlrval {
	if input == "" {
		return MLRVAL_VOID
	}

	// TODO: wrap this lib function here in this package
	_, iok := lib.TryIntFromString(input)
	if iok {
		return MlrvalFromIntString(input)
	}

	_, fok := lib.TryFloat64FromString(input)
	if fok {
		return MlrvalFromFloat64String(input)
	}

	_, bok := lib.TryBoolFromBoolString(input)
	if bok {
		return MlrvalFromBoolString(input)
	}

	return MlrvalFromString(input)
}

// ================================================================

var downcasedFloatNamesToNotInfer = map[string]bool{
	"inf":       true,
	"+inf":      true,
	"-inf":      true,
	"infinity":  true,
	"+infinity": true,
	"-infinity": true,
	"nan":       true,
}

type inferrerFunc func(input string) *Mlrval

var inferrer inferrerFunc = inferNormally

func SetInferrerNoOctal() {
	inferrer = inferWithOctalSuppress
}
func SetInferrerIntAsFloat() {
	inferrer = inferWithIntAsFloat
}
func SetInferrerStringOnly() {
	inferrer = inferStringOnly
}

func inferNormally(input string) *Mlrval {
	if input == "" {
		return MLRVAL_VOID
	}

	_, iok := lib.TryIntFromString(input)
	if iok {
		return MlrvalFromIntString(input)
	}

	if downcasedFloatNamesToNotInfer[strings.ToLower(input)] == false {
		_, fok := lib.TryFloat64FromString(input)
		if fok {
			return MlrvalFromFloat64String(input)
		}
	}

	return MlrvalFromString(input)
}

func inferWithOctalSuppress(input string) *Mlrval {
	if input == "" {
		return MLRVAL_VOID
	}

	_, iok := lib.TryIntFromStringNoOctal(input)
	if iok {
		return MlrvalFromIntString(input)
	}

	if downcasedFloatNamesToNotInfer[strings.ToLower(input)] == false {
		_, fok := lib.TryFloat64FromString(input)
		if fok {
			return MlrvalFromFloat64String(input)
		}
	}

	return MlrvalFromString(input)
}

func inferWithIntAsFloat(input string) *Mlrval {
	output := inferNormally(input)
	if output.mvtype == MT_INT {
		return MlrvalFromFloat64(float64(output.intval))
	} else {
		return output
	}
}

func inferStringOnly(input string) *Mlrval {
	if input == "" {
		return MLRVAL_VOID
	}

	return MlrvalFromString(input)
}
