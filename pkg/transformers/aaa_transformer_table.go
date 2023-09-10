package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/pkg/colorizer"
	"github.com/johnkerl/miller/pkg/lib"
)

// ----------------------------------------------------------------
var TRANSFORMER_LOOKUP_TABLE = []TransformerSetup{
	AltkvSetup,
	BarSetup,
	BootstrapSetup,
	CaseSetup,
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
	GsubSetup,
	HavingFieldsSetup,
	HeadSetup,
	HistogramSetup,
	JSONParseSetup,
	JSONStringifySetup,
	JoinSetup,
	LabelSetup,
	Latin1ToUTF8Setup,
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
	SplitSetup,
	SsubSetup,
	Stats1Setup,
	Stats2Setup,
	StepSetup,
	SubSetup,
	SummarySetup,
	TacSetup,
	TailSetup,
	TeeSetup,
	TemplateSetup,
	TopSetup,
	UTF8ToLatin1Setup,
	UnflattenSetup,
	UniqSetup,
	UnspaceSetup,
	UnsparsifySetup,
}

func ShowHelpForTransformer(verb string) bool {
	transformerSetup := LookUp(verb)
	if transformerSetup != nil {
		fmt.Println(colorizer.MaybeColorizeHelp(transformerSetup.Verb, true))
		transformerSetup.UsageFunc(os.Stdout)
		return true
	}
	return false
}

func ShowHelpForTransformerApproximate(searchString string) bool {
	found := false
	for _, transformerSetup := range TRANSFORMER_LOOKUP_TABLE {
		if strings.Contains(transformerSetup.Verb, searchString) {
			fmt.Println(colorizer.MaybeColorizeHelp(transformerSetup.Verb, true))
			transformerSetup.UsageFunc(os.Stdout)
			found = true
		}
	}
	return found
}

func LookUp(verb string) *TransformerSetup {
	for _, transformerSetup := range TRANSFORMER_LOOKUP_TABLE {
		if transformerSetup.Verb == verb {
			return &transformerSetup
		}
	}
	return nil
}

func ListVerbNamesVertically() {
	for _, transformerSetup := range TRANSFORMER_LOOKUP_TABLE {
		fmt.Printf("%s\n", transformerSetup.Verb)
	}
}

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
		transformerSetup.UsageFunc(os.Stdout)
	}
	fmt.Printf("%s\n", separator)
}
