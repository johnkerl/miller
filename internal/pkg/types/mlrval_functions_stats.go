package types

import (
	"math"

	"mlr/internal/pkg/lib"
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
func MlrvalGetVar(mn, msum, msum2 *Mlrval) *Mlrval {
	n, isInt := mn.GetIntValue()
	lib.InternalCodingErrorIf(!isInt)
	sum, isNumber := msum.GetNumericToFloatValue()
	lib.InternalCodingErrorIf(!isNumber)
	sum2, isNumber := msum2.GetNumericToFloatValue()
	lib.InternalCodingErrorIf(!isNumber)

	if n < 2 {
		return MLRVAL_VOID
	}

	mean := float64(sum) / float64(n)
	numerator := sum2 - mean*(2.0*sum-float64(n)*mean)
	if numerator < 0.0 { // round-off error
		numerator = 0.0
	}
	denominator := float64(n - 1)
	return MlrvalFromFloat64(numerator / denominator)
}

// ----------------------------------------------------------------
func MlrvalGetStddev(mn, msum, msum2 *Mlrval) *Mlrval {
	mvar := MlrvalGetVar(mn, msum, msum2)
	if mvar.IsVoid() {
		return mvar
	}
	return BIF_sqrt(mvar)
}

// ----------------------------------------------------------------
func MlrvalGetMeanEB(mn, msum, msum2 *Mlrval) *Mlrval {
	mvar := MlrvalGetVar(mn, msum, msum2)
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
func MlrvalGetSkewness(mn, msum, msum2, msum3 *Mlrval) *Mlrval {
	n, isInt := mn.GetIntValue()
	lib.InternalCodingErrorIf(!isInt)
	if n < 2 {
		return MLRVAL_VOID
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
	return MlrvalFromFloat64(numerator / denominator)
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
func MlrvalGetKurtosis(mn, msum, msum2, msum3, msum4 *Mlrval) *Mlrval {
	n, isInt := mn.GetIntValue()
	lib.InternalCodingErrorIf(!isInt)
	if n < 2 {
		return MLRVAL_VOID
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
	return MlrvalFromFloat64(numerator/denominator - 3.0)

}
