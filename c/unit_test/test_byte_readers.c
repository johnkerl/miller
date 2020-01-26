#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include "lib/mlr_arch.h"
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "lib/minunit.h"
#include "lib/mlr_test_util.h"
#include "input/byte_readers.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char* test_string_byte_reader() {
	byte_reader_t* pbr = string_byte_reader_alloc();

	int ok = pbr->popen_func(pbr, NULL, "");
	mu_assert_lf(ok == TRUE);
	// char defaults to unsigned on some platforms -- but, byte_reader_t API is
	// in terms of ints.
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	pbr->pclose_func(pbr, NULL);

	ok = pbr->popen_func(pbr, NULL, "a");
	mu_assert_lf(ok == TRUE);
	mu_assert_lf(pbr->pread_func(pbr) == 'a');
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	pbr->pclose_func(pbr, NULL);

	ok = pbr->popen_func(pbr, NULL, "abc");
	mu_assert_lf(ok == TRUE);
	mu_assert_lf(pbr->pread_func(pbr) == 'a');
	mu_assert_lf(pbr->pread_func(pbr) == 'b');
	mu_assert_lf(pbr->pread_func(pbr) == 'c');
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	pbr->pclose_func(pbr, NULL);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_stdio_byte_reader_1() {
	byte_reader_t* pbr = stdio_byte_reader_alloc();

	char* contents = "";
	char* path = write_temp_file_or_die(contents);
	int ok = pbr->popen_func(pbr, NULL, path);
	mu_assert_lf(ok == TRUE);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	unlink_file_or_die(path);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_stdio_byte_reader_2() {
	byte_reader_t* pbr = stdio_byte_reader_alloc();

	char* contents = "abcdefg";
	char* path = write_temp_file_or_die(contents);
	int ok = pbr->popen_func(pbr, NULL, path);
	mu_assert_lf(ok == TRUE);
	mu_assert_lf(pbr->pread_func(pbr) == 'a');
	mu_assert_lf(pbr->pread_func(pbr) == 'b');
	mu_assert_lf(pbr->pread_func(pbr) == 'c');
	mu_assert_lf(pbr->pread_func(pbr) == 'd');
	mu_assert_lf(pbr->pread_func(pbr) == 'e');
	mu_assert_lf(pbr->pread_func(pbr) == 'f');
	mu_assert_lf(pbr->pread_func(pbr) == 'g');
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	unlink_file_or_die(path);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_stdio_byte_reader_reuse() {
	byte_reader_t* pbr = stdio_byte_reader_alloc();

	char* contents = "abc";
	char* path = write_temp_file_or_die(contents);
	int ok = pbr->popen_func(pbr, NULL, path);
	mu_assert_lf(ok == TRUE);
	mu_assert_lf(pbr->pread_func(pbr) == 'a');
	mu_assert_lf(pbr->pread_func(pbr) == 'b');
	mu_assert_lf(pbr->pread_func(pbr) == 'c');
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	unlink_file_or_die(path);

	contents = "defg";
	path = write_temp_file_or_die(contents);
	ok = pbr->popen_func(pbr, NULL, path);
	mu_assert_lf(ok == TRUE);
	mu_assert_lf(pbr->pread_func(pbr) == 'd');
	mu_assert_lf(pbr->pread_func(pbr) == 'e');
	mu_assert_lf(pbr->pread_func(pbr) == 'f');
	mu_assert_lf(pbr->pread_func(pbr) == 'g');
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	unlink_file_or_die(path);

	return NULL;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_string_byte_reader);
	mu_run_test(test_stdio_byte_reader_1);
	mu_run_test(test_stdio_byte_reader_2);
	mu_run_test(test_stdio_byte_reader_reuse);
	return 0;
}

int main(int argc, char **argv) {
	mlr_global_init(argv[0], NULL);
	printf("TEST_BYTE_READERS ENTER\n");
	char *result = run_all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_BYTE_READERS: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
