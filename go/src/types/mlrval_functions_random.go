package types

import (
	"math"

	"miller/src/lib"
)

func MlrvalUrand(output *Mlrval) {
	output.SetFromFloat64(
		lib.RandFloat64(),
	)
}

func MlrvalUrand32(output *Mlrval) {
	output.SetFromInt(
		int(
			lib.RandUint32(),
		),
	)
}

// TODO: use a disposition matrix
func MlrvalUrandInt(input1, input2 *Mlrval) Mlrval {
	if !input1.IsLegit() {
		return *input1
	}
	if !input2.IsLegit() {
		return *input2
	}
	if !input1.IsInt() {
		return MlrvalFromError()
	}
	if !input2.IsInt() {
		return MlrvalFromError()
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

func MlrvalUrandRange(input1, input2 *Mlrval) Mlrval {
	if !input1.IsLegit() {
		return *input1
	}
	if !input2.IsLegit() {
		return *input2
	}
	a, aok := input1.GetNumericToFloatValue()
	b, bok := input2.GetNumericToFloatValue()
	if !aok {
		return MlrvalFromError()
	}
	if !bok {
		return MlrvalFromError()
	}
	return MlrvalFromFloat64(
		a + (b-a)*lib.RandFloat64(),
	)
}
