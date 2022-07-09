package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameUnsparsify = "unsparsify"

var UnsparsifySetup = TransformerSetup{
	Verb:         verbNameUnsparsify,
	UsageFunc:    transformerUnsparsifyUsage,
	ParseCLIFunc: transformerUnsparsifyParseCLI,
	IgnoresInput: false,
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

	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "--fill-with {filler string}  What to fill absent fields with. Defaults to\n")
	fmt.Fprintf(o, "                             the empty string.\n")
	fmt.Fprintf(o, "-f {a,b,c} Specify field names to be operated on. Any other fields won't be\n")
	fmt.Fprintf(o, "           modified, and operation will be streaming.\n")
	fmt.Fprintf(o, "-h|--help  Show this message.\n")

	fmt.Fprint(o,
		`Example: if the input is two records, one being 'a=1,b=2' and the other
being 'b=3,c=4', then the output is the two records 'a=1,b=2,c=' and
'a=,b=3,c=4'.
`)
}

func transformerUnsparsifyParseCLI(
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

	fillerString := ""
	var specifiedFieldNames []string = nil

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
			transformerUnsparsifyUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "--fill-with" {
			fillerString = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			specifiedFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerUnsparsifyUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerUnsparsify(
		fillerString,
		specifiedFieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerUnsparsify struct {
	fillerMlrval          *mlrval.Mlrval
	recordsAndContexts    *list.List
	fieldNamesSeen        *lib.OrderedMap
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerUnsparsify(
	fillerString string,
	specifiedFieldNames []string,
) (*TransformerUnsparsify, error) {

	fieldNamesSeen := lib.NewOrderedMap()
	for _, specifiedFieldName := range specifiedFieldNames {
		fieldNamesSeen.Put(specifiedFieldName, specifiedFieldName)
	}

	tr := &TransformerUnsparsify{
		fillerMlrval:       mlrval.FromString(fillerString),
		recordsAndContexts: list.New(),
		fieldNamesSeen:     fieldNamesSeen,
	}

	if specifiedFieldNames == nil {
		tr.recordTransformerFunc = tr.transformNonStreaming
	} else {
		tr.recordTransformerFunc = tr.transformStreaming
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerUnsparsify) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerUnsparsify) transformNonStreaming(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			key := pe.Key
			if !tr.fieldNamesSeen.Has(key) {
				tr.fieldNamesSeen.Put(key, key)
			}
		}
		tr.recordsAndContexts.PushBack(inrecAndContext)
	} else {
		for e := tr.recordsAndContexts.Front(); e != nil; e = e.Next() {
			outrecAndContext := e.Value.(*types.RecordAndContext)
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

			outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &outrecAndContext.Context))
		}

		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerUnsparsify) transformStreaming(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for pe := tr.fieldNamesSeen.Head; pe != nil; pe = pe.Next {
			if !inrec.Has(pe.Key) {
				inrec.PutCopy(pe.Key, tr.fillerMlrval)
			}
		}

		outputRecordsAndContexts.PushBack(inrecAndContext)

	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}
