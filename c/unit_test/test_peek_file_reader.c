#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include "lib/mlrutil.h"
#include "lib/minunit.h"
#include "lib/mlr_test_util.h"
#include "input/byte_readers.h"
#include "input/peek_file_reader.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char* test_empty() {
	byte_reader_t* pbr = string_byte_reader_alloc();
	int ok = pbr->popen_func(pbr, NULL, "");
	mu_assert_lf(ok == TRUE);

	peek_file_reader_t* pfr = pfr_alloc(pbr, 7);

	mu_assert_lf(pfr_peek_char(pfr) == (char)EOF); // char defaults to unsigned on some platforms
	mu_assert_lf(pfr_read_char(pfr) == (char)EOF);

	pbr->pclose_func(pbr);
	pfr_free(pfr);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_non_empty() {
	byte_reader_t* pbr = string_byte_reader_alloc();
	int ok = pbr->popen_func(pbr,

		NULL,

		"ab,cde\n"
		"123,4567\n"
	);
	mu_assert_lf(ok == TRUE);

	peek_file_reader_t* pfr = pfr_alloc(pbr, 7);

	pfr_print(pfr); mu_assert_lf(pfr_peek_char(pfr) == 'a');
	pfr_print(pfr); mu_assert_lf(pfr_read_char(pfr) == 'a');
	pfr_print(pfr); mu_assert_lf(pfr_peek_char(pfr) == 'b');
	pfr_print(pfr); mu_assert_lf(pfr_read_char(pfr) == 'b');

	pfr_print(pfr); mu_assert_lf(pfr_peek_char(pfr) == ',');
	pfr_print(pfr); mu_assert_lf(pfr_peek_char(pfr) == ',');
	pfr_print(pfr); mu_assert_lf(pfr_read_char(pfr) == ',');
	pfr_print(pfr); pfr_buffer_by(pfr, 5);
	pfr_print(pfr); pfr_advance_by(pfr, 5);
	pfr_print(pfr); mu_assert_lf(pfr_read_char(pfr) == '2');

	pfr_print(pfr); mu_assert_lf(pfr_peek_char(pfr) == '3');
	pfr_print(pfr); mu_assert_lf(pfr_peek_char(pfr) == '3');
	pfr_print(pfr); mu_assert_lf(pfr_read_char(pfr) == '3');
	pfr_print(pfr); pfr_buffer_by(pfr, 5);
	pfr_print(pfr); pfr_advance_by(pfr, 5);
	pfr_print(pfr); mu_assert_lf(pfr_read_char(pfr) == '\n');

	pbr->pclose_func(pbr);
	pfr_free(pfr);

	return NULL;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_empty);
	mu_run_test(test_non_empty);
	return 0;
}

int main(int argc, char **argv) {
	printf("TEST_PEEK_FILE_READER ENTER\n");
	char *result = run_all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_PEEK_FILE_READER: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
