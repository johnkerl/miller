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

// TestWhichScoreNameMatch verifies that a token found in the name scores higher
// than one found only in the body.
func TestWhichScoreNameMatch(t *testing.T) {
	tokens := []string{"join"}
	nameScore := whichScore(tokens, "join", "combine two streams")
	bodyOnlyScore := whichScore(tokens, "combine", "join two records")
	if nameScore <= bodyOnlyScore {
		t.Errorf("name match (%d) should outscore body-only match (%d)", nameScore, bodyOnlyScore)
	}
	if nameScore != whichNameMatchScore {
		t.Errorf("single-token name match: got %d, want %d", nameScore, whichNameMatchScore)
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
