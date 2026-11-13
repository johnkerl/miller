package input

import (
	"slices"
	"testing"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

func TestHeaderFieldNames(t *testing.T) {
	testCases := []struct {
		name   string
		header []string
		dedupe bool
		want   []string
	}{
		{"unique", []string{"a", "b", "c"}, true, []string{"a", "b", "c"}},
		{"empty", []string{}, true, []string{}},
		{"repeated", []string{"name", "4", "4"}, true, []string{"name", "4", "4_2"}},
		{"three times", []string{"a", "a", "a"}, true, []string{"a", "a_2", "a_3"}},
		{"suffix already taken", []string{"a", "a_2", "a"}, true, []string{"a", "a_2", "a_3"}},
		{"no dedupe", []string{"a", "b", "a"}, false, []string{"a", "b"}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := headerFieldNames(testCase.header, testCase.dedupe)
			if !slices.Equal(got, testCase.want) {
				t.Errorf("headerFieldNames(%v, %v) = %v, want %v",
					testCase.header, testCase.dedupe, got, testCase.want)
			}

			// The result must match the keys of a record built from the same
			// header line, since unsparsify --from-header compares the two.
			record := mlrval.NewMlrmapAsRecord()
			for _, key := range testCase.header {
				if _, err := record.PutReferenceMaybeDedupe(key, mlrval.VOID, testCase.dedupe); err != nil {
					t.Fatal(err)
				}
			}
			recordKeys := record.GetKeys()
			if !slices.Equal(got, recordKeys) {
				t.Errorf("headerFieldNames(%v, %v) = %v, but record keys are %v",
					testCase.header, testCase.dedupe, got, recordKeys)
			}
		})
	}
}
