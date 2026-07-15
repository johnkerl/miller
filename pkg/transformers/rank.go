package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/transformers/utils"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameRank = "rank"

var rankOptions = []OptionSpec{
	{Flag: "-f", Arg: "{a,b,c}", Type: "csv-list", Desc: "Field name(s) to rank."},
	{Flag: "-g", Arg: "{d,e,f}", Type: "csv-list", Desc: "Optional group-by-field name(s)."},
	{Flag: "--sorted", Type: "bool", Desc: "Promise that the input is already sorted by the field(s) being ranked (within each group, if -g is given). This computes rank in a single streaming pass and O(1) space, by comparing each record's value only to the immediately preceding one, rather than buffering all records to compute an order-independent rank. Produces wrong output if the input is not in fact sorted."},
}

var RankSetup = TransformerSetup{
	Verb:         verbNameRank,
	UsageFunc:    transformerRankUsage,
	ParseCLIFunc: transformerRankParseCLI,
	IgnoresInput: false,
	Options:      rankOptions,
}

func transformerRankUsage(
	o *os.File,
) {
	argv0 := "mlr"
	verb := verbNameRank
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "For each record's value in specified fields, computes the standard\n")
	fmt.Fprintf(o, "competition rank (1,2,2,4,...) of that value among all input records,\n")
	fmt.Fprintf(o, "optionally within groups.\n")
	fmt.Fprintf(o, "E.g. with input records x=10, x=20, x=20, and x=30, emits output records\n")
	fmt.Fprintf(o, "x=10,x_rank=1  x=20,x_rank=2  x=20,x_rank=2  and  x=30,x_rank=4.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Note: by default this is a two-pass algorithm: on the first pass it retains\n")
	fmt.Fprintf(o, "input records and their values; on the second pass it computes ranks and\n")
	fmt.Fprintf(o, "emits output records, in original input order. This means it produces no\n")
	fmt.Fprintf(o, "output until all input is read, but gives correct ranks regardless of input\n")
	fmt.Fprintf(o, "order. Use --sorted for a single-pass streaming alternative.\n")
	fmt.Fprintf(o, "\n")
	WriteVerbOptions(o, rankOptions)
	fmt.Fprintln(o, "Example: mlr rank -f x data/rank-example.csv")
	fmt.Fprintln(o, "Example: mlr rank -f x -g g data/rank-example.csv")
	fmt.Fprintln(o, "Example: mlr sort -f x then rank -f x --sorted data/rank-example.csv")
}

func transformerRankParseCLI(
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

	var rankFieldNames []string = nil
	var groupByFieldNames []string = nil
	doSorted := false

	var err error
	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		if args[argi] == "--" {
			break // All transformers must do this so main-flags can follow verb-flags
		}
		argi++

		switch opt {
		case "-h", "--help":
			transformerRankUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		case "-f":
			rankFieldNames, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "-g":
			groupByFieldNames, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "--sorted":
			doSorted = true

		default:
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	if rankFieldNames == nil {
		return nil, cli.VerbErrorf(verb, "-f field names required")
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerRank(
		rankFieldNames,
		groupByFieldNames,
		doSorted,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerRank struct {
	rankFieldNames    []string
	groupByFieldNames []string
	doSorted          bool

	// Default (unsorted) mode: two-pass. Records are retained on the first
	// pass, along with per-group-per-field percentile-keepers; on the
	// second pass (end of stream) the retained records are decorated with
	// rank fields and emitted in original input order.
	recordsAndContexts []*types.RecordAndContext
	keepers            map[string]map[string]*utils.PercentileKeeper // grouping-key -> field-name -> keeper

	// --sorted mode: single streaming pass, O(1) space. Same shape as the
	// keepers map above, but holding lightweight adjacency state instead of
	// buffered/sorted values.
	sortedStates map[string]map[string]*tRankSortedFieldState // grouping-key -> field-name -> state
}

type tRankSortedFieldState struct {
	count               int64
	rank                int64
	havePreviousValue   bool
	previousValueString string
}

func NewTransformerRank(
	rankFieldNames []string,
	groupByFieldNames []string,
	doSorted bool,
) (*TransformerRank, error) {
	return &TransformerRank{
		rankFieldNames:     rankFieldNames,
		groupByFieldNames:  groupByFieldNames,
		doSorted:           doSorted,
		recordsAndContexts: []*types.RecordAndContext{},
		keepers:            make(map[string]map[string]*utils.PercentileKeeper),
		sortedStates:       make(map[string]map[string]*tRankSortedFieldState),
	}, nil
}

func (tr *TransformerRank) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) error {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if tr.doSorted {
		tr.transformSorted(inrecAndContext, outputRecordsAndContexts)
	} else {
		tr.transformUnsorted(inrecAndContext, outputRecordsAndContexts)
	}
	return nil
}

// transformSorted computes rank in a single pass, O(1) space, by comparing
// each record's value only to the immediately preceding one within its
// group. This is only correct if the caller has ensured the input is
// already sorted by the ranked field(s), e.g. via 'mlr sort'.
func (tr *TransformerRank) transformSorted(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, hasAll := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		if hasAll {
			statesForGroup := tr.sortedStates[groupingKey]
			if statesForGroup == nil {
				statesForGroup = make(map[string]*tRankSortedFieldState)
				tr.sortedStates[groupingKey] = statesForGroup
			}
			for _, rankFieldName := range tr.rankFieldNames {
				value := inrec.Get(rankFieldName)
				if value == nil {
					continue
				}
				state := statesForGroup[rankFieldName]
				if state == nil {
					state = &tRankSortedFieldState{}
					statesForGroup[rankFieldName] = state
				}
				state.count++
				valueString := value.String() // 1, 1.0, and 1.000 are distinct
				if !state.havePreviousValue || valueString != state.previousValueString {
					state.rank = state.count
					state.previousValueString = valueString
					state.havePreviousValue = true
				}
				inrec.PutCopy(rankFieldName+"_rank", mlrval.FromInt(state.rank))
			}
		}
	}
	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
}

// transformUnsorted computes order-independent standard competition rank:
// on the first pass it retains records and accumulates per-group-per-field
// value sets; on the second pass (end of stream) it decorates the retained
// records with rank fields and emits them in original input order.
func (tr *TransformerRank) transformUnsorted(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream { // Not end of stream; pass 1
		inrec := inrecAndContext.Record

		// Append records into a single output list (so that this verb is order-preserving).
		tr.recordsAndContexts = append(tr.recordsAndContexts, inrecAndContext)

		groupingKey, hasAll := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		if hasAll {
			keepersForGroup := tr.keepers[groupingKey]
			if keepersForGroup == nil {
				keepersForGroup = make(map[string]*utils.PercentileKeeper)
				tr.keepers[groupingKey] = keepersForGroup
			}
			for _, rankFieldName := range tr.rankFieldNames {
				value := inrec.Get(rankFieldName)
				if value == nil {
					continue
				}
				keeper := keepersForGroup[rankFieldName]
				if keeper == nil {
					keeper = utils.NewPercentileKeeper(false)
					keepersForGroup[rankFieldName] = keeper
				}
				keeper.Ingest(value)
			}
		}

	} else { // End of stream; pass 2
		// Iterate over the retained records, decorating them with rank fields.
		endOfStreamContext := inrecAndContext.Context

		for _, recordAndContext := range tr.recordsAndContexts {
			outrec := recordAndContext.Record

			groupingKey, hasAll := outrec.GetSelectedValuesJoined(tr.groupByFieldNames)
			if hasAll {
				keepersForGroup := tr.keepers[groupingKey]
				for _, rankFieldName := range tr.rankFieldNames {
					value := outrec.Get(rankFieldName)
					if value == nil {
						continue
					}
					keeper := keepersForGroup[rankFieldName]
					outrec.PutCopy(rankFieldName+"_rank", keeper.EmitRank(value))
				}
			}

			*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(outrec, &endOfStreamContext))
		}
		tr.recordsAndContexts = tr.recordsAndContexts[:0]
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
	}
}
