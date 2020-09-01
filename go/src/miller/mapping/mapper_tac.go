package mapping

import (
	"container/list"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/containers"
	"miller/runtime"
)

var MapperTacSetup = MapperSetup{
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
) IRecordMapper {
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
	lrecs *list.List
}

func NewMapperTac() (*MapperTac, error) {
	return &MapperTac{
		lrecs: list.New(),
	}, nil
}

func (this *MapperTac) Map(
	inrec *containers.Lrec,
	context *runtime.Context,
	outrecs chan<- *containers.Lrec,
) {
	if inrec != nil {
		this.lrecs.PushFront(inrec)
	} else {
		// end of stream
		for e := this.lrecs.Front(); e != nil; e = e.Next() {
			outrecs <- e.Value.(*containers.Lrec)
		}
		outrecs <- nil
	}
}
