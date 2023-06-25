package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameCheck = "check"

var CheckSetup = TransformerSetup{
	Verb:         verbNameCheck,
	UsageFunc:    transformerCheckUsage,
	ParseCLIFunc: transformerCheckParseCLI,
	IgnoresInput: false,
}

func transformerCheckUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameCheck)
	fmt.Fprintf(o, "Consumes records without printing any output,\n")
	fmt.Fprintf(o, "Useful for doing a well-formatted check on input data.\n")
	fmt.Fprintf(o, "with the exception that warnings are printed to stderr.\n")
	fmt.Fprintf(o, "Current checks are:\n")
	fmt.Fprintf(o, "* If any key is the empty string\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerCheckParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

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
			transformerCheckUsage(os.Stdout)
			os.Exit(0)

		} else {
			transformerCheckUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerCheck()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerCheck struct {
	// stateless
	messagedReEmptyKey map[string]bool
}

func NewTransformerCheck() (*TransformerCheck, error) {
	return &TransformerCheck{
		messagedReEmptyKey: make(map[string]bool),
	}, nil
}

func (tr *TransformerCheck) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if pe.Key == "" {
				context := inrecAndContext.Context

				// Most Miller users are CSV users. And for CSV this will be an error on
				// *every* record, or none -- so let's not print this multiple times.
				if tr.messagedReEmptyKey[context.FILENAME] {
					continue
				}

				message := fmt.Sprintf(
					"mlr: warning: empty-string key at filename %s record number %d",
					context.FILENAME, context.NR,
				)
				fmt.Fprintln(os.Stderr, message)
				tr.messagedReEmptyKey[context.FILENAME] = true
			}
		}
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}
