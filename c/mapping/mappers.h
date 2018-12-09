#ifndef MAPPERS_H
#define MAPPERS_H
#include "containers/lhmss.h"
#include "containers/slls.h"
#include "mapping/mapper.h"

extern mapper_setup_t mapper_altkv_setup;
extern mapper_setup_t mapper_bar_setup;
extern mapper_setup_t mapper_bootstrap_setup;
extern mapper_setup_t mapper_cat_setup;
extern mapper_setup_t mapper_check_setup;
extern mapper_setup_t mapper_clean_whitespace_setup;
extern mapper_setup_t mapper_count_distinct_setup;
extern mapper_setup_t mapper_count_similar_setup;
extern mapper_setup_t mapper_cut_setup;
extern mapper_setup_t mapper_decimate_setup;
extern mapper_setup_t mapper_fill_down_setup;
extern mapper_setup_t mapper_filter_setup;
extern mapper_setup_t mapper_fraction_setup;
extern mapper_setup_t mapper_grep_setup;
extern mapper_setup_t mapper_group_by_setup;
extern mapper_setup_t mapper_group_like_setup;
extern mapper_setup_t mapper_having_fields_setup;
extern mapper_setup_t mapper_head_setup;
extern mapper_setup_t mapper_histogram_setup;
extern mapper_setup_t mapper_join_setup;
extern mapper_setup_t mapper_label_setup;
extern mapper_setup_t mapper_least_frequent_setup;
extern mapper_setup_t mapper_merge_fields_setup;
extern mapper_setup_t mapper_most_frequent_setup;
extern mapper_setup_t mapper_nest_setup;
extern mapper_setup_t mapper_nothing_setup;
extern mapper_setup_t mapper_put_setup;
extern mapper_setup_t mapper_regularize_setup;
extern mapper_setup_t mapper_remove_empty_columns_setup;
extern mapper_setup_t mapper_rename_setup;
extern mapper_setup_t mapper_reorder_setup;
extern mapper_setup_t mapper_repeat_setup;
extern mapper_setup_t mapper_reshape_setup;
extern mapper_setup_t mapper_sample_setup;
extern mapper_setup_t mapper_sec2gmt_setup;
extern mapper_setup_t mapper_sec2gmtdate_setup;
extern mapper_setup_t mapper_seqgen_setup;
extern mapper_setup_t mapper_shuffle_setup;
extern mapper_setup_t mapper_skip_trivial_records_setup;
extern mapper_setup_t mapper_sort_setup;
extern mapper_setup_t mapper_stats1_setup;
extern mapper_setup_t mapper_stats2_setup;
extern mapper_setup_t mapper_step_setup;
extern mapper_setup_t mapper_tac_setup;
extern mapper_setup_t mapper_tail_setup;
extern mapper_setup_t mapper_tee_setup;
extern mapper_setup_t mapper_top_setup;
extern mapper_setup_t mapper_uniq_setup;
extern mapper_setup_t mapper_unsparsify_setup;

// Construction is in mlrcli.c.
void mapper_chain_free(sllv_t* pmapper_chain, context_t* pctx);

#endif // MAPPERS_H
