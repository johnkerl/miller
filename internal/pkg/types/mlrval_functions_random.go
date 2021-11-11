package types

import (
	"math"

	"mlr/internal/pkg/lib"
)

func BIF_urand() *Mlrval {
	return MlrvalFromFloat64(
		lib.RandFloat64(),
	)
}

func BIF_urand32() *Mlrval {
	return MlrvalFromInt(
		int(
			lib.RandUint32(),
		),
	)
}

// TODO: use a disposition matrix
func BIF_urandint(input1, input2 *Mlrval) *Mlrval {
	if !input1.IsLegit() {
		return input1
	}
	if !input2.IsLegit() {
		return input2
	}
	if !input1.IsInt() {
		return MLRVAL_ERROR
	}
	if !input2.IsInt() {
		return MLRVAL_ERROR
	}

	a := input1.intval
	b := input2.intval

	var lo int = 0
	var hi int = 0
	if a <= b {
		lo = a
		hi = b + 1
	} else {
		lo = b
		hi = a + 1
	}
	u := int(math.Floor(float64(lo) + float64((hi-lo))*lib.RandFloat64()))
	return MlrvalFromInt(u)
}

func BIF_urandrange(input1, input2 *Mlrval) *Mlrval {
	if !input1.IsLegit() {
		return input1
	}
	if !input2.IsLegit() {
		return input2
	}
	a, aok := input1.GetNumericToFloatValue()
	b, bok := input2.GetNumericToFloatValue()
	if !aok {
		return MLRVAL_ERROR
	}
	if !bok {
		return MLRVAL_ERROR
	}
	return MlrvalFromFloat64(
		a + (b-a)*lib.RandFloat64(),
	)
}
