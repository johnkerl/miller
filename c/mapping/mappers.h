#ifndef MAPPERS_H
#define MAPPERS_H
#include "containers/lhmss.h"
#include "containers/slls.h"
#include "mapping/mapper.h"

// xxx move to a different header file
#define HAVING_FIELDS_AT_LEAST     0x4a
#define HAVING_FIELDS_WHICH_ARE    0x5b
#define HAVING_FIELDS_AT_MOST      0x6c
#define HAVING_ALL_FIELDS_MATCHING 0x7d
#define HAVING_ANY_FIELDS_MATCHING 0x8e
#define HAVING_NO_FIELDS_MATCHING  0x9f

extern mapper_setup_t mapper_bar_setup;
extern mapper_setup_t mapper_cat_setup;
extern mapper_setup_t mapper_check_setup;
extern mapper_setup_t mapper_count_distinct_setup;
extern mapper_setup_t mapper_cut_setup;
extern mapper_setup_t mapper_filter_setup;
extern mapper_setup_t mapper_group_by_setup;
extern mapper_setup_t mapper_group_like_setup;
extern mapper_setup_t mapper_having_fields_setup;
extern mapper_setup_t mapper_head_setup;
extern mapper_setup_t mapper_histogram_setup;
extern mapper_setup_t mapper_join_setup;
extern mapper_setup_t mapper_label_setup;
extern mapper_setup_t mapper_put_setup;
extern mapper_setup_t mapper_regularize_setup;
extern mapper_setup_t mapper_rename_setup;
extern mapper_setup_t mapper_reorder_setup;
extern mapper_setup_t mapper_sample_setup;
extern mapper_setup_t mapper_sec2gmt_setup;
extern mapper_setup_t mapper_sort_setup;
extern mapper_setup_t mapper_stats1_setup;
extern mapper_setup_t mapper_stats2_setup;
extern mapper_setup_t mapper_step_setup;
extern mapper_setup_t mapper_tac_setup;
extern mapper_setup_t mapper_tail_setup;
extern mapper_setup_t mapper_top_setup;
extern mapper_setup_t mapper_uniq_setup;

#endif // MAPPERS_H
