// ================================================================
// For stats1 as well as merge-fields
// ================================================================

package utils

import (
	"fmt"
	"os"
	"strings"

	"miller/src/lib"
	"miller/src/types"
)

// ----------------------------------------------------------------
type IStats1Accumulator interface {
	Ingest(value *types.Mlrval)
	Emit() *types.Mlrval
	Reset() // for merge-fields where we reset after each record instead of replace/recreate
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

// ================================================================
// TODO: comment

type Stats1NamedAccumulator struct {
	valueFieldName  string
	accumulatorName string
	accumulator     IStats1Accumulator
	outputFieldName string
}

func NewStats1NamedAccumulator(
	valueFieldName string,
	accumulatorName string,
	accumulator IStats1Accumulator,
) *Stats1NamedAccumulator {
	return &Stats1NamedAccumulator{
		valueFieldName:  valueFieldName,
		accumulatorName: accumulatorName,
		accumulator:     accumulator,
		outputFieldName: valueFieldName + "_" + accumulatorName,
	}
}

func (this *Stats1NamedAccumulator) Ingest(value *types.Mlrval) {
	this.accumulator.Ingest(value)
}

func (this *Stats1NamedAccumulator) Emit() (key string, value *types.Mlrval) {
	return this.outputFieldName, this.accumulator.Emit()
}

func (this *Stats1NamedAccumulator) Reset() {
	this.accumulator.Reset()
}

// ----------------------------------------------------------------
// If we are asked for p90 and p95 on the same column, we reuse the
// percentile-keeper object to reduce runtime memory consumption.  This
// two-level map is keyed by value-field name and grouping key.  E.g. for 'mlr
// stats1 -a median -f x,y -g a,b' there will be an entry keyed primarily by
// the string "x", and secondarily keyed by the values of a and b for a given
// record.
type Stats1AccumulatorFactory struct {
	percentileKeepers map[string]map[string]*PercentileKeeper
}

func NewStats1AccumulatorFactory() *Stats1AccumulatorFactory {
	return &Stats1AccumulatorFactory{
		percentileKeepers: make(map[string]map[string]*PercentileKeeper),
	}
}

// ----------------------------------------------------------------
func ListStats1Accumulators(o *os.File) {
	for _, info := range stats1AccumulatorInfos {
		fmt.Fprintf(o, "  %-8s %s\n", info.name, info.description)
	}
}

func ValidateStats1AccumulatorName(
	accumulatorName string,
) bool {
	// First try percentiles, which have parameterized names.
	_, ok := tryPercentileFromName(accumulatorName)
	if ok {
		return true
	}

	// Then try the lookup table.
	for _, info := range stats1AccumulatorInfos {
		if info.name == accumulatorName {
			return true
		}
	}
	return false
}

// Tries to get a percentile value from names like "p99" and "median":
// * "p95"   -> 95.0, true
// * "p99.9" -> 99.9, true
// * "p200"  -> _, false
// * "median" -> 50.0, false
// * "nonesuch" -> _, false
func tryPercentileFromName(accumulatorName string) (float64, bool) {
	if accumulatorName == "median" {
		return 50.0, true
	}
	if strings.HasPrefix(accumulatorName, "p") {
		percentile, ok := lib.TryFloat64FromString(accumulatorName[1:])
		if !ok {
			return 0.0, false
		}
		if percentile < 0.0 || percentile > 100.0 {
			return 0.0, false
		}
		return percentile, true
	}

	return 0.0, false
}

// ----------------------------------------------------------------
// For merge-fields wherein percentile-keepers are re-created on each record
func (this *Stats1AccumulatorFactory) Reset() {
	this.percentileKeepers = make(map[string]map[string]*PercentileKeeper)
}

func (this *Stats1AccumulatorFactory) MakeNamedAccumulator(
	accumulatorName string,
	groupingKey string,
	valueFieldName string,
	doInterpolatedPercentiles bool,
) *Stats1NamedAccumulator {

	accumulator := this.MakeAccumulator(
		accumulatorName,
		groupingKey,
		valueFieldName,
		doInterpolatedPercentiles,
	)
	// We don't return errors.New here. The nominal case is that the stats1
	// verb has already pre-validated accumulator names, and this is just a
	// fallback. The accumulators are instantiated for every unique combination
	// of group-by field values in the record stream, only as those values are
	// encountered: for example, with 'mlr stats1 -a count,sum -f x,y -g
	// color,shape', we make a new accumulator the first time we find a record
	// with 'color=blue,shape=square' and another the first time we find a
	// record with 'color=red,shape=circle', and so on. The right thing is to
	// pre-validate names once when the stats1 transformer is being
	// instantiated.
	lib.InternalCodingErrorIf(accumulator == nil)

	return NewStats1NamedAccumulator(
		valueFieldName,
		accumulatorName,
		accumulator,
	)
}

func (this *Stats1AccumulatorFactory) MakeAccumulator(
	accumulatorName string,
	groupingKey string,
	valueFieldName string,
	doInterpolatedPercentiles bool,
) IStats1Accumulator {
	// First try percentiles, which have parameterized names.
	percentile, ok := tryPercentileFromName(accumulatorName)
	if ok {
		percentileKeepersForValueFieldName := this.percentileKeepers[valueFieldName]
		if percentileKeepersForValueFieldName == nil {
			percentileKeepersForValueFieldName = make(map[string]*PercentileKeeper)
			this.percentileKeepers[valueFieldName] = percentileKeepersForValueFieldName
		}

		percentileKeeper := percentileKeepersForValueFieldName[groupingKey]
		isPrimary := false
		if percentileKeeper == nil {
			percentileKeeper = NewPercentileKeeper(doInterpolatedPercentiles)
			percentileKeepersForValueFieldName[groupingKey] = percentileKeeper
			isPrimary = true
		}

		// To conserve memory, percentile-keeprs on the same value-field-name
		// (and grouping-key) are shared. For example, p25,p75 on field "x".
		// This means though that each datapoint must be ingested only once
		// (e.g.  by the p25 accumulator) since it shares a percentile-keeper
		// with the p75 accumulator. We handle this by tracking the first
		// construction.
		return NewStats1PercentileAccumulator(percentileKeeper, percentile, isPrimary)
	}

	// Then try the lookup table.
	for _, info := range stats1AccumulatorInfos {
		if info.name == accumulatorName {
			return info.constructor()
		}
	}
	return nil
}

// ================================================================
type Stats1CountAccumulator struct {
	count int
}

func NewStats1CountAccumulator() IStats1Accumulator {
	return &Stats1CountAccumulator{
		count: 0,
	}
}
func (this *Stats1CountAccumulator) Ingest(value *types.Mlrval) {
	this.count++
}
func (this *Stats1CountAccumulator) Emit() *types.Mlrval {
	return types.MlrvalPointerFromInt(this.count)
}
func (this *Stats1CountAccumulator) Reset() {
	this.count = 0
}

// ----------------------------------------------------------------
type Stats1ModeAccumulator struct {
	// Needs to be an ordered map to guarantee Miller's semantics that
	// first-found breaks ties.
	countsByValue *lib.OrderedMap
}

func NewStats1ModeAccumulator() IStats1Accumulator {
	return &Stats1ModeAccumulator{
		countsByValue: lib.NewOrderedMap(),
	}
}
func (this *Stats1ModeAccumulator) Ingest(value *types.Mlrval) {
	key := value.String() // 1, 1.0, and 1.000 are distinct
	iPrevious, ok := this.countsByValue.GetWithCheck(key)
	if !ok {
		this.countsByValue.Put(key, int(1))
	} else {
		this.countsByValue.Put(key, iPrevious.(int)+1)
	}
}
func (this *Stats1ModeAccumulator) Emit() *types.Mlrval {
	if this.countsByValue.FieldCount == 0 {
		return types.MLRVAL_ERROR
	}
	maxValue := ""
	var maxCount = int(0)
	for pe := this.countsByValue.Head; pe != nil; pe = pe.Next {
		value := pe.Key
		count := pe.Value.(int)
		if maxValue == "" || count > maxCount {
			maxValue = value
			maxCount = count
		}
	}
	return types.MlrvalPointerFromString(maxValue)
}
func (this *Stats1ModeAccumulator) Reset() {
	this.countsByValue = lib.NewOrderedMap()
}

// ----------------------------------------------------------------
type Stats1AntimodeAccumulator struct {
	// Needs to be an ordered map to guarantee Miller's semantics that
	// first-found breaks ties.
	countsByValue *lib.OrderedMap
}

func NewStats1AntimodeAccumulator() IStats1Accumulator {
	return &Stats1AntimodeAccumulator{
		countsByValue: lib.NewOrderedMap(),
	}
}
func (this *Stats1AntimodeAccumulator) Ingest(value *types.Mlrval) {
	key := value.String() // 1, 1.0, and 1.000 are distinct
	iPrevious, ok := this.countsByValue.GetWithCheck(key)
	if !ok {
		this.countsByValue.Put(key, int(1))
	} else {
		this.countsByValue.Put(key, iPrevious.(int)+1)
	}
}
func (this *Stats1AntimodeAccumulator) Emit() *types.Mlrval {
	if this.countsByValue.FieldCount == 0 {
		return types.MLRVAL_ERROR
	}
	minValue := ""
	var minCount = int(0)
	for pe := this.countsByValue.Head; pe != nil; pe = pe.Next {
		value := pe.Key
		count := pe.Value.(int)
		if minValue == "" || count < minCount {
			minValue = value
			minCount = count
		}
	}
	return types.MlrvalPointerFromString(minValue)
}
func (this *Stats1AntimodeAccumulator) Reset() {
	this.countsByValue = lib.NewOrderedMap()
}

// ----------------------------------------------------------------
type Stats1SumAccumulator struct {
	sum *types.Mlrval
}

func NewStats1SumAccumulator() IStats1Accumulator {
	return &Stats1SumAccumulator{
		sum: types.MlrvalPointerFromInt(0),
	}
}
func (this *Stats1SumAccumulator) Ingest(value *types.Mlrval) {
	this.sum = types.MlrvalBinaryPlus(this.sum, value)
}
func (this *Stats1SumAccumulator) Emit() *types.Mlrval {
	return this.sum.Copy()
}
func (this *Stats1SumAccumulator) Reset() {
	this.sum = types.MlrvalPointerFromInt(0)
}

// ----------------------------------------------------------------
type Stats1MeanAccumulator struct {
	sum   *types.Mlrval
	count int
}

func NewStats1MeanAccumulator() IStats1Accumulator {
	return &Stats1MeanAccumulator{
		sum:   types.MlrvalPointerFromInt(0),
		count: 0,
	}
}
func (this *Stats1MeanAccumulator) Ingest(value *types.Mlrval) {
	this.sum = types.MlrvalBinaryPlus(this.sum, value)
	this.count++
}
func (this *Stats1MeanAccumulator) Emit() *types.Mlrval {
	if this.count == 0 {
		return types.MLRVAL_VOID
	} else {
		return types.MlrvalDivide(this.sum, types.MlrvalPointerFromInt(this.count))
	}
}
func (this *Stats1MeanAccumulator) Reset() {
	this.sum = types.MlrvalPointerFromInt(0)
	this.count = 0
}

// ----------------------------------------------------------------
type Stats1MinAccumulator struct {
	min *types.Mlrval
}

func NewStats1MinAccumulator() IStats1Accumulator {
	return &Stats1MinAccumulator{
		min: types.MLRVAL_ABSENT,
	}
}
func (this *Stats1MinAccumulator) Ingest(value *types.Mlrval) {
	this.min = types.MlrvalBinaryMin(this.min, value)
}
func (this *Stats1MinAccumulator) Emit() *types.Mlrval {
	if this.min.IsAbsent() {
		return types.MLRVAL_VOID
	} else {
		return this.min.Copy()
	}
}
func (this *Stats1MinAccumulator) Reset() {
	this.min = types.MLRVAL_ABSENT
}

// ----------------------------------------------------------------
type Stats1MaxAccumulator struct {
	max *types.Mlrval
}

func NewStats1MaxAccumulator() IStats1Accumulator {
	return &Stats1MaxAccumulator{
		max: types.MLRVAL_ABSENT,
	}
}
func (this *Stats1MaxAccumulator) Ingest(value *types.Mlrval) {
	this.max = types.MlrvalBinaryMax(this.max, value)
}
func (this *Stats1MaxAccumulator) Emit() *types.Mlrval {
	if this.max.IsAbsent() {
		return types.MLRVAL_VOID
	} else {
		return this.max.Copy()
	}
}
func (this *Stats1MaxAccumulator) Reset() {
	this.max = types.MLRVAL_ABSENT
}

// ----------------------------------------------------------------
type Stats1VarAccumulator struct {
	count int
	sum   *types.Mlrval
	sum2  *types.Mlrval
}

func NewStats1VarAccumulator() IStats1Accumulator {
	return &Stats1VarAccumulator{
		count: 0,
		sum:   types.MlrvalPointerFromInt(0),
		sum2:  types.MlrvalPointerFromInt(0),
	}
}
func (this *Stats1VarAccumulator) Ingest(value *types.Mlrval) {
	value2 := types.MlrvalTimes(value, value)
	this.count++
	this.sum = types.MlrvalBinaryPlus(this.sum, value)
	this.sum2 = types.MlrvalBinaryPlus(this.sum2, value2)
}
func (this *Stats1VarAccumulator) Emit() *types.Mlrval {
	mcount := types.MlrvalFromInt(this.count)
	return types.MlrvalGetVar(&mcount, this.sum, this.sum2)
}
func (this *Stats1VarAccumulator) Reset() {
	this.count = 0
	this.sum = types.MlrvalPointerFromInt(0)
	this.sum2 = types.MlrvalPointerFromInt(0)
}

// ----------------------------------------------------------------
type Stats1StddevAccumulator struct {
	count int
	sum   *types.Mlrval
	sum2  *types.Mlrval
}

func NewStats1StddevAccumulator() IStats1Accumulator {
	return &Stats1StddevAccumulator{
		count: 0,
		sum:   types.MlrvalPointerFromInt(0),
		sum2:  types.MlrvalPointerFromInt(0),
	}
}
func (this *Stats1StddevAccumulator) Ingest(value *types.Mlrval) {
	value2 := types.MlrvalTimes(value, value)
	this.count++
	this.sum = types.MlrvalBinaryPlus(this.sum, value)
	this.sum2 = types.MlrvalBinaryPlus(this.sum2, value2)
}
func (this *Stats1StddevAccumulator) Emit() *types.Mlrval {
	mcount := types.MlrvalFromInt(this.count)
	return types.MlrvalGetStddev(&mcount, this.sum, this.sum2)
}
func (this *Stats1StddevAccumulator) Reset() {
	this.count = 0
	this.sum = types.MlrvalPointerFromInt(0)
	this.sum2 = types.MlrvalPointerFromInt(0)
}

// ----------------------------------------------------------------
type Stats1MeanEBAccumulator struct {
	count int
	sum   *types.Mlrval
	sum2  *types.Mlrval
}

func NewStats1MeanEBAccumulator() IStats1Accumulator {
	return &Stats1MeanEBAccumulator{
		count: 0,
		sum:   types.MlrvalPointerFromInt(0),
		sum2:  types.MlrvalPointerFromInt(0),
	}
}
func (this *Stats1MeanEBAccumulator) Ingest(value *types.Mlrval) {
	value2 := types.MlrvalTimes(value, value)
	this.count++
	this.sum = types.MlrvalBinaryPlus(this.sum, value)
	this.sum2 = types.MlrvalBinaryPlus(this.sum2, value2)
}
func (this *Stats1MeanEBAccumulator) Emit() *types.Mlrval {
	mcount := types.MlrvalPointerFromInt(this.count)
	return types.MlrvalGetMeanEB(mcount, this.sum, this.sum2)
}
func (this *Stats1MeanEBAccumulator) Reset() {
	this.count = 0
	this.sum = types.MlrvalPointerFromInt(0)
	this.sum2 = types.MlrvalPointerFromInt(0)
}

// ----------------------------------------------------------------
type Stats1SkewnessAccumulator struct {
	count int
	sum   *types.Mlrval
	sum2  *types.Mlrval
	sum3  *types.Mlrval
}

func NewStats1SkewnessAccumulator() IStats1Accumulator {
	return &Stats1SkewnessAccumulator{
		count: 0,
		sum:   types.MlrvalPointerFromInt(0),
		sum2:  types.MlrvalPointerFromInt(0),
		sum3:  types.MlrvalPointerFromInt(0),
	}
}
func (this *Stats1SkewnessAccumulator) Ingest(value *types.Mlrval) {
	value2 := types.MlrvalTimes(value, value)
	value3 := types.MlrvalTimes(value, value2)
	this.count++
	this.sum = types.MlrvalBinaryPlus(this.sum, value)
	this.sum2 = types.MlrvalBinaryPlus(this.sum2, value2)
	this.sum3 = types.MlrvalBinaryPlus(this.sum3, value3)
}
func (this *Stats1SkewnessAccumulator) Emit() *types.Mlrval {
	mcount := types.MlrvalPointerFromInt(this.count)
	return types.MlrvalGetSkewness(mcount, this.sum, this.sum2, this.sum3)
}
func (this *Stats1SkewnessAccumulator) Reset() {
	this.count = 0
	this.sum = types.MlrvalPointerFromInt(0)
	this.sum2 = types.MlrvalPointerFromInt(0)
	this.sum3 = types.MlrvalPointerFromInt(0)
}

// ----------------------------------------------------------------
type Stats1KurtosisAccumulator struct {
	count int
	sum   *types.Mlrval
	sum2  *types.Mlrval
	sum3  *types.Mlrval
	sum4  *types.Mlrval
}

func NewStats1KurtosisAccumulator() IStats1Accumulator {
	return &Stats1KurtosisAccumulator{
		count: 0,
		sum:   types.MlrvalPointerFromInt(0),
		sum2:  types.MlrvalPointerFromInt(0),
		sum3:  types.MlrvalPointerFromInt(0),
		sum4:  types.MlrvalPointerFromInt(0),
	}
}
func (this *Stats1KurtosisAccumulator) Ingest(value *types.Mlrval) {
	value2 := types.MlrvalTimes(value, value)
	value3 := types.MlrvalTimes(value, value2)
	value4 := types.MlrvalTimes(value, value3)
	this.count++
	this.sum = types.MlrvalBinaryPlus(this.sum, value)
	this.sum2 = types.MlrvalBinaryPlus(this.sum2, value2)
	this.sum3 = types.MlrvalBinaryPlus(this.sum3, value3)
	this.sum4 = types.MlrvalBinaryPlus(this.sum4, value4)
}
func (this *Stats1KurtosisAccumulator) Emit() *types.Mlrval {
	mcount := types.MlrvalPointerFromInt(this.count)
	return types.MlrvalGetKurtosis(mcount, this.sum, this.sum2, this.sum3, this.sum4)
}
func (this *Stats1KurtosisAccumulator) Reset() {
	this.count = 0
	this.sum = types.MlrvalPointerFromInt(0)
	this.sum2 = types.MlrvalPointerFromInt(0)
	this.sum3 = types.MlrvalPointerFromInt(0)
	this.sum4 = types.MlrvalPointerFromInt(0)
}

// ----------------------------------------------------------------
// To conserve memory, percentile-keeprs on the same value-field-name (and
// grouping-key) are shared. For example, p25,p75 on field "x".  This means
// though that each datapoint must be ingested only once (e.g.  by the p25
// accumulator) since it shares a percentile-keepr with the p75 accumulator.
// The isPrimary flag tracks this.
type Stats1PercentileAccumulator struct {
	percentileKeeper *PercentileKeeper
	percentile       float64
	isPrimary        bool
}

func NewStats1PercentileAccumulator(
	percentileKeeper *PercentileKeeper,
	percentile float64,
	isPrimary bool,
) IStats1Accumulator {
	return &Stats1PercentileAccumulator{
		percentileKeeper: percentileKeeper,
		percentile:       percentile,
		isPrimary:        isPrimary,
	}
}

func (this *Stats1PercentileAccumulator) Ingest(value *types.Mlrval) {
	if this.isPrimary {
		this.percentileKeeper.Ingest(value)
	}
}

func (this *Stats1PercentileAccumulator) Emit() *types.Mlrval {
	return this.percentileKeeper.Emit(this.percentile)
}

func (this *Stats1PercentileAccumulator) Reset() {
	if this.isPrimary {
		this.percentileKeeper.Reset()
	}
}
