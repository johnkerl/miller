package bifs

import (
	"math"

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

func BIF_collection_count(collection *mlrval.Mlrval) *mlrval.Mlrval {
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

func BIF_collection_null_count(collection *mlrval.Mlrval) *mlrval.Mlrval {
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

func BIF_collection_distinct_count(collection *mlrval.Mlrval) *mlrval.Mlrval {
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

func BIF_collection_mode(collection *mlrval.Mlrval) *mlrval.Mlrval {
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

func BIF_collection_antimode(collection *mlrval.Mlrval) *mlrval.Mlrval {
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

func BIF_collection_sum(collection *mlrval.Mlrval) *mlrval.Mlrval {
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

func BIF_collection_sum2(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	f := func(element *mlrval.Mlrval) *mlrval.Mlrval {
		return BIF_times(element, element)
	}
	return collection_sum_of_function(collection, f)
}

func BIF_collection_sum3(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	f := func(element *mlrval.Mlrval) *mlrval.Mlrval {
		return BIF_times(element, BIF_times(element, element))
	}
	return collection_sum_of_function(collection, f)
}

func BIF_collection_sum4(collection *mlrval.Mlrval) *mlrval.Mlrval {
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

func BIF_collection_mean(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	n := BIF_collection_count(collection)
	sum := BIF_collection_sum(collection)
	return BIF_divide(sum, n)
}

func BIF_collection_meaneb(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	n := BIF_collection_count(collection)
	sum := BIF_collection_sum(collection)
	sum2 := BIF_collection_sum2(collection)
	return BIF_finalize_mean_eb(n, sum, sum2)
}

func BIF_collection_variance(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	n := BIF_collection_count(collection)
	sum := BIF_collection_sum(collection)
	sum2 := BIF_collection_sum2(collection)
	return BIF_finalize_variance(n, sum, sum2)
}

func BIF_collection_stddev(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	n := BIF_collection_count(collection)
	sum := BIF_collection_sum(collection)
	sum2 := BIF_collection_sum2(collection)
	return BIF_finalize_stddev(n, sum, sum2)
}

func BIF_collection_skewness(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	n := BIF_collection_count(collection)
	sum := BIF_collection_sum(collection)
	sum2 := BIF_collection_sum2(collection)
	sum3 := BIF_collection_sum3(collection)
	return BIF_finalize_skewness(n, sum, sum2, sum3)
}

func BIF_collection_kurtosis(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, value_if_not := check_collection(collection)
	if !ok {
		return value_if_not
	}
	n := BIF_collection_count(collection)
	sum := BIF_collection_sum(collection)
	sum2 := BIF_collection_sum2(collection)
	sum3 := BIF_collection_sum3(collection)
	sum4 := BIF_collection_sum4(collection)
	return BIF_finalize_kurtosis(n, sum, sum2, sum3, sum4)
}

func BIF_collection_min(collection *mlrval.Mlrval) *mlrval.Mlrval {
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

func BIF_collection_max(collection *mlrval.Mlrval) *mlrval.Mlrval {
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

func BIF_collection_minlen(collection *mlrval.Mlrval) *mlrval.Mlrval {
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

func BIF_collection_maxlen(collection *mlrval.Mlrval) *mlrval.Mlrval {
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

//	"sort"
//func BIF_sort_in_place(collection *mlrval.Mlrval) {
//	// XXX what if map?
//	if xxx {
//		sort.Slice(array, func(i, j int) bool {
//			return mlrval.LessThan(array[i], array[j])
//		})
//		keeper.sorted = true
//	}
//}

// * count    Count instances of fields
// * sum      Compute sums of specified fields
// * mean     Compute averages (sample means) of specified fields
// * meaneb   Estimate error bars for averages (assuming no sample autocorrelation)
// * var      Compute sample variance of specified fields
// * stddev   Compute sample standard deviation of specified fields
// * skewness Compute sample skewness of specified fields
// * kurtosis Compute sample kurtosis of specified fields

// * distinct_count Count number of distinct values per field
// * null_count Count number of empty-string/JSON-null instances per field

// * min      Compute minimum values of specified fields
// * max      Compute maximum values of specified fields

// * minlen   Compute minimum string-lengths of specified fields
// * maxlen   Compute maximum string-lengths of specified fields

// * mode     Find most-frequently-occurring values for fields; first-found wins tie
// * antimode Find least-frequently-occurring values for fields; first-found wins tie

//   p10 p25.2 p50 p98 p100 etc.
//   median   This is the same as p50
