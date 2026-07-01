// Machine-readable (JSON) accessors over the verb (transformer) catalog, for
// `mlr help --as-json` and similar tooling.
//
// Tier-1 caveat: unlike functions and flags, verb options are not held in any
// structured form -- each verb hand-writes a UsageFunc that prints prose. So
// here we expose the verb name, a one-line summary, and the captured raw usage
// text. Structured per-verb options (flag/arg/type) are a planned follow-on
// (an optional Options field on TransformerSetup); when present they can be
// emitted alongside UsageText.

package transformers

import (
	"bytes"
	"io"
	"os"
	"strings"
)

// VerbInfoForJSON is the structured view of a single verb. Summary is the first
// non-"Usage:" line of the usage text; UsageText is the verb's full usage
// output verbatim (the Tier-1 fallback for not-yet-structured options).
type VerbInfoForJSON struct {
	Name         string `json:"name"`
	Summary      string `json:"summary"`
	IgnoresInput bool   `json:"ignores_input"`
	UsageText    string `json:"usage_text"`
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
