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
func MlrvalUrandInt(output, input1, input2 *Mlrval) {
	if !input1.IsLegit() {
		output.CopyFrom(input1)
		return
	}
	if !input2.IsLegit() {
		output.CopyFrom(input2)
		return
	}
	if !input1.IsInt() {
		output.SetFromError()
		return
	}
	if !input2.IsInt() {
		output.SetFromError()
		return
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
	output.SetFromInt(u)
}

func MlrvalUrandRange(output, input1, input2 *Mlrval) {
	if !input1.IsLegit() {
		output.CopyFrom(input1)
		return
	}
	if !input2.IsLegit() {
		output.CopyFrom(input2)
		return
	}
	a, aok := input1.GetNumericToFloatValue()
	b, bok := input2.GetNumericToFloatValue()
	if !aok {
		output.SetFromError()
		return
	}
	if !bok {
		output.SetFromError()
		return
	}
	output.SetFromFloat64(
		a + (b-a)*lib.RandFloat64(),
	)
}
