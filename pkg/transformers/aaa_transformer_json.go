// Machine-readable (JSON) accessors over the verb (transformer) catalog, for
// `mlr help --as-json` and similar tooling.
//
// Tier-1: every verb exposes Summary and UsageText (prose fallback).
// Tier-2: verbs that have been migrated populate Options on their
// TransformerSetup; those entries appear as a structured option list in the
// JSON and agents no longer need to scrape the prose.

package transformers

import (
	"bytes"
	"io"
	"os"
	"strings"
)

// VerbInfoForJSON is the structured view of a single verb.
//   - Options is non-nil for Tier-2 verbs that have been migrated; agents
//     should prefer it when available.
//   - UsageText is always present as the Tier-1 prose fallback.
type VerbInfoForJSON struct {
	Name         string       `json:"name"`
	Summary      string       `json:"summary"`
	IgnoresInput bool         `json:"ignores_input"`
	Options      []OptionSpec `json:"options"`
	UsageText    string       `json:"usage_text"`
}

// captureUsageFunc runs a verb's UsageFunc against a pipe and returns what it
// printed. The UsageFunc signature takes an *os.File, so we hand it the
// write-end of a pipe directly -- no global-stdout swap needed. A goroutine
// drains the read-end so we never block on the OS pipe buffer.
func captureUsageFunc(usageFunc TransformerUsageFunc) string {
	r, w, err := os.Pipe()
	if err != nil {
		return ""
	}

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		done <- buf.String()
	}()

	usageFunc(w)
	w.Close()
	s := <-done
	r.Close()
	return s
}

// summarizeUsageText returns the first line that isn't blank or a "Usage:"
// banner -- the closest thing each verb has to a one-line description.
func summarizeUsageText(usageText string) string {
	for _, line := range strings.Split(usageText, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "Usage:") {
			continue
		}
		return trimmed
	}
	return ""
}

func makeVerbInfoForJSON(setup *TransformerSetup) *VerbInfoForJSON {
	usageText := captureUsageFunc(setup.UsageFunc)
	return &VerbInfoForJSON{
		Name:         setup.Verb,
		Summary:      summarizeUsageText(usageText),
		IgnoresInput: setup.IgnoresInput,
		Options:      setup.Options, // nil for unmigrated verbs; omitted from JSON via omitempty
		UsageText:    strings.TrimRight(usageText, "\n"),
	}
}

// GetVerbInfosForJSON returns the full verb catalog in table order.
func GetVerbInfosForJSON() []*VerbInfoForJSON {
	infos := make([]*VerbInfoForJSON, 0, len(TRANSFORMER_LOOKUP_TABLE))
	for i := range TRANSFORMER_LOOKUP_TABLE {
		infos = append(infos, makeVerbInfoForJSON(&TRANSFORMER_LOOKUP_TABLE[i]))
	}
	return infos
}

// GetVerbInfoForJSON returns the structured view of a single verb, or nil if
// there is no such verb.
func GetVerbInfoForJSON(verb string) *VerbInfoForJSON {
	setup := LookUp(verb)
	if setup == nil {
		return nil
	}
	return makeVerbInfoForJSON(setup)
}
