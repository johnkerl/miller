// ================================================================
// For stats1 as well as merge-fields
// ================================================================

package transformers

import (
	"fmt"
	"os"

	"miller/types"
)

// ----------------------------------------------------------------
type IStats1Accumulator interface {
	Ingest(value *types.Mlrval)
	Emit() types.Mlrval
}

// ----------------------------------------------------------------
type newStats1AccumulatorFunc func() IStats1Accumulator

type stats1AccumulatorInfo struct {
	name        string
	description string
	constructor newStats1AccumulatorFunc
}

var stats1AccumulatorInfos []stats1AccumulatorInfo = []stats1AccumulatorInfo{

	{
		"count",
		"Count instances of fields",
		NewStats1CountAccumulator,
	},

	{
		"mode",
		"Find most-frequently-occurring values for fields; first-found wins tie",
		NewStats1ModeAccumulator,
	},

	{
		"antimode",
		"Find least-frequently-occurring values for fields; first-found wins tie",
		NewStats1AntimodeAccumulator,
	},

	{
		"sum",
		"Compute sums of specified fields",
		NewStats1SumAccumulator,
	},
	{
		"mean",
		"Compute averages (sample means) of specified fields",
		NewStats1MeanAccumulator,
	},

	{
		"var",
		"Compute sample variance of specified fields",
		NewStats1VarAccumulator,
	},
	{
		"stddev",
		"Compute sample standard deviation of specified fields",
		NewStats1StddevAccumulator,
	},
	{
		"meaneb",
		"Estimate error bars for averages (assuming no sample autocorrelation)",
		NewStats1MeanEBAccumulator,
	},
	{
		"skewness",
		"Compute sample skewness of specified fields",
		NewStats1SkewnessAccumulator,
	},
	{
		"kurtosis",
		"Compute sample kurtosis of specified fields",
		NewStats1KurtosisAccumulator,
	},

	{
		"min",
		"Compute minimum values of specified fields",
		NewStats1MinAccumulator,
	},
	{
		"max",
		"Compute maximum values of specified fields",
		NewStats1MaxAccumulator,
	},
}

// ----------------------------------------------------------------
type Stats1AccumulatorFactory struct {
	// If we are asked for p90 and p95 on the same column, we reuse the
	// percentile-keeper object to reduce runtime memory consumption.
	percentileKeepers map[string]*PercentileKeeper
}

func NewStats1AccumulatorFactory() *Stats1AccumulatorFactory {
	return &Stats1AccumulatorFactory{
		percentileKeepers: make(map[string]*PercentileKeeper),
	}
}

// ----------------------------------------------------------------
func listStats1Accumulators(o *os.File) {
	for _, info := range stats1AccumulatorInfos {
		fmt.Fprintf(o, "  %-8s %s\n", info.name, info.description)
	}
}

// ----------------------------------------------------------------
func (this *Stats1AccumulatorFactory) Make(
	accumulatorName string,
	valueFieldName string,
) IStats1Accumulator {
	if accumulatorName == "median" {
		// TODO: if ssccanf p%f ...
		percentileKeeper := this.percentileKeepers[valueFieldName]
		if percentileKeeper == nil {
			// xxx have remembered interpolation preference
			percentileKeeper = NewPercentileKeeper()
			this.percentileKeepers[valueFieldName] = percentileKeeper
		}
		return NewStats1PercentileAccumulator(percentileKeeper, accumulatorName)
	}

	for _, info := range stats1AccumulatorInfos {
		if info.name == accumulatorName {
			return info.constructor()
		}
	}
	return nil
}

// ================================================================
type Stats1CountAccumulator struct {
	count int64
}

func NewStats1CountAccumulator() IStats1Accumulator {
	return &Stats1CountAccumulator{
		count: 0,
	}
}
func (this *Stats1CountAccumulator) Ingest(value *types.Mlrval) {
	this.count++
}
func (this *Stats1CountAccumulator) Emit() types.Mlrval {
	return types.MlrvalFromInt64(this.count)
}

// ----------------------------------------------------------------
type Stats1ModeAccumulator struct {
	countsByValue map[string]int64
}

func NewStats1ModeAccumulator() IStats1Accumulator {
	return &Stats1ModeAccumulator{
		countsByValue: make(map[string]int64),
	}
}
func (this *Stats1ModeAccumulator) Ingest(value *types.Mlrval) {
	key := value.String() // 1, 1.0, and 1.000 are distinct
	this.countsByValue[key] = this.countsByValue[key] + 1
}
func (this *Stats1ModeAccumulator) Emit() types.Mlrval {
	if len(this.countsByValue) == 0 {
		return types.MlrvalFromError()
	}
	maxValue := ""
	var maxCount = int64(0)
	for value, count := range this.countsByValue {
		if maxValue == "" || count > maxCount {
			maxValue = value
			maxCount = count
		}
	}
	return types.MlrvalFromString(maxValue)
}

// ----------------------------------------------------------------
type Stats1AntimodeAccumulator struct {
	countsByValue map[string]int64
}

func NewStats1AntimodeAccumulator() IStats1Accumulator {
	return &Stats1AntimodeAccumulator{
		countsByValue: make(map[string]int64),
	}
}
func (this *Stats1AntimodeAccumulator) Ingest(value *types.Mlrval) {
	key := value.String() // 1, 1.0, and 1.000 are distinct
	this.countsByValue[key] = this.countsByValue[key] + 1
}
func (this *Stats1AntimodeAccumulator) Emit() types.Mlrval {
	if len(this.countsByValue) == 0 {
		return types.MlrvalFromError()
	}
	maxValue := ""
	var maxCount = int64(0)
	for value, count := range this.countsByValue {
		if maxValue == "" || count < maxCount {
			maxValue = value
			maxCount = count
		}
	}
	return types.MlrvalFromString(maxValue)
}

// ----------------------------------------------------------------
type Stats1SumAccumulator struct {
	sum types.Mlrval
}

func NewStats1SumAccumulator() IStats1Accumulator {
	return &Stats1SumAccumulator{
		sum: types.MlrvalFromInt64(0),
	}
}
func (this *Stats1SumAccumulator) Ingest(value *types.Mlrval) {
	this.sum = types.MlrvalBinaryPlus(&this.sum, value)
}
func (this *Stats1SumAccumulator) Emit() types.Mlrval {
	return *this.sum.Copy()
}

// ----------------------------------------------------------------
type Stats1MeanAccumulator struct {
	sum   types.Mlrval
	count int64
}

func NewStats1MeanAccumulator() IStats1Accumulator {
	return &Stats1MeanAccumulator{
		sum:   types.MlrvalFromInt64(0),
		count: 0,
	}
}
func (this *Stats1MeanAccumulator) Ingest(value *types.Mlrval) {
	this.sum = types.MlrvalBinaryPlus(&this.sum, value)
	this.count++
}
func (this *Stats1MeanAccumulator) Emit() types.Mlrval {
	if this.count == 0 {
		return types.MlrvalFromVoid()
	} else {
		mcount := types.MlrvalFromInt64(this.count)
		return types.MlrvalDivide(&this.sum, &mcount)
	}
}

// ----------------------------------------------------------------
type Stats1MinAccumulator struct {
	min types.Mlrval
}

func NewStats1MinAccumulator() IStats1Accumulator {
	return &Stats1MinAccumulator{
		min: types.MlrvalFromAbsent(),
	}
}
func (this *Stats1MinAccumulator) Ingest(value *types.Mlrval) {
	this.min = types.MlrvalBinaryMin(&this.min, value)
}
func (this *Stats1MinAccumulator) Emit() types.Mlrval {
	return *this.min.Copy()
}

// ----------------------------------------------------------------
type Stats1MaxAccumulator struct {
	max types.Mlrval
}

func NewStats1MaxAccumulator() IStats1Accumulator {
	return &Stats1MaxAccumulator{
		max: types.MlrvalFromAbsent(),
	}
}
func (this *Stats1MaxAccumulator) Ingest(value *types.Mlrval) {
	this.max = types.MlrvalBinaryMax(&this.max, value)
}
func (this *Stats1MaxAccumulator) Emit() types.Mlrval {
	return *this.max.Copy()
}

// ----------------------------------------------------------------
type Stats1VarAccumulator struct {
	count int64
	sum   types.Mlrval
	sum2  types.Mlrval
}

func NewStats1VarAccumulator() IStats1Accumulator {
	return &Stats1VarAccumulator{
		count: 0,
		sum:   types.MlrvalFromInt64(0),
		sum2:  types.MlrvalFromInt64(0),
	}
}
func (this *Stats1VarAccumulator) Ingest(value *types.Mlrval) {
	value2 := types.MlrvalTimes(value, value)
	this.count++
	this.sum = types.MlrvalBinaryPlus(&this.sum, value)
	this.sum2 = types.MlrvalBinaryPlus(&this.sum2, &value2)
}
func (this *Stats1VarAccumulator) Emit() types.Mlrval {
	mcount := types.MlrvalFromInt64(this.count)
	return types.MlrvalGetVar(&mcount, &this.sum, &this.sum2)
}

// ----------------------------------------------------------------
type Stats1StddevAccumulator struct {
	count int64
	sum   types.Mlrval
	sum2  types.Mlrval
}

func NewStats1StddevAccumulator() IStats1Accumulator {
	return &Stats1StddevAccumulator{
		count: 0,
		sum:   types.MlrvalFromInt64(0),
		sum2:  types.MlrvalFromInt64(0),
	}
}
func (this *Stats1StddevAccumulator) Ingest(value *types.Mlrval) {
	value2 := types.MlrvalTimes(value, value)
	this.count++
	this.sum = types.MlrvalBinaryPlus(&this.sum, value)
	this.sum2 = types.MlrvalBinaryPlus(&this.sum2, &value2)
}
func (this *Stats1StddevAccumulator) Emit() types.Mlrval {
	mcount := types.MlrvalFromInt64(this.count)
	return types.MlrvalGetStddev(&mcount, &this.sum, &this.sum2)
}

// ----------------------------------------------------------------
type Stats1MeanEBAccumulator struct {
	count int64
	sum   types.Mlrval
	sum2  types.Mlrval
}

func NewStats1MeanEBAccumulator() IStats1Accumulator {
	return &Stats1MeanEBAccumulator{
		count: 0,
		sum:   types.MlrvalFromInt64(0),
		sum2:  types.MlrvalFromInt64(0),
	}
}
func (this *Stats1MeanEBAccumulator) Ingest(value *types.Mlrval) {
	value2 := types.MlrvalTimes(value, value)
	this.count++
	this.sum = types.MlrvalBinaryPlus(&this.sum, value)
	this.sum2 = types.MlrvalBinaryPlus(&this.sum2, &value2)
}
func (this *Stats1MeanEBAccumulator) Emit() types.Mlrval {
	mcount := types.MlrvalFromInt64(this.count)
	return types.MlrvalGetMeanEB(&mcount, &this.sum, &this.sum2)
}

// ----------------------------------------------------------------
type Stats1SkewnessAccumulator struct {
	count int64
	sum   types.Mlrval
	sum2  types.Mlrval
	sum3  types.Mlrval
}

func NewStats1SkewnessAccumulator() IStats1Accumulator {
	return &Stats1SkewnessAccumulator{
		count: 0,
		sum:   types.MlrvalFromInt64(0),
		sum2:  types.MlrvalFromInt64(0),
		sum3:  types.MlrvalFromInt64(0),
	}
}
func (this *Stats1SkewnessAccumulator) Ingest(value *types.Mlrval) {
	value2 := types.MlrvalTimes(value, value)
	value3 := types.MlrvalTimes(value, &value2)
	this.count++
	this.sum = types.MlrvalBinaryPlus(&this.sum, value)
	this.sum2 = types.MlrvalBinaryPlus(&this.sum2, &value2)
	this.sum3 = types.MlrvalBinaryPlus(&this.sum3, &value3)
}
func (this *Stats1SkewnessAccumulator) Emit() types.Mlrval {
	mcount := types.MlrvalFromInt64(this.count)
	return types.MlrvalGetSkewness(&mcount, &this.sum, &this.sum2, &this.sum3)
}

// ----------------------------------------------------------------
type Stats1KurtosisAccumulator struct {
	count int64
	sum   types.Mlrval
	sum2  types.Mlrval
	sum3  types.Mlrval
	sum4  types.Mlrval
}

func NewStats1KurtosisAccumulator() IStats1Accumulator {
	return &Stats1KurtosisAccumulator{
		count: 0,
		sum:   types.MlrvalFromInt64(0),
		sum2:  types.MlrvalFromInt64(0),
		sum3:  types.MlrvalFromInt64(0),
		sum4:  types.MlrvalFromInt64(0),
	}
}
func (this *Stats1KurtosisAccumulator) Ingest(value *types.Mlrval) {
	value2 := types.MlrvalTimes(value, value)
	value3 := types.MlrvalTimes(value, &value2)
	value4 := types.MlrvalTimes(value, &value3)
	this.count++
	this.sum = types.MlrvalBinaryPlus(&this.sum, value)
	this.sum2 = types.MlrvalBinaryPlus(&this.sum2, &value2)
	this.sum3 = types.MlrvalBinaryPlus(&this.sum3, &value3)
	this.sum4 = types.MlrvalBinaryPlus(&this.sum4, &value4)
}
func (this *Stats1KurtosisAccumulator) Emit() types.Mlrval {
	mcount := types.MlrvalFromInt64(this.count)
	return types.MlrvalGetKurtosis(&mcount, &this.sum, &this.sum2, &this.sum3, &this.sum4)
}

// ----------------------------------------------------------------
// TODO: refactor this for use with more than one percentile

type Stats1PercentileAccumulator struct {
	percentileKeeper *PercentileKeeper
	percentile float64
}

func NewStats1PercentileAccumulator(
	percentileKeeper *PercentileKeeper,
	name string, // "median" or things like "p99.9"
) IStats1Accumulator {
	return &Stats1PercentileAccumulator{
		percentileKeeper: percentileKeeper,
		percentile: 50.0, // TEMP
	}
}

func (this *Stats1PercentileAccumulator) Ingest(value *types.Mlrval) {
	this.percentileKeeper.Ingest(value)
}

func (this *Stats1PercentileAccumulator) Emit(
//	char* valueFieldName,
//	char* stats1_acc_name,
//	lrec_t* outrec,
) types.Mlrval {

	//	if (stats1_acc_name[0] == 'm') { // Pre-validated to be either p{number} or median.
	//		p = 50.0;
	//	} else {
	//		// TODO: do the sscanf once at alloc time and store the double in the state struct for a minor perf gain.
	//		(void)sscanf(stats1_acc_name, "p%lf", &p); // Assuming this was range-checked earlier on to be in [0,100].
	//	}

	// TODO: parameterize
	v := this.percentileKeeper.EmitNonInterpolated(this.percentile)

	return v
}
