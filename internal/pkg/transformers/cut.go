package transformers

import (
	"container/list"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameCut = "cut"

var CutSetup = TransformerSetup{
	Verb:         verbNameCut,
	UsageFunc:    transformerCutUsage,
	ParseCLIFunc: transformerCutParseCLI,
	IgnoresInput: false,
}

func transformerCutUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameCut)
	fmt.Fprintf(o, "Passes through input records with specified fields included/excluded.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, " -f {a,b,c} Comma-separated field names for cut, e.g. a,b,c.\n")
	fmt.Fprintf(o, " -o Retain fields in the order specified here in the argument list.\n")
	fmt.Fprintf(o, "    Default is to retain them in the order found in the input data.\n")
	fmt.Fprintf(o, " -x|--complement  Exclude, rather than include, field names specified by -f.\n")
	fmt.Fprintf(o, " -r Treat field names as regular expressions. \"ab\", \"a.*b\" will\n")
	fmt.Fprintf(o, "   match any field name containing the substring \"ab\" or matching\n")
	fmt.Fprintf(o, "   \"a.*b\", respectively; anchors of the form \"^ab$\", \"^a.*b$\" may\n")
	fmt.Fprintf(o, "   be used. The -o flag is ignored when -r is present.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(o, "  %s %s -f hostname,status\n", "mlr", verbNameCut)
	fmt.Fprintf(o, "  %s %s -x -f hostname,status\n", "mlr", verbNameCut)
	fmt.Fprintf(o, "  %s %s -r -f '^status$,sda[0-9]'\n", "mlr", verbNameCut)
	fmt.Fprintf(o, "  %s %s -r -f '^status$,\"sda[0-9]\"'\n", "mlr", verbNameCut)
	fmt.Fprintf(o, "  %s %s -r -f '^status$,\"sda[0-9]\"i' (this is case-insensitive)\n", "mlr", verbNameCut)
}

func transformerCutParseCLI(
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

	var fieldNames []string = nil
	doArgOrder := false
	doComplement := false
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
			transformerCutUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-o" {
			doArgOrder = true

		} else if opt == "-x" {
			doComplement = true

		} else if opt == "--complement" {
			doComplement = true

		} else if opt == "-r" {
			doRegexes = true

		} else {
			transformerCutUsage(os.Stderr)
			os.Exit(1)
		}
	}

	if fieldNames == nil {
		transformerCutUsage(os.Stderr)
		os.Exit(1)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerCut(
		fieldNames,
		doArgOrder,
		doComplement,
		doRegexes,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerCut struct {
	fieldNameList []string
	fieldNameSet  map[string]bool

	doComplement bool
	regexes      []*regexp.Regexp

	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerCut(
	fieldNames []string,
	doArgOrder bool,
	doComplement bool,
	doRegexes bool,
) (*TransformerCut, error) {

	tr := &TransformerCut{}

	if !doRegexes {
		tr.fieldNameList = fieldNames
		tr.fieldNameSet = lib.StringListToSet(fieldNames)
		if !doComplement {
			if !doArgOrder {
				tr.recordTransformerFunc = tr.includeWithInputOrder
			} else {
				tr.recordTransformerFunc = tr.includeWithArgOrder
			}
		} else {
			tr.recordTransformerFunc = tr.exclude
		}
	} else {
		tr.doComplement = doComplement
		tr.regexes = make([]*regexp.Regexp, len(fieldNames))
		for i, regexString := range fieldNames {
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
			tr.regexes[i] = regex
		}
		tr.recordTransformerFunc = tr.processWithRegexes
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerCut) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

// ----------------------------------------------------------------
// mlr cut -f a,b,c
func (tr *TransformerCut) includeWithInputOrder(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		outrec := mlrval.NewMlrmap()
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			fieldName := pe.Key
			_, wanted := tr.fieldNameSet[fieldName]
			if wanted {
				outrec.PutReference(fieldName, pe.Value) // inrec will be GC'ed
			}
		}
		outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
		outputRecordsAndContexts.PushBack(outrecAndContext)
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}

// ----------------------------------------------------------------
// mlr cut -o -f a,b,c
func (tr *TransformerCut) includeWithArgOrder(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		outrec := mlrval.NewMlrmap()
		for _, fieldName := range tr.fieldNameList {
			value := inrec.Get(fieldName)
			if value != nil {
				outrec.PutReference(fieldName, value) // inrec will be GC'ed
			}
		}
		outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
		outputRecordsAndContexts.PushBack(outrecAndContext)
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}

// ----------------------------------------------------------------
// mlr cut -x -f a,b,c
func (tr *TransformerCut) exclude(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for _, fieldName := range tr.fieldNameList {
			if inrec.Has(fieldName) {
				inrec.Remove(fieldName)
			}
		}
	}
	outputRecordsAndContexts.PushBack(inrecAndContext)
}

// ----------------------------------------------------------------
func (tr *TransformerCut) processWithRegexes(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		newrec := mlrval.NewMlrmapAsRecord()
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			matchesAny := false
			for _, regex := range tr.regexes {
				if regex.MatchString(pe.Key) {
					matchesAny = true
					break
				}
			}
			// Boolean XOR is spelt '!=' in Go
			if matchesAny != tr.doComplement {
				// Pointer-motion is OK since the inrec is being hereby discarded.
				// We're simply transferring ownership to the newrec.
				newrec.PutReference(pe.Key, pe.Value)
			}
		}
		outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}
