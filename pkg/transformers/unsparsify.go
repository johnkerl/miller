package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameUnsparsify = "unsparsify"

var unsparsifyOptions = []OptionSpec{
	{Flag: "--fill-with", Arg: "{filler string}", Type: "string", Desc: "What to fill absent fields with. Defaults to the empty string."},
	{Flag: "-f", Arg: "{a,b,c}", Type: "csv-list", Desc: "Specify field names to be operated on; others are not modified and operation is streaming."},
	{Flag: "--from-header", Type: "bool", Desc: "Operate on the field names from the header line of the current input file; others are not modified and operation is streaming. Use this with --ragged to fill in the trailing fields of data lines which are shorter than the header line. Applies to CSV, TSV, CSV-lite, PPRINT, and markdown input, and does nothing for input formats without a header line. Can be used with -f."},
}

var UnsparsifySetup = TransformerSetup{
	Verb:         verbNameUnsparsify,
	UsageFunc:    transformerUnsparsifyUsage,
	ParseCLIFunc: transformerUnsparsifyParseCLI,
	IgnoresInput: false,
	Options:      unsparsifyOptions,
}

func transformerUnsparsifyUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameUnsparsify)
	fmt.Fprint(o,
		`Prints records with the union of field names over all input records.
For field names absent in a given record but present in others, fills in
a value. This verb retains all input before producing any output.
`)

	WriteVerbOptions(o, unsparsifyOptions)

	fmt.Fprint(o,
		`Example: if the input is two records, one being 'a=1,b=2' and the other
being 'b=3,c=4', then the output is the two records 'a=1,b=2,c=' and
'a=,b=3,c=4'.
Example: if the CSV input has header line 'a,b,c' and data lines '1,2' and
'4,5,6', then 'mlr --csv --ragged unsparsify --from-header' gives 'a,b,c',
'1,2,', and '4,5,6'.
`)
}

func transformerUnsparsifyParseCLI(
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

	fillerString := ""
	var specifiedFieldNames []string = nil
	fromHeader := false

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
			transformerUnsparsifyUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		case "--fill-with":
			fillerString, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "-f":
			specifiedFieldNames, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "--from-header":
			fromHeader = true

		default:
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerUnsparsify(
		fillerString,
		specifiedFieldNames,
		fromHeader,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerUnsparsify struct {
	fillerMlrval          *mlrval.Mlrval
	recordsAndContexts    []*types.RecordAndContext
	fieldNamesSeen        *lib.OrderedMap[string]
	fromHeader            bool
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerUnsparsify(
	fillerString string,
	specifiedFieldNames []string,
	fromHeader bool,
) (*TransformerUnsparsify, error) {

	fieldNamesSeen := lib.NewOrderedMap[string]()
	for _, specifiedFieldName := range specifiedFieldNames {
		fieldNamesSeen.Put(specifiedFieldName, specifiedFieldName)
	}

	tr := &TransformerUnsparsify{
		fillerMlrval:       mlrval.FromString(fillerString),
		recordsAndContexts: []*types.RecordAndContext{},
		fieldNamesSeen:     fieldNamesSeen,
		fromHeader:         fromHeader,
	}

	if specifiedFieldNames == nil && !fromHeader {
		tr.recordTransformerFunc = tr.transformNonStreaming
	} else {
		tr.recordTransformerFunc = tr.transformStreaming
	}

	return tr, nil
}

func (tr *TransformerUnsparsify) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) error {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	return tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

func (tr *TransformerUnsparsify) transformNonStreaming(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) error {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			key := pe.Key
			if !tr.fieldNamesSeen.Has(key) {
				tr.fieldNamesSeen.Put(key, key)
			}
		}
		tr.recordsAndContexts = append(tr.recordsAndContexts, inrecAndContext)
	} else {
		for _, outrecAndContext := range tr.recordsAndContexts {
			outrec := outrecAndContext.Record

			newrec := mlrval.NewMlrmapAsRecord()
			for pe := tr.fieldNamesSeen.Head; pe != nil; pe = pe.Next {
				fieldName := pe.Key
				if !outrec.Has(fieldName) {
					newrec.PutCopy(fieldName, tr.fillerMlrval)
				} else {
					newrec.PutReference(fieldName, outrec.Get(fieldName))
				}
			}

			*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(newrec, &outrecAndContext.Context))
		}

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
	}
	return nil
}

func (tr *TransformerUnsparsify) transformStreaming(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) error {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		// With --ragged, a CSV data line shorter than the header line yields a
		// record without the trailing keys. The reader puts the header in the
		// context so we can fill them in here, without retaining records.
		if tr.fromHeader {
			for _, key := range inrecAndContext.Context.HeaderFieldNames {
				if !inrec.Has(key) {
					inrec.PutCopy(key, tr.fillerMlrval)
				}
			}
		}

		for pe := tr.fieldNamesSeen.Head; pe != nil; pe = pe.Next {
			if !inrec.Has(pe.Key) {
				inrec.PutCopy(pe.Key, tr.fillerMlrval)
			}
		}

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
	}
	return nil
}
