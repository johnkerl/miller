// ================================================================
// Constructors
// ================================================================

package types

import (
	"strings"

	"mlr/internal/pkg/lib"
)

// ----------------------------------------------------------------

// MlrvalFromPrevalidatedIntString is for situations where the string has
// already been determined to be parseable as int. For example, int literals in
// the Miller DSL fall into this category, as the parser has already matched
// them.
func MlrvalFromPrevalidatedIntString(input string) *Mlrval {
	ival, ok := lib.TryIntFromString(input)
	lib.InternalCodingErrorIf(!ok)
	return &Mlrval{
		mvtype:        MT_INT,
		printrep:      input,
		printrepValid: true,
		intval:        ival,
	}
}

// MlrvalFromPrevalidatedFloat64String is for situations where the string has
// already been determined to be parseable as int. For example, int literals in
// the Miller DSL fall into this category, as the parser has already matched
// them.
func MlrvalFromPrevalidatedFloat64String(input string) *Mlrval {
	fval, ok := lib.TryFloat64FromString(input)
	lib.InternalCodingErrorIf(!ok)
	return &Mlrval{
		mvtype:        MT_FLOAT,
		printrep:      input,
		printrepValid: true,
		floatval:      fval,
	}
}

// ----------------------------------------------------------------

// MlrvalTryPointerFromFloatString is used by MlrvalFormatter (fmtnum DSL
// function, format-values verb, etc).  Each mlrval has printrep and a
// printrepValid for its original string, then a type-code like MT_INT or
// MT_FLOAT, and type-specific storage like intval or floatval.
//
// If the user has taken a mlrval with original string "3.14" and formatted it
// with "%.4f" then its printrep will be "3.1400" but its type should still be
// MT_FLOAT.
//
// If on the other hand the user has formatted the same mlrval with "[[%.4f]]"
// then its printrep will be "[[3.1400]]" and it will be MT_STRING.
// This function supports that.
func MlrvalTryPointerFromFloatString(input string) *Mlrval {
	_, fok := lib.TryFloat64FromString(input)
	if fok {
		return MlrvalFromPrevalidatedFloat64String(input)
	} else {
		return MlrvalFromString(input)
	}
}

// MlrvalTryPointerFromIntString is  used by MlrvalFormatter (fmtnum DSL
// function, format-values verb, etc).  Each mlrval has printrep and a
// printrepValid for its original string, then a type-code like MT_INT or
// MT_FLOAT, and type-specific storage like intval or floatval.
//
// If the user has taken a mlrval with original string "314" and formatted it
// with "0x%04x" then its printrep will be "0x013a" but its type should still be
// MT_INT.
//
// If on the other hand the user has formatted the same mlrval with
// "[[%0x04x]]" then its printrep will be "[[0x013a]]" and it will be
// MT_STRING.  This function supports that.
func MlrvalTryPointerFromIntString(input string) *Mlrval {
	_, iok := lib.TryIntFromString(input)
	if iok {
		return MlrvalFromPrevalidatedIntString(input)
	} else {
		return MlrvalFromString(input)
	}
}

// ----------------------------------------------------------------

// When loading data files, don't scan these words into floats -- even though
// the Go library is willing to do so.
var downcasedFloatNamesToNotInfer = map[string]bool{
	"inf":       true,
	"+inf":      true,
	"-inf":      true,
	"infinity":  true,
	"+infinity": true,
	"-infinity": true,
	"nan":       true,
}

// ----------------------------------------------------------------
type tInferrer func(input string, inferBool bool) *Mlrval

var inferrer tInferrer = inferNormally

func SetInferrerNoOctal() {
	inferrer = inferWithOctalSuppress
}
func SetInferrerIntAsFloat() {
	inferrer = inferWithIntAsFloat
}
func SetInferrerStringOnly() {
	inferrer = inferStringOnly
}

// MlrvalFromInferredTypeForDataFiles is for parsing field values directly from
// data files (except JSON, which is typed -- "true" and true are distinct).
// Mostly the same as MlrvalFromInferredType, except it doesn't auto-infer
// true/false to bool; don't auto-infer NaN/Inf to float; etc.
func MlrvalFromInferredTypeForDataFiles(input string) *Mlrval {
	return inferrer(input, false)
}

// MlrvalFromInferredType is for parsing field values not directly from data
// files.  Mostly the same as MlrvalFromInferredTypeForDataFiles, except it
// auto-infers true/false to bool; don't auto-infer NaN/Inf to float; etc.
func MlrvalFromInferredType(input string) *Mlrval {
	return inferrer(input, true)
}

func inferNormally(input string, inferBool bool) *Mlrval {
	if input == "" {
		return MLRVAL_VOID
	}

	_, iok := lib.TryIntFromString(input)
	if iok {
		return MlrvalFromPrevalidatedIntString(input)
	}

	if downcasedFloatNamesToNotInfer[strings.ToLower(input)] == false {
		_, fok := lib.TryFloat64FromString(input)
		if fok {
			return MlrvalFromPrevalidatedFloat64String(input)
		}
	}

	if inferBool {
		_, bok := lib.TryBoolFromBoolString(input)
		if bok {
			return MlrvalFromBoolString(input)
		}
	}

	return MlrvalFromString(input)
}

func inferWithOctalSuppress(input string, inferBool bool) *Mlrval {
	output := inferNormally(input, inferBool)
	if output.mvtype != MT_INT {
		return output
	}

	if input[0] == '0' && len(input) > 1 {
		c := input[1]
		if c != 'x' && c != 'X' && c != 'b' && c != 'B' {
			return MlrvalFromString(input)
		}
	}
	if strings.HasPrefix(input, "-0") && len(input) > 2 {
		c := input[2]
		if c != 'x' && c != 'X' && c != 'b' && c != 'B' {
			return MlrvalFromString(input)
		}
	}

	return output
}

func inferWithIntAsFloat(input string, inferBool bool) *Mlrval {
	output := inferNormally(input, inferBool)
	if output.mvtype == MT_INT {
		return &Mlrval{
			mvtype:        MT_FLOAT,
			printrepValid: true,
			printrep:      input,
			floatval:      float64(output.intval),
		}
	} else {
		return output
	}
}

func inferStringOnly(input string, inferBool bool) *Mlrval {
	if input == "" {
		return MLRVAL_VOID
	}

	return MlrvalFromString(input)
}
