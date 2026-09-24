// Support for issue #233: if the input format is numerically indexed --
// records read with field names 1,2,3,... via --nidx / -T / --inidx /
// --implicit-csv-header and friends -- and the first verb in the then-chain
// asks for fields by names none of which are numeric, then (almost certainly)
// nothing will match: e.g. 'mlr -T cut -f last_first' when the user intended
// 'mlr --tsv cut -f last_first'. Here we print a warning to stderr -- not an
// error, since output and exit code must remain unchanged: there are
// legitimate (if unusual) uses, e.g. 'cut -x -f name' as an intentional no-op
// guard, and we cannot be fully sure of the user's intent.
//
// To keep false positives to a minimum, we warn only when all of the
// following hold:
//
// * Miller is actually reading input (not 'mlr -n', not a seqgen-led chain).
// * The input format gives positional field names 1,2,3,...
// * The verb is the *first* in the then-chain: verbs downstream of label,
//   put, rename, etc. can legitimately see non-numeric field names even with
//   --nidx input.
// * The verb takes field names (not regexes) via its flags: 'cut -r' and the
//   having-fields matching flags are excluded since a regex like "1$" can
//   match positional names.
// * Every requested field name is non-numeric. A single numeric name like
//   '-f 1,foo' suggests the user knows the input is positionally named.

package climain

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
)

// verbFlagsForPositionalInputHint maps verb names to the verb flags which
// take comma-separated field names (to be looked up in the input records) as
// their argument.
var verbFlagsForPositionalInputHint = map[string]map[string]bool{
	"cut":           {"-f": true},
	"having-fields": {"--at-least": true, "--which-are": true, "--at-most": true},
	"nest":          {"-f": true},
}

// verbRegexFlagsForPositionalInputHint maps verb names to the verb flags
// whose presence means field names are treated as regular expressions, in
// which case we make no attempt to guess what they might match.
var verbRegexFlagsForPositionalInputHint = map[string]map[string]bool{
	"cut":           {"-r": true},
	"having-fields": {"--all-matching": true, "--any-matching": true, "--none-matching": true},
}

// maybeWarnOnNamedFieldsForPositionalInput prints a warning (only) to stderr
// as described at the top of this file. It never modifies options or
// transformers, and never exits.
func maybeWarnOnNamedFieldsForPositionalInput(
	options *cli.TOptions,
	verbSequences [][]string,
) {
	if options.NoInput {
		return
	}

	// Positionally-named inputs: NIDX format, or header-having formats with
	// --implicit-csv-header et al.
	if options.ReaderOptions.InputFileFormat != "nidx" && !options.ReaderOptions.UseImplicitHeader {
		return
	}

	if len(verbSequences) == 0 {
		return
	}
	// Only the first verb in the then-chain sees the reader's positional
	// field names directly.
	verbSequence := verbSequences[0]
	verb := verbSequence[0]

	nameFlags := verbFlagsForPositionalInputHint[verb]
	if nameFlags == nil {
		return
	}
	regexFlags := verbRegexFlagsForPositionalInputHint[verb]

	fieldNames := []string{}
	for i := 1; i < len(verbSequence); i++ {
		if regexFlags[verbSequence[i]] {
			return
		}
		if nameFlags[verbSequence[i]] && i+1 < len(verbSequence) {
			fieldNames = append(fieldNames, strings.Split(verbSequence[i+1], ",")...)
			i++
		}
	}
	if len(fieldNames) == 0 {
		return
	}

	for _, fieldName := range fieldNames {
		if _, err := strconv.Atoi(fieldName); err == nil {
			// At least one field name is a positional index: assume the user
			// knows what they're doing.
			return
		}
	}

	fmt.Fprintf(os.Stderr,
		"mlr: warning: %s was given field name(s) \"%s\", but the input format is numerically indexed with positional field names 1,2,3,... -- so no fields will match.\n",
		verb, strings.Join(fieldNames, ","),
	)
	fmt.Fprintf(os.Stderr,
		"mlr: warning: if your data has a header line, use a header-aware format flag such as --csv or --tsv in place of --nidx/-T/--implicit-csv-header; or, refer to fields positionally, e.g. -f 1,2.\n",
	)
}
