package mlrval

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/scan"
)

// TODO: comment no infer-bool from data files. Always false in this path.

// It's essential that we use mv.Type() not mv.mvtype since types are
// JIT-computed on first access for most data-file values. See type.go for more
// information.

func (mv *Mlrval) Type() MVType {
	if mv.mvtype == MT_PENDING {
		packageLevelInferrer(mv, false)
	}
	return mv.mvtype
}

// Support for mlr -S, mlr -A, mlr -O.
type tInferrer func(mv *Mlrval, inferBool bool) *Mlrval

// xxx temp
var packageLevelInferrer tInferrer = inferWithOctalAsString

//var packageLevelInferrer tInferrer = inferTemp

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
func inferWithOctalAsString(mv *Mlrval, inferBool bool) *Mlrval {
	inferWithOctalAsInt(mv, inferBool)
	if mv.mvtype != MT_INT && mv.mvtype != MT_FLOAT {
		return mv
	}

	if octalDetector.MatchString(mv.printrep) {
		return mv.SetFromString(mv.printrep)
	} else {
		return mv
	}
}

// inferWithOctalAsInt is for mlr -O.
func inferWithOctalAsInt(mv *Mlrval, inferBool bool) *Mlrval {
	if mv.printrep == "" {
		return mv.SetFromVoid()
	}

	intval, iok := lib.TryIntFromString(mv.printrep)
	if iok {
		return mv.SetFromPrevalidatedIntString(mv.printrep, intval)
	}

	if downcasedFloatNamesToNotInfer[strings.ToLower(mv.printrep)] == false {
		floatval, fok := lib.TryFloatFromString(mv.printrep)
		if fok {
			return mv.SetFromPrevalidatedFloatString(mv.printrep, floatval)
		}
	}

	if inferBool {
		boolval, bok := lib.TryBoolFromBoolString(mv.printrep)
		if bok {
			return mv.SetFromPrevalidatedBoolString(mv.printrep, boolval)
		}
	}
	return mv.SetFromString(mv.printrep)
}

// inferWithIntAsFloat is for mlr -A.
func inferWithIntAsFloat(mv *Mlrval, inferBool bool) *Mlrval {
	inferWithOctalAsString(mv, inferBool)
	if mv.Type() == MT_INT {
		mv.floatval = float64(mv.intval)
		mv.mvtype = MT_FLOAT
	}
	return mv
}

// inferStringOnly is for mlr -S.
func inferStringOnly(mv *Mlrval, inferBool bool) *Mlrval {
	return mv.SetFromString(mv.printrep)
}

// ----------------------------------------------------------------
// experimental

// TODO: comment
func inferFromDecimalInt(mv *Mlrval, inferBool bool) *Mlrval {
	intval, err := strconv.ParseInt(mv.printrep, 10, 64)
	if err == nil {
		return mv.SetFromPrevalidatedIntString(mv.printrep, int(intval))
	} else {
		return mv.SetFromString(mv.printrep)
	}
}

// TODO: comment
func inferFromOctalInt(mv *Mlrval, inferBool bool) *Mlrval {
	intval, err := strconv.ParseInt(mv.printrep, 8, 64)
	if err == nil {
		return mv.SetFromPrevalidatedIntString(mv.printrep, int(intval))
	} else {
		return mv.SetFromString(mv.printrep)
	}
}

// TODO: comment
func inferFromHexInt(mv *Mlrval, inferBool bool) *Mlrval {
	intval, err := strconv.ParseInt(mv.printrep, 16, 64)
	if err == nil {
		return mv.SetFromPrevalidatedIntString(mv.printrep, int(intval))
	} else {
		return mv.SetFromString(mv.printrep)
	}
}

// TODO: comment
func inferFromBinaryInt(mv *Mlrval, inferBool bool) *Mlrval {
	intval, err := strconv.ParseInt(mv.printrep, 2, 64)
	// xxx to do: length check & overflow/uint check
	if err == nil {
		return mv.SetFromPrevalidatedIntString(mv.printrep, int(intval))
	} else {
		return mv.SetFromString(mv.printrep)
	}
}

// TODO: comment
func inferFromMaybeFloat(mv *Mlrval, inferBool bool) *Mlrval {
	floatval, err := strconv.ParseFloat(mv.printrep, 64)
	if err == nil {
		return mv.SetFromPrevalidatedFloatString(mv.printrep, floatval)
	} else {
		return mv.SetFromString(mv.printrep)
	}
}

// TODO: comment
func inferFromBool(mv *Mlrval, inferBool bool) *Mlrval {
	if mv.printrep == "true" {
		return mv.SetFromPrevalidatedBoolString(mv.printrep, true)
	} else {
		return mv.SetFromPrevalidatedBoolString(mv.printrep, false)
	}
}

// const (
//     scanTypeString     ScanType = 0
//     scanTypeDecimalInt          = 1
//     scanTypeOctalInt            = 2
//     scanTypeHexInt              = 3
//     scanTypeBinaryInt           = 4
//     scanTypeMaybeFloat          = 5
//     scanTypeBool                = 6
// )

var tempScanTypeInferrerTable []tInferrer = []tInferrer{
	inferStringOnly,
	inferFromDecimalInt,
	inferFromOctalInt,
	inferFromHexInt,
	inferFromBinaryInt,
	inferFromMaybeFloat,
	inferFromBool,
}

func inferTemp(mv *Mlrval, inferBool bool) *Mlrval {
	scanType := scan.FindScanType(mv.printrep)
	return tempScanTypeInferrerTable[scanType](mv, inferBool)
}
