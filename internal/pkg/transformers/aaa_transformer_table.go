package transformers

import (
	"fmt"
	"os"

	"mlr/internal/pkg/colorizer"
	"mlr/internal/pkg/lib"
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
	TemplateSetup,
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
func ListVerbNamesVertically() {
	for _, transformerSetup := range TRANSFORMER_LOOKUP_TABLE {
		fmt.Printf("%s\n", transformerSetup.Verb)
	}
}

// ----------------------------------------------------------------
func ListVerbNamesAsParagraph() {
	verbNames := make([]string, len(TRANSFORMER_LOOKUP_TABLE))

	for i, transformerSetup := range TRANSFORMER_LOOKUP_TABLE {
		verbNames[i] = transformerSetup.Verb
	}

	lib.PrintWordsAsParagraph(verbNames)
}

// ----------------------------------------------------------------
func UsageVerbs() {
	separator := "================================================================"

	for i, transformerSetup := range TRANSFORMER_LOOKUP_TABLE {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s\n", separator)
		lib.InternalCodingErrorIf(transformerSetup.UsageFunc == nil)
		fmt.Println(colorizer.MaybeColorizeHelp(transformerSetup.Verb, true))
		transformerSetup.UsageFunc(os.Stdout, false, 0)
	}
	fmt.Printf("%s\n", separator)
	os.Exit(0)
}
