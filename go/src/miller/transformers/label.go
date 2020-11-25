package mappers

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var LabelSetup = transforming.TransformerSetup{
	Verb:         "label",
	ParseCLIFunc: mapperLabelParseCLI,
	IgnoresInput: false,
}

func mapperLabelParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperLabelUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	// Get the label field names from the command line
	if argi >= argc {
		flagSet.Usage()
		os.Exit(1)
	}
	newNames := lib.SplitString(args[argi], ",")

	argi += 1

	transformer, err := NewTransformerLabel(
		newNames,
	)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return nil
	}

	*pargi = argi
	return transformer
}

func mapperLabelUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {new1,new2,new3,...}\n", argv0, verb)
	fmt.Fprintf(o,
		`Given n comma-separated names, renames the first n fields of each record to
have the respective name. (Fields past the nth are left with their original
names.) Particularly useful with --inidx or --implicit-csv-header, to give
useful names to otherwise integer-indexed fields.
`)
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperLabel struct {
	newNames []string
}

func NewTransformerLabel(
	newNames []string,
) (*MapperLabel, error) {
	// TODO: make this a library function.
	uniquenessChecker := make(map[string]bool)
	for _, newName := range newNames {
		_, ok := uniquenessChecker[newName]
		if ok {
			return nil, errors.New(
				fmt.Sprintf(
					"mlr label: labels must be unique; got duplicate \"%s\"\n",
					newName,
				),
			)
		}
		uniquenessChecker[newName] = true
	}

	this := &MapperLabel{
		newNames: newNames,
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperLabel) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		inrec.Label(this.newNames)
	}
	outputChannel <- inrecAndContext // end-of-stream marker
}
