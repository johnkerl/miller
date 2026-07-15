package transformers

import (
	"testing"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

func TestParseStepperCount(t *testing.T) {
	cases := []struct {
		stepperName string
		baseName    string
		wantCount   int
		wantOK      bool
	}{
		{"shift_lag", "shift_lag", 1, true},
		{"shift_lag_1", "shift_lag", 1, true},
		{"shift_lag_12", "shift_lag", 12, true},
		{"shift_lag_007", "shift_lag", 7, true},
		{"shift_lag_0", "shift_lag", 0, false},
		{"shift_lag_-3", "shift_lag", 0, false},
		{"shift_lag_+3", "shift_lag", 0, false},
		{"shift_lag_", "shift_lag", 0, false},
		{"shift_lag_x", "shift_lag", 0, false},
		{"shift_lag_3_4", "shift_lag", 0, false},
		{"shift_lead", "shift_lag", 0, false},
		{"shift_lag", "shift", 0, false}, // "lag" is not a count
		{"shift_12", "shift", 12, true},
		{"delta_4", "delta", 4, true},
		{"ratio_4", "ratio", 4, true},
	}

	for _, tc := range cases {
		count, ok := parseStepperCount(tc.stepperName, tc.baseName)
		if ok != tc.wantOK || (ok && count != tc.wantCount) {
			t.Errorf(
				"parseStepperCount(%q, %q) = (%d, %v); want (%d, %v)",
				tc.stepperName, tc.baseName, count, ok, tc.wantCount, tc.wantOK,
			)
		}
	}
}

func TestStepperNameHasBadCount(t *testing.T) {
	cases := []struct {
		stepperName string
		want        bool
	}{
		{"delta_0", true},
		{"delta_-1", true},
		{"delta_x", true},
		{"delta_3_4", true},
		{"shift_lag_0", true},
		{"shift_lead_-2", true},
		{"ratio_", true},
		{"shift_", true},
		{"delta", false},   // valid; never reaches the error path anyway
		{"delta_7", false}, // valid; never reaches the error path anyway
		{"foo", false},
		{"foo_0", false},
		{"slwin_x_y", false}, // slwin has its own handling
	}

	for _, tc := range cases {
		got := stepperNameHasBadCount(tc.stepperName)
		if got != tc.want {
			t.Errorf("stepperNameHasBadCount(%q) = %v; want %v", tc.stepperName, got, tc.want)
		}
	}
}

func TestStepperInputFromNameWithCounts(t *testing.T) {
	cases := []struct {
		stepperName  string
		wantFound    bool
		wantBackward int
		wantForward  int
	}{
		{"shift_lag", true, 0, 0},
		{"shift_lag_12", true, 0, 0},
		{"shift_lead", true, 0, 1},
		{"shift_lead_4", true, 0, 4},
		{"shift", true, 0, 0},
		{"shift_3", true, 0, 0},
		{"delta", true, 0, 0},
		{"delta_2", true, 0, 0},
		{"ratio_2", true, 0, 0},
		{"slwin_2_3", true, 2, 3},
		{"shift_lag_0", false, 0, 0},
		{"delta_x", false, 0, 0},
		{"nonesuch", false, 0, 0},
	}

	for _, tc := range cases {
		stepperInput := stepperInputFromName(tc.stepperName)
		if tc.wantFound {
			if stepperInput == nil {
				t.Errorf("stepperInputFromName(%q) = nil; want non-nil", tc.stepperName)
				continue
			}
			if stepperInput.numRecordsBackward != tc.wantBackward ||
				stepperInput.numRecordsForward != tc.wantForward {
				t.Errorf(
					"stepperInputFromName(%q) = {backward:%d, forward:%d}; want {backward:%d, forward:%d}",
					tc.stepperName,
					stepperInput.numRecordsBackward, stepperInput.numRecordsForward,
					tc.wantBackward, tc.wantForward,
				)
			}
		} else {
			if stepperInput != nil {
				t.Errorf("stepperInputFromName(%q) = %v; want nil", tc.stepperName, stepperInput)
			}
		}
	}
}

func TestAllocateStepperOutputFieldNames(t *testing.T) {
	cases := []struct {
		stepperName         string
		wantOutputFieldName string
	}{
		{"shift", "x_shift"},
		{"shift_lag", "x_shift_lag"},
		{"shift_lag_12", "x_shift_lag_12"},
		{"shift_lead", "x_shift_lead"},
		{"shift_lead_4", "x_shift_lead_4"},
		{"delta", "x_delta"},
		{"delta_2", "x_delta_2"},
		{"ratio", "x_ratio"},
		{"ratio_2", "x_ratio_2"},
	}

	for _, tc := range cases {
		stepperInput := stepperInputFromName(tc.stepperName)
		if stepperInput == nil {
			t.Errorf("stepperInputFromName(%q) = nil; want non-nil", tc.stepperName)
			continue
		}
		stepper, err := allocateStepper(stepperInput, "x", nil, nil)
		if err != nil {
			t.Errorf("allocateStepper for %q: %v", tc.stepperName, err)
			continue
		}
		if stepper == nil {
			t.Errorf("allocateStepper for %q = nil; want non-nil", tc.stepperName)
			continue
		}
		var got string
		switch s := stepper.(type) {
		case *tStepperShiftLag:
			got = s.outputFieldName
		case *tStepperShiftLead:
			got = s.outputFieldName
		case *tStepperDelta:
			got = s.outputFieldName
		case *tStepperRatio:
			got = s.outputFieldName
		default:
			t.Errorf("allocateStepper for %q: unexpected stepper type %T", tc.stepperName, stepper)
			continue
		}
		if got != tc.wantOutputFieldName {
			t.Errorf(
				"allocateStepper for %q: output field name %q; want %q",
				tc.stepperName, got, tc.wantOutputFieldName,
			)
		}
	}
}

func TestValueRing(t *testing.T) {
	ring := newValueRing(3)

	// Fewer than 3 values seen: has must be false.
	for i := 1; i <= 3; i++ {
		nBack, has := ring.push(mlrval.FromInt(int64(i)))
		if has {
			t.Errorf("push %d: has = true; want false", i)
		}
		if nBack != nil {
			t.Errorf("push %d: nBack = %v; want nil", i, nBack)
		}
	}

	// From the 4th value on, the value from 3 records back is returned.
	for i := 4; i <= 8; i++ {
		nBack, has := ring.push(mlrval.FromInt(int64(i)))
		if !has {
			t.Errorf("push %d: has = false; want true", i)
		}
		if nBack == nil {
			t.Errorf("push %d: nBack = nil; want non-nil", i)
			continue
		}
		want := int64(i - 3)
		got, ok := nBack.GetIntValue()
		if !ok || got != want {
			t.Errorf("push %d: nBack = %v; want %d", i, nBack, want)
		}
	}

	// A nil push (absent field) occupies a slot.
	ring.push(nil)                              // 9th slot: nil
	ring.push(mlrval.FromInt(10))               // 10th
	ring.push(mlrval.FromInt(11))               // 11th
	nBack, has := ring.push(mlrval.FromInt(12)) // 12th: 3 back is the nil slot
	if !has {
		t.Errorf("push 12: has = false; want true")
	}
	if nBack != nil {
		t.Errorf("push 12: nBack = %v; want nil", nBack)
	}
}
