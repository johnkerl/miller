package cli

import (
	"flag"
	"fmt"
	"os"

	"miller/transformers"
	"miller/transforming"
)

// ----------------------------------------------------------------
var MAPPER_LOOKUP_TABLE = []transforming.TransformerSetup{
	transformers.AltkvSetup,
	transformers.BootstrapSetup,
	transformers.CatSetup,
	transformers.CheckSetup,
	transformers.CountSetup,
	transformers.CountSimilarSetup,
	transformers.CutSetup,
	transformers.DecimateSetup,
	transformers.FillDownSetup,
	transformers.FilterSetup,
	transformers.FlattenSetup,
	transformers.GapSetup,
	transformers.GrepSetup,
	transformers.GroupBySetup,
	transformers.GroupLikeSetup,
	transformers.HeadSetup,
	transformers.JoinSetup,
	transformers.LabelSetup,
	transformers.NothingSetup,
	transformers.PutSetup,
	transformers.RegularizeSetup,
	transformers.RemoveEmptyColumnsSetup,
	transformers.RenameSetup,
	transformers.ReorderSetup,
	transformers.SampleSetup,
	transformers.Sec2GMTSetup,
	transformers.SeqgenSetup,
	transformers.ShuffleSetup,
	transformers.SkipTrivialRecordsSetup,
	transformers.SortSetup,
	transformers.SortWithinRecordsSetup,
	transformers.StepSetup,
	transformers.TacSetup,
	transformers.TailSetup,
	transformers.UnflattenSetup,
	transformers.UnsparsifySetup,
}

//	&transformer_altkv_setup,
//	&transformer_bar_setup,
//	&transformer_bootstrap_setup,
//	&transformer_cat_setup,
//	&transformer_check_setup,
//	&transformer_clean_whitespace_setup,
//	&transformer_count_setup,
//	&transformer_count_distinct_setup,
//	&transformer_count_similar_setup,
//	&transformer_cut_setup,
//	&transformer_decimate_setup,
//	&transformer_fill_down_setup,
//	&transformer_filter_setup,
//	&transformer_format_values_setup,
//	&transformer_fraction_setup,
//	&transformer_grep_setup,
//	&transformer_group_by_setup,
//	&transformer_group_like_setup,
//	&transformer_having_fields_setup,
//	&transformer_head_setup,
//	&transformer_histogram_setup,
//	&transformer_join_setup,
//	&transformer_label_setup,
//	&transformer_least_frequent_setup,
//	&transformer_merge_fields_setup,
//	&transformer_most_frequent_setup,
//	&transformer_nest_setup,
//	&transformer_nothing_setup,
//	&transformer_put_setup,
//	&transformer_regularize_setup,
//	&transformer_remove_empty_columns_setup,
//	&transformer_rename_setup,
//	&transformer_reorder_setup,
//	&transformer_repeat_setup,
//	&transformer_reshape_setup,
//	&transformer_sample_setup,
//	&transformer_sec2gmt_setup,
//	&transformer_sec2gmtdate_setup,
//	&transformer_seqgen_setup,
//	&transformer_shuffle_setup,
//	&transformer_skip_trivial_records_setup,
//	&transformer_sort_setup,
//	// xxx temp for 5.4.0 -- will continue work after
//	// &transformer_sort_within_records_setup,
//	&transformer_stats1_setup,
//	&transformer_stats2_setup,
//	&transformer_step_setup,
//	&transformer_tac_setup,
//	&transformer_tail_setup,
//	&transformer_tee_setup,
//	&transformer_top_setup,
//	&transformer_uniq_setup,
//	&transformer_unsparsify_setup,

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
		args := [3]string{os.Args[0], transformerSetup.Verb, "--help"}
		argi := 1
		transformerSetup.ParseCLIFunc(
			&argi,
			3,
			args[:],
			flag.ContinueOnError,
			nil,
			nil,
		)
		fmt.Printf("\n")
	}
	fmt.Printf("%s\n", separator)
	os.Exit(0)
}
