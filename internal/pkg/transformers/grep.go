package transformers

import (
	"container/list"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameGrep = "grep"

var GrepSetup = TransformerSetup{
	Verb:         verbNameGrep,
	UsageFunc:    transformerGrepUsage,
	ParseCLIFunc: transformerGrepParseCLI,
	IgnoresInput: false,
}

func transformerGrepUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {regular expression}\n", "mlr", verbNameGrep)
	fmt.Fprintf(o, "Passes through records which match the regular expression.\n")

	fmt.Fprint(o, "Options:\n")
	fmt.Fprint(o, "-i  Use case-insensitive search.\n")
	fmt.Fprint(o, "-v  Invert: pass through records which do not match the regex.\n")
	fmt.Fprint(o, "-a  Only grep for values, not keys and values.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	fmt.Fprintf(o, `Note that "%s filter" is more powerful, but requires you to know field names.
By contrast, "%s grep" allows you to regex-match the entire record. It does this
by formatting each record in memory as DKVP (or NIDX, if -a is supplied), using
OFS "," and OPS "=", and matching the resulting line against the regex specified
here. In particular, the regex is not applied to the input stream: if you have
CSV with header line "x,y,z" and data line "1,2,3" then the regex will be
matched, not against either of these lines, but against the DKVP line
"x=1,y=2,z=3".  Furthermore, not all the options to system grep are supported,
and this command is intended to be merely a keystroke-saver. To get all the
features of system grep, you can do
  "%s --odkvp ... | grep ... | %s --idkvp ..."
`, "mlr", "mlr", "mlr", "mlr")
}

func transformerGrepParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	ignoreCase := false
	invert := false
	valuesOnly := false

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		if args[argi] == "--" {
			break // All transformers must do this so main-flags can follow verb-flags
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerGrepUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-i" {
			ignoreCase = true

		} else if opt == "-v" {
			invert = true

		} else if opt == "-a" {
			valuesOnly = true

		} else {
			transformerGrepUsage(os.Stderr)
			os.Exit(1)
		}
	}

	// Get the regex from the command line
	if argi >= argc {
		transformerGrepUsage(os.Stderr)
		os.Exit(1)
	}
	pattern := args[argi]
	argi++

	if ignoreCase {
		pattern = "(?i)" + pattern
	}

	// TODO: maybe CompilePOSIX
	regexp, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %s: couldn't compile regex \"%s\"\n",
			"mlr", verb, pattern)
		os.Exit(1)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerGrep(
		regexp,
		invert,
		valuesOnly,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerGrep struct {
	regexp     *regexp.Regexp
	invert     bool
	valuesOnly bool
}

func NewTransformerGrep(
	regexp *regexp.Regexp,
	invert bool,
	valuesOnly bool,
) (*TransformerGrep, error) {
	tr := &TransformerGrep{
		regexp:     regexp,
		invert:     invert,
		valuesOnly: valuesOnly,
	}
	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerGrep) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		var inrecAsString string
		if tr.valuesOnly {
			inrecAsString = inrec.ToNIDXString()
		} else {
			inrecAsString = inrec.ToDKVPString()
		}
		matches := tr.regexp.Match([]byte(inrecAsString))
		if tr.invert {
			if !matches {
				outputRecordsAndContexts.PushBack(inrecAndContext)
			}
		} else {
			if matches {
				outputRecordsAndContexts.PushBack(inrecAndContext)
			}
		}
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}
