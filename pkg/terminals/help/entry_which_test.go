package help

import (
	"testing"
)

// TestWhichTokenize verifies that query tokenization drops stopwords and deduplicates.
func TestWhichTokenize(t *testing.T) {
	cases := []struct {
		query  string
		expect []string
	}{
		{"join two files on a key", []string{"join", "files", "key"}},
		{"count distinct values", []string{"count", "distinct", "values"}},
		{"group records by field", []string{"group", "records", "field"}},
		{"", nil},
		{"a an the", nil},
	}
	for _, tc := range cases {
		got := whichTokenize(tc.query)
		if len(got) != len(tc.expect) {
			t.Errorf("whichTokenize(%q): got %v, want %v", tc.query, got, tc.expect)
			continue
		}
		for i, tok := range got {
			if tok != tc.expect[i] {
				t.Errorf("whichTokenize(%q)[%d]: got %q, want %q", tc.query, i, tok, tc.expect[i])
			}
		}
	}
}

// TestWhichScoreNameMatch verifies that a name hit scores higher than a body-only
// hit and that nameHit is set correctly.
func TestWhichScoreNameMatch(t *testing.T) {
	tokens := []string{"join"}
	nameScore, nameHit := whichScore(tokens, "join", "combine two streams")
	bodyScore, bodyHit := whichScore(tokens, "combine", "join two records")
	if nameScore <= bodyScore {
		t.Errorf("name match (%d) should outscore body-only match (%d)", nameScore, bodyScore)
	}
	if nameScore != whichNameMatchScore {
		t.Errorf("single-token name match: got %d, want %d", nameScore, whichNameMatchScore)
	}
	if !nameHit {
		t.Error("nameHit should be true when token matches name")
	}
	if bodyHit {
		t.Error("nameHit should be false when token matches only body")
	}
}

// TestWhichExitCodeRequiresNameHit verifies that body-only matches (however
// many) do not trigger the confident-match exit code. This guards against the
// failure mode where 4 body-only token hits (4×5=20) equal whichNameMatchScore
// and would incorrectly signal a confident match.
func TestWhichExitCodeRequiresNameHit(t *testing.T) {
	// Tokens that appear in the body but not in the name.
	tokens := []string{"statistics", "aggregate", "compute", "average"}
	name := "xyz-verb"
	body := "statistics aggregate compute average"
	score, nameHit := whichScore(tokens, name, body)
	if nameHit {
		t.Error("body-only hits should not set nameHit")
	}
	if score < whichNameMatchScore {
		t.Errorf("expected body score >= %d to demonstrate the false-positive risk, got %d", whichNameMatchScore, score)
	}
}

// TestWhichSearchReturnsJoinVerb checks that "join" as a query surfaces the join verb.
func TestWhichSearchReturnsJoinVerb(t *testing.T) {
	tokens := whichTokenize("join two files on a key")
	results := whichSearch(tokens)
	if len(results) == 0 {
		t.Fatal("expected at least one result for 'join'")
	}
	top := results[0]
	if top.Name != "join" || top.Kind != "verb" {
		t.Errorf("expected top result to be verb:join, got kind=%s name=%s", top.Kind, top.Name)
	}
	if top.Score < whichNameMatchScore {
		t.Errorf("top result score %d below confident threshold %d", top.Score, whichNameMatchScore)
	}
}

// TestWhichSearchReturnsCountVerb checks that "count records" surfaces the count verb.
func TestWhichSearchReturnsCountVerb(t *testing.T) {
	tokens := whichTokenize("count records")
	results := whichSearch(tokens)
	found := false
	for _, r := range results {
		if r.Kind == "verb" && r.Name == "count" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'count' verb in results for 'count records'")
	}
}

// TestWhichSearchAllStopwords ensures an all-stopword query returns nothing.
func TestWhichSearchAllStopwords(t *testing.T) {
	tokens := whichTokenize("a an the to")
	if len(tokens) != 0 {
		t.Errorf("expected empty tokens for all-stopword query, got %v", tokens)
	}
}
