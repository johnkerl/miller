package transformers

// ================================================================
// WIDE:
//          time           X          Y           Z
// 1  2009-01-01  0.65473572  2.4520609 -1.46570942
// 2  2009-01-02 -0.89248112  0.2154713 -2.05357735
// 3  2009-01-03  0.98012375  1.3179287  4.64248357
// 4  2009-01-04  0.35397376  3.3765645 -0.25237774
// 5  2009-01-05  2.19357813  1.3477511  0.09719105

// LONG:
//          time  item       price
// 1  2009-01-01     X  0.65473572
// 2  2009-01-02     X -0.89248112
// 3  2009-01-03     X  0.98012375
// 4  2009-01-04     X  0.35397376
// 5  2009-01-05     X  2.19357813
// 6  2009-01-01     Y  2.45206093
// 7  2009-01-02     Y  0.21547134
// 8  2009-01-03     Y  1.31792866
// 9  2009-01-04     Y  3.37656453
// 10 2009-01-05     Y  1.34775108
// 11 2009-01-01     Z -1.46570942
// 12 2009-01-02     Z -2.05357735
// 13 2009-01-03     Z  4.64248357
// 14 2009-01-04     Z -0.25237774
// 15 2009-01-05     Z  0.09719105

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameReshape = "reshape"

var ReshapeSetup = TransformerSetup{
	Verb:         verbNameReshape,
	UsageFunc:    transformerReshapeUsage,
	ParseCLIFunc: transformerReshapeParseCLI,
	IgnoresInput: false,
}

func transformerReshapeUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	argv0 := "mlr"
	verb := verbNameReshape

	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)

	fmt.Fprintf(o, "Wide-to-long options:\n")
	fmt.Fprintf(o, "  -i {input field names}   -o {key-field name,value-field name}\n")
	fmt.Fprintf(o, "  -r {input field regexes} -o {key-field name,value-field name}\n")
	fmt.Fprintf(o, "  These pivot/reshape the input data such that the input fields are removed\n")
	fmt.Fprintf(o, "  and separate records are emitted for each key/value pair.\n")
	fmt.Fprintf(o, "  Note: this works with tail -f and produces output records for each input\n")
	fmt.Fprintf(o, "  record seen.\n")
	fmt.Fprintf(o, "Long-to-wide options:\n")
	fmt.Fprintf(o, "  -s {key-field name,value-field name}\n")
	fmt.Fprintf(o, "  These pivot/reshape the input data to undo the wide-to-long operation.\n")
	fmt.Fprintf(o, "  Note: this does not work with tail -f; it produces output records only after\n")
	fmt.Fprintf(o, "  all input records have been read.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "  Input file \"wide.txt\":\n")
	fmt.Fprintf(o, "    time       X           Y\n")
	fmt.Fprintf(o, "    2009-01-01 0.65473572  2.4520609\n")
	fmt.Fprintf(o, "    2009-01-02 -0.89248112 0.2154713\n")
	fmt.Fprintf(o, "    2009-01-03 0.98012375  1.3179287\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "  %s --pprint %s -i X,Y -o item,value wide.txt\n", argv0, verb)
	fmt.Fprintf(o, "    time       item value\n")
	fmt.Fprintf(o, "    2009-01-01 X    0.65473572\n")
	fmt.Fprintf(o, "    2009-01-01 Y    2.4520609\n")
	fmt.Fprintf(o, "    2009-01-02 X    -0.89248112\n")
	fmt.Fprintf(o, "    2009-01-02 Y    0.2154713\n")
	fmt.Fprintf(o, "    2009-01-03 X    0.98012375\n")
	fmt.Fprintf(o, "    2009-01-03 Y    1.3179287\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "  %s --pprint %s -r '[A-Z]' -o item,value wide.txt\n", argv0, verb)
	fmt.Fprintf(o, "    time       item value\n")
	fmt.Fprintf(o, "    2009-01-01 X    0.65473572\n")
	fmt.Fprintf(o, "    2009-01-01 Y    2.4520609\n")
	fmt.Fprintf(o, "    2009-01-02 X    -0.89248112\n")
	fmt.Fprintf(o, "    2009-01-02 Y    0.2154713\n")
	fmt.Fprintf(o, "    2009-01-03 X    0.98012375\n")
	fmt.Fprintf(o, "    2009-01-03 Y    1.3179287\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "  Input file \"long.txt\":\n")
	fmt.Fprintf(o, "    time       item value\n")
	fmt.Fprintf(o, "    2009-01-01 X    0.65473572\n")
	fmt.Fprintf(o, "    2009-01-01 Y    2.4520609\n")
	fmt.Fprintf(o, "    2009-01-02 X    -0.89248112\n")
	fmt.Fprintf(o, "    2009-01-02 Y    0.2154713\n")
	fmt.Fprintf(o, "    2009-01-03 X    0.98012375\n")
	fmt.Fprintf(o, "    2009-01-03 Y    1.3179287\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "  %s --pprint %s -s item,value long.txt\n", argv0, verb)
	fmt.Fprintf(o, "    time       X           Y\n")
	fmt.Fprintf(o, "    2009-01-01 0.65473572  2.4520609\n")
	fmt.Fprintf(o, "    2009-01-02 -0.89248112 0.2154713\n")
	fmt.Fprintf(o, "    2009-01-03 0.98012375  1.3179287\n")
	fmt.Fprintf(o, "See also %s nest.\n", argv0)

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerReshapeParseCLI(
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
	var inputFieldNames []string = nil
	var inputFieldRegexStrings []string = nil
	var outputFieldNames []string = nil
	var splitOutFieldNames []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerReshapeUsage(os.Stdout, true, 0)

		} else if opt == "-i" {
			inputFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
		} else if opt == "-r" {
			inputFieldRegexStrings = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
		} else if opt == "-o" {
			outputFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
		} else if opt == "-s" {
			splitOutFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerReshapeUsage(os.Stderr, true, 1)
		}
	}

	outputKeyFieldName := ""
	outputValueFieldName := ""
	splitOutKeyFieldName := ""
	splitOutValueFieldName := ""

	if splitOutFieldNames == nil {
		// wide to long
		if inputFieldNames == nil && inputFieldRegexStrings == nil {
			transformerReshapeUsage(os.Stderr, true, 1)
		}

		if outputFieldNames == nil {
			transformerReshapeUsage(os.Stderr, true, 1)
		}
		if len(outputFieldNames) != 2 {
			transformerReshapeUsage(os.Stderr, true, 1)
		}
		outputKeyFieldName = outputFieldNames[0]
		outputValueFieldName = outputFieldNames[1]
	} else {
		// long to wide
		if len(splitOutFieldNames) != 2 {
			transformerReshapeUsage(os.Stderr, true, 1)
		}
		splitOutKeyFieldName = splitOutFieldNames[0]
		splitOutValueFieldName = splitOutFieldNames[1]
	}

	transformer, err := NewTransformerReshape(
		inputFieldNames,
		inputFieldRegexStrings,
		outputKeyFieldName,
		outputValueFieldName,
		splitOutKeyFieldName,
		splitOutValueFieldName,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerReshape struct {
	// for wide-to-long:
	inputFieldNames      []string
	inputFieldRegexes    []*regexp.Regexp
	outputKeyFieldName   string
	outputValueFieldName string

	// for long-to-wide:
	splitOutKeyFieldName            string
	splitOutValueFieldName          string
	otherKeysToOtherValuesToBuckets *lib.OrderedMap

	recordTransformerFunc RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerReshape(
	inputFieldNames []string,
	inputFieldRegexStrings []string,
	outputKeyFieldName string,
	outputValueFieldName string,
	splitOutKeyFieldName string,
	splitOutValueFieldName string,
) (*TransformerReshape, error) {

	tr := &TransformerReshape{
		inputFieldNames:      inputFieldNames,
		outputKeyFieldName:   outputKeyFieldName,
		outputValueFieldName: outputValueFieldName,

		splitOutKeyFieldName:            splitOutKeyFieldName,
		splitOutValueFieldName:          splitOutValueFieldName,
		otherKeysToOtherValuesToBuckets: lib.NewOrderedMap(),
	}

	if inputFieldRegexStrings != nil {
		tr.inputFieldRegexes = make([]*regexp.Regexp, len(inputFieldRegexStrings))
		// TODO: make a library function for string-array to regex-array
		for i, inputFieldRegexString := range inputFieldRegexStrings {
			regex, err := lib.CompileMillerRegex(inputFieldRegexString)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"%s %s: cannot compile regex [%s]\n",
					"mlr", verbNameReshape, inputFieldRegexString,
				)
				os.Exit(1)
			}
			tr.inputFieldRegexes[i] = regex
		}
	}

	if splitOutKeyFieldName == "" {
		if tr.inputFieldRegexes == nil {
			tr.recordTransformerFunc = tr.wideToLongNoRegex
		} else {
			tr.recordTransformerFunc = tr.wideToLongRegex
		}
	} else {
		tr.recordTransformerFunc = tr.longToWide
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerReshape) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerReshape) wideToLongNoRegex(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		pairs := types.NewMlrmap()
		for _, inputFieldName := range tr.inputFieldNames {
			value := inrec.Get(inputFieldName)
			if value != nil {
				// Reference, not copy, since the inrec will be freed here, or everything else will
				pairs.PutReference(inputFieldName, value)
			}
		}

		// Unset the record keys after iterating over them, rather than during
		for pe := pairs.Head; pe != nil; pe = pe.Next {
			inrec.Remove(pe.Key)
		}

		if pairs.IsEmpty() {
			outputChannel <- inrecAndContext
		} else {
			for pf := pairs.Head; pf != nil; pf = pf.Next {
				outrec := inrec.Copy()
				outrec.PutReference(tr.outputKeyFieldName, types.MlrvalFromString(pf.Key))
				outrec.PutReference(tr.outputValueFieldName, pf.Value)
				outputChannel <- types.NewRecordAndContext(outrec, &inrecAndContext.Context)
			}
		}

	} else {
		outputChannel <- inrecAndContext // emit end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerReshape) wideToLongRegex(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		pairs := types.NewMlrmap()

		for pd := inrec.Head; pd != nil; pd = pd.Next {
			for _, inputFieldRegex := range tr.inputFieldRegexes {
				if inputFieldRegex.MatchString(pd.Key) {
					// Reference, not copy, since the inrec will be freed here, or everything else will
					pairs.PutReference(pd.Key, pd.Value)
					break
				}
			}
		}

		// Unset the record keys after iterating over them, rather than during
		for pe := pairs.Head; pe != nil; pe = pe.Next {
			inrec.Remove(pe.Key)
		}

		if pairs.IsEmpty() {
			outputChannel <- inrecAndContext
		} else {
			for pf := pairs.Head; pf != nil; pf = pf.Next {
				outrec := inrec.Copy()
				outrec.PutReference(tr.outputKeyFieldName, types.MlrvalFromString(pf.Key))
				outrec.PutReference(tr.outputValueFieldName, pf.Value)
				outputChannel <- types.NewRecordAndContext(outrec, &inrecAndContext.Context)
			}
		}

	} else {
		outputChannel <- inrecAndContext // emit end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerReshape) longToWide(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		splitOutKeyFieldValue := inrec.Get(tr.splitOutKeyFieldName)
		splitOutValueFieldValue := inrec.Get(tr.splitOutValueFieldName)
		if splitOutKeyFieldValue == nil || splitOutValueFieldValue == nil {
			outputChannel <- inrecAndContext
			return
		}

		inrec.Remove(tr.splitOutKeyFieldName)
		inrec.Remove(tr.splitOutValueFieldName)

		// Don't unset tr.fieldName in the record, so we can implode in-place at the end.
		otherKeysJoined := inrec.GetKeysJoined()
		var otherValuesToBuckets *lib.OrderedMap = nil

		iOtherValuesToBuckets := tr.otherKeysToOtherValuesToBuckets.Get(otherKeysJoined)
		if iOtherValuesToBuckets == nil {
			otherValuesToBuckets = lib.NewOrderedMap()
			tr.otherKeysToOtherValuesToBuckets.Put(otherKeysJoined, otherValuesToBuckets)
		} else {
			otherValuesToBuckets = iOtherValuesToBuckets.(*lib.OrderedMap)
		}

		otherValuesJoined := inrec.GetValuesJoined()
		var bucket *tReshapeBucket = nil
		iBucket := otherValuesToBuckets.Get(otherValuesJoined)
		if iBucket == nil {
			bucket = newReshapeBucket(inrec)
			otherValuesToBuckets.Put(otherValuesJoined, bucket)
		} else {
			bucket = iBucket.(*tReshapeBucket)
		}

		bucket.pairs.PutCopy(splitOutKeyFieldValue.String(), splitOutValueFieldValue)

	} else {

		for pe := tr.otherKeysToOtherValuesToBuckets.Head; pe != nil; pe = pe.Next {
			otherValuesToBuckets := pe.Value.(*lib.OrderedMap)
			for pf := otherValuesToBuckets.Head; pf != nil; pf = pf.Next {
				bucket := pf.Value.(*tReshapeBucket)
				outrec := bucket.representative
				bucket.representative = nil // ownership transfer

				for pg := bucket.pairs.Head; pg != nil; pg = pg.Next {
					outrec.PutReference(pg.Key, pg.Value)
				}

				outputChannel <- types.NewRecordAndContext(outrec, &inrecAndContext.Context)
			}
		}

		outputChannel <- inrecAndContext // emit end-of-stream marker
	}
}

type tReshapeBucket struct {
	representative *types.Mlrmap
	pairs          *types.Mlrmap
}

func newReshapeBucket(representative *types.Mlrmap) *tReshapeBucket {
	return &tReshapeBucket{
		representative: representative,
		pairs:          types.NewMlrmap(),
	}
}
