package transformers

import (
	"testing"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

// makeStats1TestRecord builds a record from alternating key/value pairs, e.g.
// makeStats1TestRecord("x", "1", "g", "a").
func makeStats1TestRecord(kvPairs ...string) *types.RecordAndContext {
	record := mlrval.NewMlrmapAsRecord()
	for i := 0; i < len(kvPairs); i += 2 {
		record.PutReference(kvPairs[i], mlrval.FromInferredType(kvPairs[i+1]))
	}
	context := types.NewNilContext()
	return types.NewRecordAndContext(record, context)
}

func runStats1TestTransformer(
	tr RecordTransformer,
	inputs []*types.RecordAndContext,
) []*types.RecordAndContext {
	inputDownstreamDoneChannel := make(chan bool, 1)
	outputDownstreamDoneChannel := make(chan bool, 1)
	outputs := make([]*types.RecordAndContext, 0)
	for _, input := range inputs {
		tr.Transform(input, &outputs, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	}
	// End-of-stream marker
	eos := types.NewRecordAndContext(nil, types.NewNilContext())
	eos.EndOfStream = true
	tr.Transform(eos, &outputs, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	return outputs
}

// TestStats1SlidingWindow exercises 'mlr stats1 -a sum,mean,min,max -f x -w 3'.
func TestStats1SlidingWindow(t *testing.T) {
	tr, err := NewTransformerStats1(
		[]string{"sum", "mean", "min", "max"},
		[]string{"x"},
		[]string{},   // groupByFieldNameList
		false, false, // doRegexValueFieldNames, doRegexGroupByFieldNames
		false, false, // invertRegexValueFieldNames, invertRegexGroupByFieldNames
		false, // doInterpolatedPercentiles
		false, // doIterativeStats
		3,     // slidingWindowSize
	)
	if err != nil {
		t.Fatal(err)
	}

	inputs := []*types.RecordAndContext{
		makeStats1TestRecord("x", "1"),
		makeStats1TestRecord("x", "2"),
		makeStats1TestRecord("x", "3"),
		makeStats1TestRecord("x", "4"),
		makeStats1TestRecord("x", "5"),
	}
	outputs := runStats1TestTransformer(tr, inputs)

	// 5 data records plus end-of-stream marker
	if len(outputs) != 6 {
		t.Fatalf("got %d outputs, want 6", len(outputs))
	}

	expectations := []struct{ sum, mean, min, max string }{
		{"1", "1", "1", "1"},   // window [1]
		{"3", "1.5", "1", "2"}, // window [1,2]
		{"6", "2", "1", "3"},   // window [1,2,3]
		{"9", "3", "2", "4"},   // window [2,3,4]
		{"12", "4", "3", "5"},  // window [3,4,5]
	}
	for i, expectation := range expectations {
		record := outputs[i].Record
		for name, want := range map[string]string{
			"x_sum":  expectation.sum,
			"x_mean": expectation.mean,
			"x_min":  expectation.min,
			"x_max":  expectation.max,
		} {
			value := record.Get(name)
			if value == nil {
				t.Fatalf("record %d: field %s missing", i, name)
			}
			if value.String() != want {
				t.Errorf("record %d: %s got %s, want %s", i, name, value.String(), want)
			}
		}
	}

	if !outputs[5].EndOfStream {
		t.Errorf("last output should be the end-of-stream marker")
	}
}

// TestStats1SlidingWindowGrouped exercises 'mlr stats1 -a count,sum -f x -g g -w 2':
// windows are maintained per group.
func TestStats1SlidingWindowGrouped(t *testing.T) {
	tr, err := NewTransformerStats1(
		[]string{"count", "sum"},
		[]string{"x"},
		[]string{"g"}, // groupByFieldNameList
		false, false,  // doRegexValueFieldNames, doRegexGroupByFieldNames
		false, false, // invertRegexValueFieldNames, invertRegexGroupByFieldNames
		false, // doInterpolatedPercentiles
		false, // doIterativeStats
		2,     // slidingWindowSize
	)
	if err != nil {
		t.Fatal(err)
	}

	inputs := []*types.RecordAndContext{
		makeStats1TestRecord("g", "a", "x", "1"),
		makeStats1TestRecord("g", "b", "x", "10"),
		makeStats1TestRecord("g", "a", "x", "2"),
		makeStats1TestRecord("g", "b", "x", "20"),
		makeStats1TestRecord("g", "a", "x", "3"),
	}
	outputs := runStats1TestTransformer(tr, inputs)

	// 5 data records plus end-of-stream marker
	if len(outputs) != 6 {
		t.Fatalf("got %d outputs, want 6", len(outputs))
	}

	expectations := []struct{ count, sum string }{
		{"1", "1"},  // group a, window [1]
		{"1", "10"}, // group b, window [10]
		{"2", "3"},  // group a, window [1,2]
		{"2", "30"}, // group b, window [10,20]
		{"2", "5"},  // group a, window [2,3]
	}
	for i, expectation := range expectations {
		record := outputs[i].Record
		count := record.Get("x_count")
		sum := record.Get("x_sum")
		if count == nil || sum == nil {
			t.Fatalf("record %d: missing x_count or x_sum", i)
		}
		if count.String() != expectation.count {
			t.Errorf("record %d: x_count got %s, want %s", i, count.String(), expectation.count)
		}
		if sum.String() != expectation.sum {
			t.Errorf("record %d: x_sum got %s, want %s", i, sum.String(), expectation.sum)
		}
	}
}
