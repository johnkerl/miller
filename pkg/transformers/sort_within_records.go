package transformers

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

// ----------------------------------------------------------------
const verbNameSortWithinRecords = "sort-within-records"

var SortWithinRecordsSetup = TransformerSetup{
	Verb:         verbNameSortWithinRecords,
	UsageFunc:    transformerSortWithinRecordsUsage,
	ParseCLIFunc: transformerSortWithinRecordsParseCLI,
	IgnoresInput: false,
}

func transformerSortWithinRecordsUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameSortWithinRecords)
	fmt.Fprintln(o, "Outputs records sorted lexically ascending by keys.")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {names}   Sort only these keys; others preserve record order.\n")
	fmt.Fprintf(o, "-r {names}   Like -f but use regular expressions to match field names.\n");
	fmt.Fprintf(o, "             Example: -f '^[xy]' -r sorts keys starting with x or y.\n")
	fmt.Fprintf(o, "             Without -f, -r recursively sorts subobjects/submaps (e.g. for JSON input).\n")
	fmt.Fprintf(o, "-h|--help    Show this message.\n")
}

func transformerSortWithinRecordsParseCLI(
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
	doRecurse := false
	var fieldNames []string = nil
	doRegexes := false

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
			transformerSortWithinRecordsUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-r" {
			if fieldNames != nil {
				doRegexes = true
			} else {
				doRecurse = true
			}

		} else if opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerSortWithinRecordsUsage(os.Stderr)
			os.Exit(1)
		}
	}

	// -r with -f means regex; -r without -f means recurse
	if fieldNames != nil && doRecurse {
		doRegexes = true
		doRecurse = false
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerSortWithinRecords(doRecurse, fieldNames, doRegexes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerSortWithinRecords struct {
	doRecurse   bool
	fieldNames  []string
	fieldSet    map[string]bool
	regex       *regexp.Regexp
	doRegexes   bool
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerSortWithinRecords(
	doRecurse bool,
	fieldNames []string,
	doRegexes bool,
) (*TransformerSortWithinRecords, error) {

	tr := &TransformerSortWithinRecords{
		doRecurse: doRecurse,
		fieldNames: fieldNames,
		doRegexes: doRegexes,
	}

	if fieldNames != nil {
		if doRegexes {
			// Handles "a.*b"i Miller case-insensitive-regex specification
			regexString := fieldNames[0]
			regex, err := lib.CompileMillerRegex(regexString)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"%s %s: cannot compile regex [%s]\n",
					"mlr", verbNameSortWithinRecords, regexString,
				)
				os.Exit(1)
			}
			tr.regex = regex
		} else {
			tr.fieldSet = lib.StringListToSet(fieldNames)
		}
	}

	if fieldNames != nil {
		if doRecurse {
			tr.recordTransformerFunc = tr.transformSelectiveRecursively
		} else {
			tr.recordTransformerFunc = tr.transformSelective
		}
	} else if doRecurse {
		tr.recordTransformerFunc = tr.transformRecursively
	} else {
		tr.recordTransformerFunc = tr.transformNonrecursively
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerSortWithinRecords) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerSortWithinRecords) keyMatches(key string) bool {
	if tr.doRegexes {
		return tr.regex.MatchString(key)
	}
	return tr.fieldSet[key]
}

// ----------------------------------------------------------------
func (tr *TransformerSortWithinRecords) transformSelective(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		var matchingKeys []string
		var restEntries []*mlrval.MlrmapEntry
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if tr.keyMatches(pe.Key) {
				matchingKeys = append(matchingKeys, pe.Key)
			} else {
				restEntries = append(restEntries, pe)
			}
		}
		lib.SortStrings(matchingKeys)
		other := mlrval.NewMlrmapAsRecord()
		for _, key := range matchingKeys {
			other.PutReference(key, inrec.Get(key))
		}
		for _, pe := range restEntries {
			other.PutReference(pe.Key, pe.Value)
		}
		*inrec = *other
	}
	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
}

// ----------------------------------------------------------------
func (tr *TransformerSortWithinRecords) transformSelectiveRecursively(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		var matchingKeys []string
		var restEntries []*mlrval.MlrmapEntry
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if tr.keyMatches(pe.Key) {
				matchingKeys = append(matchingKeys, pe.Key)
			} else {
				restEntries = append(restEntries, pe)
			}
		}
		lib.SortStrings(matchingKeys)
		other := mlrval.NewMlrmapAsRecord()
		for _, key := range matchingKeys {
			val := inrec.Get(key)
			if val != nil {
				if m := val.GetMap(); m != nil {
					m.SortByKeyRecursively()
				}
			}
			other.PutReference(key, val)
		}
		for _, pe := range restEntries {
			if pe.Value != nil {
				if m := pe.Value.GetMap(); m != nil {
					m.SortByKeyRecursively()
				}
			}
			other.PutReference(pe.Key, pe.Value)
		}
		*inrec = *other
	}
	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
}

// ----------------------------------------------------------------
func (tr *TransformerSortWithinRecords) transformNonrecursively(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		inrec.SortByKey()
	}
	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // including end-of-stream marker
}

// ----------------------------------------------------------------
func (tr *TransformerSortWithinRecords) transformRecursively(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		inrec.SortByKeyRecursively()
	}
	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // including end-of-stream marker
}
