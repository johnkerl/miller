// Renders the "Options:" block of a verb's usage message from its structured
// []OptionSpec, so that the prose usage text and the JSON catalog (Tier-2)
// stay in sync -- each verb's option list is written once, in its OptionSpec.

package transformers

import (
	"fmt"
	"os"
	"strings"
)

// verbUsageLineWidth is the wrap width for generated option descriptions.
const verbUsageLineWidth = 80

const helpFlagHead = "-h|--help"

// WriteVerbOptions writes the "Options:" block for a verb's usage message
// from its structured option list. Flag heads are aligned into a single
// column; descriptions are word-wrapped with continuation lines indented to
// the description column. Since every verb supports -h|--help, that line is
// appended uniformly.
func WriteVerbOptions(o *os.File, options []OptionSpec) {
	fmt.Fprintf(o, "Options:\n")

	heads := make([]string, len(options))
	maxHeadLen := len(helpFlagHead)
	for i := range options {
		opt := &options[i]
		head := opt.Flag
		if len(opt.Aliases) > 0 {
			head = opt.Flag + "|" + strings.Join(opt.Aliases, "|")
		}
		if opt.Arg != "" {
			head += " " + opt.Arg
		}
		heads[i] = head
		if len(head) > maxHeadLen {
			maxHeadLen = len(head)
		}
	}

	descColumn := maxHeadLen + 1
	for i := range options {
		writeVerbOptionLine(o, heads[i], options[i].Desc, descColumn)
	}
	writeVerbOptionLine(o, helpFlagHead, "Show this message.", descColumn)
}

// writeVerbOptionLine writes one option's head and word-wrapped description,
// with continuation lines indented to descColumn.
func writeVerbOptionLine(o *os.File, head string, desc string, descColumn int) {
	if desc == "" {
		fmt.Fprintf(o, "%s\n", head)
		return
	}
	indent := strings.Repeat(" ", descColumn)
	line := fmt.Sprintf("%-*s", descColumn, head)
	lineHasDesc := false
	for word := range strings.FieldsSeq(desc) {
		if lineHasDesc && len(line)+1+len(word) > verbUsageLineWidth {
			fmt.Fprintf(o, "%s\n", line)
			line = indent
			lineHasDesc = false
		}
		if lineHasDesc {
			line += " "
		}
		line += word
		lineHasDesc = true
	}
	fmt.Fprintf(o, "%s\n", line)
}
