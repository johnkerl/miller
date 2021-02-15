package transformers

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameCut = "cut"

var CutSetup = transforming.TransformerSetup{
	Verb:         verbNameCut,
	UsageFunc:    transformerCutUsage,
	ParseCLIFunc: transformerCutParseCLI,
	IgnoresInput: false,
}

func transformerCutUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameCut)
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
	fmt.Fprintf(o, "  %s %s -f hostname,status\n", lib.MlrExeName(), verbNameCut)
	fmt.Fprintf(o, "  %s %s -x -f hostname,status\n", lib.MlrExeName(), verbNameCut)
	fmt.Fprintf(o, "  %s %s -r -f '^status$,sda[0-9]'\n", lib.MlrExeName(), verbNameCut)
	fmt.Fprintf(o, "  %s %s -r -f '^status$,\"sda[0-9]\"'\n", lib.MlrExeName(), verbNameCut)
	fmt.Fprintf(o, "  %s %s -r -f '^status$,\"sda[0-9]\"i' (this is case-insensitive)\n", lib.MlrExeName(), verbNameCut)

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerCutParseCLI(
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

	var fieldNames []string = nil
	doArgOrder := false
	doComplement := false
	doRegexes := false

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerCutUsage(os.Stdout, true, 0)

		} else if opt == "-f" {
			fieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-o" {
			doArgOrder = true

		} else if opt == "-x" {
			doComplement = true

		} else if opt == "--complement" {
			doComplement = true

		} else if opt == "-r" {
			doRegexes = true

		} else {
			transformerCutUsage(os.Stderr, true, 1)
		}
	}

	if fieldNames == nil {
		transformerCutUsage(os.Stderr, true, 1)
	}

	transformer, _ := NewTransformerCut(
		fieldNames,
		doArgOrder,
		doComplement,
		doRegexes,
	)

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerCut struct {
	fieldNameList []string
	fieldNameSet  map[string]bool

	doComplement bool
	regexes      []*regexp.Regexp

	recordTransformerFunc transforming.RecordTransformerFunc
}

func NewTransformerCut(
	fieldNames []string,
	doArgOrder bool,
	doComplement bool,
	doRegexes bool,
) (*TransformerCut, error) {

	this := &TransformerCut{}

	if !doRegexes {
		this.fieldNameList = fieldNames
		this.fieldNameSet = lib.StringListToSet(fieldNames)
		if !doComplement {
			if !doArgOrder {
				this.recordTransformerFunc = this.includeWithInputOrder
			} else {
				this.recordTransformerFunc = this.includeWithArgOrder
			}
		} else {
			this.recordTransformerFunc = this.exclude
		}
	} else {
		this.doComplement = doComplement
		this.regexes = make([]*regexp.Regexp, len(fieldNames))
		for i, regexString := range fieldNames {
			// Handles "a.*b"i Miller case-insensitive-regex specification
			regex, err := lib.CompileMillerRegex(regexString)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"%s %s: cannot compile regex [%s]\n",
					lib.MlrExeName(), verbNameCut, regexString,
				)
				os.Exit(1)
			}
			this.regexes[i] = regex
		}
		this.recordTransformerFunc = this.processWithRegexes
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerCut) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
// mlr cut -f a,b,c
func (this *TransformerCut) includeWithInputOrder(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		outrec := types.NewMlrmap()
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			fieldName := pe.Key
			_, wanted := this.fieldNameSet[fieldName]
			if wanted {
				outrec.PutReference(fieldName, pe.Value) // inrec will be GC'ed
			}
		}
		outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
		outputChannel <- outrecAndContext
	} else {
		outputChannel <- inrecAndContext
	}
}

// ----------------------------------------------------------------
// mlr cut -o -f a,b,c
func (this *TransformerCut) includeWithArgOrder(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		outrec := types.NewMlrmap()
		for _, fieldName := range this.fieldNameList {
			value := inrec.Get(fieldName)
			if value != nil {
				outrec.PutReference(fieldName, value) // inrec will be GC'ed
			}
		}
		outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
		outputChannel <- outrecAndContext
	} else {
		outputChannel <- inrecAndContext
	}
}

// ----------------------------------------------------------------
// mlr cut -x -f a,b,c
func (this *TransformerCut) exclude(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for _, fieldName := range this.fieldNameList {
			if inrec.Has(fieldName) {
				inrec.Remove(fieldName)
			}
		}
	}
	outputChannel <- inrecAndContext
}

// ----------------------------------------------------------------
func (this *TransformerCut) processWithRegexes(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		newrec := types.NewMlrmapAsRecord()
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			matchesAny := false
			for _, regex := range this.regexes {
				if regex.MatchString(pe.Key) {
					matchesAny = true
					break
				}
			}
			// Boolean XOR is spelt '!=' in Go
			if matchesAny != this.doComplement {
				// Pointer-motion is OK since the inrec is being hereby discarded.
				// We're simply transferring ownership to the newrec.
				newrec.PutReference(pe.Key, pe.Value)
			}
		}
		outputChannel <- types.NewRecordAndContext(newrec, &inrecAndContext.Context)
	} else {
		outputChannel <- inrecAndContext
	}
}
