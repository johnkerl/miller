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
const verbNameBootstrap = "bootstrap"

var BootstrapSetup = TransformerSetup{
	Verb:         verbNameBootstrap,
	UsageFunc:    transformerBootstrapUsage,
	ParseCLIFunc: transformerBootstrapParseCLI,
	IgnoresInput: false,
}

func transformerBootstrapUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameBootstrap)
	fmt.Fprintf(o,
		`Emits an n-sample, with replacement, of the input records.
See also %s sample and %s shuffle.
`, "mlr", "mlr")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o,
		` -n Number of samples to output. Defaults to number of input records.
    Must be non-negative.
`)
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerBootstrapParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	nout := -1

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerBootstrapUsage(os.Stdout, true, 0)

		} else if opt == "-n" {
			nout = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerBootstrapUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerBootstrap(nout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerBootstrap struct {
	recordsAndContexts *list.List
	nout               int
}

func NewTransformerBootstrap(nout int) (*TransformerBootstrap, error) {
	tr := &TransformerBootstrap{
		recordsAndContexts: list.New(),
		nout:               nout,
	}
	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerBootstrap) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	// Not end of input stream: retain the record, and emit nothing until end of stream.
	if !inrecAndContext.EndOfStream {
		tr.recordsAndContexts.PushBack(inrecAndContext)
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
	nin := tr.recordsAndContexts.Len()
	nout := tr.nout
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
		head := tr.recordsAndContexts.Front()
		if head == nil {
			break
		}
		recordArray[i] = head.Value.(*types.RecordAndContext)
		tr.recordsAndContexts.Remove(head)
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
