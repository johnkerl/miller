package transformers

import (
	"encoding/json"
	"testing"
)

// TestMigratedVerbsHaveOptions checks that verbs which have been given
// structured Options emit them in their JSON catalog entry (non-nil, non-empty
// for verbs that actually have options).
func TestMigratedVerbsHaveOptions(t *testing.T) {
	cases := []struct {
		verb            string
		wantNonNilOpts  bool
		wantMinOptCount int
	}{
		{"nothing", true, 0}, // explicitly migrated, no verb-specific options
		{"cat", true, 5},
		{"head", true, 2},
		{"tail", true, 2},
		{"tee", true, 2},
	}

	for _, tc := range cases {
		info := GetVerbInfoForJSON(tc.verb)
		if info == nil {
			t.Errorf("verb %q not found in catalog", tc.verb)
			continue
		}
		if tc.wantNonNilOpts && info.Options == nil {
			t.Errorf("verb %q: Options is nil, want non-nil (migrated)", tc.verb)
			continue
		}
		if len(info.Options) < tc.wantMinOptCount {
			t.Errorf("verb %q: got %d options, want >= %d", tc.verb, len(info.Options), tc.wantMinOptCount)
		}
	}
}

// TestUnmigratedVerbHasNoOptions checks that an unmigrated verb still emits
// nil Options (omitted from JSON) so agents fall back to usage_text.
func TestUnmigratedVerbHasNoOptions(t *testing.T) {
	// stats1 is intentionally not migrated in this PR.
	info := GetVerbInfoForJSON("stats1")
	if info == nil {
		t.Fatal("stats1 not found in catalog")
	}
	if info.Options != nil {
		t.Errorf("stats1: Options should be nil for unmigrated verb, got %v", info.Options)
	}
	if info.UsageText == "" {
		t.Error("stats1: UsageText should be non-empty for unmigrated verb")
	}
}

// TestOptionSpecFieldsPopulated verifies that migrated OptionSpec entries have
// the required fields filled in.
func TestOptionSpecFieldsPopulated(t *testing.T) {
	info := GetVerbInfoForJSON("cat")
	if info == nil {
		t.Fatal("cat not found")
	}
	for _, opt := range info.Options {
		if opt.Flag == "" {
			t.Errorf("cat: OptionSpec with empty Flag: %+v", opt)
		}
		if opt.Type == "" {
			t.Errorf("cat: OptionSpec %q has empty Type", opt.Flag)
		}
		if opt.Desc == "" {
			t.Errorf("cat: OptionSpec %q has empty Desc", opt.Flag)
		}
	}
}

// TestOptionsRoundTripJSON verifies that OptionSpec survives JSON
// marshal/unmarshal with all fields intact.
func TestOptionsRoundTripJSON(t *testing.T) {
	info := GetVerbInfoForJSON("cat")
	if info == nil {
		t.Fatal("cat not found")
	}

	b, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var roundTrip VerbInfoForJSON
	if err := json.Unmarshal(b, &roundTrip); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(roundTrip.Options) != len(info.Options) {
		t.Errorf("options count: got %d, want %d", len(roundTrip.Options), len(info.Options))
	}
	for i, orig := range info.Options {
		got := roundTrip.Options[i]
		if got.Flag != orig.Flag || got.Type != orig.Type || got.Desc != orig.Desc {
			t.Errorf("option[%d] mismatch: got %+v, want %+v", i, got, orig)
		}
	}
}

// TestUnmigratedVerbHasNilOptionsInJSON verifies that nil Options are omitted
// from JSON (not serialized as null), so agents can use key presence to test
// Tier-2 availability.
func TestUnmigratedVerbHasNilOptionsInJSON(t *testing.T) {
	info := GetVerbInfoForJSON("stats1")
	if info == nil {
		t.Fatal("stats1 not found")
	}
	b, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// The JSON must not contain an "options" key for an unmigrated verb.
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
	if _, ok := raw["options"]; ok {
		t.Error("unmigrated verb stats1 should not have \"options\" key in JSON")
	}
}

// TestVerbOptionsNilCheckRuns ensures VerbOptionsNilCheck doesn't panic and
// reports expected migrated verbs as migrated.
func TestVerbOptionsNilCheckRuns(t *testing.T) {
	// Just verify it doesn't panic; output goes to stdout which test harness ignores.
	// The actual count is asserted via the regression test case.
	VerbOptionsNilCheck()
}
