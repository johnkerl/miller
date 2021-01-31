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
const verbNameGroupLike = "group-like"

var GroupLikeSetup = transforming.TransformerSetup{
	Verb:         verbNameGroupLike,
	ParseCLIFunc: transformerGroupLikeParseCLI,
	UsageFunc:    transformerGroupLikeUsage,
	IgnoresInput: false,
}

func transformerGroupLikeParseCLI(
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
			transformerGroupLikeUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else {
			transformerGroupLikeUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	transformer, _ := NewTransformerGroupLike()

	*pargi = argi
	return transformer
}

func transformerGroupLikeUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s, with no options.\n", os.Args[0], verbNameGroupLike)
	fmt.Fprint(o,
		`Outputs records in batches having identical field names.
`)

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerGroupLike struct {
	// map from string to *list.List
	recordListsByGroup *lib.OrderedMap
}

func NewTransformerGroupLike() (*TransformerGroupLike, error) {

	this := &TransformerGroupLike{
		recordListsByGroup: lib.NewOrderedMap(),
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerGroupLike) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey := inrec.GetKeysJoined()

		recordListForGroup := this.recordListsByGroup.Get(groupingKey)
		if recordListForGroup == nil { // first time
			recordListForGroup = list.New()
			this.recordListsByGroup.Put(groupingKey, recordListForGroup)
		}

		recordListForGroup.(*list.List).PushBack(inrecAndContext)

	} else {
		for outer := this.recordListsByGroup.Head; outer != nil; outer = outer.Next {
			recordListForGroup := outer.Value.(*list.List)
			for inner := recordListForGroup.Front(); inner != nil; inner = inner.Next() {
				outputChannel <- inner.Value.(*types.RecordAndContext)
			}
		}
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
