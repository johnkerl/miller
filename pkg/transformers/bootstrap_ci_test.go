package transformers

import (
	"testing"

	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

func TestBootstrapCIRejectsBadAccumulatorName(t *testing.T) {
	_, err := NewTransformerBootstrapCI(
		[]string{"nonesuch"},
		[]string{"x"},
		[]string{},
		100,
		0.95,
		false,
	)
	if err == nil {
		t.Fatal("expected error for unknown accumulator name, got nil")
	}
}

func TestBootstrapCIParseCLIValidation(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"missing -f", []string{"bootstrap-ci"}},
		{"non-positive -n", []string{"bootstrap-ci", "-f", "x", "-n", "0"}},
		{"confidence level too low", []string{"bootstrap-ci", "-f", "x", "-c", "0"}},
		{"confidence level too high", []string{"bootstrap-ci", "-f", "x", "-c", "1"}},
		{"unknown accumulator", []string{"bootstrap-ci", "-f", "x", "-a", "nonesuch"}},
		{"unknown option", []string{"bootstrap-ci", "-f", "x", "-q"}},
	}

	for _, tc := range cases {
		argi := 0
		_, err := transformerBootstrapCIParseCLI(&argi, len(tc.args), tc.args, nil, true)
		if err == nil {
			t.Errorf("%s: expected parse error, got nil", tc.name)
		}
	}
}

// feedBootstrapCI passes records with the given x-values (and a constant
// group-by field g=one) through the transformer, followed by an end-of-stream
// marker, returning the transformer's output records.
func feedBootstrapCI(tr *TransformerBootstrapCI, xValues []int64) []*types.RecordAndContext {
	context := types.NewContext()
	inputDownstreamDoneChannel := make(chan bool, 1)
	outputDownstreamDoneChannel := make(chan bool, 1)

	outputRecordsAndContexts := make([]*types.RecordAndContext, 0)
	for _, xValue := range xValues {
		record := mlrval.NewMlrmapAsRecord()
		record.PutCopy("g", mlrval.FromString("one"))
		record.PutCopy("x", mlrval.FromInt(xValue))
		_ = tr.Transform(
			types.NewRecordAndContext(record, context),
			&outputRecordsAndContexts,
			inputDownstreamDoneChannel,
			outputDownstreamDoneChannel,
		)
	}
	_ = tr.Transform(
		types.NewEndOfStreamMarker(context),
		&outputRecordsAndContexts,
		inputDownstreamDoneChannel,
		outputDownstreamDoneChannel,
	)
	return outputRecordsAndContexts
}

func TestBootstrapCIEndToEnd(t *testing.T) {
	lib.SeedRandom(12345)

	tr, err := NewTransformerBootstrapCI(
		[]string{"mean"},
		[]string{"x"},
		[]string{"g"},
		500,
		0.95,
		false,
	)
	if err != nil {
		t.Fatal(err)
	}

	outputRecordsAndContexts := feedBootstrapCI(tr, []int64{1, 2, 3, 4, 5})

	// One stats record plus the end-of-stream marker
	if len(outputRecordsAndContexts) != 2 {
		t.Fatalf("expected 2 output records-and-contexts, got %d", len(outputRecordsAndContexts))
	}
	if !outputRecordsAndContexts[1].EndOfStream {
		t.Fatal("expected second output to be the end-of-stream marker")
	}

	outrec := outputRecordsAndContexts[0].Record
	if got := outrec.Get("g"); got == nil || got.String() != "one" {
		t.Errorf("expected group-by field g=one, got %v", got)
	}

	pointEstimate := outrec.Get("x_mean")
	lo := outrec.Get("x_mean_lo")
	hi := outrec.Get("x_mean_hi")
	if pointEstimate == nil || lo == nil || hi == nil {
		t.Fatalf("missing output fields; record is %s", outrec.String())
	}

	if got, ok := pointEstimate.GetNumericToFloatValue(); !ok || got != 3.0 {
		t.Errorf("expected point estimate x_mean=3, got %v", pointEstimate)
	}

	loValue, ok := lo.GetNumericToFloatValue()
	if !ok {
		t.Fatalf("x_mean_lo is not numeric: %v", lo)
	}
	hiValue, ok := hi.GetNumericToFloatValue()
	if !ok {
		t.Fatalf("x_mean_hi is not numeric: %v", hi)
	}

	// The resampled means all lie within [min, max] of the data, and the
	// confidence interval must bracket the point estimate.
	if loValue < 1.0 || hiValue > 5.0 {
		t.Errorf("confidence interval [%v, %v] out of data range [1, 5]", loValue, hiValue)
	}
	if loValue > 3.0 || hiValue < 3.0 {
		t.Errorf("confidence interval [%v, %v] does not bracket point estimate 3", loValue, hiValue)
	}
	if loValue >= hiValue {
		t.Errorf("expected lo < hi, got [%v, %v]", loValue, hiValue)
	}
}

func TestBootstrapCISingleValueDegenerate(t *testing.T) {
	lib.SeedRandom(12345)

	tr, err := NewTransformerBootstrapCI(
		[]string{"mean"},
		[]string{"x"},
		[]string{"g"},
		100,
		0.95,
		false,
	)
	if err != nil {
		t.Fatal(err)
	}

	outputRecordsAndContexts := feedBootstrapCI(tr, []int64{7})

	if len(outputRecordsAndContexts) != 2 {
		t.Fatalf("expected 2 output records-and-contexts, got %d", len(outputRecordsAndContexts))
	}
	outrec := outputRecordsAndContexts[0].Record

	// With a single datum, every resample is that datum, so the interval is
	// degenerate at the point estimate.
	for _, key := range []string{"x_mean", "x_mean_lo", "x_mean_hi"} {
		value := outrec.Get(key)
		if value == nil {
			t.Fatalf("missing output field %s; record is %s", key, outrec.String())
		}
		if got, ok := value.GetNumericToFloatValue(); !ok || got != 7.0 {
			t.Errorf("expected %s=7, got %v", key, value)
		}
	}
}
