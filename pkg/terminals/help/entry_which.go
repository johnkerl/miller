// mlr which -- intent-to-capability router for AI agents and interactive users.
//
// Usage:
//
//	mlr which "natural language query"
//
// Searches verb names and summaries (and other catalog items) for query-word
// matches, emits a ranked JSON array, and exits with:
//
//	0  — at least one result whose name contains a query token (confident match)
//	1  — usage / argument error
//	2  — no result scored a name-level match (low confidence)
//
// The exit-code contract lets an agent branch on status rather than parsing
// the prose, while the result array is still useful even on exit code 2.
//
// Lives here (alongside the other --as-json help machinery) rather than in its
// own package because it imports the same four catalog registries and shares
// the firstLine helper with entry_json.go.

package help

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/dsl/cst"
	"github.com/johnkerl/miller/v6/pkg/transformers"
)

// WhichResultEntry is one ranked match returned by `mlr which`.
type WhichResultEntry struct {
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Score   int    `json:"score"`
	Summary string `json:"summary"`
}

// WhichMain is the entrypoint called by the terminals dispatcher for `mlr which`.
func WhichMain(args []string) int {
	args = args[1:] // strip "which"

	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintf(os.Stderr, "Usage: mlr which \"natural language query\"\n")
		fmt.Fprintf(os.Stderr, "Searches verb names, summaries, and other catalog items for query-word matches.\n")
		fmt.Fprintf(os.Stderr, "Emits a JSON array of {kind, name, score, summary} sorted by descending score.\n")
		fmt.Fprintf(os.Stderr, "Exit codes: 0=confident match (name hit), 2=no confident match.\n")
		return 1
	}

	query := strings.Join(args, " ")
	tokens := whichTokenize(query)
	if len(tokens) == 0 {
		fmt.Fprintf(os.Stderr, "mlr which: empty query\n")
		return 1
	}

	results := whichSearch(tokens)

	bytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr which: could not render JSON: %v\n", err)
		return 1
	}
	fmt.Println(string(bytes))

	if len(results) > 0 && results[0].Score >= whichNameMatchScore {
		return 0
	}
	return 2
}

// whichNameMatchScore is the per-token weight for a name match; used as the
// confidence threshold for exit-code 0.
const whichNameMatchScore = 20

// whichSearch scores every catalog item against tokens and returns matches in
// descending score order. Items with score 0 are omitted.
func whichSearch(tokens []string) []WhichResultEntry {
	results := make([]WhichResultEntry, 0)

	for _, v := range transformers.GetVerbInfosForJSON() {
		if s := whichScore(tokens, v.Name, v.Summary+" "+v.UsageText); s > 0 {
			results = append(results, WhichResultEntry{
				Kind: "verb", Name: v.Name, Score: s, Summary: v.Summary,
			})
		}
	}

	for _, f := range cst.BuiltinFunctionManagerInstance.GetFunctionInfosForJSON() {
		if s := whichScore(tokens, f.Name, f.Help); s > 0 {
			results = append(results, WhichResultEntry{
				Kind: "function", Name: f.Name, Score: s, Summary: firstLine(f.Help),
			})
		}
	}

	for _, fl := range cli.FLAG_TABLE.GetFlagInfosForJSON() {
		nameText := fl.Name + " " + strings.Join(fl.AltNames, " ")
		if s := whichScore(tokens, nameText, fl.Help); s > 0 {
			results = append(results, WhichResultEntry{
				Kind: "flag", Name: fl.Name, Score: s, Summary: fl.Help,
			})
		}
	}

	for _, kw := range cst.GetKeywordInfosForJSON() {
		if s := whichScore(tokens, kw.Name, kw.Help); s > 0 {
			results = append(results, WhichResultEntry{
				Kind: "keyword", Name: kw.Name, Score: s, Summary: firstLine(kw.Help),
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		// Stable tiebreak: verbs first (most useful for agents), then alphabetical name.
		ki, kj := whichKindRank(results[i].Kind), whichKindRank(results[j].Kind)
		if ki != kj {
			return ki < kj
		}
		return results[i].Name < results[j].Name
	})

	return results
}

// whichScore sums per-token weights: whichNameMatchScore per token found in
// name, 5 per token found in body. Matching is case-insensitive substring.
func whichScore(tokens []string, name, body string) int {
	lname := strings.ToLower(name)
	lbody := strings.ToLower(body)
	total := 0
	for _, tok := range tokens {
		if strings.Contains(lname, tok) {
			total += whichNameMatchScore
		} else if strings.Contains(lbody, tok) {
			total += 5
		}
	}
	return total
}

// whichTokenize lowercases and splits a query into non-trivial words, dropping
// single-character tokens and common stopwords that carry no discriminating
// signal against Miller's catalog.
func whichTokenize(query string) []string {
	stopwords := map[string]bool{
		"a": true, "an": true, "the": true, "to": true, "of": true,
		"in": true, "on": true, "at": true, "by": true, "for": true,
		"and": true, "or": true, "is": true, "it": true, "do": true,
		"with": true, "from": true, "into": true, "how": true,
		"get": true, "use": true, "two": true, "my": true,
	}
	words := strings.FieldsFunc(strings.ToLower(query), func(r rune) bool {
		return !('a' <= r && r <= 'z') && !('0' <= r && r <= '9') && r != '-' && r != '_'
	})
	var tokens []string
	seen := map[string]bool{}
	for _, w := range words {
		if len(w) <= 1 || stopwords[w] || seen[w] {
			continue
		}
		seen[w] = true
		tokens = append(tokens, w)
	}
	return tokens
}

func whichKindRank(kind string) int {
	switch kind {
	case "verb":
		return 0
	case "function":
		return 1
	case "flag":
		return 2
	case "keyword":
		return 3
	}
	return 4
}
