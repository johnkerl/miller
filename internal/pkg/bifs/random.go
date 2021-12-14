package types

import (
	"math"

	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

func BIF_urand() *mlrval.Mlrval {
	return mlrval.FromFloat(
		lib.RandFloat64(),
	)
}

func BIF_urand32() *mlrval.Mlrval {
	return mlrval.FromInt(
		int(
			lib.RandUint32(),
		),
	)
}

// TODO: use a disposition matrix
func BIF_urandint(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsLegit() {
		return input1
	}
	if !input2.IsLegit() {
		return input2
	}
	if !input1.IsInt() {
		return mlrval.ERROR
	}
	if !input2.IsInt() {
		return mlrval.ERROR
	}

	a := input1.AcquireIntValue()
	b := input2.AcquireIntValue()

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
	return mlrval.FromInt(u)
}

func BIF_urandrange(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsLegit() {
		return input1
	}
	if !input2.IsLegit() {
		return input2
	}
	a, aok := input1.GetNumericToFloatValue()
	b, bok := input2.GetNumericToFloatValue()
	if !aok {
		return mlrval.ERROR
	}
	if !bok {
		return mlrval.ERROR
	}
	return mlrval.FromFloat(
		a + (b-a)*lib.RandFloat64(),
	)
}
