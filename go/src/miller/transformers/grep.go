package transformers

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"miller/clitypes"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameGrep = "grep"

var GrepSetup = transforming.TransformerSetup{
	Verb:         verbNameGrep,
	ParseCLIFunc: transformerGrepParseCLI,
	UsageFunc:    transformerGrepUsage,
	IgnoresInput: false,
}

func transformerGrepParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	ignoreCase := false
	invert := false

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerGrepUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else if args[argi] == "-i" {
			ignoreCase = true
			argi++

		} else if args[argi] == "-v" {
			invert = true
			argi++

		} else {
			transformerGrepUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	// Get the regex from the command line
	if argi >= argc {
		transformerGrepUsage(os.Stderr, true, 1)
	}
	pattern := args[argi]
	argi += 1

	if ignoreCase {
		pattern = "(?i)" + pattern
	}

	// TODO: maybe CompilePOSIX
	regexp, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %s: couldn't compile regex \"%s\"\n",
			args[0], verb, pattern)
		os.Exit(1)
	}

	transformer, _ := NewTransformerGrep(
		regexp,
		invert,
	)

	*pargi = argi
	return transformer
}

func transformerGrepUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {regular expression}\n", os.Args[0], verbNameGrep)
	fmt.Fprintf(o, "Passes through records which match the regular expression.\n")

	fmt.Fprint(o, "Options:\n")
	fmt.Fprint(o, "-i  Use case-insensitive search.\n")
	fmt.Fprint(o, "-v  Invert: pass through records which do not match the regex.\n")

	fmt.Fprintf(o, `Note that "%s filter" is more powerful, but requires you to know field names.
By contrast, "%s grep" allows you to regex-match the entire record. It does
this by formatting each record in memory as DKVP, using command-line-specified
ORS/OFS/OPS, and matching the resulting line against the regex specified
here. In particular, the regex is not applied to the input stream: if you
have CSV with header line "x,y,z" and data line "1,2,3" then the regex will
be matched, not against either of these lines, but against the DKVP line
"x=1,y=2,z=3".  Furthermore, not all the options to system grep are supported,
and this command is intended to be merely a keystroke-saver. To get all the
features of system grep, you can do
  "%s --odkvp ... | grep ... | %s --idkvp ..."
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0])

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerGrep struct {
	regexp *regexp.Regexp
	invert bool
}

func NewTransformerGrep(
	regexp *regexp.Regexp,
	invert bool,
) (*TransformerGrep, error) {
	this := &TransformerGrep{
		regexp: regexp,
		invert: invert,
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerGrep) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		inrecAsString := inrec.ToDKVPString()
		matches := this.regexp.Match([]byte(inrecAsString))
		if this.invert {
			if !matches {
				outputChannel <- inrecAndContext
			}
		} else {
			if matches {
				outputChannel <- inrecAndContext
			}
		}
	} else {
		outputChannel <- inrecAndContext
	}
}
