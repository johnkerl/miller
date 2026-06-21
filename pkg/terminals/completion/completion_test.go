package completion

import (
	"slices"
	"testing"
)

// split builds a words slice and cword from a command line plus the
// already-typed current word. The cursor word is always the last element.
func wordsAndCword(words ...string) ([]string, int) {
	return words, len(words) - 1
}

func TestContextDirectives(t *testing.T) {
	tests := []struct {
		name      string
		words     []string
		wantDir   Directive
		mustHave  []string // candidates that must be present
		mustLack  []string // candidates that must be absent
		emptyCand bool     // candidate list must be empty
	}{
		{
			name:     "bare mlr offers verbs",
			words:    []string{"mlr", ""},
			wantDir:  DirectiveCandidates,
			mustHave: []string{"cat", "sort", "put"},
			mustLack: []string{"--icsv"},
		},
		{
			name:     "dash in main region offers main flags",
			words:    []string{"mlr", "--ic"},
			wantDir:  DirectiveCandidates,
			mustHave: []string{"--icsv"},
			mustLack: []string{"cat"},
		},
		{
			name:     "after format flag offers verbs",
			words:    []string{"mlr", "--icsv", ""},
			wantDir:  DirectiveCandidates,
			mustHave: []string{"cat", "head"},
		},
		{
			name:      "main flag taking a filename arg yields file completion",
			words:     []string{"mlr", "--from", ""},
			wantDir:   DirectiveFiles,
			emptyCand: true,
		},
		{
			name:     "inside verb, dash offers that verb's flags",
			words:    []string{"mlr", "cat", "-"},
			wantDir:  DirectiveCandidates,
			mustHave: []string{"-n", "-N", "-g", "--filename"},
			mustLack: []string{"-f"}, // -f is not a cat flag
		},
		{
			name:     "inside verb, non-flag offers then plus files",
			words:    []string{"mlr", "cat", ""},
			wantDir:  DirectiveDefault,
			mustHave: []string{"then"},
		},
		{
			name:     "after then offers verbs",
			words:    []string{"mlr", "sort", "-f", "a,b", "then", ""},
			wantDir:  DirectiveCandidates,
			mustHave: []string{"cat", "tac", "head"},
		},
		{
			name:      "verb flag taking arg yields file completion",
			words:     []string{"mlr", "sort", "-f", ""},
			wantDir:   DirectiveFiles,
			emptyCand: true,
		},
		{
			name:      "trailing filename region",
			words:     []string{"mlr", "cat", "data.csv", ""},
			wantDir:   DirectiveFiles,
			emptyCand: true,
		},
		{
			name:     "double-dash separator returns to main flags",
			words:    []string{"mlr", "cat", "--", "--oj"},
			wantDir:  DirectiveCandidates,
			mustHave: []string{"--ojson"},
		},
		{
			name:     "second verb in chain completes its own flags",
			words:    []string{"mlr", "cat", "-n", "then", "head", "-"},
			wantDir:  DirectiveCandidates,
			mustHave: []string{"-n", "-g"},
		},
		{
			name:     "prefix filtering inside verb flags",
			words:    []string{"mlr", "cut", "--comp"},
			wantDir:  DirectiveCandidates,
			mustHave: []string{"--complement"},
			mustLack: []string{"-f"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			words, cword := wordsAndCword(tc.words...)
			got := Complete(words, cword)
			if got.Directive != tc.wantDir {
				t.Errorf("directive: got %q, want %q (candidates=%v)",
					got.Directive, tc.wantDir, got.Candidates)
			}
			if tc.emptyCand && len(got.Candidates) != 0 {
				t.Errorf("expected no candidates, got %v", got.Candidates)
			}
			for _, want := range tc.mustHave {
				if !slices.Contains(got.Candidates, want) {
					t.Errorf("missing expected candidate %q in %v", want, got.Candidates)
				}
			}
			for _, lack := range tc.mustLack {
				if slices.Contains(got.Candidates, lack) {
					t.Errorf("unexpected candidate %q in %v", lack, got.Candidates)
				}
			}
		})
	}
}

// TestTerminalCompletion verifies that terminal subcommands (mlr help, mlr
// version, ...) and the top-level terminal flags (-h, --version, ...) are
// offered as completions -- and only where they are valid.
func TestTerminalCompletion(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		mustHave []string
		mustLack []string
	}{
		{
			name:     "first word offers terminals alongside verbs",
			words:    []string{"mlr", ""},
			mustHave: []string{"cat", "help", "version", "repl", "completion"},
		},
		{
			name:     "first word prefix matches both verb and terminal",
			words:    []string{"mlr", "he"},
			mustHave: []string{"head", "help"},
		},
		{
			name:     "leading dash offers terminal flags",
			words:    []string{"mlr", "-"},
			mustHave: []string{"-h", "--help", "--version", "--bare-version", "-L", "-F"},
		},
		{
			name:     "version flag prefix",
			words:    []string{"mlr", "--ve"},
			mustHave: []string{"--version"},
			mustLack: []string{"--bare-version"},
		},
		{
			name:     "terminals not offered after then",
			words:    []string{"mlr", "cat", "then", ""},
			mustHave: []string{"head", "sort"},
			mustLack: []string{"help", "version", "repl"},
		},
		{
			name:     "terminals not offered in main region after a verb",
			words:    []string{"mlr", "cat", "--", ""},
			mustLack: []string{"help", "version"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Complete(tc.words, len(tc.words)-1)
			for _, want := range tc.mustHave {
				if !slices.Contains(got.Candidates, want) {
					t.Errorf("missing expected candidate %q in %v", want, got.Candidates)
				}
			}
			for _, lack := range tc.mustLack {
				if slices.Contains(got.Candidates, lack) {
					t.Errorf("unexpected candidate %q in %v", lack, got.Candidates)
				}
			}
		})
	}
}

// TestHelpTopicCompletion verifies completion of `mlr help` topics and their
// name arguments, plus `mlr completion` subcommands.
func TestHelpTopicCompletion(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		mustHave []string
		mustLack []string
	}{
		{
			name:     "help topics",
			words:    []string{"mlr", "help", ""},
			mustHave: []string{"flags", "verb", "function", "keyword", "list-verbs"},
		},
		{
			name:     "help topic prefix",
			words:    []string{"mlr", "help", "list-v"},
			mustHave: []string{"list-verbs"},
			mustLack: []string{"flags"},
		},
		{
			name:     "help verb takes verb names",
			words:    []string{"mlr", "help", "verb", ""},
			mustHave: []string{"cat", "sort", "put"},
		},
		{
			name:     "help function takes function names",
			words:    []string{"mlr", "help", "function", "strl"},
			mustHave: []string{"strlen"},
		},
		{
			name:     "help keyword takes keyword names",
			words:    []string{"mlr", "help", "keyword", ""},
			mustHave: []string{"ENV", "FILENAME"},
		},
		{
			name:     "help flag takes flag names",
			words:    []string{"mlr", "help", "flag", "--ic"},
			mustHave: []string{"--icsv"},
		},
		{
			name:     "help works after a main flag",
			words:    []string{"mlr", "--icsv", "help", ""},
			mustHave: []string{"flags", "verb"},
		},
		{
			name:     "completion subcommands",
			words:    []string{"mlr", "completion", ""},
			mustHave: []string{"bash", "zsh"},
		},
		{
			name:     "help topic with no name-argument yields nothing",
			words:    []string{"mlr", "help", "list-verbs", ""},
			mustLack: []string{"cat", "flags"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Complete(tc.words, len(tc.words)-1)
			for _, want := range tc.mustHave {
				if !slices.Contains(got.Candidates, want) {
					t.Errorf("missing expected candidate %q in %v", want, got.Candidates)
				}
			}
			for _, lack := range tc.mustLack {
				if slices.Contains(got.Candidates, lack) {
					t.Errorf("unexpected candidate %q in %v", lack, got.Candidates)
				}
			}
		})
	}
}

// TestAdversarialFlagValue verifies that a verb flag's argument value which
// happens to look like the chain keyword `then` is treated as a value, not as a
// verb-chain separator. This requires correct per-verb arity.
func TestAdversarialFlagValue(t *testing.T) {
	// `mlr cut -f then -<cursor>`: `then` is the value of cut's -f, so the
	// cursor is still inside cut's flag region.
	words := []string{"mlr", "cut", "-f", "then", "-"}
	got := Complete(words, len(words)-1)
	if got.Directive != DirectiveCandidates {
		t.Fatalf("directive: got %q, want %q", got.Directive, DirectiveCandidates)
	}
	if !slices.Contains(got.Candidates, "-o") {
		t.Errorf("expected cut flag -o among candidates, got %v", got.Candidates)
	}
}

// TestPutArityOverride verifies the override table: put's `-s` takes an
// argument even though its usage text lacks the `{...}` convention.
func TestPutArityOverride(t *testing.T) {
	found, takesArg := verbFlagTakesArg("put", "-s")
	if !found || !takesArg {
		t.Errorf("put -s: got found=%v takesArg=%v, want true/true", found, takesArg)
	}
	// And it should be offered as a candidate.
	if !slices.Contains(verbFlagNames("put"), "-s") {
		t.Errorf("put -s missing from candidates %v", verbFlagNames("put"))
	}
}

// TestFlagValueCompletion verifies enum-value completion for arg-taking main
// flags: file formats for -i/-o/--io and separator aliases for --ifs etc.,
// with a fallback to filename completion for non-enum flags.
func TestFlagValueCompletion(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		wantDir  Directive
		mustHave []string
		mustLack []string
	}{
		{
			name:     "format flag -i offers file formats",
			words:    []string{"mlr", "-i", ""},
			wantDir:  DirectiveCandidates,
			mustHave: []string{"csv", "json", "tsv", "pprint"},
			mustLack: []string{"comma"},
		},
		{
			name:     "format flag --io offers file formats",
			words:    []string{"mlr", "--io", "js"},
			wantDir:  DirectiveCandidates,
			mustHave: []string{"json"},
			mustLack: []string{"csv"},
		},
		{
			name:     "separator flag --ifs offers aliases",
			words:    []string{"mlr", "--ifs", ""},
			wantDir:  DirectiveCandidates,
			mustHave: []string{"comma", "tab", "pipe", "semicolon"},
		},
		{
			name:     "separator flag --ofs prefix-filters aliases",
			words:    []string{"mlr", "--ofs", "co"},
			wantDir:  DirectiveCandidates,
			mustHave: []string{"colon", "comma"},
			mustLack: []string{"tab", "pipe"},
		},
		{
			name:     "regex separator flag offers regex aliases",
			words:    []string{"mlr", "--ifs-regex", ""},
			wantDir:  DirectiveCandidates,
			mustHave: []string{"spaces", "tabs", "whitespace"},
		},
		{
			name:    "non-enum flag --ofmt falls back to files",
			words:   []string{"mlr", "--ofmt", ""},
			wantDir: DirectiveFiles,
		},
		{
			// A verb flag spelled like a main value-flag must NOT inherit the
			// main flag's value set: `uniq -o` takes a field name, not a format.
			name:    "verb flag colliding with main format flag falls back to files",
			words:   []string{"mlr", "uniq", "-o", ""},
			wantDir: DirectiveFiles,
		},
		{
			name:    "filename flag --from falls back to files",
			words:   []string{"mlr", "--from", ""},
			wantDir: DirectiveFiles,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Complete(tc.words, len(tc.words)-1)
			if got.Directive != tc.wantDir {
				t.Errorf("directive: got %q, want %q (candidates=%v)",
					got.Directive, tc.wantDir, got.Candidates)
			}
			for _, want := range tc.mustHave {
				if !slices.Contains(got.Candidates, want) {
					t.Errorf("missing expected value %q in %v", want, got.Candidates)
				}
			}
			for _, lack := range tc.mustLack {
				if slices.Contains(got.Candidates, lack) {
					t.Errorf("unexpected value %q in %v", lack, got.Candidates)
				}
			}
		})
	}
}

// TestVerbFlagScrape sanity-checks the usage-text scraper on representative
// verbs.
func TestVerbFlagScrape(t *testing.T) {
	cases := []struct {
		verb     string
		flag     string
		takesArg bool
	}{
		{"head", "-n", true},
		{"head", "-g", true},
		{"cat", "-n", false},
		{"cat", "-N", true},
		{"cut", "-f", true},
		{"cut", "-o", false},
		{"cut", "--complement", false},
	}
	for _, c := range cases {
		found, takesArg := verbFlagTakesArg(c.verb, c.flag)
		if !found {
			t.Errorf("%s %s: not found", c.verb, c.flag)
			continue
		}
		if takesArg != c.takesArg {
			t.Errorf("%s %s: takesArg got %v want %v", c.verb, c.flag, takesArg, c.takesArg)
		}
	}
}
