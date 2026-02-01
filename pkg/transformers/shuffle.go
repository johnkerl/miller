package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/types"
)

// ----------------------------------------------------------------
const verbNameShuffle = "shuffle"

var ShuffleSetup = TransformerSetup{
	Verb:         verbNameShuffle,
	UsageFunc:    transformerShuffleUsage,
	ParseCLIFunc: transformerShuffleParseCLI,
	IgnoresInput: false,
}

func transformerShuffleUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameShuffle)
	fmt.Fprintf(o, "Outputs records randomly permuted. No output records are produced until\n")
	fmt.Fprintf(o, "all input records are read. See also %s bootstrap and %s sample.\n",
		"mlr", "mlr",
	)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerShuffleParseCLI(
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
			transformerShuffleUsage(os.Stdout)
			os.Exit(0)

		} else {
			transformerShuffleUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerShuffle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerShuffle struct {
	recordsAndContexts []*types.RecordAndContext
}

func NewTransformerShuffle() (*TransformerShuffle, error) {

	tr := &TransformerShuffle{
		recordsAndContexts: make([]*types.RecordAndContext, 0),
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerShuffle) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	// Not end of input stream: retain the record, and emit nothing until end of stream.
	if !inrecAndContext.EndOfStream {
		tr.recordsAndContexts = append(tr.recordsAndContexts, inrecAndContext)

	} else { // end of record stream
		// Knuth shuffle:
		// * Initial permutation is identity.
		// * Make a pseudorandom permutation using pseudorandom swaps in the image map.
		n := int64(len(tr.recordsAndContexts))
		images := make([]int64, n)
		for i := int64(0); i < n; i++ {
			images[i] = i
		}
		unusedStart := int64(0)
		numUnused := n
		for i := int64(0); i < n; i++ {
			// Select a pseudorandom element from the pool of unused images.
			u := lib.RandRange(unusedStart, unusedStart+numUnused)
			temp := images[u]
			images[u] = images[i]
			images[i] = temp
			// Decrease the size of the pool by 1.  (Yes, unusedStart and k always have the same value.
			// Using two variables wastes negligible memory and makes the code easier to understand.)
			unusedStart++
			numUnused--
		}

		// Move the record-pointers from slice to array.
		array := tr.recordsAndContexts

		// Transfer from input array to output list. Because permutations are one-to-one maps,
		// all input records have ownership transferred exactly once. So, there are no
		// records to copy here.
		for i := int64(0); i < n; i++ {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, array[images[i]])
		}

		// Emit the stream-terminating null record
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
	}
}
