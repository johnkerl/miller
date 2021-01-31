package transformers

import (
	"container/list"
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameBootstrap = "bootstrap"

var BootstrapSetup = transforming.TransformerSetup{
	Verb:         verbNameBootstrap,
	ParseCLIFunc: transformerBootstrapParseCLI,
	UsageFunc:    transformerBootstrapUsage,
	IgnoresInput: false,
}

func transformerBootstrapParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi] // xxx port more
	argi++

	nout := -1

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerBootstrapUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else if args[argi] == "-n" {
			nout = clitypes.VerbGetIntArgOrDie(verb, args, &argi, argc)

		} else {
			transformerBootstrapUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	transformer, _ := NewTransformerBootstrap(nout)

	*pargi = argi
	return transformer
}

func transformerBootstrapUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", os.Args[0], verbNameBootstrap)
	fmt.Fprintf(o,
		`Emits an n-sample, with replacement, of the input records.
See also %s sample and %s shuffle.
`, os.Args[0], os.Args[0])
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o,
		` -n Number of samples to output. Defaults to number of input records.
    Must be non-negative.
`)

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerBootstrap struct {
	recordsAndContexts *list.List
	nout               int
}

func NewTransformerBootstrap(nout int) (*TransformerBootstrap, error) {
	this := &TransformerBootstrap{
		recordsAndContexts: list.New(),
		nout:               nout,
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerBootstrap) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	// Not end of input stream: retain the record, and emit nothing until end of stream.
	if !inrecAndContext.EndOfStream {
		this.recordsAndContexts.PushBack(inrecAndContext)
		return
	}

	// Else end of record stream

	// Given nin input records, we produce nout output records, but
	// sampling with replacement.
	//
	// About memory management:
	//
	// Normally in Miller transformers we pass through pointers to records.
	// Here, though, since we do sampling with replacement, a record could
	// be emitted twice or more. To avoid producing multiple records in the
	// output stream pointing to the same memory, we would have to copy the
	// second one. In the original C (single-threaded) version of this
	// code, that was the case.
	//
	// However, in Go, there is concurrent processing.  It would be
	// possible for us to emit a pointer to a particular record without
	// copying, then when emitting that saem record a second time, copy it.
	// But due to concurrency, the pointed-to record could have already
	// been mutated downstream. We wouldn't be copying our input as we
	// received it -- we'd be copying something potentially modified.
	//
	// For that reason, this transformer must copy all output.

	// TODO: Go list Len() maxes at 2^31. We should track this ourselves in an int.
	nin := this.recordsAndContexts.Len()
	nout := this.nout
	if nout == -1 {
		nout = nin
	}

	if nout == 0 {
		// Emit the stream-terminating null record
		outputChannel <- inrecAndContext
		return
	}

	// Make an array of pointers into the input list.
	recordArray := make([]*types.RecordAndContext, nin)
	for i := 0; i < nin; i++ {
		head := this.recordsAndContexts.Front()
		if head == nil {
			break
		}
		recordArray[i] = head.Value.(*types.RecordAndContext)
		this.recordsAndContexts.Remove(head)
	}

	// Do the sample-with-replacment, reading from random indices in the input
	// array and emitting output.
	for i := 0; i < nout; i++ {
		index := lib.RandRange(0, nin)
		recordAndContext := recordArray[index]
		// Already emitted once; copy
		outputChannel <- recordAndContext.Copy()
	}

	// Emit the stream-terminating null record
	outputChannel <- inrecAndContext
}
