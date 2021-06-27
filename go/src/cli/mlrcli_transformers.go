package cli

import (
	"fmt"
	"os"

	"miller/src/lib"
	"miller/src/transformers"
	"miller/src/transforming"
)

// ----------------------------------------------------------------
var MAPPER_LOOKUP_TABLE = []transforming.TransformerSetup{
	transformers.AltkvSetup,
	transformers.BarSetup,
	transformers.BootstrapSetup,
	transformers.CatSetup,
	transformers.CheckSetup,
	transformers.CleanWhitespaceSetup,
	transformers.CountDistinctSetup,
	transformers.CountSetup,
	transformers.CountSimilarSetup,
	transformers.CutSetup,
	transformers.DecimateSetup,
	transformers.FillDownSetup,
	transformers.FillEmptySetup,
	transformers.FilterSetup,
	transformers.FlattenSetup,
	transformers.FormatValuesSetup,
	transformers.FractionSetup,
	transformers.GapSetup,
	transformers.GrepSetup,
	transformers.GroupBySetup,
	transformers.GroupLikeSetup,
	transformers.HavingFieldsSetup,
	transformers.HeadSetup,
	transformers.HistogramSetup,
	transformers.JSONParseSetup,
	transformers.JSONStringifySetup,
	transformers.JoinSetup,
	transformers.LabelSetup,
	transformers.LeastFrequentSetup,
	transformers.MergeFieldsSetup,
	transformers.MostFrequentSetup,
	transformers.NestSetup,
	transformers.NothingSetup,
	transformers.PutSetup,
	transformers.RegularizeSetup,
	transformers.RemoveEmptyColumnsSetup,
	transformers.RenameSetup,
	transformers.ReorderSetup,
	transformers.RepeatSetup,
	transformers.ReshapeSetup,
	transformers.SampleSetup,
	transformers.Sec2GMTDateSetup,
	transformers.Sec2GMTSetup,
	transformers.SeqgenSetup,
	transformers.ShuffleSetup,
	transformers.SkipTrivialRecordsSetup,
	transformers.SortSetup,
	transformers.SortWithinRecordsSetup,
	transformers.Stats1Setup,
	transformers.Stats2Setup,
	transformers.StepSetup,
	transformers.TacSetup,
	transformers.TailSetup,
	transformers.TeeSetup,
	transformers.TopSetup,
	transformers.UnflattenSetup,
	transformers.UniqSetup,
	transformers.UnsparsifySetup,
}

// ----------------------------------------------------------------
func lookUpTransformerSetup(verb string) *transforming.TransformerSetup {
	for _, transformerSetup := range MAPPER_LOOKUP_TABLE {
		if transformerSetup.Verb == verb {
			return &transformerSetup
		}
	}
	return nil
}

func listAllVerbsRaw(o *os.File) {
	for _, transformerSetup := range MAPPER_LOOKUP_TABLE {
		fmt.Fprintf(o, "%s\n", transformerSetup.Verb)
	}
}

// TODO: move to help package
func listAllVerbs(o *os.File, leader string) {
	separator := " "

	leaderlen := len(leader)
	separatorlen := len(separator)
	linelen := leaderlen
	j := 0

	for _, transformerSetup := range MAPPER_LOOKUP_TABLE {
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

func usageAllVerbs(argv0 string) {
	separator := "================================================================"

	for _, transformerSetup := range MAPPER_LOOKUP_TABLE {
		fmt.Printf("%s\n", separator)
		lib.InternalCodingErrorIf(transformerSetup.UsageFunc == nil)
		transformerSetup.UsageFunc(os.Stdout, false, 0)
	}
	fmt.Printf("%s\n", separator)
	os.Exit(0)
}
