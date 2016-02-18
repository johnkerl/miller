#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "input/lrec_readers.h"
#include "containers/join_bucket_keeper.h"
#include "containers/mixutil.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

static int tjbk_verbose = FALSE;

// ----------------------------------------------------------------
static sllv_t* make_records_empty() {
	sllv_t* precords = sllv_alloc();
	return precords;
}

// ----------------------------------------------------------------
static sllv_t* make_records_113335() {
	sllv_t* precords = sllv_alloc();
	sllv_append(precords, lrec_literal_2("l","1", "b","10"));
	sllv_append(precords, lrec_literal_2("l","1", "b","11"));
	sllv_append(precords, lrec_literal_2("l","3", "b","12"));
	sllv_append(precords, lrec_literal_2("l","3", "b","13"));
	sllv_append(precords, lrec_literal_2("l","3", "b","14"));
	sllv_append(precords, lrec_literal_2("l","5", "b","15"));
	return precords;
}

// ----------------------------------------------------------------
static sllv_t* make_records_het() {
	sllv_t* precords = sllv_alloc();
	sllv_append(precords, lrec_literal_2("x","100", "b","10"));
	sllv_append(precords, lrec_literal_2("l","1",   "b","11"));
	sllv_append(precords, lrec_literal_2("l","1",   "b","12"));
	sllv_append(precords, lrec_literal_2("x","200", "b","13"));
	sllv_append(precords, lrec_literal_2("l","3",   "b","14"));
	sllv_append(precords, lrec_literal_2("l","3",   "b","15"));
	sllv_append(precords, lrec_literal_2("x","300", "b","16"));
	sllv_append(precords, lrec_literal_2("l","5",   "b","17"));
	sllv_append(precords, lrec_literal_2("l","5",   "b","18"));

	return precords;
}

// ----------------------------------------------------------------
static void set_up(
	sllv_t* precords,
	slls_t** ppleft_field_names,
	lrec_reader_t** ppreader)
{
	slls_t* pleft_field_names = slls_alloc();
	slls_append_no_free(pleft_field_names, "l");

	lrec_reader_t* preader = lrec_reader_in_memory_alloc(precords);
	printf("left records:\n");
	lrec_print_list_with_prefix(precords, "  ");
	printf("\n");

	*ppleft_field_names = pleft_field_names;
	*ppreader = preader;
}

static void emit(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** pprecords_paired, sllv_t** pprecords_left_unpaired)
{
	if (tjbk_verbose) {
		printf("BEFORE EMIT\n");
		join_bucket_keeper_print(pkeeper);
		printf("\n");
	}

	join_bucket_keeper_emit(pkeeper, pright_field_values, pprecords_paired, pprecords_left_unpaired);

	if (tjbk_verbose) {
		printf("AFTER EMIT\n");
		join_bucket_keeper_print(pkeeper);
		printf("\n");
	}
}

static int list_is_null(sllv_t* plist, char* desc, char* rval) {
	if (plist == NULL) {
		printf("%-8s is null with rval=\"%s\"; ok\n", desc, rval);
		return TRUE;
	} else {
		printf("%-8s should be null with rval=\"%s\" and is not (length=%llu):\n",
			desc, rval, plist->length);
		lrec_print_list_with_prefix(plist, "  ");
		return FALSE;
	}
}

static int list_has_length(sllv_t* plist, unsigned long long length, char* desc, char* rval) {
	if (plist == NULL) {
		printf("%-8s is null with rval=\"%s\" and should not be\n", desc, rval);
		return FALSE;
	}

	if (plist->length == length) {
		printf("%-8s length=%llu with rval=\"%s\"; ok:\n", desc, length, rval);
		lrec_print_list_with_prefix(plist, "  ");
		return TRUE;
	} else {
		printf("%-8s length=%llu with rval=\"%s\" but should be %llu:\n", desc, plist->length, rval, length);
		lrec_print_list_with_prefix(plist, "  ");
		return FALSE;
	}
}

// ----------------------------------------------------------------
static char* test_left_empty_right_empty() {
	printf("----------------------------------------------------------------\n");
	printf("test_left_empty_right_empty enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_empty(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, NULL, pleft_field_names);
	sllv_t* precords_paired;
	sllv_t* precords_left_unpaired;

	emit(pkeeper, NULL, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(precords_left_unpaired, 0, "unpaired", "(eof)"));
	printf("\n");

	printf("test_left_empty_right_empty exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test_left_empty() {
	printf("----------------------------------------------------------------\n");
	printf("test_left_empty enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_empty(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, NULL, pleft_field_names);
	sllv_t* precords_paired;
	sllv_t* precords_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("0");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(precords_left_unpaired, 0, "unpaired", "(eof)"));
	printf("\n");

	printf("test_left_empty exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test_right_empty() {
	printf("----------------------------------------------------------------\n");
	printf("test_right_empty enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_113335(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, NULL, pleft_field_names);
	sllv_t* precords_paired;
	sllv_t* precords_left_unpaired;

	emit(pkeeper, NULL, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(precords_left_unpaired, 6, "unpaired", "(eof)"));
	printf("\n");

	printf("test_right_empty exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test_unpaired_before_left_start() {
	printf("----------------------------------------------------------------\n");
	printf("test_unpaired_before_left_start enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_113335(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, NULL, pleft_field_names);
	sllv_t* precords_paired;
	sllv_t* precords_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("0");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(precords_left_unpaired, 6, "unpaired", "(eof)"));
	printf("\n");

	printf("test_unpaired_before_left_start exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test_unpaired_after_left_end() {
	printf("----------------------------------------------------------------\n");
	printf("test_unpaired_after_left_end enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_113335(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, NULL, pleft_field_names);
	sllv_t* precords_paired;
	sllv_t* precords_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("6");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_has_length(precords_left_unpaired, 6, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(precords_left_unpaired, 0, "unpaired", "(eof)"));
	printf("\n");

	printf("test_unpaired_after_left_end exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test_middle_pairings() {
	printf("----------------------------------------------------------------\n");
	printf("test_middle_pairings enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_113335(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, NULL, pleft_field_names);
	sllv_t* precords_paired;
	sllv_t* precords_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("0");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("1");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_has_length(precords_paired, 2, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("1");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_has_length(precords_paired, 2, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(precords_left_unpaired, 4, "unpaired", "(eof)"));
	printf("\n");

	printf("test_middle_pairings exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test_middle() {
	printf("----------------------------------------------------------------\n");
	printf("test_middle enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_113335(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, NULL, pleft_field_names);
	sllv_t* precords_paired;
	sllv_t* precords_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("2");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_has_length(precords_left_unpaired, 2, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("2");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("3");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_has_length(precords_paired, 3, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(precords_left_unpaired, 1, "unpaired", "(eof)"));
	printf("\n");

	printf("test_middle exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test_walk_through_all() {
	printf("----------------------------------------------------------------\n");
	printf("test_walk_through_all enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_113335(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, NULL, pleft_field_names);
	sllv_t* precords_paired;
	sllv_t* precords_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("0");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("0");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("1");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_has_length(precords_paired, 2, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("1");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_has_length(precords_paired, 2, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("2");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("2");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("3");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_has_length(precords_paired, 3, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("3");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_has_length(precords_paired, 3, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("4");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("4");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("5");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_has_length(precords_paired, 1, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("5");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_has_length(precords_paired, 1, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("6");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("6");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(precords_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(precords_left_unpaired, 0, "unpaired", pright_field_values->phead->value));
	printf("\n");

	printf("test_walk_through_all exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test_het_unpaired_before_left_start() {
	printf("----------------------------------------------------------------\n");
	printf("test_het_unpaired_before_left_start enter\n");
	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_het(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, NULL, pleft_field_names);
	sllv_t* precords_paired;
	sllv_t* precords_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("0");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_has_length(precords_left_unpaired, 2, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(precords_left_unpaired, 7, "unpaired", "(eof)"));
	printf("\n");
	printf("test_het_unpaired_before_left_start exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test_het_unpaired_after_left_end() {
	printf("----------------------------------------------------------------\n");
	printf("test_het_unpaired_after_left_end enter\n");
	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_het(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, NULL, pleft_field_names);
	sllv_t* precords_paired;
	sllv_t* precords_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("6");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_has_length(precords_left_unpaired, 9, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(precords_left_unpaired, 0, "unpaired", "(eof)"));
	printf("\n");
	printf("test_het_unpaired_after_left_end exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test_het_initial_pairing() {
	printf("----------------------------------------------------------------\n");
	printf("test_het_initial_pairing enter\n");
	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_het(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, NULL, pleft_field_names);
	sllv_t* precords_paired;
	sllv_t* precords_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("1");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_has_length(precords_paired, 2, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_has_length(precords_left_unpaired, 2, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(precords_left_unpaired, 5, "unpaired", "(eof)"));

	printf("\n");
	printf("test_het_initial_pairing exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test_het_middle_pairing() {
	printf("----------------------------------------------------------------\n");
	printf("test_het_middle_pairing enter\n");
	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_het(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, NULL, pleft_field_names);
	sllv_t* precords_paired;
	sllv_t* precords_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("3");
	emit(pkeeper, pright_field_values, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_has_length(precords_paired, 2, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_has_length(precords_left_unpaired, 5, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &precords_paired, &precords_left_unpaired);
	mu_assert_lf(list_is_null(precords_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(precords_left_unpaired, 2, "unpaired", "(eof)"));
	printf("\n");
	printf("test_het_middle_pairing exit\n");
	printf("\n");
	return 0;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_left_empty_right_empty);
	mu_run_test(test_left_empty);
	mu_run_test(test_right_empty);
	mu_run_test(test_unpaired_before_left_start);
	mu_run_test(test_unpaired_after_left_end);
	mu_run_test(test_middle_pairings);
	mu_run_test(test_middle);
	mu_run_test(test_walk_through_all);
	mu_run_test(test_het_unpaired_before_left_start);
	mu_run_test(test_het_unpaired_after_left_end);
	mu_run_test(test_het_initial_pairing);
	mu_run_test(test_het_middle_pairing);
	printf("----------------------------------------------------------------\n");
	return 0;
}

int main(int argc, char **argv) {
	printf("TEST_JOIN_BUCKET_KEEPER ENTER\n");
	for (int argi = 1; argi < argc; argi++) {
		if (streq(argv[argi], "-v"))
			tjbk_verbose = TRUE;
	}

	char *result = run_all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_JOIN_BUCKET_KEEPER: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
