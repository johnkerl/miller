// Structured error output for `mlr --errors-json`.
//
// When --errors-json is set (or MLR_ERRORS_JSON is truthy), parse-time errors
// are emitted as a JSON object to stderr instead of a plain text message.
// This lets AI agents and scripts branch on error kind rather than regex-
// matching English prose.
//
// JSON shape:
//   {
//     "error":        "mlr: verb \"flitre\" not found",
//     "kind":         "unknown-verb",
//     "token":        "flitre",
//     "verb":         "",          // set for verb-option errors
//     "hint":         "Run 'mlr filter --help' for usage.",
//     "did_you_mean": ["filter", "flatten"]
//   }

package climain

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/dsl/cst"
	"github.com/johnkerl/miller/v6/pkg/transformers"
)

// ----------------------------------------------------------------
// CLIError is the typed error returned by the pass-one parser for
// errors that carry structured metadata (unknown verb, unknown flag).
// Other parse errors still flow as plain errors; EmitStructuredError
// handles both.

type CLIError struct {
	Kind  string // "unknown-verb" | "unknown-flag"
	Token string // the unrecognized token
	Verb  string // enclosing verb context, if any
	Msg   string // human-readable message (same as what prose path would print)
}

func (e *CLIError) Error() string { return e.Msg }

// ----------------------------------------------------------------
// StructuredError is the JSON DTO emitted to stderr.

type StructuredError struct {
	Error      string   `json:"error"`
	Kind       string   `json:"kind"`
	Token      string   `json:"token,omitempty"`
	Verb       string   `json:"verb,omitempty"`
	Hint       string   `json:"hint,omitempty"`
	DidYouMean []string `json:"did_you_mean,omitempty"`
}

// ----------------------------------------------------------------
// Opt-in detection

// WantErrorsJSON returns true when the caller has opted in via the
// --errors-json flag anywhere in args, or a truthy MLR_ERRORS_JSON env var.
// It is called before ParseCommandLine so it can affect how parse errors
// are reported.
func WantErrorsJSON(args []string) bool {
	if isTruthyEnv(os.Getenv("MLR_ERRORS_JSON")) {
		return true
	}
	for _, arg := range args {
		if arg == "--errors-json" {
			return true
		}
	}
	return false
}

func isTruthyEnv(v string) bool {
	switch v {
	case "1", "true", "True", "TRUE", "yes", "Yes", "YES":
		return true
	}
	return false
}

// ----------------------------------------------------------------
// Error categorization and emission

// EmitStructuredError writes a JSON error document to stderr and exits 1.
// It is called in place of printError when --errors-json is active.
func EmitStructuredError(err error) {
	se := categorize(err)
	se.DidYouMean = nearMatches(se)

	b, jerr := json.MarshalIndent(se, "", "  ")
	if jerr != nil {
		// Fallback: plain text if we somehow can't marshal
		fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
		return
	}
	fmt.Fprintln(os.Stderr, string(b))
}

// categorize turns any error into a StructuredError, using the typed CLIError
// fields when available and string-matching as a fallback.
func categorize(err error) StructuredError {
	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		se := StructuredError{
			Error: cliErr.Msg,
			Kind:  cliErr.Kind,
			Token: cliErr.Token,
			Verb:  cliErr.Verb,
		}
		switch cliErr.Kind {
		case "unknown-verb":
			se.Hint = "Run 'mlr -l' for a list of verbs, or 'mlr help verb <name>' for details."
		case "unknown-flag":
			se.Hint = "Run 'mlr --help' for a list of main flags."
		case "verb-option-error":
			if cliErr.Verb != "" {
				se.Hint = "Run 'mlr " + cliErr.Verb + " --help' for a list of options."
			}
		}
		return se
	}

	// Fallback: categorize by message pattern.
	msg := err.Error()

	// "mlr {verb}: option "{flag}" not recognized"
	if strings.Contains(msg, ": option ") && strings.Contains(msg, "not recognized") {
		verb := extractVerbFromMsg(msg)
		token := extractQuotedToken(msg)
		return StructuredError{
			Error: msg,
			Kind:  "verb-option-error",
			Token: token,
			Verb:  verb,
			Hint:  fmt.Sprintf("Run 'mlr %s --help' for a list of options.", verb),
		}
	}

	// DSL parse errors. The "parse error:" prefix comes from the DSL parser
	// (pkg/parsing/parser); "cannot parse DSL"/"DSL expression" are the wrapper
	// messages. (The CSV reader's "parse error on line ..." is a stream-time
	// error and never reaches this command-line-parse categorizer.)
	if strings.Contains(msg, "cannot parse DSL") ||
		strings.Contains(msg, "DSL expression") ||
		strings.Contains(msg, "parse error:") {
		return StructuredError{
			Error: msg,
			Kind:  "dsl-parse-error",
			Hint:  "Run 'mlr put --help' for DSL syntax reference.",
		}
	}

	return StructuredError{
		Error: msg,
		Kind:  "generic",
	}
}

// nearMatches populates did_you_mean for the structured error.
func nearMatches(se StructuredError) []string {
	if se.Token == "" {
		return nil
	}
	query := strings.ToLower(se.Token)

	var candidates []string
	switch se.Kind {
	case "unknown-verb":
		candidates = transformers.GetVerbNames()

	case "unknown-flag":
		candidates = cli.FLAG_TABLE.GetFlagNames()

	case "verb-option-error":
		// Use structured OptionSpec when available (from PR3 catalog).
		if se.Verb != "" {
			if info := transformers.GetVerbInfoForJSON(se.Verb); info != nil {
				for _, opt := range info.Options {
					candidates = append(candidates, opt.Flag)
				}
			}
		}
		if len(candidates) == 0 {
			candidates = cli.FLAG_TABLE.GetFlagNames()
		}

	case "dsl-parse-error":
		candidates = append(cst.BuiltinFunctionManagerInstance.GetBuiltinFunctionNames(),
			cst.GetKeywordNames()...)
	}

	return topMatches(query, candidates, 3, levenshteinThreshold(query))
}

// topMatches returns up to n candidates from the list sorted by edit distance,
// keeping only those within maxDist of the query.
func topMatches(query string, candidates []string, n, maxDist int) []string {
	type scored struct {
		name string
		dist int
	}
	var scored_ []scored
	for _, c := range candidates {
		d := levenshtein(query, strings.ToLower(c))
		if d <= maxDist {
			scored_ = append(scored_, scored{c, d})
		}
	}
	sort.Slice(scored_, func(i, j int) bool {
		if scored_[i].dist != scored_[j].dist {
			return scored_[i].dist < scored_[j].dist
		}
		return scored_[i].name < scored_[j].name
	})
	result := make([]string, 0, n)
	for i, s := range scored_ {
		if i >= n {
			break
		}
		result = append(result, s.name)
	}
	return result
}

// levenshteinThreshold returns the maximum edit distance to consider a match.
// Shorter queries use a tighter threshold.
func levenshteinThreshold(query string) int {
	n := len(query)
	switch {
	case n <= 3:
		return 1
	case n <= 6:
		return 2
	default:
		return 3
	}
}

// ----------------------------------------------------------------
// Levenshtein edit distance (Wagner-Fischer, O(m*n) time, O(n) space)

func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := range prev {
		prev[j] = j
	}
	for i, ca := range ra {
		curr[0] = i + 1
		for j, cb := range rb {
			cost := 1
			if ca == cb {
				cost = 0
			}
			curr[j+1] = min3(curr[j]+1, prev[j+1]+1, prev[j]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// ----------------------------------------------------------------
// Helpers for string-based categorization

// extractTokenFromVerbError extracts the quoted token from a verb error
// such as `mlr cut: option "--foo" not recognized`.
// Used in mlrcli_parse.go to populate CLIError.Token from the error string.
func extractTokenFromVerbError(msg string) string {
	return extractQuotedToken(msg)
}

// extractVerbFromMsg extracts the verb from "mlr {verb}: ..." error messages.
func extractVerbFromMsg(msg string) string {
	// Pattern: "mlr {verb}: ..."
	msg = strings.TrimPrefix(msg, "mlr ")
	if i := strings.Index(msg, ":"); i > 0 {
		return msg[:i]
	}
	return ""
}

// extractQuotedToken extracts the first double-quoted token from a message.
func extractQuotedToken(msg string) string {
	start := strings.Index(msg, "\"")
	if start < 0 {
		return ""
	}
	end := strings.Index(msg[start+1:], "\"")
	if end < 0 {
		return ""
	}
	return msg[start+1 : start+1+end]
}
