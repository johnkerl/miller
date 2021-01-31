package types

import (
	"math"

	"miller/lib"
)

func MlrvalUrand() Mlrval {
	return MlrvalFromFloat64(
		lib.RandFloat64(),
	)
}

func MlrvalUrand32() Mlrval {
	return MlrvalFromInt(
		int(
			lib.RandUint32(),
		),
	)
}

// TODO: use a disposition matrix
func MlrvalUrandInt(ma, mb *Mlrval) Mlrval {
	if !ma.IsLegit() {
		return *ma
	}
	if !mb.IsLegit() {
		return *mb
	}
	if !ma.IsInt() {
		return MlrvalFromError()
	}
	if !mb.IsInt() {
		return MlrvalFromError()
	}

	a := ma.intval
	b := mb.intval

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

func MlrvalUrandRange(ma, mb *Mlrval) Mlrval {
	if !ma.IsLegit() {
		return *ma
	}
	if !mb.IsLegit() {
		return *mb
	}
	a, aok := ma.GetNumericToFloatValue()
	b, bok := mb.GetNumericToFloatValue()
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
