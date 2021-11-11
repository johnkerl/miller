package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameFraction = "fraction"

var FractionSetup = TransformerSetup{
	Verb:         verbNameFraction,
	UsageFunc:    transformerFractionUsage,
	ParseCLIFunc: transformerFractionParseCLI,
	IgnoresInput: false,
}

func transformerFractionUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	argv0 := "mlr"
	verb := verbNameFraction
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "For each record's value in specified fields, computes the ratio of that\n")
	fmt.Fprintf(o, "value to the sum of values in that field over all input records.\n")
	fmt.Fprintf(o, "E.g. with input records  x=1  x=2  x=3  and  x=4, emits output records\n")
	fmt.Fprintf(o, "x=1,x_fraction=0.1  x=2,x_fraction=0.2  x=3,x_fraction=0.3  and  x=4,x_fraction=0.4\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Note: this is internally a two-pass algorithm: on the first pass it retains\n")
	fmt.Fprintf(o, "input records and accumulates sums; on the second pass it computes quotients\n")
	fmt.Fprintf(o, "and emits output records. This means it produces no output until all input is read.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {a,b,c}    Field name(s) for fraction calculation\n")
	fmt.Fprintf(o, "-g {d,e,f}    Optional group-by-field name(s) for fraction counts\n")
	fmt.Fprintf(o, "-p            Produce percents [0..100], not fractions [0..1]. Output field names\n")
	fmt.Fprintf(o, "              end with \"_percent\" rather than \"_fraction\"\n")
	fmt.Fprintf(o, "-c            Produce cumulative distributions, i.e. running sums: each output\n")
	fmt.Fprintf(o, "              value folds in the sum of the previous for the specified group\n")
	fmt.Fprintf(o, "              E.g. with input records  x=1  x=2  x=3  and  x=4, emits output records\n")
	fmt.Fprintf(o, "              x=1,x_cumulative_fraction=0.1  x=2,x_cumulative_fraction=0.3\n")
	fmt.Fprintf(o, "              x=3,x_cumulative_fraction=0.6  and  x=4,x_cumulative_fraction=1.0\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerFractionParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	var fractionFieldNames []string = nil
	var groupByFieldNames []string = nil
	doPercents := false
	doCumu := false

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerFractionUsage(os.Stdout, true, 0)

		} else if opt == "-f" {
			fractionFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-p" {
			doPercents = true

		} else if opt == "-c" {
			doCumu = true

		} else {
			transformerFractionUsage(os.Stderr, true, 1)
		}
	}

	if fractionFieldNames == nil {
		transformerFractionUsage(os.Stderr, true, 1)
	}

	transformer, err := NewTransformerFraction(
		fractionFieldNames,
		groupByFieldNames,
		doPercents,
		doCumu,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerFraction struct {
	fractionFieldNames []string
	groupByFieldNames  []string
	doCumu             bool

	recordsAndContexts *list.List
	// Two-level map: Group-by field names are the first keyset;
	// fraction field names are keys into the second.
	sums  map[string]map[string]*types.Mlrval
	cumus map[string]map[string]*types.Mlrval

	outputFieldNameSuffix string        // "_fraction" or "_percent"
	multiplier            *types.Mlrval // 1.0 for fraction or 100.0 for percent
	zero                  *types.Mlrval
}

// ----------------------------------------------------------------
func NewTransformerFraction(
	fractionFieldNames []string,
	groupByFieldNames []string,
	doPercents bool,
	doCumu bool,
) (*TransformerFraction, error) {

	recordsAndContexts := list.New()
	sums := make(map[string]map[string]*types.Mlrval)
	cumus := make(map[string]map[string]*types.Mlrval)

	var multiplier *types.Mlrval
	var outputFieldNameSuffix string
	if doPercents {
		multiplier = types.MlrvalFromInt(100)
		if doCumu {
			outputFieldNameSuffix = "_cumulative_percent"
		} else {
			outputFieldNameSuffix = "_percent"
		}
	} else {
		multiplier = types.MlrvalFromInt(1)
		if doCumu {
			outputFieldNameSuffix = "_cumulative_fraction"
		} else {
			outputFieldNameSuffix = "_fraction"
		}
	}

	zero := types.MlrvalFromInt(0)

	return &TransformerFraction{
		fractionFieldNames:    fractionFieldNames,
		groupByFieldNames:     groupByFieldNames,
		doCumu:                doCumu,
		recordsAndContexts:    recordsAndContexts,
		sums:                  sums,
		cumus:                 cumus,
		outputFieldNameSuffix: outputFieldNameSuffix,
		multiplier:            multiplier,
		zero:                  zero,
	}, nil
}

// ----------------------------------------------------------------

func (tr *TransformerFraction) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream { // Not end of stream; pass 1
		inrec := inrecAndContext.Record

		// Append records into a single output list (so that this verb is order-preserving).
		tr.recordsAndContexts.PushBack(inrecAndContext)

		// Accumulate sums of fraction-field values grouped by group-by field names
		groupingKey, hasAll := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)

		if hasAll {
			sumsForGroup := tr.sums[groupingKey]
			var cumusForGroup map[string]*types.Mlrval = nil
			if sumsForGroup == nil {
				sumsForGroup = make(map[string]*types.Mlrval)
				tr.sums[groupingKey] = sumsForGroup
				cumusForGroup = make(map[string]*types.Mlrval)
				tr.cumus[groupingKey] = cumusForGroup
			}
			for _, fractionFieldName := range tr.fractionFieldNames {
				value := inrec.Get(fractionFieldName)
				if value != nil {
					value.AssertNumeric() // may fatal the process
					sum := sumsForGroup[fractionFieldName]
					if sum == nil { // First value for group
						sumsForGroup[fractionFieldName] = value.Copy()
						cumusForGroup[fractionFieldName] = tr.zero
					} else {
						sumsForGroup[fractionFieldName] = types.BIF_plus_binary(sum, value)
					}
				}
			}
		}

	} else { // End of stream; pass 2
		// Iterate over the retained records, decorating them with fraction fields.
		endOfStreamContext := inrecAndContext.Context

		for {
			element := tr.recordsAndContexts.Front()
			if element == nil {
				break
			}
			tr.recordsAndContexts.Remove(element)
			recordAndContext := element.Value.(*types.RecordAndContext)
			outrec := recordAndContext.Record

			groupingKey, hasAll := outrec.GetSelectedValuesJoined(tr.groupByFieldNames)
			if hasAll {
				sumsForGroup := tr.sums[groupingKey]
				cumusForGroup := tr.cumus[groupingKey]
				lib.InternalCodingErrorIf(sumsForGroup == nil) // should have been populated on pass 1

				for _, fractionFieldName := range tr.fractionFieldNames {
					value := outrec.Get(fractionFieldName)
					if value != nil {
						value.AssertNumeric() // may fatal the process

						var numerator *types.Mlrval = nil
						var cumu *types.Mlrval = nil
						var outputValue *types.Mlrval = nil

						if tr.doCumu {
							cumu = cumusForGroup[fractionFieldName]
							numerator = types.BIF_plus_binary(value, cumu)
						} else {
							numerator = value
						}

						denominator := sumsForGroup[fractionFieldName]
						if !types.MlrvalEqualsAsBool(value, tr.zero) {
							outputValue = types.BIF_divide(numerator, denominator)
							outputValue = types.BIF_times(outputValue, tr.multiplier)
						} else {
							outputValue = types.MLRVAL_ERROR
						}

						outrec.PutCopy(
							fractionFieldName+tr.outputFieldNameSuffix,
							outputValue,
						)

						if tr.doCumu {
							cumusForGroup[fractionFieldName] = types.BIF_plus_binary(cumu, value)
						}
					}
				}
			}

			outputChannel <- types.NewRecordAndContext(outrec, &endOfStreamContext)
		}
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
