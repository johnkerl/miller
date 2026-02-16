package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameGroupLike = "group-like"

var GroupLikeSetup = TransformerSetup{
	Verb:         verbNameGroupLike,
	UsageFunc:    transformerGroupLikeUsage,
	ParseCLIFunc: transformerGroupLikeParseCLI,
	IgnoresInput: false,
}

func transformerGroupLikeUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameGroupLike)
	fmt.Fprintln(o, "Outputs records in batches having identical field names.")
	fmt.Fprintln(o, "Options:")
	fmt.Fprintln(o, "-h|--help Show this message.")
}

func transformerGroupLikeParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

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
			transformerGroupLikeUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		}
		return nil, cli.VerbErrorf(verbNameGroupLike, "option \"%s\" not recognized", opt)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerGroupLike()
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerGroupLike struct {
	// map from string to record slices
	recordListsByGroup *lib.OrderedMap[*[]*types.RecordAndContext]
}

func NewTransformerGroupLike() (*TransformerGroupLike, error) {

	tr := &TransformerGroupLike{
		recordListsByGroup: lib.NewOrderedMap[*[]*types.RecordAndContext](),
	}

	return tr, nil
}

func (tr *TransformerGroupLike) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey := inrec.GetKeysJoined()

		recordListForGroup := tr.recordListsByGroup.Get(groupingKey)
		if recordListForGroup == nil { // first time
			records := []*types.RecordAndContext{}
			recordListForGroup = &records
			tr.recordListsByGroup.Put(groupingKey, recordListForGroup)
		}

		*recordListForGroup = append(*recordListForGroup, inrecAndContext)

	} else {
		for outer := tr.recordListsByGroup.Head; outer != nil; outer = outer.Next {
			recordListForGroup := outer.Value
			for _, recordAndContext := range *recordListForGroup {
				*outputRecordsAndContexts = append(*outputRecordsAndContexts, recordAndContext)
			}
		}
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
	}
}
