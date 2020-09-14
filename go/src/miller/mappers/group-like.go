package mappers

import (
	"container/list"
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/mapping"
	"miller/types"
)

// ----------------------------------------------------------------
var GroupLikeSetup = mapping.MapperSetup{
	Verb:         "group-like",
	ParseCLIFunc: mapperGroupLikeParseCLI,
	IgnoresInput: false,
}

func mapperGroupLikeParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) mapping.IRecordMapper {

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperGroupLikeUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentioally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperGroupLike()

	*pargi = argi
	return mapper
}

func mapperGroupLikeUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Outputs records in batches having identical field names.
`)
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperGroupLike struct {
	recordListsByGroup map[string]*list.List
}

func NewMapperGroupLike() (*MapperGroupLike, error) {

	this := &MapperGroupLike{
		recordListsByGroup: make(map[string]*list.List),
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperGroupLike) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream

		groupByKey := inrec.GetKeysJoined()

		recordListForGroup, present := this.recordListsByGroup[groupByKey]
		if !present { // first time
			recordListForGroup = list.New()
			this.recordListsByGroup[groupByKey] = recordListForGroup
		}

		recordListForGroup.PushBack(inrecAndContext)

	} else {
		for _, recordListForGroup := range this.recordListsByGroup {
			for entry := recordListForGroup.Front(); entry != nil; entry = entry.Next() {
				outputChannel <- entry.Value.(*types.RecordAndContext)
			}
		}
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
