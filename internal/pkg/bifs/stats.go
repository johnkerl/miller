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
// STATS ROUTINES -- other than min/max which are placed separately.

// This is a helper function for BIFs which operate only on array or map.
// It shorthands what values to return for non-collection inputs.
func check_collection(c *mlrval.Mlrval, funcname string) (bool, *mlrval.Mlrval) {
	vtype := c.Type()
	switch vtype {
	case mlrval.MT_ARRAY:
		return true, c
	case mlrval.MT_MAP:
		return true, c
	case mlrval.MT_ABSENT:
		return false, mlrval.ABSENT
	case mlrval.MT_ERROR:
		return false, c
	default:
		return false, mlrval.FromNotCollectionError(funcname, c)
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

func BIF_count(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "count")
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

func BIF_null_count(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "null_count")
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

func BIF_distinct_count(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "distinct_count")
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

func BIF_mode(collection *mlrval.Mlrval) *mlrval.Mlrval {
	return bif_mode_or_antimode(collection, "mode", func(a, b int) bool { return a > b })
}

func BIF_antimode(collection *mlrval.Mlrval) *mlrval.Mlrval {
	return bif_mode_or_antimode(collection, "antimode", func(a, b int) bool { return a < b })
}

func bif_mode_or_antimode(
	collection *mlrval.Mlrval,
	funcname string,
	cmp func(int, int) bool,
) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, funcname)
	if !ok {
		return value_if_not
	}

	// Do not use a Go map[string]int as that makes the output in the case of ties
	// (e.g. input = [3,3,4,4]) non-determinstic. That's bad for unit tests and also
	// simply bad UX.
	counts := lib.NewOrderedMap()

	// We use stringification to detect uniqueness. Yet we want the output to be typed,
	// e.g. mode of an array of ints should be an int, not a string. Here we store
	// a reference to one representative for each equivalence class.
	reps := lib.NewOrderedMap()

	if collection.IsArray() {
		a := collection.AcquireArrayValue()
		if len(a) == 0 {
			return mlrval.VOID
		}
		for _, e := range a {
			valueString := e.OriginalString()
			if counts.Has(valueString) {
				counts.Put(valueString, counts.Get(valueString).(int)+1)
			} else {
				counts.Put(valueString, 1)
				reps.Put(valueString, e)
			}
		}
	} else {
		m := collection.AcquireMapValue()
		if m.Head == nil {
			return mlrval.VOID
		}
		for pe := m.Head; pe != nil; pe = pe.Next {
			valueString := pe.Value.OriginalString()
			if counts.Has(valueString) {
				counts.Put(valueString, counts.Get(valueString).(int)+1)
			} else {
				counts.Put(valueString, 1)
				reps.Put(valueString, pe.Value)
			}
		}
	}
	first := true
	maxk := ""
	maxv := -1
	for pf := counts.Head; pf != nil; pf = pf.Next {
		k := pf.Key
		v := pf.Value.(int)
		if first || cmp(v, maxv) {
			maxk = k
			maxv = v
			first = false
		}
	}
	// OrderedMap has interface{} values, so dereference as Mlrval. Then, copy the Mlrval
	// so we're not returning a pointer to input data.
	return reps.Get(maxk).(*mlrval.Mlrval).Copy()
}

func BIF_sum(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "sum")
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

func BIF_sum2(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "sum2")
	if !ok {
		return value_if_not
	}
	f := func(element *mlrval.Mlrval) *mlrval.Mlrval {
		return BIF_times(element, element)
	}
	return collection_sum_of_function(collection, f)
}

func BIF_sum3(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "sum3")
	if !ok {
		return value_if_not
	}
	f := func(element *mlrval.Mlrval) *mlrval.Mlrval {
		return BIF_times(element, BIF_times(element, element))
	}
	return collection_sum_of_function(collection, f)
}

func BIF_sum4(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "sum4")
	if !ok {
		return value_if_not
	}
	f := func(element *mlrval.Mlrval) *mlrval.Mlrval {
		sq := BIF_times(element, element)
		return BIF_times(sq, sq)
	}
	return collection_sum_of_function(collection, f)
}

func BIF_mean(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "mean")
	if !ok {
		return value_if_not
	}
	n := BIF_count(collection)
	if n.AcquireIntValue() == 0 {
		return mlrval.VOID
	}
	sum := BIF_sum(collection)
	return BIF_divide(sum, n)
}

func BIF_meaneb(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "meaneb")
	if !ok {
		return value_if_not
	}
	n := BIF_count(collection)
	sum := BIF_sum(collection)
	sum2 := BIF_sum2(collection)
	return BIF_finalize_mean_eb(n, sum, sum2)
}

func BIF_variance(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "variance")
	if !ok {
		return value_if_not
	}
	n := BIF_count(collection)
	sum := BIF_sum(collection)
	sum2 := BIF_sum2(collection)
	return BIF_finalize_variance(n, sum, sum2)
}

func BIF_stddev(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "stddev")
	if !ok {
		return value_if_not
	}
	n := BIF_count(collection)
	sum := BIF_sum(collection)
	sum2 := BIF_sum2(collection)
	return BIF_finalize_stddev(n, sum, sum2)
}

func BIF_skewness(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "skewness")
	if !ok {
		return value_if_not
	}
	n := BIF_count(collection)
	sum := BIF_sum(collection)
	sum2 := BIF_sum2(collection)
	sum3 := BIF_sum3(collection)
	return BIF_finalize_skewness(n, sum, sum2, sum3)
}

func BIF_kurtosis(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "kurtosis")
	if !ok {
		return value_if_not
	}
	n := BIF_count(collection)
	sum := BIF_sum(collection)
	sum2 := BIF_sum2(collection)
	sum3 := BIF_sum3(collection)
	sum4 := BIF_sum4(collection)
	return BIF_finalize_kurtosis(n, sum, sum2, sum3, sum4)
}

func BIF_minlen(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "minlen")
	if !ok {
		return value_if_not
	}
	if collection.IsArray() {
		return BIF_minlen_variadic(collection.AcquireArrayValue())
	} else {
		return BIF_minlen_within_map_values(collection.AcquireMapValue())
	}
}

func BIF_maxlen(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "maxlen")
	if !ok {
		return value_if_not
	}
	if collection.IsArray() {
		return BIF_maxlen_variadic(collection.AcquireArrayValue())
	} else {
		return BIF_maxlen_within_map_values(collection.AcquireMapValue())
	}
}

func BIF_sort_collection(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, "sort_collection")
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

func BIF_median(
	collection *mlrval.Mlrval,
) *mlrval.Mlrval {
	return bif_percentile_with_options_aux(collection, mlrval.FromFloat(50.0), nil, "median")
}

func BIF_median_with_options(
	collection *mlrval.Mlrval,
	options *mlrval.Mlrval,
) *mlrval.Mlrval {
	return bif_percentile_with_options_aux(collection, mlrval.FromFloat(50.0), options, "median")
}

func BIF_percentile(
	collection *mlrval.Mlrval,
	percentile *mlrval.Mlrval,
) *mlrval.Mlrval {
	return bif_percentile_with_options_aux(collection, percentile, nil, "percentile")
}

func BIF_percentile_with_options(
	collection *mlrval.Mlrval,
	percentile *mlrval.Mlrval,
	options *mlrval.Mlrval,
) *mlrval.Mlrval {
	return bif_percentile_with_options_aux(collection, percentile, options, "percentile")
}

func BIF_percentiles(
	collection *mlrval.Mlrval,
	percentiles *mlrval.Mlrval,
) *mlrval.Mlrval {
	return bif_percentiles_with_options_aux(collection, percentiles, nil, "percentiles")
}

func BIF_percentiles_with_options(
	collection *mlrval.Mlrval,
	percentiles *mlrval.Mlrval,
	options *mlrval.Mlrval,
) *mlrval.Mlrval {
	return bif_percentiles_with_options_aux(collection, percentiles, options, "percentiles")
}

func bif_percentile_with_options_aux(
	collection *mlrval.Mlrval,
	percentile *mlrval.Mlrval,
	options *mlrval.Mlrval,
	funcname string,
) *mlrval.Mlrval {
	percentiles := mlrval.FromSingletonArray(percentile)
	outputs := bif_percentiles_with_options_aux(collection, percentiles, options, funcname)

	// Check for error/absent returns from the main impl body
	ok, value_if_not := check_collection(outputs, funcname)
	if !ok {
		return value_if_not
	}

	return outputs.AcquireMapValue().Head.Value
}

func bif_percentiles_with_options_aux(
	collection *mlrval.Mlrval,
	percentiles *mlrval.Mlrval,
	options *mlrval.Mlrval,
	funcname string,
) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection, funcname)
	if !ok {
		return value_if_not
	}

	array_is_sorted := false
	interpolate_linearly := false
	output_array_not_map := false

	if options != nil {
		om := options.GetMap()
		if om == nil { // not a map
			return type_error_named_argument(funcname, "map", "options", options)
		}
		for pe := om.Head; pe != nil; pe = pe.Next {
			if pe.Key == "array_is_sorted" || pe.Key == "ais" {
				if mlrval.Equals(pe.Value, mlrval.TRUE) {
					array_is_sorted = true
				} else if mlrval.Equals(pe.Value, mlrval.FALSE) {
					array_is_sorted = false
				} else {
					return type_error_named_argument(funcname, "boolean", pe.Key, pe.Value)
				}
			} else if pe.Key == "interpolate_linearly" || pe.Key == "il" {
				if mlrval.Equals(pe.Value, mlrval.TRUE) {
					interpolate_linearly = true
				} else if mlrval.Equals(pe.Value, mlrval.FALSE) {
					interpolate_linearly = false
				} else {
					return type_error_named_argument(funcname, "boolean", pe.Key, pe.Value)
				}
			} else if pe.Key == "output_array_not_map" || pe.Key == "oa" {
				if mlrval.Equals(pe.Value, mlrval.TRUE) {
					output_array_not_map = true
				} else if mlrval.Equals(pe.Value, mlrval.FALSE) {
					output_array_not_map = false
				} else {
					return type_error_named_argument(funcname, "boolean", pe.Key, pe.Value)
				}
			}
		}
	}

	var sorted_array *mlrval.Mlrval
	if array_is_sorted {
		if !collection.IsArray() {
			return mlrval.FromNotArrayError(funcname+" collection", collection)
		}
		sorted_array = collection
	} else {
		sorted_array = BIF_sort_collection(collection)
	}

	return bif_percentiles_impl(
		sorted_array.AcquireArrayValue(),
		percentiles,
		interpolate_linearly,
		output_array_not_map,
		funcname,
	)
}

func bif_percentiles_impl(
	sorted_array []*mlrval.Mlrval,
	percentiles *mlrval.Mlrval,
	interpolate_linearly bool,
	output_array_not_map bool,
	funcname string,
) *mlrval.Mlrval {

	ps := percentiles.GetArray()
	if ps == nil { // not an array
		return mlrval.FromNotArrayError(funcname+" percentiles", percentiles)
	}

	outputs := make([]*mlrval.Mlrval, len(ps))

	for i, _ := range ps {
		p, ok := ps[i].GetNumericToFloatValue()
		if !ok {
			outputs[i] = type_error_named_argument(funcname, "numeric", "percentile", ps[i])
		} else if len(sorted_array) == 0 {
			outputs[i] = mlrval.VOID
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
