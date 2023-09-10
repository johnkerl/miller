// ================================================================
// TODO: comment here
// ================================================================

package utils

import (
	"fmt"
	"sort"

	"github.com/johnkerl/miller/pkg/bifs"
	"github.com/johnkerl/miller/pkg/mlrval"
)

type PercentileKeeper struct {
	data                      []*mlrval.Mlrval
	sorted                    bool
	doInterpolatedPercentiles bool
}

// Lower outer fence, lower inner fence, upper inner fence, upper outer fence.
var fenceInnerK = mlrval.FromFloat(1.5)
var fenceOuterK = mlrval.FromFloat(3.0)

// ----------------------------------------------------------------
func NewPercentileKeeper(
	doInterpolatedPercentiles bool,
) *PercentileKeeper {
	return &PercentileKeeper{
		data:                      make([]*mlrval.Mlrval, 0, 1000),
		sorted:                    false,
		doInterpolatedPercentiles: doInterpolatedPercentiles,
	}
}

func (keeper *PercentileKeeper) Reset() {
	keeper.data = make([]*mlrval.Mlrval, 0, 1000)
	keeper.sorted = false
}

// ----------------------------------------------------------------
func (keeper *PercentileKeeper) Ingest(value *mlrval.Mlrval) {
	if len(keeper.data) >= cap(keeper.data) {
		newData := make([]*mlrval.Mlrval, len(keeper.data), 2*cap(keeper.data))
		copy(newData, keeper.data)
		keeper.data = newData
	}

	n := len(keeper.data)

	keeper.data = keeper.data[0 : n+1]
	keeper.data[n] = value.Copy()

	keeper.sorted = false
}

// ----------------------------------------------------------------
func (keeper *PercentileKeeper) sortIfNecessary() {
	if !keeper.sorted {
		sort.Slice(keeper.data, func(i, j int) bool {
			return mlrval.LessThan(keeper.data[i], keeper.data[j])
		})
		keeper.sorted = true
	}
}

// ----------------------------------------------------------------
func (keeper *PercentileKeeper) Emit(percentile float64) *mlrval.Mlrval {
	if keeper.doInterpolatedPercentiles {
		return keeper.EmitLinearlyInterpolated(percentile)
	} else {
		return keeper.EmitNonInterpolated(percentile)
	}
}

func (keeper *PercentileKeeper) EmitNonInterpolated(percentile float64) *mlrval.Mlrval {
	if len(keeper.data) == 0 {
		return mlrval.VOID
	}
	keeper.sortIfNecessary()
	return bifs.GetPercentileNonInterpolated(keeper.data, int(len(keeper.data)), percentile)
}

func (keeper *PercentileKeeper) EmitLinearlyInterpolated(percentile float64) *mlrval.Mlrval {
	if len(keeper.data) == 0 {
		return mlrval.VOID
	}
	keeper.sortIfNecessary()
	return bifs.GetPercentileLinearlyInterpolated(keeper.data, int(len(keeper.data)), percentile)
}

// ----------------------------------------------------------------
// TODO: COMMENT
func (keeper *PercentileKeeper) EmitNamed(name string) *mlrval.Mlrval {
	if name == "min" {
		return keeper.EmitNonInterpolated(0.0)
	} else if name == "p25" {
		return keeper.EmitNonInterpolated(25.0)
	} else if name == "median" {
		return keeper.EmitNonInterpolated(50.0)
	} else if name == "p75" {
		return keeper.EmitNonInterpolated(75.0)
	} else if name == "max" {
		return keeper.EmitNonInterpolated(100.0)

	} else if name == "iqr" {
		p25 := keeper.EmitNonInterpolated(25.0)
		p75 := keeper.EmitNonInterpolated(75.0)
		if p25.IsNumeric() && p75.IsNumeric() {
			return bifs.BIF_minus_binary(p75, p25)
		} else {
			return mlrval.VOID
		}

	} else if name == "lof" {
		p25 := keeper.EmitNonInterpolated(25.0)
		iqr := keeper.EmitNamed("iqr")
		if p25.IsNumeric() && iqr.IsNumeric() {
			return bifs.BIF_minus_binary(p25, bifs.BIF_times(fenceOuterK, iqr))
		} else {
			return mlrval.VOID
		}

	} else if name == "lif" {
		p25 := keeper.EmitNonInterpolated(25.0)
		iqr := keeper.EmitNamed("iqr")
		if p25.IsNumeric() && iqr.IsNumeric() {
			return bifs.BIF_minus_binary(p25, bifs.BIF_times(fenceInnerK, iqr))
		} else {
			return mlrval.VOID
		}

	} else if name == "uif" {
		p75 := keeper.EmitNonInterpolated(75.0)
		iqr := keeper.EmitNamed("iqr")
		if p75.IsNumeric() && iqr.IsNumeric() {
			return bifs.BIF_plus_binary(p75, bifs.BIF_times(fenceInnerK, iqr))
		} else {
			return mlrval.VOID
		}

	} else if name == "uof" {
		p75 := keeper.EmitNonInterpolated(75.0)
		iqr := keeper.EmitNamed("iqr")
		if p75.IsNumeric() && iqr.IsNumeric() {
			return bifs.BIF_plus_binary(p75, bifs.BIF_times(fenceOuterK, iqr))
		} else {
			return mlrval.VOID
		}

	} else {
		return mlrval.FromError(
			fmt.Errorf(`stats1: unrecognized percentilename "%s"`, name),
		)
	}
}

// ----------------------------------------------------------------
func (keeper *PercentileKeeper) Dump() {
	fmt.Printf("percentile_keeper dump:\n")
	for i, datum := range keeper.data {
		ival, ok := datum.GetIntValue()
		if ok {
			fmt.Printf("[%02d] %d\n", i, ival)
		}
		fval, ok := datum.GetFloatValue()
		if ok {
			fmt.Printf("[%02d] %.8f\n", i, fval)
		}
	}
}
