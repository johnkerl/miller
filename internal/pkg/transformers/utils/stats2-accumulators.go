// ================================================================
// For stats2
// ================================================================

package utils

import (
	"fmt"
	"math"
	"os"

	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
type IStats2Accumulator interface {
	Ingest(
		x float64,
		y float64,
	)

	Populate(
		valueFieldName1 string,
		valueFieldName2 string,
		outrec *types.Mlrmap,
	)

	Fit(
		x float64,
		y float64,
		outrec *types.Mlrmap,
	)
}

type newStats2AccumulatorFunc func(
	valueFieldName1 string,
	valueFieldName2 string,
	accumulatorName string,
	doVerbose bool,
) IStats2Accumulator

type stats2AccumulatorInfo struct {
	name        string
	description string
	constructor newStats2AccumulatorFunc
}

var stats2AccumulatorInfos []stats2AccumulatorInfo = []stats2AccumulatorInfo{
	{
		"linreg-ols",
		"Linear regression using ordinary least squares",
		NewStats2LinRegOLSAccumulator,
	},
	{
		"linreg-pca",
		"Linear regression using principal component analysis",
		NewStats2LinRegPCAAccumulator,
	},
	{
		"r2",
		"Quality metric for linreg-ols (linreg-pca emits its own)",
		NewStats2R2Accumulator,
	},
	{
		"logireg",
		"Logistic regression",
		NewStats2LogiRegAccumulator,
	},
	{
		"corr",
		"Sample correlation",
		NewStats2CorrAccumulator,
	},
	{
		"cov",
		"Sample covariance",
		NewStats2CovAccumulator,
	},
	{
		"covx",
		"Sample-covariance matrix",
		NewStats2CovXAccumulator,
	},
}

// ----------------------------------------------------------------
type Stats2AccumulatorFactory struct {
}

func NewStats2AccumulatorFactory() *Stats2AccumulatorFactory {
	return &Stats2AccumulatorFactory{}
}

func ListStats2Accumulators(o *os.File) {
	for _, info := range stats2AccumulatorInfos {
		fmt.Fprintf(o, "  %-8s %s\n", info.name, info.description)
	}
}

func ValidateStats2AccumulatorName(
	accumulatorName string,
) bool {
	for _, info := range stats2AccumulatorInfos {
		if info.name == accumulatorName {
			return true
		}
	}
	return false
}

func (factory *Stats2AccumulatorFactory) Make(
	valueFieldName1 string,
	valueFieldName2 string,
	accumulatorName string,
	doVerbose bool,
) IStats2Accumulator {
	// TODO: hashmapify the lookup table
	for _, info := range stats2AccumulatorInfos {
		if info.name == accumulatorName {
			return info.constructor(valueFieldName1, valueFieldName2, accumulatorName, doVerbose)
		}
	}
	return nil
}

// ================================================================
type Stats2LinRegOLSAccumulator struct {
	count              int
	sumx               float64
	sumy               float64
	sumx2              float64
	sumxy              float64
	mOutputFieldName   string
	bOutputFieldName   string
	nOutputFieldName   string
	fitOutputFieldName string
	fitReady           bool
	m                  float64
	b                  float64
}

func NewStats2LinRegOLSAccumulator(
	valueFieldName1 string,
	valueFieldName2 string,
	accumulatorName string,
	doVerbose bool,
) IStats2Accumulator {
	prefix := valueFieldName1 + "_" + valueFieldName2 + "_"
	return &Stats2LinRegOLSAccumulator{
		count:              0,
		sumx:               0.0,
		sumy:               0.0,
		sumx2:              0.0,
		sumxy:              0.0,
		mOutputFieldName:   prefix + "ols_m",
		bOutputFieldName:   prefix + "ols_b",
		nOutputFieldName:   prefix + "ols_n",
		fitOutputFieldName: prefix + "ols_fit",
		fitReady:           false,
		m:                  -999.0,
		b:                  -999.0,
	}
}

func (acc *Stats2LinRegOLSAccumulator) Ingest(
	x float64,
	y float64,
) {
	acc.count++
	acc.sumx += x
	acc.sumy += y
	acc.sumx2 += x * x
	acc.sumxy += x * y
}

func (acc *Stats2LinRegOLSAccumulator) Populate(
	valueFieldName1 string,
	valueFieldName2 string,
	outrec *types.Mlrmap,
) {
	if acc.count < 2 {
		outrec.PutCopy(acc.mOutputFieldName, types.MLRVAL_VOID)
		outrec.PutCopy(acc.bOutputFieldName, types.MLRVAL_VOID)
	} else {

		m, b := lib.GetLinearRegressionOLS(acc.count, acc.sumx, acc.sumx2, acc.sumxy, acc.sumy)

		outrec.PutReference(acc.mOutputFieldName, types.MlrvalFromFloat64(m))
		outrec.PutReference(acc.bOutputFieldName, types.MlrvalFromFloat64(b))
	}
	outrec.PutReference(acc.nOutputFieldName, types.MlrvalFromInt(acc.count))
}

func (acc *Stats2LinRegOLSAccumulator) Fit(
	x float64,
	y float64,
	outrec *types.Mlrmap,
) {

	if !acc.fitReady {
		// Idea for hold-and-fit in stats2.go is:
		// * We've ingested say 10,000 records
		// * After the end of those we compute m and b
		// * Then for all 10,000 records we compute y = m*x + b
		// The fitReady flag keeps us from recomputing the linear fit 10,000 times
		acc.m, acc.b = lib.GetLinearRegressionOLS(acc.count, acc.sumx, acc.sumx2, acc.sumxy, acc.sumy)
		acc.fitReady = true
	}

	if acc.count < 2 {
		outrec.PutCopy(acc.fitOutputFieldName, types.MLRVAL_VOID)
	} else {
		yfit := acc.m*x + acc.b
		outrec.PutReference(acc.fitOutputFieldName, types.MlrvalFromFloat64(yfit))
	}
}

// ================================================================
const LOGIREG_DVECTOR_INITIAL_SIZE = 16

type Stats2LogiRegAccumulator struct {
	xs                 []float64
	ys                 []float64
	mOutputFieldName   string
	bOutputFieldName   string
	nOutputFieldName   string
	fitOutputFieldName string
	fitReady           bool
	m                  float64
	b                  float64
}

func NewStats2LogiRegAccumulator(
	valueFieldName1 string,
	valueFieldName2 string,
	accumulatorName string,
	doVerbose bool,
) IStats2Accumulator {
	prefix := valueFieldName1 + "_" + valueFieldName2 + "_"
	return &Stats2LogiRegAccumulator{
		xs:                 make([]float64, 0, LOGIREG_DVECTOR_INITIAL_SIZE),
		ys:                 make([]float64, 0, LOGIREG_DVECTOR_INITIAL_SIZE),
		mOutputFieldName:   prefix + "logistic_m",
		bOutputFieldName:   prefix + "logistic_b",
		nOutputFieldName:   prefix + "logistic_n",
		fitOutputFieldName: prefix + "logistic_fit",
		fitReady:           false,
		m:                  -999.0,
		b:                  -999.0,
	}
}

func (acc *Stats2LogiRegAccumulator) Ingest(
	x float64,
	y float64,
) {
	acc.xs = append(acc.xs, x) // append is smart about cap-increase via doubling
	acc.ys = append(acc.ys, y) // append is smart about cap-increase via doubling
}

func (acc *Stats2LogiRegAccumulator) Populate(
	valueFieldName1 string,
	valueFieldName2 string,
	outrec *types.Mlrmap,
) {

	if len(acc.xs) < 2 {
		outrec.PutCopy(acc.mOutputFieldName, types.MLRVAL_VOID)
		outrec.PutCopy(acc.bOutputFieldName, types.MLRVAL_VOID)
	} else {
		m, b := lib.LogisticRegression(acc.xs, acc.ys)
		outrec.PutCopy(acc.mOutputFieldName, types.MlrvalFromFloat64(m))
		outrec.PutCopy(acc.bOutputFieldName, types.MlrvalFromFloat64(b))
	}
	outrec.PutReference(acc.nOutputFieldName, types.MlrvalFromInt(len(acc.xs)))
}

func (acc *Stats2LogiRegAccumulator) Fit(
	x float64,
	y float64,
	outrec *types.Mlrmap,
) {

	if !acc.fitReady {
		// Idea for hold-and-fit in stats2.go is:
		// * We've ingested say 10,000 records
		// * After the end of those we compute m and b
		// * Then for all 10,000 records we compute y = m*x + b
		// The fitReady flag keeps us from recomputing the linear fit 10,000 times
		acc.m, acc.b = lib.LogisticRegression(acc.xs, acc.ys)
		acc.fitReady = true
	}

	if len(acc.xs) < 2 {
		outrec.PutCopy(acc.fitOutputFieldName, types.MLRVAL_VOID)
	} else {
		yfit := 1.0 / (1.0 + math.Exp(-acc.m*x-acc.b))
		outrec.PutReference(acc.fitOutputFieldName, types.MlrvalFromFloat64(yfit))
	}
}

// ================================================================
// http://en.wikipedia.org/wiki/Pearson_product-moment_correlation_coefficient
// Alternatively, just use sqrt(corr) as defined above.

type Stats2R2Accumulator struct {
	count             int
	sumx              float64
	sumy              float64
	sumx2             float64
	sumxy             float64
	sumy2             float64
	r2OutputFieldName string
}

func NewStats2R2Accumulator(
	valueFieldName1 string,
	valueFieldName2 string,
	accumulatorName string,
	doVerbose bool,
) IStats2Accumulator {
	prefix := valueFieldName1 + "_" + valueFieldName2 + "_"
	return &Stats2R2Accumulator{
		count:             0,
		sumx:              0.0,
		sumy:              0.0,
		sumx2:             0.0,
		sumxy:             0.0,
		sumy2:             0.0,
		r2OutputFieldName: prefix + "r2",
	}
}

func (acc *Stats2R2Accumulator) Ingest(
	x float64,
	y float64,
) {
	acc.count++
	acc.sumx += x
	acc.sumy += y
	acc.sumx2 += x * x
	acc.sumxy += x * y
	acc.sumy2 += y * y
}

func (acc *Stats2R2Accumulator) Populate(
	valueFieldName1 string,
	valueFieldName2 string,
	outrec *types.Mlrmap,
) {

	if acc.count < 2 {
		outrec.PutCopy(acc.r2OutputFieldName, types.MLRVAL_VOID)
	} else {
		n := float64(acc.count)
		sumx := acc.sumx
		sumy := acc.sumy
		sumx2 := acc.sumx2
		sumy2 := acc.sumy2
		sumxy := acc.sumxy
		numerator := n*sumxy - sumx*sumy
		numerator = numerator * numerator
		denominator := (n*sumx2 - sumx*sumx) * (n*sumy2 - sumy*sumy)
		output := numerator / denominator
		outrec.PutReference(acc.r2OutputFieldName, types.MlrvalFromFloat64(output))
	}
}

// Trivial function; there is no fit-feature here
func (acc *Stats2R2Accumulator) Fit(
	x float64,
	y float64,
	outrec *types.Mlrmap,
) {
}

// ================================================================
// Shared code for Corr, Cov, CovX, and LinRegPCA.
// Corr(X,Y) = Cov(X,Y) / sigma_X sigma_Y.

type BivarMeasure int

const (
	DO_CORR BivarMeasure = iota
	DO_COV
	DO_COVX
	DO_LINREG_PCA
)

type Stats2CorrCovAccumulator struct {
	count int
	sumx  float64
	sumy  float64
	sumx2 float64
	sumxy float64
	sumy2 float64

	doWhich   BivarMeasure
	doVerbose bool

	corrOutputFieldName string

	covOutputFieldName string

	covx00OutputFieldName string
	covx01OutputFieldName string
	covx10OutputFieldName string
	covx11OutputFieldName string

	pca_mOutputFieldName   string
	pca_bOutputFieldName   string
	pca_nOutputFieldName   string
	pca_qOutputFieldName   string
	pca_l1OutputFieldName  string
	pca_l2OutputFieldName  string
	pca_v11OutputFieldName string
	pca_v12OutputFieldName string
	pca_v21OutputFieldName string
	pca_v22OutputFieldName string
	pca_fitOutputFieldName string

	fitReady bool
	m        float64
	b        float64
	q        float64
}

// ----------------------------------------------------------------
func NewStats2CorrCovAccumulator(
	valueFieldName1 string,
	valueFieldName2 string,
	accumulatorName string,
	doVerbose bool,
	doWhich BivarMeasure,
) IStats2Accumulator {
	prefix := valueFieldName1 + "_" + valueFieldName2 + "_"
	return &Stats2CorrCovAccumulator{
		count: 0,

		sumx:      0.0,
		sumy:      0.0,
		sumx2:     0.0,
		sumxy:     0.0,
		sumy2:     0.0,
		doWhich:   doWhich,
		doVerbose: doVerbose,

		corrOutputFieldName: prefix + "corr",

		covOutputFieldName: prefix + "cov",

		covx00OutputFieldName: valueFieldName1 + "_" + valueFieldName1 + "_covx",
		covx01OutputFieldName: valueFieldName1 + "_" + valueFieldName2 + "_covx",
		covx10OutputFieldName: valueFieldName2 + "_" + valueFieldName1 + "_covx",
		covx11OutputFieldName: valueFieldName2 + "_" + valueFieldName2 + "_covx",

		pca_mOutputFieldName:   prefix + "pca_m",
		pca_bOutputFieldName:   prefix + "pca_b",
		pca_nOutputFieldName:   prefix + "pca_n",
		pca_qOutputFieldName:   prefix + "pca_quality",
		pca_l1OutputFieldName:  prefix + "pca_eival1",
		pca_l2OutputFieldName:  prefix + "pca_eival2",
		pca_v11OutputFieldName: prefix + "pca_eivec11",
		pca_v12OutputFieldName: prefix + "pca_eivec12",
		pca_v21OutputFieldName: prefix + "pca_eivec21",
		pca_v22OutputFieldName: prefix + "pca_eivec22",
		pca_fitOutputFieldName: prefix + "pca_fit",

		fitReady: false,
		m:        -999.0,
		b:        -999.0,
	}
}

func (acc *Stats2CorrCovAccumulator) Ingest(
	x float64,
	y float64,
) {
	acc.count++
	acc.sumx += x
	acc.sumy += y
	acc.sumx2 += x * x
	acc.sumxy += x * y
	acc.sumy2 += y * y
}

func (acc *Stats2CorrCovAccumulator) Populate(
	valueFieldName1 string,
	valueFieldName2 string,
	outrec *types.Mlrmap,
) {

	if acc.doWhich == DO_COVX {
		key00 := acc.covx00OutputFieldName
		key01 := acc.covx01OutputFieldName
		key10 := acc.covx10OutputFieldName
		key11 := acc.covx11OutputFieldName
		if acc.count < 2 {
			outrec.PutCopy(key00, types.MLRVAL_VOID)
			outrec.PutCopy(key01, types.MLRVAL_VOID)
			outrec.PutCopy(key10, types.MLRVAL_VOID)
			outrec.PutCopy(key11, types.MLRVAL_VOID)
		} else {
			Q := lib.GetCovMatrix(
				acc.count,
				acc.sumx,
				acc.sumx2,
				acc.sumy,
				acc.sumy2,
				acc.sumxy,
			)
			outrec.PutReference(key00, types.MlrvalFromFloat64(Q[0][0]))
			outrec.PutReference(key01, types.MlrvalFromFloat64(Q[0][1]))
			outrec.PutReference(key10, types.MlrvalFromFloat64(Q[1][0]))
			outrec.PutReference(key11, types.MlrvalFromFloat64(Q[1][1]))
		}

	} else if acc.doWhich == DO_LINREG_PCA {
		keym := acc.pca_mOutputFieldName
		keyb := acc.pca_bOutputFieldName
		keyn := acc.pca_nOutputFieldName
		keyq := acc.pca_qOutputFieldName

		keyl1 := acc.pca_l1OutputFieldName
		keyl2 := acc.pca_l2OutputFieldName
		keyv11 := acc.pca_v11OutputFieldName
		keyv12 := acc.pca_v12OutputFieldName
		keyv21 := acc.pca_v21OutputFieldName
		keyv22 := acc.pca_v22OutputFieldName

		if acc.count < 2 {
			outrec.PutCopy(keym, types.MLRVAL_VOID)
			outrec.PutCopy(keyb, types.MLRVAL_VOID)
			outrec.PutCopy(keyn, types.MLRVAL_VOID)
			outrec.PutCopy(keyq, types.MLRVAL_VOID)

			if acc.doVerbose {

				outrec.PutCopy(keyl1, types.MLRVAL_VOID)
				outrec.PutCopy(keyl2, types.MLRVAL_VOID)
				outrec.PutCopy(keyv11, types.MLRVAL_VOID)
				outrec.PutCopy(keyv12, types.MLRVAL_VOID)
				outrec.PutCopy(keyv21, types.MLRVAL_VOID)
				outrec.PutCopy(keyv22, types.MLRVAL_VOID)
			}
		} else {
			Q := lib.GetCovMatrix(
				acc.count,
				acc.sumx,
				acc.sumx2,
				acc.sumy,
				acc.sumy2,
				acc.sumxy,
			)

			l1, l2, v1, v2 := lib.GetRealSymmetricEigensystem(Q)

			x_mean := acc.sumx / float64(acc.count)
			y_mean := acc.sumy / float64(acc.count)
			m, b, q := lib.GetLinearRegressionPCA(l1, l2, v1, v2, x_mean, y_mean)

			outrec.PutReference(keym, types.MlrvalFromFloat64(m))
			outrec.PutReference(keyb, types.MlrvalFromFloat64(b))
			outrec.PutReference(keyn, types.MlrvalFromInt(acc.count))
			outrec.PutReference(keyq, types.MlrvalFromFloat64(q))

			if acc.doVerbose {
				outrec.PutReference(keyl1, types.MlrvalFromFloat64(l1))
				outrec.PutReference(keyl2, types.MlrvalFromFloat64(l2))
				outrec.PutReference(keyv11, types.MlrvalFromFloat64(v1[0]))
				outrec.PutReference(keyv12, types.MlrvalFromFloat64(v1[1]))
				outrec.PutReference(keyv21, types.MlrvalFromFloat64(v2[0]))
				outrec.PutReference(keyv22, types.MlrvalFromFloat64(v2[1]))
			}
		}
	} else {
		key := acc.corrOutputFieldName
		if acc.doWhich == DO_COV {
			key = acc.covOutputFieldName
		}
		if acc.count < 2 {
			outrec.PutCopy(key, types.MLRVAL_VOID)
		} else {
			output := lib.GetCov(acc.count, acc.sumx, acc.sumy, acc.sumxy)
			if acc.doWhich == DO_CORR {
				sigmax := math.Sqrt(lib.GetVar(acc.count, acc.sumx, acc.sumx2))
				sigmay := math.Sqrt(lib.GetVar(acc.count, acc.sumy, acc.sumy2))
				output = output / sigmax / sigmay
			}
			outrec.PutReference(key, types.MlrvalFromFloat64(output))
		}
	}
}

func (acc *Stats2CorrCovAccumulator) Fit(
	x float64,
	y float64,
	outrec *types.Mlrmap,
) {
	if acc.doWhich != DO_LINREG_PCA {
		return
	}

	if !acc.fitReady {
		// Idea for hold-and-fit in stats2.go is:
		// * We've ingested say 10,000 records
		// * After the end of those we compute m and b
		// * Then for all 10,000 records we compute y = m*x + b
		// The fitReady flag keeps us from recomputing the linear fit 10,000 times
		Q := lib.GetCovMatrix(acc.count, acc.sumx, acc.sumx2, acc.sumy, acc.sumy2, acc.sumxy)

		l1, l2, v1, v2 := lib.GetRealSymmetricEigensystem(Q)

		x_mean := acc.sumx / float64(acc.count)
		y_mean := acc.sumy / float64(acc.count)
		acc.m, acc.b, acc.q = lib.GetLinearRegressionPCA(l1, l2, v1, v2, x_mean, y_mean)

		acc.fitReady = true
	}
	if acc.count < 2 {
		outrec.PutCopy(acc.pca_fitOutputFieldName, types.MLRVAL_VOID)
	} else {
		yfit := acc.m*x + acc.b
		outrec.PutCopy(acc.pca_fitOutputFieldName, types.MlrvalFromFloat64(yfit))
	}
}

// ================================================================
func NewStats2CorrAccumulator(
	valueFieldName1 string,
	valueFieldName2 string,
	accumulatorName string,
	doVerbose bool,
) IStats2Accumulator {
	return NewStats2CorrCovAccumulator(
		valueFieldName1,
		valueFieldName2,
		accumulatorName,
		doVerbose,
		DO_CORR,
	)
}

func NewStats2CovAccumulator(
	valueFieldName1 string,
	valueFieldName2 string,
	accumulatorName string,
	doVerbose bool,
) IStats2Accumulator {
	return NewStats2CorrCovAccumulator(
		valueFieldName1,
		valueFieldName2,
		accumulatorName,
		doVerbose,
		DO_COV,
	)
}

func NewStats2CovXAccumulator(
	valueFieldName1 string,
	valueFieldName2 string,
	accumulatorName string,
	doVerbose bool,
) IStats2Accumulator {
	return NewStats2CorrCovAccumulator(
		valueFieldName1,
		valueFieldName2,
		accumulatorName,
		doVerbose,
		DO_COVX,
	)
}

func NewStats2LinRegPCAAccumulator(
	valueFieldName1 string,
	valueFieldName2 string,
	accumulatorName string,
	doVerbose bool,
) IStats2Accumulator {
	return NewStats2CorrCovAccumulator(
		valueFieldName1,
		valueFieldName2,
		accumulatorName,
		doVerbose,
		DO_LINREG_PCA,
	)
}
