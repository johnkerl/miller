package bifs

import (
	"math"
	"sort"

	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

// ----------------------------------------------------------------
// We would need a second pass through the data to compute the error-bars given
// the data and the m and the b.
//
//	# Young 1962, pp. 122-124.  Compute sample variance of linear
//	# approximations, then variances of m and b.
//	var_z = 0.0
//	for i in range(0, N):
//		var_z += (m * xs[i] + b - ys[i])**2
//	var_z /= N
//
//	var_m = (N * var_z) / D
//	var_b = (var_z * sumx2) / D
//
//	output = [m, b, math.sqrt(var_m), math.sqrt(var_b)]

// ----------------------------------------------------------------
func BIF_finalize_variance(mn, msum, msum2 *mlrval.Mlrval) *mlrval.Mlrval {
	n, isInt := mn.GetIntValue()
	lib.InternalCodingErrorIf(!isInt)
	sum, isNumber := msum.GetNumericToFloatValue()
	lib.InternalCodingErrorIf(!isNumber)
	sum2, isNumber := msum2.GetNumericToFloatValue()
	lib.InternalCodingErrorIf(!isNumber)

	if n < 2 {
		return mlrval.VOID
	}

	mean := float64(sum) / float64(n)
	numerator := sum2 - mean*(2.0*sum-float64(n)*mean)
	if numerator < 0.0 { // round-off error
		numerator = 0.0
	}
	denominator := float64(n - 1)
	return mlrval.FromFloat(numerator / denominator)
}

// ----------------------------------------------------------------
func BIF_finalize_stddev(mn, msum, msum2 *mlrval.Mlrval) *mlrval.Mlrval {
	mvar := BIF_finalize_variance(mn, msum, msum2)
	if mvar.IsVoid() {
		return mvar
	}
	return BIF_sqrt(mvar)
}

// ----------------------------------------------------------------
func BIF_finalize_mean_eb(mn, msum, msum2 *mlrval.Mlrval) *mlrval.Mlrval {
	mvar := BIF_finalize_variance(mn, msum, msum2)
	if mvar.IsVoid() {
		return mvar
	}
	return BIF_sqrt(BIF_divide(mvar, mn))
}

// ----------------------------------------------------------------
// Unbiased estimator:
//    (1/n)   sum{(xi-mean)**3}
//  -----------------------------
// [(1/(n-1)) sum{(xi-mean)**2}]**1.5

// mean = sumx / n; n mean = sumx

// sum{(xi-mean)^3}
//   = sum{xi^3 - 3 mean xi^2 + 3 mean^2 xi - mean^3}
//   = sum{xi^3} - 3 mean sum{xi^2} + 3 mean^2 sum{xi} - n mean^3
//   = sumx3 - 3 mean sumx2 + 3 mean^2 sumx - n mean^3
//   = sumx3 - 3 mean sumx2 + 3n mean^3 - n mean^3
//   = sumx3 - 3 mean sumx2 + 2n mean^3
//   = sumx3 - mean*(3 sumx2 + 2n mean^2)

// sum{(xi-mean)^2}
//   = sum{xi^2 - 2 mean xi + mean^2}
//   = sum{xi^2} - 2 mean sum{xi} + n mean^2
//   = sumx2 - 2 mean sumx + n mean^2
//   = sumx2 - 2 n mean^2 + n mean^2
//   = sumx2 - n mean^2

// ----------------------------------------------------------------
func BIF_finalize_skewness(mn, msum, msum2, msum3 *mlrval.Mlrval) *mlrval.Mlrval {
	n, isInt := mn.GetIntValue()
	lib.InternalCodingErrorIf(!isInt)
	if n < 2 {
		return mlrval.VOID
	}
	fn := float64(n)
	sum, isNumber := msum.GetNumericToFloatValue()
	lib.InternalCodingErrorIf(!isNumber)
	sum2, isNumber := msum2.GetNumericToFloatValue()
	lib.InternalCodingErrorIf(!isNumber)
	sum3, isNumber := msum3.GetNumericToFloatValue()
	lib.InternalCodingErrorIf(!isNumber)

	mean := sum / fn
	numerator := sum3 - mean*(3.0*sum2-2.0*fn*mean*mean)
	numerator = numerator / fn
	denominator := (sum2 - fn*mean*mean) / (fn - 1.0)
	denominator = math.Pow(denominator, 1.5)
	return mlrval.FromFloat(numerator / denominator)
}

// Unbiased:
//  (1/n) sum{(x-mean)**4}
//  ----------------------- - 3
// [(1/n) sum{(x-mean)**2}]**2

// sum{(xi-mean)^4}
//   = sum{xi^4 - 4 mean xi^3 + 6 mean^2 xi^2 - 4 mean^3 xi + mean^4}
//   = sum{xi^4} - 4 mean sum{xi^3} + 6 mean^2 sum{xi^2} - 4 mean^3 sum{xi} + n mean^4
//   = sum{xi^4} - 4 mean sum{xi^3} + 6 mean^2 sum{xi^2} - 4 n mean^4 + n mean^4
//   = sum{xi^4} - 4 mean sum{xi^3} + 6 mean^2 sum{xi^2} - 3 n mean^4
//   = sum{xi^4} - mean*(4 sum{xi^3} - 6 mean sum{xi^2} + 3 n mean^3)
//   = sumx4 - mean*(4 sumx3 - 6 mean sumx2 + 3 n mean^3)
//   = sumx4 - mean*(4 sumx3 - mean*(6 sumx2 - 3 n mean^2))

// ----------------------------------------------------------------
func BIF_finalize_kurtosis(mn, msum, msum2, msum3, msum4 *mlrval.Mlrval) *mlrval.Mlrval {
	n, isInt := mn.GetIntValue()
	lib.InternalCodingErrorIf(!isInt)
	if n < 2 {
		return mlrval.VOID
	}
	fn := float64(n)
	sum, isNumber := msum.GetNumericToFloatValue()
	lib.InternalCodingErrorIf(!isNumber)
	sum2, isNumber := msum2.GetNumericToFloatValue()
	lib.InternalCodingErrorIf(!isNumber)
	sum3, isNumber := msum3.GetNumericToFloatValue()
	lib.InternalCodingErrorIf(!isNumber)
	sum4, isNumber := msum4.GetNumericToFloatValue()
	lib.InternalCodingErrorIf(!isNumber)

	mean := sum / fn

	numerator := sum4 - mean*(4.0*sum3-mean*(6.0*sum2-3.0*fn*mean*mean))
	numerator = numerator / fn
	denominator := (sum2 - fn*mean*mean) / fn
	denominator = denominator * denominator
	return mlrval.FromFloat(numerator/denominator - 3.0)

}

// ================================================================
// XXX TEMP

// XXX COMMENT
func check_collection(c *mlrval.Mlrval) (bool, *mlrval.Mlrval) {
	vtype := c.Type()
	switch vtype {
	case mlrval.MT_ARRAY:
		return true, c
	case mlrval.MT_MAP:
		return true, c
	case mlrval.MT_ABSENT:
		return false, mlrval.ABSENT
	default:
		return false, mlrval.ERROR
	}
}

// collection_sum_of_function sums f(value) for value in the array or map:
// e.g. sum of values, sum of squares of values, etc.
func collection_sum_of_function(
	collection *mlrval.Mlrval,
	f func(element *mlrval.Mlrval) *mlrval.Mlrval,
) *mlrval.Mlrval {
	return mlrval.CollectionFold(
		collection,
		mlrval.FromInt(0),
		func(a, b *mlrval.Mlrval) *mlrval.Mlrval {
			return BIF_plus_binary(a, f(b))
		},
	)
}

func BIF_stats_count(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	if collection.IsArray() {
		arrayval := collection.AcquireArrayValue()
		return mlrval.FromInt(int64(len(arrayval)))
	} else {
		mapval := collection.AcquireMapValue()
		return mlrval.FromInt(mapval.FieldCount)
	}
}

func BIF_stats_null_count(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	f := func(element *mlrval.Mlrval) *mlrval.Mlrval {
		if element.IsVoid() || element.IsNull() {
			return mlrval.FromInt(1)
		} else {
			return mlrval.FromInt(0)
		}
	}
	return mlrval.CollectionFold(
		collection,
		mlrval.FromInt(0),
		func(a, b *mlrval.Mlrval) *mlrval.Mlrval {
			return BIF_plus_binary(a, f(b))
		},
	)
}

func BIF_stats_distinct_count(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	counts := make(map[string]int)
	if collection.IsArray() {
		a := collection.AcquireArrayValue()
		for _, e := range a {
			valueString := e.OriginalString()
			counts[valueString] += 1
		}
	} else {
		m := collection.AcquireMapValue()
		for pe := m.Head; pe != nil; pe = pe.Next {
			valueString := pe.Value.OriginalString()
			counts[valueString] += 1
		}
	}
	return mlrval.FromInt(int64(len(counts)))
}

func BIF_stats_mode(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	counts := make(map[string]int)
	if collection.IsArray() {
		a := collection.AcquireArrayValue()
		for _, e := range a {
			valueString := e.OriginalString()
			counts[valueString] += 1
		}
	} else {
		m := collection.AcquireMapValue()
		for pe := m.Head; pe != nil; pe = pe.Next {
			valueString := pe.Value.OriginalString()
			counts[valueString] += 1
		}
	}
	maxk := ""
	maxv := -1
	for k, v := range counts {
		if v > maxv {
			maxk = k
			maxv = v
		}
	}
	return mlrval.FromString(maxk)
}

func BIF_stats_antimode(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	counts := make(map[string]int)
	if collection.IsArray() {
		a := collection.AcquireArrayValue()
		for _, e := range a {
			valueString := e.OriginalString()
			counts[valueString] += 1
		}
	} else {
		m := collection.AcquireMapValue()
		for pe := m.Head; pe != nil; pe = pe.Next {
			valueString := pe.Value.OriginalString()
			counts[valueString] += 1
		}
	}
	first := true
	maxk := ""
	maxv := -1
	for k, v := range counts {
		if first || v < maxv {
			maxk = k
			maxv = v
			first = false
		}
	}
	return mlrval.FromString(maxk)
}

func BIF_stats_sum(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	return collection_sum_of_function(
		collection,
		func(e *mlrval.Mlrval) *mlrval.Mlrval {
			return e
		},
	)
}

func BIF_stats_sum2(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	f := func(element *mlrval.Mlrval) *mlrval.Mlrval {
		return BIF_times(element, element)
	}
	return collection_sum_of_function(collection, f)
}

func BIF_stats_sum3(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	f := func(element *mlrval.Mlrval) *mlrval.Mlrval {
		return BIF_times(element, BIF_times(element, element))
	}
	return collection_sum_of_function(collection, f)
}

func BIF_stats_sum4(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	f := func(element *mlrval.Mlrval) *mlrval.Mlrval {
		sq := BIF_times(element, element)
		return BIF_times(sq, sq)
	}
	return collection_sum_of_function(collection, f)
}

func BIF_stats_mean(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	n := BIF_stats_count(collection)
	sum := BIF_stats_sum(collection)
	return BIF_divide(sum, n)
}

func BIF_stats_meaneb(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	n := BIF_stats_count(collection)
	sum := BIF_stats_sum(collection)
	sum2 := BIF_stats_sum2(collection)
	return BIF_finalize_mean_eb(n, sum, sum2)
}

func BIF_stats_variance(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	n := BIF_stats_count(collection)
	sum := BIF_stats_sum(collection)
	sum2 := BIF_stats_sum2(collection)
	return BIF_finalize_variance(n, sum, sum2)
}

func BIF_stats_stddev(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	n := BIF_stats_count(collection)
	sum := BIF_stats_sum(collection)
	sum2 := BIF_stats_sum2(collection)
	return BIF_finalize_stddev(n, sum, sum2)
}

func BIF_stats_skewness(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	n := BIF_stats_count(collection)
	sum := BIF_stats_sum(collection)
	sum2 := BIF_stats_sum2(collection)
	sum3 := BIF_stats_sum3(collection)
	return BIF_finalize_skewness(n, sum, sum2, sum3)
}

func BIF_stats_kurtosis(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	n := BIF_stats_count(collection)
	sum := BIF_stats_sum(collection)
	sum2 := BIF_stats_sum2(collection)
	sum3 := BIF_stats_sum3(collection)
	sum4 := BIF_stats_sum4(collection)
	return BIF_finalize_kurtosis(n, sum, sum2, sum3, sum4)
}

func BIF_stats_min(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	if collection.IsArray() {
		return BIF_min_variadic(collection.AcquireArrayValue())
	} else {
		return BIF_min_over_map_values(collection.AcquireMapValue())
	}
}

func BIF_stats_max(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	if collection.IsArray() {
		return BIF_max_variadic(collection.AcquireArrayValue())
	} else {
		return BIF_max_over_map_values(collection.AcquireMapValue())
	}
}

func BIF_stats_minlen(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	if collection.IsArray() {
		return BIF_minlen_variadic(collection.AcquireArrayValue())
	} else {
		return BIF_minlen_over_map_values(collection.AcquireMapValue())
	}
}

func BIF_stats_maxlen(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	if collection.IsArray() {
		return BIF_maxlen_variadic(collection.AcquireArrayValue())
	} else {
		return BIF_maxlen_over_map_values(collection.AcquireMapValue())
	}
}

func BIF_sort_collection(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}

	var array []*mlrval.Mlrval
	if collection.IsArray() {
		arrayval := collection.AcquireArrayValue()
		n := len(arrayval)
		array = make([]*mlrval.Mlrval, n)
		for i := 0; i < n; i++ {
			array[i] = arrayval[i].Copy()
		}
	} else {
		mapval := collection.AcquireMapValue()
		n := mapval.FieldCount
		array = make([]*mlrval.Mlrval, n)
		i := 0
		for pe := mapval.Head; pe != nil; pe = pe.Next {
			array[i] = pe.Value.Copy()
			i++
		}
	}

	sort.Slice(array, func(i, j int) bool {
		return mlrval.LessThan(array[i], array[j])
	})

	return mlrval.FromArray(array)
}

func BIF_stats_median(
	collection *mlrval.Mlrval,
) *mlrval.Mlrval {
	return BIF_stats_percentile(collection, mlrval.FromFloat(50.0))
}

func BIF_stats_median_with_options(
	collection *mlrval.Mlrval,
	options *mlrval.Mlrval,
) *mlrval.Mlrval {
	return BIF_stats_percentile_with_options(collection, mlrval.FromFloat(50.0), options)
}

func BIF_stats_percentile(
	collection *mlrval.Mlrval,
	percentile *mlrval.Mlrval,
) *mlrval.Mlrval {
	return BIF_stats_percentile_with_options(collection, percentile, nil)
}

func BIF_stats_percentile_with_options(
	collection *mlrval.Mlrval,
	percentile *mlrval.Mlrval,
	options *mlrval.Mlrval,
) *mlrval.Mlrval {
	percentiles := mlrval.FromSingletonArray(percentile)
	outputs := BIF_stats_percentiles_with_options(collection, percentiles, options)
	return outputs.AcquireMapValue().Head.Value
}

func BIF_stats_percentiles(
	collection *mlrval.Mlrval,
	percentiles *mlrval.Mlrval,
) *mlrval.Mlrval {
	return BIF_stats_percentiles_with_options(collection, percentiles, nil)
}

func BIF_stats_percentiles_with_options(
	collection *mlrval.Mlrval,
	percentiles *mlrval.Mlrval,
	options *mlrval.Mlrval,
) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}

	array_is_sorted := false
	interpolate_linearly := false
	output_array_not_map := false

	if options != nil {
		om := options.GetMap()
		if om == nil { // not a map
			return mlrval.ERROR
		}
		for pe := om.Head; pe != nil; pe = pe.Next {
			if pe.Key == "array_is_sorted" {
				if mlrval.Equals(pe.Value, mlrval.TRUE) {
					array_is_sorted = true
				} else if mlrval.Equals(pe.Value, mlrval.FALSE) {
					array_is_sorted = false
				} else {
					return mlrval.ERROR
				}
			} else if pe.Key == "interpolate_linearly" {
				if mlrval.Equals(pe.Value, mlrval.TRUE) {
					interpolate_linearly = true
				} else if mlrval.Equals(pe.Value, mlrval.FALSE) {
					interpolate_linearly = false
				} else {
					return mlrval.ERROR
				}
			} else if pe.Key == "output_array_not_map" {
				if mlrval.Equals(pe.Value, mlrval.TRUE) {
					output_array_not_map = true
				} else if mlrval.Equals(pe.Value, mlrval.FALSE) {
					output_array_not_map = false
				} else {
					return mlrval.ERROR
				}
			}
		}
	}

	var sorted_array *mlrval.Mlrval
	if array_is_sorted {
		if !collection.IsArray() {
			return mlrval.ERROR
		}
		sorted_array = collection
	} else {
		sorted_array = BIF_sort_collection(collection)
	}

	return bif_percentiles(
		sorted_array.AcquireArrayValue(),
		percentiles,
		interpolate_linearly,
		output_array_not_map,
	)
}

func bif_percentiles(
	sorted_array []*mlrval.Mlrval,
	percentiles *mlrval.Mlrval,
	interpolate_linearly bool,
	output_array_not_map bool,
) *mlrval.Mlrval {

	ps := percentiles.GetArray()
	if ps == nil { // not an array
		return mlrval.ERROR
	}

	outputs := make([]*mlrval.Mlrval, len(ps))

	for i, _ := range ps {
		p, ok := ps[i].GetNumericToFloatValue()
		if !ok {
			outputs[i] = mlrval.ERROR.Copy()
		} else {
			if interpolate_linearly {
				outputs[i] = GetPercentileLinearlyInterpolated(sorted_array, len(sorted_array), p)
			} else {
				outputs[i] = GetPercentileNonInterpolated(sorted_array, len(sorted_array), p)
			}
		}
	}

	if output_array_not_map {
		return mlrval.FromArray(outputs)
	} else {
		m := mlrval.NewMlrmap()
		for i, _ := range ps {
			sp := ps[i].String()
			m.PutCopy(sp, outputs[i])
		}
		return mlrval.FromMap(m)
	}
}
