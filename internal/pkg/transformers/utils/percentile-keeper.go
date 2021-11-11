// ================================================================
// TODO: comment here
// ================================================================

package utils

import (
	"fmt"
	"math"
	"sort"

	"mlr/internal/pkg/types"
)

type PercentileKeeper struct {
	data                      []*types.Mlrval
	sorted                    bool
	doInterpolatedPercentiles bool
}

// ----------------------------------------------------------------
func NewPercentileKeeper(
	doInterpolatedPercentiles bool,
) *PercentileKeeper {
	return &PercentileKeeper{
		data:                      make([]*types.Mlrval, 0, 1000),
		sorted:                    false,
		doInterpolatedPercentiles: doInterpolatedPercentiles,
	}
}

func (keeper *PercentileKeeper) Reset() {
	keeper.data = make([]*types.Mlrval, 0, 1000)
	keeper.sorted = false
}

// ----------------------------------------------------------------
func (keeper *PercentileKeeper) Ingest(value *types.Mlrval) {
	if len(keeper.data) >= cap(keeper.data) {
		newData := make([]*types.Mlrval, len(keeper.data), 2*cap(keeper.data))
		copy(newData, keeper.data)
		keeper.data = newData
	}

	n := len(keeper.data)

	keeper.data = keeper.data[0 : n+1]
	keeper.data[n] = value.Copy()

	keeper.sorted = false
}

// ================================================================
// Non-interpolated percentiles (see also https://en.wikipedia.org/wiki/Percentile)

// ----------------------------------------------------------------
// OPTION 1: int index = p*n/100.0;
//
// x
// 0
// 20
// 40
// 60
// 80
// 100
//
// x_p00 0 x_p10  0 x_p20 20 x_p30 20 x_p40 40 x_p50 60 x_p60 60 x_p70 80 x_p80  80 x_p90 100 x_p100 100
// x_p01 0 x_p11  0 x_p21 20 x_p31 20 x_p41 40 x_p51 60 x_p61 60 x_p71 80 x_p81  80 x_p91 100
// x_p02 0 x_p12  0 x_p22 20 x_p32 20 x_p42 40 x_p52 60 x_p62 60 x_p72 80 x_p82  80 x_p92 100
// x_p03 0 x_p13  0 x_p23 20 x_p33 20 x_p43 40 x_p53 60 x_p63 60 x_p73 80 x_p83  80 x_p93 100
// x_p04 0 x_p14  0 x_p24 20 x_p34 40 x_p44 40 x_p54 60 x_p64 60 x_p74 80 x_p84 100 x_p94 100
// x_p05 0 x_p15  0 x_p25 20 x_p35 40 x_p45 40 x_p55 60 x_p65 60 x_p75 80 x_p85 100 x_p95 100
// x_p06 0 x_p16  0 x_p26 20 x_p36 40 x_p46 40 x_p56 60 x_p66 60 x_p76 80 x_p86 100 x_p96 100
// x_p07 0 x_p17 20 x_p27 20 x_p37 40 x_p47 40 x_p57 60 x_p67 80 x_p77 80 x_p87 100 x_p97 100
// x_p08 0 x_p18 20 x_p28 20 x_p38 40 x_p48 40 x_p58 60 x_p68 80 x_p78 80 x_p88 100 x_p98 100
// x_p09 0 x_p19 20 x_p29 20 x_p39 40 x_p49 40 x_p59 60 x_p69 80 x_p79 80 x_p89 100 x_p99 100
//
// x
// 0
// 25
// 50
// 75
// 100
//
// x_p00 0 x_p10 0 x_p20 25 x_p30 25 x_p40 50 x_p50 50 x_p60 75 x_p70 75 x_p80 100 x_p90 100 x_p100 100
// x_p01 0 x_p11 0 x_p21 25 x_p31 25 x_p41 50 x_p51 50 x_p61 75 x_p71 75 x_p81 100 x_p91 100
// x_p02 0 x_p12 0 x_p22 25 x_p32 25 x_p42 50 x_p52 50 x_p62 75 x_p72 75 x_p82 100 x_p92 100
// x_p03 0 x_p13 0 x_p23 25 x_p33 25 x_p43 50 x_p53 50 x_p63 75 x_p73 75 x_p83 100 x_p93 100
// x_p04 0 x_p14 0 x_p24 25 x_p34 25 x_p44 50 x_p54 50 x_p64 75 x_p74 75 x_p84 100 x_p94 100
// x_p05 0 x_p15 0 x_p25 25 x_p35 25 x_p45 50 x_p55 50 x_p65 75 x_p75 75 x_p85 100 x_p95 100
// x_p06 0 x_p16 0 x_p26 25 x_p36 25 x_p46 50 x_p56 50 x_p66 75 x_p76 75 x_p86 100 x_p96 100
// x_p07 0 x_p17 0 x_p27 25 x_p37 25 x_p47 50 x_p57 50 x_p67 75 x_p77 75 x_p87 100 x_p97 100
// x_p08 0 x_p18 0 x_p28 25 x_p38 25 x_p48 50 x_p58 50 x_p68 75 x_p78 75 x_p88 100 x_p98 100
// x_p09 0 x_p19 0 x_p29 25 x_p39 25 x_p49 50 x_p59 50 x_p69 75 x_p79 75 x_p89 100 x_p99 100
//
// ----------------------------------------------------------------
// OPTION 2: int index = p*(n-1)/100.0;
//
// x
// 0
// 20
// 40
// 60
// 80
// 100
//
// x_p00 0 x_p10 0 x_p20 20 x_p30 20 x_p40 40 x_p50 40 x_p60 60 x_p70 60 x_p80 80 x_p90 80 x_p100 100
// x_p01 0 x_p11 0 x_p21 20 x_p31 20 x_p41 40 x_p51 40 x_p61 60 x_p71 60 x_p81 80 x_p91 80
// x_p02 0 x_p12 0 x_p22 20 x_p32 20 x_p42 40 x_p52 40 x_p62 60 x_p72 60 x_p82 80 x_p92 80
// x_p03 0 x_p13 0 x_p23 20 x_p33 20 x_p43 40 x_p53 40 x_p63 60 x_p73 60 x_p83 80 x_p93 80
// x_p04 0 x_p14 0 x_p24 20 x_p34 20 x_p44 40 x_p54 40 x_p64 60 x_p74 60 x_p84 80 x_p94 80
// x_p05 0 x_p15 0 x_p25 20 x_p35 20 x_p45 40 x_p55 40 x_p65 60 x_p75 60 x_p85 80 x_p95 80
// x_p06 0 x_p16 0 x_p26 20 x_p36 20 x_p46 40 x_p56 40 x_p66 60 x_p76 60 x_p86 80 x_p96 80
// x_p07 0 x_p17 0 x_p27 20 x_p37 20 x_p47 40 x_p57 40 x_p67 60 x_p77 60 x_p87 80 x_p97 80
// x_p08 0 x_p18 0 x_p28 20 x_p38 20 x_p48 40 x_p58 40 x_p68 60 x_p78 60 x_p88 80 x_p98 80
// x_p09 0 x_p19 0 x_p29 20 x_p39 20 x_p49 40 x_p59 40 x_p69 60 x_p79 60 x_p89 80 x_p99 80
//
// x
// 0
// 25
// 50
// 75
// 100
//
// x_p00 0 x_p10 0 x_p20  0 x_p30 25 x_p40 25 x_p50 50 x_p60 50 x_p70 50 x_p80 75 x_p90 75 x_p100 100
// x_p01 0 x_p11 0 x_p21  0 x_p31 25 x_p41 25 x_p51 50 x_p61 50 x_p71 50 x_p81 75 x_p91 75
// x_p02 0 x_p12 0 x_p22  0 x_p32 25 x_p42 25 x_p52 50 x_p62 50 x_p72 50 x_p82 75 x_p92 75
// x_p03 0 x_p13 0 x_p23  0 x_p33 25 x_p43 25 x_p53 50 x_p63 50 x_p73 50 x_p83 75 x_p93 75
// x_p04 0 x_p14 0 x_p24  0 x_p34 25 x_p44 25 x_p54 50 x_p64 50 x_p74 50 x_p84 75 x_p94 75
// x_p05 0 x_p15 0 x_p25 25 x_p35 25 x_p45 25 x_p55 50 x_p65 50 x_p75 75 x_p85 75 x_p95 75
// x_p06 0 x_p16 0 x_p26 25 x_p36 25 x_p46 25 x_p56 50 x_p66 50 x_p76 75 x_p86 75 x_p96 75
// x_p07 0 x_p17 0 x_p27 25 x_p37 25 x_p47 25 x_p57 50 x_p67 50 x_p77 75 x_p87 75 x_p97 75
// x_p08 0 x_p18 0 x_p28 25 x_p38 25 x_p48 25 x_p58 50 x_p68 50 x_p78 75 x_p88 75 x_p98 75
// x_p09 0 x_p19 0 x_p29 25 x_p39 25 x_p49 25 x_p59 50 x_p69 50 x_p79 75 x_p89 75 x_p99 75
//
// ----------------------------------------------------------------
// OPTION 3: int index = (int)ceil(p*(n-1)/100.0);
//
// x
// 0
// 20
// 40
// 60
// 80
// 100
//
// x_p00  0 x_p10 20 x_p20 20 x_p30 40 x_p40 40 x_p50 60 x_p60 60 x_p70 80 x_p80  80 x_p90 100 x_p100 100
// x_p01 20 x_p11 20 x_p21 40 x_p31 40 x_p41 60 x_p51 60 x_p61 80 x_p71 80 x_p81 100 x_p91 100
// x_p02 20 x_p12 20 x_p22 40 x_p32 40 x_p42 60 x_p52 60 x_p62 80 x_p72 80 x_p82 100 x_p92 100
// x_p03 20 x_p13 20 x_p23 40 x_p33 40 x_p43 60 x_p53 60 x_p63 80 x_p73 80 x_p83 100 x_p93 100
// x_p04 20 x_p14 20 x_p24 40 x_p34 40 x_p44 60 x_p54 60 x_p64 80 x_p74 80 x_p84 100 x_p94 100
// x_p05 20 x_p15 20 x_p25 40 x_p35 40 x_p45 60 x_p55 60 x_p65 80 x_p75 80 x_p85 100 x_p95 100
// x_p06 20 x_p16 20 x_p26 40 x_p36 40 x_p46 60 x_p56 60 x_p66 80 x_p76 80 x_p86 100 x_p96 100
// x_p07 20 x_p17 20 x_p27 40 x_p37 40 x_p47 60 x_p57 60 x_p67 80 x_p77 80 x_p87 100 x_p97 100
// x_p08 20 x_p18 20 x_p28 40 x_p38 40 x_p48 60 x_p58 60 x_p68 80 x_p78 80 x_p88 100 x_p98 100
// x_p09 20 x_p19 20 x_p29 40 x_p39 40 x_p49 60 x_p59 60 x_p69 80 x_p79 80 x_p89 100 x_p99 100
//
// x
// 0
// 25
// 50
// 75
// 100
//
// x_p00  0 x_p10 25 x_p20 25 x_p30 50 x_p40 50 x_p50 50 x_p60 75 x_p70  75 x_p80 100 x_p90 100 x_p100 100
// x_p01 25 x_p11 25 x_p21 25 x_p31 50 x_p41 50 x_p51 75 x_p61 75 x_p71  75 x_p81 100 x_p91 100
// x_p02 25 x_p12 25 x_p22 25 x_p32 50 x_p42 50 x_p52 75 x_p62 75 x_p72  75 x_p82 100 x_p92 100
// x_p03 25 x_p13 25 x_p23 25 x_p33 50 x_p43 50 x_p53 75 x_p63 75 x_p73  75 x_p83 100 x_p93 100
// x_p04 25 x_p14 25 x_p24 25 x_p34 50 x_p44 50 x_p54 75 x_p64 75 x_p74  75 x_p84 100 x_p94 100
// x_p05 25 x_p15 25 x_p25 25 x_p35 50 x_p45 50 x_p55 75 x_p65 75 x_p75  75 x_p85 100 x_p95 100
// x_p06 25 x_p16 25 x_p26 50 x_p36 50 x_p46 50 x_p56 75 x_p66 75 x_p76 100 x_p86 100 x_p96 100
// x_p07 25 x_p17 25 x_p27 50 x_p37 50 x_p47 50 x_p57 75 x_p67 75 x_p77 100 x_p87 100 x_p97 100
// x_p08 25 x_p18 25 x_p28 50 x_p38 50 x_p48 50 x_p58 75 x_p68 75 x_p78 100 x_p88 100 x_p98 100
// x_p09 25 x_p19 25 x_p29 50 x_p39 50 x_p49 50 x_p59 75 x_p69 75 x_p79 100 x_p89 100 x_p99 100
//
// ----------------------------------------------------------------
// OPTION 4: int index = (int)ceil(-0.5 + p*(n-1)/100.0);
//
// x
// 0
// 20
// 40
// 60
// 80
// 100
//
// x_p00 0 x_p10  0 x_p20 20 x_p30 20 x_p40 40 x_p50 40 x_p60 60 x_p70 60 x_p80 80 x_p90  80 x_p100 100
// x_p01 0 x_p11 20 x_p21 20 x_p31 40 x_p41 40 x_p51 60 x_p61 60 x_p71 80 x_p81 80 x_p91 100
// x_p02 0 x_p12 20 x_p22 20 x_p32 40 x_p42 40 x_p52 60 x_p62 60 x_p72 80 x_p82 80 x_p92 100
// x_p03 0 x_p13 20 x_p23 20 x_p33 40 x_p43 40 x_p53 60 x_p63 60 x_p73 80 x_p83 80 x_p93 100
// x_p04 0 x_p14 20 x_p24 20 x_p34 40 x_p44 40 x_p54 60 x_p64 60 x_p74 80 x_p84 80 x_p94 100
// x_p05 0 x_p15 20 x_p25 20 x_p35 40 x_p45 40 x_p55 60 x_p65 60 x_p75 80 x_p85 80 x_p95 100
// x_p06 0 x_p16 20 x_p26 20 x_p36 40 x_p46 40 x_p56 60 x_p66 60 x_p76 80 x_p86 80 x_p96 100
// x_p07 0 x_p17 20 x_p27 20 x_p37 40 x_p47 40 x_p57 60 x_p67 60 x_p77 80 x_p87 80 x_p97 100
// x_p08 0 x_p18 20 x_p28 20 x_p38 40 x_p48 40 x_p58 60 x_p68 60 x_p78 80 x_p88 80 x_p98 100
// x_p09 0 x_p19 20 x_p29 20 x_p39 40 x_p49 40 x_p59 60 x_p69 60 x_p79 80 x_p89 80 x_p99 100
//
// x
// 0
// 25
// 50
// 75
// 100
//
// x_p00 0 x_p10  0 x_p20 25 x_p30 25 x_p40 50 x_p50 50 x_p60 50 x_p70 75 x_p80  75 x_p90 100 x_p100 100
// x_p01 0 x_p11  0 x_p21 25 x_p31 25 x_p41 50 x_p51 50 x_p61 50 x_p71 75 x_p81  75 x_p91 100
// x_p02 0 x_p12  0 x_p22 25 x_p32 25 x_p42 50 x_p52 50 x_p62 50 x_p72 75 x_p82  75 x_p92 100
// x_p03 0 x_p13 25 x_p23 25 x_p33 25 x_p43 50 x_p53 50 x_p63 75 x_p73 75 x_p83  75 x_p93 100
// x_p04 0 x_p14 25 x_p24 25 x_p34 25 x_p44 50 x_p54 50 x_p64 75 x_p74 75 x_p84  75 x_p94 100
// x_p05 0 x_p15 25 x_p25 25 x_p35 25 x_p45 50 x_p55 50 x_p65 75 x_p75 75 x_p85  75 x_p95 100
// x_p06 0 x_p16 25 x_p26 25 x_p36 25 x_p46 50 x_p56 50 x_p66 75 x_p76 75 x_p86  75 x_p96 100
// x_p07 0 x_p17 25 x_p27 25 x_p37 25 x_p47 50 x_p57 50 x_p67 75 x_p77 75 x_p87  75 x_p97 100
// x_p08 0 x_p18 25 x_p28 25 x_p38 50 x_p48 50 x_p58 50 x_p68 75 x_p78 75 x_p88 100 x_p98 100
// x_p09 0 x_p19 25 x_p29 25 x_p39 50 x_p49 50 x_p59 50 x_p69 75 x_p79 75 x_p89 100 x_p99 100
//
// ----------------------------------------------------------------
// CONCLUSION:
// * I like option 2 for its simplicity ...
// * ... but option 1 matches R's quantile with type=1.
// * (Note that Miller's interpolated percentiles match match R's quantile with type=7)
// ----------------------------------------------------------------

func computeIndexNoninterpolated(n int, p float64) int {
	index := int(p * float64(n) / 100.0)
	//index := p * (float64(float64(n)) - 1) / 100.0
	//index := int(ceil(p * (float64(n) - 1) / 100.0))
	//index := int(ceil(-0.5 + p*(float64(n)-1)/100.0))
	if index >= n {
		index = n - 1
	}
	if index < 0 {
		index = 0
	}
	return index
}

// xxx pending pointer-output refactor
func getPercentileLinearlyInterpolated(array []*types.Mlrval, n int, p float64) types.Mlrval {
	findex := (p / 100.0) * (float64(n) - 1)
	if findex < 0.0 {
		findex = 0.0
	}
	iindex := int(math.Floor(findex))
	if iindex >= n-1 {
		return *array[iindex].Copy()
	} else {
		// array[iindex] + frac * (array[iindex+1] - array[iindex])
		// TODO: just do this in float64.
		frac := types.MlrvalFromFloat64(findex - float64(iindex))
		diff := types.BIF_minus_binary(array[iindex+1], array[iindex])
		prod := types.BIF_times(frac, diff)
		return *types.BIF_plus_binary(array[iindex], prod)
	}
}

// ----------------------------------------------------------------
func (keeper *PercentileKeeper) sortIfNecessary() {
	if !keeper.sorted {
		sort.Slice(keeper.data, func(i, j int) bool {
			return types.MlrvalLessThanAsBool(keeper.data[i], keeper.data[j])
		})
		keeper.sorted = true
	}
}

// ----------------------------------------------------------------
func (keeper *PercentileKeeper) Emit(percentile float64) *types.Mlrval {
	if keeper.doInterpolatedPercentiles {
		return keeper.EmitLinearlyInterpolated(percentile)
	} else {
		return keeper.EmitNonInterpolated(percentile)
	}
}

func (keeper *PercentileKeeper) EmitNonInterpolated(percentile float64) *types.Mlrval {
	if len(keeper.data) == 0 {
		return types.MLRVAL_VOID
	}
	keeper.sortIfNecessary()
	return keeper.data[computeIndexNoninterpolated(int(len(keeper.data)), percentile)].Copy()
}

func (keeper *PercentileKeeper) EmitLinearlyInterpolated(percentile float64) *types.Mlrval {
	if len(keeper.data) == 0 {
		return types.MLRVAL_VOID
	}
	keeper.sortIfNecessary()
	output := getPercentileLinearlyInterpolated(keeper.data, int(len(keeper.data)), percentile)
	return output.Copy()
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
