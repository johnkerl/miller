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
const verbNameShuffle = "shuffle"

var ShuffleSetup = transforming.TransformerSetup{
	Verb:         verbNameShuffle,
	ParseCLIFunc: transformerShuffleParseCLI,
	UsageFunc:    transformerShuffleUsage,
	IgnoresInput: false,
}

func transformerShuffleParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerShuffleUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else {
			transformerShuffleUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	transformer, _ := NewTransformerShuffle()

	*pargi = argi
	return transformer
}

func transformerShuffleUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s, with no options.\n", os.Args[0], verbNameShuffle)
	fmt.Fprintf(o,
		`Outputs records randomly permuted. No output records are produced until
all input records are read. See also %s bootstrap and %s sample.
`, os.Args[0], os.Args[0])

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerShuffle struct {
	recordsAndContexts *list.List
}

func NewTransformerShuffle() (*TransformerShuffle, error) {

	this := &TransformerShuffle{
		recordsAndContexts: list.New(),
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerShuffle) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	// Not end of input stream: retain the record, and emit nothing until end of stream.
	if !inrecAndContext.EndOfStream {
		this.recordsAndContexts.PushBack(inrecAndContext)

	} else { // end of record stream
		// Knuth shuffle:
		// * Initial permutation is identity.
		// * Make a pseudorandom permutation using pseudorandom swaps in the image map.
		// TODO: Go list Len() maxes at 2^31. We should track this ourselves in an int.
		n := this.recordsAndContexts.Len()
		images := make([]int, n)
		for i := 0; i < n; i++ {
			images[i] = i
		}
		unusedStart := 0
		numUnused := n
		for i := 0; i < n; i++ {
			// Select a pseudorandom element from the pool of unused images.
			u := lib.RandRange(unusedStart, unusedStart+numUnused)
			temp := images[u]
			images[u] = images[i]
			images[i] = temp
			// Decrease the size of the pool by 1.  (Yes, unusedStart and k always have the same value.
			// Using two variables wastes neglible memory and makes the code easier to understand.)
			unusedStart++
			numUnused--
		}

		// Move the record-pointers from linked list to array.
		array := make([]*types.RecordAndContext, n)
		for i := 0; i < n; i++ {
			head := this.recordsAndContexts.Front()
			if head == nil {
				break
			}
			array[i] = head.Value.(*types.RecordAndContext)
			this.recordsAndContexts.Remove(head)
		}

		// Transfer from input array to output list. Because permutations are one-to-one maps,
		// all input records have ownership transferred exactly once. So, there are no
		// records to copy here.
		for i := 0; i < n; i++ {
			outputChannel <- array[images[i]]
		}

		// Emit the stream-terminating null record
		outputChannel <- inrecAndContext
	}
}
