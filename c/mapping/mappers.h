#ifndef MAPPERS_H
#define MAPPERS_H
#include "containers/lhmss.h"
#include "containers/slls.h"
#include "mapping/mapper.h"

#define HAVING_FIELDS_AT_LEAST  0x4a
#define HAVING_FIELDS_WHICH_ARE 0x5b
#define HAVING_FIELDS_AT_MOST   0x6c

// Used by count-distinct
mapper_t* mapper_uniq_alloc(slls_t* pgroup_by_field_names, int show_counts);

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
extern mapper_setup_t mapper_label_setup;
extern mapper_setup_t mapper_put_setup;
extern mapper_setup_t mapper_rename_setup;
extern mapper_setup_t mapper_reorder_setup;
extern mapper_setup_t mapper_sort_setup;
extern mapper_setup_t mapper_stats1_setup;
extern mapper_setup_t mapper_stats2_setup;
extern mapper_setup_t mapper_step_setup;
extern mapper_setup_t mapper_tac_setup;
extern mapper_setup_t mapper_tail_setup;
extern mapper_setup_t mapper_top_setup;
extern mapper_setup_t mapper_uniq_setup;

#endif // MAPPERS_H
