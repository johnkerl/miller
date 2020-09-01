package mappers

import (
	"container/list"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/mapping"
	"miller/containers"
)

var TacSetup = mapping.MapperSetup{
	Verb:         "tac",
	ParseCLIFunc: mapperTacParseCLIFunc,
	UsageFunc:    mapperTacUsageFunc,
	IgnoresInput: false,
}

func mapperTacParseCLIFunc(
	pargi *int,
	argc int,
	args []string,
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) mapping.IRecordMapper {
	*pargi += 1

	// xxx temp err keep or no
	mapper, _ := NewMapperTac()
	return mapper
}

func mapperTacUsageFunc(
	o *os.File,
	argv0 string,
	verb string,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Prints records in reverse order from the order in which they were encountered.\n")
}

// ----------------------------------------------------------------
type MapperTac struct {
	lrecsAndContexts *list.List
}

func NewMapperTac() (*MapperTac, error) {
	return &MapperTac{
		lrecsAndContexts: list.New(),
	}, nil
}

func (this *MapperTac) Map(
	inrecAndContext *containers.LrecAndContext,
	outrecsAndContexts chan<- *containers.LrecAndContext,
) {
	if inrecAndContext.Lrec != nil {
		this.lrecsAndContexts.PushFront(inrecAndContext)
	} else {
		// end of stream
		for e := this.lrecsAndContexts.Front(); e != nil; e = e.Next() {
			outrecsAndContexts <- e.Value.(*containers.LrecAndContext)
		}
		outrecsAndContexts <- containers.NewLrecAndContext(
			nil, // signals end of input record stream
			&inrecAndContext.Context,
		)
	}
}
