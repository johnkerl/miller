// Package completion implements shell tab-completion for Miller.
//
// The hard part of completing Miller's command line is that, unlike most CLIs
// which have a single flag set, a Miller command line is a sequence of
// contexts:
//
//	mlr [main flags] verb1 [verb1 flags] then verb2 [verb2 flags] ... [files]
//
// So completing the word under the cursor requires knowing *which* context the
// cursor is in. We determine that with a tolerant left-to-right walk of the
// words before the cursor (Layer A), then generate candidates appropriate to
// that context (Layer B).
//
// The walk needs to know, for each flag, whether it consumes a following
// argument value -- otherwise it can't tell a flag's value apart from a verb
// name, a `then`, or a filename. Main-flag arity comes exactly from the
// flag-table (see cli.FlagTable.FlagTakesArg). Per-verb flag names and arity
// are scraped from each verb's usage text (see verbflags.go), with a small
// override table for the few verbs whose usage text doesn't follow the
// `-f {arg}` convention.
//
// We deliberately do NOT drive the real verb parsers to do the walk: several of
// them call os.Exit on incomplete or unrecognized input (e.g. `mlr subs -f x`
// with no replacement text yet), which would kill the completion subprocess on
// very common mid-typing command lines. The tolerant walk never exits and
// degrades gracefully instead.
package completion

import (
	"sort"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/terminals/help"
	"github.com/johnkerl/miller/v6/pkg/terminals/registry"
	"github.com/johnkerl/miller/v6/pkg/transformers"
)

// Directive tells the shell shim how to use the candidate list.
type Directive string

const (
	// DirectiveCandidates: offer exactly the candidate words, nothing else.
	DirectiveCandidates Directive = "candidates"
	// DirectiveFiles: ignore the candidate words; do filename completion.
	DirectiveFiles Directive = "files"
	// DirectiveDefault: offer the candidate words AND do filename completion.
	DirectiveDefault Directive = "default"
)

// Result is what Complete returns: a directive plus a prefix-filtered list of
// candidate words.
type Result struct {
	Directive  Directive
	Candidates []string
}

// contextKind classifies the position of the cursor within the command line.
type contextKind int

const (
	// ctxMainOrVerb: in the main-flags region -- either before the first verb,
	// or after a `--` separator. A main flag or a (first) verb may come next.
	ctxMainOrVerb contextKind = iota
	// ctxExpectVerb: immediately after `then`/`+`; a verb name comes next.
	ctxExpectVerb
	// ctxVerbFlags: inside a verb's flag region. A verb flag, `then`, or a
	// filename may come next.
	ctxVerbFlags
	// ctxFlagValue: the cursor word is the argument value for the immediately
	// preceding arg-taking flag.
	ctxFlagValue
	// ctxFiles: in the trailing data-file-names region.
	ctxFiles
	// ctxTerminalArgs: after a terminal subcommand (mlr help, mlr completion,
	// ...), which consumes the rest of the command line.
	ctxTerminalArgs
)

type context struct {
	kind         contextKind
	verb         string   // set when kind == ctxVerbFlags
	flag         string   // set when kind == ctxFlagValue: the arg-taking flag being valued
	sawVerb      bool     // for kind == ctxMainOrVerb: whether a verb has appeared yet
	terminal     string   // set when kind == ctxTerminalArgs: the terminal name
	terminalArgs []string // for kind == ctxTerminalArgs: words after the terminal, before the cursor
}

// Complete is the entry point for the engine. words is the full argv as the
// shell sees it (words[0] is the program name, e.g. "mlr"); cword is the
// zero-based index of the word the cursor is on (matching bash's COMP_CWORD).
func Complete(words []string, cword int) Result {
	cur := ""
	if cword >= 0 && cword < len(words) {
		cur = words[cword]
	}

	// Walk the words strictly before the cursor to classify the cursor's
	// context.
	end := min(max(cword, 0), len(words))
	ctx := walk(words, end)

	switch ctx.kind {

	case ctxMainOrVerb:
		if strings.HasPrefix(cur, "-") {
			// Main flags plus the top-level terminal flags (-h, --help,
			// --version, the help shorthands, ...).
			cands := sortedUnion(mainFlagNames(), terminalFlagNames())
			return Result{DirectiveCandidates, filterByPrefix(cands, cur)}
		}
		// Verb names, plus terminal subcommand names (help, version, ...) when
		// no verb has appeared yet -- terminals are valid only as the first
		// non-flag token.
		if !ctx.sawVerb {
			cands := sortedUnion(verbNames(), terminalNames())
			return Result{DirectiveCandidates, filterByPrefix(cands, cur)}
		}
		return Result{DirectiveCandidates, filterByPrefix(verbNames(), cur)}

	case ctxExpectVerb:
		return Result{DirectiveCandidates, filterByPrefix(verbNames(), cur)}

	case ctxVerbFlags:
		if strings.HasPrefix(cur, "-") {
			return Result{DirectiveCandidates, filterByPrefix(verbFlagNames(ctx.verb), cur)}
		}
		// A non-flag here is either the verb-chain keyword `then` or the
		// beginning of filenames.
		return Result{DirectiveDefault, filterByPrefix([]string{"then"}, cur)}

	case ctxFlagValue:
		// For a main flag whose argument is a known enumerated set (file
		// formats, separator aliases), offer those values. Verb-flag values
		// are not yet completed (issue #2097) -- and we must not apply the
		// main-flag value sets to an identically-spelled verb flag (e.g.
		// `mlr uniq -o` takes a field name, not a format).
		if ctx.verb == "" {
			if values := cli.FlagValueCandidates(ctx.flag); values != nil {
				return Result{DirectiveCandidates, filterByPrefix(values, cur)}
			}
		}
		return Result{DirectiveFiles, nil}

	case ctxFiles:
		return Result{DirectiveFiles, nil}

	case ctxTerminalArgs:
		return completeTerminalArgs(ctx.terminal, ctx.terminalArgs, cur)
	}

	return Result{DirectiveFiles, nil}
}

// walk scans words[1:end] left to right and reports the context that the word
// at index `end` (the cursor word) sits in. It mirrors the segmentation done by
// climain's pass-one parser, but never exits and tolerates incomplete input.
func walk(words []string, end int) context {
	i := 1
	inVerb := false
	curVerb := ""
	// sawVerb records whether a verb has been seen yet. Terminal subcommands
	// (mlr help, mlr version, ...) are valid only as the first non-flag token,
	// i.e. before any verb, so they should be offered only while !sawVerb.
	sawVerb := false

	for i < end {
		tok := words[i]

		if !inVerb {
			// Main-flags region (also the slot for the first verb).
			if strings.HasPrefix(tok, "-") {
				if tok == "--" {
					// Separator between a verb and a following main flag.
					i++
					continue
				}
				found, takesArg := cli.FLAG_TABLE.FlagTakesArg(tok)
				if found && takesArg {
					if i+1 >= end {
						// The value is the cursor word.
						return context{kind: ctxFlagValue, flag: tok}
					}
					i += 2
					continue
				}
				// Arity-0 or unrecognized main flag: consume just the flag.
				i++
				continue
			}
			if tok == "then" || tok == "+" {
				i++
				if i >= end {
					return context{kind: ctxExpectVerb}
				}
				curVerb = words[i]
				inVerb = true
				sawVerb = true
				i++
				continue
			}
			// First non-flag token. A terminal subcommand (mlr help, mlr
			// version, ...) is valid only here; everything after it belongs to
			// that terminal.
			if !sawVerb && isTerminalName(tok) {
				return context{kind: ctxTerminalArgs, terminal: tok, terminalArgs: words[i+1 : end]}
			}
			// First verb.
			curVerb = tok
			inVerb = true
			sawVerb = true
			i++
			continue
		}

		// Inside verb curVerb's flag region.
		if tok == "then" || tok == "+" {
			i++
			if i >= end {
				return context{kind: ctxExpectVerb}
			}
			curVerb = words[i]
			i++
			continue
		}
		if tok == "--" {
			// Back to the main-flags region; main flags may follow a verb.
			inVerb = false
			i++
			continue
		}
		if strings.HasPrefix(tok, "-") {
			found, takesArg := verbFlagTakesArg(curVerb, tok)
			if found && takesArg {
				if i+1 >= end {
					return context{kind: ctxFlagValue, verb: curVerb, flag: tok}
				}
				i += 2
				continue
			}
			i++
			continue
		}
		// A non-flag, non-keyword token inside a verb region begins the
		// trailing data-file names (Miller puts filenames last).
		return context{kind: ctxFiles}
	}

	// Reached the cursor word.
	if inVerb {
		return context{kind: ctxVerbFlags, verb: curVerb}
	}
	return context{kind: ctxMainOrVerb, sawVerb: sawVerb}
}

// filterByPrefix returns the candidates that have cur as a prefix, preserving
// input order.
func filterByPrefix(candidates []string, cur string) []string {
	if cur == "" {
		return candidates
	}
	out := make([]string, 0, len(candidates))
	for _, c := range candidates {
		if strings.HasPrefix(c, cur) {
			out = append(out, c)
		}
	}
	return out
}

// mainFlagNames returns all main-flag spellings, sorted for a navigable
// display. This includes the format-conversion keystroke-saver flags (--c2j,
// --x2y, ...).
func mainFlagNames() []string {
	names := cli.FLAG_TABLE.GetFlagNames()
	sort.Strings(names)
	return names
}

// verbNames returns all verb names, sorted.
func verbNames() []string {
	names := transformers.GetVerbNames()
	sort.Strings(names)
	return names
}

// terminalNames returns the terminal subcommand names (help, version, ...).
func terminalNames() []string {
	return registry.Names
}

// terminalFlagNames returns the top-level terminal flags: -h/--help and the
// help shorthands, plus the version flags.
func terminalFlagNames() []string {
	return append(help.GetTerminalFlagNames(), registry.VersionFlagNames...)
}

// sortedUnion merges the given name lists, de-duplicates, and sorts the result.
func sortedUnion(lists ...[]string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0)
	for _, list := range lists {
		for _, name := range list {
			if !seen[name] {
				seen[name] = true
				out = append(out, name)
			}
		}
	}
	sort.Strings(out)
	return out
}
