// Per-verb flag discovery for shell completion.
//
// Miller's verbs do not expose a structured flag table the way main flags do;
// each verb hand-rolls its own CLI parser. Rather than maintain a parallel
// hand-written table that could drift from the parsers, we scrape each verb's
// usage text -- which is the same source of truth shown to users by
// `mlr <verb> --help` -- for its flag names and arity.
//
// The scrape relies on Miller's consistent usage-text convention:
//
//	-f {comma-separated field names}  ...   (takes an argument: brace follows)
//	-n                                ...   (no argument)
//	-x|--complement                   ...   (alternate spellings, '|'-separated)
//
// A handful of verbs document a flag without the `{...}` convention (e.g.
// `put -s name=value`); those are corrected via verbFlagArityOverrides. Getting
// arity slightly wrong only affects the tolerant command-line walk in rare
// adversarial cases (a flag value that looks like `then` or a verb name); it
// never crashes and never produces a wrong command line.

package completion

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/transformers"
)

// verbFlagInfo is the scraped flag metadata for one verb.
type verbFlagInfo struct {
	names    []string        // flag spellings, in usage-text order
	takesArg map[string]bool // spelling -> whether it consumes an argument
}

// verbFlagArityOverrides corrects arity for flags whose usage text does not use
// the `-flag {arg}` convention. Keyed by verb, then by flag spelling, with the
// value being whether the flag takes an argument.
var verbFlagArityOverrides = map[string]map[string]bool{
	// `put`/`filter` document `-s name=value` and `-e {expression}` etc.; the
	// `-s` form has no braces, so the scraper would under-detect its arity.
	"put":    {"-s": true},
	"filter": {"-s": true},
}

// verbFlagCache memoizes scrapes within a single process invocation.
var verbFlagCache = map[string]*verbFlagInfo{}

// verbFlagNames returns the flag spellings to offer as completion candidates
// for the given verb.
func verbFlagNames(verb string) []string {
	return getVerbFlagInfo(verb).names
}

// verbFlagTakesArg reports whether `flag` is a known flag of `verb` and whether
// it consumes a following argument value.
func verbFlagTakesArg(verb string, flag string) (found bool, takesArg bool) {
	info := getVerbFlagInfo(verb)
	ta, ok := info.takesArg[flag]
	return ok, ta
}

func getVerbFlagInfo(verb string) *verbFlagInfo {
	if info, ok := verbFlagCache[verb]; ok {
		return info
	}
	info := scrapeVerbFlagInfo(verb)
	verbFlagCache[verb] = info
	return info
}

func scrapeVerbFlagInfo(verb string) *verbFlagInfo {
	info := &verbFlagInfo{
		names:    []string{},
		takesArg: map[string]bool{},
	}

	usage, ok := captureVerbUsage(verb)
	if ok {
		parseUsageFlags(usage, info)
	}

	// Apply arity overrides, and ensure overridden flags are offered as
	// candidates even if scraping missed them.
	if overrides, ok := verbFlagArityOverrides[verb]; ok {
		for flag, takesArg := range overrides {
			if _, seen := info.takesArg[flag]; !seen {
				info.names = append(info.names, flag)
			}
			info.takesArg[flag] = takesArg
		}
	}

	return info
}

// parseUsageFlags extracts flag spellings and arity from verb usage text.
func parseUsageFlags(usage string, info *verbFlagInfo) {
	for _, line := range strings.Split(usage, "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		if !strings.HasPrefix(trimmed, "-") {
			continue
		}
		fields := strings.Fields(trimmed)
		if len(fields) == 0 {
			continue
		}
		// A following `{...}` token signals an argument-taking flag.
		takesArg := len(fields) >= 2 && strings.HasPrefix(fields[1], "{")

		// The leading token may bundle alternate spellings, e.g. `-h|--help`
		// or `-x|--complement` or `-tr|-rt`.
		for _, name := range strings.Split(fields[0], "|") {
			name = strings.TrimRight(name, ":,")
			if name == "-" || name == "--" || !strings.HasPrefix(name, "-") {
				continue
			}
			if _, seen := info.takesArg[name]; !seen {
				info.names = append(info.names, name)
			}
			// If any spelling on the line takes an argument, all do.
			if takesArg || info.takesArg[name] {
				info.takesArg[name] = true
			} else {
				info.takesArg[name] = false
			}
		}
	}
}

// captureVerbUsage runs a verb's UsageFunc and returns its text. UsageFunc
// writes to an *os.File, so we capture it through an os.Pipe.
func captureVerbUsage(verb string) (string, bool) {
	setup := transformers.LookUp(verb)
	if setup == nil || setup.UsageFunc == nil {
		return "", false
	}

	r, w, err := os.Pipe()
	if err != nil {
		return "", false
	}

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	setup.UsageFunc(w)
	_ = w.Close()
	text := <-done
	_ = r.Close()

	return text, true
}
