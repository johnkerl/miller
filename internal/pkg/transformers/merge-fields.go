package transformers

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/transformers/utils"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameMergeFields = "merge-fields"

var MergeFieldsSetup = TransformerSetup{
	Verb:         verbNameMergeFields,
	UsageFunc:    transformerMergeFieldsUsage,
	ParseCLIFunc: transformerMergeFieldsParseCLI,
	IgnoresInput: false,
}

type mergeByType int

const (
	e_MERGE_BY_NAME_LIST mergeByType = iota
	e_MERGE_BY_NAME_REGEX
	e_MERGE_BY_COLLAPSING
	e_MERGE_UNSPECIFIED
)

func transformerMergeFieldsUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	argv0 := "mlr"
	verb := verbNameMergeFields
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Computes univariate statistics for each input record, accumulated across\n")
	fmt.Fprintf(o, "specified fields.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-a {sum,count,...}  Names of accumulators. One or more of:\n")
	utils.ListStats1Accumulators(o)
	fmt.Fprintf(o, "-f {a,b,c}  Value-field names on which to compute statistics. Requires -o.\n")
	fmt.Fprintf(o, "-r {a,b,c}  Regular expressions for value-field names on which to compute\n")
	fmt.Fprintf(o, "            statistics. Requires -o.\n")
	fmt.Fprintf(o, "-c {a,b,c}  Substrings for collapse mode. All fields which have the same names\n")
	fmt.Fprintf(o, "            after removing substrings will be accumulated together. Please see\n")
	fmt.Fprintf(o, "            examples below.\n")
	fmt.Fprintf(o, "-i          Use interpolated percentiles, like R's type=7; default like type=1.\n")
	fmt.Fprintf(o, "            Not sensical for string-valued fields.\n")
	fmt.Fprintf(o, "-o {name}   Output field basename for -f/-r.\n")
	fmt.Fprintf(o, "-k          Keep the input fields which contributed to the output statistics;\n")
	fmt.Fprintf(o, "            the default is to omit them.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "String-valued data make sense unless arithmetic on them is required,\n")
	fmt.Fprintf(o, "e.g. for sum, mean, interpolated percentiles, etc. In case of mixed data,\n")
	fmt.Fprintf(o, "numbers are less than strings.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Example input data: \"a_in_x=1,a_out_x=2,b_in_y=4,b_out_x=8\".\n")
	fmt.Fprintf(o, "Example: %s %s -a sum,count -f a_in_x,a_out_x -o foo\n", argv0, verb)
	fmt.Fprintf(o, "  produces \"b_in_y=4,b_out_x=8,foo_sum=3,foo_count=2\" since \"a_in_x,a_out_x\" are\n")
	fmt.Fprintf(o, "  summed over.\n")
	fmt.Fprintf(o, "Example: %s %s -a sum,count -r in_,out_ -o bar\n", argv0, verb)
	fmt.Fprintf(o, "  produces \"bar_sum=15,bar_count=4\" since all four fields are summed over.\n")
	fmt.Fprintf(o, "Example: %s %s -a sum,count -c in_,out_\n", argv0, verb)
	fmt.Fprintf(o, "  produces \"a_x_sum=3,a_x_count=2,b_y_sum=4,b_y_count=1,b_x_sum=8,b_x_count=1\"\n")
	fmt.Fprintf(o, "  since \"a_in_x\" and \"a_out_x\" both collapse to \"a_x\", \"b_in_y\" collapses to\n")
	fmt.Fprintf(o, "  \"b_y\", and \"b_out_x\" collapses to \"b_x\".\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerMergeFieldsParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	accumulatorNameList := make([]string, 0)
	valueFieldNameList := make([]string, 0)
	outputFieldBasename := ""
	doWhich := e_MERGE_UNSPECIFIED
	keepInputFields := false
	doInterpolatedPercentiles := false

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerMergeFieldsUsage(os.Stdout, true, 0)

		} else if opt == "-a" {
			accumulatorNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			valueFieldNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			doWhich = e_MERGE_BY_NAME_LIST

		} else if opt == "-r" {
			valueFieldNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			doWhich = e_MERGE_BY_NAME_REGEX

		} else if opt == "-c" {
			valueFieldNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			doWhich = e_MERGE_BY_COLLAPSING

		} else if opt == "-o" {
			outputFieldBasename = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-k" {
			keepInputFields = true

		} else if opt == "-i" {
			doInterpolatedPercentiles = true

		} else if opt == "-S" {
			// No-op pass-through for backward compatibility with Miller 5

		} else if opt == "-F" {
			// No-op pass-through for backward compatibility with Miller 5

		} else {
			transformerMergeFieldsUsage(os.Stderr, true, 1)
		}
	}

	// TODO: libify for use across verbs.
	if len(accumulatorNameList) == 0 {
		fmt.Fprintf(os.Stderr, "%s %s: -a option is required.\n", "mlr", verbNameMergeFields)
		fmt.Fprintf(os.Stderr, "Please see %s %s --help for more information.\n", "mlr", verbNameMergeFields)
		os.Exit(1)
	}
	if len(valueFieldNameList) == 0 {
		fmt.Fprintf(os.Stderr, "%s %s: -f option is required.\n", "mlr", verbNameMergeFields)
		fmt.Fprintf(os.Stderr, "Please see %s %s --help for more information.\n", "mlr", verbNameMergeFields)
		os.Exit(1)
	}
	if outputFieldBasename == "" {
		if doWhich == e_MERGE_BY_NAME_LIST || doWhich == e_MERGE_BY_NAME_REGEX {
			transformerMergeFieldsUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerMergeFields(
		accumulatorNameList,
		valueFieldNameList,
		outputFieldBasename,
		doWhich,
		doInterpolatedPercentiles,
		keepInputFields,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
// Given: accumulate count,sum on values x,y group by a,b.
//
// Example input:       Example output:
//   a b x y            a b x_count x_sum y_count y_sum
//   s t 1 2            s t 2       6     2       8
//   u v 3 4            u v 1       3     1       4
//   s t 5 6            u w 1       7     1       9
//   u w 7 9
//
// Multilevel hashmap structure:
// {
//   "s,t" : {                <--- group-by field names
//     "x" : {                  <--- value field name
//       "count" : Stats1CountAccumulator object,
//       "sum"   : Stats1SumAccumulator  object
//     },
//     "y" : {
//       "count" : Stats1CountAccumulator object,
//       "sum"   : Stats1SumAccumulator  object
//     },
//   },
//   "u,v" : {
//     "x" : {
//       "count" : Stats1CountAccumulator object,
//       "sum"   : Stats1SumAccumulator  object
//     },
//     "y" : {
//       "count" : Stats1CountAccumulator object,
//       "sum"   : Stats1SumAccumulator  object
//     },
//   },
//   "u,w" : {
//     "x" : {
//       "count" : Stats1CountAccumulator object,
//       "sum"   : Stats1SumAccumulator  object
//     },
//     "y" : {
//       "count" : Stats1CountAccumulator object,
//       "sum"   : Stats1SumAccumulator  object
//     },
//   },
// }

type TransformerMergeFields struct {
	// Input:
	accumulatorNameList       []string
	valueFieldNameList        []string
	outputFieldBasename       string
	doInterpolatedPercentiles bool
	keepInputFields           bool

	// State:
	accumulatorFactory    *utils.Stats1AccumulatorFactory
	valueFieldNameRegexes []*regexp.Regexp

	// Ordered map from accumulator name to accumulator
	namedAccumulators *lib.OrderedMap

	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerMergeFields(
	accumulatorNameList []string,
	valueFieldNameList []string,
	outputFieldBasename string,
	doWhich mergeByType,
	doInterpolatedPercentiles bool,
	keepInputFields bool,
) (*TransformerMergeFields, error) {

	for _, accumulatorName := range accumulatorNameList {
		if !utils.ValidateStats1AccumulatorName(accumulatorName) {
			return nil, errors.New(
				fmt.Sprintf(
					"%s %s: accumulator \"%s\" not found.\n",
					"mlr", verbNameMergeFields, accumulatorName,
				),
			)
		}
	}

	tr := &TransformerMergeFields{
		accumulatorNameList:       accumulatorNameList,
		valueFieldNameList:        valueFieldNameList,
		outputFieldBasename:       outputFieldBasename,
		doInterpolatedPercentiles: doInterpolatedPercentiles,
		keepInputFields:           keepInputFields,
		accumulatorFactory:        utils.NewStats1AccumulatorFactory(),
		namedAccumulators:         lib.NewOrderedMap(),
	}

	tr.valueFieldNameRegexes = make([]*regexp.Regexp, len(valueFieldNameList))
	for i, regexString := range valueFieldNameList {
		// Handles "a.*b"i Miller case-insensitive-regex specification
		regex, err := lib.CompileMillerRegex(regexString)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"%s %s: cannot compile regex [%s]\n",
				"mlr", verbNameCut, regexString,
			)
			os.Exit(1)
		}
		tr.valueFieldNameRegexes[i] = regex
	}

	for _, accumulatorName := range accumulatorNameList {
		accumulator := tr.accumulatorFactory.MakeNamedAccumulator(
			accumulatorName,
			"", // grouping-key used for stats1, not here
			outputFieldBasename,
			doInterpolatedPercentiles,
		)
		tr.namedAccumulators.Put(accumulatorName, accumulator)
	}

	if doWhich == e_MERGE_BY_NAME_LIST {
		tr.recordTransformerFunc = tr.transformByNameList
	} else if doWhich == e_MERGE_BY_NAME_REGEX {
		tr.recordTransformerFunc = tr.transformByNameRegex
	} else if doWhich == e_MERGE_BY_COLLAPSING {
		tr.recordTransformerFunc = tr.transformByCollapsing
	} else {
		lib.InternalCodingErrorIf(true)
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerMergeFields) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerMergeFields) transformByNameList(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if inrecAndContext.EndOfStream {
		outputChannel <- inrecAndContext // end-of-stream marker
		return
	}

	inrec := inrecAndContext.Record

	for pa := tr.namedAccumulators.Head; pa != nil; pa = pa.Next {
		accumulator := pa.Value.(*utils.Stats1NamedAccumulator)
		accumulator.Reset() // re-use from one record to the next
	}

	for _, valueFieldName := range tr.valueFieldNameList {
		mvalue := inrec.Get(valueFieldName)
		if mvalue == nil { // key not present
			continue
		}

		if mvalue.IsEmpty() { // key present with empty value
			if !tr.keepInputFields {
				inrec.Remove(valueFieldName)
			}
			continue
		}

		for pa := tr.namedAccumulators.Head; pa != nil; pa = pa.Next {
			accumulator := pa.Value.(*utils.Stats1NamedAccumulator)
			accumulator.Ingest(mvalue)
		}

		if !tr.keepInputFields {
			inrec.Remove(valueFieldName)
		}
	}

	for pa := tr.namedAccumulators.Head; pa != nil; pa = pa.Next {
		accumulator := pa.Value.(*utils.Stats1NamedAccumulator)
		key, value := accumulator.Emit()
		inrec.PutReference(key, value)
	}

	outputChannel <- inrecAndContext
}

// ----------------------------------------------------------------
func (tr *TransformerMergeFields) transformByNameRegex(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if inrecAndContext.EndOfStream {
		outputChannel <- inrecAndContext // end-of-stream marker
		return
	}

	inrec := inrecAndContext.Record

	for pa := tr.namedAccumulators.Head; pa != nil; pa = pa.Next {
		accumulator := pa.Value.(*utils.Stats1NamedAccumulator)
		accumulator.Reset() // re-use from one record to the next
	}

	for pe := inrec.Head; pe != nil; /* increment inside loop*/ {
		valueFieldName := pe.Key

		matched := false
		for _, valueFieldNameRegex := range tr.valueFieldNameRegexes {
			if valueFieldNameRegex.MatchString(pe.Key) {
				matched = true
				break
			}
		}
		if !matched {
			pe = pe.Next
			continue
		}

		mvalue := inrec.Get(valueFieldName)

		if mvalue == nil { // Key not present
			pe = pe.Next
			continue
		}

		if mvalue.IsEmpty() { // Key present with empty value
			if !tr.keepInputFields { // We are modifying the record while iterating over it.
				next := pe.Next
				inrec.Unlink(pe)
				pe = next
			} else {
				pe = pe.Next
			}
			continue
		}

		for pa := tr.namedAccumulators.Head; pa != nil; pa = pa.Next {
			accumulator := pa.Value.(*utils.Stats1NamedAccumulator)
			accumulator.Ingest(mvalue)
		}

		if !tr.keepInputFields { // We are modifying the record while iterating over it.
			next := pe.Next
			inrec.Unlink(pe)
			pe = next
		} else {
			pe = pe.Next
		}
	}

	for pa := tr.namedAccumulators.Head; pa != nil; pa = pa.Next {
		accumulator := pa.Value.(*utils.Stats1NamedAccumulator)
		key, value := accumulator.Emit()
		inrec.PutReference(key, value)
	}

	outputChannel <- inrecAndContext
}

// ----------------------------------------------------------------
// mlr merge-fields -c in_,out_ -a sum
// a_in_x  1     a_sum_x 3
// a_out_x 2     b_sum_y 4
// b_in_y  4     b_sum_x 8
// b_out_x 8

func (tr *TransformerMergeFields) transformByCollapsing(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if inrecAndContext.EndOfStream {
		outputChannel <- inrecAndContext // end-of-stream marker
		return
	}

	inrec := inrecAndContext.Record
	tr.accumulatorFactory.Reset() // discard cached percentile-keepers

	// Ordered map from short name to accumulator name to accumulator
	collapseAccumulators := lib.NewOrderedMap()

	for pe := inrec.Head; pe != nil; /* increment inside loop */ {
		valueFieldName := pe.Key

		matched := false
		shortName := ""
		for _, valueFieldNameRegex := range tr.valueFieldNameRegexes {
			matched = valueFieldNameRegex.MatchString(pe.Key)
			if matched {
				// TODO: comment re matrix
				shortName = lib.RegexSubCompiled(valueFieldName, valueFieldNameRegex, "", nil)
				break
			}
		}
		if !matched {
			pe = pe.Next
			continue
		}

		mvalue := inrec.Get(valueFieldName)
		if mvalue == nil { // Key aesent
			pe = pe.Next
			continue
		}

		var namedAccumulators *lib.OrderedMap
		iNamedAccumulators := collapseAccumulators.Get(shortName)
		if iNamedAccumulators == nil {
			namedAccumulators = lib.NewOrderedMap()
			for _, accumulatorName := range tr.accumulatorNameList {
				accumulator := tr.accumulatorFactory.MakeNamedAccumulator(
					accumulatorName,
					"", // grouping-key used for stats1, not here
					shortName,
					tr.doInterpolatedPercentiles,
				)
				namedAccumulators.Put(accumulatorName, accumulator)
			}
			collapseAccumulators.Put(shortName, namedAccumulators)
		} else {
			namedAccumulators = iNamedAccumulators.(*lib.OrderedMap)
		}

		// The accumulator has been initialized with default values; continue
		// here. (If we were to continue before namedAccumulators.Put(...) we
		// would be failing to construct the accumulator.)
		if mvalue.IsEmpty() { // key present with empty value
			if !tr.keepInputFields { // We are modifying the record while iterating over it.
				next := pe.Next
				inrec.Unlink(pe)
				pe = next
			} else {
				pe = pe.Next
			}
			continue
		}

		for pa := namedAccumulators.Head; pa != nil; pa = pa.Next {
			accumulator := pa.Value.(*utils.Stats1NamedAccumulator)
			accumulator.Ingest(mvalue)
		}

		if !tr.keepInputFields { // We are modifying the record while iterating over it.
			next := pe.Next
			inrec.Unlink(pe)
			pe = next
		} else {
			pe = pe.Next
		}
	}

	for ps := collapseAccumulators.Head; ps != nil; ps = ps.Next {
		namedAccumulators := ps.Value.(*lib.OrderedMap)
		for pa := namedAccumulators.Head; pa != nil; pa = pa.Next {
			accumulator := pa.Value.(*utils.Stats1NamedAccumulator)
			key, value := accumulator.Emit()
			inrec.PutReference(key, value)
		}
	}

	outputChannel <- inrecAndContext
}
