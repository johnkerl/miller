package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameFraction = "fraction"

var FractionSetup = transforming.TransformerSetup{
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
	argv0 := lib.MlrExeName()
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
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

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
			fractionFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

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

	transformer, _ := NewTransformerFraction(
		fractionFieldNames,
		groupByFieldNames,
		doPercents,
		doCumu,
	)

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
		multiplier = types.MlrvalPointerFromInt(100)
		if doCumu {
			outputFieldNameSuffix = "_cumulative_percent"
		} else {
			outputFieldNameSuffix = "_percent"
		}
	} else {
		multiplier = types.MlrvalPointerFromInt(1)
		if doCumu {
			outputFieldNameSuffix = "_cumulative_fraction"
		} else {
			outputFieldNameSuffix = "_fraction"
		}
	}

	zero := types.MlrvalPointerFromInt(0)

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
func (this *TransformerFraction) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream { // Not end of stream; pass 1
		inrec := inrecAndContext.Record

		// Append records into a single output list (so that this verb is order-preserving).
		this.recordsAndContexts.PushBack(inrecAndContext)

		// Accumulate sums of fraction-field values grouped by group-by field names
		groupingKey, hasAll := inrec.GetSelectedValuesJoined(this.groupByFieldNames)

		if hasAll {
			sumsForGroup := this.sums[groupingKey]
			var cumusForGroup map[string]*types.Mlrval = nil
			if sumsForGroup == nil {
				sumsForGroup = make(map[string]*types.Mlrval)
				this.sums[groupingKey] = sumsForGroup
				cumusForGroup = make(map[string]*types.Mlrval)
				this.cumus[groupingKey] = cumusForGroup
			}
			for _, fractionFieldName := range this.fractionFieldNames {
				value := inrec.Get(fractionFieldName)
				if value != nil {
					value.AssertNumeric() // may fatal the process
					sum := sumsForGroup[fractionFieldName]
					if sum == nil { // First value for group
						sumsForGroup[fractionFieldName] = value.Copy()
						cumusForGroup[fractionFieldName] = this.zero
					} else {
						sumsForGroup[fractionFieldName] = types.MlrvalBinaryPlus(sum, value)
					}
				}
			}
		}

	} else { // End of stream; pass 2
		// Iterate over the retained records, decorating them with fraction fields.
		endOfStreamContext := inrecAndContext.Context

		for {
			element := this.recordsAndContexts.Front()
			if element == nil {
				break
			}
			this.recordsAndContexts.Remove(element)
			recordAndContext := element.Value.(*types.RecordAndContext)
			outrec := recordAndContext.Record

			groupingKey, hasAll := outrec.GetSelectedValuesJoined(this.groupByFieldNames)
			if hasAll {
				sumsForGroup := this.sums[groupingKey]
				cumusForGroup := this.cumus[groupingKey]
				lib.InternalCodingErrorIf(sumsForGroup == nil) // should have been populated on pass 1

				for _, fractionFieldName := range this.fractionFieldNames {
					value := outrec.Get(fractionFieldName)
					if value != nil {
						value.AssertNumeric() // may fatal the process

						var numerator *types.Mlrval = nil
						var cumu *types.Mlrval = nil
						var outputValue *types.Mlrval = nil

						if this.doCumu {
							cumu = cumusForGroup[fractionFieldName]
							numerator = types.MlrvalBinaryPlus(value, cumu)
						} else {
							numerator = value
						}

						denominator := sumsForGroup[fractionFieldName]
						if !types.MlrvalEqualsAsBool(value, this.zero) {
							outputValue = types.MlrvalDivide(numerator, denominator)
							outputValue = types.MlrvalTimes(outputValue, this.multiplier)
						} else {
							outputValue = types.MLRVAL_ERROR
						}

						outrec.PutCopy(
							fractionFieldName+this.outputFieldNameSuffix,
							outputValue,
						)

						if this.doCumu {
							cumusForGroup[fractionFieldName] = types.MlrvalBinaryPlus(cumu, value)
						}
					}
				}
			}

			outputChannel <- types.NewRecordAndContext(outrec, &endOfStreamContext)
		}
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
