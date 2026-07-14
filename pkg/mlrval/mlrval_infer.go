package mlrval

import (
	"strconv"

	"github.com/johnkerl/miller/v6/pkg/scan"
)

// Note: we do not infer bools from data files by default; see the opt-in
// --infer-booleans flag (issue #965).

// It's essential that we use mv.Type() not mv.mvtype since types are
// JIT-computed on first access for most data-file values. See type.go for more
// information.

func (mv *Mlrval) Type() MVType {
	if mv.mvtype == MT_PENDING {
		packageLevelInferrer(mv)
	}
	return mv.mvtype
}

// Support for mlr -S, mlr -A, mlr -O.
type tInferrer func(mv *Mlrval) *Mlrval

var packageLevelInferrer tInferrer = inferNormally

// SetInferrerOctalAsInt is for mlr -O.
func SetInferrerOctalAsInt() {
	packageLevelInferrer = inferWithOctalAsInt
}

// SetInferrerIntAsFloat is for mlr -F.
func SetInferrerIntAsFloat() {
	packageLevelInferrer = inferWithIntAsFloat
}

// SetInferrerStringOnly is for mlr -S.
func SetInferrerStringOnly() {
	packageLevelInferrer = inferString
}

// Support for mlr --infer-booleans and mlr --infer-special-floats (issue
// #965). Unlike the -S/-A/-O options above, these don't replace the
// package-level inferrer; they augment it, so they compose with -A and -O.
// They intentionally have no effect with -S / --infer-none, which disables
// type inference entirely.
var inferTrueFalseAsBool bool = false
var inferInfNanAsFloat bool = false

// SetInferrerBooleans is for mlr --infer-booleans.
func SetInferrerBooleans() {
	inferTrueFalseAsBool = true
}

// SetInferrerSpecialFloats is for mlr --infer-special-floats.
func SetInferrerSpecialFloats() {
	inferInfNanAsFloat = true
}

func inferNormally(mv *Mlrval) *Mlrval {
	if inferTrueFalseAsBool || inferInfNanAsFloat {
		if pmv := tryOptInLiteralInference(mv); pmv != nil {
			return pmv
		}
	}
	scanType := scan.FindScanType(mv.printrep)
	return normalInferrerTable[scanType](mv)
}

// inferWithOctalAsInt uses the leading-zero-as-int inferrer table (mlr -O).
func inferWithOctalAsInt(mv *Mlrval) *Mlrval {
	if inferTrueFalseAsBool || inferInfNanAsFloat {
		if pmv := tryOptInLiteralInference(mv); pmv != nil {
			return pmv
		}
	}
	scanType := scan.FindScanType(mv.printrep)
	return leadingZeroAsIntInferrerTable[scanType](mv)
}

// tryOptInLiteralInference handles opt-in inference of fixed literals --
// "true"/"false" for --infer-booleans, and "Inf"/"Infinity"/"NaN" for
// --infer-special-floats -- all of which the scan package classifies as
// strings. It returns nil when no literal matches, in which case the caller
// proceeds with normal scan-type dispatch.
func tryOptInLiteralInference(mv *Mlrval) *Mlrval {
	if len(mv.printrep) == 0 {
		return nil
	}

	if inferTrueFalseAsBool {
		// Exactly "true" and "false", matching the Miller DSL's boolean
		// literals as well as JSON's.
		if mv.printrep == "true" {
			return mv.SetFromBool(true)
		}
		if mv.printrep == "false" {
			return mv.SetFromBool(false)
		}
	}

	if inferInfNanAsFloat {
		// strconv.ParseFloat accepts "inf", "infinity", and "nan",
		// case-insensitively, with optional leading sign for the first two.
		// The first-character gate keeps ordinary numeric strings like "-123"
		// on the normal inference path, so that they still infer as int.
		rest := mv.printrep
		if rest[0] == '+' || rest[0] == '-' {
			rest = rest[1:]
		}
		if len(rest) > 0 {
			r0 := rest[0]
			if r0 == 'i' || r0 == 'I' || r0 == 'n' || r0 == 'N' {
				floatval, err := strconv.ParseFloat(mv.printrep, 64)
				if err == nil {
					return mv.SetFromPrevalidatedFloatString(mv.printrep, floatval)
				}
			}
		}
	}

	return nil
}

// inferWithIntAsFloat is for mlr -A.
func inferWithIntAsFloat(mv *Mlrval) *Mlrval {
	inferNormally(mv)
	if mv.Type() == MT_INT {
		mv.intf = float64(mv.intf.(int64))
		mv.mvtype = MT_FLOAT
	}
	return mv
}

// inferString is for mlr -S.
func inferString(mv *Mlrval) *Mlrval {
	return mv.SetFromString(mv.printrep)
}

// Important: synchronize this with the type-ordering in the scan package.
var normalInferrerTable []tInferrer = []tInferrer{
	inferString,
	inferDecimalInt,
	inferString, // inferLeadingZeroDecimalIntAsInt,
	inferOctalInt,
	inferString, // inferFromLeadingZeroOctalIntAsInt,
	inferHexInt,
	inferBinaryInt,
	inferMaybeFloat,
}

// Important: synchronize this with the type-ordering in the scan package.
var leadingZeroAsIntInferrerTable []tInferrer = []tInferrer{
	inferString,
	inferDecimalInt,
	inferLeadingZeroDecimalIntAsInt,
	inferOctalInt,
	inferFromLeadingZeroOctalIntAsInt,
	inferHexInt,
	inferBinaryInt,
	inferMaybeFloat,
}

// inferDecimalInt parses a base-10 integer or keeps the value as a string.
func inferDecimalInt(mv *Mlrval) *Mlrval {
	intval, err := strconv.ParseInt(mv.printrep, 10, 64)
	if err == nil {
		return mv.SetFromPrevalidatedIntString(mv.printrep, intval)
	}
	return mv.SetFromString(mv.printrep)
}

// inferLeadingZeroDecimalIntAsInt parses base-10 integers when leading zeros are allowed.
func inferLeadingZeroDecimalIntAsInt(mv *Mlrval) *Mlrval {
	intval, err := strconv.ParseInt(mv.printrep, 10, 64)
	if err == nil {
		return mv.SetFromPrevalidatedIntString(mv.printrep, intval)
	}
	return mv.SetFromString(mv.printrep)
}

// inferOctalInt parses explicit 0o-prefixed octal integers.
// E.g. explicit 0o377, not 0377
func inferOctalInt(mv *Mlrval) *Mlrval {
	return inferBaseInt(mv, 8)
}

// inferFromLeadingZeroOctalIntAsInt parses 0-prefixed octal integers.
func inferFromLeadingZeroOctalIntAsInt(mv *Mlrval) *Mlrval {
	intval, err := strconv.ParseInt(mv.printrep, 8, 64)
	if err == nil {
		return mv.SetFromPrevalidatedIntString(mv.printrep, intval)
	}
	return mv.SetFromString(mv.printrep)
}

// inferHexInt parses 0x-prefixed hex integers with two's-complement handling.
func inferHexInt(mv *Mlrval) *Mlrval {
	var input string
	var negate bool
	// Skip known leading 0x or -0x prefix
	switch mv.printrep[0] {
	case '-':
		input = mv.printrep[3:]
		negate = true
	case '+':
		input = mv.printrep[3:]
		negate = false
	default:
		input = mv.printrep[2:]
		negate = false
	}

	// Following twos-complement formatting familiar from all manner of
	// languages, including C which was Miller's original implementation
	// language, we want to allow 0x00....00 through 0x7f....ff as positive
	// 64-bit integers and 0x80....00 through 0xff....ff as negative ones. Go's
	// signed-int parsing explicitly doesn't allow that, but we don't want Go
	// semantics to dictate Miller semantics.  So, we try signed-int parsing
	// for 0x00....00 through 0x7f....ff, as well as positive or negative
	// decimal. Failing that, we try unsigned-int parsing for 0x80....00
	// through 0xff....ff.

	i0 := input[0]
	if len(input) == 16 && ('8' <= i0 && i0 <= 'f') {
		uintval, err := strconv.ParseUint(input, 16, 64)
		intval := int64(uintval)
		if negate {
			intval = -intval
		}
		if err == nil {
			return mv.SetFromPrevalidatedIntString(mv.printrep, intval)
		}
		return mv.SetFromString(mv.printrep)
	}
	intval, err := strconv.ParseInt(input, 16, 64)
	if negate {
		intval = -intval
	}
	if err == nil {
		return mv.SetFromPrevalidatedIntString(mv.printrep, intval)
	} else {
		return mv.SetFromString(mv.printrep)
	}

}

// inferBinaryInt parses 0b-prefixed binary integers.
func inferBinaryInt(mv *Mlrval) *Mlrval {
	return inferBaseInt(mv, 2)
}

// inferMaybeFloat parses floating-point values or keeps the value as a string.
func inferMaybeFloat(mv *Mlrval) *Mlrval {
	floatval, err := strconv.ParseFloat(mv.printrep, 64)
	if err == nil {
		return mv.SetFromPrevalidatedFloatString(mv.printrep, floatval)
	}
	return mv.SetFromString(mv.printrep)
}

// inferBaseInt is shared code for parsing 0o/0b integers.
func inferBaseInt(mv *Mlrval, base int) *Mlrval {
	var input string
	var negate bool
	// Skip known leading 0x or -0x prefix
	switch mv.printrep[0] {
	case '-':
		input = mv.printrep[3:]
		negate = true
	case '+':
		input = mv.printrep[3:]
		negate = false
	default:
		input = mv.printrep[2:]
		negate = false
	}
	intval, err := strconv.ParseInt(input, base, 64)
	if err == nil {
		if negate {
			intval = -intval
		}
		return mv.SetFromPrevalidatedIntString(mv.printrep, intval)
	}
	return mv.SetFromString(mv.printrep)
}
