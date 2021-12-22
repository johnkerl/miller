package mlrval

import (
	"regexp"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/lib"
)

// TODO: comment no infer-bool from data files. Always false in this path.

// It's essential that we use mv.Type() not mv.mvtype since types are
// JIT-computed on first access for most data-file values. See type.go for more
// information.

func (mv *Mlrval) Type() MVType {
	if mv.mvtype == MT_PENDING {
		packageLevelInferrer(mv, mv.printrep, false)
	}
	return mv.mvtype
}

// Support for mlr -S, mlr -A, mlr -O.
type tInferrer func(mv *Mlrval, input string, inferBool bool) *Mlrval

var packageLevelInferrer tInferrer = inferWithOctalAsString

// SetInferrerOctalAsInt is for default behavior.
func SetInferrerOctalAsString() {
	packageLevelInferrer = inferWithOctalAsString
}

// SetInferrerOctalAsInt is for mlr -O.
func SetInferrerOctalAsInt() {
	packageLevelInferrer = inferWithOctalAsInt
}

// SetInferrerStringOnly is for mlr -A.
func SetInferrerIntAsFloat() {
	packageLevelInferrer = inferWithIntAsFloat
}

// SetInferrerStringOnly is for mlr -S.
func SetInferrerStringOnly() {
	packageLevelInferrer = inferStringOnly
}

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

var octalDetector = regexp.MustCompile("^-?0[0-9]+")

// inferWithOctalAsString is for default behavior.
func inferWithOctalAsString(mv *Mlrval, input string, inferBool bool) *Mlrval {
	inferWithOctalAsInt(mv, input, inferBool)
	if mv.mvtype != MT_INT && mv.mvtype != MT_FLOAT {
		return mv
	}

	if octalDetector.MatchString(mv.printrep) {
		return mv.SetFromString(input)
	} else {
		return mv
	}
}

// inferWithOctalAsInt is for mlr -O.
func inferWithOctalAsInt(mv *Mlrval, input string, inferBool bool) *Mlrval {
	if input == "" {
		return mv.SetFromVoid()
	}

	intval, iok := lib.TryIntFromString(input)
	if iok {
		return mv.SetFromPrevalidatedIntString(input, intval)
	}

	if downcasedFloatNamesToNotInfer[strings.ToLower(input)] == false {
		floatval, fok := lib.TryFloatFromString(input)
		if fok {
			return mv.SetFromPrevalidatedFloatString(input, floatval)
		}
	}

	if inferBool {
		boolval, bok := lib.TryBoolFromBoolString(input)
		if bok {
			return mv.SetFromPrevalidatedBoolString(input, boolval)
		}
	}
	return mv.SetFromString(input)
}

// inferWithIntAsFloat is for mlr -A.
func inferWithIntAsFloat(mv *Mlrval, input string, inferBool bool) *Mlrval {
	inferWithOctalAsString(mv, input, inferBool)
	if mv.Type() == MT_INT {
		mv.floatval = float64(mv.intval)
		mv.mvtype = MT_FLOAT
	}
	return mv
}

// inferStringOnly is for mlr -S.
func inferStringOnly(mv *Mlrval, input string, inferBool bool) *Mlrval {
	return mv.SetFromString(input)
}
