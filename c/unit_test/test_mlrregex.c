#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char * test_interpolate_regex_captures() {
	int was_allocated = FALSE;

	char* output = interpolate_regex_captures("hello", NULL, &was_allocated);
	mu_assert_lf(streq(output, "hello"));
	mu_assert_lf(was_allocated == FALSE);

	string_array_t* psa = string_array_alloc(0);
	output = interpolate_regex_captures("hello", psa, &was_allocated);
	mu_assert_lf(streq(output, "hello"));
	mu_assert_lf(was_allocated == FALSE);
	string_array_free(psa);

	// captures are indexed 1-up so the X is a placeholder at index 0
	psa = string_array_from_line(mlr_strdup_or_die("X,a,b,c"), ',');
	output = interpolate_regex_captures("hello", psa, &was_allocated);
	mu_assert_lf(streq(output, "hello"));
	mu_assert_lf(was_allocated == FALSE);
	string_array_free(psa);

	psa = string_array_from_line(mlr_strdup_or_die("X,a,b,c"), ',');
	output = interpolate_regex_captures("h\\3ello", psa, &was_allocated);
	printf("output=[%s]\n", output);
	mu_assert_lf(streq(output, "hcello"));
	mu_assert_lf(was_allocated == TRUE);
	string_array_free(psa);

	psa = string_array_from_line(mlr_strdup_or_die("X,a,b,c"), ',');
	output = interpolate_regex_captures("h\\1ello", psa, &was_allocated);
	printf("output=[%s]\n", output);
	mu_assert_lf(streq(output, "haello"));
	mu_assert_lf(was_allocated == TRUE);
	string_array_free(psa);

	psa = string_array_from_line(mlr_strdup_or_die("X,a,b,c"), ',');
	output = interpolate_regex_captures("h\\4ello", psa, &was_allocated);
	printf("output=[%s]\n", output);
	mu_assert_lf(streq(output, "h\\4ello"));
	mu_assert_lf(was_allocated == FALSE);
	string_array_free(psa);

	psa = string_array_from_line(mlr_strdup_or_die("X,a,b,c"), ',');
	output = interpolate_regex_captures("h\\0ello", psa, &was_allocated);
	printf("output=[%s]\n", output);
	mu_assert_lf(streq(output, "h\\0ello"));
	mu_assert_lf(was_allocated == FALSE);
	string_array_free(psa);

	psa = string_array_from_line(mlr_strdup_or_die("X,a,b,c"), ',');
	output = interpolate_regex_captures("h\\3e\\1l\\2l\\4o", psa, &was_allocated);
	printf("output=[%s]\n", output);
	mu_assert_lf(streq(output, "hcealbl\\4o"));
	mu_assert_lf(was_allocated == TRUE);
	string_array_free(psa);

	return 0;
}


// ================================================================
static char * all_tests() {
	mu_run_test(test_interpolate_regex_captures);
	return 0;
}

int main(int argc, char **argv) {
	printf("TEST_MLRREGEX ENTER\n");
	char *result = all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_MLRREGEX: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
