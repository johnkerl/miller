package mappers

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"miller/clitypes"
	"miller/mapping"
	"miller/types"
)

// ----------------------------------------------------------------
var GrepSetup = mapping.MapperSetup{
	Verb:         "grep",
	ParseCLIFunc: mapperGrepParseCLI,
	IgnoresInput: false,
}

func mapperGrepParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) mapping.IRecordMapper {

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)

	pIgnoreCase := flagSet.Bool(
		"i",
		false,
		`Use case-insensitive search`,
	)

	pInvert := flagSet.Bool(
		"v",
		false,
		`Invert: pass through records which do not match the regex.`,
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperGrepUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	// Get the regex from the command line
	if argi >= argc {
		flagSet.Usage()
		os.Exit(1)
	}
	pattern := args[argi]
	argi += 1

	if *pIgnoreCase {
		pattern = "(?i)" + pattern
	}

	// TODO: maybe CompilePOSIX
	regexp, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %s: couldn't compile regex \"%s\"\n",
			args[0], verb, pattern)
		os.Exit(1)
	}

	mapper, _ := NewMapperGrep(
		regexp,
		*pInvert,
	)

	*pargi = argi
	return mapper
}

func mapperGrepUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {regular expression}\n", argv0, verb)
	fmt.Fprintf(o, "Passes through records which match the regular expression.\n")

	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	fmt.Fprint(o, "Options:\n")
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})

	fmt.Fprint(o, `Note that "mlr filter" is more powerful, but requires you to know field names.
By contrast, "mlr grep" allows you to regex-match the entire record. It does
this by formatting each record in memory as DKVP, using command-line-specified
ORS/OFS/OPS, and matching the resulting line against the regex specified
here. In particular, the regex is not applied to the input stream: if you
have CSV with header line "x,y,z" and data line "1,2,3" then the regex will
be matched, not against either of these lines, but against the DKVP line
"x=1,y=2,z=3".  Furthermore, not all the options to system grep are supported,
and this command is intended to be merely a keystroke-saver. To get all the
features of system grep, you can do
  "mlr --odkvp ... | grep ... | mlr --idkvp ..."
`)

}

// ----------------------------------------------------------------
type MapperGrep struct {
	regexp *regexp.Regexp
	invert bool
}

func NewMapperGrep(
	regexp *regexp.Regexp,
	invert bool,
) (*MapperGrep, error) {
	this := &MapperGrep{
		regexp: regexp,
		invert: invert,
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperGrep) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
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
