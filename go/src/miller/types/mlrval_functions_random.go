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
	return MlrvalFromInt64(
		int64(
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
	return MlrvalFromInt64(u)
}

func MlrvalUrandRange(ma, mb *Mlrval) Mlrval {
	if !ma.IsLegit() {
		return *ma
	}
	if !mb.IsLegit() {
		return *mb
	}
	a, aok := ma.GetFloatValue()
	b, bok := mb.GetFloatValue()
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
