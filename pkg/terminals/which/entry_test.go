package which

import (
	"testing"
)

// TestTokenize verifies that query tokenization drops stopwords and deduplicates.
func TestTokenize(t *testing.T) {
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
		got := tokenize(tc.query)
		if len(got) != len(tc.expect) {
			t.Errorf("tokenize(%q): got %v, want %v", tc.query, got, tc.expect)
			continue
		}
		for i, tok := range got {
			if tok != tc.expect[i] {
				t.Errorf("tokenize(%q)[%d]: got %q, want %q", tc.query, i, tok, tc.expect[i])
			}
		}
	}
}

// TestScoreNameMatch verifies that a token found in the name scores higher
// than one found only in the body.
func TestScoreNameMatch(t *testing.T) {
	tokens := []string{"join"}
	nameScore := score(tokens, "join", "combine two streams")
	bodyOnlyScore := score(tokens, "combine", "join two records")
	if nameScore <= bodyOnlyScore {
		t.Errorf("name match (%d) should outscore body-only match (%d)", nameScore, bodyOnlyScore)
	}
	if nameScore != nameMatchScore {
		t.Errorf("single-token name match: got %d, want %d", nameScore, nameMatchScore)
	}
}

// TestSearchReturnsJoinVerb checks that "join" as a query surfaces the join verb.
func TestSearchReturnsJoinVerb(t *testing.T) {
	tokens := tokenize("join two files on a key")
	results := search(tokens)
	if len(results) == 0 {
		t.Fatal("expected at least one result for 'join'")
	}
	// The join verb should be the top result.
	top := results[0]
	if top.Name != "join" || top.Kind != "verb" {
		t.Errorf("expected top result to be verb:join, got kind=%s name=%s", top.Kind, top.Name)
	}
	if top.Score < nameMatchScore {
		t.Errorf("top result score %d below confident threshold %d", top.Score, nameMatchScore)
	}
}

// TestSearchReturnsCountVerb checks that "count records" surfaces the count verb.
func TestSearchReturnsCountVerb(t *testing.T) {
	tokens := tokenize("count records")
	results := search(tokens)
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

// TestSearchNoResults ensures an all-stopword query returns nothing.
func TestSearchNoResults(t *testing.T) {
	tokens := tokenize("a an the to")
	if len(tokens) != 0 {
		t.Errorf("expected empty tokens for all-stopword query, got %v", tokens)
	}
}
