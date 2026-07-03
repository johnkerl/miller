package climain

import (
	"fmt"
	"testing"
)

// ----------------------------------------------------------------
// Levenshtein

func TestLevenshteinIdentical(t *testing.T) {
	if d := levenshtein("filter", "filter"); d != 0 {
		t.Errorf("identical strings: got %d, want 0", d)
	}
}

func TestLevenshteinEmpty(t *testing.T) {
	if d := levenshtein("", "abc"); d != 3 {
		t.Errorf("empty a: got %d, want 3", d)
	}
	if d := levenshtein("abc", ""); d != 3 {
		t.Errorf("empty b: got %d, want 3", d)
	}
}

func TestLevenshteinSubstitution(t *testing.T) {
	// "fliter" → "filter": swap l↔i at positions 1-2 → distance 2
	d := levenshtein("fliter", "filter")
	if d != 2 {
		t.Errorf("fliter/filter: got %d, want 2", d)
	}
}

func TestLevenshteinInsertion(t *testing.T) {
	// "--jsonn" → "--json": 1 deletion
	if d := levenshtein("--jsonn", "--json"); d != 1 {
		t.Errorf("--jsonn/--json: got %d, want 1", d)
	}
}

func TestLevenshteinDeletion(t *testing.T) {
	// "--complemment" → "--complement": 1 extra m
	if d := levenshtein("--complemment", "--complement"); d != 1 {
		t.Errorf("--complemment/--complement: got %d, want 1", d)
	}
}

// ----------------------------------------------------------------
// topMatches

func TestTopMatchesBasic(t *testing.T) {
	// "filtter" is distance 1 from "filter" (extra t), clearly the top match.
	candidates := []string{"filter", "flatten", "fraction", "grep", "head"}
	matches := topMatches("filtter", candidates, 3, 2)
	if len(matches) == 0 {
		t.Fatal("expected at least one match for 'filtter'")
	}
	if matches[0] != "filter" {
		t.Errorf("expected 'filter' as top match, got %q", matches[0])
	}
}

func TestTopMatchesNoneWithinThreshold(t *testing.T) {
	candidates := []string{"aaaaa", "bbbbb", "ccccc"}
	matches := topMatches("xyz", candidates, 3, 1)
	if len(matches) != 0 {
		t.Errorf("expected no matches, got %v", matches)
	}
}

func TestTopMatchesRespectsCap(t *testing.T) {
	// All candidates are within distance 1 of "a"; cap at 2
	candidates := []string{"aa", "ab", "ac", "ad", "ae"}
	matches := topMatches("a", candidates, 2, 2)
	if len(matches) > 2 {
		t.Errorf("expected at most 2 matches, got %d", len(matches))
	}
}

// ----------------------------------------------------------------
// levenshteinThreshold

func TestThresholdShortQuery(t *testing.T) {
	if levenshteinThreshold("ab") != 1 {
		t.Error("2-char query should have threshold 1")
	}
}

func TestThresholdMediumQuery(t *testing.T) {
	if levenshteinThreshold("filter") != 2 {
		t.Error("6-char query should have threshold 2")
	}
}

func TestThresholdLongQuery(t *testing.T) {
	if levenshteinThreshold("--complement") != 3 {
		t.Error("12-char query should have threshold 3")
	}
}

// ----------------------------------------------------------------
// WantErrorsJSON

func TestWantErrorsJSONFlag(t *testing.T) {
	if !WantErrorsJSON([]string{"mlr", "--errors-json", "flitre"}) {
		t.Error("should detect --errors-json")
	}
	if !WantErrorsJSON([]string{"mlr", "flitre", "--errors-json"}) {
		t.Error("should detect --errors-json in any position")
	}
	if WantErrorsJSON([]string{"mlr", "flitre"}) {
		t.Error("should not detect --errors-json when absent")
	}
}

// ----------------------------------------------------------------
// CLIError round-trip

func TestCLIErrorInterface(t *testing.T) {
	err := &CLIError{Kind: "unknown-verb", Token: "foo", Msg: "mlr: verb \"foo\" not found"}
	if err.Error() != err.Msg {
		t.Errorf("Error() should return Msg, got %q", err.Error())
	}
}

// ----------------------------------------------------------------
// categorize

func TestCategorizeUnknownVerb(t *testing.T) {
	err := &CLIError{Kind: "unknown-verb", Token: "flitre", Msg: "mlr: verb \"flitre\" not found"}
	se := categorize(err)
	if se.Kind != "unknown-verb" {
		t.Errorf("kind: got %q, want unknown-verb", se.Kind)
	}
	if se.Token != "flitre" {
		t.Errorf("token: got %q, want flitre", se.Token)
	}
	if se.Hint == "" {
		t.Error("hint should be non-empty for unknown-verb")
	}
}

func TestCategorizeUnknownFlag(t *testing.T) {
	err := &CLIError{Kind: "unknown-flag", Token: "--jsonn", Msg: "mlr: option \"--jsonn\" not recognized"}
	se := categorize(err)
	if se.Kind != "unknown-flag" {
		t.Errorf("kind: got %q, want unknown-flag", se.Kind)
	}
	if se.Hint == "" {
		t.Error("hint should be non-empty for unknown-flag")
	}
}

func TestCategorizeVerbOptionError(t *testing.T) {
	err := &CLIError{
		Kind: "verb-option-error", Token: "--bad", Verb: "cut",
		Msg: "mlr cut: option \"--bad\" not recognized",
	}
	se := categorize(err)
	if se.Kind != "verb-option-error" {
		t.Errorf("kind: got %q, want verb-option-error", se.Kind)
	}
	if se.Verb != "cut" {
		t.Errorf("verb: got %q, want cut", se.Verb)
	}
	if se.Hint == "" {
		t.Error("hint should be non-empty for verb-option-error")
	}
}

func TestCategorizeGenericFallback(t *testing.T) {
	err := fmt.Errorf("some unexpected error")
	se := categorize(err)
	if se.Kind != "generic" {
		t.Errorf("kind: got %q, want generic", se.Kind)
	}
}

func TestCategorizeDSLParseError(t *testing.T) {
	// The DSL parser (pkg/parsing/parser) emits bare "parse error: ..."
	// messages; these should categorize as dsl-parse-error, e.g. for
	// `mlr put --explain` with --errors-json.
	err := fmt.Errorf("parse error: unexpected equals (\"=\")")
	se := categorize(err)
	if se.Kind != "dsl-parse-error" {
		t.Errorf("kind: got %q, want dsl-parse-error", se.Kind)
	}
	if se.Hint == "" {
		t.Error("hint should be non-empty for dsl-parse-error")
	}
}

func TestCategorizeDSLParseErrorNotCSV(t *testing.T) {
	// The CSV reader's "parse error on line ..." is a stream-time error and
	// should not be mistaken for a DSL parse error by the substring match.
	err := fmt.Errorf("parse error on line 3, column 5: bare \" in non-quoted-field")
	se := categorize(err)
	if se.Kind == "dsl-parse-error" {
		t.Errorf("kind: got dsl-parse-error, want non-DSL categorization for a CSV parse error")
	}
}
