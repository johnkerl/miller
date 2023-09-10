package bifs

import (
	"math"

	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
)

func BIF_urand() *mlrval.Mlrval {
	return mlrval.FromFloat(
		lib.RandFloat64(),
	)
}

func BIF_urand32() *mlrval.Mlrval {
	return mlrval.FromInt(
		int64(
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
		return mlrval.FromNotIntError("urandint", input1)
	}
	if !input2.IsInt() {
		return mlrval.FromNotIntError("urandint", input2)
	}

	a := input1.AcquireIntValue()
	b := input2.AcquireIntValue()

	var lo int64 = 0
	var hi int64 = 0
	if a <= b {
		lo = a
		hi = b + 1
	} else {
		lo = b
		hi = a + 1
	}
	u := int64(math.Floor(float64(lo) + float64((hi-lo))*lib.RandFloat64()))
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
		return mlrval.FromNotNumericError("urandrange", input1)
	}
	if !bok {
		return mlrval.FromNotNumericError("urandrange", input2)
	}
	return mlrval.FromFloat(
		a + (b-a)*lib.RandFloat64(),
	)
}

func BIF_urandelement(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	arrayval := input1.GetArray()
	if arrayval == nil { // not an array
		return mlrval.FromNotArrayError("urandelement", input1)
	}
	if len(arrayval) == 0 {
		return mlrval.FromErrorString("urandelement: received a zero-length array as input")
	}

	// lo is inclusive, hi is exclusive
	index := lib.RandRange(0, int64(len(arrayval)))
	return arrayval[index].Copy()
}
