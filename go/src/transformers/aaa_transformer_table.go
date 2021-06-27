package transformers

import (
	"fmt"
	"os"

	"miller/src/lib"
)

// ----------------------------------------------------------------
var TRANSFORMER_LOOKUP_TABLE = []TransformerSetup{
	AltkvSetup,
	BarSetup,
	BootstrapSetup,
	CatSetup,
	CheckSetup,
	CleanWhitespaceSetup,
	CountDistinctSetup,
	CountSetup,
	CountSimilarSetup,
	CutSetup,
	DecimateSetup,
	FillDownSetup,
	FillEmptySetup,
	FilterSetup,
	FlattenSetup,
	FormatValuesSetup,
	FractionSetup,
	GapSetup,
	GrepSetup,
	GroupBySetup,
	GroupLikeSetup,
	HavingFieldsSetup,
	HeadSetup,
	HistogramSetup,
	JSONParseSetup,
	JSONStringifySetup,
	JoinSetup,
	LabelSetup,
	LeastFrequentSetup,
	MergeFieldsSetup,
	MostFrequentSetup,
	NestSetup,
	NothingSetup,
	PutSetup,
	RegularizeSetup,
	RemoveEmptyColumnsSetup,
	RenameSetup,
	ReorderSetup,
	RepeatSetup,
	ReshapeSetup,
	SampleSetup,
	Sec2GMTDateSetup,
	Sec2GMTSetup,
	SeqgenSetup,
	ShuffleSetup,
	SkipTrivialRecordsSetup,
	SortSetup,
	SortWithinRecordsSetup,
	Stats1Setup,
	Stats2Setup,
	StepSetup,
	TacSetup,
	TailSetup,
	TeeSetup,
	TopSetup,
	UnflattenSetup,
	UniqSetup,
	UnsparsifySetup,
}

// ----------------------------------------------------------------
func LookUp(verb string) *TransformerSetup {
	for _, transformerSetup := range TRANSFORMER_LOOKUP_TABLE {
		if transformerSetup.Verb == verb {
			return &transformerSetup
		}
	}
	return nil
}

// ----------------------------------------------------------------
func ListAllVerbNames(o *os.File) {
	for _, transformerSetup := range TRANSFORMER_LOOKUP_TABLE {
		fmt.Fprintf(o, "%s\n", transformerSetup.Verb)
	}
}

// ----------------------------------------------------------------
func ListAllVerbs(o *os.File, leader string) {
	separator := " "

	leaderlen := len(leader)
	separatorlen := len(separator)
	linelen := leaderlen
	j := 0

	for _, transformerSetup := range TRANSFORMER_LOOKUP_TABLE {
		verb := transformerSetup.Verb
		verblen := len(verb)
		linelen += separatorlen + verblen
		if linelen >= 80 {
			fmt.Fprintf(o, "\n")
			linelen = leaderlen + separatorlen + verblen
			j = 0
		}
		if j == 0 {
			fmt.Fprintf(o, "%s", leader)
		}
		fmt.Fprintf(o, "%s%s", separator, verb)
		j++
	}

	fmt.Fprintf(o, "\n")
}

// ----------------------------------------------------------------
func UsageAllVerbs(argv0 string) {
	separator := "================================================================"

	for _, transformerSetup := range TRANSFORMER_LOOKUP_TABLE {
		fmt.Printf("%s\n", separator)
		lib.InternalCodingErrorIf(transformerSetup.UsageFunc == nil)
		transformerSetup.UsageFunc(os.Stdout, false, 0)
	}
	fmt.Printf("%s\n", separator)
	os.Exit(0)
}
