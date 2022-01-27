package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/transformers/utils"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameTop = "top"
const verbTopDefaultOutputFieldName = "top_idx"

var TopSetup = TransformerSetup{
	Verb:         verbNameTop,
	UsageFunc:    transformerTopUsage,
	ParseCLIFunc: transformerTopParseCLI,
	IgnoresInput: false,
}

func transformerTopUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	argv0 := "mlr"
	verb := verbNameTop
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "-f {a,b,c}    Value-field names for top counts.\n")
	fmt.Fprintf(o, "-g {d,e,f}    Optional group-by-field names for top counts.\n")
	fmt.Fprintf(o, "-n {count}    How many records to print per category; default 1.\n")
	fmt.Fprintf(o, "-a            Print all fields for top-value records; default is\n")
	fmt.Fprintf(o, "              to print only value and group-by fields. Requires a single\n")
	fmt.Fprintf(o, "              value-field name only.\n")
	fmt.Fprintf(o, "--min         Print top smallest values; default is top largest values.\n")
	fmt.Fprintf(o, "-F            Keep top values as floats even if they look like integers.\n")
	fmt.Fprintf(o, "-o {name}     Field name for output indices. Default \"%s\".\n", verbTopDefaultOutputFieldName)
	fmt.Fprintf(o, "              This is ignored if -a is used.\n")

	fmt.Fprintf(o, "Prints the n records with smallest/largest values at specified fields,\n")
	fmt.Fprintf(o, "optionally by category. If -a is given, then the top records are emitted\n")
	fmt.Fprintf(o, "with the same fields as they appeared in the input. Without -a, only fields\n")
	fmt.Fprintf(o, "from -f, fields from -g, and the top-index field are emitted. For more information\n")
	fmt.Fprintf(o, "please see https://miller.readthedocs.io/en/latest/reference-verbs#top\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerTopParseCLI(
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

	// Parse local flags
	topCount := int64(1)
	var valueFieldNames []string = nil
	var groupByFieldNames []string = nil
	showFullRecords := false
	doMax := true
	outputFieldName := verbTopDefaultOutputFieldName

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
			transformerTopUsage(os.Stdout, true, 0)

		} else if opt == "-n" {
			topCount = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)
		} else if opt == "-f" {
			valueFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
		} else if opt == "-a" {
			showFullRecords = true
		} else if opt == "--max" {
			doMax = true
		} else if opt == "--min" {
			doMax = false
		} else if opt == "-F" {
			// Ignored in Miller 6; allowed for command-line backward compatibility
		} else if opt == "-o" {
			outputFieldName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerTopUsage(os.Stderr, true, 1)
		}
	}

	if valueFieldNames == nil {
		transformerTopUsage(os.Stderr, true, 1)
	}
	if len(valueFieldNames) > 1 && showFullRecords {
		transformerTopUsage(os.Stderr, true, 1)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, _ := NewTransformerTop(
		topCount,
		valueFieldNames,
		groupByFieldNames,
		showFullRecords,
		doMax,
		outputFieldName,
	)

	return transformer
}

// ----------------------------------------------------------------
type TransformerTop struct {
	topCount          int64
	valueFieldNames   []string
	groupByFieldNames []string
	showFullRecords   bool
	doMax             bool
	outputFieldName   string

	// Two-level map from grouping key (string of joined-together group-by field values),
	// to string value-field name, to *utils.TopKeeper
	groups                           *lib.OrderedMap
	groupingKeysToGroupByFieldValues map[string][]*mlrval.Mlrval
}

// ----------------------------------------------------------------
func NewTransformerTop(
	topCount int64,
	valueFieldNames []string,
	groupByFieldNames []string,
	showFullRecords bool,
	doMax bool,
	outputFieldName string,
) (*TransformerTop, error) {

	tr := &TransformerTop{
		topCount:          topCount,
		valueFieldNames:   valueFieldNames,
		groupByFieldNames: groupByFieldNames,
		showFullRecords:   showFullRecords,
		doMax:             doMax,
		outputFieldName:   outputFieldName,

		groups:                           lib.NewOrderedMap(),
		groupingKeysToGroupByFieldValues: make(map[string][]*mlrval.Mlrval),
	}

	return tr, nil
}

func (tr *TransformerTop) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		tr.ingest(inrecAndContext)
	} else {
		tr.emit(inrecAndContext, outputRecordsAndContexts)
	}
}

func (tr *TransformerTop) ingest(
	inrecAndContext *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record

	// ["s", "t"]
	valueFieldValues, fok := inrec.ReferenceSelectedValues(tr.valueFieldNames)
	groupingKey, groupByFieldValues, gok := inrec.GetSelectedValuesAndJoined(tr.groupByFieldNames)

	// Heterogeneous-data case -- not all sought fields were present in record
	if !fok || !gok {
		return
	}
	iSecondLevel := tr.groups.Get(groupingKey)
	var secondLevel *lib.OrderedMap = nil
	if iSecondLevel == nil {
		secondLevel = lib.NewOrderedMap()
		tr.groups.Put(groupingKey, secondLevel)
		tr.groupingKeysToGroupByFieldValues[groupingKey] = groupByFieldValues
	} else {
		secondLevel = iSecondLevel.(*lib.OrderedMap)
	}

	// for "x", "y" and "1", "2"
	for i := range tr.valueFieldNames {
		valueFieldName := tr.valueFieldNames[i]
		valueFieldValue := valueFieldValues[i]

		iTopKeeper := secondLevel.Get(valueFieldName)
		var topKeeper *utils.TopKeeper
		if iTopKeeper == nil {
			topKeeper = utils.NewTopKeeper(tr.topCount, tr.doMax)
			secondLevel.Put(valueFieldName, topKeeper)
		} else {
			topKeeper = iTopKeeper.(*utils.TopKeeper)
		}

		var maybeRecordAndContext *types.RecordAndContext = nil
		if tr.showFullRecords {
			maybeRecordAndContext = inrecAndContext
		}
		topKeeper.Add(
			valueFieldValue,
			maybeRecordAndContext,
		)
	}

}

// ----------------------------------------------------------------
func (tr *TransformerTop) emit(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
) {
	for pa := tr.groups.Head; pa != nil; pa = pa.Next {
		groupingKey := pa.Key
		secondLevel := pa.Value.(*lib.OrderedMap)
		groupByFieldValues := tr.groupingKeysToGroupByFieldValues[groupingKey]

		// Above we required that there be only one value field in the
		// show-full-records case. That's because here, we print each record at most
		// once, which would need a change in the format presented as output.
		if tr.showFullRecords {
			for pb := secondLevel.Head; pb != nil; pb = pb.Next {
				topKeeper := pb.Value.(*utils.TopKeeper)
				for i := int64(0); i < topKeeper.GetSize(); i++ {
					outputRecordsAndContexts.PushBack(topKeeper.TopRecordsAndContexts[i].Copy())
				}
			}

		} else {

			for i := int64(0); i < tr.topCount; i++ {
				newrec := mlrval.NewMlrmapAsRecord()

				// Add in a=s,b=t fields:
				for j := range tr.groupByFieldNames {
					newrec.PutCopy(tr.groupByFieldNames[j], groupByFieldValues[j])
				}

				// Add in fields such as x_top_1=#
				// for "x", "y"
				for pb := secondLevel.Head; pb != nil; pb = pb.Next {
					valueFieldName := pb.Key
					topKeeper := pb.Value.(*utils.TopKeeper)
					key := valueFieldName + "_top"
					if i < topKeeper.GetSize() {
						newrec.PutReference(tr.outputFieldName, mlrval.FromInt(i+1))
						newrec.PutReference(key, topKeeper.TopValues[i].Copy())
					} else {
						newrec.PutReference(tr.outputFieldName, mlrval.FromInt(i+1))
						newrec.PutCopy(key, mlrval.VOID)
					}
				}

				outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
			}
		}
	}

	outputRecordsAndContexts.PushBack(inrecAndContext) // emit end-of-stream marker
}
