package cli

import (
	"flag"
	"fmt"
	"os"

	"miller/mappers"
	"miller/mapping"
)

// ----------------------------------------------------------------
var MAPPER_LOOKUP_TABLE = []mapping.MapperSetup{
	mappers.CatSetup,
	mappers.CutSetup,
	mappers.GroupBySetup,
	mappers.GroupLikeSetup,
	mappers.HeadSetup,
	mappers.LabelSetup,
	mappers.NothingSetup,
	mappers.PutSetup,
	mappers.RenameSetup,
	mappers.TacSetup,
	mappers.TailSetup,
}

//	&mapper_altkv_setup,
//	&mapper_bar_setup,
//	&mapper_bootstrap_setup,
//	&mapper_cat_setup,
//	&mapper_check_setup,
//	&mapper_clean_whitespace_setup,
//	&mapper_count_setup,
//	&mapper_count_distinct_setup,
//	&mapper_count_similar_setup,
//	&mapper_cut_setup,
//	&mapper_decimate_setup,
//	&mapper_fill_down_setup,
//	&mapper_filter_setup,
//	&mapper_format_values_setup,
//	&mapper_fraction_setup,
//	&mapper_grep_setup,
//	&mapper_group_by_setup,
//	&mapper_group_like_setup,
//	&mapper_having_fields_setup,
//	&mapper_head_setup,
//	&mapper_histogram_setup,
//	&mapper_join_setup,
//	&mapper_label_setup,
//	&mapper_least_frequent_setup,
//	&mapper_merge_fields_setup,
//	&mapper_most_frequent_setup,
//	&mapper_nest_setup,
//	&mapper_nothing_setup,
//	&mapper_put_setup,
//	&mapper_regularize_setup,
//	&mapper_remove_empty_columns_setup,
//	&mapper_rename_setup,
//	&mapper_reorder_setup,
//	&mapper_repeat_setup,
//	&mapper_reshape_setup,
//	&mapper_sample_setup,
//	&mapper_sec2gmt_setup,
//	&mapper_sec2gmtdate_setup,
//	&mapper_seqgen_setup,
//	&mapper_shuffle_setup,
//	&mapper_skip_trivial_records_setup,
//	&mapper_sort_setup,
//	// xxx temp for 5.4.0 -- will continue work after
//	// &mapper_sort_within_records_setup,
//	&mapper_stats1_setup,
//	&mapper_stats2_setup,
//	&mapper_step_setup,
//	&mapper_tac_setup,
//	&mapper_tail_setup,
//	&mapper_tee_setup,
//	&mapper_top_setup,
//	&mapper_uniq_setup,
//	&mapper_unsparsify_setup,

// ----------------------------------------------------------------
func lookUpMapperSetup(verb string) *mapping.MapperSetup {
	for _, mapperSetup := range MAPPER_LOOKUP_TABLE {
		if mapperSetup.Verb == verb {
			return &mapperSetup
		}
	}
	return nil
}

func listAllVerbsRaw(o *os.File) {
	for _, mapperSetup := range MAPPER_LOOKUP_TABLE {
		fmt.Fprintf(o, "%s\n", mapperSetup.Verb)
	}
}

func listAllVerbs(o *os.File, leader string) {
	separator := " "

	leaderlen := len(leader)
	separatorlen := len(separator)
	linelen := leaderlen
	j := 0

	for _, mapperSetup := range MAPPER_LOOKUP_TABLE {
		verb := mapperSetup.Verb
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

	for _, mapperSetup := range MAPPER_LOOKUP_TABLE {
		fmt.Printf("%s\n", separator)
		args := [3]string{os.Args[0], mapperSetup.Verb, "--help"}
		argi := 1
		mapperSetup.ParseCLIFunc(
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
