package cli

import (
	"slices"
	"testing"
)

func TestFlagValueCandidates(t *testing.T) {
	// Format flags offer file-format names.
	for _, flag := range []string{"-i", "-o", "--io"} {
		got := FlagValueCandidates(flag)
		if !slices.Contains(got, "csv") || !slices.Contains(got, "json") {
			t.Errorf("%s: expected file-format names, got %v", flag, got)
		}
	}

	// Separator flags offer separator aliases.
	got := FlagValueCandidates("--ifs")
	if !slices.Contains(got, "comma") || !slices.Contains(got, "pipe") {
		t.Errorf("--ifs: expected separator aliases, got %v", got)
	}

	// Regex-separator flags offer regex aliases.
	got = FlagValueCandidates("--ifs-regex")
	if !slices.Contains(got, "whitespace") {
		t.Errorf("--ifs-regex: expected regex aliases, got %v", got)
	}

	// Flags without an enumerable value set return nil.
	for _, flag := range []string{"--ofmt", "--from", "-n", "--not-a-flag"} {
		if got := FlagValueCandidates(flag); got != nil {
			t.Errorf("%s: expected nil, got %v", flag, got)
		}
	}
}
