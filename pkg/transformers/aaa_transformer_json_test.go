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

// TestAllVerbsFullyMigrated asserts the full-migration invariant: every verb
// in the catalog has a non-nil Options slice and a non-empty UsageText.
// This is the Tier-2 completion check; it fails if a new verb is added without
// populating Options.
func TestAllVerbsFullyMigrated(t *testing.T) {
	for i := range TRANSFORMER_LOOKUP_TABLE {
		setup := &TRANSFORMER_LOOKUP_TABLE[i]
		if setup.Options == nil {
			t.Errorf("verb %q has nil Options (not yet migrated to Tier-2)", setup.Verb)
		}
		info := GetVerbInfoForJSON(setup.Verb)
		if info == nil {
			t.Errorf("verb %q: GetVerbInfoForJSON returned nil", setup.Verb)
			continue
		}
		if info.UsageText == "" {
			t.Errorf("verb %q: UsageText is empty", setup.Verb)
		}
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

// TestAllVerbsHaveOptionsKeyInJSON verifies that every fully-migrated verb
// emits an "options" key in its JSON output (even when the list is empty),
// so agents can rely on key presence as the Tier-2 signal.
func TestAllVerbsHaveOptionsKeyInJSON(t *testing.T) {
	for i := range TRANSFORMER_LOOKUP_TABLE {
		setup := &TRANSFORMER_LOOKUP_TABLE[i]
		if setup.Options == nil {
			continue // not yet migrated; skip rather than double-fail
		}
		info := GetVerbInfoForJSON(setup.Verb)
		if info == nil {
			t.Errorf("verb %q: GetVerbInfoForJSON returned nil", setup.Verb)
			continue
		}
		b, err := json.Marshal(info)
		if err != nil {
			t.Errorf("verb %q: marshal error: %v", setup.Verb, err)
			continue
		}
		var raw map[string]any
		if err := json.Unmarshal(b, &raw); err != nil {
			t.Errorf("verb %q: unmarshal error: %v", setup.Verb, err)
			continue
		}
		if _, ok := raw["options"]; !ok {
			t.Errorf("verb %q: missing \"options\" key in JSON", setup.Verb)
		}
	}
}
