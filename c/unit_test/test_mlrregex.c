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
static char * test_copy_regex_captures() {

	const size_t nmatchmax = 10; // Capture-groups \1 through \9 supported, along with entire-string match
	regmatch_t matches[nmatchmax];
	string_array_t* pregex_captures_1_up = string_array_alloc(0);
	regex_t regex;

	char* input  = "abcde";
	char* sregex = "abcde";
	regcomp_or_die(&regex, sregex, 0);
	regmatch_or_die(&regex, input, nmatchmax, matches);
	copy_regex_captures(pregex_captures_1_up, input, matches, nmatchmax);
	mu_assert_lf(pregex_captures_1_up->length == 1);
	mu_assert_lf(pregex_captures_1_up->strings[0] == NULL);
	regfree(&regex);

	input  = "abcde";
	sregex = "a(.*)e";
	regcomp_or_die(&regex, sregex, 0);
	regmatch_or_die(&regex, input, nmatchmax, matches);
	copy_regex_captures(pregex_captures_1_up, input, matches, nmatchmax);
	mu_assert_lf(pregex_captures_1_up->length == 2);
	mu_assert_lf(pregex_captures_1_up->strings[0] == NULL);
	mu_assert_lf(streq(pregex_captures_1_up->strings[1], "bcd"));
	regfree(&regex);

	input  = "abcde";
	sregex = "a(b)(.)(d)e";
	regcomp_or_die(&regex, sregex, 0);
	regmatch_or_die(&regex, input, nmatchmax, matches);
	copy_regex_captures(pregex_captures_1_up, input, matches, nmatchmax);
	mu_assert_lf(pregex_captures_1_up->length == 4);
	mu_assert_lf(pregex_captures_1_up->strings[0] == NULL);
	mu_assert_lf(streq(pregex_captures_1_up->strings[1], "b"));
	mu_assert_lf(streq(pregex_captures_1_up->strings[2], "c"));
	mu_assert_lf(streq(pregex_captures_1_up->strings[3], "d"));
	regfree(&regex);

	input  = "abcdefghij";
	sregex = "(a)(b)(c)(d)(e)(f)(g)(h)(i)";
	regcomp_or_die(&regex, sregex, 0);
	regmatch_or_die(&regex, input, nmatchmax, matches);
	copy_regex_captures(pregex_captures_1_up, input, matches, nmatchmax);
	printf("X=%d\n", pregex_captures_1_up->length);
	mu_assert_lf(pregex_captures_1_up->length == 10);
	mu_assert_lf(pregex_captures_1_up->strings[0] == NULL);
	mu_assert_lf(streq(pregex_captures_1_up->strings[1], "a"));
	mu_assert_lf(streq(pregex_captures_1_up->strings[2], "b"));
	mu_assert_lf(streq(pregex_captures_1_up->strings[3], "c"));
	mu_assert_lf(streq(pregex_captures_1_up->strings[4], "d"));
	mu_assert_lf(streq(pregex_captures_1_up->strings[5], "e"));
	mu_assert_lf(streq(pregex_captures_1_up->strings[6], "f"));
	mu_assert_lf(streq(pregex_captures_1_up->strings[7], "g"));
	mu_assert_lf(streq(pregex_captures_1_up->strings[8], "h"));
	mu_assert_lf(streq(pregex_captures_1_up->strings[9], "i"));
	regfree(&regex);

	string_array_free(pregex_captures_1_up);

	return 0;
}

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
	mu_assert_lf(streq(output, "hello"));
	mu_assert_lf(was_allocated == TRUE);
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
	mu_assert_lf(streq(output, "hcealblo"));
	mu_assert_lf(was_allocated == TRUE);
	string_array_free(psa);

	return 0;
}


// ================================================================
static char * all_tests() {
	mu_run_test(test_copy_regex_captures);
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
