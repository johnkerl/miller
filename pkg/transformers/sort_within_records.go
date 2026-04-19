package transformers

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/facette/natsort"
	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

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
	fmt.Fprintf(o, "-r {regex}   Sort only keys matching this regex; others preserve record order.\n")
	fmt.Fprintf(o, "             Example: -r '^[xy]' sorts keys starting with x or y.\n")
	fmt.Fprintf(o, "             With no regex argument, -r recursively sorts subobjects/submaps\n")
	fmt.Fprintf(o, "             (e.g. for JSON input), or combines with -f to treat names as regex.\n")
	fmt.Fprintf(o, "-n           Sort field names naturally (e.g. 2 before 12). Combines with -f/-r.\n")
	fmt.Fprintf(o, "-h|--help    Show this message.\n")
}

func transformerSortWithinRecordsParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++
	doRecurse := false
	var fieldNames []string = nil
	doRegexes := false
	doNatural := false

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
			return nil, cli.ErrHelpRequested

		} else if opt == "-r" {
			// If the next token exists and isn't another flag, consume it as
			// the regex pattern. Otherwise -r is arity-0: combined with a
			// preceding -f it means regex mode; standalone it means recursive.
			if argi < argc && !strings.HasPrefix(args[argi], "-") {
				fieldNames = []string{args[argi]}
				argi++
				doRegexes = true
			} else if fieldNames != nil {
				doRegexes = true
			} else {
				doRecurse = true
			}

		} else if opt == "-f" {
			names, err := cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			fieldNames = names

		} else if opt == "-n" {
			doNatural = true

		} else {
			return nil, cli.VerbErrorf(verbNameSortWithinRecords, "option \"%s\" not recognized", opt)
		}
	}

	// -r with -f means regex; -r without -f means recurse
	if fieldNames != nil && doRecurse {
		doRegexes = true
		doRecurse = false
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerSortWithinRecords(doRecurse, fieldNames, doRegexes, doNatural)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerSortWithinRecords struct {
	doRecurse             bool
	doNatural             bool
	fieldNames            []string
	fieldSet              map[string]bool
	regex                 *regexp.Regexp
	doRegexes             bool
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerSortWithinRecords(
	doRecurse bool,
	fieldNames []string,
	doRegexes bool,
	doNatural bool,
) (*TransformerSortWithinRecords, error) {

	tr := &TransformerSortWithinRecords{
		doRecurse:  doRecurse,
		doNatural:  doNatural,
		fieldNames: fieldNames,
		doRegexes:  doRegexes,
	}

	if fieldNames != nil {
		if doRegexes {
			if len(fieldNames) > 1 {
				return nil, cli.VerbErrorf(
					verbNameSortWithinRecords,
					"regex mode takes a single pattern; got %d names: %s. "+
						"Use alternation in the regex (e.g. 'a|b') instead of a comma-list.",
					len(fieldNames), strings.Join(fieldNames, ","),
				)
			}
			// Handles "a.*b"i Miller case-insensitive-regex specification
			regexString := fieldNames[0]
			regex, err := lib.CompileMillerRegex(regexString)
			if err != nil {
				return nil, cli.VerbErrorf(
					verbNameSortWithinRecords,
					"cannot compile regex [%s]", regexString,
				)
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

func (tr *TransformerSortWithinRecords) sortKeys(keys []string) {
	if tr.doNatural {
		sort.SliceStable(keys, func(i, j int) bool {
			return natsort.Compare(keys[i], keys[j])
		})
	} else {
		slices.Sort(keys)
	}
}

func (tr *TransformerSortWithinRecords) sortMapByKey(m *mlrval.Mlrmap) {
	keys := m.GetKeys()
	tr.sortKeys(keys)
	other := mlrval.NewMlrmapAsRecord()
	for _, key := range keys {
		other.PutReference(key, m.Get(key))
	}
	*m = *other
}

func (tr *TransformerSortWithinRecords) sortMapByKeyRecursively(m *mlrval.Mlrmap) {
	keys := m.GetKeys()
	tr.sortKeys(keys)
	other := mlrval.NewMlrmapAsRecord()
	for _, key := range keys {
		val := m.Get(key)
		if val != nil {
			if sub := val.GetMap(); sub != nil {
				tr.sortMapByKeyRecursively(sub)
			}
		}
		other.PutReference(key, val)
	}
	*m = *other
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
		tr.sortKeys(matchingKeys)
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
		tr.sortKeys(matchingKeys)
		other := mlrval.NewMlrmapAsRecord()
		for _, key := range matchingKeys {
			val := inrec.Get(key)
			if val != nil {
				if m := val.GetMap(); m != nil {
					tr.sortMapByKeyRecursively(m)
				}
			}
			other.PutReference(key, val)
		}
		for _, pe := range restEntries {
			if pe.Value != nil {
				if m := pe.Value.GetMap(); m != nil {
					tr.sortMapByKeyRecursively(m)
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
		if tr.doNatural {
			tr.sortMapByKey(inrec)
		} else {
			inrec.SortByKey()
		}
	}
	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // including end-of-stream marker
}

func (tr *TransformerSortWithinRecords) transformRecursively(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		if tr.doNatural {
			tr.sortMapByKeyRecursively(inrec)
		} else {
			inrec.SortByKeyRecursively()
		}
	}
	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // including end-of-stream marker
}
